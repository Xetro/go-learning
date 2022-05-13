package main

import (
	"basic-api/controllers"
	"log"
	"net/http"
)

func main() {
	api := controllers.NewAPI()

	if err := http.ListenAndServe(":8080", api); err != nil {
		log.Println(err)
	}

}
