package main

import (
	"bufio"
	"controller/utils"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var stdin = bufio.NewReader(os.Stdin)

func getStdinStr(prompt string) (string, error) {
	fmt.Print(prompt)
	line, err := stdin.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimRight(line, "\r\n"), nil
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	username, err := utils.GitHubAuth()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nHey %s!\nWelcome to Eigen Cloud\n\n", username)

	fmt.Println("What would you like to do?")
	fmt.Println("(1) Create a new instance")
	fmt.Println("(2) Start an instance")
	fmt.Println("(3) Stop an instance")
	fmt.Println("(4) List instances")

	choice, err := getStdinStr("\nChoice: ")
	if err != nil {
		log.Fatal(err)
	}

	switch choice {
	case "1":
		if err := utils.CreateInstance(username); err != nil {
			log.Fatal(err)
		}
	case "2":
		if err := utils.StartInstance(username); err != nil {
			log.Fatal(err)
		}
	case "3":
		if err := utils.StopInstance(username); err != nil {
			log.Fatal(err)
		}
	case "4":
		jails, err := utils.ListInstances(username)
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
