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
}

type SeriesSource struct {
	Voicecover string
	Season     int
	Series     int
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
	GetSeries(id string) Series
}
