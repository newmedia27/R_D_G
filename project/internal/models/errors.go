package models

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserDuplicate      = errors.New("user already exists")
	ErrSessionNotFound    = errors.New("session not found")
	ErrSessionDuplicate   = errors.New("session already exists")
	ErrInvalidPassword    = errors.New("invalid password")
	ErrInvalidDevice      = errors.New("invalid device")
	ErrSessionExpired     = errors.New("session expired")
	ErrSuspiciousActivity = errors.New("suspicious activity")
	ErrChatDuplicate      = errors.New("chat already exists")
	ErrChatNotFound       = errors.New("chat not found")
	ErrInvalidUserIds     = errors.New("invalid user ids")
	ErrMessageDuplicate   = errors.New("message already exists")
	ErrMessageNotFound    = errors.New("message not found")
	ErrForbidden          = errors.New("forbidden")
	ErrIncorrectId        = errors.New("incorrect id")
	ErrEmptyId            = errors.New("empty id")
)
