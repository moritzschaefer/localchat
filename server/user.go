package main

import "math"

type Degrees float64
type Radians float64

func (g Degrees) toRadians() float64 {
	return float64(math.Pi * g / 180.0)
}

type Position struct {
	latitude  Degrees
	longitude Degrees
}

func (p Position) InRadius(other Position, distance float64) bool {
	// TODO: Calculate distance
	//

	R := 6371.0
	φ1 := p.latitude.toRadians()
	φ2 := other.latitude.toRadians()
	Δφ := (other.latitude - p.latitude).toRadians()
	Δλ := (other.longitude - p.longitude).toRadians()

	a := math.Sin(Δφ/2)*math.Sin(Δφ/2) + math.Cos(φ1)*math.Cos(φ2)*math.Sin(Δλ/2)*math.Sin(Δλ/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	d := R * c
	return d < distance
}

type User struct {
	username string
	position Position
}
