// Package models contains business models description.
package entity

import "time"

// ShortURL is main entity for system.
type ShortURL struct {
	OriginalURL   string    `json:"url"`            // original URL that was shortened
	ID            string    `json:"id"`             // unique ID of the short URL.
	DeletedAt     time.Time `json:"deleted_at"`     // is used to mark a record as deleted
	CreatedByID   string    `json:"created_by"`     // ID of the user who created the short URL
	CorrelationID string    `json:"correlation_id"` // CorrelationID is used for matching original and shorten urls in shorten batch operation
}

// исходный URL, который был сокращён
// уникальный идентификатор короткого URL.
// используется для отметки записи как удалённой
// Идентификатор пользователя, создавшего короткий URL
// CorrelationID используется для сопоставления исходного и сокращённого URL при пакетном сокращении
