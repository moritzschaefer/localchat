package main

import (
	"fmt"
	"strconv"
	"strings"
)

func ParsePosition(serialInput string) (Position, error) {
	// Returns the input string, which should be of form "latitude, longitude" as a Position-struct

	latlong := strings.Split(serialInput, ",")
	if len(latlong) != 2 {
		return Position{}, fmt.Errorf("Expected two values separated by one comma")
	}
	lat, err := strconv.ParseFloat(latlong[0], 64)
	if err != nil {
		return Position{}, err
	}
	long, err := strconv.ParseFloat(latlong[1], 64)
	if err != nil {
		return Position{}, err
	}
	return Position{lat, long}, nil

}
