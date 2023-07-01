package request

type (
	AuthRequestMakeMagicLink struct {
		Email string `json:"email" form:"email" binding:"required"`
	}

	AuthRequestMakeSession struct {
		JWT  string `json:"jwt" form:"jwt" binding:"required"`
		Type string `json:"type" form:"jwt" binding:"required"`
	}

	AuthRequestInvite struct {
		Email string `json:"email" form:"email" binding:"required"`
	}

	EventRequestMakeNew struct {
		GoogleFormID    string `json:"google_form_id" form:"google_form_id" binding:"required"`
		Name            string `json:"name" form:"name" binding:"required"`
		PreregisterDate string `json:"preregister_date" form:"preregister_date" binding:"required"`
		EventDate       string `json:"event_date" form:"event_date" binding:"required"`
		Location        string `json:"location" form:"location" binding:"required"`
	}

	EventRequestUpdateParticipant struct {
		Status         string `json:"status" form:"status" binding:"required"`
		DeclinedReason string `json:"declined_reason,omitempty" form:"declined_reason,omitempty"`
	}

	EventValidationRequest struct {
		GoogleFormID string `json:"google_form_id" form:"google_form_id" binding:"required"`
	}
)
