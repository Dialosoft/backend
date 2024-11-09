package request

type NewPost struct {
	UserID  string `json:"userID" example:"8213280e-2000-403a-b375-cdcda6488450"`
	ForumID string `json:"forumID" example:"8213280e-2000-403a-b375-cdcda6488451"`
	Title   string `json:"title" example:"some title"`
	Content string `json:"content" example:"some content"`
}

type UpdatePostTitle struct {
	Title  string `json:"title"`
	PostID string `json:"postID"`
}

type UpdatePostContent struct {
	Content string `json:"content"`
	PostID  string `json:"postID"`
}

type LikeOrUnlikePost struct {
	PostID string `json:"postID"`
	UserID string `json:"userID"`
}
