package moviedb

import (
	"database/sql"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/jamesclonk-io/moviedb-backend/modules/database"
)

type MovieDB interface {
	GetMovie(id string) (*Movie, error)
	DeleteMovie(id string) (int64, error)
	AddMovie(*Movie) error
	SaveMovie(*Movie) error
	GetMovieListings(...MovieListingOptions) ([]*MovieListing, error)
	GetLanguagesByMovie(id string) ([]*Language, error)
	GetGenresByMovie(id string) ([]*Genre, error)
	GetActorsByMovie(id string) ([]*Person, error)
	GetDirectorsByMovie(id string) ([]*Person, error)
	GetLanguages() ([]*Language, error)
	GetGenres() ([]*Genre, error)
	GetPerson(id string) (*Person, error)
	GetActors() ([]*Person, error)
	GetDirectors() ([]*Person, error)
	GetStatistics() (*Statistics, error)
}

type movieDB struct {
	*sql.DB
	DatabaseType string
}

func NewMovieDB(adapter *database.Adapter) MovieDB {
	return &movieDB{adapter.Database, adapter.Type}
}

func (mdb *movieDB) GetLanguagesByMovie(id string) ([]*Language, error) {
	stmt, err := mdb.Prepare(`
		select ml.id, ml.name, ml.country, ml.native_name
		from movie_language ml
		join movie_link_language mll on (mll.language_id = ml.id)
		where mll.movie_id = $1
		order by ml.name asc`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(id)
	if err != nil {
		return nil, err
	}

	ls := []*Language{}
	for rows.Next() {
		var l Language
		if err := rows.Scan(&l.Id, &l.Name, &l.Country, &l.NativeName); err != nil {
			return nil, err
		}
		ls = append(ls, &l)
	}
	return ls, nil
}

func (mdb *movieDB) GetGenresByMovie(id string) ([]*Genre, error) {
	stmt, err := mdb.Prepare(`
		select mg.id, mg.name 
		from movie_genre mg
		join movie_link_genre mlg on (mlg.genre_id = mg.id)
		where mlg.movie_id = $1
		order by mg.name asc`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(id)
	if err != nil {
		return nil, err
	}

	gs := []*Genre{}
	for rows.Next() {
		var g Genre
		if err := rows.Scan(&g.Id, &g.Name); err != nil {
			return nil, err
		}
		gs = append(gs, &g)
	}
	return gs, nil
}

func (mdb *movieDB) GetActorsByMovie(id string) ([]*Person, error) {
	stmt, err := mdb.Prepare(`
		select distinct mp.id, mp.name 
		from movie_people mp
		join movie_link_actor mla on (mla.person_id = mp.id)
		where mla.movie_id = $1
		order by mp.name asc`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(id)
	if err != nil {
		return nil, err
	}

	ps := []*Person{}
	for rows.Next() {
		var p Person
		if err := rows.Scan(&p.Id, &p.Name); err != nil {
			return nil, err
		}
		ps = append(ps, &p)
	}
	return ps, nil
}

func (mdb *movieDB) GetDirectorsByMovie(id string) ([]*Person, error) {
	stmt, err := mdb.Prepare(`
		select distinct mp.id, mp.name 
		from movie_people mp
		join movie_link_director mld on (mld.person_id = mp.id)
		where mld.movie_id = $1
		order by mp.name asc`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(id)
	if err != nil {
		return nil, err
	}

	ps := []*Person{}
	for rows.Next() {
		var p Person
		if err := rows.Scan(&p.Id, &p.Name); err != nil {
			return nil, err
		}
		ps = append(ps, &p)
	}
	return ps, nil
}

func (mdb *movieDB) GetLanguages() ([]*Language, error) {
	rows, err := mdb.Query(`select ml.id, ml.name, ml.country, ml.native_name 
		from movie_language ml order by ml.name asc`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ls := []*Language{}
	for rows.Next() {
		var l Language
		if err := rows.Scan(&l.Id, &l.Name, &l.Country, &l.NativeName); err != nil {
			return nil, err
		}
		ls = append(ls, &l)
	}
	return ls, nil
}

func (mdb *movieDB) GetGenres() ([]*Genre, error) {
	rows, err := mdb.Query(`select mg.id, mg.name from movie_genre mg order by mg.name asc`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	gs := []*Genre{}
	for rows.Next() {
		var g Genre
		if err := rows.Scan(&g.Id, &g.Name); err != nil {
			return nil, err
		}
		gs = append(gs, &g)
	}
	return gs, nil
}

func (mdb *movieDB) GetPerson(id string) (*Person, error) {
	stmt, err := mdb.Prepare(`select id, name from movie_people where id = $1`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	p := &Person{}
	if err := stmt.QueryRow(id).Scan(&p.Id, &p.Name); err != nil {
		return nil, err
	}
	return p, nil
}

func (mdb *movieDB) GetActors() ([]*Person, error) {
	rows, err := mdb.Query(`select distinct mp.id, mp.name 
		from movie_people mp join movie_link_actor mla on (mla.person_id = mp.id) order by mp.name asc`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ps := []*Person{}
	for rows.Next() {
		var p Person
		if err := rows.Scan(&p.Id, &p.Name); err != nil {
			return nil, err
		}
		ps = append(ps, &p)
	}
	return ps, nil
}

func (mdb *movieDB) GetDirectors() ([]*Person, error) {
	rows, err := mdb.Query(`select distinct mp.id, mp.name 
		from movie_people mp join movie_link_director mld on (mld.person_id = mp.id) order by mp.name asc`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ps := []*Person{}
	for rows.Next() {
		var p Person
		if err := rows.Scan(&p.Id, &p.Name); err != nil {
			return nil, err
		}
		ps = append(ps, &p)
	}
	return ps, nil
}

func (mdb *movieDB) GetStatistics() (*Statistics, error) {
	stats := Statistics{}

	// -----------------------------------------------------------------
	// dates & count
	rows0, err := mdb.Query(`select id, date from movie_dbdate`)
	if err != nil {
		return nil, err
	}
	defer rows0.Close()

	for rows0.Next() {
		var id int
		var date time.Time
		if err := rows0.Scan(&id, &date); err != nil {
			return nil, err
		}
		if id == 1 {
			stats.GroundZero = date
		} else {
			stats.LastUpdate = date
		}
	}

	if err := mdb.QueryRow(`select count(*) from movie_movie`).Scan(&stats.Count); err != nil {
		return nil, err
	}

	// -----------------------------------------------------------------
	// general statistics
	rows1, err := mdb.Query(`select disk_type, sum(disks), sum(length), count(*) 
		from movie_movie group by disk_type order by disk_type desc`)
	if err != nil {
		return nil, err
	}
	defer rows1.Close()

	mcs := []*MovieType{}
	for rows1.Next() {
		var mc MovieType
		if err := rows1.Scan(&mc.DiskType, &mc.Disks, &mc.Length, &mc.Count); err != nil {
			return nil, err
		}
		mcs = append(mcs, &mc)
	}
	stats.Movies = mcs

	if err := mdb.QueryRow(`select sum(actors) as actors, sum(directors) as directors, sum(people) as people from (		
			select 0 as actors, 0 as directors, count(*) as people from movie_people
			union
			select count(*) as actors, 0 as directors, 0 as people 
				from (select distinct mp.id, mp.name from movie_people mp join movie_link_actor mla on (mla.person_id = mp.id)) actors
			union
			select 0 as actors, count(*) as directors, 0 as people 
				from (select distinct mp.id, mp.name from movie_people mp join movie_link_director mld on (mld.person_id = mp.id)) directors
		) final`).Scan(&stats.Actors, &stats.Directors, &stats.People); err != nil {
		return nil, err
	}

	// -----------------------------------------------------------------
	// actor, director and actor&director statistics
	rows2, err := mdb.Query(`select mp.id, mp.name, count(*) 
		from movie_people mp join movie_link_actor mla on (mla.person_id = mp.id) 
		group by mp.id, mp.name order by 3 desc, 2 asc, 1 asc limit 5`)
	if err != nil {
		return nil, err
	}
	defer rows2.Close()

	as := []*PersonWithCount{}
	for rows2.Next() {
		var p PersonWithCount
		if err := rows2.Scan(&p.Id, &p.Name, &p.Count); err != nil {
			return nil, err
		}
		as = append(as, &p)
	}
	stats.TopActors = as

	rows3, err := mdb.Query(`select mp.id, mp.name, count(*) 
		from movie_people mp join movie_link_director mld on (mld.person_id = mp.id) 
		group by mp.id, mp.name order by 3 desc, 2 asc, 1 asc limit 5`)
	if err != nil {
		return nil, err
	}
	defer rows3.Close()

	ds := []*PersonWithCount{}
	for rows3.Next() {
		var p PersonWithCount
		if err := rows3.Scan(&p.Id, &p.Name, &p.Count); err != nil {
			return nil, err
		}
		ds = append(ds, &p)
	}
	stats.TopDirectors = ds

	rows4, err := mdb.Query(`select mp.id, mp.name, 
			(select count(*) from movie_link_actor where person_id = mp.id)
            + (select count(*) from movie_link_director where person_id = mp.id) as count
        from movie_people mp
            join movie_link_actor mla on (mla.person_id = mp.id)
			join movie_link_director mld on (mld.person_id = mp.id)
        group by mp.id, mp.name order by 3 desc, 2 asc, 1 asc limit 5`)
	if err != nil {
		return nil, err
	}
	defer rows4.Close()

	ads := []*PersonWithCount{}
	for rows4.Next() {
		var p PersonWithCount
		if err := rows4.Scan(&p.Id, &p.Name, &p.Count); err != nil {
			return nil, err
		}
		ads = append(ads, &p)
	}
	stats.TopActorsAndDirectors = ads

	// -----------------------------------------------------------------
	// region, score and rating statistics
	rows5, err := mdb.Query(`select disk_region, count(*) 
		from movie_movie group by disk_region order by disk_region asc`)
	if err != nil {
		return nil, err
	}
	defer rows5.Close()

	rs := []*TypeCount{}
	for rows5.Next() {
		var t TypeCount
		if err := rows5.Scan(&t.Type, &t.Count); err != nil {
			return nil, err
		}
		rs = append(rs, &t)
	}
	stats.Regions = rs

	rows6, err := mdb.Query(`select score, count(*) from movie_movie group by score order by score desc`)
	if err != nil {
		return nil, err
	}
	defer rows6.Close()

	ss := []*TypeCount{}
	for rows6.Next() {
		var t TypeCount
		if err := rows6.Scan(&t.Type, &t.Count); err != nil {
			return nil, err
		}
		ss = append(ss, &t)
	}
	stats.Scores = ss

	rows7, err := mdb.Query(`select rating, count(*) from movie_movie group by rating order by rating desc`)
	if err != nil {
		return nil, err
	}
	defer rows7.Close()

	rss := []*TypeCount{}
	for rows7.Next() {
		var t TypeCount
		if err := rows7.Scan(&t.Type, &t.Count); err != nil {
			return nil, err
		}
		rss = append(rss, &t)
	}
	stats.Ratings = rss

	// -----------------------------------------------------------------
	// numbers

	// calculating avg movies per day and currently estimated total new movies
	avgMoviesPerDay := (float64(stats.Count) / (stats.LastUpdate.Sub(stats.GroundZero).Hours() / 24))
	daysSinceLastUpdate := (time.Now().Sub(stats.LastUpdate).Hours() / 24)
	stats.NewMoviesEstimate = round((daysSinceLastUpdate * avgMoviesPerDay), 1)
	stats.AvgMoviesPerDay = round(avgMoviesPerDay, 2)

	// calculating runtimes
	stats.TotalLength = stats.Movies[0].Length + stats.Movies[1].Length
	stats.AvgLengthPerMovie = int(round(float64(stats.TotalLength/stats.Count), 0))
	stats.AvgLengthPerDisk = int(round(float64(stats.TotalLength/(stats.Movies[0].Disks+stats.Movies[1].Disks)), 0))

	// counts
	stats.DvdMovies = stats.Movies[0].Count
	stats.BlurayMovies = stats.Movies[1].Count
	stats.DvdDisks = stats.Movies[0].Disks
	stats.BlurayDisks = stats.Movies[1].Disks

	return &stats, nil
}

func (mdb *movieDB) GetMovie(id string) (*Movie, error) {
	stmt, err := mdb.Prepare(`select id, title, alttitle, year, description, format, length, 
		disk_region, rating, disks, score, picture, disk_type from movie_movie where id = $1`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var m Movie
	if err := stmt.QueryRow(id).Scan(&m.Id, &m.Title, &m.Alttitle, &m.Year, &m.Description, &m.Format, &m.Length,
		&m.Region, &m.Rating, &m.Disks, &m.Score, &m.Picture, &m.Type); err != nil {
		return nil, err
	}

	languages, err := mdb.GetLanguagesByMovie(id)
	if err != nil {
		return nil, err
	}
	m.Languages = languages

	genres, err := mdb.GetGenresByMovie(id)
	if err != nil {
		return nil, err
	}
	m.Genres = genres

	actors, err := mdb.GetActorsByMovie(id)
	if err != nil {
		return nil, err
	}
	m.Actors = actors

	directors, err := mdb.GetDirectorsByMovie(id)
	if err != nil {
		return nil, err
	}
	m.Directors = directors

	return &m, nil
}

func (mdb *movieDB) AddMovie(movie *Movie) error {
	// first get next/new movie_id
	var newId int
	row := mdb.QueryRow(`select max(id)+1 from movie_movie`)
	if err := row.Scan(&newId); err != nil {
		return err
	}
	if newId < 900 {
		return errors.New(fmt.Sprintf("new movie_id impossible! [%v]", newId))
	}

	// set movie_id and call SaveMovie, which will check if movie already exists (it won't) and then inserts it
	movie.Id = newId
	return mdb.SaveMovie(movie)
}

func (mdb *movieDB) SaveMovie(movie *Movie) error {
	tx, err := mdb.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// check if movie already exists
	var exists string
	rows, err := tx.Query("select 'yes' from movie_movie where id = $1", movie.Id)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		if err := rows.Scan(&exists); err != nil {
			return err
		}
	}

	if exists == "yes" {
		// update movie
		stmt, err := tx.Prepare(`UPDATE movie_movie
			set title = $1,
			alttitle = $2,
			year = $3,
			description = $4,
			format = $5,
			length = $6,
			disk_region = $7,
			rating = $8,
			disks = $9,
			score = $10,
			picture = $11,
			disk_type = $12
			where id = $13
			`)
		if err != nil {
			return err
		}
		defer stmt.Close()

		if _, err := stmt.Exec(movie.Title, movie.Alttitle, movie.Year, movie.Description, movie.Format,
			movie.Length, movie.Region, movie.Rating, movie.Disks, movie.Score, movie.Picture, movie.Type, movie.Id); err != nil {
			return err
		}

	} else {
		// insert movie
		stmt, err := tx.Prepare(`INSERT INTO movie_movie
			(id, title, alttitle, year, description, format, length, disk_region, rating, disks, score, picture, disk_type) 
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`)
		if err != nil {
			return err
		}
		defer stmt.Close()

		if _, err := stmt.Exec(movie.Id, movie.Title, movie.Alttitle, movie.Year, movie.Description, movie.Format,
			movie.Length, movie.Region, movie.Rating, movie.Disks, movie.Score, movie.Picture, movie.Type); err != nil {
			return err
		}
	}

	if err := saveLanguages(tx, movie); err != nil {
		return err
	}

	if err := saveGenres(tx, movie); err != nil {
		return err
	}

	if err := saveActors(tx, movie); err != nil {
		return err
	}

	if err := saveDirectors(tx, movie); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func saveLanguages(tx *sql.Tx, movie *Movie) error {
	for idx, language := range movie.Languages {
		// check if language already exists
		var exists string
		var id int
		rows, err := tx.Query("select 'yes', id from movie_language where name = $1", language.Name)
		if err != nil {
			return err
		}
		defer rows.Close()
		if rows.Next() {
			if err := rows.Scan(&exists, &id); err != nil {
				return err
			}
		}

		// insert to language table
		if exists == "yes" {
			// update only
			movie.Languages[idx].Id = id

		} else {
			// insert
			// first get next/new language_id
			var newId int
			row := tx.QueryRow(`select max(id)+1 from movie_language`)
			if err := row.Scan(&newId); err != nil {
				return err
			}
			if newId < 10 {
				return errors.New(fmt.Sprintf("new language_id impossible! [%v]", newId))
			}

			stmt, err := tx.Prepare(`INSERT INTO movie_language (id, name, country, native_name) VALUES ($1,$2,$3,$4)`)
			if err != nil {
				return err
			}
			defer stmt.Close()

			if _, err := stmt.Exec(newId, language.Name, language.Country, language.NativeName); err != nil {
				return err
			}
			movie.Languages[idx].Id = newId
		}

		// check if language link already exists
		rows, err = tx.Query("select 'yes' from movie_link_language where movie_id = $1 and language_id = $2",
			movie.Id, movie.Languages[idx].Id)
		if err != nil {
			return err
		}
		defer rows.Close()
		if !rows.Next() {
			// insert to link table
			stmt, err := tx.Prepare(`INSERT INTO movie_link_language (movie_id, language_id) VALUES ($1,$2)`)
			if err != nil {
				return err
			}
			defer stmt.Close()

			if _, err := stmt.Exec(movie.Id, movie.Languages[idx].Id); err != nil {
				return err
			}
		}
	}

	return nil
}

func saveGenres(tx *sql.Tx, movie *Movie) error {
	for idx, genre := range movie.Genres {
		// check if genre already exists
		var exists string
		var id int
		rows, err := tx.Query("select 'yes', id from movie_genre where name = $1", genre.Name)
		if err != nil {
			return err
		}
		defer rows.Close()
		if rows.Next() {
			if err := rows.Scan(&exists, &id); err != nil {
				return err
			}
		}

		// insert to genre table
		if exists == "yes" {
			// update only
			movie.Genres[idx].Id = id

		} else {
			// insert
			// first get next/new genre_id
			var newId int
			row := tx.QueryRow(`select max(id)+1 from movie_genre`)
			if err := row.Scan(&newId); err != nil {
				return err
			}
			if newId < 20 {
				return errors.New(fmt.Sprintf("new genre_id impossible! [%v]", newId))
			}

			stmt, err := tx.Prepare(`INSERT INTO movie_genre (id, name) VALUES ($1,$2)`)
			if err != nil {
				return err
			}
			defer stmt.Close()

			if _, err := stmt.Exec(newId, genre.Name); err != nil {
				return err
			}
			movie.Genres[idx].Id = newId
		}

		// check if language link already exists
		rows, err = tx.Query("select 'yes' from movie_link_genre where movie_id = $1 and genre_id = $2",
			movie.Id, movie.Genres[idx].Id)
		if err != nil {
			return err
		}
		defer rows.Close()
		if !rows.Next() {
			// insert to link table
			stmt, err := tx.Prepare(`INSERT INTO movie_link_genre (movie_id, genre_id) VALUES ($1,$2)`)
			if err != nil {
				return err
			}
			defer stmt.Close()

			if _, err := stmt.Exec(movie.Id, movie.Genres[idx].Id); err != nil {
				return err
			}
		}
	}

	return nil
}

func saveActors(tx *sql.Tx, movie *Movie) error {
	for idx, person := range movie.Actors {
		// check if person already exists
		var exists string
		var id int
		rows, err := tx.Query("select 'yes', id from movie_people where name = $1", person.Name)
		if err != nil {
			return err
		}
		defer rows.Close()
		if rows.Next() {
			if err := rows.Scan(&exists, &id); err != nil {
				return err
			}
		}

		// insert to people table
		if exists == "yes" {
			// update only
			movie.Actors[idx].Id = id

		} else {
			// insert
			// first get next/new person_id
			var newId int
			row := tx.QueryRow(`select max(id)+1 from movie_people`)
			if err := row.Scan(&newId); err != nil {
				return err
			}
			if newId < 5000 {
				return errors.New(fmt.Sprintf("new person_id impossible! [%v]", newId))
			}

			stmt, err := tx.Prepare(`INSERT INTO movie_people (id, name) VALUES ($1,$2)`)
			if err != nil {
				return err
			}
			defer stmt.Close()

			if _, err := stmt.Exec(newId, person.Name); err != nil {
				return err
			}
			movie.Actors[idx].Id = newId
		}

		// check if actor link already exists
		rows, err = tx.Query("select 'yes' from movie_link_actor where movie_id = $1 and person_id = $2",
			movie.Id, movie.Actors[idx].Id)
		if err != nil {
			return err
		}
		defer rows.Close()
		if !rows.Next() {
			// insert to link table
			stmt, err := tx.Prepare(`INSERT INTO movie_link_actor (movie_id, person_id) VALUES ($1,$2)`)
			if err != nil {
				return err
			}
			defer stmt.Close()

			if _, err := stmt.Exec(movie.Id, movie.Actors[idx].Id); err != nil {
				return err
			}
		}
	}

	return nil
}

func saveDirectors(tx *sql.Tx, movie *Movie) error {
	for idx, person := range movie.Directors {
		// check if person already exists
		var exists string
		var id int
		rows, err := tx.Query("select 'yes', id from movie_people where name = $1", person.Name)
		if err != nil {
			return err
		}
		defer rows.Close()
		if rows.Next() {
			if err := rows.Scan(&exists, &id); err != nil {
				return err
			}
		}

		// insert to people table
		if exists == "yes" {
			// update only
			movie.Directors[idx].Id = id

		} else {
			// insert
			// first get next/new person_id
			var newId int
			row := tx.QueryRow(`select max(id)+1 from movie_people`)
			if err := row.Scan(&newId); err != nil {
				return err
			}
			if newId < 5000 {
				return errors.New(fmt.Sprintf("new person_id impossible! [%v]", newId))
			}

			stmt, err := tx.Prepare(`INSERT INTO movie_people (id, name) VALUES ($1,$2)`)
			if err != nil {
				return err
			}
			defer stmt.Close()

			if _, err := stmt.Exec(newId, person.Name); err != nil {
				return err
			}
			movie.Directors[idx].Id = newId
		}

		// check if actor link already exists
		rows, err = tx.Query("select 'yes' from movie_link_director where movie_id = $1 and person_id = $2",
			movie.Id, movie.Directors[idx].Id)
		if err != nil {
			return err
		}
		defer rows.Close()
		if !rows.Next() {
			// insert to link table
			stmt, err := tx.Prepare(`INSERT INTO movie_link_director (movie_id, person_id) VALUES ($1,$2)`)
			if err != nil {
				return err
			}
			defer stmt.Close()

			if _, err := stmt.Exec(movie.Id, movie.Directors[idx].Id); err != nil {
				return err
			}
		}
	}

	return nil
}

func (mdb *movieDB) DeleteMovie(id string) (int64, error) {
	var rowsDeleted int64

	tx, err := mdb.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	sqls := []string{
		`delete from movie_movie where id = $1`,
		`delete from movie_link_actor where movie_id = $1`,
		`delete from movie_link_director where movie_id = $1`,
		`delete from movie_link_genre where movie_id = $1`,
		`delete from movie_link_language where movie_id = $1`,
	}

	for _, sql := range sqls {
		rows, err := deleteWithinTransaction(tx, id, sql)
		if err != nil {
			return 0, err
		}
		rowsDeleted = rowsDeleted + rows
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return rowsDeleted, nil
}

func deleteWithinTransaction(tx *sql.Tx, id string, sql string) (int64, error) {
	stmt, err := tx.Prepare(sql)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(id)
	if err != nil {
		return 0, err
	}

	rowsDeleted, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsDeleted, nil
}

func (mdb *movieDB) GetMovieListings(opt ...MovieListingOptions) ([]*MovieListing, error) {
	var options MovieListingOptions
	if len(opt) > 0 {
		options = opt[0]
	}

	sql := `select mm.id, mm.title, mm.year, mm.score, mm.rating from movie_movie mm `
	var params []interface{} // all bind variables in here
	paramCounter := 1

	if len(options.Query) > 0 {
		for _, query := range options.Query {
			switch {
			case query.Query() == "language":
				sql += fmt.Sprintf("join movie_link_language mll on (mll.movie_id = mm.id and mll.language_id = $%d) ", paramCounter)
				params = append(params, query.Value())
				paramCounter += 1
			case query.Query() == "genre":
				sql += fmt.Sprintf("join movie_link_genre mlg on (mlg.movie_id = mm.id and mlg.genre_id = $%d) ", paramCounter)
				params = append(params, query.Value())
				paramCounter += 1
			case query.Query() == "actor":
				sql += fmt.Sprintf("join movie_link_actor mla on (mla.movie_id = mm.id and mla.person_id = $%d) ", paramCounter)
				params = append(params, query.Value())
				paramCounter += 1
			case query.Query() == "director":
				sql += fmt.Sprintf("join movie_link_director mld on (mld.movie_id = mm.id and mld.person_id = $%d) ", paramCounter)
				params = append(params, query.Value())
				paramCounter += 1
			}
		}
		sql += "where 1 = 1 "
		for _, query := range options.Query {
			switch {
			case query.Query() == "char" && query.Value() == "num":
				sql += "and substr(title,1,1) in ('1','2','3','4','5','6','7','8','9','0') "
			case query.Query() == "char":
				sql += fmt.Sprintf("and upper(substr(title,1,1)) = upper($%d) ", paramCounter)
				params = append(params, query.Value())
				paramCounter += 1
			case query.Query() == "search":
				if mdb.DatabaseType == "postgres" {
					sql += fmt.Sprintf("and (title ilike $%d or alttitle ilike $%d or description ilike $%d) ",
						paramCounter, paramCounter, paramCounter)
				} else {
					sql += fmt.Sprintf("and (title like $%d or alttitle like $%d or description like $%d) ",
						paramCounter, paramCounter, paramCounter)
				}
				params = append(params, fmt.Sprintf("%%%s%%", query.Value()))
				paramCounter += 1
			case query.Query() != "language" &&
				query.Query() != "genre" &&
				query.Query() != "actor" &&
				query.Query() != "director":
				sql += fmt.Sprintf("and %s = $%d ", query.Query(), paramCounter)
				params = append(params, query.Value())
				paramCounter += 1
			}
		}
	}

	if len(options.Sort) > 0 {
		sql += "order by "
		for i, sort := range options.Sort {
			if i > 0 {
				sql += ", "
			}
			sql += fmt.Sprintf("%s %s", sort.Field(), sort.Order())
		}
	} else {
		sql += "order by title asc"
	}

	stmt, err := mdb.Prepare(sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Err() != nil {
		return nil, err
	}

	ms := []*MovieListing{}
	for rows.Next() {
		var m MovieListing
		if err := rows.Scan(&m.Id, &m.Title, &m.Year, &m.Score, &m.Rating); err != nil {
			return nil, err
		}
		ms = append(ms, &m)
	}
	return ms, nil
}

func round(val float64, prec int) float64 {
	var rounder float64
	pow := math.Pow(10, float64(prec))
	intermed := val * pow

	if intermed < 0.0 {
		intermed -= 0.5
	} else {
		intermed += 0.5
	}
	rounder = float64(int64(intermed))

	return rounder / float64(pow)
}
