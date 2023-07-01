package service

import (
	"context"
	"fmt"
	"github.com/aasumitro/tix/common"
	"github.com/aasumitro/tix/internal/domain/entity"
	"github.com/aasumitro/tix/pkg/mailer"
	"github.com/aasumitro/tix/pkg/mailer/template"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
	"gopkg.in/gomail.v2"
	"os"
)

func (service *tixService) GenerateTicket(
	ctx context.Context,
	googleFormID string,
	participantID int32,
) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	event, err := service.postgreSQLRepository.GetEventByGoogleFormID(
		ctx, googleFormID)
	if err != nil {
		return err
	}

	participant, err := service.postgreSQLRepository.GetParticipantByIDAndEventID(
		ctx, participantID, event.ID)
	if err != nil {
		return err
	}

	if err := service.generatePDFTicket(event, participant); err != nil {
		return err
	}

	service.sendTicketViaEmail(event.ID, participant.ID, event.Name, participant.Name, participant.Email)

	return nil
}

func (service *tixService) generatePDFTicket(
	event *entity.Event,
	participant *entity.Participant,
) error {
	m := pdf.NewMaroto(consts.Landscape, consts.A4)
	m.SetPageMargins(common.PdfMarginLeft, common.PdfMarginTop, common.PdfMarginRight)
	m.RegisterHeader(func() {})
	m.RegisterFooter(func() {})

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

	attachment := fmt.Sprintf("./temps/exports/gen%d%dtix.pdf",
		event.ID, participant.ID)
	if err := m.OutputFileAndClose(attachment); err != nil {
		return fmt.Errorf("⚠️ could not save pdf: %s", err.Error())
	}

	return nil
}

func (service *tixService) sendTicketViaEmail(
	eventID, participantID int32,
	eventName, participantName, targetEmail string,
) {
	filePath := "temps/exports"
	attachmentName := fmt.Sprintf("gen%d%dtix.pdf", eventID, participantID)
	title := fmt.Sprintf("Ticket for %s", eventName)

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
			Name:   participantName,
			Intros: []string{"Please find attached the requested ticket of event!"},
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
