package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
)

func mashSteps(r *http.Request, mashbox *MashStepBox) ([]*MashStep, error) {
	var mashSteps []*MashStep
	count := 1
	for {
		mashTemp := r.FormValue("mashtemp" + strconv.Itoa(count))
		mashTime := r.FormValue("mashtime" + strconv.Itoa(count))
		if mashTime == "" || mashTemp == "" {
			break
		}
		mashStep := &MashStep{
			Temp: Mustfloat(func() (float64, error) {
				return strconv.ParseFloat(mashTemp, 64)
			}),
			Time: Mustfloat(func() (float64, error) {
				return strconv.ParseFloat(mashTime, 64)
			}),
		}
		err := validate.Struct(mashStep)
		if err != nil {
			return mashSteps, err
		}
		_, err = mashbox.Put(mashStep)
		if err != nil {
			return mashSteps, err
		}
		mashSteps = append(mashSteps, mashStep)
		count++
	}
	return mashSteps, nil
}

func formToUsedMalts(r *http.Request, resbox *ResourceBox, uresbox *UsedResourceBox) ([]*UsedResource, error) {
	var usedResources []*UsedResource
	type resourceVal struct {
		Id  uint64
		EBC float64
		ISO float64
	}
	var unmarshalVal resourceVal
	for count := 0; count < 100; count++ {
		//get the id of used Mals and try to get corresponding Resources from ObjectBox
		data := r.FormValue("selMalt" + strconv.Itoa(count))
		if data == "" {
			continue
		}
		json.Unmarshal([]byte(data), &unmarshalVal)
		resource, err := resbox.Get(unmarshalVal.Id)
		if err != nil {
			return usedResources, err
		}
		usedResource := &UsedResource{
			Resource: resource,
			Proportion: Mustfloat(func() (float64, error) {
				return strconv.ParseFloat(r.FormValue("maltprop"+strconv.Itoa(count)), 64)
			}),
			CookingTime: 0,
		}
		err = validate.Struct(usedResource)
		if err != nil {
			return usedResources, err
		}
		_, err = uresbox.Put(usedResource)
		if err != nil {
			return usedResources, err
		}
		usedResources = append(usedResources, usedResource)
	}
	return usedResources, nil
}
func formToUsedHops(r *http.Request, resbox *ResourceBox, uresbox *UsedResourceBox) ([]*UsedResource, error) {
	var usedResources []*UsedResource
	type resourceVal struct {
		Id  uint64
		EBC float64
		ISO float64
	}
	var unmarshalVal resourceVal
	for count := 1; count < 100; count++ {
		//get the id of used Mals and try to get corresponding Resources from ObjectBox
		data := r.FormValue("selHop" + strconv.Itoa(count))
		if data == "" {
			continue
		}
		json.Unmarshal([]byte(data), &unmarshalVal)
		resource, err := resbox.Get(unmarshalVal.Id)
		if err != nil {
			return usedResources, err
		}
		usedResource := &UsedResource{
			Resource: resource,
			Proportion: Mustfloat(func() (float64, error) {
				return strconv.ParseFloat(r.FormValue("hopperl"+strconv.Itoa(count)), 64)
			}),
			CookingTime: Mustfloat(func() (float64, error) {
				return strconv.ParseFloat(r.FormValue("hoptime"+strconv.Itoa(count)), 64)
			}),
		}
		err = validate.Struct(usedResource)
		if err != nil {
			return usedResources, err
		}
		_, err = uresbox.Put(usedResource)
		if err != nil {
			return usedResources, err
		}
		usedResources = append(usedResources, usedResource)
	}
	return usedResources, nil
}
func formToUsedYeasts(r *http.Request, resbox *ResourceBox, uresbox *UsedResourceBox) ([]*UsedResource, error) {
	var usedResources []*UsedResource
	for count := 1; count < 100; count++ {
		data := r.FormValue("selYeast" + strconv.Itoa(count))
		if data == "" {
			continue
		}
		Id := MustUInt(func() (uint64, error) {
			return strconv.ParseUint(data, 10, 64)
		})

		resource, err := resbox.Get(Id)
		if err != nil {
			return usedResources, err
		}
		usedResource := &UsedResource{
			Resource: resource,
			Proportion: Mustfloat(func() (float64, error) {
				return strconv.ParseFloat(r.FormValue("yestprop"+strconv.Itoa(count)), 64)
			}),
			CookingTime: 0,
		}
		err = validate.Struct(usedResource)
		if err != nil {
			return usedResources, err
		}
		_, err = uresbox.Put(usedResource)
		if err != nil {
			return usedResources, err
		}
		usedResources = append(usedResources, usedResource)
	}
	return usedResources, nil
}

func formToRecipe(r *http.Request, resbox *ResourceBox, uresbox *UsedResourceBox, mashbox *MashStepBox) (*Recipe, error) {
	validate = validator.New()
	mashSteps, err := mashSteps(r, mashbox)
	if err != nil {
		return nil, err
	}
	usedMalts, err := formToUsedMalts(r, resbox, uresbox)
	if err != nil {
		return nil, err
	}
	usedHops, err := formToUsedHops(r, resbox, uresbox)
	if err != nil {
		return nil, err
	}
	usedYeasts, err := formToUsedYeasts(r, resbox, uresbox)
	if err != nil {
		return nil, err
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
		return recipe, err
	}
	return recipe, err

}
