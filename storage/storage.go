package storage

import (
	"math"
	"os"

	"github.com/jedib0t/go-pretty/table"
	"gonum.org/v1/gonum/stat/distuv"
)

type SPECStorage map[string]SPECProject

type SPECProject map[string]*SPECTask

type SPECTask struct {
	Name       string  `json:"name"`
	Estimation float64 `json:"estimation"`
	Confidence float64 `json:"confidemce"`
	Mean       float64 `json:"mean"`
}

func NewStorage() SPECStorage {
	return make(map[string]SPECProject)
}

func NewProject() SPECProject {
	return make(map[string]*SPECTask)
}

func NewTask(name string, estimation float64, confidemce float64) *SPECTask {
	dist := distuv.LogNormal{math.Log(estimation), confidemce, nil}
	return &SPECTask{name, estimation, confidemce, dist.Mean()}
}

func PrintTasks(tasks map[string]*SPECTask) {
	t := GetTasksTable(tasks)
	t.Render()
}

func GetTasksTable(tasks map[string]*SPECTask) table.Writer {
	t := getTasksTableHeader()
	for k := range tasks {
		t.AppendRow([]interface{}{
			tasks[k].Name,
			tasks[k].Estimation,
			tasks[k].Confidence,
			floatRound(tasks[k].Mean),
		})
	}
	return t
}

func getTasksTableHeader() table.Writer {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Name", "Estimation", "Confidence", "Mean"})
	return t
}

func PrintProjects(dataStore *SPECStorage) {
	t := GetProjectsTable(dataStore)
	t.Render()
}

func GetProjectsTable(dataStore *SPECStorage) table.Writer {
	t := getProjectsTableHeader()
	var meanSum float64
	for k := range *dataStore {
		meanSum = 0
		for j := range (*dataStore)[k] {
			meanSum += floatRound((*dataStore)[k][j].Mean)
		}
		t.AppendRow([]interface{}{
			k,
			len((*dataStore)[k]),
			meanSum,
		})
	}
	return t
}

func getProjectsTableHeader() table.Writer {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Name", "# of Tasks", "Mean Sum"})
	return t
}

func floatRound(number float64) float64 {
	return math.Round(number*100) / 100
}
