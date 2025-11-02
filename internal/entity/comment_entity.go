package entity

import "github.com/google/uuid"

type Comment struct {
	ID uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`

	Section  string `json:"section" gorm:"not null"`
	Comment  string `json:"comment" gorm:"not null"`
	Baseline string `json:"baseline" gorm:"not null"`

	DocumentID     uuid.UUID `json:"document_id" gorm:"not null"`
	UserID         uuid.UUID `json:"user_id" gorm:"not null"`
	CommentReplyID uuid.UUID `json:"comment_reply_id" gorm:"not null"`

	Timestamp

	Document *Document `json:"document,omitempty" gorm:"foreignKey:DocumentID"`
}
