package transaction

import (
	enMovie "abishar/internal/entity/movies"
	enVotes "abishar/internal/entity/votes"
	"context"
	"errors"
	"fmt"
)

type (
	voteRepository interface {
		Vote(ctx context.Context, form enVotes.VoteRequest) (int64, error)
		UnVote(ctx context.Context, userID, movieID int64) error
		GetMostVotedMovie(ctx context.Context) (*enVotes.VotedMovies, error)
		GetVotedMovieByUser(ctx context.Context, userID int64) ([]enVotes.VotedMovies, error)
		GetVote(ctx context.Context, userID int64, movieID int64) (*enVotes.Vote, error)
	}

	movieRepository interface {
		GetMovie(ctx context.Context, movieID int64) (*enMovie.Movies, error)
	}
)

type Usecase struct {
	voteRepo  voteRepository
	movieRepo movieRepository
}

func NewUsecase(
	voteRepo voteRepository,
	movieRepo movieRepository,
) *Usecase {
	return &Usecase{
		voteRepo:  voteRepo,
		movieRepo: movieRepo,
	}
}

func (uc *Usecase) Vote(ctx context.Context, form enVotes.VoteRequest) (int64, error) {
	movie, err := uc.movieRepo.GetMovie(ctx, form.MovieID)
	if err != nil {
		return 0, errors.New(fmt.Sprintf("Failed to get movie. err: %+v", err))
	}

	if movie.ID == 0 {
		return 0, errors.New(fmt.Sprintf("movie not exist. err: %+v", err))
	}

	vote, err := uc.voteRepo.GetVote(ctx, form.UserID, form.MovieID)
	if err != nil {
		return 0, errors.New(fmt.Sprintf("Failed to get vote. err: %+v", err))
	}

	if vote.ID != 0 {
		return 0, errors.New(fmt.Sprintf("user already vote this movie"))
	}

	voteID, err := uc.voteRepo.Vote(ctx, form)
	if err != nil {
		return 0, errors.New(fmt.Sprintf("failed to create transaction. err: %+v", err))
	}

	return voteID, nil
}

func (uc *Usecase) UnVote(ctx context.Context, form enVotes.VoteRequest) error {
	movie, err := uc.movieRepo.GetMovie(ctx, form.MovieID)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to get movie. err: %+v", err))
	}

	if movie.ID == 0 {
		return errors.New(fmt.Sprintf("movie not exist. err: %+v", err))
	}

	vote, err := uc.voteRepo.GetVote(ctx, form.UserID, form.MovieID)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to get vote. err: %+v", err))
	}

	if vote.ID == 0 {
		return errors.New(fmt.Sprintf("user are not voting this movie"))
	}

	err = uc.voteRepo.UnVote(ctx, form.UserID, form.MovieID)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to unvote. err: %+v", err))
	}

	return nil
}

func (uc *Usecase) GetMostVotedMovie(ctx context.Context) (*enVotes.VotedMovies, error) {
	movie, err := uc.voteRepo.GetMostVotedMovie(ctx)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed get most voted movies. err: %+v", err))
	}

	return movie, nil
}

func (uc *Usecase) GetVotedMovieByUser(ctx context.Context, userID int64) ([]enVotes.VotedMovies, error) {
	movies, err := uc.voteRepo.GetVotedMovieByUser(ctx, userID)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed get most voted movies. err: %+v", err))
	}

	return movies, nil
}
