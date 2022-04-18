package main

import (
	"net/http"
	"strconv"
)

type AbstractResource struct {
	Name   string  `validate: "alphanum"`
	Amount float64 `validate: "numeric"`
}

type Resource interface {
	SetName(request *http.Request)
	SetAmount(request *http.Request)
	SetCharacteristics(request *http.Request)
	GetID() uint64
	GetName() string
	GetAmount() float64
	Query(box *BoxFor) (ListOf, error)
}

type ListOf struct {
	Hop   []*Hop
	Malt  []*Malt
	Yeast []*Yeast
}

type InvalidResource struct {
	AbstractResource
}

func makeResource(request *http.Request) Resource {
	switch request.FormValue("type") {
	case "hop":
		return new(Hop)
	case "malt":
		return new(Malt)
	case "yeast":
		return new(Yeast)
	default:
		invalid := new(InvalidResource)
		return invalid
	}
}

//Setting Methods
//Abstract Type Methods
func (resource *AbstractResource) SetName(request *http.Request) {
	resource.Name = request.FormValue("name")
}
func (resource *AbstractResource) SetAmount(request *http.Request) {
	resource.Amount = Mustfloat(func() (float64, error) {
		return strconv.ParseFloat(request.FormValue("amount"), 64)
	})
}
func (resource *AbstractResource) GetName() string {
	return resource.Name
}
func (resource *AbstractResource) GetAmount() float64 {
	return resource.Amount
}

//Concrete Type Methods
func (resource *Malt) SetCharacteristics(request *http.Request) {
	resource.EBC = Mustfloat(func() (float64, error) {
		return strconv.ParseFloat(request.FormValue("ebc"), 64)
	})
}
func (resource *Hop) SetCharacteristics(request *http.Request) {
	resource.ISO = Mustfloat(func() (float64, error) {
		return strconv.ParseFloat(request.FormValue("iso"), 64)
	})
}
func (resource *Yeast) SetCharacteristics(request *http.Request) {
	resource.MinTemp = Mustfloat(func() (float64, error) {
		return strconv.ParseFloat(request.FormValue("minTemp"), 64)
	})
	resource.MaxTemp = Mustfloat(func() (float64, error) {
		return strconv.ParseFloat(request.FormValue("maxTemp"), 64)
	})
	resource.OberG = Mustbool(func() (bool, error) {
		return strconv.ParseBool(request.FormValue("og"))
	})
}
func (resource *InvalidResource) SetCharacteristics(request *http.Request) {

}
func (resource *Hop) GetID() uint64 {
	return resource.Id
}
func (resource *Malt) GetID() uint64 {
	return resource.Id
}
func (resource *Yeast) GetID() uint64 {
	return resource.Id
}
func (resource *InvalidResource) GetID() uint64 {
	return 0
}
func (queryResource *InvalidResource) Query(box *BoxFor) {
}

//other Methods
func (resource *Hop) GetISO() float64 {
	return resource.ISO
}
func (resource *Malt) GetEBC() float64 {
	return resource.EBC
}
func (resource *Yeast) GetOberG() bool {
	return resource.OberG
}
func (resource *Yeast) GetMinTemp() float64 {
	return resource.MinTemp
}
func (resource *Yeast) GetMaxTemp() float64 {
	return resource.MaxTemp
}
func (resource *Yeast) SetQueryMinTemp() *TempRange {
	if resource.MinTemp != 0 {
		return &TempRange{resource.MinTemp - 1, resource.MinTemp + 1}
	} else {
		return &TempRange{0, 100}
	}
}
func (resource *Yeast) SetQueryMaxTemp() *TempRange {
	if resource.MaxTemp != 0 {
		return &TempRange{resource.MaxTemp - 1, resource.MaxTemp + 1}
	} else {
		return &TempRange{0, 100}
	}
}
