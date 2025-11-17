package models

import "errors"

var (
	ErrInsufficientStock       = errors.New("insufficient stock")
	ErrInvalidReservationState = errors.New("reservation not pending")
	ErrReservationNotFound     = errors.New("reservation not found")
	ErrProductNotFound         = errors.New("product not found")
)
