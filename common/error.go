package common

import "errors"

var (
	ErrUserNotFound            = errors.New("account with given email is not found. Please contact the administrator to invite you as a user to continue using this application")
	ErrUserAlreadyInvited      = errors.New("user with the given email has already been invited. Please give instruction to check their email to continue using this application")
	ErrRateLimitingPushQueue   = errors.New("you can make this request once every minute")
	ErrDeclineReasonNotProvide = errors.New("please provide decline status")
)
