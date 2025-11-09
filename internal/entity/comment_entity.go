package entity

import "github.com/google/uuid"

type CommentStatus string

const (
	CommentStatusAccepted CommentStatus = "ACCEPTED"
	CommentStatusReject   CommentStatus = "REJECT"
)

type Comment struct {
	ID uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`

	Section           string         `json:"section" gorm:"not null"`
	Comment           string         `json:"comment" gorm:"not null"`
	Baseline          string         `json:"baseline" gorm:"not null"`
	IsCloseOutComment bool           `json:"is_close_out_comment" gorm:""`
	Status            *CommentStatus `json:"comment_status" gorm:""`

	AreaOfConcernID uuid.UUID  `json:"area_of_concern_id" gorm:"not null"`
	DocumentID      uuid.UUID  `json:"document_id" gorm:"not null"`
	UserID          uuid.UUID  `json:"user_id" gorm:"not null"`
	CommentReplyID  *uuid.UUID `json:"comment_reply_id" gorm:""`

	Timestamp

	AreaOfConcern *AreaOfConcern `json:"area_of_concern,omitempty" gorm:"foreignKey:AreaOfConcernID"`
	User          *User          `json:"user" gorm:"foreignKey:UserID"`
	Document      *Document      `json:"document,omitempty" gorm:"foreignKey:DocumentID"`
	CommentReply  *Comment       `json:"comment_reply,omitempty" gorm:"foreignKey:CommentReplyID"`
}
