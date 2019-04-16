package storage

import (
	"math"

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
