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

func Mustfloat(fn func() (float64, error)) float64 {
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

func resourceQuery(qres *Resource) ([]*Resource, error) {
	ob := initObjectBox()
	rob := BoxForResource(ob)
	var query = rob.Query()
	switch qres.Type {
	case "malt":
		query = maltQuery(rob, qres)
	case "hop":
		query = hopQuery(rob, qres)
	case "yeast":
		query = yeastQuery(rob, qres)
	default:
		log.Panic("Help the Resource has no type!")
	}
	ls, err := query.Find()
	return ls, err

}

func formToResource(r *http.Request) (*Resource, error) {
	validate = validator.New()
	res := &Resource{
		Type: r.FormValue("type"),
		Name: r.FormValue("name"),
		Amount: Mustfloat(func() (float64, error) {
			return strconv.ParseFloat(r.FormValue("amount"), 64)
		}),
		EBC: Mustfloat(func() (float64, error) {
			return strconv.ParseFloat(r.FormValue("ebc"), 64)
		}),
		ISO: Mustfloat(func() (float64, error) {
			return strconv.ParseFloat(r.FormValue("iso"), 64)
		}),
		MinTemp: Mustfloat(func() (float64, error) {
			return strconv.ParseFloat(r.FormValue("minTemp"), 64)
		}),
		MaxTemp: Mustfloat(func() (float64, error) {
			return strconv.ParseFloat(r.FormValue("maxTemp"), 64)
		}),
		OberG: Mustbool(func() (bool, error) {
			return strconv.ParseBool(r.FormValue("og"))
		}),
	}
	err := validate.Struct(res)
	return res, err

}

func maltQuery(rob *ResourceBox, qres *Resource) *ResourceQuery {
	query := rob.Query(Resource_.Name.Contains(qres.Name, false),
		Resource_.Amount.GreaterOrEqual(qres.Amount))
	if qres.EBC != 0 {
		query = rob.Query(Resource_.Name.Contains(qres.Name, false),
			Resource_.Amount.GreaterOrEqual(qres.Amount),
			Resource_.EBC.Between(qres.EBC-1, qres.EBC+1))
	}
	return query
}

func hopQuery(rob *ResourceBox, qres *Resource) *ResourceQuery {
	query := rob.Query(Resource_.Name.Contains(qres.Name, false),
		Resource_.Amount.GreaterOrEqual(qres.Amount))
	if qres.ISO != 0 {
		query = rob.Query(Resource_.Name.Contains(qres.Name, false),
			Resource_.Amount.GreaterOrEqual(qres.Amount),
			Resource_.ISO.Between(qres.ISO-0.1, qres.ISO+0.1))
	}
	return query
}

func yeastQuery(rob *ResourceBox, qres *Resource) *ResourceQuery {
	query := rob.Query(Resource_.Name.Contains(qres.Name, false),
		Resource_.Amount.GreaterOrEqual(qres.Amount))
	if qres.MinTemp != 0 && qres.MaxTemp != 0 {
		query = rob.Query(Resource_.Name.Contains(qres.Name, false),
			Resource_.Amount.GreaterOrEqual(qres.Amount),
			Resource_.OberG.Equals(qres.OberG),
			Resource_.MinTemp.Between(qres.MinTemp-1, qres.MinTemp+1),
			Resource_.MaxTemp.Between(qres.MaxTemp-1, qres.MaxTemp+1))
		return query
	} else if qres.MinTemp != 0 {
		query = rob.Query(Resource_.Name.Contains(qres.Name, false),
			Resource_.OberG.Equals(qres.OberG),
			Resource_.Amount.GreaterOrEqual(qres.Amount),
			Resource_.MinTemp.Between(qres.MinTemp-1, qres.MinTemp+1))
		return query
	} else if qres.MaxTemp != 0 {
		query = rob.Query(Resource_.Name.Contains(qres.Name, false),
			Resource_.OberG.Equals(qres.OberG),
			Resource_.Amount.GreaterOrEqual(qres.Amount),
			Resource_.MaxTemp.Between(qres.MaxTemp-1, qres.MaxTemp+1))
		return query
	} else {
		return query
	}
}
