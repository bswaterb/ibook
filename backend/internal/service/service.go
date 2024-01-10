package service

import "github.com/google/wire"

// ServiceProviderSet is data providers.
var ServiceProviderSet = wire.NewSet(NewUserService, NewArticleService)

type User struct {
	Id          int64
	Email       string
	NickName    string
	PhoneNumber string
	PassWord    string
}

type ArticleAuthor struct {
	Article
}

type ArticleReader struct {
	Article
}

type ArticleStatus uint8

const (
	ArticleStatusUnknown = iota
	ArticleStatusUnpublished
	ArticleStatusPublished
	ArticleStatusPrivate
)

type Article struct {
	Id          int64
	Title       string
	Abstract    string
	Content     string
	Status      ArticleStatus
	Author      Author
	UpdatedTime int64
	CreatedTime int64
}

type Author struct {
	Id   int64
	Name string
}

func (a *Article) GenAbstract() string {
	contentString := []rune(a.Content)
	if len(contentString) < 100 {
		a.Abstract = a.Content
	}
	return string(contentString[:100])
}
