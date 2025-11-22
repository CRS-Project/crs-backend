package dto

type (
	StatisticAOCAndCommentChart struct {
		Name                 string `json:"name"`
		TotalAreaOfConcern   int    `json:"total_area_of_concern"`
		TotalDocuments       int    `json:"total_documents"`
		TotalComments        int    `json:"total_comments"`
		TotalCommentRejected int    `json:"total_comment_rejected"`
	}

	StatisticAOCAndCommentCard struct {
		TotalAreaOfConcern           int `json:"total_area_of_concern"`
		TotalDocuments               int `json:"total_documents"`
		TotalComments                int `json:"total_comments"`
		TotalCommentRejected         int `json:"total_comment_rejected"`
		TotalDocumentsWithoutComment int `json:"total_documents_without_comment"`
	}

	StatisticCommentUsersChart struct {
		Name          string `json:"name"`
		CommentClosed int    `json:"comment_closed"`
		TotalComment  int    `json:"total_comment"`
	}

	StatisticCommentUsersData struct {
		ID            string `json:"id"`
		Name          string `json:"name"`
		CommentClosed int    `json:"comment_closed"`
		TotalComment  int    `json:"total_comment"`
	}
)
