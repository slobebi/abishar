package votes

type Vote struct {
	ID      int64 `json:"id" db:"id"`
	UserID  int64 `json:"userID" db:"user_id"`
	MovieID int64 `json:"movieID" db:"movie_id"`
}

type VotedMovies struct {
	ID          int64  `json:"id" db:"id"`
	Title       string `json:"title" db:"title"`
	Description string `json:"description" db:"description"`
	Duration    int64  `json:"duration" db:"duration"`
	Artists     string `json:"artists" db:"artists"`
	Genres      string `json:"genres" db:"genres"`
	WatchURL    string `json:"watchUrl" db:"watch_url"`
	Votes       int64  `json:"votes" db:"votes"`
}
type VoteRequest struct {
	UserID  int64 `json:"userID"`
	MovieID int64 `json:"movieID"`
}
