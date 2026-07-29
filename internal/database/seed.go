package database

import (
	"database/sql"
	"fmt"
)

type seedMovie struct {
	title    string
	year     int
	duration int
	genres   []int // genre ids
	actors   []int // actor ids
}

type seedActor struct {
	name      string
	birthDate string
}

// Seed inserts sample data when the database is empty. It is a no-op if
// any movies already exist, so restarting the server never duplicates data.
func Seed(db *sql.DB) error {
	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM movies`).Scan(&count); err != nil {
		return fmt.Errorf("seed: count movies: %w", err)
	}
	if count > 0 {
		return nil
	}

	genres := []string{
		"Action",
		"Sci-Fi",
		"Drama",
		"Comedy",
		"Thriller",
		"Romance",
		"Adventure",
		"Fantasy",
		"Crime",
		"Horror",
	}

	actors := []seedActor{
		{"Keanu Reeves", "1964-09-02"},
		{"Leonardo DiCaprio", "1974-11-11"},
		{"Tom Hardy", "1977-09-15"},
		{"Elliot Page", "1987-02-21"},
		{"Bruce Willis", "1955-03-19"},
		{"Tom Hanks", "1956-07-09"},
		{"Meryl Streep", "1949-06-22"},
		{"Anne Hathaway", "1982-11-12"},
		{"Kate Winslet", "1975-10-05"},
		{"Matthew McConaughey", "1969-11-04"},
		{"Christian Bale", "1974-01-30"},
		{"Samuel L. Jackson", "1948-12-21"},
		{"Emma Stone", "1988-11-06"},
		{"Ryan Gosling", "1980-11-12"},
		{"Julia Roberts", "1967-10-28"},
		{"Hugh Grant", "1960-09-09"},
		{"Amy Adams", "1974-08-20"},
		{"Daniel Craig", "1968-03-02"},
		{"Sandra Bullock", "1964-07-26"},
		{"Charlize Theron", "1975-08-07"},
		{"Tim Robbins", "1958-10-16"},
		{"Morgan Freeman", "1937-06-01"},
		{"Joseph Gordon-Levitt", "1981-02-17"},
		{"Marion Cotillard", "1975-09-30"},
		{"John Travolta", "1954-02-18"},
		{"Uma Thurman", "1970-04-29"},
		{"Jonah Hill", "1983-12-20"},
		{"Mark Ruffalo", "1967-11-22"},
		{"Robert Downey Jr.", "1965-04-04"},
		{"Scarlett Johansson", "1984-11-22"},
		{"Heath Ledger", "1979-04-04"},
		{"Aaron Eckhart", "1968-03-12"},
		{"Christopher Walken", "1943-03-31"},
		{"Chadwick Boseman", "1976-11-29"},
		{"Michael B. Jordan", "1987-02-09"},
		{"Timothee Chalamet", "1995-12-27"},
		{"Zendaya", "1996-09-01"},
		{"Al Pacino", "1940-04-25"},
		{"Russell Crowe", "1964-04-07"},
		{"Cillian Murphy", "1976-05-25"},
		{"Jodie Foster", "1962-11-19"},
		{"Anthony Hopkins", "1937-12-31"},
		{"Matt Damon", "1970-10-08"},
		{"Robin Williams", "1951-07-21"},
		{"Miles Teller", "1987-02-20"},
		{"J.K. Simmons", "1955-01-09"},
		{"Ralph Fiennes", "1962-12-22"},
		{"Hugh Jackman", "1968-10-12"},
		{"Kenneth Branagh", "1960-12-10"},
		{"Michelle Yeoh", "1962-08-06"},
		{"Ke Huy Quan", "1971-08-20"},
	}

	movies := []seedMovie{
		{"The Matrix", 1999, 136, []int{1, 2}, []int{1}},
		{"Inception", 2010, 148, []int{1, 2, 5}, []int{2, 4, 23, 24}},
		{"Die Hard", 1988, 132, []int{1, 5}, []int{5}},
		{"Forrest Gump", 1994, 142, []int{3, 4, 6}, []int{6}},
		{"The Devil Wears Prada", 2006, 109, []int{4, 3}, []int{7, 8}},
		{"Mamma Mia!", 2008, 108, []int{4, 6}, []int{7}},
		{"The Iron Lady", 2011, 105, []int{3}, []int{7}},
		{"Titanic", 1997, 195, []int{3, 6}, []int{2, 9}},
		{"Interstellar", 2014, 169, []int{2, 3}, []int{8, 10}},
		{"The Dark Knight", 2008, 152, []int{1, 5, 9}, []int{11, 22, 31, 32}},
		{"Pulp Fiction", 1994, 154, []int{9, 5, 3}, []int{25, 12, 26, 5}},
		{"The Shawshank Redemption", 1994, 142, []int{3}, []int{21, 22}},
		{"Catch Me If You Can", 2002, 141, []int{3, 4, 9}, []int{2, 6, 33}},
		{"The Wolf of Wall Street", 2013, 180, []int{4, 3, 9}, []int{2, 27}},
		{"La La Land", 2016, 128, []int{6, 4, 3}, []int{13, 14}},
		{"Mad Max: Fury Road", 2015, 120, []int{1, 2, 7}, []int{3, 20}},
		{"Notting Hill", 1999, 124, []int{6, 4}, []int{15, 16}},
		{"Speed", 1994, 116, []int{1, 5}, []int{1, 19}},
		{"Arrival", 2016, 116, []int{2, 3}, []int{17}},
		{"Knives Out", 2019, 130, []int{4, 5, 9}, []int{18}},
		{"John Wick", 2014, 101, []int{1, 5}, []int{1}},
		{"Black Panther", 2018, 134, []int{1, 2, 7}, []int{34, 35}},
		{"Dune", 2021, 155, []int{2, 7}, []int{36, 37}},
		{"Iron Man", 2008, 126, []int{1, 2, 7}, []int{29}},
		{"The Avengers", 2012, 143, []int{1, 2, 7}, []int{29, 30, 28}},
		{"The Godfather", 1972, 175, []int{9, 3}, []int{38}},
		{"Gladiator", 2000, 155, []int{1, 3, 7}, []int{39}},
		{"Oppenheimer", 2023, 180, []int{3, 5}, []int{40, 29}},
		{"The Silence of the Lambs", 1991, 118, []int{5, 9, 10}, []int{41, 42}},
		{"Good Will Hunting", 1997, 126, []int{3, 6}, []int{43, 44}},
		{"Whiplash", 2014, 106, []int{3, 5}, []int{45, 46}},
		{"The Grand Budapest Hotel", 2014, 99, []int{4, 3, 7}, []int{47}},
		{"The Prestige", 2006, 130, []int{2, 5, 3}, []int{11, 48}},
		{"Dunkirk", 2017, 106, []int{3, 7, 1}, []int{49, 40}},
		{"Everything Everywhere All at Once", 2022, 139, []int{2, 4, 3}, []int{50, 51}},
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("seed: begin transaction: %w", err)
	}
	defer tx.Rollback() //this is for the 3 loops below, in successful scenario, commit with close the transaction already.

	for _, name := range genres {
		if _, err := tx.Exec(`INSERT INTO genres (name) VALUES (?)`, name); err != nil {
			return fmt.Errorf("seed: insert genre %q: %w", name, err)
		}
	}
	for _, a := range actors {
		if _, err := tx.Exec(`INSERT INTO actors (name, birth_date) VALUES (?, ?)`, a.name, a.birthDate); err != nil {
			return fmt.Errorf("seed: insert actor %q: %w", a.name, err)
		}
	}
	for _, m := range movies {
		res, err := tx.Exec(
			`INSERT INTO movies (title, release_year, duration) VALUES (?, ?, ?)`,
			m.title, m.year, m.duration,
		)
		if err != nil {
			return fmt.Errorf("seed: insert movie %q: %w", m.title, err)
		}
		movieID, err := res.LastInsertId()
		if err != nil {
			return fmt.Errorf("seed: movie id for %q: %w", m.title, err)
		}
		for _, gid := range m.genres {
			if _, err := tx.Exec(`INSERT INTO movie_genre (movie_id, genre_id) VALUES (?, ?)`, movieID, gid); err != nil {
				return fmt.Errorf("seed: link movie %q to genre %d: %w", m.title, gid, err)
			}
		}
		for _, aid := range m.actors {
			if _, err := tx.Exec(`INSERT INTO movie_actor (movie_id, actor_id) VALUES (?, ?)`, movieID, aid); err != nil {
				return fmt.Errorf("seed: link movie %q to actor %d: %w", m.title, aid, err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("seed: commit: %w", err)
	}
	return nil
}
