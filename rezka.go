package alpclib

import (
	"encoding/base64"
	"encoding/json"
	"errors"
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
		pageurl, _ := s.Find("a").Attr("href")
		id := base64.StdEncoding.EncodeToString([]byte(pageurl))
		poster, _ := s.Find("img").Attr("src")
		name := s.Find(".b-content__inline_item-link").Find("a").Text()
		descrow := strings.Split(s.Find(".b-content__inline_item-link").Find("div").Text(), ", ")
		year, _ := strconv.Atoi(descrow[0])
		country := descrow[1]
		films = append(films, Film{
			ID:        id,
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
		pageurl, _ := s.Find("a").Attr("href")
		id := base64.StdEncoding.EncodeToString([]byte(pageurl))
		poster, _ := s.Find("img").Attr("src")
		name := s.Find(".b-content__inline_item-link").Find("a").Text()
		descrow := strings.Split(s.Find(".b-content__inline_item-link").Find("div").Text(), ", ")
		year, _ := strconv.Atoi(descrow[0])
		country := descrow[1]
		series = append(series, Series{
			ID:        id,
			Name:      name,
			PosterURL: poster,
			EndYear:   year,
			Country:   country,
		})
	})
	return series, nil
}

func (*Rezka) AllowedOrders() []string {
	return []string{
		"last",
		"popular",
		"watching",
	}
}

func (*Rezka) AllowedGenres() []string {
	return []string{}
}

func (r *Rezka) GetFilm(id string) (Film, error) {
	_pageurl, err := base64.StdEncoding.DecodeString(id)
	if err != nil {
		return Film{}, err
	}
	pageurl := string(_pageurl)
	resp, err := http.DefaultClient.Get(pageurl)
	if err != nil {
		return Film{}, err
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return Film{}, err
	}
	year := 0
	country := ""
	name := doc.Find(".b-post__title").Find("h1").Text()
	description := doc.Find(".b-post__description_text").Text()
	poster, _ := doc.Find(".b-sidecover").Find("img").Attr("src")
	doc.Find(".b-post__info").Find("tr").Each(func(i int, s *goquery.Selection) {
		nodehtml, _ := s.Html()
		if strings.Contains(nodehtml, "Дата выхода") {
			yearstr := strings.ReplaceAll(s.Find("a").Text(), " года", "")
			_year, _ := strconv.Atoi(yearstr)
			year = _year
		}
		if strings.Contains(nodehtml, "Страна") {
			country = s.Find("a").First().Text()
		}
	})
	sources := []FilmSource{}
	doc.Find(".b-translator__item").Each(func(i int, s *goquery.Selection) {
		voicecover, _ := s.Attr("title")
		_id, _ := s.Attr("data-id")
		_tid, _ := s.Attr("data-translator_id")
		_camrip, _ := s.Attr("data-camrip")
		_ads, _ := s.Attr("data-ads")
		_director, _ := s.Attr("data-director")
		data := url.Values{}
		data.Set("id", _id)
		data.Set("translator_id", _tid)
		data.Set("is_camrip", _camrip)
		data.Set("is_ads", _ads)
		data.Set("is_director", _director)
		data.Set("action", "get_movie")
		req, err := http.NewRequest("POST", "https://rezka.ag/ajax/get_cdn_series/", strings.NewReader(data.Encode()))
		if err != nil {
			return
		}
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return
		}
		defer res.Body.Close()
		var respdata map[string]interface{}
		json.NewDecoder(res.Body).Decode(&respdata)
		rawsources := strings.Split(respdata["url"].(string), ",")
		for _, rawsource := range rawsources {
			quality, source, err := r.RawSourceQuality(rawsource)
			if err != nil {
				continue
			}
			videourls := strings.Split(source, " or ")
			videourl := videourls[len(videourls)-1]
			sources = append(sources, FilmSource{
				Voicecover: voicecover,
				Quality:    quality,
				URL:        videourl,
			})
		}
	})
	return Film{
		ID:          id,
		Name:        name,
		Description: description,
		PosterURL:   poster,
		Year:        year,
		Country:     country,
		Sources:     sources,
	}, nil
}

func (*Rezka) GetSeries(id string) Series {
	return Series{}
}

func (*Rezka) RawSourceQuality(rawsource string) (string, string, error) {
	if strings.Contains(rawsource, "[360p]") {
		return "360p", strings.ReplaceAll(rawsource, "[360p]", ""), nil
	} else if strings.Contains(rawsource, "[480p]") {
		return "480p", strings.ReplaceAll(rawsource, "[480p]", ""), nil
	} else if strings.Contains(rawsource, "[720p]") {
		return "720p", strings.ReplaceAll(rawsource, "[720p]", ""), nil
	} else if strings.Contains(rawsource, "[1080p]") {
		return "1080p", strings.ReplaceAll(rawsource, "[1080p]", ""), nil
	} else if strings.Contains(rawsource, "[1080p]") {
		return "1080p", strings.ReplaceAll(rawsource, "[1080p]", ""), nil
	} else if strings.Contains(rawsource, "[1080p Ultra]") {
		return "1080p Ultra", strings.ReplaceAll(rawsource, "[1080p Ultra]", ""), nil
	}
	return "", "", errors.New("Quality is not determined")
}
