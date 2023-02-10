package votes

import (
	enVote "abishar/internal/entity/votes"
	"abishar/internal/pkg/redigo"
	"context"
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

func (r *Repository) Vote(ctx context.Context, form enVote.VoteRequest) (int64, error) {
	result, err := r.database.ExecContext(ctx, `
    insert into votes 
      (user_id, movie_id)
    values ($1, $2 )
  `, form.UserID, form.MovieID)
	if err != nil {
		log.Printf("failed to vote. err: %+v", err)
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("failed to vote. err: %+v", err)
		return 0, err
	}

	return id, nil
}

func (r *Repository) UnVote(ctx context.Context, userID, movieID int64) error {
	_, err := r.database.ExecContext(ctx, `
    delete from votes 
    where user_id = $1 and movie_id = $2
  `, userID, movieID)
	if err != nil {
		log.Printf("failed to unvote. err: %+v", err)
		return err
	}

	return nil
}

func (r *Repository) GetVote(ctx context.Context, userID int64, movieID int64) (*enVote.Vote, error) {
	vote := &enVote.Vote{}

	err := r.database.SelectContext(ctx, vote, `
    select
     id, user_id, movie_id
    from vote
    where user_id=$1 and movie_id=$2
    `, userID, movieID)
	if err != nil {
		log.Printf("failed to get most voted. err: %+v", err)
		return vote, err
	}

	return vote, nil

}

func (r *Repository) GetMostVotedMovie(ctx context.Context) (*enVote.VotedMovies, error) {
	movie := &enVote.VotedMovies{}

	err := r.database.SelectContext(ctx, movie, `
    select
      (select count(*) from votes where movie_id=m.id) as votes,
      m.id, m.title, m.description, m.duration, m.artists, m.genres, m.watch_url
    from movies m
    inner join votes v on v.movie_id = m.id
    order by (select count(*) from votes where movie_id=m.id) desc
    limit 1
    `)
	if err != nil {
		log.Printf("failed to get most voted. err: %+v", err)
		return movie, err
	}

	return movie, nil
}

func (r *Repository) GetVotedMovieByUser(ctx context.Context, userID int64) ([]enVote.VotedMovies, error) {
	movies := make([]enVote.VotedMovies, 0)

	err := r.database.SelectContext(ctx, &movies, `
    select
      (select count(*) from votes where movie_id=m.id) as votes,
      m.id, m.title, m.description, m.duration, m.artists, m.genres, m.watch_url
    from movies m
    inner join votes v on v.movie_id = m.id and v.user_id = $1
    order by (select count(*) from votes where movie_id=m.id) desc
  `, userID)
	if err != nil {
		log.Printf("[GetTransactionsByUser] failed to get transaction for user_id %d. err: %+v", userID, err)
		return movies, err
	}

	return movies, nil
}
