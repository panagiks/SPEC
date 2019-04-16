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
			math.Round(tasks[k].Mean*100) / 100,
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
