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
	RemoveByID(id uint64, uniBox *BoxFor)
	GetName() string
	GetAmount() float64
	Query(box *BoxFor) (*ResourceLists, error)
	PutInBox(box *BoxFor) (uint64, error)
}

type ResourceLists struct {
	Hop   []*Hop
	Malt  []*Malt
	Yeast []*Yeast
}

type InvalidResource struct {
	AbstractResource
}

func makeResource(resourceType string) Resource {
	switch resourceType {
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
	resource.OberG = request.FormValue("og")
}
func (resource *UsedResource) SetCharacteristics(request *http.Request) {
	resource.Proportion = Mustfloat(func() (float64, error) {
		return strconv.ParseFloat(request.FormValue(""), 64)
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
func (resource *Hop) RemoveByID(id uint64, uniBox *BoxFor) {
	uniBox.Hop.RemoveId(id)
}
func (resource *Malt) RemoveByID(id uint64, uniBox *BoxFor) {
	uniBox.Malt.RemoveId(id)
}
func (resource *Yeast) RemoveByID(id uint64, uniBox *BoxFor) {
	uniBox.Yeast.RemoveId(id)
}
func (resource *InvalidResource) RemoveByID(id uint64, uniBox *BoxFor) {
	return
}
func (queryResource *InvalidResource) Query(box *BoxFor) (*ResourceLists, error) {
	return new(ResourceLists), nil
}

func (resource *Hop) PutInBox(box *BoxFor) (uint64, error) {
	id, err := box.Hop.Put(resource)
	return id, err
}
func (resource *Malt) PutInBox(box *BoxFor) (uint64, error) {
	id, err := box.Malt.Put(resource)
	return id, err
}
func (resource *Yeast) PutInBox(box *BoxFor) (uint64, error) {
	id, err := box.Yeast.Put(resource)
	return id, err
}
func (resource *InvalidResource) PutInBox(box *BoxFor) (uint64, error) {
	return 0, nil
}

//other Methods
func (resource *Hop) GetISO() float64 {
	return resource.ISO
}
func (resource *Malt) GetEBC() float64 {
	return resource.EBC
}
func (resource *Yeast) GetOberG() string {
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

// Methods for UsedResource

func (usedResource *UsedResource) ValidateAndPut(uniBox *BoxFor) error {
	if err := validate.Struct(usedResource); err != nil {
		return err
	}
	if _, err := uniBox.UsedResource.Put(usedResource); err != nil {
		return err
	}
	return nil
}
func (usedResource *UsedResource) SetProportion(formName string) {
	usedResource.Proportion = Mustfloat(func() (float64, error) {
		return strconv.ParseFloat(formName, 64)
	})
}
func (usedResource *UsedResource) SetCookingTime(formName string) {
	usedResource.CookingTime = Mustfloat(func() (float64, error) {
		return strconv.ParseFloat(formName, 64)
	})
}

// Methods for MashStep

func (mashStep *MashStep) SetAndValidate(mashTemp string, mashTime string) error {
	mashStep.Temp = Mustfloat(func() (float64, error) {
		return strconv.ParseFloat(mashTemp, 64)
	})
	mashStep.Time = Mustfloat(func() (float64, error) {
		return strconv.ParseFloat(mashTime, 64)
	})
	return validate.Struct(mashStep)
}

//Methods for ResourceList
