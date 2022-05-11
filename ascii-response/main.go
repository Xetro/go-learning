package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
)

var poke_types []string = []string{"fire", "water", "rock", "electric", "grass"}
var piece_chars []string = []string{"luffy", "shanks", "nami", "zoro", "sanji"}
var pokemons map[string]Character = make(map[string]Character)
var onepieces map[string]Character = make(map[string]Character)

type Character struct {
	name string
	term []byte
	html []byte
}

func (p *Character) loadImage() {
	d, err := os.ReadFile("./term/" + p.name)
	check(err)
	p.term = d
	d, err = os.ReadFile("./html/" + p.name)
	check(err)
	p.html = d
}

type App struct{}

func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var next http.Handler
	var head string

	head, r.URL.Path = shiftPath(r.URL.Path)

	switch head {
	case "pokemon":
		next = http.HandlerFunc(pokemon)
	case "onepiece":
		next = http.HandlerFunc(onepiece)
	default:
		next = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Nothing to see here.", http.StatusNotFound)
		})
	}

	next.ServeHTTP(w, r)
}

func pokemon(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = shiftPath(r.URL.Path)

	if head == "" {
		w.Write([]byte("Need pokemon type as part of URL path\n"))
		return
	}

	p, ok := pokemons[head]
	if !ok {
		w.Write([]byte("This pokemon type is missing :(\n"))
		return
	}
	userAgent := r.UserAgent()
	if userAgent[:4] == "curl" {
		w.Write(p.term)
	} else {
		w.Write(p.html)
	}
}

func onepiece(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = shiftPath(r.URL.Path)

	if head == "" {
		w.Write([]byte("Need character name as part of URL path\n"))
		return
	}

	p, ok := onepieces[head]
	if !ok {
		w.Write([]byte("This character is missing :(\n"))
		return
	}

	userAgent := r.UserAgent()
	if userAgent[:4] == "curl" {
		w.Write(p.term)
	} else {
		w.Write(p.html)
	}
}

func shiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	for _, v := range poke_types {
		p := Character{v, make([]byte, 128), make([]byte, 128)}
		p.loadImage()
		pokemons[v] = p
	}

	for _, v := range piece_chars {
		p := Character{v, make([]byte, 128), make([]byte, 128)}
		p.loadImage()
		onepieces[v] = p
	}
	app := &App{}
	fmt.Println("Server started on port 8080")
	err := http.ListenAndServe(":8080", app)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Println("Server closed")
	} else if err != nil {
		fmt.Println("error starting server: ", err)
	}
}
