package main

import (
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
)
func mashSteps(r *http.Request, ob *ObjectBox) ([]*MashStep, error){
	var mashSteps []*MashStep
	mashBox := BoxForMashStep(ob)
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
		_, err := mashBox.Put(mashStep)
		if err != nil {
			return mashSteps, err
		}
		mashSteps = append(mashSteps, mashStep)
		count++
	}
	return mashSteps, nil
}

func UsedMalt(r *http.Request, ob *ObjectBox) (*UsedResource, err){
	var usedMalts []*UsedResource
	resBox := BoxForResource(ob)
	count := 1
	type maltVal struct{
		Id uint64
		EBC float64
	}
	var maltVal maltVal
	for {
		var usedMalt UsedResource
		//get the id of used Mals and try to get corresponding Resources from ObjectBox
		data := r.FormValue("selMalt" + strconv.Itoa(count))
		json.Unmarshal([]byte(data), &maltVal)
		id := MustUInt(func() (uint64, error) {
			return strconv.ParseUint(maltVal.Id, 64)
		})
		malt, err := resBox.Get(id)
		if err != nil {
			return usedMalt, err
		}
		//fill ther usedMalt with Form and malt from query

		}
	}
}

func formToRecipe(r *http.Request, ob *ObjectBox) (*Resource, error) {
	validate = validator.New()
	mashSteps, err := mashSteps(r)
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
		Malts:
		Hops:
		Yeasts:
	}
	err := validate.Struct(res)
	return res, err

}
