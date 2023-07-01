package common

import "time"

const (
	ContextTimeout           = 25
	EventRemovalScheduleTime = 10
	EventSyncScheduleTime    = 30

	EventDataCacheTimeDuration        = 15 * time.Minute
	EventParticipantCacheTimeDuration = 3 * time.Minute
	GoogleFormCacheTimeDuration       = time.Hour * 24

	LastWeekDay = 7
)

const (
	SwaggerDefaultModelsExpandDepth = 4

	// DBMaxOpenConnection the default is 0 (unlimited)
	DBMaxOpenConnection = 100
	DBMaxIdleConnection = 10
	// DBMaxLifetimeConnection the default is 0 (connections are reused forever)
	DBMaxLifetimeConnection = 2

	SupabaseAuthEndpoint = "auth/v1"

	AccessTokenCookieKey = "access_token"

	EmptyPath = ""

	AutoSyncEventKey        = "event_auto_sync"
	ReqSyncEventQueueKey    = "req_sync_event_queue"
	ReqGenEventTixQueueKey  = "req_gen_event_tix_queue"
	ReqExpEventDataQueueKey = "req_exp_event_data_queue"
)

type EventParticipantStatus string

const (
	ParticipantStatusNone      EventParticipantStatus = "none"
	ParticipantRequestApproved EventParticipantStatus = "approved"
	ParticipantRequestDeclined EventParticipantStatus = "declined"
	ParticipantRequestWaiting  EventParticipantStatus = "waiting"
)

type EventExportType string

const (
	ExportTypePDF EventExportType = "pdf"
	ExportTypeXLS EventExportType = "xls"
)

const (
	MsgWaitGenTix = "Please wait a moment while we send the generated ticket to the intended recipient."
	MsgWaitSync   = "Please wait a moment while we sync the event data."
	MsgWaitExport = `
		Please wait a moment while we export the event data and send it to your email. 
		Please check your email periodically within 1-3 minutes after making these request.
    `
)

const (
	SymDash  = "–"
	SymCheck = "√"
)

const (
	ExcelTitleSize       = 21
	ExcelSubtitleSize    = 16
	ExcelTitleRow        = 1
	ExcelSubtitleRow     = 2
	ExcelTitleHeight     = 40
	ExcelSubtitleHeight  = 30
	ExcelTableStartIndex = 11
)

const (
	PdfMarginLeft        = 20
	PdfMarginRight       = 20
	PdfMarginTop         = 10
	PdfHeaderRowHeight   = 20
	PdfHeaderColWidth    = 12
	PdfTitleSize         = 18
	PdfSubtitleMarginTop = 9
	PdfSubtitleSize      = 14

	PdfLineSpaceHeight = 1.0
	PdfLineWidth       = 0.5

	PdfEventDataRowHeight     = 4
	PdfEventDataTitleColWidth = 3
	PdfEventDataItemColWidth  = 8
	PdfEventDataSize          = 8

	PdfTableTitleRowHeight = 20
	PdfTableTitleRowWidth  = 12
	PdfTableTitleSize      = 12
	PdfTableTitleMarginTop = 8

	PdfFooterTitleSize      = 8
	PdfFooterTitleMarginTop = 12
)
