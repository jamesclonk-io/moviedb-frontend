package moviedb

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/jamesclonk-io/moviedb-backend/modules/database"
	"github.com/stretchr/testify/assert"
)

var (
	movieTestDbFile     string = "../../_fixtures/test.db"
	movieTestDbFileCopy string = "../../_fixtures/test_copy.db"
)

func init() {
	os.Setenv("JCIO_DATABASE_TYPE", "sqlite")
	os.Setenv("JCIO_DATABASE_URI", fmt.Sprintf("sqlite3://%s", movieTestDbFileCopy))

	copyFile(movieTestDbFile, movieTestDbFileCopy)
}

func copyFile(from, to string) {
	in, err := os.Open(from)
	if err != nil {
		panic(err)
	}
	defer in.Close()

	out, err := os.Create(to)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		panic(err)
	}
	if err := out.Close(); err != nil {
		panic(err)
	}
}

func getMovieDB() *movieDB {
	return NewMovieDB(database.NewAdapter()).(*movieDB)
}

func Test_MovieDB_Connection(t *testing.T) {
	mdb := getMovieDB()
	defer mdb.Close()

	if err := mdb.Ping(); err != nil {
		t.Fatal(err)
	}
}

func Test_MovieDB_GetMovie(t *testing.T) {
	mdb := getMovieDB()
	defer mdb.Close()

	movie, err := mdb.GetMovie("1")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "[1] Face/Off (1997)", movie.String())

	expectedMovie := &Movie{
		Id:      2,
		Title:   "Fight Club",
		Year:    1999,
		Score:   5,
		Rating:  16,
		Region:  "2",
		Format:  "16:9",
		Disks:   1,
		Type:    "DVD",
		Length:  135,
		Picture: "fight_club.jpg",
		Languages: []*Language{
			&Language{Id: 1, Name: "Deutsch", Country: "Schweiz", NativeName: "Deutsch"},
			&Language{Id: 2, Name: "Englisch", Country: "USA", NativeName: "English"},
		},
		Genres: []*Genre{
			&Genre{Id: 6, Name: "Drama"},
			&Genre{Id: 22, Name: "Mystery"},
			&Genre{Id: 21, Name: "Special\u0026nbsp;Style"},
			&Genre{Id: 4, Name: "Thriller"},
		},
		Actors: []*Person{
			&Person{Id: 7, Name: "Brad Pitt"},
			&Person{Id: 8, Name: "Edward Norton"},
			&Person{Id: 9, Name: "Helena Bonham Carter"},
			&Person{Id: 819, Name: "Jared Leto"},
			&Person{Id: 10, Name: "Meat Loaf"},
		},
		Directors: []*Person{
			&Person{Id: 11, Name: "David Fincher"},
		},
	}
	movie, err = mdb.GetMovie("2")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expectedMovie.Id, movie.Id)
	assert.Equal(t, expectedMovie.Title, movie.Title)
	assert.Equal(t, expectedMovie.Alttitle, movie.Alttitle)
	assert.Equal(t, expectedMovie.Year, movie.Year)
	assert.Equal(t, expectedMovie.Score, movie.Score)
	assert.Equal(t, expectedMovie.Rating, movie.Rating)
	assert.Equal(t, expectedMovie.Region, movie.Region)
	assert.Equal(t, expectedMovie.Format, movie.Format)
	assert.Equal(t, expectedMovie.Disks, movie.Disks)
	assert.Equal(t, expectedMovie.Type, movie.Type)
	assert.Equal(t, expectedMovie.Length, movie.Length)
	assert.Equal(t, expectedMovie.Picture, movie.Picture)

	assert.Equal(t, expectedMovie.Languages, movie.Languages)
	assert.Equal(t, expectedMovie.Genres, movie.Genres)
	assert.Equal(t, expectedMovie.Actors, movie.Actors)
	assert.Equal(t, expectedMovie.Directors, movie.Directors)
}

func Test_MovieDB_DeleteMovie(t *testing.T) {
	copyFile(movieTestDbFile, movieTestDbFileCopy)
	mdb := getMovieDB()
	defer mdb.Close()
	defer copyFile(movieTestDbFile, movieTestDbFileCopy)

	movie, err := mdb.GetMovie("7")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "[7] Austin Powers 2 (1999)", movie.String())

	rows, err := mdb.DeleteMovie("7")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, int64(10), rows)
}

