package taskservice

import (
	"context"
	"time"

	"github.com/powerman/structlog"
)

type Task struct {
	Id        string
	Author    string
	Assignee  string
	Title     string
	Boby      string
	Url       string
	StartDate time.Time
}
type TaskBot interface {
	GetTasksLoop(ctx context.Context, inTasksChannel chan Task, log *structlog.Logger)
}
