package web

type ArticleEditReq struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type ArticleEditReply struct {
	Id       int64  `json:"id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	AuthorId int64  `json:"authorId"`
}

func (req ArticleEditReq) validate() bool {
	return true
}
