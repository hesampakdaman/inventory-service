package models

import "errors"

var (
	ErrInsufficientStock       = errors.New("insufficient stock")
	ErrInvalidReservationState = errors.New("reservation not pending")
	ErrReservationNotFound     = errors.New("Reservation not found")
	ErrProductNotFound         = errors.New("Product not found")
)
