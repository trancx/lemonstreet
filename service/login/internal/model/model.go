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

type LoginInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}