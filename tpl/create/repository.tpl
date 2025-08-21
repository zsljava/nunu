package repository

import (
    "context"
    "{{ .ProjectName }}/common/base/repository"
	"{{ .ProjectName }}/internal/{{ .BasePkgName}}/model"
)

type {{ .StructName }}Repository interface {
	Get{{ .StructName }}(ctx context.Context, id int64) (*model.{{ .StructName }}, error)
}

func New{{ .StructName }}Repository(
	repository *repository.Repository,
) {{ .StructName }}Repository {
	return &{{ .StructNameLowerFirst }}Repository{
		Repository: repository,
	}
}

type {{ .StructNameLowerFirst }}Repository struct {
	*repository.Repository
}

func (r *{{ .StructNameLowerFirst }}Repository) Get{{ .StructName }}(ctx context.Context, id int64) (*model.{{ .StructName }}, error) {
	var {{ .StructNameLowerFirst }} model.{{ .StructName }}

	return &{{ .StructNameLowerFirst }}, nil
}
