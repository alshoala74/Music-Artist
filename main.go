package main 

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"io"
	"strconv"
	"html/template"
)
	
type Artist struct {
	ID int `json:"id"`
	Image  string `json:"image"`
	Name   string `json:"name"`
	Members []string `json:"members"`
	CreationDate int `json:"creationDate"`
	FirstAlbum string `json:"firstAlbum"`
	Locations []string `json:"locations"`
	ConcertDates []string `json:"concertDates"`
	Relations map[string][]string `json:"relations"`
}

type Locations struct {
	ID int `json:"id"`
	Location []string `json:"locations"`
	Dates string `json: "dates"`
}
type Dates struct {
	ID int `json:"id"`
	Dates []string `json: "dates"`
}
type Relation struct {
	ID int `json:"id"`
	DatesLocations map[string][]string `json: "datesLocations"`
}

var Respons []Artist

func main(){
	//get the data from the API
	respons,err:=http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		log.Fatal(err)
	}

	//read the body and save it as []byte
	responsData,err:=io.ReadAll(respons.Body)
		if err != nil {
			log.Fatal(err)
		}

	//convert []byte to []string
	json.Unmarshal(responsData,&Respons)


	for i,Art:=range Respons{
		id := Art.ID
		idstring:=strconv.Itoa(id)
		loca:=FetchLocations(idstring)
		date:=FetchDates(idstring)
		rel:=FetchRelation(idstring)
		Respons[i].Locations=loca.Location
		Respons[i].ConcertDates=date.Dates
		Respons[i].Relations=rel.DatesLocations

	 }


	 http.HandleFunc("/",HomeHandler)
	 http.HandleFunc("/artist",ArtistHandler)
	 http.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir("templates"))))
	 fmt.Println("server conected @ http://localhost:8080")
	 http.ListenAndServe(":8080", nil)
	
}

func FetchLocations(idstring string) Locations{
	url:="https://groupietrackers.herokuapp.com/api/locations/"+idstring
	respons,err:=http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	//read the body and save it as []byte
	responsData,err:=io.ReadAll(respons.Body)
		if err != nil {
			log.Fatal(err)
		}

	var Respons Locations	
	//convert []byte to []string
	json.Unmarshal(responsData,&Respons)

	return Respons
}


func FetchDates(idstring string) Dates{
	url:="https://groupietrackers.herokuapp.com/api/dates/"+idstring
	respons,err:=http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	//read the body and save it as []byte
	responsData,err:=io.ReadAll(respons.Body)
		if err != nil {
			log.Fatal(err)
		}
	
	var Respons Dates
	//convert []byte to []string
	json.Unmarshal(responsData,&Respons)

	return Respons
}

func FetchRelation(idstring string) Relation{
	url:="https://groupietrackers.herokuapp.com/api/relation/"+idstring
	respons,err:=http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	//read the body and save it as []byte
	responsData,err:=io.ReadAll(respons.Body)
		if err != nil {
			log.Fatal(err)
		}

	var Respons Relation
	//convert []byte to []string
	json.Unmarshal(responsData, &Respons)

	return Respons
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		http.ServeFile(w, r, "templates/404.html")
		return
	}
	tmpl, err := template.ParseFiles("templates/index.html") 
	if err != nil {    
		                         
		w.WriteHeader(http.StatusInternalServerError)
		http.ServeFile(w, r, "templates/500.html")
        return
	}
	err1 := tmpl.Execute(w, Respons) 
	if err1 != nil {      
		   
		w.WriteHeader(http.StatusInternalServerError)
		http.ServeFile(w, r, "templates/500.html")
        return
	}
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
}

func ArtistHandler(w http.ResponseWriter, r *http.Request){
	id := r.URL.Query().Get("id")
	
	idint,err:=strconv.Atoi(id)
		if err != nil{
		w.WriteHeader(http.StatusBadRequest)
		http.ServeFile(w, r, "templates/500.html")
		return
		}
	
	var artist Artist

	for _,a:=range Respons{
		if a.ID == idint{
			artist=a
		}
	}

	if artist.ID == 0{
		w.WriteHeader(http.StatusNotFound)
		http.ServeFile(w, r, "templates/404.html")
		return
	}

	tmpl, err := template.ParseFiles("templates/artist.html")
    if err != nil {
		w.WriteHeader(http.StatusNotFound)
        http.ServeFile(w, r, "templates/404.html")
        return
    }
    err = tmpl.Execute(w, []Artist{artist})
    if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.ServeFile(w, r, "templates/404.html")
        return
    }
	

}







