package task

import (
	"context"
	"{{ .ProjectName }}/internal/{{ .BasePkgName}}/repository"
	"time"
)

type {{ .StructName }}Task interface {
	Execute(ctx context.Context) error
}

func New{{ .StructName }}Task(
	task *Task,
	{{ .StructNameLowerFirst }}Repo repository.{{ .StructName }}Repository,
) {{ .StructName }}Task {
	return &{{ .StructNameLowerFirst }}Task{
		{{ .StructNameLowerFirst }}Repo: {{ .StructNameLowerFirst }}Repo,
		Task: task,
	}
}

type {{ .StructNameLowerFirst }}Task struct {
	{{ .StructNameLowerFirst }}Repo repository.{{ .StructName }}Repository
	*Task
}

func (t {{ .StructNameLowerFirst }}Task) Execute(ctx context.Context) error {
	// do something
	for {
		t.logger.Info("execute task")
	}
}
