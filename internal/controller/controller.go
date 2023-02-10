package controller

import (
	"abishar/internal/controller/movie"
	"abishar/internal/controller/user"
	"abishar/internal/controller/vote"
)

type Controllers struct {
	User  *user.Controller
	Movie *movie.Controller
	Vote  *vote.Controller
}

func NewControllers(
	user *user.Controller,
	movie *movie.Controller,
	vote *vote.Controller,
) *Controllers {
	return &Controllers{
		User:  user,
		Movie: movie,
		Vote:  vote,
	}
}
