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
	IsCloseOutComment bool           `json:"is_close_out_comment" gorm:"default:false"`
	AttachFileUrl     *string        `json:"attach_file_url" gorm:""`
	Status            *CommentStatus `json:"comment_status" gorm:""`

	DisciplineListDocumentID uuid.UUID  `json:"discipline_list_document_id" gorm:"not null"`
	UserID                   uuid.UUID  `json:"user_id" gorm:"not null"`
	CommentReplyID           *uuid.UUID `json:"comment_reply_id" gorm:""`

	DeletedBy uuid.UUID `json:"deleted_by" gorm:"type:uuid"`
	UpdatedBy uuid.UUID `json:"updated_by" gorm:"type:uuid"`
	Timestamp

	DisciplineListDocument *DisciplineListDocument `json:"discipline_list_document,omitempty" gorm:"foreignKey:DisciplineListDocumentID"`
	User                   *User                   `json:"user" gorm:"foreignKey:UserID"`
	CommentReply           *Comment                `json:"comment_reply,omitempty" gorm:"foreignKey:CommentReplyID"`
	CommentReplies         []Comment               `json:"comment_replies,omitempty" gorm:"foreignKey:CommentReplyID"`
}
