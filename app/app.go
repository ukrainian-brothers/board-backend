package app

import "github.com/ukrainian-brothers/board-backend/app/board"

type Commands struct {
	AddAdvert board.AddAdvert
	AddUser   board.AddUser
}

type Queries struct {
	GetAdvert          board.GetAdvert
	GetUserByLogin     board.GetUserByLogin
	UserExists         board.UserExists
	VerifyUserPassword board.VerifyUserPassword
}

type Application struct {
	Commands Commands
	Queries  Queries
}
