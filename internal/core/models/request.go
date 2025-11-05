package models

import "github.com/google/uuid"

type RequstID uuid.UUID

func (r RequstID) String() string {
	return uuid.UUID(r).String()
}
