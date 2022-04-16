package main

import (
	"fmt"
	"testing"
)

func TestGetGenres(m *testing.T) {

	res := getGenres("es")
	if res != 0 {
		fmt.Print("error")
	}

	for k, v := range mapGenres {
		fmt.Sprintf("%s %s", k, v)
	}
}

func TestGetMovies(m *testing.T) {
	res, movies := getMovies("sonic", "es")
	if res.Result == 0 {
		for _, v := range movies {
			fmt.Sprintf("%s %s", v.Titulo, v.Sinopsis)

		}
	}
}
