package domain

import (
    "context"
    "{{ .ProjectName }}/common/base/service"
	"{{ .ProjectName }}/internal/{{ .BasePkgName}}/model"
	"{{ .ProjectName }}/internal/{{ .BasePkgName}}/repository"
)

type {{ .StructName }}DomainService interface {
	Get{{ .StructName }}(ctx context.Context, id int64) (*model.{{ .StructName }}, error)
}
func New{{ .StructName }}DomainService(
    service *Service,
    {{ .StructNameLowerFirst }}Repository repository.{{ .StructName }}Repository,
) {{ .StructName }}Service {
	return &{{ .StructNameLowerFirst }}Service{
		Service:        service,
		{{ .StructNameLowerFirst }}Repository: {{ .StructNameLowerFirst }}Repository,
	}
}

type {{ .StructNameLowerFirst }}DomainService struct {
	*Service
	{{ .StructNameLowerFirst }}Repository repository.{{ .StructName }}Repository
}

func (s *{{ .StructNameLowerFirst }}DomainService) Get{{ .StructName }}(ctx context.Context, id int64) (*model.{{ .StructName }}, error) {
	return s.{{ .StructNameLowerFirst }}Repository.Get{{ .StructName }}(ctx, id)
}
