package model

import "{{ .ProjectName }}/common"

type {{ .StructName }}Response struct {
    common.Response
    Data {{ .StructName }}ResponseData
}

type {{ .StructName }}ResponseData struct {
    // Nickname string `json:"nickname" example:"alan"`
}
