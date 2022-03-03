package app

import "github.com/ukrainian-brothers/board-backend/app/board"

type Commands struct {
	AddAdvert board.AddAdvert
	EditAdvert board.EditAdvert
}

type Queries struct {
	GetAdvert board.GetAdvert
}

type Application struct {
	Commands Commands
	Queries Queries
}
