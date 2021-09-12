package model

type User struct {
	Id      int		`redis:"id"`
	Balance int64	`redis:"balance"`
}
