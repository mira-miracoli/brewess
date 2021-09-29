package main

//go:generate go run github.com/objectbox/objectbox-go/cmd/objectbox-gogen

// recipe model
type Resource struct {
	Id      uint64
	Type    string
	Name    string
	Amount  float64
	EBC     float64
	MinTemp float64
	MaxTemp float64
	OberG   bool
	ISO     float64
}

type Recipe struct {
	Id          uint64
	Name        string
	Description string // short text to descripe and add any comments
	Malts       []*Resource
	IsoAlpha    float64
	Hops        []*Resource // specifies hop-resources to use and its proportion
	HopSugg     []string
	DryHop      []string // used for hopping examples
	SHA         float64
	Yeasts      []*Resource
	AlcTarget   float64 // specifies the targeted amount alc by vol
	OGTarget    float64 // specifies the targeted original gravity in %sacc

}
