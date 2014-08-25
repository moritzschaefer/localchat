package main

type Position struct {
	latitude  float64
	longitude float64
}

type User struct {
	username string
	position Position
}
