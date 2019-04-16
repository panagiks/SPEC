package main

import (
	"fmt"
	"strconv"
	"strings"

	prompt "github.com/c-bata/go-prompt"
	"github.com/panagiks/SPEC/storage"
)

var (
	dataStore storage.SPECStorage
)

func executor(in string) {
	in = strings.TrimSpace(in)
	blocks := strings.Split(in, " ")
	switch blocks[0] {
	case "add":
		if len(blocks) < 2 {
			fmt.Println("Must select what to add; Options {project, task}!")
			return
		}
		switch blocks[1] {
		case "task":
			if len(blocks) < 6 {
				fmt.Println("Usage: \"add task {project} {name} {estimation} {confidence}\"!")
				return
			}
			if median, err := strconv.ParseFloat(blocks[4], 64); err == nil {
				if sigma, err := strconv.ParseFloat(blocks[5], 64); err == nil {
					_, ok := dataStore[blocks[2]]
					if !ok {
						fmt.Printf("No project with name: \"%s\"! Please create one with \"add project\" command!\n", blocks[2])
						return
					}
					dataStore[blocks[2]][blocks[3]] = storage.NewTask(blocks[3], median, sigma)
				}
			}
			return
		case "project":
			if len(blocks) < 3 {
				fmt.Println("Usage: \"add project {name}\"!")
				return
			}
			dataStore[blocks[2]] = storage.NewProject()
		}
	case "list":
		if len(blocks) < 2 {
			fmt.Println("Must select what to list; Options {project, task}!")
			return
		}
		switch blocks[1] {
		case "project":
			fmt.Println("Projects:")
			for k := range dataStore {
				fmt.Printf("\t%s\n", k)
			}
		case "task":
			if len(blocks) < 3 {
				fmt.Println("When listing tasks you must define a project!")
				return
			}
			_, ok := dataStore[blocks[2]]
			if !ok {
				fmt.Printf("No project with name: \"%s\"! Please create one with \"add project\" command!\n", blocks[2])
				return
			}
			storage.PrintTasks(dataStore[blocks[2]])
		}
	}
	return
}

func completer(t prompt.Document) []prompt.Suggest {
	return []prompt.Suggest{
		{Text: "add"},
		{Text: "add project"},
		{Text: "add task"},
		{Text: "list"},
		{Text: "list task"},
		{Text: "list project"},
	}
}

func main() {
	dataStore = storage.NewStorage()
	fmt.Println("Please use `exit` or `Ctrl-D` to exit this program.")
	defer fmt.Println("Bye!")
	p := prompt.New(
		executor,
		completer,
		prompt.OptionTitle("SPEC: Software Project Estimation Calculator"),
	)
	p.Run()
}
