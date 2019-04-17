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

func executorAdd(blocks []string) {
	if len(blocks) < 1 {
		fmt.Println("Must select what to add; Options {project, task}!")
		return
	}
	switch blocks[0] {
	case "task":
		executorAddTask(blocks[1:])
	case "project":
		executorAddProject(blocks[1:])
	}
}

func executorAddTask(blocks []string) {
	if len(blocks) < 4 {
		fmt.Println("Usage: \"add task {project} {name} {estimation} {confidence}\"!")
		return
	}
	if median, err := strconv.ParseFloat(blocks[2], 64); err == nil {
		if sigma, err := strconv.ParseFloat(blocks[3], 64); err == nil {
			_, ok := dataStore[blocks[0]]
			if !ok {
				fmt.Printf("No project with name: \"%s\"! Please create one with \"add project\" command!\n", blocks[0])
				return
			}
			dataStore[blocks[0]][blocks[1]] = storage.NewTask(blocks[1], median, sigma)
		}
	}
}

func executorAddProject(blocks []string) {
	if len(blocks) < 1 {
		fmt.Println("Usage: \"add project {name}\"!")
		return
	}
	dataStore[blocks[0]] = storage.NewProject()
}

func executorList(blocks []string) {
	if len(blocks) < 1 {
		fmt.Println("Must select what to list; Options {project, task}!")
		return
	}
	switch blocks[0] {
	case "project":
		executorListProject()
	case "task":
		executorListTask(blocks[1:])
	}
}

func executorListTask(blocks []string) {
	if len(blocks) < 1 {
		fmt.Println("When listing tasks you must define a project!")
		return
	}
	_, ok := dataStore[blocks[0]]
	if !ok {
		fmt.Printf("No project with name: \"%s\"! Please create one with \"add project\" command!\n", blocks[0])
		return
	}
	storage.PrintTasks(dataStore[blocks[0]])
}

func executorListProject() {
	storage.PrintProjects(&dataStore)
}

func executor(in string) {
	in = strings.TrimSpace(in)
	blocks := strings.Split(in, " ")
	switch blocks[0] {
	case "add":
		executorAdd(blocks[1:])
	case "list":
		executorList(blocks[1:])
	}
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
