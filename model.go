package main

//go:generate go run github.com/objectbox/objectbox-go/cmd/objectbox-gogen

// recipe model
type Resource struct {
	Id      uint64
	Type    string  `validate: "required, oneof= 'malt' 'hop' 'yeast'"`
	Name    string  `validate: "alphanum"`
	Amount  float64 `validate: "numeric"`
	EBC     float64 `validate: "numeric, gte=0"`
	MinTemp float64 `validate: "numeric, gte=0"`
	MaxTemp float64 `validate: "numeric, lte=40"`
	OberG   bool    `validate: "boolean"`
	ISO     float64 `validate: "numeric, lte=100"`
}

type MashStep struct {
	Id   uint64
	Temp float64 `validate: "required, numeric, gte=0, lte=100"`
	Time float64 `validate: "required, numeric, gte=0"`
}

type UsedResource struct {
	Id          uint64
	Resource    *Resource
	Proportion  float64 `validate: "numeric, gte=0, lte=100"`
	CookingTime float64 `validate: "numeric, gte=0"`
}

type Recipe struct {
	Id               uint64
	Name             string  `validate: "alphanum"`
	BasicInfo        string  // short text to descripe and add any comments
	HopInfo          string  `validate: "alphanum"`
	MaltInfo         string  `validate: "alphanum"`
	MashInfo         string  `validate: "alphanum"`
	FermentationInfo string  `validate: "alphanum"`
	IBU              float64 `validate: "numeric, gte=0"`
	EBC              float64 `validate: "numeric, gte=0"`
	OGTarget         float64 `validate: "numeric, gte=0"` // specifies the targeted original gravity in %sacc
	CastWorth        float64 `validate: "numeric, gte=0"`
	MashSteps        []*MashStep
	CookingTime      float64 `validate: "numeric, gte=0"`
	SHA              float64 `validate: "numeric, gte=0, lte=100"`
	Malts            []*UsedResource
	Hops             []*UsedResource
	Yeasts           []*UsedResource
}
