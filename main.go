package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/objectbox/objectbox-go/objectbox"
)

var sem = make(chan int, 1)

var validate *validator.Validate

type BoxFor struct {
	Hop          *HopBox
	Malt         *MaltBox
	Yeast        *YeastBox
	Recipe       *RecipeBox
	UsedResource *UsedResourceBox
	MashStep     *MashStepBox
}

var templates, templateError = template.ParseFiles("./html/newres.html",
	"./html/searchres.html",
	"./html/badsearch.html", "./html/resultsres.html", "./html/editrecipe.html")

func initObjectBox() *objectbox.ObjectBox {
	objectBox, err := objectbox.NewBuilder().Model(ObjectBoxModel()).Build()
	if err != nil {
		panic(err)
	}
	return objectBox
}

func main() {
	objectBox := initObjectBox()
	defer objectBox.Close()
	uniBox := &BoxFor{
		Hop:          BoxForHop(objectBox),
		Malt:         BoxForMalt(objectBox),
		Yeast:        BoxForYeast(objectBox),
		Recipe:       BoxForRecipe(objectBox),
		UsedResource: BoxForUsedResource(objectBox),
		MashStep:     BoxForMashStep(objectBox),
	}

	renderSingle := func(w http.ResponseWriter, templateName string) {
		err := templates.ExecuteTemplate(w, templateName+".html", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	renderList := func(responseWriter http.ResponseWriter, tmpl string, resources *ResourceLists) {
		err := templates.ExecuteTemplate(responseWriter, tmpl+".html", resources)
		if err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		}
	}

	staticHandler := func(w http.ResponseWriter, request *http.Request) {
		fmt.Println(request.URL.Path)
		switch request.URL.Path {
		case "/search-resource/":
			renderSingle(w, "searchres")
		case "/new-resource/":
			renderSingle(w, "newres")
		case "/edit-recipe/":
			renderSingle(w, "editrecipe")
		case "/home/":
			http.ServeFile(w, request, "./html/home.html")
		default:

			http.Redirect(w, request, "http://localhost:8080/home/", http.StatusNotFound)

		}
	}

	searchResultsResourceHandler := func(responseWriter http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			http.Error(responseWriter, "Please use the search form", http.StatusBadRequest)
			return
		}
		searchResource, err := FormToResource(request)
		if err != nil {
			http.ServeFile(responseWriter, request, "./html/badsearch.html")
			return
		}
		// call OB qery function
		foundResources, queryErr := searchResource.Query(uniBox)
		if queryErr != nil {
			http.ServeFile(responseWriter, request, "./html/badsearch.html")
			return
		} else {
			renderList(responseWriter, "resultsres", foundResources)
		}

		// render Template with Result List if List is empty return search failed page
	}

	saveResourceHandler := func(responseWriter http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			http.Error(responseWriter, "Please use the resource creation form", http.StatusBadRequest)
			return
		}
		res, err := FormToResource(request)
		if err != nil {
			http.Redirect(responseWriter, request, "/new-resource/", http.StatusInternalServerError)
			return
		}
		if _, err := res.PutInBox(uniBox); err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		}
	}

	deleteResourceHandler := func(responseWriter http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			http.Error(responseWriter, "No valid delete request", http.StatusBadRequest)
			return
		}
		resource, _ := FormToResource(request)
		resourceID, err := strconv.ParseUint(request.FormValue("ajax_post_data"), 10, 64)
		if err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		}
		resource.RemoveByID(resourceID, uniBox)
	}
	ResourceMarshal := func(responseWriter http.ResponseWriter, request *http.Request) {

		if request.Method != http.MethodGet {
			http.Error(responseWriter, "Invalid AJAX request", http.StatusBadRequest)
			return
		}
		resource, err := FormToResource(request)
		if err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		}
		foundResources, err := resource.Query(uniBox)
		if err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		}
		marshalledResources, err := json.Marshal(foundResources)
		if err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		}
		fmt.Fprint(responseWriter, string(marshalledResources))
	}

	saveRecipeHandler := func(responseWriter http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			http.Error(responseWriter, "Please use the resource creation form", http.StatusBadRequest)
			return
		}
		if err := formToRecipe(request, uniBox); err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		}
	}

	if templateError != nil {
		fmt.Print(templateError.Error())
	}
	http.HandleFunc("/", staticHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/search-results/", searchResultsResourceHandler)
	http.HandleFunc("/save-resource/", saveResourceHandler)
	http.HandleFunc("/delete-resource/", deleteResourceHandler)
	http.HandleFunc("/save-recipe/", saveRecipeHandler)
	sem <- 1
	go http.HandleFunc("/get-json/", ResourceMarshal)
	<-sem
	log.Fatal(http.ListenAndServe(":8080", nil))
}
