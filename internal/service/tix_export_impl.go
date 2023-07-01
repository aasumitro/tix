package service

import (
	"context"
	"fmt"
	"github.com/aasumitro/tix/common"
	"github.com/aasumitro/tix/internal/domain/entity"
	"github.com/aasumitro/tix/pkg/mailer"
	"github.com/aasumitro/tix/pkg/mailer/template"
	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
	"github.com/xuri/excelize/v2"
	"gopkg.in/gomail.v2"
	"os"
	"strconv"
	"strings"
	"time"
)

func (service *tixService) ExportEvent(
	ctx context.Context,
	googleFormID, exportFileType, targetEmail string,
) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	event, err := service.postgreSQLRepository.GetEventByGoogleFormID(ctx, googleFormID)
	if err != nil {
		return err
	}

	participants, err := service.postgreSQLRepository.GetAllParticipants(
		ctx, event.ID, "", 0, 0, 0, "name", "ASC")
	if err != nil {
		return err
	}

	totalApproved := service.postgreSQLRepository.CountParticipants(
		ctx, event.ID, common.ParticipantRequestApproved, 0, 0)
	totalDeclined := service.postgreSQLRepository.CountParticipants(
		ctx, event.ID, common.ParticipantRequestDeclined, 0, 0)
	totalWaitingApproval := service.postgreSQLRepository.CountParticipants(
		ctx, event.ID, common.ParticipantRequestWaiting, 0, 0)

	if strings.EqualFold(string(common.ExportTypeXLS), strings.ToLower(exportFileType)) {
		service.exportEventToExcel(
			event, participants, totalApproved,
			totalDeclined, totalWaitingApproval, targetEmail)
	}

	if strings.EqualFold(string(common.ExportTypePDF), strings.ToLower(exportFileType)) {
		service.exportEventToPDF(
			event, participants, totalApproved,
			totalDeclined, totalWaitingApproval, targetEmail)
	}

	return nil
}

func (service *tixService) exportEventToExcel(
	event *entity.Event,
	participants []*entity.Participant,
	totalApproved, totalDeclined, totalWaiting int,
	targetEmail string,
) {
	f := excelize.NewFile()
	defer func() { _ = f.Close() }()

	// HEADER
	_ = f.MergeCell("Sheet1", "A1", "I1")
	_ = f.MergeCell("Sheet1", "A2", "I2")
	_ = f.SetCellValue("Sheet1", "A1", "TIX")
	_ = f.SetCellValue("Sheet1", "A2", "Manage Events Participants and Tickets")
	titleStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
		},
		Font: &excelize.Font{
			Bold: true,
			Size: common.ExcelTitleSize,
		},
	})
	subtitleStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
		},
		Font: &excelize.Font{
			Bold:   false,
			Italic: true,
			Size:   common.ExcelSubtitleSize,
		},
	})
	_ = f.SetCellStyle("Sheet1", "A1", "A1", titleStyle)
	_ = f.SetCellStyle("Sheet1", "A2", "A2", subtitleStyle)
	_ = f.SetRowHeight("Sheet1", common.ExcelTitleRow, common.ExcelTitleHeight)
	_ = f.SetRowHeight("Sheet1", common.ExcelSubtitleRow, common.ExcelSubtitleHeight)

	// EVENT DATA
	_ = f.SetCellValue("Sheet1", "A5", "Nama:")
	_ = f.SetCellValue("Sheet1", "B5", event.Name)
	_ = f.SetCellValue("Sheet1", "A6", "Lokasi:")
	_ = f.SetCellValue("Sheet1", "B6", event.Location)
	_ = f.SetCellValue("Sheet1", "A7", "Tanggal:")
	_ = f.SetCellValue("Sheet1", "B7", func() string {
		ts := time.Unix(int64(event.EventDate), 0)
		return fmt.Sprintf("%d %s %d", ts.Day(), ts.Month().String(), ts.Year())
	}())
	_ = f.SetCellValue("Sheet1", "A8", "Total Peserta:")
	_ = f.SetCellValue("Sheet1", "B8", fmt.Sprintf(
		"%d –– %d diterima | %d menunggu | %d ditolak ––",
		event.TotalParticipants, totalApproved, totalWaiting, totalDeclined))
	_ = f.MergeCell("Sheet1", "B5", "I5")
	_ = f.MergeCell("Sheet1", "B6", "I6")
	_ = f.MergeCell("Sheet1", "B7", "I7")
	_ = f.MergeCell("Sheet1", "B8", "I8")

	// PARTICIPANT DATA

	rowsData := [][]interface{}{{
		"No", "Nama", "Email", "No Telp.", "Pekerjaan",
		"Tanggal Lahir", "Diterima", "Ditolak", "Alasan Ditolak"}}
	for i, participant := range participants {
		rowsData = append(rowsData, []interface{}{
			i + 1, participant.Name, participant.Email,
			participant.Phone, participant.Job, participant.DoB,
			func() string {
				if participant.ApprovedAt.Valid {
					return common.SymCheck
				}
				return common.SymDash
			}(),
			func() string {
				if participant.DeclinedAt.Valid {
					return common.SymCheck
				}
				return common.SymDash
			}(),
			func() string {
				if participant.DeclinedReason.Valid {
					return participant.DeclinedReason.String
				}
				return common.SymDash
			}(),
		})
	}
	borderStyle, _ := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{
				Type:  "left",
				Color: "000000",
				Style: 1,
			},
			{
				Type:  "top",
				Color: "000000",
				Style: 1,
			},
			{
				Type:  "right",
				Color: "000000",
				Style: 1,
			},
			{
				Type:  "bottom",
				Color: "000000",
				Style: 1,
			},
		},
	})
	for idx, row := range rowsData {
		rowData := row
		cellStart, _ := excelize.CoordinatesToCellName(1, idx+common.ExcelTableStartIndex)
		cellEnd, _ := excelize.CoordinatesToCellName(len(rowData), idx+common.ExcelTableStartIndex)
		_ = f.SetSheetRow("Sheet1", cellStart, &rowData)
		_ = f.SetCellStyle("Sheet1", cellStart, cellEnd, borderStyle)
		// Expand the width of cells in each row based on content
		for colIdx, cellValue := range rowData {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, idx+common.ExcelTableStartIndex)
			var valueStr string
			// Convert cell value to string
			switch val := cellValue.(type) {
			case string:
				valueStr = val
			case int:
				valueStr = strconv.Itoa(val)
			default:
				// Handle other data types if needed
				valueStr = ""
			}
			cellWidth := len(valueStr) + 2 // Add some padding
			// Get current column width
			currentWidth, _ := f.GetColWidth("Sheet1", cell[:1])
			if float64(cellWidth) > currentWidth {
				_ = f.SetColWidth("Sheet1", cell[:1], cell[:1], float64(cellWidth))
			}
		}
	}

	// Save spreadsheet by the given path.
	if err := f.SaveAs(fmt.Sprintf("./temps/exports/%s.xlsx", event.GoogleFormID)); err != nil {
		fmt.Println(err)
		return
	}

	service.sendViaEmail(event.Name, event.GoogleFormID, "xlsx", targetEmail)
}

