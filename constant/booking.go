package constant

import "errors"

var (
	ErrSoldOut          = errors.New("event is sold out")
	ErrBookingNotFound  = errors.New("booking not found")
	ErrAlreadyCancelled = errors.New("booking already cancelled")
	ErrNotPending       = errors.New("booking is not pending")
)

const (
	BOOKING_STATUS_PENDING   = 1
	BOOKING_STATUS_CONFIRMED = 2
	BOOKING_STATUS_CANCELLED = 3
)

const (
	PAYMENT_STATUS_PENDING = 1
	PAYMENT_STATUS_SUCCESS = 2
	PAYMENT_STATUS_FAILED  = 3
)

const SYSTEM_USER_ID uint64 = 1000
