package httpfs

import (
	"io/fs"
	"net/http"
)

func (h HTTPFileServer) handleFileRequest(requestedFilePath string) FileServerResponse {
	content, err := fs.ReadFile(h.fs, requestedFilePath)
	res := FileServerResponse{
		Status:  http.StatusOK,
		Content: content,
		Err:     err,
	}
	if res.Err != nil {
		res.Status = http.StatusInternalServerError
		res.Content = getErrorContentForStatus(&res)
	}

	return res
}

func (h HTTPFileServer) handleFileNotExistRequest(requestedFilePath string) FileServerResponse {
	res := FileServerResponse{Status: http.StatusNotFound}

	// Try finding a html file with the same name
	res.Content, _ = fs.ReadFile(h.fs, requestedFilePath+".html")
	if res.Content != nil {
		res.Status = http.StatusOK
	} else {
		res.Content = getErrorContentForStatus(&res)
	}

	return res
}

func (h HTTPFileServer) handleFileForbiddenRequest() FileServerResponse {
	res := FileServerResponse{Status: http.StatusForbidden}
	res.Content = getErrorContentForStatus(&res)

	return res
}

func (h HTTPFileServer) handleOtherErrorRequest(err error) FileServerResponse {
	res := FileServerResponse{Status: http.StatusInternalServerError, Err: err}
	res.Content = getErrorContentForStatus(&res)

	return res
}
