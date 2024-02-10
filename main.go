package main

import (
	"fmt"
	"github.com/azvaliev/hawk/cli"
	"github.com/azvaliev/hawk/httpfs"
	"github.com/azvaliev/hawk/resolve-path"
	"log"
	"net"
	"net/http"
	"os"
)

func main() {
	args := cli.MustParseCLIArgs()

	fullPath, resolvePathErr := resolvepath.GetDirAbsolutePath(args.Dir)
	if resolvePathErr != nil {
		log.Fatalf("Failed to resolve basePath: %s\n", resolvePathErr.Context)
	}
	fileServer := httpfs.GetHTTPFileServer(os.DirFS(fullPath))

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", args.Port))
	if err != nil {
		log.Fatalf("Failed to start server: %s", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	fmt.Printf("Listening on http://localhost:%d\n", port)

	err = http.Serve(listener, fileServer)
	if err != nil {
		log.Fatalf("Fatal Error: %s\n", err)
	}
}
