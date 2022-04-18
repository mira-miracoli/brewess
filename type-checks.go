package main

import (
	"errors"
	"log"
	"strconv"
	"time"
)

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
