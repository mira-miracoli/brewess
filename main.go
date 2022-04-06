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

var sem = make(chan int, 1)

var validate *validator.Validate

type Resourcels struct {
	Resources []*Resource
}

var templates, terr = template.ParseFiles("./html/newres.html",
	"./html/searchres.html",
	"./html/badsearch.html", "./html/resultsres.html", "./html/editrecipe.html")

func main() {
	ob := initObjectBox()
	defer ob.Close()
	renderSingle := func(w http.ResponseWriter, tmpl string, r *Resource) {
		err := templates.ExecuteTemplate(w, tmpl+".html", r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	renderList := func(w http.ResponseWriter, tmpl string, r *Resourcels) {
		err := templates.ExecuteTemplate(w, tmpl+".html", r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	staticHandler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL.Path)
		res := &Resource{}
		switch r.URL.Path {
		case "/search-resource/":
			renderSingle(w, "searchres", res)
		case "/new-resource/":
			renderSingle(w, "newres", res)
		case "/edit-recipe/":
			renderSingle(w, "editrecipe", res)
		case "/home/":
			http.ServeFile(w, r, "./html/home.html")
		default:

			http.Redirect(w, r, "http://localhost:8080/home/", http.StatusNotFound)

		}
	}

	resultsResHandler := func(w http.ResponseWriter, r *http.Request) {
		box := BoxForResource(ob)

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
		ls, qerr := resourceQuery(qres, box)
		if qerr != nil {
			http.ServeFile(w, r, "./html/badsearch.html")
			return
		} else {
			ls := &Resourcels{Resources: ls}
			renderList(w, "resultsres", ls)
		}

		// render Template with Result List if List is empty return search failed page
	}

	saveResHandler := func(w http.ResponseWriter, r *http.Request) {
		box := BoxForResource(ob)
		if r.Method != http.MethodPost {
			http.Error(w, "Please use the resource creation form", http.StatusBadRequest)
			return
		}
		res, err := formToResource(r)
		if err != nil {
			http.Redirect(w, r, "/new-resource/", http.StatusInternalServerError)
			return
		}
		_, err = box.Put(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	deleteResHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "No valid delete request", http.StatusBadRequest)
			return
		}
		id_str := r.FormValue("ajax_post_data")
		id, err := strconv.ParseUint(id_str, 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		rb := BoxForResource(ob)
		err = rb.RemoveId(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	resourceMarshal := func(w http.ResponseWriter, r *http.Request) {
		box := BoxForResource(ob)

		if r.Method != http.MethodGet {
			http.Error(w, "Invalid AJAX request", http.StatusBadRequest)
			return
		}
		res_string := r.FormValue("ajax_post_data")
		qres := new(Resource)
		qres.Type = res_string
		res_ls, err := resourceQuery(qres, box)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		res_json, err := json.Marshal(res_ls)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		fmt.Fprint(w, string(res_json))
	}

	saveRecipeHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Please use the resource creation form", http.StatusBadRequest)
			return
		}
		resbox := BoxForResource(ob)
		uresbox := BoxForUsedResource(ob)
		recipebox := BoxForRecipe(ob)
		mashbox := BoxForMashStep(ob)
		recipe, err := formToRecipe(r, resbox, uresbox, mashbox)
		//to do
		if err != nil {
			http.Redirect(w, r, "/new-resource/", http.StatusInternalServerError)
			return
		}
		_, err = recipebox.Put(recipe)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	if terr != nil {
		fmt.Print(terr.Error())
	}
	http.HandleFunc("/", staticHandler)
	http.HandleFunc("/search-results/", resultsResHandler)
	http.HandleFunc("/save-resource/", saveResHandler)
	http.HandleFunc("/delete-resource/", deleteResHandler)
	http.HandleFunc("/save-recipe/", saveRecipeHandler)
	sem <- 1
	go http.HandleFunc("/get-json/", resourceMarshal)
	<-sem
	log.Fatal(http.ListenAndServe(":8080", nil))
}
