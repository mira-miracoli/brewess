package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

var templates, terr = template.ParseFiles("./html/new.html",
	"./html/view.html", "./html/edit.html", "./html/searchres.html",
	"./html/badsearch.html")

func main() {
	if terr != nil {
		fmt.Print(terr.Error())
	}
	http.HandleFunc("/home/", homeHandler)
	http.HandleFunc("/search-resource/", searchResHandler)
	http.HandleFunc("/search-results/", resultsResHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func renderTemplate(w http.ResponseWriter, tmpl string, r *Resource) {
	err := templates.ExecuteTemplate(w, tmpl+".html", r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./html/home.html")
}

func searchResHandler(w http.ResponseWriter, r *http.Request) {
	res := &Resource{}
	renderTemplate(w, "searchres", res)

}

func resultsResHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Please use the search form", http.StatusBadRequest)
		return
	}
	qres, err := formToResource(r)
	if err != nil {
		http.ServeFile(w, r, "./html/badsearch.html")
		return
	}
	// call OB qery function
	_, qerr := resourceQuery(qres)
	if qerr != nil {
		http.ServeFile(w, r, "./html/badsearch.html")
		return
	} else {
		http.Error(w, "search worked, no error", http.StatusOK)
	}

	// render Template with Result List if List is empty return search failed page
}
