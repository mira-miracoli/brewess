package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/objectbox/objectbox-go/objectbox"
)

type TempRange struct {
	Min float64
	Max float64
}

func Mustfloat(fn func() (float64, error)) float64 {
	v, err := fn()
	if !errors.Is(err, strconv.ErrSyntax) && (err != nil) {
		log.Fatalln(err)
	}
	return v
}

func MustUInt(fn func() (uint64, error)) uint64 {
	v, err := fn()
	if !errors.Is(err, strconv.ErrSyntax) && (err != nil) {
		log.Fatalln(err)
	}
	return v
}

func Mustbool(fn func() (bool, error)) bool {
	v, err := fn()
	if !errors.Is(err, strconv.ErrSyntax) && (err != nil) {
		log.Fatalln(err)
	}
	return v
}

func Mustdate(fn func() (time.Time, error)) time.Time {
	v, err := fn()
	if !errors.Is(err, strconv.ErrSyntax) && (err != nil) {
		log.Fatalln(err)
	}
	return v
}

func initObjectBox() *objectbox.ObjectBox {
	objectBox, err := objectbox.NewBuilder().Model(ObjectBoxModel()).Build()
	if err != nil {
		panic(err)
	}
	return objectBox
}

func resourceQuery(queryResource *Resource, box *ResourceBox) ([]*Resource, error) {

	var query = box.Query()
	switch queryResource.Type {
	case "malt":
		query = maltQuery(box, queryResource)
	case "hop":
		query = hopQuery(box, queryResource)
	case "yeast":
		query = yeastQuery(box, queryResource)
	default:
		log.Panic("Help, the Resource has no type!")
	}
	ls, err := query.Find()
	return ls, err

}

func formToResource(request *http.Request) (*Resource, error) {
	res := new(Resource)
	res.SetType(request)
	res.SetName(request)
	res.SetAmount(request)
	res.SetMinTemp(request)
	res.SetMaxTemp(request)
	res.SetOberG(request)
	res.SetISO(request)
	res.SetEBC(request)
	err := validator.New().Struct(res)
	return res, err
}

func maltQuery(rBox *ResourceBox, queryResource *Resource) *ResourceQuery {
	query := rBox.Query(Resource_.Type.Equals("malt", true), Resource_.Name.Contains(queryResource.Name, false),
		Resource_.Amount.GreaterOrEqual(queryResource.Amount))
	if queryResource.EBC != 0 {
		query = rBox.Query(Resource_.Name.Contains(queryResource.Name, false),
			Resource_.Amount.GreaterOrEqual(queryResource.Amount),
			Resource_.EBC.Between(queryResource.EBC-1, queryResource.EBC+1))
	}
	return query
}

func hopQuery(rBox *ResourceBox, queryResource *Resource) *ResourceQuery {
	query := rBox.Query(Resource_.Type.Equals("hop", true), Resource_.Name.Contains(queryResource.Name, false),
		Resource_.Amount.GreaterOrEqual(queryResource.Amount))
	if queryResource.ISO != 0 {
		query = rBox.Query(Resource_.Name.Contains(queryResource.Name, false),
			Resource_.Amount.GreaterOrEqual(queryResource.Amount),
			Resource_.ISO.Between(queryResource.ISO-0.1, queryResource.ISO+0.1))
	}
	return query
}

//Please split into smaller functions!
func yeastQuery(rBox *ResourceBox, queryResource *Resource) *ResourceQuery {
	minTemp := queryResource.SetQueryMinTemp()
	maxTemp := queryResource.SetQueryMaxTemp()
	query := rBox.Query(Resource_.Type.Equals("yeast", true),
		Resource_.Name.Contains(queryResource.Name, false),
		Resource_.Amount.GreaterOrEqual(queryResource.Amount),
		Resource_.OberG.Equals(queryResource.OberG),
		Resource_.MinTemp.Between(minTemp.Min, minTemp.Max),
		Resource_.MinTemp.Between(maxTemp.Min, maxTemp.Max))
	return query
}
