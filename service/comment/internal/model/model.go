package model

// Kratos hello kratos.
type Kratos struct {
	Hello string
}

type Article struct {
	ID int64
	Content string
	Author string
}

type Comment struct {
	ID	int64
	UserID	int64
	ArtID	int64
	Content	string
}