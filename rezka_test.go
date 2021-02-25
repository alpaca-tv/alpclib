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
	log.Println(films, err)
}
