package alpclib

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var streamsregex = regexp.MustCompile(`(?m)"streams":"([\[\]\:\\\/\.\,\w ]+)"`)

type Rezka struct{}

type rezkaEpisode struct {
	ID     int
	Season int
}

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
	doc, err := goquery.NewDocument(listurl.String())
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
	doc, err := goquery.NewDocument(pageurl)
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
		req, _ := http.NewRequest("POST", "https://rezka.ag/ajax/get_cdn_series/", strings.NewReader(data.Encode()))
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
			quality, source, err := r.rawSourceQuality(rawsource)
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
	if len(sources) == 0 {
		pagehtml, _ := doc.Html()
		streams := streamsregex.FindAllString(pagehtml, -1)
		if len(streams) != 0 {
			rawsources := strings.Split(streams[0], ",")
			for _, rawsource := range rawsources {
				quality, source, err := r.rawSourceQuality(rawsource)
				if err != nil {
					continue
				}
				videourls := strings.Split(source, " or ")
				videourl := videourls[len(videourls)-1]
				videourl = strings.ReplaceAll(videourl, "\\", "")
				sources = append(sources, FilmSource{
					Voicecover: "Default",
					Quality:    quality,
					URL:        videourl,
				})
			}
		}
	}
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

func (r *Rezka) GetSeries(id string, season int, episode int) (Series, error) {
	_pageurl, err := base64.StdEncoding.DecodeString(id)
	if err != nil {
		return Series{}, err
	}
	pageurl := string(_pageurl)
	doc, err := goquery.NewDocument(pageurl)
	if err != nil {
		return Series{}, err
	}
	dataidstr, _ := doc.Find(".b-simple_episode__item").First().Attr("data-id")
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
	sources := []SeriesSource{}
	doc.Find(".b-translator__item").Each(func(i int, s *goquery.Selection) {
		voicecover, _ := s.Attr("title")
		tidstr, _ := s.Attr("data-translator_id")
		data := url.Values{}
		data.Set("id", dataidstr)
		data.Set("translator_id", tidstr)
		data.Set("action", "get_episodes")
		req, _ := http.NewRequest("POST", "https://rezka.ag/ajax/get_cdn_series/", strings.NewReader(data.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return
		}
		defer res.Body.Close()
		var respdata map[string]interface{}
		json.NewDecoder(res.Body).Decode(&respdata)
		epnode, _ := goquery.NewDocumentFromReader(strings.NewReader(respdata["episodes"].(string)))
		epnode.Find(".b-simple_episode__item").Each(func(i int, s *goquery.Selection) {
			seasonstr, _ := s.Attr("data-season_id")
			episodestr, _ := s.Attr("data-episode_id")
			_season, _ := strconv.Atoi(seasonstr)
			_episode, _ := strconv.Atoi(episodestr)
			if (season == 0 || season == _season) && (episode == 0 || episode == _episode) {
				data := url.Values{}
				data.Set("id", dataidstr)
				data.Set("translator_id", tidstr)
				data.Set("season", seasonstr)
				data.Set("episode", episodestr)
				data.Set("action", "get_stream")
				req, _ := http.NewRequest("POST", "https://rezka.ag/ajax/get_cdn_series/", strings.NewReader(data.Encode()))
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
					quality, source, err := r.rawSourceQuality(rawsource)
					if err != nil {
						continue
					}
					videourls := strings.Split(source, " or ")
					videourl := videourls[len(videourls)-1]
					sources = append(sources, SeriesSource{
						Voicecover: voicecover,
						Season:     _season,
						Episode:    _episode,
						Quality:    quality,
						URL:        videourl,
					})
				}
			} else {
				sources = append(sources, SeriesSource{
					Voicecover: voicecover,
					Season:     _season,
					Episode:    _episode,
					Quality:    "",
					URL:        "",
				})
			}

		})
	})
	return Series{
		ID:          id,
		Name:        name,
		Description: description,
		PosterURL:   poster,
		StartYear:   year,
		Country:     country,
		Sources:     sources,
	}, nil
}

func (*Rezka) rawSourceQuality(rawsource string) (string, string, error) {
	if strings.Contains(rawsource, "[360p]") {
		return "240p", strings.ReplaceAll(rawsource, "[360p]", ""), nil
	} else if strings.Contains(rawsource, "[480p]") {
		return "360p", strings.ReplaceAll(rawsource, "[480p]", ""), nil
	} else if strings.Contains(rawsource, "[720p]") {
		return "480p", strings.ReplaceAll(rawsource, "[720p]", ""), nil
	} else if strings.Contains(rawsource, "[1080p]") {
		return "720p", strings.ReplaceAll(rawsource, "[1080p]", ""), nil
	} else if strings.Contains(rawsource, "[1080p Ultra]") {
		return "1080p", strings.ReplaceAll(rawsource, "[1080p Ultra]", ""), nil
	}
	return "", "", errors.New("Quality is not determined")
}
