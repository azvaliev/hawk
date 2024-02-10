package main

import (
	"fmt"
	"github.com/azvaliev/hawk/httpfs"
	"github.com/azvaliev/hawk/resolve-path"
	"log"
	"net/http"
	"os"
)

const port = 4580
const serveDir = "static"

func main() {
	fullPath, resolvePathErr := resolvepath.GetDirAbsolutePath(serveDir)
	if resolvePathErr != nil {
		log.Fatalf("Failed to resolve basePath: %s\n", resolvePathErr.Context)
	}

	fileServer := httpfs.GetHTTPFileServer(os.DirFS(fullPath))

	fmt.Printf("Listening on http://localhost:%d\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), fileServer)
	if err != nil {
		log.Fatalf("Fatal Error: %s\n", err)
	}
}
