package model

// Kratos hello kratos.
type Kratos struct {
	Hello string
}

type Article struct {
	ID          	int64
	//UserID			int64
	//Author      	string
	//Title			string
	//CreatedTime		string
	//Description 	string
	Content     	string
}

type ArticleBaseInfo struct {
	ID				int64
	UserID			int64
	Author      	string
	Title			string
	Description 	string
}