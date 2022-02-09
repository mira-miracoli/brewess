# Brewess
A tool for brewing recipe management and calculation
...which is currently under construction  
  
TODO:
- [X] add the "Create Recipe" Page
- [X] build a working brewing calculator
- [X] get resources from database and fill selectors
- [ ] adding "save recipe" feature and corresponding model
- [ ] adding "Brew it" button that substracts the needed amounts from resource entries in database
- [ ] fix the spaghetti-code
- [ ] fix CSS and nav-bar

## Installation
* Tested on GNU/Linux only
* You need at least go 1.17
* Install Objectbox using
```
bash <(curl -s https://raw.githubusercontent.com/objectbox/objectbox-go/main/install.sh)
```
(you also need gcc or clang for this)
* Clone the repository
* Run the following commands in "brewess" folder:
```
go mod tidy
go run .
```
* You can access the app with a browser via [localhost:8080/home](url)

## Resource Management
Brewess supports the holy trinity: malt, hop and yeast (water would be rather complex and for most hobby brewesses and brewers hard to manipulate so I replaced it with yeast).
You can use the web interface to create all your resources and and search for it using the search mask.
I'll explain how the search fields work:
#### General Properties
* Name: has to be part of the name sting (e.g. searching for Cara will result both CaraAmber and CaraDunkel) empty matches all
* Amount: is, for practical reasons a minimum - all resource entries greater or eq will be shown
#### Malt-specific
* EBC: Unit for beer coloring property of malt, the higher, the darker. Shows all Malts +/- 1 EBC (to find alternatives easily)
#### Hop-specific
* Alpha-acids: the fraction of alpha-acids in hops tell, simplified, how bitter the beer will taste, the more alpha-acid, the more bitter. (+/- 0.1 % tolerance in results)
* (maybe I'll add a list for main flavors)
#### Yeast-specific
* Min. Temperature: the minimum fermenting temperature, according to yeast data sheets. +/- 1 °C tolerance in search results
* Max. Temperature: the maximum fermenting temperature, according to yeast data sheets. +/- 1 °C tolerance in search results
* Yeast Breakage Behavior: indicates if the yeast is top or bottom fermenting. search will (yet) only filter if at least one other field is used simultaneously
* (maybe I'll add the optimum fermenting temperature and a string for main flavors) 

