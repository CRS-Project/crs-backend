package dto

type (
	CommentRequest struct {
		ID              string `json:"-"`
		Section         string `json:"section" binding:"required"`
		Comment         string `json:"comment" binding:"required"`
		Baseline        string `json:"baseline" binding:"required"`
		DocumentId      string `json:"document_id" binding:"required"`
		AreaOfConcernId string `json:"-"`
		UserId          string `json:"-"`
		ReplyId         string `json:"-"`
	}

	UpdateCommentRequest struct {
		ID              string  `json:"-"`
		Section         string  `json:"section" binding:"required"`
		Comment         string  `json:"comment" binding:"required"`
		Baseline        string  `json:"baseline" binding:"required"`
		Status          *string `json:"status"  binding:""`
		DocumentId      string  `json:"document_id" binding:"required"`
		AreaOfConcernId string  `json:"-"`
		UserId          string  `json:"-"`
		ReplyId         string  `json:"-"`
	}

	CommentResponse struct {
		ID          string       `json:"id"`
		Section     string       `json:"section"`
		Comment     string       `json:"comment"`
		Baseline    string       `json:"baseline"`
		Status      *string      `json:"status"`
		CommentAt   string       `json:"comment_at"`
		UserComment *UserComment `json:"user_comment,omitempty"`
	}
)
