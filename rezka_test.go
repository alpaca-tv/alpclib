package alpclib

import (
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
