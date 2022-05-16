package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
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

var templates, templateError = template.ParseFiles("./html/newres.html", "./html/editrecipe_templ.html",
	"./html/searchres.html", "./html/recipe_search.html", "./html/recipe_results.html",
	"./html/badsearch.html", "./html/resultsres.html", "./html/editrecipe.html")

var validPath = regexp.MustCompile("^/(edit-recipe)/([0-9]+)$")

func initObjectBox() *objectbox.ObjectBox {
	objectBox, err := objectbox.NewBuilder().Model(ObjectBoxModel()).Build()
	if err != nil {
		panic(err)
	}
	return objectBox
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
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

	renderResourceList := func(responseWriter http.ResponseWriter, tmpl string, resources *ResourceLists) {
		err := templates.ExecuteTemplate(responseWriter, tmpl+".html", resources)
		if err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		}
	}
	renderRecipeList := func(responseWriter http.ResponseWriter, tmpl string, recipes []*Recipe) {
		err := templates.ExecuteTemplate(responseWriter, tmpl+".html", recipes)
		if err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		}
	}
	renderRecipeSingle := func(responseWriter http.ResponseWriter, tmpl string, recipe *Recipe) {
		err := templates.ExecuteTemplate(responseWriter, tmpl+".html", recipe)
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
		case "/search-recipe/":
			renderSingle(w, "recipe_search")
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
			renderResourceList(responseWriter, "resultsres", foundResources)
		}

		// render Template with Result List if List is empty return search failed page
	}
	searchResultsRecipeHandler := func(responseWriter http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			http.Error(responseWriter, "Please use the search form", http.StatusBadRequest)
			return
		}
		searchRecipe, err := new(Recipe).FormToRecipe(request, uniBox)
		if err != nil {
			http.ServeFile(responseWriter, request, "./html/badsearch.html")
			return
		}
		// call OB qery function
		foundRecipes, queryErr := searchRecipe.Query(request, uniBox)
		if queryErr != nil {
			http.ServeFile(responseWriter, request, "./html/badsearch.html")
			return
		} else {
			renderRecipeList(responseWriter, "recipe_results", foundRecipes)
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
	deleteRecipeHandler := func(responseWriter http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			http.Error(responseWriter, "No valid delete request", http.StatusBadRequest)
			return
		}
		recipeID, err := strconv.ParseUint(request.FormValue("ajax_post_data"), 10, 64)
		if err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		}
		err = uniBox.Recipe.RemoveId(recipeID)
		if err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		}
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
		recipe, err := new(Recipe).FormToRecipe(request, uniBox)
		if err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		}
		_, err = uniBox.Recipe.Put(recipe)
		if err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		}

	}
	editRecipeHandler := func(responseWriter http.ResponseWriter, request *http.Request, path string) {
		if request.Method != http.MethodGet {
			http.Error(responseWriter, "Not a GET Request", http.StatusBadRequest)
		}
		recipe, err := uniBox.Recipe.Get(MustUInt(func() (uint64, error) {
			return strconv.ParseUint(path, 10, 64)
		}))
		if err != nil || recipe == nil {
			renderSingle(responseWriter, "editrecipe")
			http.Redirect(responseWriter, request, "/edit-recipe/", http.StatusFound)
		} else {
			renderRecipeSingle(responseWriter, "editrecipe_templ", recipe)
		}
	}

	if templateError != nil {
		fmt.Print(templateError.Error())
	}
	http.HandleFunc("/", staticHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/resource-search-results/", searchResultsResourceHandler)
	http.HandleFunc("/recipe-search-results/", searchResultsRecipeHandler)
	http.HandleFunc("/save-resource/", saveResourceHandler)
	http.HandleFunc("/delete-resource/", deleteResourceHandler)
	http.HandleFunc("/save-recipe/", saveRecipeHandler)
	http.HandleFunc("/delete-recipe/", deleteRecipeHandler)
	http.HandleFunc("/edit-recipe/", makeHandler(editRecipeHandler))
	sem <- 1
	go http.HandleFunc("/get-json/", ResourceMarshal)
	<-sem
	log.Fatal(http.ListenAndServe(":8080", nil))
}
