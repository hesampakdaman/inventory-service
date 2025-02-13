package models

import "github.com/google/uuid"

type ProductID uuid.UUID

type Product struct {
	ID          ProductID
	Available   int
	Reserved    int
	Title       string
	Description string
}
