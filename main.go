package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	// "strconv"
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

type Response struct {
	Results []struct {
		Geometry struct {
			Location struct {
				Lat float64
				Lng float64
			}
		}
	}
}

func main() {
	placeMarkers("lausanne-switzerland")

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
	http.HandleFunc("/artist", getArtist)
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

func getArtist(w http.ResponseWriter, r *http.Request) {
	temp, _ := template.ParseFiles("templates/test.html")
	if r.Method == "GET" {
		fmt.Println("here")
		ID, _ := strconv.Atoi(r.FormValue("uid"))
		API.ID = ID - 1
		fmt.Println(API.ID)
		temp.Execute(w, API)
	}
}

func placeMarkers(a string) {
	fmt.Println("here")
	// c, err := maps.NewClient(maps.WithAPIKey("AIzaSyDji8r-zQbC7DIfHWpPaTUX0uwtFGT6_eo"))
	// if err != nil {
	// 	log.Fatalf("fatal error: %s", err)
	// }

	// r := &maps.GeocodingRequest{
	// 	Address: a,
	// }

	// coor, err := c.Geocode(context.Background(), r)
	// if err != nil {
	// 	log.Fatalf("fatal erro: %s", err)
	// }

	// fmt.Println(coor)
	safeAddr := url.QueryEscape(a)
	apiKey := "AIzaSyDji8r-zQbC7DIfHWpPaTUX0uwtFGT6_eo"
	fullURL := fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?address=%s&key=%s", safeAddr, apiKey)
	resp, err := http.Get(fullURL)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	var res Response

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		log.Println(err)
	}

	lat := res.Results[0].Geometry.Location.Lat
	lng := res.Results[0].Geometry.Location.Lng

	fmt.Println(lat)
	fmt.Println(lng)
}
