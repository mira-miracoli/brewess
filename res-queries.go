package main

import (
	"errors"
	"net/http"
	"reflect"

	"github.com/go-playground/validator/v10"
)

type TempRange struct {
	Min float64
	Max float64
}

func formToResource(request *http.Request) (Resource, error) {
	resource := makeResource(request)
	if reflect.TypeOf(resource).Kind().String() == "invalidResource" {
		return resource, errors.New("No such Resource type")
	}
	resource.SetName(request)
	resource.SetAmount(request)
	resource.SetCharacteristics(request)
	err := validator.New().Struct(resource)
	return resource, err
}

func (queryResource *Hop) Query(box *BoxFor) *ResourceQuery {
	query := box.Hop.Query(Hop_.AbstractResource_Name.Contains(queryResource.Name, false),
		Hop_.AbstractResource_Amount.GreaterOrEqual(queryResource.Amount))
	if queryResource.GetISO() != 0 {
		query = box.Hop.Query(Hop_.AbstractResource_Name.Contains(queryResource.Name, false),
			Hop_.AbstractResource_Amount.GreaterOrEqual(queryResource.Amount),
			Hop_.ISO.Between(queryResource.ISO-1, queryResource.ISO+1))
	}
	return query
}

func (queryResource *Malt) Query(box *BoxFor) *ResourceQuery {
	query := box.Malt.Query(
		Malt_.AbstractResource_Name.Contains(queryResource.Name, false),
		Malt_.AbstractResource_Amount.GreaterOrEqual(queryResource.Amount))
	if queryResource.GetEBC() != 0 {
		query = box.Malt.Query(Malt_.AbstractResource_Name.Contains(queryResource.Name, false),
			Malt_.AbstractResource_Amount.GreaterOrEqual(queryResource.Amount),
			Malt_.EBC.Between(queryResource.EBC-0.1, queryResource.EBC+0.1))
	}
	return query
}

//Please split into smaller functions!
func (queryResource *Yeast) Query(box *BoxFor) (*ListOf, error) {
	minTemp := queryResource.SetQueryMinTemp()
	maxTemp := queryResource.SetQueryMaxTemp()
	query := box.Yeast.Query(
		Yeast_.AbstractResource_Name.Contains(queryResource.Name, false),
		Yeast_.AbstractResource_Amount.GreaterOrEqual(queryResource.Amount),
		Yeast_.OberG.Equals(queryResource.OberG),
		Yeast_.MinTemp.Between(minTemp.Min, minTemp.Max),
		Yeast_.MaxTemp.Between(maxTemp.Min, maxTemp.Max))
	result, err := query.Find()
	foundItems := &ListOf{Yeast: result}
	return foundItems, err
}
