package model

import (
	"comment/api/artapi"
	"comment/api/cmtapi"
)

// Kratos hello kratos.
type Kratos struct {
	Hello string
}

type Article struct {
	ID int64
	Content string
	Author string
}

type PostComment struct {
	ABI artapi.ArticleBaseInfo	`json:"ainfo"`
	Comment cmtapi.Comment		`json:"comment"`
}