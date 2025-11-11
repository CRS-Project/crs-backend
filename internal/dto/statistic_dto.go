package dto

type (
	StatisticAOCAndCommentCharts struct {
		Date                 string `json:"date"`
		TotalAreaOfConcern   int    `json:"total_area_of_concern"`
		TotalDocuments       int    `json:"total_documents"`
		TotalComments        int    `json:"total_comments"`
		TotalCommentRejected int    `json:"total_comment_rejected"`
	}

	StatisticAOCAndCommentCard struct {
		TotalAreaOfConcern   int `json:"total_area_of_concern"`
		TotalDocuments       int `json:"total_documents"`
		TotalComments        int `json:"total_comments"`
		TotalCommentRejected int `json:"total_comment_rejected"`
	}

	StatisticCommentUsersCharts struct {
		SMEInitial    int `json:"sme_initial"`
		CommentClosed int `json:"comment_closed"`
		TotalComment  int `json:"total_comment"`
	}
)
