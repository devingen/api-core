package dto

import "time"

type UpdateEntryResponse struct {
	ID        string    `json:"_id"`
	UpdatedAt time.Time `json:"_updated"`
	Revision  int       `json:"_revision"`
}
