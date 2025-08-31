package types

import "errors"

// Common errors
var (
	ErrInvalidWinstonValue = errors.New("invalid winston value")
	ErrInvalidSigner       = errors.New("invalid signer")
	ErrInvalidConfig       = errors.New("invalid configuration")
	ErrUploadFailed        = errors.New("upload failed")
	ErrSigningFailed       = errors.New("signing failed")
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrInvalidDataItem     = errors.New("invalid data item")
	ErrNetworkError        = errors.New("network error")
	ErrAuthenticationFailed = errors.New("authentication failed")
	ErrInvalidResponse     = errors.New("invalid response from server")
)
