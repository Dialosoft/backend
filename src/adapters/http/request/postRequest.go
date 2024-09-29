package request

type NewPost struct {
	UserID  string `json:"userID"`
	Title   string `json:"title"`
	Content string `json:"content"`
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
