package main

import (
	"fmt"
	"log"

	"controller/utils"
)

func getStdinStr(prompt string) string {
	fmt.Print(prompt)
	var input string
	_, err := fmt.Scanln(&input)
	if err != nil {
		log.Fatal(err)
	}
	return input
}

func main() {
	fmt.Println("What would you like to do?")
	fmt.Println("(1) Create a new instance")
	fmt.Println("(2) Start an instance")
	fmt.Println("(3) Stop an instance")
	fmt.Println("(4) List instances")

	choice := getStdinStr("\nChoice: ")
	customerId := getStdinStr("Customer ID: ")

	switch choice {
	case "1":
		if err := utils.CreateInstance(customerId); err != nil {
			log.Fatal(err)
		}
	case "2":
		if err := utils.StartInstance(customerId); err != nil {
			log.Fatal(err)
		}
	case "3":
		if err := utils.StopInstance(customerId); err != nil {
			log.Fatal(err)
		}
	case "4":
		jails, err := utils.ListInstances(customerId)
		if err != nil {
			log.Fatal(err)
		}
		for _, j := range jails {
			fmt.Println(j)
		}
	default:
		fmt.Println("Invalid choice, try again.")
	}
}