func (service *tixService) exportEventToPDF(
	event *entity.Event,
	participants []*entity.Participant,
	totalApproved, totalDeclined, totalWaiting int,
	targetEmail string,
) {
	m := pdf.NewMaroto(consts.Landscape, consts.A4)
	m.SetPageMargins(common.PdfMarginLeft, common.PdfMarginTop, common.PdfMarginRight)
	m.RegisterHeader(func() {})
	m.RegisterFooter(func() {})

	// HEADER
	m.Row(common.PdfHeaderRowHeight, func() {
		m.Col(common.PdfHeaderColWidth, func() {
			m.Text("TIX", props.Text{
				Size:  common.PdfTitleSize,
				Style: consts.Bold,
				Align: consts.Center,
			})
			m.Text("Manage Events Participants and Tickets", props.Text{
				Top:   common.PdfSubtitleMarginTop,
				Size:  common.PdfSubtitleSize,
				Style: consts.Italic,
				Align: consts.Center,
			})
		})
	})
	m.Line(common.PdfLineSpaceHeight, props.Line{Width: common.PdfLineWidth})
	m.Row(common.PdfEventDataRowHeight, func() {
		m.Col(common.PdfEventDataTitleColWidth, func() {
			m.Text("Nama", props.Text{
				Size:  common.PdfEventDataSize,
				Style: consts.Bold,
				Align: consts.Left,
			})
		})
		m.Col(common.PdfEventDataItemColWidth, func() {
			m.Text(": "+event.Name, props.Text{
				Size:  common.PdfEventDataSize,
				Style: consts.Normal,
				Align: consts.Left,
			})
		})
	})
	m.Row(common.PdfEventDataRowHeight, func() {
		m.Col(common.PdfEventDataTitleColWidth, func() {
			m.Text("Lokasi", props.Text{
				Size:  common.PdfEventDataSize,
				Style: consts.Bold,
				Align: consts.Left,
			})
		})
		m.Col(common.PdfEventDataItemColWidth, func() {
			m.Text(": "+event.Location, props.Text{
				Size:  common.PdfEventDataSize,
				Style: consts.Normal,
				Align: consts.Left,
			})
		})
	})
	m.Row(common.PdfEventDataRowHeight, func() {
		m.Col(common.PdfEventDataTitleColWidth, func() {
			m.Text("Tanggal", props.Text{
				Size:  common.PdfEventDataSize,
				Style: consts.Bold,
				Align: consts.Left,
			})
		})
		m.Col(common.PdfEventDataItemColWidth, func() {
			m.Text(": "+func() string {
				ts := time.Unix(int64(event.EventDate), 0)
				return fmt.Sprintf("%d %s %d", ts.Day(), ts.Month().String(), ts.Year())
			}(), props.Text{
				Size:  common.PdfEventDataSize,
				Style: consts.Normal,
				Align: consts.Left,
			})
		})
	})
	m.Row(common.PdfEventDataRowHeight, func() {
		m.Col(common.PdfEventDataTitleColWidth, func() {
			m.Text("Total Peserta", props.Text{
				Size:  common.PdfEventDataSize,
				Style: consts.Bold,
				Align: consts.Left,
			})
		})
		m.Col(common.PdfEventDataItemColWidth, func() {
			m.Text(": "+fmt.Sprintf(
				"%d –– %d diterima | %d menunggu | %d ditolak ––",
				event.TotalParticipants, totalApproved, totalWaiting, totalDeclined,
			), props.Text{
				Size:  common.PdfEventDataSize,
				Style: consts.Normal,
				Align: consts.Left,
			})
		})
	})
	m.Line(common.PdfLineSpaceHeight, props.Line{Width: common.PdfLineWidth})

	// CONTENT
	m.Row(common.PdfTableTitleRowHeight, func() {
		m.Col(common.PdfTableTitleRowWidth, func() {
			m.Text("Daftar Peserta", props.Text{
				Size:  common.PdfTableTitleSize,
				Style: consts.Bold,
				Align: consts.Left,
				Top:   common.PdfTableTitleMarginTop,
			})
		})
	})
	tableHeader := []string{"Nama", "Tanggal Lahir", "Email", "No Telp.", "Pekerjaan", "Status"}
	var tableContents [][]string
	for _, participant := range participants {
		tableContents = append(tableContents, []string{
			participant.Name, participant.DoB, participant.Email,
			participant.Phone, participant.Job, func() string {
				if participant.ApprovedAt.Valid {
					return "diterima"
				}
				if participant.DeclinedAt.Valid {
					return "ditolak"
				}
				return "menunggu"
			}(),
		})
	}
	m.TableList(tableHeader, tableContents, props.TableList{
		ContentProp: props.TableListContent{
			GridSizes: []uint{2, 2, 2, 2, 2, 2},
			Size:      8,
		},
		HeaderProp: props.TableListContent{
			GridSizes: []uint{2, 2, 2, 2, 2, 2},
			Color:     color.Color{Red: 100},
		},
		Align: consts.Left,
		Line:  true,
		LineProp: props.Line{
			Color: color.Color{
				Red:   0,
				Green: 0,
				Blue:  0,
			},
			Style: consts.Solid,
		},
	})

	// FOOTER
	m.Row(common.PdfTableTitleRowHeight, func() {
		m.Col(common.PdfTableTitleRowWidth, func() {
			m.Text("Generated "+time.Now().String()+" via tix.bakode.xyz", props.Text{
				Style: consts.Italic,
				Size:  common.PdfFooterTitleSize,
				Align: consts.Center,
				Top:   common.PdfFooterTitleMarginTop,
			})
		})
	})

	if err := m.OutputFileAndClose(fmt.Sprintf("./temps/exports/%s.pdf", event.GoogleFormID)); err != nil {
		fmt.Println("⚠️  Could not save PDF:", err)
		return
	}

	service.sendViaEmail(event.Name, event.GoogleFormID, "pdf", targetEmail)
}

