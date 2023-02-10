package movies

import (
	enMovies "abishar/internal/entity/movies"
	"abishar/internal/pkg/redigo"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

type (
	redis interface {
		Del(keys ...string) error
		Get(key string) *redigo.Result
		Keys(key string) *redigo.Result
		Setex(key string, expireTime int, value interface{}) error
	}
)

type Repository struct {
	database *sqlx.DB
	redis    redis
}

func NewRepository(
	db *sqlx.DB,
	redis redis,
) *Repository {
	return &Repository{
		database: db,
		redis:    redis,
	}
}

const (
	expireTime = 3600
	movieKey   = "movie-%d"
)

func (r *Repository) InsertMovie(ctx context.Context, form enMovies.MovieRequest) (int64, error) {
	result, err := r.database.ExecContext(ctx, `
    insert into movies 
      (title, description, duration, artists, genres, watch_url)
    values ($1, $2, $3, $4, $5, $6)
  `, form.Title, form.Description, form.Duration, form.Artists, form.Genres, form.WatchURL)
	if err != nil {
		log.Printf("[InsertMovies] failed to insert movie. err: %+v", err)
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("[InsertMovies] failed to insert movei. err: %+v", err)
		return 0, err
	}

	return id, nil
}

func (r *Repository) UpdateMovies(ctx context.Context, form enMovies.Movies) error {
	_, err := r.database.ExecContext(ctx, `
    update movies 
     set title=$1, description=$2, duration=$3, artists=$4, genres=$5, watch_url=$6
    where id = $7
  `, form.Title, form.Description, form.Duration, form.Artists, form.Genres, form.WatchURL, form.ID)
	if err != nil {
		log.Printf("[InsertMovie] failed to update movie. err: %+v", err)
		return err
	}

	key := fmt.Sprintf(movieKey, form.ID)

	err = r.redis.Del(key)
	if err != nil {
		log.Printf("Failed to delete redis for key %s. err: %v", key, err)
	}

	return nil
}

func (r *Repository) UpdateViews(ctx context.Context, movieID int64) error {
	_, err := r.database.ExecContext(ctx, `
    update movies 
     set views=views+1
    where id = $1
  `, movieID)
	if err != nil {
		log.Printf("[UpdateViews] failed to add views. err: %+v", err)
		return err
	}

	key := fmt.Sprintf(movieKey, movieID)

	err = r.redis.Del(key)
	if err != nil {
		log.Printf("Failed to delete redis for key %s. err: %v", key, err)
	}

	return nil
}

func (r *Repository) GetMovies(ctx context.Context, limit, offset int) ([]enMovies.Movies, error) {
	movies := make([]enMovies.Movies, 0)

	err := r.database.SelectContext(ctx, &movies, `
    select
      id, title, description, duration, artists, genres, watch_url
    from movies 
    limit $1
    offset $2
  `, limit, offset)
	if err != nil {
		if err == sql.ErrNoRows {
			return movies, nil
		}

		log.Printf("[GetMovies] Failed to get movies. err: %+v", err)
		return movies, err
	}

	return movies, nil
}

func (r *Repository) SearchMovies(ctx context.Context, title, description, artists, genres string, limit, offset int) ([]enMovies.Movies, error) {
	products := make([]enMovies.Movies, 0)

	err := r.database.SelectContext(ctx, &products, `
    select
      id,title, description, duration, artists, genres, watch_url
    from movies 
    where title ilike '%$1%'
    and description ilike '%$2%'
    and artists ilike '%$3%'
    and genres ilike '%$4%'
    limit $5
    offset $6
  `, title, description, artists, genres, limit, offset)
	if err != nil {
		if err == sql.ErrNoRows {
			return products, nil
		}

		log.Printf("[GetProducts] Failed to get movies. err: %+v", err)
		return products, err
	}

	return products, nil
}

func (r *Repository) GetMovie(ctx context.Context, movieID int64) (*enMovies.Movies, error) {

	movie := &enMovies.Movies{}

	key := fmt.Sprintf(movieKey, movieID)
	cache := r.redis.Get(key)
	cacheByte, ok := cache.Value.([]byte)
	if ok {
		err := json.Unmarshal(cacheByte, movie)
		if err != nil {
			log.Printf("failed to unmarshal movie")
			return movie, nil
		}
	}

	err := r.database.GetContext(ctx, movie, `
    select * from movies where id=$1
  `, movieID)

	if err != nil {
		if err == sql.ErrNoRows {
			return movie, nil
		}

		log.Printf("failed to get movie. err: %+v", err)
		return movie, err
	}

	movieByte, _ := json.Marshal(*movie)

	err = r.redis.Setex(key, expireTime, movieByte)
	if err != nil {
		log.Printf("Failed to get redis for key %s. err: %v", key, err)
	}

	return movie, nil
}
