package main

import (
	"fmt"
	"time"
)

type Acme struct {
	Name        string
	YearOfBirth int
	Employees   []*Employee
	Address     *Address
	Vat         Vat
}

type Address struct {
	Street   string
	Number   int
	City     string
	Zip      int
	Location *string
}

type Vat struct {
	Number string
	Type   string
}

type Employee struct {
	Name     string
	Position string
	Wage     int
	JoinedAt time.Time
	Projects []Project
}

func (e *Employee) String() string {
	return fmt.Sprintf("%s - %s - %v - %s - %v", e.Name, e.Position, e.Wage, e.JoinedAt, e.Projects)
}

type Project struct {
	Name       string
	Value      int
	StartedAt  time.Time
	FinishedAt time.Time
}
