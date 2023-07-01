package response

type (
	AutoSyncRespond struct {
		FormID    string `json:"form_id"`
		EventDate int32  `json:"event_date"`
	}

	ServiceSingleRespond struct {
		Code    int
		Message any
	}

	GoogleFormQuestion struct {
		ID    string `json:"id"`
		Title string `json:"title"`
	}

	GoogleFormRespond struct {
		RespondID         string                   `json:"respond_id"`
		LastSubmittedTime string                   `json:"last_submitted_time"`
		CreateTime        string                   `json:"create_time"`
		Answer            *GoogleFormRespondAnswer `json:"answers"`
	}

	GoogleFormRespondAnswer struct {
		PoP   string `json:"pop"` // Proof of payment - bukti transfer
		DoB   string `json:"dob"` // Date of birth
		Email string `json:"email"`
		Name  string `json:"name"`
		Phone string `json:"phone"`
		Job   string `json:"job"`
	}

	SupabaseRespond struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	EventResponse struct {
		ID                int32  `json:"id"`
		GoogleFormID      string `json:"google_form_id"`
		Name              string `json:"name"`
		Location          string `json:"location"`
		PreregisterDate   int32  `json:"preregister_date"`
		EventDate         int32  `json:"event_date"`
		TotalParticipants int32  `json:"total_participants"`
		IsActive          bool   `json:"is_active"`
	}

	ParticipantResponse struct {
		ID             int32  `json:"id"`
		EventID        int32  `json:"event_id"`
		Name           string `json:"name"`
		Email          string `json:"email"`
		Phone          string `json:"phone"`
		Job            string `json:"job"`
		PoP            string `json:"prof_of_payment"`
		DoB            string `json:"date_of_birth"`
		ApprovedAt     *int32 `json:"approved_at"`
		DeclinedAt     *int32 `json:"declined_at"`
		DeclinedReason string `json:"declined_reason"`
		Status         string `json:"status"`
	}

	WeeklyOverviewResponse struct {
		Name  string `json:"name"`
		Total int    `json:"total"`
	}

	EventOverviewResponse struct {
		*EventResponse
		TotalApprovedParticipant        int                       `json:"total_approved_participant"`
		TotalWaitingApprovalParticipant int                       `json:"total_waiting_approval_participant"`
		TotalDeclinedParticipant        int                       `json:"total_declined_participant"`
		WeeklyOverview                  []*WeeklyOverviewResponse `json:"weekly_overview"`
		LatestRespondents               []*ParticipantResponse    `json:"latest_respondents"`
	}

	UserResponse struct {
		ID         int32  `json:"id"`
		UUID       string `json:"uuid"`
		Username   string `json:"username"`
		Email      string `json:"email"`
		IsVerified bool   `json:"is_verified"`
	}
)