func (service *tixService) sendViaEmail(
	eventName, eventFormID, exportType, targetEmail string,
) {
	filePath := "temps/exports"
	attachmentName := fmt.Sprintf("%s.%s", eventFormID, exportType)
	title := fmt.Sprintf("Export Data for %s with type %s", eventName, exportType)

	// EMAIL FROM TEMPLATE
	m := mailer.Mailer{
		Theme: new(template.Default),
		Product: mailer.Product{
			Name: "TIX",
			Link: "https://tix.bakode.xyz/",
			Logo: "https://avatars.githubusercontent.com/u/105574217?s=400&u=81ba732eec2ca291da7654906168eb38a391ea22&v=4",
		},
	}
	e := mailer.Email{
		Body: mailer.Body{
			Name:   "Tix User",
			Intros: []string{"Please find attached the requested export of event data. Thank you for using tix app.!"},
		},
	}
	txtBody, err := m.GenerateHTML(&e)
	if err != nil {
		fmt.Println(err)
		return
	}

	// BUILD EMAIL
	mail := gomail.NewMessage()
	mail.SetHeader("From", "BAKODE SUPPORT <support@bakode.xyz>")
	mail.SetHeader("To", targetEmail)
	mail.SetHeader("Subject", title)
	mail.SetBody("text/html", txtBody)
	mail.Attach(fmt.Sprintf("%s/%s", filePath, attachmentName), gomail.Rename(attachmentName))

	// SEND EMAIL
	if err := service.mailer.DialAndSend(mail); err != nil {
		fmt.Println(err)
		return
	}

	// REMOVE FILE
	if err := os.Remove(fmt.Sprintf("%s/%s", filePath, attachmentName)); err != nil {
		fmt.Println("Error removing file:", err)
		return
	}
}
