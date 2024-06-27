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
	Nicknames   map[string]*Employee
	Equity      map[*Employee]int
}

type Address struct {
	Street   string
	Number   int
	City     string
	Zip      int
	Location *string
}

func (a *Address) String() string {
	if a.Location != nil {
		return fmt.Sprintf("{%s %d, %s %d (%s)}", a.Street, a.Number, a.City, a.Zip, *a.Location)
	}
	return fmt.Sprintf("{%s %d, %s %d}", a.Street, a.Number, a.City, a.Zip)
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

func (e *Employee) KeyForChanges() string {
	return e.Name
}

type Project struct {
	Name       string
	Value      int
	StartedAt  time.Time
	FinishedAt time.Time
	SeqID      []byte
}
