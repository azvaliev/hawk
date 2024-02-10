package httpfs

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"net/http"
)

//go:embed error.gohtml
var errorTemplateEmbed embed.FS

type errorTemplateParams struct {
	Status            int
	StatusMessage     string
	AdditionalMessage string
}

func getErrorContentForStatus(res *FileServerResponse) []byte {
	if res.Status < 400 {
		return res.Content
	}

	templateParams := errorTemplateParams{
		Status:        res.Status,
		StatusMessage: http.StatusText(res.Status),
	}

	switch res.Status {
	case http.StatusNotFound:
		templateParams.AdditionalMessage = "These aren't the droids you're looking for"
	case http.StatusForbidden:
		templateParams.AdditionalMessage = "You don't have permission to access this resource"
	default:
		if res.Err == nil {
			templateParams.AdditionalMessage = "An unknown error occurred"
		} else {
			templateParams.AdditionalMessage = fmt.Sprint(res.Err)
		}
	}

	t, err := template.ParseFS(errorTemplateEmbed, "error.gohtml")
	if err != nil {
		return res.Content
	}

	var contentBuffer bytes.Buffer
	err = t.Execute(&contentBuffer, templateParams)

	if err == nil {
		return contentBuffer.Bytes()
	}
	return res.Content
}
