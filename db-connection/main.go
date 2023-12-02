package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Scanner interface {
	Scan(dest ...any) error
}

type Person struct {
	ID    int64
	First string
	Last  string
	Email string
}

func main() {
	os.Remove("./sample.db")

	db, err := sql.Open("sqlite3", "./sample.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := createTables(db); err != nil {
		log.Fatal(err)
	}

	if err := insertRows(db); err != nil {
		log.Fatal(err)
	}

	p, err := queryDemo(db, "email_3")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(p)

	all, err := getAll(db)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("all results")
	for _, p = range all {
		fmt.Println(p)
	}
}

func createTables(db *sql.DB) error {
	qry := `CREATE TABLE test(
    id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
    fname TEXT,
    lname TEXT,
    email TEXT
  );`

	if _, err := db.Exec(qry); err != nil {
		return err
	}
	return nil
}

func insertRows(db *sql.DB) error {

	qry := `
    INSERT INTO test(fname, lname, email)
    VALUES(?, ?, ?);
  `
	ps, err := db.Prepare(qry)
	if err != nil {
		return err
	}

	defer ps.Close()

	for i := 0; i < 10; i++ {
		fname := fmt.Sprintf("fname_%d", i)
		lname := fmt.Sprintf("lname_%d", i)
		email := fmt.Sprintf("email_%d", i)

		if _, err := ps.Exec(fname, lname, email); err != nil {
			return err
		}
	}
	return nil
}

func queryDemo(db *sql.DB, email string) (p Person, err error) {
	row := db.QueryRow(`SELECT * FROM test WHERE email = ?`, email)
	if err != nil {
		return
	}
	err = scan(row, &p)
	return
}

func getAll(db *sql.DB) (all []Person, err error) {
	rows, err := db.Query(`SELECT * FROM test`)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		entity := Person{}
		if err = scan(rows, &entity); err != nil {
			return
		}
		all = append(all, entity)
	}

	err = rows.Err()
	return
}

func scan(r Scanner, entity *Person) error {
	return r.Scan(
		&entity.ID,
		&entity.First,
		&entity.Last,
		&entity.Email,
	)
}
