package product

import (
	"context"
	"errors"
	"fmt"
	"log"

	enMovie "abishar/internal/entity/movies"
)


type (
  movieRepository interface {
    InsertMovie(ctx context.Context, form enMovie.MovieRequest) (int64, error)
    UpdateMovies(ctx context.Context, form enMovie.Movies) error
    UpdateViews(ctx context.Context, movieID int64) error
    GetMovies(ctx context.Context, limit, offset int) ([]enMovie.Movies, error)
    GetMovie(ctx context.Context, movieID int64) (*enMovie.Movies, error)
    SearchMovies(ctx context.Context, title, description, artists, genres string, limit, offset int) ([]enMovie.Movies, error)
  }
)

type Usecase struct {
  movieRepo movieRepository
}

func NewUsecase(
  movieRepo movieRepository,
) *Usecase {
  return &Usecase{
    movieRepo: movieRepo,
  }
}

func (uc *Usecase) InsertMovie(ctx context.Context, form enMovie.MovieRequest) (int64, error) {
  id, err := uc.movieRepo.InsertMovie(ctx, form)
  if err != nil {
    return 0, errors.New(fmt.Sprintf("Failed to insert movie. err: %+v", err))
  }

  return id, nil
}

func (uc *Usecase) UpdateMovie(ctx context.Context, form enMovie.Movies) error {
  // check existing product
  movie, err := uc.movieRepo.GetMovie(ctx, form.ID)
  if err != nil {
    return errors.New(fmt.Sprintf("Failed to check movie. err: %+v", err))
  }

  if movie.ID == 0 {
    return errors.New(fmt.Sprintf("Movie not existed"))
  }

  err = uc.movieRepo.UpdateMovies(ctx, form)
  if err != nil {
    return errors.New(fmt.Sprintf("Failed to update product. err: %+v", err))
  }

  return nil
}

func (uc *Usecase) GetMovies(ctx context.Context, page, limit int) ([]enMovie.Movies, error) {
  offset := (page - 1) * limit

  movies, err := uc.movieRepo.GetMovies(ctx, limit, offset)
  if err != nil {
    return make([]enMovie.Movies, 0), errors.New(fmt.Sprintf("Failed to get all products. err: %+v", err))
  }

  for _, v := range movies {
    err := uc.movieRepo.UpdateViews(ctx, v.ID)
    if err != nil {
      log.Printf("failed to add views. err: %v", err)
    }
  }

  return movies, nil
}


func (uc *Usecase) SearchMovies(ctx context.Context, page, limit int, title, description, artists, genres string) ([]enMovie.Movies, error) {
  offset := (page - 1) * limit

  movies, err := uc.movieRepo.SearchMovies(ctx, title, description, artists, genres, limit, offset)
  if err != nil {
    return make([]enMovie.Movies, 0), errors.New(fmt.Sprintf("Failed to get movies. err: %+v", err))
  }

  for _, v := range movies {
    err := uc.movieRepo.UpdateViews(ctx, v.ID)
    if err != nil {
      log.Printf("failed to add views. err: %v", err)
    }
  }
  return movies, nil

}
