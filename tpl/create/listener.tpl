package listener

import (
	"context"
	"{{ .ProjectName }}/internal/{{ .BasePkgName}}/repository"
	"time"
)

type {{ .StructName }}Listener interface {
	KafkaConsumer(ctx context.Context) error
}

func New{{ .StructName }}Listener(
	listener *Listener,
	{{ .StructNameLowerFirst }}Repo repository.{{ .StructName }}Repository,
) {{ .StructName }}Listener {
	return &{{ .StructNameLowerFirst }}Listener{
		{{ .StructNameLowerFirst }}Repo: {{ .StructNameLowerFirst }}Repo,
		Listener: listener,
	}
}

type {{ .StructNameLowerFirst }}Listener struct {
	{{ .StructNameLowerFirst }}Repo repository.{{ .StructName }}Repository
	*Listener
}

func (t {{ .StructNameLowerFirst }}Listener) KafkaConsumer(ctx context.Context) error {
	// do something
	for {
		t.logger.Info("KafkaConsumer")
		time.Sleep(time.Second * 5)
	}
}
