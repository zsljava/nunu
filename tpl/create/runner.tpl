package runner

import (
	"context"
	"{{ .ProjectName }}/internal/{{ .BasePkgName}}/repository"
	"time"
)

type {{ .StructName }}Runner interface {
	Execute(ctx context.Context) error
}

func New{{ .StructName }}Runner(
	runner *Runner,
	{{ .StructNameLowerFirst }}Repo repository.{{ .StructName }}Repository,
) {{ .StructName }}Runner {
	return &{{ .StructNameLowerFirst }}Runner{
		{{ .StructNameLowerFirst }}Repo: {{ .StructNameLowerFirst }}Repo,
		Runner: runner,
	}
}

type {{ .StructNameLowerFirst }}Runner struct {
	{{ .StructNameLowerFirst }}Repo repository.{{ .StructName }}Repository
	*Runner
}

func (t {{ .StructNameLowerFirst }}Runner) Execute(ctx context.Context) error {
	// do something
	for {
		t.logger.Info("execute Runner")
		time.Sleep(time.Second * 5)
	}
}
