package main

import (
	"fmt"

	"github.com/billcobbler/heb-to-go/pkg/heb-api/v1"
)

func main() {
	heb := heb.NewClient()
	stores, err := heb.LocateStores("78741", 3)
	if err != nil {
		fmt.Println(err)
	}
	for _, s := range stores {
		fmt.Printf("ID: %s, Name: %s, Zip: %s\n", s.ID, s.Name, s.PostalCode)
		fmt.Println("==========================")
		timeslots, err := heb.GetStoreTimeslots(s.ID)
		if err != nil {
			fmt.Println(err)
		}
		for _, t := range timeslots {
			fmt.Printf("Date: %s, Type: %s, Time: %s\n", t.Date, t.FullfillmentType, t.StartTime)
		}
		fmt.Println()
	}

}
