package alpclib

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Rezka struct{}

func (*Rezka) ListFilms(p *ListParameters) ([]Film, error) {
	var listurl *url.URL
	if p.Search != "" {
		_url, _ := url.Parse("https://rezka.ag/search")
		q := _url.Query()
		q.Set("do", "search")
		q.Set("subaction", "search")
		q.Set("q", p.Search)
		_url.RawQuery = q.Encode()
		listurl = _url
	} else {
		_url, _ := url.Parse("https://rezka.ag/films/")
		q := _url.Query()
		if p.OrderBy != "" {
			q.Set("filter", p.OrderBy)
		}
		_url.RawQuery = q.Encode()
		listurl = _url
	}
	resp, err := http.DefaultClient.Get(listurl.String())
	if err != nil {
		return []Film{}, err
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return []Film{}, err
	}
	films := []Film{}
	doc.Find(".b-content__inline_item").Each(func(i int, s *goquery.Selection) {
		entrytype := s.Find(".entity").Text()
		if entrytype != "Фильм" {
			return
		}
		poster, _ := s.Find("img").Attr("src")
		name := s.Find(".b-content__inline_item-link").Find("a").Text()
		descrow := strings.Split(s.Find(".b-content__inline_item-link").Find("div").Text(), ", ")
		year, _ := strconv.Atoi(descrow[0])
		country := descrow[1]
		films = append(films, Film{
			Name:      name,
			PosterURL: poster,
			Year:      year,
			Country:   country,
		})
	})
	return films, nil
}

func (*Rezka) ListSeries(p *ListParameters) ([]Series, error) {
	var listurl *url.URL
	if p.Search != "" {
		_url, _ := url.Parse("https://rezka.ag/search")
		q := _url.Query()
		q.Set("do", "search")
		q.Set("subaction", "search")
		q.Set("q", p.Search)
		_url.RawQuery = q.Encode()
		listurl = _url
	} else {
		_url, _ := url.Parse("https://rezka.ag/series/")
		q := _url.Query()
		if p.OrderBy != "" {
			q.Set("filter", p.OrderBy)
		}
		_url.RawQuery = q.Encode()
		listurl = _url
	}
	resp, err := http.DefaultClient.Get(listurl.String())
	if err != nil {
		return []Series{}, err
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return []Series{}, err
	}
	series := []Series{}
	doc.Find(".b-content__inline_item").Each(func(i int, s *goquery.Selection) {
		entrytype := s.Find(".entity").Text()
		if entrytype != "Сериал" {
			return
		}
		poster, _ := s.Find("img").Attr("src")
		name := s.Find(".b-content__inline_item-link").Find("a").Text()
		descrow := strings.Split(s.Find(".b-content__inline_item-link").Find("div").Text(), ", ")
		year, _ := strconv.Atoi(descrow[0])
		country := descrow[1]
		series = append(series, Series{
			Name:      name,
			PosterURL: poster,
			EndYear:   year,
			Country:   country,
		})
	})
	return series, nil
}

func AllowedOrders() []string {
	return []string{
		"last",
		"popular",
		"watching",
	}
}

func AllowedGenres() []string {
	return []string{}
}

func GetFilm(id string) Film {
	return Film{}
}

func GetSeries(id string) Series {
	return Series{}
}
