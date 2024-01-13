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
	AuthorId int64 `json:"authorId"`
	Offset   int64 `json:"offset"`
	Limit    int64 `json:"limit"`
}

type ArticleListReply struct {
	Articles []*Article `json:"articles"`
}

type GetArticleReply struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type LikeArticleReq struct {
	ArticleId int64 `json:"articleId"`
	// 1 -> 从不喜欢改为喜欢  0 -> 从喜欢改为不喜欢
	Like int64 `json:"like"`
}

type LikeArticleReply struct {
	OK bool `json:"ok"`
	// "like" or "normal" or "dislike"
	// 喜欢 / 默认 / 踩
	CurrentStatus string `json:"currentStatus"`
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
