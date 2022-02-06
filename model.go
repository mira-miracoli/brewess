package main

import "time"

//go:generate go run github.com/objectbox/objectbox-go/cmd/objectbox-gogen

// recipe model
type Resource struct {
	Id      uint64
	Type    string    `validate: "required, oneof= 'malt' 'hop' 'yeast'"`
	Name    string    `validate: "alphanum"`
	Amount  float64   `validate: "numeric`
	EBC     float64   `validate: "numeric, gte=0"`
	MinTemp float64   `validate: "numeric, gte=0"`
	MaxTemp float64   `validate: "numeric, lte=40"`
	OberG   bool      `validate: "boolean"`
	ISO     float64   `validate: "numeric, lte=100"`
	Opened  time.Time `validate: "datetime=2006-01-02"`
	Exp     time.Time `validate: "datetime=2006-01-02"`
}

type Recipe struct {
	Id          uint64
	Name        string
	Description string // short text to descripe and add any comments
	//Malts       map[float64]*Resource
	IsoAlpha float64
	//Hops        map[float64]*Resource // specifies hop-resources to use and its proportion
	HopSugg []string
	DryHop  []string // used for hopping examples
	SHA     float64
	//Yeasts      map[float64]*Resource
	AlcTarget float64 // specifies the targeted amount alc by vol
	OGTarget  float64 // specifies the targeted original gravity in %sacc

}