func Test_MovieDB_AddMovie(t *testing.T) {
	copyFile(movieTestDbFile, movieTestDbFileCopy)
	mdb := getMovieDB()
	defer mdb.Close()
	defer copyFile(movieTestDbFile, movieTestDbFileCopy)

	expectedMovie := &Movie{
		Id:       0,
		Title:    "Testfilm",
		Alttitle: sql.NullString{"The ultimate movie!", true},
		Year:     2029,
		Score:    4,
		Rating:   12,
		Region:   "2",
		Format:   "16:9",
		Disks:    2,
		Type:     "BluRay",
		Length:   123,
		Picture:  "testfilm.jpg",
		Languages: []*Language{
			&Language{Name: "Deutsch"},
			&Language{Name: "Englisch"},
			&Language{Name: "Serbokroatisch"},
		},
		Genres: []*Genre{
			&Genre{Name: "Deutsche Soap"},
			&Genre{Name: "Drama"},
			&Genre{Name: "Mystery"},
			&Genre{Name: "Thriller"},
		},
		Actors: []*Person{
			&Person{Name: "Brad Pitt"},
			&Person{Name: "Edward Norton"},
			&Person{Name: "Looize de Testador"},
		},
		Directors: []*Person{
			&Person{Name: "David Fincher"},
			&Person{Name: "Senõr Spielbergo"},
		},
	}

	rows, err := mdb.Query(`select * from movie_people where name = 'Looize de Testador'`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	if rows.Next() {
		t.Fatal(errors.New("Looize de Testador should not yet exist"))
	}

	if err := mdb.AddMovie(expectedMovie); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 915, expectedMovie.Id)

	movie, err := mdb.GetMovie(strconv.Itoa(expectedMovie.Id))
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expectedMovie.Id, movie.Id)
	assert.Equal(t, expectedMovie.Title, movie.Title)
	assert.Equal(t, expectedMovie.Alttitle.String, movie.Alttitle.String)
	assert.Equal(t, expectedMovie.Year, movie.Year)
	assert.Equal(t, expectedMovie.Score, movie.Score)
	assert.Equal(t, expectedMovie.Rating, movie.Rating)
	assert.Equal(t, expectedMovie.Region, movie.Region)
	assert.Equal(t, expectedMovie.Format, movie.Format)
	assert.Equal(t, expectedMovie.Disks, movie.Disks)
	assert.Equal(t, expectedMovie.Type, movie.Type)
	assert.Equal(t, expectedMovie.Length, movie.Length)
	assert.Equal(t, expectedMovie.Picture, movie.Picture)

	assert.Equal(t, expectedMovie.Languages[0].Id, movie.Languages[0].Id)
	assert.Equal(t, expectedMovie.Languages[0].Name, movie.Languages[0].Name)
	assert.Equal(t, expectedMovie.Languages[1].Id, movie.Languages[1].Id)
	assert.Equal(t, expectedMovie.Languages[1].Name, movie.Languages[1].Name)
	assert.Equal(t, expectedMovie.Languages[2].Id, movie.Languages[2].Id)
	assert.Equal(t, expectedMovie.Languages[2].Name, movie.Languages[2].Name)
	assert.Equal(t, expectedMovie.Genres[0].Id, movie.Genres[0].Id)
	assert.Equal(t, expectedMovie.Genres[0].Name, movie.Genres[0].Name)
	assert.Equal(t, expectedMovie.Genres[1].Id, movie.Genres[1].Id)
	assert.Equal(t, expectedMovie.Genres[1].Name, movie.Genres[1].Name)
	assert.Equal(t, expectedMovie.Genres[2].Id, movie.Genres[2].Id)
	assert.Equal(t, expectedMovie.Genres[2].Name, movie.Genres[2].Name)
	assert.Equal(t, expectedMovie.Genres[3].Id, movie.Genres[3].Id)
	assert.Equal(t, expectedMovie.Genres[3].Name, movie.Genres[3].Name)
	assert.Equal(t, expectedMovie.Actors[0].Id, movie.Actors[0].Id)
	assert.Equal(t, expectedMovie.Actors[0].Name, movie.Actors[0].Name)
	assert.Equal(t, expectedMovie.Actors[1].Id, movie.Actors[1].Id)
	assert.Equal(t, expectedMovie.Actors[1].Name, movie.Actors[1].Name)
	assert.Equal(t, expectedMovie.Actors[2].Id, movie.Actors[2].Id)
	assert.Equal(t, expectedMovie.Actors[2].Name, movie.Actors[2].Name)
	assert.Equal(t, expectedMovie.Directors[0].Id, movie.Directors[0].Id)
	assert.Equal(t, expectedMovie.Directors[0].Name, movie.Directors[0].Name)
	assert.Equal(t, expectedMovie.Directors[1].Id, movie.Directors[1].Id)
	assert.Equal(t, expectedMovie.Directors[1].Name, movie.Directors[1].Name)

	// person checks -----------------------------------------------------------------------
	rows, err = mdb.Query(`select * from movie_people where name = 'Senõr Spielbergo'`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	if !rows.Next() {
		t.Error(errors.New("Senõr Spielbergo should exist now"))
	}

	rows, err = mdb.Query(`select * from movie_people mp
	join movie_link_director mld on (mld.person_id = mp.id)
	where mp.name = 'David Fincher'
	and mld.movie_id = '915'`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	if !rows.Next() {
		t.Error(errors.New("David Fincher should be a director of newly inserted movie"))
	}

	rows, err = mdb.Query(`select * from movie_people mp
	join movie_link_actor mla on (mla.person_id = mp.id)
	where mp.name = 'Looize de Testador'
	and mla.movie_id = '915'`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	if !rows.Next() {
		t.Error(errors.New("Looize de Testador should be an actor of newly inserted movie"))
	}

	rows, err = mdb.Query(`select * from movie_people mp
	join movie_link_actor mla on (mla.person_id = mp.id)
	where mp.name = 'Brad Pitt'
	and mla.movie_id = '915'`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	if !rows.Next() {
		t.Error(errors.New("Brad Pitt should be an actor of newly inserted movie"))
	}

	// language checks -----------------------------------------------------------------------
	rows, err = mdb.Query(`select * from movie_language ml
	join movie_link_language mll on (mll.language_id = ml.id)
	where ml.name = 'Serbokroatisch'
	and mll.movie_id = '915'`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	if !rows.Next() {
		t.Error(errors.New("Serbokroatisch should be a language of newly inserted movie"))
	}

	rows, err = mdb.Query(`select * from movie_language ml
	join movie_link_language mll on (mll.language_id = ml.id)
	where ml.name = 'Deutsch'
	and mll.movie_id = '915'`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	if !rows.Next() {
		t.Error(errors.New("Deutsch should be a language of newly inserted movie"))
	}

	// genre checks -----------------------------------------------------------------------
	rows, err = mdb.Query(`select * from movie_genre mg
	join movie_link_genre mlg on (mlg.genre_id = mg.id)
	where mg.name = 'Deutsche Soap'
	and mlg.movie_id = '915'`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	if !rows.Next() {
		t.Error(errors.New("Deutsche Soap should be a genre of newly inserted movie"))
	}

	rows, err = mdb.Query(`select * from movie_genre mg
	join movie_link_genre mlg on (mlg.genre_id = mg.id)
	where mg.name = 'Mystery'
	and mlg.movie_id = '915'`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	if !rows.Next() {
		t.Error(errors.New("Mystery should be a genre of newly inserted movie"))
	}
}

func Test_MovieDB_UpdateMovie(t *testing.T) {
	copyFile(movieTestDbFile, movieTestDbFileCopy)
	mdb := getMovieDB()
	defer mdb.Close()
	defer copyFile(movieTestDbFile, movieTestDbFileCopy)

	expectedMovie := &Movie{
		Id:       3,
		Title:    "Testfilm",
		Alttitle: sql.NullString{"The ultimate movie!", true},
		Year:     2029,
		Score:    4,
		Rating:   12,
		Region:   "2",
		Format:   "16:9",
		Disks:    2,
		Type:     "BluRay",
		Length:   123,
		Picture:  "testfilm.jpg",
		Languages: []*Language{
			&Language{Name: "Deutsch"},
			&Language{Name: "Englisch"},
			&Language{Name: "Serbokroatisch"},
		},
		Genres: []*Genre{
			&Genre{Name: "Deutsche Soap"},
			&Genre{Name: "Drama"},
			&Genre{Name: "Mystery"},
			&Genre{Name: "Thriller"},
		},
		Actors: []*Person{
			&Person{Name: "Brad Pitt"},
			&Person{Name: "Edward Norton"},
			&Person{Name: "Looize de Testador"},
		},
		Directors: []*Person{
			&Person{Name: "David Fincher"},
			&Person{Name: "Senõr Spielbergo"},
		},
	}

	rows, err := mdb.Query(`select * from movie_people where name = 'Looize de Testador'`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	if rows.Next() {
		t.Fatal(errors.New("Looize de Testador should not yet exist"))
	}

	if err := mdb.SaveMovie(expectedMovie); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, expectedMovie.Id, 3, "movie should have id 3")

	movie, err := mdb.GetMovie(strconv.Itoa(expectedMovie.Id))
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, expectedMovie.Id, movie.Id)
	assert.Equal(t, expectedMovie.Title, movie.Title)
	assert.Equal(t, expectedMovie.Alttitle.String, movie.Alttitle.String)
	assert.Equal(t, expectedMovie.Year, movie.Year)
	assert.Equal(t, expectedMovie.Score, movie.Score)
	assert.Equal(t, expectedMovie.Rating, movie.Rating)
	assert.Equal(t, expectedMovie.Region, movie.Region)
	assert.Equal(t, expectedMovie.Format, movie.Format)
	assert.Equal(t, expectedMovie.Disks, movie.Disks)
	assert.Equal(t, expectedMovie.Type, movie.Type)
	assert.Equal(t, expectedMovie.Length, movie.Length)
	assert.Equal(t, expectedMovie.Picture, movie.Picture)

	assert.Equal(t, expectedMovie.Languages[2].Id, movie.Languages[3].Id)
	assert.Equal(t, expectedMovie.Languages[2].Name, movie.Languages[3].Name)
	assert.Equal(t, expectedMovie.Genres[0].Id, movie.Genres[2].Id)
	assert.Equal(t, expectedMovie.Genres[0].Name, movie.Genres[2].Name)

	assert.Equal(t, expectedMovie.Actors[0].Id, movie.Actors[1].Id)
	assert.Equal(t, expectedMovie.Actors[0].Name, movie.Actors[1].Name)
	assert.Equal(t, expectedMovie.Actors[2].Id, movie.Actors[4].Id)
	assert.Equal(t, expectedMovie.Actors[2].Name, movie.Actors[4].Name)
	assert.Equal(t, expectedMovie.Directors[0].Id, movie.Directors[0].Id)
	assert.Equal(t, expectedMovie.Directors[0].Name, movie.Directors[0].Name)
	assert.Equal(t, expectedMovie.Directors[1].Id, movie.Directors[2].Id)
	assert.Equal(t, expectedMovie.Directors[1].Name, movie.Directors[2].Name)

	// person checks -----------------------------------------------------------------------
	rows, err = mdb.Query(`select * from movie_people where name = 'Senõr Spielbergo'`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	if !rows.Next() {
		t.Error(errors.New("Senõr Spielbergo should exist now"))
	}

	rows, err = mdb.Query(`select * from movie_people mp
	join movie_link_director mld on (mld.person_id = mp.id)
	where mp.name = 'David Fincher'
	and mld.movie_id = '3'`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	if !rows.Next() {
		t.Error(errors.New("David Fincher should be a director of updated movie"))
	}

	rows, err = mdb.Query(`select * from movie_people mp
	join movie_link_actor mla on (mla.person_id = mp.id)
	where mp.name = 'Looize de Testador'
	and mla.movie_id = '3'`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	if !rows.Next() {
		t.Error(errors.New("Looize de Testador should be an actor of updated movie"))
	}

	rows, err = mdb.Query(`select * from movie_people mp
	join movie_link_actor mla on (mla.person_id = mp.id)
	where mp.name = 'Brad Pitt'
	and mla.movie_id = '3'`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	if !rows.Next() {
		t.Error(errors.New("Brad Pitt should be an actor of updated movie"))
	}

	// language checks -----------------------------------------------------------------------
	rows, err = mdb.Query(`select * from movie_language ml
	join movie_link_language mll on (mll.language_id = ml.id)
	where ml.name = 'Serbokroatisch'
	and mll.movie_id = '3'`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	if !rows.Next() {
		t.Error(errors.New("Serbokroatisch should be a language of updated movie"))
	}

	rows, err = mdb.Query(`select * from movie_language ml
	join movie_link_language mll on (mll.language_id = ml.id)
	where ml.name = 'Deutsch'
	and mll.movie_id = '3'`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	if !rows.Next() {
		t.Error(errors.New("Deutsch should be a language of updatedmovie"))
	}

	// genre checks -----------------------------------------------------------------------
	rows, err = mdb.Query(`select * from movie_genre mg
	join movie_link_genre mlg on (mlg.genre_id = mg.id)
	where mg.name = 'Deutsche Soap'
	and mlg.movie_id = '3'`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	if !rows.Next() {
		t.Error(errors.New("Deutsche Soap should be a genre of updated movie"))
	}

	rows, err = mdb.Query(`select * from movie_genre mg
	join movie_link_genre mlg on (mlg.genre_id = mg.id)
	where mg.name = 'Mystery'
	and mlg.movie_id = '3'`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	if !rows.Next() {
		t.Error(errors.New("Mystery should be a genre of updated movie"))
	}
}

func Test_MovieDB_MovieListing(t *testing.T) {
	copyFile(movieTestDbFile, movieTestDbFileCopy)
	mdb := getMovieDB()
	defer mdb.Close()

	movies, err := mdb.GetMovieListings(MovieListingOptions{Sort: []Sort{NewSort("id", "asc")}})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, true, len(movies) > 910)

	expected := &MovieListing{
		Id:     2,
		Title:  "Fight Club",
		Year:   1999,
		Score:  5,
		Rating: 16,
	}
	assert.Equal(t, expected, movies[1])

	expected = &MovieListing{
		Id:     914,
		Title:  "Argo",
		Year:   2012,
		Score:  5,
		Rating: 12,
	}
	assert.Equal(t, expected, movies[len(movies)-1])

	movies, err = mdb.GetMovieListings(MovieListingOptions{Sort: []Sort{NewSort("score", "desc"), NewSort("title", "desc")}})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, true, len(movies) > 910)

	expected = &MovieListing{
		Id:     362,
		Title:  "The Wire (5)",
		Year:   2006,
		Score:  5,
		Rating: 16,
	}
	assert.Equal(t, expected, movies[1])

	movies, err = mdb.GetMovieListings(MovieListingOptions{Sort: []Sort{NewSort("score", "desc"), NewSort("title", "asc")}})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, true, len(movies) > 910)

	expected = &MovieListing{
		Id:     175,
		Title:  "24 (2)",
		Year:   2002,
		Score:  5,
		Rating: 16,
	}
	assert.Equal(t, expected, movies[0])

	movies, err = mdb.GetMovieListings(MovieListingOptions{
		Query: []Query{NewQuery("year", "2013"), NewQuery("score", "4")},
		Sort:  []Sort{NewSort("title", "desc")},
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 2, len(movies))

	expected = &MovieListing{
		Id:     898,
		Title:  "Star Trek Into Darkness",
		Year:   2013,
		Score:  4,
		Rating: 12,
	}
	assert.Equal(t, expected, movies[0])

	movies, err = mdb.GetMovieListings(MovieListingOptions{
		Query: []Query{NewQuery("year", "2013"), NewQuery("language", "3"), NewQuery("genre", "9")},
		Sort:  []Sort{NewSort("title", "desc")},
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 3, len(movies))

	expected = &MovieListing{
		Id:     890,
		Title:  "We\u0026#039;re the Millers",
		Year:   2013,
		Score:  3,
		Rating: 12,
	}
	assert.Equal(t, expected, movies[0])

	movies, err = mdb.GetMovieListings(MovieListingOptions{
		Query: []Query{NewQuery("actor", "331"), NewQuery("director", "331")},
		Sort:  []Sort{NewSort("title", "asc")},
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 2, len(movies))

	expected = &MovieListing{
		Id:     914,
		Title:  "Argo",
		Year:   2012,
		Score:  5,
		Rating: 12,
	}
	assert.Equal(t, expected, movies[0])

	movies, err = mdb.GetMovieListings(MovieListingOptions{
		Query: []Query{NewQuery("char", "num")},
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 23, len(movies))

	expected = &MovieListing{
		Id:     405,
		Title:  "10.000 BC",
		Year:   2008,
		Score:  2,
		Rating: 12,
	}
	assert.Equal(t, expected, movies[0])

	movies, err = mdb.GetMovieListings(MovieListingOptions{
		Query: []Query{NewQuery("char", "c")},
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 44, len(movies))

	expected = &MovieListing{
		Id:     543,
		Title:  "Californication (1)",
		Year:   2007,
		Score:  4,
		Rating: 18,
	}
	assert.Equal(t, expected, movies[0])

	movies, err = mdb.GetMovieListings(MovieListingOptions{
		Query: []Query{NewQuery("search", "Minutes")},
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 4, len(movies))

	expected = &MovieListing{
		Id:     481,
		Title:  "Blade Runner",
		Year:   1982,
		Score:  5,
		Rating: 16,
	}
	assert.Equal(t, expected, movies[2])
}

func Test_MovieDB_LanguagesByMovie(t *testing.T) {
	mdb := getMovieDB()
	defer mdb.Close()

	expected := []*Language{
		&Language{Id: 2, Name: "Englisch", Country: "USA", NativeName: "English"},
		&Language{Id: 3, Name: "Franz\u0026#246;sisch", Country: "Frankreich", NativeName: "Fran\u0026#231;ais"},
	}
	languages, err := mdb.GetLanguagesByMovie("3")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, expected, languages)
}

func Test_MovieDB_GenresByMovie(t *testing.T) {
	mdb := getMovieDB()
	defer mdb.Close()

	expected := []*Genre{
		&Genre{Id: 9, Name: "Comedy"},
		&Genre{Id: 23, Name: "Crime"},
		&Genre{Id: 13, Name: "Film\u0026nbsp;Noir"},
	}
	genres, err := mdb.GetGenresByMovie("3")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, expected, genres)
}

func Test_MovieDB_ActorsByMovie(t *testing.T) {
	mdb := getMovieDB()
	defer mdb.Close()

	expected := []*Person{
		&Person{Id: 13, Name: "Benicio del Toro"},
		&Person{Id: 7, Name: "Brad Pitt"},
		&Person{Id: 14, Name: "Jason Statham"},
		&Person{Id: 15, Name: "Vinnie Jones"},
	}
	actors, err := mdb.GetActorsByMovie("3")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, expected, actors)
}

func Test_MovieDB_DirectorsByMovie(t *testing.T) {
	mdb := getMovieDB()
	defer mdb.Close()

	expected := []*Person{
		&Person{Id: 12, Name: "Guy Ritchie"},
	}
	directors, err := mdb.GetDirectorsByMovie("3")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, expected, directors)
}

func Test_MovieDB_Genres(t *testing.T) {
	mdb := getMovieDB()
	defer mdb.Close()

	expected := "Thriller"
	genres, err := mdb.GetGenres()
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for i := range genres {
		if genres[i].Name == expected {
			found = true
		}
	}
	if !found {
		t.Errorf("Expected genres to contain [%v]", expected)
	}

	unexpected := "Girlie Movie"
	found = false
	for i := range genres {
		if genres[i].Name == unexpected {
			found = true
		}
	}
	if found {
		t.Errorf("Expected genres to not contain [%v]", unexpected)
	}
}

func Test_MovieDB_Languages(t *testing.T) {
	mdb := getMovieDB()
	defer mdb.Close()

	expected := "Englisch"
	languages, err := mdb.GetLanguages()
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for i := range languages {
		if languages[i].Name == expected {
			found = true
		}
	}
	if !found {
		t.Errorf("Expected languages to contain [%v]", expected)
	}

	unexpected := "Klingonisch"
	found = false
	for i := range languages {
		if languages[i].Name == unexpected {
			found = true
		}
	}
	if found {
		t.Errorf("Expected languages to not contain [%v]", unexpected)
	}

	expected = "T&#252;rkei"
	found = false
	for i := range languages {
		if languages[i].Country == expected {
			found = true
		}
	}
	if !found {
		t.Errorf("Expected languages to contain [%v]", expected)
	}

	expected = "Nab'ee Maya' Tzij"
	found = false
	for i := range languages {
		if languages[i].NativeName == expected {
			found = true
		}
	}
	if !found {
		t.Errorf("Expected languages to contain [%v]", expected)
	}
}

func Test_MovieDB_Person(t *testing.T) {
	mdb := getMovieDB()
	defer mdb.Close()

	expected := &Person{Id: 470, Name: "Roger Moore"}
	person, err := mdb.GetPerson("470")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, expected, person)
}

func Test_MovieDB_Actors(t *testing.T) {
	mdb := getMovieDB()
	defer mdb.Close()

	expected := "Sylvester Stallone"
	actors, err := mdb.GetActors()
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for i := range actors {
		if actors[i].Name == expected {
			found = true
		}
	}
	if !found {
		t.Errorf("Expected actors to contain [%v]", expected)
	}

	unexpected := "Steven Spielberg"
	found = false
	for i := range actors {
		if actors[i].Name == unexpected {
			found = true
		}
	}
	if found {
		t.Errorf("Expected actors to not contain [%v]", unexpected)
	}
}

func Test_MovieDB_Directors(t *testing.T) {
	mdb := getMovieDB()
	defer mdb.Close()

	expected := "Steven Spielberg"
	directors, err := mdb.GetDirectors()
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for i := range directors {
		if directors[i].Name == expected {
			found = true
		}
	}
	if !found {
		t.Errorf("Expected directors to contain [%v]", expected)
	}

	unexpected := "Bruce Willis"
	found = false
	for i := range directors {
		if directors[i].Name == unexpected {
			found = true
		}
	}
	if found {
		t.Errorf("Expected directors to not contain [%v]", unexpected)
	}
}

func Test_MovieDB_Statistics(t *testing.T) {
	mdb := getMovieDB()
	defer mdb.Close()

	// dates
	stats, err := mdb.GetStatistics()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 912, stats.Count)

	gz, err := time.Parse(time.UnixDate, "Sun Aug 01 00:13:37 UTC 1999")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, gz, stats.GroundZero)

	lu, err := time.Parse(time.UnixDate, "Wed Jan 01 17:11:36 UTC 2014")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, lu, stats.LastUpdate)

	// stats
	assert.Equal(t, 2, len(stats.Movies))
	assert.Equal(t, "DVD", stats.Movies[0].DiskType)
	assert.Equal(t, 1238, stats.Movies[0].Disks)
	assert.Equal(t, 154253, stats.Movies[0].Length)
	assert.Equal(t, 602, stats.Movies[0].Count)
	assert.Equal(t, "BluRay", stats.Movies[1].DiskType)
	assert.Equal(t, 493, stats.Movies[1].Disks)
	assert.Equal(t, 61691, stats.Movies[1].Length)
	assert.Equal(t, 310, stats.Movies[1].Count)

	assert.Equal(t, 4696, stats.Actors)
	assert.Equal(t, 578, stats.Directors)
	assert.Equal(t, 5244, stats.People)

	assert.Equal(t, "Bud Spencer", stats.TopActors[0].Name)
	assert.Equal(t, 338, stats.TopActors[1].Id)
	assert.Equal(t, 21, stats.TopActors[2].Count)

	assert.Equal(t, "Kenji Kamiyama", stats.TopDirectors[0].Name)
	assert.Equal(t, 127, stats.TopDirectors[1].Id)
	assert.Equal(t, 11, stats.TopDirectors[2].Count)

	assert.Equal(t, "Clint Eastwood", stats.TopActorsAndDirectors[0].Name)
	assert.Equal(t, 189, stats.TopActorsAndDirectors[1].Id)
	assert.Equal(t, 12, stats.TopActorsAndDirectors[2].Count)

	assert.Equal(t, "0", stats.Regions[0].Type)
	assert.Equal(t, 74, stats.Regions[1].Count)
	assert.Equal(t, "B", stats.Regions[4].Type)
	assert.Equal(t, 311, stats.Regions[4].Count)

	assert.Equal(t, "5", stats.Scores[0].Type)
	assert.Equal(t, 214, stats.Scores[1].Count)

	assert.Equal(t, "21", stats.Ratings[0].Type)
	assert.Equal(t, 132, stats.Ratings[1].Count)

	// numbers
	assert.True(t, stats.AvgMoviesPerDay > 0.15 && stats.AvgMoviesPerDay < 0.3)
	assert.True(t, stats.NewMoviesEstimate > 90)
	assert.Equal(t, 602, stats.DvdMovies)
	assert.Equal(t, 310, stats.BlurayMovies)
	assert.Equal(t, 1238, stats.DvdDisks)
	assert.Equal(t, 493, stats.BlurayDisks)
	assert.Equal(t, 215944, stats.TotalLength)
	assert.Equal(t, 236, stats.AvgLengthPerMovie)
	assert.Equal(t, 124, stats.AvgLengthPerDisk)
}
