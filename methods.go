package main

import (
	"net/http"
	"strconv"
)

//Setting Methods
func (resource *Resource) SetType(request *http.Request) {
	resource.Type = request.FormValue("type")
}
func (resource *Resource) SetName(request *http.Request) {
	resource.Name = request.FormValue("name")
}
func (resource *Resource) SetEBC(request *http.Request) {
	resource.EBC = Mustfloat(func() (float64, error) {
		return strconv.ParseFloat(request.FormValue("ebc"), 64)
	})
}
func (resource *Resource) SetISO(request *http.Request) {
	resource.ISO = Mustfloat(func() (float64, error) {
		return strconv.ParseFloat(request.FormValue("iso"), 64)
	})
}
func (resource *Resource) SetAmount(request *http.Request) {
	resource.Amount = Mustfloat(func() (float64, error) {
		return strconv.ParseFloat(request.FormValue("amount"), 64)
	})
}
func (resource *Resource) SetMinTemp(request *http.Request) {
	resource.MinTemp = Mustfloat(func() (float64, error) {
		return strconv.ParseFloat(request.FormValue("minTemp"), 64)
	})
}
func (resource *Resource) SetMaxTemp(request *http.Request) {
	resource.MaxTemp = Mustfloat(func() (float64, error) {
		return strconv.ParseFloat(request.FormValue("maxTemp"), 64)
	})
}
func (resource *Resource) SetOberG(request *http.Request) {
	resource.OberG = Mustbool(func() (bool, error) {
		return strconv.ParseBool(request.FormValue("og"))
	})
}

//other Methods
func (resource *Resource) SetQueryMinTemp() (tempRange *TempRange) {
	if resource.MinTemp != 0 {
		return &TempRange{resource.MinTemp - 1, resource.MinTemp + 1}
	} else {
		return &TempRange{0, 100}
	}
}

func (resource *Resource) SetQueryMaxTemp() (tempRange *TempRange) {
	if resource.MaxTemp != 0 {
		return &TempRange{resource.MaxTemp - 1, resource.MaxTemp + 1}
	} else {
		return &TempRange{0, 100}
	}
}
