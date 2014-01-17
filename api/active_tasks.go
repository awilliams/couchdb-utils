package api

import (
	"fmt"
	"github.com/awilliams/couchdb-utils/util"
)

type ActiveTask struct {
	Type           string
	Pid            string
	Database       Database
	Progress       int
	DesignDocument DesignDoc `json:"design_document"`
	StartedOn      int       `json:"started_on"`
	Source         string
	Target         string
	Continuous     bool
}

func (a ActiveTask) PP(printer util.Printer) {
	progress := fmt.Sprintf("%02d%%", a.Progress)
	var addInfo string
	switch a.Type {
	case "replication":
		addInfo = fmt.Sprintf("%s → %s", a.Source, a.Target)
		if a.Continuous {
			addInfo += " (continuous)"
		}
	case "indexer":
		addInfo = fmt.Sprintf("%s/%s", a.Database.String(), a.DesignDocument.String())
	default:
		addInfo = a.Database.String()
	}
	printer.Print("[%s %s]\n %s", progress, a.Type, addInfo)
}

type ActiveTasks []ActiveTask

func (a *ActiveTasks) path() string {
	return "_active_tasks"
}

func (a *ActiveTasks) ByType(t string) ActiveTasks {
	var filtered ActiveTasks
	for _, activeTask := range *a {
		if activeTask.Type == t {
			filtered = append(filtered, activeTask)
		}
	}
	return filtered
}

func (a ActiveTasks) PP(printer util.Printer) {
	for _, activeTask := range a {
		activeTask.PP(printer)
	}
}

func (c *Couchdb) GetActiveTasks() (ActiveTasks, error) {
	activeTasks := new(ActiveTasks)
	err := c.getJson(activeTasks, activeTasks.path())
	return *activeTasks, err
}
