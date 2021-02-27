package alpclib

import (
	"log"
	"testing"
)

func TestRezkaListFilms(t *testing.T) {
	r := Rezka{}
	films, err := r.ListFilms(&ListParameters{
		Search: "Harry Potter",
	})
	if err != nil {
		t.Error(err)
	}
	if len(films) == 0 {
		t.Error("Empty films")
	}
}

func TestRezkaListSeries(t *testing.T) {
	r := Rezka{}
	series, err := r.ListSeries(&ListParameters{
		Search: "Game Of Thrones",
	})
	if err != nil {
		t.Error(err)
	}
	if len(series) == 0 {
		t.Error("Empty series")
	}
}

func TestRezkaGetFilm(t *testing.T) {
	r := Rezka{}
	film, err := r.GetFilm("aHR0cHM6Ly9yZXprYS5hZy9maWxtcy9hY3Rpb24vMjM5LWdhcnJpLXBvdHRlci1pLWRhcnktc21lcnRpLWNoYXN0LWktMjAxMC5odG1s")
	if err != nil {
		t.Error(err)
	}
	if len(film.Sources) == 0 {
		t.Error("Empty sources")
	}
}

func TestRezkaGetSeries(t *testing.T) {
	r := Rezka{}
	series, err := r.GetSeries("aHR0cHM6Ly9yZXprYS5hZy9zZXJpZXMvZmFudGFzeS80NS1pZ3JhLXByZXN0b2xvdi0yMDExLWdvdC1vbmxpbmUuaHRtbA==", 1, 1)
	if err != nil {
		t.Error(err)
	}
	if len(series.Sources) == 0 {
		t.Error("Empty sources")
	}
	log.Println(series)
}
