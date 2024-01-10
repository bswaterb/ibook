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

type ArticlePublishReq struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type ArticlePublishReply struct {
	Id int64 `json:"id"`
	OK bool  `json:"ok"`
}

type ArticleWithdrawReq struct {
	Id int64 `json:"id"`
}

type ArticleWithdrawReply struct {
	Id int64 `json:"id"`
	OK bool  `json:"ok"`
}

type ArticleListReq struct {
	Offset int64 `json:"offset"`
	Limit  int64 `json:"limit"`
}

type ArticleListReply struct {
	Articles []Article `json:"articles"`
}

type GetArticleReply struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Article struct {
	Id       int64  `json:"id"`
	Title    string `json:"title"`
	Abstract string `json:"abstract"`
	// Content string `json:"content"`
	Status      uint8  `json:"status"`
	AuthorId    int64  `json:"authorId"`
	AuthorName  string `json:"authorName"`
	CreatedTime string `json:"createdTime"`
	UpdatedTime string `json:"updatedTime"`
}

func (req ArticleEditReq) validate() bool {
	return true
}
