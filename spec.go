package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"

	prompt "github.com/c-bata/go-prompt"
	"github.com/panagiks/SPEC/storage"
)

type cliState struct {
	executor           func(in []string)
	completer          func(t prompt.Document) []prompt.Suggest
	LivePrefix         string
	IsLivePrefixEnable bool
	project            string
	task               string
}

var (
	dataStore storage.SPECStorage
	cli       []cliState
)

func updateTask(project string, name string, actual string) {
	if actualF, err := strconv.ParseFloat(actual, 64); err == nil {
		_, ok := dataStore[project]
		if !ok {
			fmt.Printf("No project with name: \"%s\"! Please create one!\n", project)
			return
		}
		_, ok = dataStore[project][name]
		if !ok {
			fmt.Printf("No task with name: \"%s\"! Please create one with \"add\" command!\n", name)
			return
		}
		dataStore[project][name].Actual = actualF
	}
}

func addTask(project string, name string, median string, sigma string) {
	if medianF, err := strconv.ParseFloat(median, 64); err == nil {
		if sigmaF, err := strconv.ParseFloat(sigma, 64); err == nil {
			_, ok := dataStore[project]
			if !ok {
				fmt.Printf("No project with name: \"%s\"! Please create one with \"add\" command!\n", project)
				return
			}
			dataStore[project][name] = storage.NewTask(name, medianF, sigmaF)
		}
	}
}

func listTask(project string) {
	_, ok := dataStore[project]
	if !ok {
		fmt.Printf("No project with name: \"%s\"! Please create one!\n", project)
		return
	}
	storage.PrintTasks(dataStore[project])
}

func addProject(name string) {
	dataStore[name] = storage.NewProject()
}

func listProject() {
	dataStore.PrintProjects()
}

func selectTask(project string, name string) {
	_, ok := dataStore[project][name]
	if !ok {
		fmt.Printf("No task with name: \"%s\"! Please create one with \"add\" command!\n", name)
		return
	}
	cli = append(
		cli,
		cliState{
			executor:           taskExecutor,
			completer:          taskCompleter,
			IsLivePrefixEnable: true,
			LivePrefix:         fmt.Sprintf("[P:%s][T:%s] >", project, name),
			project:            project,
			task:               name,
		},
	)
}

func selectProject(project string) {
	_, ok := dataStore[project]
	if !ok {
		fmt.Printf("No project with name: \"%s\"! Please create one with \"add\" command!\n", project)
		return
	}
	cli = append(
		cli,
		cliState{
			executor:           projectExecutor,
			completer:          projectCompleter,
			IsLivePrefixEnable: true,
			LivePrefix:         fmt.Sprintf("[P:%s] >", project),
			project:            project,
		},
	)
}

func taskExecutor(blocks []string) {
	switch blocks[0] {
	case "actual":
		if len(blocks) < 2 {
			fmt.Println("Usage: \"actual {duration}\"!")
			return
		}
		updateTask(cli[len(cli)-1].project, cli[len(cli)-1].task, blocks[1])
	}
}

func taskCompleter(t prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "actual", Description: "Set the selected task's actual duration; Syntax: \"actual {duration}\""},
	}
	return s
}

func projectExecutor(blocks []string) {
	switch blocks[0] {
	case "add", "a":
		if len(blocks) < 4 {
			fmt.Println("Usage: \"add {name} {estimation} {confidence}\"!")
			return
		}
		addTask(cli[len(cli)-1].project, blocks[1], blocks[2], blocks[3])
	case "list", "l":
		listTask(cli[len(cli)-1].project)
	default:
		selectTask(cli[len(cli)-1].project, blocks[0])
	}
}

func projectCompleter(t prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "add", Description: "Add a new task; Syntax: \"add {taskName} {estimation} {confidence}\""},
		{Text: "list", Description: "List tasks within the project"},
	}
	for k := range dataStore[cli[len(cli)-1].project] {
		s = append(s, prompt.Suggest{Text: k, Description: "Select task for further options"})
	}
	return s
}

func baseExecutor(blocks []string) {
	switch blocks[0] {
	case "add", "a":
		if len(blocks) < 2 {
			fmt.Println("Usage: \"add {projectName}\"!")
			return
		}
		addProject(blocks[1])
	case "list", "l":
		listProject()
	case "save":
		if len(blocks) < 2 {
			fmt.Println("Usage: \"save {pathToFile}\"!")
			return
		}
		dataStore.Save(blocks[1])
	case "load":
		if len(blocks) < 2 {
			fmt.Println("Usage: \"load {pathToFile}\"!")
			return
		}
		dataStore.Load(blocks[1])
	default:
		selectProject(blocks[0])
	}
}

func baseCompleter(t prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "add", Description: "Add a new project; Syntax: \"add {projectName}\""},
		{Text: "list", Description: "List projects"},
		{Text: "save", Description: "Save current state to a file; Syntax: \"save {pathToFile}\""},
		{Text: "load", Description: "Load state from a file; Syntax: \"load {pathToFile}\""},
	}
	for k := range dataStore {
		s = append(s, prompt.Suggest{Text: k, Description: "Select project for further options"})
	}
	return s
}

func goBack() {
	if len(cli) == 1 {
		runtime.Goexit()
	}
	cli = cli[:len(cli)-1]
}

func executor(in string) {
	in = strings.TrimSpace(in)
	if in == "" {
		goBack()
		return
	}
	blocks := strings.Split(in, " ")
	cli[len(cli)-1].executor(blocks)
}

func completer(t prompt.Document) []prompt.Suggest {
	s := cli[len(cli)-1].completer(t)
	return prompt.FilterHasPrefix(s, t.GetWordBeforeCursor(), true)
}

func changeLivePrefix() (string, bool) {
	return cli[len(cli)-1].LivePrefix, cli[len(cli)-1].IsLivePrefixEnable
}

func main() {
	dataStore = storage.NewStorage()
	cli = append(
		cli,
		cliState{executor: baseExecutor, completer: baseCompleter},
	)
	fmt.Println("Please use `exit` or `Ctrl-D` to exit this program.")
	defer os.Exit(0)
	defer fmt.Println("Bye!")
	p := prompt.New(
		executor,
		completer,
		prompt.OptionTitle("SPEC: Software Project Estimation Calculator"),
		prompt.OptionPrefix("> "),
		prompt.OptionLivePrefix(changeLivePrefix),
	)
	p.Run()
}
