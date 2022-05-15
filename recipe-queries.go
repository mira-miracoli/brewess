package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
)

type unmarshalResource struct {
	Id  uint64
	EBC float64
	ISO float64
}

func (res *unmarshalResource) JSONToId(data string) (uint64, error) {
	err := json.Unmarshal([]byte(data), res)
	return res.Id, err
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
		if data == "" {
			continue
		}
		resource, err := uniBox.Malt.Get(new(unmarshalResource).JSONToId(data))
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
		resource, err := uniBox.Malt.Get(new(unmarshalResource).JSONToId(data))
		if err != nil {
			return usedResources, err
		}
		usedResource := &UsedResource{
			ResourceID: resource.GetID(),
		}
		usedResource.SetProportion("Hopperl" + strconv.Itoa(count))
		usedResource.SetCookingTime("Hoptime" + strconv.Itoa(count))
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
			ResourceID: resource.GetID(),
		}
		usedResource.SetProportion("yestprop" + strconv.Itoa(count))
		if err := usedResource.ValidateAndPut(uniBox); err != nil {
			return usedResources, err
		}
		usedResources = append(usedResources, usedResource)
	}
	return usedResources, nil
}

func formToRecipe(r *http.Request, uniBox *BoxFor) error {
	validate = validator.New()
	mashSteps, err := mashSteps(r, uniBox)
	if err != nil {
		return err
	}
	usedMalts, err := formToUsedMalts(r, uniBox)
	if err != nil {
		return err
	}
	usedHops, err := formToUsedHops(r, uniBox)
	if err != nil {
		return err
	}
	usedYeasts, err := formToUsedYeasts(r, uniBox)
	if err != nil {
		return err
	}
	recipe := &Recipe{
		Name:             r.FormValue("name"),
		BasicInfo:        r.FormValue("destext"),
		HopInfo:          r.FormValue("hopping-notes"),
		MaltInfo:         r.FormValue("grist-notes"),
		MashInfo:         r.FormValue("mashing-notes"),
		FermentationInfo: r.FormValue("fermentation-notes"),
		CastWorth: Mustfloat(func() (float64, error) {
			return strconv.ParseFloat(r.FormValue("castwort"), 64)
		}),
		EBC: Mustfloat(func() (float64, error) {
			return strconv.ParseFloat(r.FormValue("EBC"), 64)
		}),
		IBU: Mustfloat(func() (float64, error) {
			return strconv.ParseFloat(r.FormValue("IBU"), 64)
		}),
		OGTarget: Mustfloat(func() (float64, error) {
			return strconv.ParseFloat(r.FormValue("OG"), 64)
		}),
		SHA: Mustfloat(func() (float64, error) {
			return strconv.ParseFloat(r.FormValue("SHA"), 64)
		}),
		CookingTime: Mustfloat(func() (float64, error) {
			return strconv.ParseFloat(r.FormValue("CookingTime"), 64)
		}),
		MashSteps: mashSteps,
		Malts:     usedMalts,
		Hops:      usedHops,
		Yeasts:    usedYeasts,
	}
	err = validate.Struct(recipe)
	if err != nil {
		return err
	}
	if _, err = uniBox.Recipe.Put(recipe); err != nil {
		return err
	}
	return nil
}
