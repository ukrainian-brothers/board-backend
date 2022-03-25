package app

import "github.com/ukrainian-brothers/board-backend/app/board"

type Commands struct {
	AddAdvert board.AddAdvert
	AddUser   board.AddUser
}

type Queries struct {
	GetAdvert board.GetAdvert
}

type Application struct {
	Commands Commands
	Queries  Queries
}
