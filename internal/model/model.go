package model

//go:generate go run github.com/objectbox/objectbox-go/cmd/objectbox-gogen

// recipe model
type Malt struct {
	Id     uint64
	Name   string
	Amount int
	EBC    int
}

type Hop struct {
	Id     uint64
	Name   string
	Amount int
	Iso    float32
}

type Yeast struct {
	Id      uint64
	Name    string
	Amount  int
	MinTemp float32
	MaxTemp float32
	OberG   bool // if true, the yeast is top-fermenting
}
type Recipe struct {
	Id          uint64
	Name        string
	Description string // short text to descripe and add any comments
	Malts       []*Malt
	IsoAlpha    float32
	Hops        []*Hop // specifies hop-resources to use and its proportion
	HopSugg     []string
	DryHop      []string // used for hopping examples
	SHA         float32
	Yeasts      []*Yeast
	AlcTarget   float32 // specifies the targeted amount alc by vol
	OGTarget    float32 // specifies the targeted original gravity in %sacc

}
