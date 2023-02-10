package movies

type Movies struct {
	ID          int64  `json:"id" db:"id"`
	Title       string `json:"title" db:"title"`
	Description string `json:"description" db:"description"`
	Duration    int64  `json:"duration" db:"duration"`
	Artists     string `json:"artists" db:"artists"`
	Genres      string `json:"genres" db:"genres"`
	WatchURL    string `json:"watchUrl" db:"watch_url"`
}

type MovieRequest struct {
	Title       string `json:"title" db:"title"`
	Description string `json:"description"`
	Duration    int64  `json:"duration" db:"duration"`
	Artists     string `json:"artists" db:"artists"`
	Genres      string `json:"genres" db:"genres"`
	WatchURL    string `json:"watchUrl" db:"watch_url"`
}
