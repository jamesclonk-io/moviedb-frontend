package moviedb

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"
)

type Statistics struct {
	GroundZero            time.Time          `json:"ground_zero" xml:"ground_zero"`
	LastUpdate            time.Time          `json:"last_update" xml:"last_update"`
	Count                 int                `json:"count" xml:"count"`
	Movies                []*MovieType       `json:"movie_types" xml:"movie_types"`
	Actors                int                `json:"actors" xml:"actors"`
	Directors             int                `json:"directors" xml:"directors"`
	People                int                `json:"people_total" xml:"people_total"`
	TopActors             []*PersonWithCount `json:"top5_actors" xml:"top5_actors"`
	TopDirectors          []*PersonWithCount `json:"top5_directors" xml:"top5_directors"`
	TopActorsAndDirectors []*PersonWithCount `json:"top5_actors_and_directors" xml:"top5_actors_and_directors"`
	Regions               []*TypeCount       `json:"regions" xml:"regions"`
	Scores                []*TypeCount       `json:"scores" xml:"scores"`
	Ratings               []*TypeCount       `json:"ratings" xml:"ratings"`
	AvgMoviesPerDay       float64
	NewMoviesEstimate     float64
	DvdMovies             int
	BlurayMovies          int
	DvdDisks              int
	BlurayDisks           int
	TotalLength           int
	AvgLengthPerMovie     int
	AvgLengthPerDisk      int
}

type MovieType struct {
	DiskType string `json:"type" xml:"type"`
	Disks    int    `json:"disks" xml:"disks"`
	Length   int    `json:"length" xml:"length"`
	Count    int    `json:"count" xml:"count"`
}

type PersonWithCount struct {
	Id    int    `json:"id" xml:"id,attr"`
	Name  string `json:"name" xml:"name"`
	Count int    `json:"count" xml:"count"`
}

type TypeCount struct {
	Type  string `json:"type" xml:"type"`
	Count int    `json:"count" xml:"count"`
}

type Movie struct {
	Id          int            `json:"id" xml:"id,attr"`
	Title       string         `json:"title" xml:"title"`
	Alttitle    sql.NullString `json:"alttitle" xml:"alttitle"`
	Year        int            `json:"year" xml:"year"`
	Description string         `json:"description" xml:"description"`
	Format      string         `json:"format" xml:"format"`
	Length      int            `json:"length" xml:"length"`
	Region      string         `json:"region" xml:"region"`
	Rating      int            `json:"rating" xml:"rating"`
	Disks       int            `json:"disks" xml:"disks"`
	Score       int            `json:"score" xml:"year"`
	Picture     string         `json:"picture" xml:"picture"`
	Type        string         `json:"type" xml:"type"`
	Languages   []*Language    `json:"languages" xml:"languages"`
	Genres      []*Genre       `json:"genres" xml:"genres"`
	Actors      []*Person      `json:"actors" xml:"actors"`
	Directors   []*Person      `json:"directors" xml:"directors"`
}

func (m *Movie) String() string {
	return fmt.Sprintf("[%d] %s (%d)", m.Id, m.Title, m.Year)
}

type Language struct {
	Id         int    `json:"id" xml:"id,attr"`
	Name       string `json:"name" xml:"name"`
	Country    string `json:"country" xml:"country"`
	NativeName string `json:"native_name" xml:"native_name"`
}

type Genre struct {
	Id   int    `json:"id" xml:"id,attr"`
	Name string `json:"name" xml:"name"`
}

type Person struct {
	Id   int    `json:"id" xml:"id,attr"`
	Name string `json:"name" xml:"name"`
}

type MovieListing struct {
	Id     int    `json:"id" xml:"id,attr"`
	Title  string `json:"title" xml:"title"`
	Year   int    `json:"year" xml:"year"`
	Score  int    `json:"score" xml:"year"`
	Rating int    `json:"rating" xml:"rating"`
}

func (m *MovieListing) String() string {
	return fmt.Sprintf("[%d] %s (%d)", m.Id, m.Title, m.Year)
}

type MovieListingOptions struct {
	Sort  []Sort
	Query []Query
}

func ParseMovieListingOptions(req *http.Request) MovieListingOptions {
	var options MovieListingOptions
	q := req.URL.Query()

	var sortlist []Sort
	sort := q["sort"]
	by := q["by"]
	if len(sort) == len(by) {
		for i := range sort {
			sortlist = append(sortlist, NewSort(sort[i], by[i]))
		}
		options.Sort = sortlist
	}

	var querylist []Query
	query := q["query"]
	value := q["value"]
	if len(query) == len(value) {
		for i := range query {
			querylist = append(querylist, NewQuery(query[i], value[i]))
		}
		options.Query = querylist
	}

	return options
}

type Sort interface {
	Field() string
	Order() string
}

type sort struct {
	field string
	order string
}

func (s *sort) Field() string {
	return s.field
}

func (s *sort) Order() string {
	return s.order
}

func (s *sort) SetOrderBy(field string, order string) {
	switch {
	case field == "title" || field == "year" ||
		field == "score" || field == "rating" ||
		field == "format" || field == "disk_region" ||
		field == "length" || field == "disks" || field == "disk_type":
		s.field = field
	default:
		s.field = "id"
	}

	switch order {
	case "desc":
		s.order = "desc"
	default:
		s.order = "asc"
	}
}

func NewSort(field string, order string) Sort {
	var s sort
	s.SetOrderBy(field, order)
	return &s
}

type Query interface {
	Query() string
	Value() string
}

type query struct {
	query string
	value string
}

func (q *query) Query() string {
	return q.query
}

func (q *query) Value() string {
	return q.value
}

func (q *query) SetQuery(query string, value string) {
	switch {
	case query == "title" || query == "year" ||
		query == "score" || query == "rating" ||
		query == "disk_region" || query == "disk_type" ||
		query == "language" || query == "genre" ||
		query == "format" || query == "disks" ||
		query == "char" || query == "search" ||
		query == "actor" || query == "director" || query == "length":
		q.query = query
	default:
		q.query = "id"
	}
	q.value = value
}

func NewQuery(key string, value string) Query {
	var q query
	q.SetQuery(key, value)
	return &q
}
