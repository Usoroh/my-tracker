package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

var temp *template.Template

type Locations struct {
	Index []struct {
		ID        int
		Locations []string
	}
}

type Relations struct {
	Index []struct {
		ID             int
		DatesLocations map[string][]string
	}
}

var API struct {
	ID            int
	Artist        []Performer
	LocationsHtml Locations
	RelationsHtml Relations
}

type Performer struct {
	ID           int
	Image        string
	Name         string
	Members      []string
	CreationDate int
	FirstAlbum   string
}

func main() {

	//create each group, location strict -> then put data
	artists, _ := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	artistsBytes, _ := ioutil.ReadAll(artists.Body)
	artists.Body.Close()
	json.Unmarshal(artistsBytes, &API.Artist)

	locations, _ := http.Get("https://groupietrackers.herokuapp.com/api/locations")
	locationsBytes, _ := ioutil.ReadAll(locations.Body)
	locations.Body.Close()
	json.Unmarshal(locationsBytes, &API.LocationsHtml)

	relations, _ := http.Get("https://groupietrackers.herokuapp.com/api/relation")
	relationsBytes, _ := ioutil.ReadAll(relations.Body)
	relations.Body.Close()
	json.Unmarshal(relationsBytes, &API.RelationsHtml)

	//static data, css, js
	static := http.FileServer(http.Dir("public"))
	//secure, not access another files
	http.Handle("/public/", http.StripPrefix("/public/", static))

	http.HandleFunc("/", mainPage)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

func mainPage(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("templates/index.html")
	if r.Method == "GET" {
		if err != nil {
			fmt.Println("Internal Server Error")
			return
		}
		temp.Execute(w, API.Artist)
	}
}
