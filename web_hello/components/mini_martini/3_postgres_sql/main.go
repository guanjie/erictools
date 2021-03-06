package main

import (
	"database/sql"
	"github.com/codegangsta/martini"
	_ "github.com/lib/pq"
	"github.com/martini-contrib/render"
	"net/http"
)

type Book struct {
	Title       string
	Author      string
	Description string
}

func SetupDB() *sql.DB {
	// Adding sslmode=disable to disable the panics
	db, err := sql.Open("postgres", "dbname=lesson4 sslmode=disable")
	PanicIf(err)
	return db
}

func PanicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	m := martini.Classic()
	// Directly mapping the sql data into the database
	m.Map(SetupDB())
	m.Use(render.Renderer(render.Options{
		Layout: "layout",
	}))

	m.Get("/", ShowBooks)
	m.Run()
}

func ShowBooks(ren render.Render, r *http.Request, db *sql.DB) {
	rows, err := db.Query("SELECT title, author, description FROM books")
	PanicIf(err)
	defer rows.Close()

	// Get the books slice
	books := []Book{}
	for rows.Next() {
		b := Book{}
		err := rows.Scan(&b.Title, &b.Author, &b.Description)
		PanicIf(err)
		books = append(books, b)
	}

	// Put books slice into the books template
	ren.HTML(200, "books", books)
}
