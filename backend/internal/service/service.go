package service

import "github.com/google/wire"

// ServiceProviderSet is data providers.
var ServiceProviderSet = wire.NewSet(NewUserService, NewArticleService)

type User struct {
	Id          int64
	Email       string
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
	Id      int64
	Title   string
	Content string
	Status  ArticleStatus
	Author  Author
}

type Author struct {
	Id   int64
	Name string
}
