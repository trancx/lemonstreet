package model

import (
	accapi "article/api/accapi"
	"article/api/artapi"
	"article/api/cmtapi"
)

type Article struct {
	ID          	int64	`json:"id"`
	//UserID			int64
	//Author      	string
	//Title			string
	//CreatedTime		string
	//Description 	string
	Content     	string	`json:"content"`
}

type ArticleBaseInfo struct {
	ID				int64
	UserID			int64
	Author      	string
	Title			string
	Description 	string
}

type PostArticle struct {
	Title		string			`json:"title"`
	Content		string			`json:"content"`
}

// acceleration
type ArticleInfo struct {
	 UInfo *accapi.UserBaseInfo	`json:"user"`
	 ABI *artapi.ArticleBaseInfo	`json:"ainfo"`
	 Content *Article			`json:"content"`
	 Comments []*cmtapi.Comment	`json:"comments"`
}