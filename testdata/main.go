package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func main() {
	// use a fixed date as now
	now := time.Date(2023, 10, 30, 13, 14, 15, 0, time.UTC)

	acme := Acme{
		Name:        "Acme Inc.",
		YearOfBirth: 2018,
		Employees: []*Employee{
			{
				Name:     "John Doe",
				Position: "CEO",
				Wage:     100000,
				JoinedAt: now,
				Projects: []Project{
					{
						Name:       "Project 1",
						Value:      100000,
						StartedAt:  now,
						FinishedAt: now.Add(time.Hour * 24 * 30),
					},
					{
						Name:       "Project 2",
						Value:      200000,
						StartedAt:  now,
						FinishedAt: now.Add(time.Hour * 24 * 30),
					},
				},
			},
			{
				Name:     "Jane Doe",
				Position: "CTO",
				Wage:     80000,
				JoinedAt: now,
				Projects: []Project{
					{
						Name:       "Project 3",
						Value:      300000,
						StartedAt:  now,
						FinishedAt: now.Add(time.Hour * 24 * 30),
					},
					{
						Name:       "Project 4",
						Value:      400000,
						StartedAt:  now,
						FinishedAt: now.Add(time.Hour * 24 * 30),
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

	t := &testing.T{}

	assert.True(t, mutator.SetAddress(nil))
	assert.True(t, mutator.SetAddress(newAddr))
	assert.True(t, mutator.MutateVat().MutateType("Company Ltd."))
	assert.True(t, mutator.MutateAddress().MutateStreet("Baker Street"))
	assert.True(t, mutator.MutateAddress().SetLocation(&uk))
	assert.True(t, mutator.MutateEmployeesAt(0).MutateName("John Smith"))
	assert.True(t, mutator.MutateEmployeesAt(0).MutateProjectsAt(0).MutateName("Project 1 - Updated"))
	assert.False(t, mutator.MutateEmployeesAt(0).MutateName("John Smith"))
	mutator.AppendEmployees(&Employee{
		Name:     "Roger Smith",
		Position: "CFO",
		Wage:     90000,
		JoinedAt: now,
		Projects: nil,
	})
	assert.True(t, mutator.SetNicknames(map[string]*Employee{
		"Johnny": acme.Employees[0],
	}))
	assert.True(t, mutator.InsertNicknames("Janey", acme.Employees[1]))
	assert.False(t, mutator.InsertNicknames("Janey", acme.Employees[1]))
	assert.True(t, mutator.RemoveNicknames("Johnny"))
	assert.False(t, mutator.RemoveNicknames("Roger Ramjet"))
	assert.True(t, mutator.MutateNicknamesWithKey("Janey").MutateWage(50000))
	assert.True(t, mutator.InsertEquity(acme.Employees[1], 1000))

	for _, change := range mutator.FormatChanges() {
		fmt.Println(change)
	}

	assert.Equal(t, 2019, acme.YearOfBirth)
	assert.Equal(t, newAddr, acme.Address)
	assert.Equal(t, "Company Ltd.", acme.Vat.Type)
	assert.Equal(t, "Baker Street", acme.Address.Street)
	assert.Equal(t, "UK", *acme.Address.Location)
	assert.Equal(t, "John Smith", acme.Employees[0].Name)
	assert.Equal(t, "Project 1 - Updated", acme.Employees[0].Projects[0].Name)
	assert.Equal(t, acme.Employees[1], acme.Nicknames["Janey"])
	assert.Nil(t, acme.Nicknames["Johnny"])
	assert.Equal(t, 50000, acme.Employees[1].Wage)
}
