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

type Format struct {
	Method string	`json:"method"`
	API string	`json:"api"`
	Params interface{}	`json:"params"`
}