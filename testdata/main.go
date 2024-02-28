package main

import (
	"fmt"
	"log"
	"time"
)

func assertBool(expected bool, obtained bool) {
	if expected != obtained {
		log.Fatalf("expected %+v", expected)
	}
}

func assertEqual[T comparable](expected T, obtained T) {
	if expected != obtained {
		log.Fatalf("expected %+v, got %+v", expected, obtained)
	}
}

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
	mutator.SetYearOfBirth(2019)

	newAddr := &Address{
		Street: "Liverpool Street",
		Number: 300,
		City:   "London",
		Zip:    45001,
	}

	uk := "UK"

	assertBool(true, mutator.SetAddress(nil))
	assertBool(true, mutator.SetAddress(newAddr))
	assertBool(true, mutator.Vat().SetType("Company Ltd."))
	assertBool(true, mutator.Address().SetStreet("Baker Street"))
	assertBool(true, mutator.Address().SetLocation(&uk))
	assertBool(true, mutator.EmployeesAt(0).SetName("John Smith"))
	assertBool(true, mutator.EmployeesAt(0).ProjectsAt(0).SetName("Project 1 - Updated"))
	assertBool(true, mutator.EmployeesAt(0).ProjectsAt(0).SetSeqID([]byte("1234567890")))
	assertBool(false, mutator.EmployeesAt(0).ProjectsAt(0).SetSeqID([]byte("1234567890")))
	assertBool(true, mutator.EmployeesAt(0).ProjectsAt(0).SetSeqID([]byte("123456789")))
	assertBool(true, mutator.EmployeesAt(0).ProjectsAt(1).SetSeqID([]byte("1234")))
	assertBool(true, mutator.EmployeesAt(0).ProjectsAt(1).SetSeqID(nil))
	assertBool(false, mutator.EmployeesAt(0).SetName("John Smith"))
	assertBool(true, mutator.EmployeesByPtr(acme.Employees[1]) != nil)
	assertBool(true, mutator.EmployeesByPtr(&Employee{}) == nil)
	mutator.AppendEmployees(&Employee{
		Name:     "Roger Smith",
		Position: "CFO",
		Wage:     90000,
		JoinedAt: now,
		Projects: nil,
	})
	assertBool(true, mutator.SetNicknames(map[string]*Employee{
		"Johnny": acme.Employees[0],
	}))
	assertBool(true, mutator.InsertNicknames("Janey", acme.Employees[1]))
	assertBool(false, mutator.InsertNicknames("Janey", acme.Employees[1]))
	assertBool(true, mutator.RemoveNicknames("Johnny"))
	assertBool(false, mutator.RemoveNicknames("Roger Ramjet"))
	assertBool(true, mutator.NicknamesWithKey("Janey").SetWage(50000))
	assertBool(true, mutator.InsertEquity(acme.Employees[1], 1000))

	for _, change := range mutator.FormatChanges() {
		fmt.Println(change)
	}

	assertEqual(2019, acme.YearOfBirth)
	assertEqual(newAddr, acme.Address)
	assertEqual("Company Ltd.", acme.Vat.Type)
	assertEqual("Baker Street", acme.Address.Street)
	assertEqual("UK", *acme.Address.Location)
	assertEqual("John Smith", acme.Employees[0].Name)
	assertEqual("Project 1 - Updated", acme.Employees[0].Projects[0].Name)
	assertEqual(acme.Employees[1], acme.Nicknames["Janey"])
	assertBool(true, acme.Nicknames["Johnny"] == nil)
	assertEqual(50000, acme.Employees[1].Wage)
}
