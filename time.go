package util

import (
	"time"
)

// CST get china location
// if error will get local location
func CST() *time.Location {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		loc = time.Local
	}
	return loc
}

// CSTE Get CST location with error
func CSTE() (*time.Location, error) {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return nil, err
	}
	return loc, nil
}
