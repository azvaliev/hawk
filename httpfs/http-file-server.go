package httpfs

import (
	resolvePath "github.com/azvaliev/hawk/resolve-path"
	"io/fs"
	"log"
	"net/http"
)

type HTTPFileServer struct {
	fs fs.FS
}
type FileServerResponse struct {
	Status  int
	Content []byte
	Err     error
}

// GetHTTPFileServer get an instance of HTTPFileServer
// File resolution is relative to the working directory of the fs instance
func GetHTTPFileServer(fs fs.FS) HTTPFileServer {
	return HTTPFileServer{
		fs: fs,
	}
}

func (h HTTPFileServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	fsResponse := h.getResponse(req)

	res.WriteHeader(fsResponse.Status)
	if fsResponse.Err != nil {
		log.Printf("Error: %s\n", fsResponse.Err)
	}
	if fsResponse.Content != nil {
		_, err := res.Write(fsResponse.Content)
		if err != nil {
			log.Printf("Error writing response: %s\n", err)
		}
	}

	log.Printf(
		"%s %s - %d %s\n",
		req.Method,
		req.URL.Path,
		fsResponse.Status,
		http.StatusText(fsResponse.Status),
	)
}

func (h HTTPFileServer) getResponse(req *http.Request) FileServerResponse {
	if req.Method != "GET" {
		return FileServerResponse{
			Status:  http.StatusMethodNotAllowed,
			Content: nil,
			Err:     nil,
		}
	}

	// Remove slash prefix
	requestedPath := req.URL.Path[1:]
	if requestedPath == "" {
		requestedPath = "."
	}
	pathType, err := resolvePath.GetPathType(fs.Stat(h.fs, requestedPath))

	switch pathType {
	case resolvePath.PathIsDir:
		return h.handleDirRequest(requestedPath)
	case resolvePath.PathIsFile:
		return h.handleFileRequest(requestedPath)
	case resolvePath.PathDoesNotExist:
		return h.handleFileNotExistRequest(requestedPath)
	case resolvePath.PathForbidden:
		return FileServerResponse{
			Status:  http.StatusForbidden,
			Content: nil,
			Err:     nil,
		}
	default:
		return FileServerResponse{
			Status:  http.StatusInternalServerError,
			Content: nil,
			Err:     err,
		}
	}

}
