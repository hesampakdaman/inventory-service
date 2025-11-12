package models

import "github.com/google/uuid"

type RequestID uuid.UUID

func (r RequestID) String() string {
	return uuid.UUID(r).String()
}
