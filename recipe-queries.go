package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
)

type unmarshalResource struct {
	Id  uint64
	EBC float64
	ISO float64
}

func (res *unmarshalResource) JSONToId(data string) uint64 {
	json.Unmarshal([]byte(data), res)
	return res.Id
}
func mashSteps(r *http.Request, uniBox *BoxFor) ([]*MashStep, error) {
	var mashSteps []*MashStep
	count := 1
	for {
		mashTemp := r.FormValue("mashtemp" + strconv.Itoa(count))
		mashTime := r.FormValue("mashtime" + strconv.Itoa(count))
		if mashTime == "" || mashTemp == "" {
			return mashSteps, nil
		}
		mashStep := new(MashStep)
		if err := mashStep.SetAndValidate(mashTemp, mashTime); err != nil {
			return mashSteps, err
		}
		if _, err := uniBox.MashStep.Put(mashStep); err != nil {
			return mashSteps, err
		}
		mashSteps = append(mashSteps, mashStep)
		count++
	}
}

func formToUsedMalts(r *http.Request, uniBox *BoxFor) ([]*UsedResource, error) {
	var usedResources []*UsedResource
	for count := 0; count < 100; count++ {
		//get the id of used Mals and try to get corresponding Resources from ObjectBox
		data := r.FormValue("selMalt" + strconv.Itoa(count))
		if data != "" {
			id := MustUInt(func() (uint64, error) {
				return strconv.ParseUint(data, 10, 64)
			})
			resource, err := uniBox.Malt.Get(id)
			if err != nil {
				return usedResources, err
			}
			usedResource := &UsedResource{
				ResourceID: resource.GetID(),
				Proportion: Mustfloat(func() (float64, error) {
					return strconv.ParseFloat(r.FormValue("maltprop"+strconv.Itoa(count)), 64)
				}),
				CookingTime: 0,
			}
			if err := usedResource.ValidateAndPut(uniBox); err != nil {
				return usedResources, err
			}
			usedResources = append(usedResources, usedResource)
		}
	}
	return usedResources, nil
}
func formToUsedHops(r *http.Request, uniBox *BoxFor) ([]*UsedResource, error) {
	var usedResources []*UsedResource
	for count := 1; count < 100; count++ {
		//get the id of used Mals and try to get corresponding Resources from ObjectBox
		data := r.FormValue("selHop" + strconv.Itoa(count))
		if data == "" {
			continue
		}
		id := MustUInt(func() (uint64, error) {
			return strconv.ParseUint(data, 10, 64)
		})
		resource, err := uniBox.Malt.Get(id)
		if err != nil {
			return usedResources, err
		}
		usedResource := &UsedResource{
			ResourceID: resource.GetID(),
		}
		usedResource.SetProportion(r.FormValue("hopperl" + strconv.Itoa(count)))
		usedResource.SetCookingTime(r.FormValue("hoptime" + strconv.Itoa(count)))
		fmt.Println(usedResource.Proportion)
		if err := usedResource.ValidateAndPut(uniBox); err != nil {
			return usedResources, err
		}
		usedResources = append(usedResources, usedResource)
	}
	return usedResources, nil
}
func formToUsedYeasts(r *http.Request, uniBox *BoxFor) ([]*UsedResource, error) {
	var usedResources []*UsedResource
	for count := 1; count < 100; count++ {
		data := r.FormValue("selYeast" + strconv.Itoa(count))
		if data == "" {
			continue
		}
		Id := MustUInt(func() (uint64, error) {
			return strconv.ParseUint(data, 10, 64)
		})
		resource, err := uniBox.Yeast.Get(Id)
		if err != nil {
			return usedResources, err
		}
		usedResource := &UsedResource{
			ResourceID:  resource.GetID(),
			CookingTime: 0,
		}
		fmt.Println(r.FormValue("yestprop" + strconv.Itoa(count)))
		usedResource.SetProportion(r.FormValue(("yestprop" + strconv.Itoa(count))))
		fmt.Println(usedResource.Proportion)
		if err := usedResource.ValidateAndPut(uniBox); err != nil {
			return usedResources, err
		}
		usedResources = append(usedResources, usedResource)
	}
	return usedResources, nil
}

func (recipe *Recipe) FormToRecipe(r *http.Request, uniBox *BoxFor) (*Recipe, error) {
	validate = validator.New()
	mashSteps, err := mashSteps(r, uniBox)
	if err != nil {
		return new(Recipe), err
	}
	usedMalts, err := formToUsedMalts(r, uniBox)
	if err != nil {
		return new(Recipe), err
	}
	usedHops, err := formToUsedHops(r, uniBox)
	if err != nil {
		return new(Recipe), err
	}
	usedYeasts, err := formToUsedYeasts(r, uniBox)
	if err != nil {
		return new(Recipe), err
	}
	recipe.Name = r.FormValue("name")
	recipe.BasicInfo = r.FormValue("destext")
	recipe.HopInfo = r.FormValue("hopping-notes")
	recipe.MaltInfo = r.FormValue("grist-notes")
	recipe.MashInfo = r.FormValue("mashing-notes")
	recipe.FermentationInfo = r.FormValue("fermentation-notes")
	recipe.CastWorth = Mustfloat(func() (float64, error) {
		return strconv.ParseFloat(r.FormValue("castwort"), 64)
	})
	recipe.EBC = Mustfloat(func() (float64, error) {
		return strconv.ParseFloat(r.FormValue("EBC"), 64)
	})
	recipe.IBU = Mustfloat(func() (float64, error) {
		return strconv.ParseFloat(r.FormValue("IBU"), 64)
	})
	recipe.OGTarget = Mustfloat(func() (float64, error) {
		return strconv.ParseFloat(r.FormValue("OG"), 64)
	})
	recipe.SHA = Mustfloat(func() (float64, error) {
		return strconv.ParseFloat(r.FormValue("SHA"), 64)
	})
	recipe.CookingTime = Mustfloat(func() (float64, error) {
		return strconv.ParseFloat(r.FormValue("CookingTime"), 64)
	})
	recipe.MashSteps = mashSteps
	recipe.Malts = usedMalts
	recipe.Hops = usedHops
	recipe.Yeasts = usedYeasts

	err = validate.Struct(recipe)
	if err != nil {
		return recipe, err
	}
	return recipe, nil
}

func (recipe *Recipe) Query(r *http.Request, uniBox *BoxFor) ([]*Recipe, error) {
	query := uniBox.Recipe.Query(Recipe_.Name.Contains(recipe.Name, false))
	return query.Find()
}
