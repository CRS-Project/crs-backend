package dto

type (
	CommentRequest struct {
		ID                string `json:"-"`
		Section           string `json:"section" binding:""`
		Comment           string `json:"comment" binding:"required"`
		Baseline          string `json:"baseline" binding:""`
		DocumentId        string `json:"document_id" binding:"required"`
		IsCloseOutComment bool   `json:"is_close_out_comment" binding:""`
		AreaOfConcernId   string `json:"-"`
		UserId            string `json:"-"`
		ReplyId           string `json:"-"`
	}

	UpdateCommentRequest struct {
		ID                string  `json:"-"`
		Section           string  `json:"section" binding:""`
		Comment           string  `json:"comment" binding:"required"`
		Baseline          string  `json:"baseline" binding:""`
		Status            *string `json:"status"  binding:""`
		DocumentId        string  `json:"document_id" binding:"required"`
		IsCloseOutComment bool    `json:"is_close_out_comment" binding:""`
		AreaOfConcernId   string  `json:"-"`
		UserId            string  `json:"-"`
		ReplyId           string  `json:"-"`
	}

	CommentResponse struct {
		ID                    string            `json:"id"`
		Section               string            `json:"section"`
		Comment               string            `json:"comment"`
		Baseline              string            `json:"baseline"`
		Status                *string           `json:"status"`
		DocumentID            string            `json:"document_id"`
		CommentAt             string            `json:"comment_at"`
		CompanyDocumentNumber string            `json:"company_document_number"`
		UserComment           *UserComment      `json:"user_comment,omitempty"`
		CommentReplies        []CommentResponse `json:"comment_replies"`
	}
)
