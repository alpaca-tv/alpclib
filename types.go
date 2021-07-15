package alpclib

type ListParameters struct {
	OrderBy string
	Genres  []string
	Search  string
	Page    int
}

type Film struct {
	ID          string
	Name        string
	Description string
	PosterURL   string
	Year        int
	Rating      float64
	Country     string
	Genres      []string
	Related     []Film // List of ids
	Sources     []FilmSource
}

type FilmSource struct {
	Voicecover string
	Quality    string
	URL        string
}

type Series struct {
	ID          string
	Name        string
	Description string
	PosterURL   string
	StartYear   int
	EndYear     int
	Rating      float64
	Country     string
	Genres      []string
	Sources     []SeriesSource
}

type SeriesSource struct {
	Voicecover string
	Season     int
	Episode    int
	Quality    string
	URL        string
}

type Source interface {

	// Listings
	ListFilms(p *ListParameters) []Film
	ListSeries(p *ListParameters) []Series

	// Allowed listing parameters
	AllowedOrders() []string
	AllowedGenres() []string

	// Get exact
	GetFilm(id string) Film
	// Provide season and/or episode, if you want to reduce execution time
	// (Source URL and quality will be ignored). Use 0 for all, -1 for meta info only
	GetSeries(id string, season int, episode int) Series
}
