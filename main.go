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

type Resources struct {
	Resources []*Resource
}

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

	renderSingle := func(w http.ResponseWriter, templateName string, resources *Resource) {
		err := templates.ExecuteTemplate(w, templateName+".html", resources)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	renderList := func(responseWriter http.ResponseWriter, tmpl string, resources *Resources) {
		err := templates.ExecuteTemplate(responseWriter, tmpl+".html", resources)
		if err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		}
	}

	staticHandler := func(w http.ResponseWriter, request *http.Request) {
		fmt.Println(request.URL.Path)
		res := &Resource{}
		switch request.URL.Path {
		case "/search-resource/":
			renderSingle(w, "searchres", res)
		case "/new-resource/":
			renderSingle(w, "newres", res)
		case "/edit-recipe/":
			renderSingle(w, "editrecipe", res)
		case "/home/":
			http.ServeFile(w, request, "./html/home.html")
		default:

			http.Redirect(w, request, "http://localhost:8080/home/", http.StatusNotFound)

		}
	}

	searchResultsResourceHandler := func(responseWriter http.ResponseWriter, request *http.Request) {
		resBox := BoxForResource(objectBox)

		if request.Method != http.MethodPost {
			http.Error(responseWriter, "Please use the search form", http.StatusBadRequest)
			return
		}
		searchResource, err := formToResource(request)
		if err != nil {
			http.ServeFile(responseWriter, request, "./html/badsearch.html")
			return
		}
		// call OB qery function
		foundResources, queryError := resourceQuery(searchResource, resBox)
		if queryError != nil {
			http.ServeFile(responseWriter, request, "./html/badsearch.html")
			return
		} else {
			renderList(responseWriter, "resultsres", &Resources{Resources: foundResources})
		}

		// render Template with Result List if List is empty return search failed page
	}

	saveResourceHandler := func(responseWriter http.ResponseWriter, request *http.Request) {
		resBox := BoxForResource(objectBox)
		if request.Method != http.MethodPost {
			http.Error(responseWriter, "Please use the resource creation form", http.StatusBadRequest)
			return
		}
		res, err := formToResource(request)
		if err != nil {
			http.Redirect(responseWriter, request, "/new-resource/", http.StatusInternalServerError)
			return
		}
		_, err = resBox.Put(res)
		if err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		}
	}

	deleteResourceHandler := func(responseWriter http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			http.Error(responseWriter, "No valid delete request", http.StatusBadRequest)
			return
		}
		resourceID, err := strconv.ParseUint(request.FormValue("ajax_post_data"), 10, 64)
		if err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		}
		resBox := BoxForResource(objectBox)
		err = resBox.RemoveId(resourceID)
		if err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		}
	}
	resourceMarshal := func(responseWriter http.ResponseWriter, request *http.Request) {
		box := BoxForResource(objectBox)

		if request.Method != http.MethodGet {
			http.Error(responseWriter, "Invalid AJAX request", http.StatusBadRequest)
			return
		}
		requestedType := request.FormValue("ajax_post_data")
		searchResource := new(Resource)
		searchResource.Type = requestedType
		foundResources, err := resourceQuery(searchResource, box)
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
		resBox := BoxForResource(objectBox)
		usedResBox := BoxForUsedResource(objectBox)
		recipeBox := BoxForRecipe(objectBox)
		mashBox := BoxForMashStep(objectBox)
		recipe, err := formToRecipe(request, resBox, usedResBox, mashBox)
		//to do
		if err != nil {
			http.Redirect(responseWriter, request, "/new-resource/", http.StatusInternalServerError)
			return
		}
		_, err = recipeBox.Put(recipe)
		if err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		}
	}

	if templateError != nil {
		fmt.Print(templateError.Error())
	}
	http.HandleFunc("/", staticHandler)
	http.HandleFunc("/search-results/", searchResultsResourceHandler)
	http.HandleFunc("/save-resource/", saveResourceHandler)
	http.HandleFunc("/delete-resource/", deleteResourceHandler)
	http.HandleFunc("/save-recipe/", saveRecipeHandler)
	sem <- 1
	go http.HandleFunc("/get-json/", resourceMarshal)
	<-sem
	log.Fatal(http.ListenAndServe(":8080", nil))
}
