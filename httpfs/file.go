package httpfs

import (
	"io/fs"
	"net/http"
)

func (h HTTPFileServer) handleFileRequest(requestedFilePath string) FileServerResponse {
	content, err := fs.ReadFile(h.fs, requestedFilePath)
	status := http.StatusOK
	if err != nil {
		status = http.StatusInternalServerError
	}

	return FileServerResponse{
		Status:  status,
		Content: content,
		Err:     err,
	}
}

func (h HTTPFileServer) handleFileNotExistRequest(originalRequestFilePath string) FileServerResponse {
	status := http.StatusNotFound
	// Try finding an html file with the same name
	content, err := fs.ReadFile(h.fs, originalRequestFilePath+".html")
	if err == nil {
		status = http.StatusOK
	}

	return FileServerResponse{
		Status:  status,
		Content: content,
		Err:     nil,
	}
}
