package model

import "github.com/zsljava/gokit/common/response"

type {{ .StructName }}Response struct {
    response.Response
    Data {{ .StructName }}ResponseData
}

type {{ .StructName }}ResponseData struct {
    // Nickname string `json:"nickname" example:"alan"`
}
