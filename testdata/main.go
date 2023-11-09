package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func main() {
	acme := Acme{
		Name:        "Acme Inc.",
		YearOfBirth: 2018,
		Employees: []*Employee{
			{
				Name:     "John Doe",
				Position: "CEO",
				Wage:     100000,
				JoinedAt: time.Now(),
				Projects: []Project{
					{
						Name:       "Project 1",
						Value:      100000,
						StartedAt:  time.Now(),
						FinishedAt: time.Now().Add(time.Hour * 24 * 30),
					},
					{
						Name:       "Project 2",
						Value:      200000,
						StartedAt:  time.Now(),
						FinishedAt: time.Now().Add(time.Hour * 24 * 30),
					},
				},
			},
			{
				Name:     "Jane Doe",
				Position: "CTO",
				Wage:     80000,
				JoinedAt: time.Now(),
				Projects: []Project{
					{
						Name:       "Project 3",
						Value:      300000,
						StartedAt:  time.Now(),
						FinishedAt: time.Now().Add(time.Hour * 24 * 30),
					},
					{
						Name:       "Project 4",
						Value:      400000,
						StartedAt:  time.Now(),
						FinishedAt: time.Now().Add(time.Hour * 24 * 30),
					},
				},
			},
		},
		Address: &Address{
			Street: "Main Street",
			Number: 123,
			City:   "New York",
			Zip:    10001,
		},
		Vat: Vat{
			Number: "123456789",
			Type:   "Company",
		},
	}

	mutator := NewMutatorAcme(&acme)
	mutator.MutateYearOfBirth(2019)

	newAddr := &Address{
		Street: "Liverpool Street",
		Number: 300,
		City:   "London",
		Zip:    45001,
	}

	uk := "UK"

	mutator.SetAddress(nil)
	mutator.SetAddress(newAddr)
	mutator.MutateVat().MutateType("Company Ltd.")
	mutator.MutateAddress().MutateStreet("Baker Street")
	mutator.MutateAddress().SetLocation(&uk)
	mutator.MutateEmployeesAt(0).MutateName("John Smith")
	mutator.MutateEmployeesAt(0).MutateProjectsAt(0).MutateName("Project 1 - Updated")
	mutator.MutateEmployeesAt(0).MutateName("John Smith")
	mutator.AppendEmployees(&Employee{
		Name:     "Roger Smith",
		Position: "CFO",
		Wage:     90000,
		JoinedAt: time.Now(),
		Projects: nil,
	})

	for _, change := range mutator.FormatChanges() {
		println(change)
	}

	t := &testing.T{}

	assert.Equal(t, 2019, acme.YearOfBirth)
	assert.Equal(t, newAddr, acme.Address)
	assert.Equal(t, "Company Ltd.", acme.Vat.Type)
	assert.Equal(t, "Baker Street", acme.Address.Street)
	assert.Equal(t, "UK", *acme.Address.Location)
	assert.Equal(t, "John Smith", acme.Employees[0].Name)
	assert.Equal(t, "Project 1 - Updated", acme.Employees[0].Projects[0].Name)
}
