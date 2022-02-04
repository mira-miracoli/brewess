package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

type Resourcels struct {
	Resources []*Resource
}

var templates, terr = template.ParseFiles("./html/newres.html",
	"./html/searchres.html",
	"./html/badsearch.html", "./html/resultsres.html")

func main() {
	if terr != nil {
		fmt.Print(terr.Error())
	}
	http.HandleFunc("/", staticHandler)
	http.HandleFunc("/search-results/", resultsResHandler)
	http.HandleFunc("/save-resource/", saveResHandler)
	http.HandleFunc("/delete-resource/", deleteResHandler)
	http.HandleFunc("/get-json/", resourceMarshal)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func renderSingle(w http.ResponseWriter, tmpl string, r *Resource) {
	err := templates.ExecuteTemplate(w, tmpl+".html", r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func renderList(w http.ResponseWriter, tmpl string, r *Resourcels) {
	err := templates.ExecuteTemplate(w, tmpl+".html", r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	switch r.URL.Path {
	case "/search-resource/":
		res := &Resource{}
		renderSingle(w, "searchres", res)
	case "/new-resource/":
		res := &Resource{}
		renderSingle(w, "newres", res)
	case "/home/":
		http.ServeFile(w, r, "./html/home.html")
	default:

		http.Redirect(w, r, "http://localhost:8080/home/", http.StatusNotFound)

	}
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
	ls, qerr := resourceQuery(qres)
	if qerr != nil {
		http.ServeFile(w, r, "./html/badsearch.html")
		return
	} else {
		ls := &Resourcels{Resources: ls}
		renderList(w, "resultsres", ls)
	}

	// render Template with Result List if List is empty return search failed page
}

func resourceMarshal(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid AJAX request", http.StatusBadRequest)
		return
	}
	res_sting := r.FormValue("ajax_post_data")
	qres := &Resource{Type: res_sting}
	res_ls, err := resourceQuery(qres)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	res_json, err := json.Marshal(res_ls)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fmt.Fprint(w, string(res_json))
}

func saveResHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Please use the resource creation form", http.StatusBadRequest)
		return
	}
	res, err := formToResource(r)
	if err != nil {
		http.Redirect(w, r, "/new-resource/", http.StatusInternalServerError)
		return
	}
	_, err = saveResource(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func deleteResHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "No valid delete request", http.StatusBadRequest)
		return
	}
	id_str := r.FormValue("ajax_post_data")
	id, err := strconv.ParseUint(id_str, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	ob := initObjectBox()
	rb := BoxForResource(ob)
	defer ob.Close()
	err = rb.RemoveId(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
