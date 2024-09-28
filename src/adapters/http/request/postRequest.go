package request

type NewPost struct {
	UserID   string `json:"userID"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Views    uint32 `json:"views"`
	Comments uint32 `json:"comments"`
}
