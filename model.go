package main

import "time"

type Empoyee struct {
	ID           uint64
	FirstName    string
	LastName     string
	PlaceOfBirth string
	DateOfBirth  time.Time
	Address      string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
