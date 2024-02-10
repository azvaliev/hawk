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
	errLog := log.New(os.Stderr, "", 0)

	// Get the directory from args
	fullPath, resolvePathErr := resolvepath.GetDirAbsolutePath(args.Dir)
	if resolvePathErr != nil {
		errLog.Fatalf("Failed to resolve path \"%s\": %s\n", args.Dir, resolvePathErr.Context)
	}
	fileServer := httpfs.GetHTTPFileServer(os.DirFS(fullPath))

	// Setup listener on the port
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", args.Port))
	if err != nil {
		errLog.Fatalf("Failed to start server: %s", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	fmt.Printf("Listening on http://localhost:%d\n", port)

	// Start the static file server
	err = http.Serve(listener, fileServer)
	if err != nil {
		errLog.Fatalf("Fatal Error: %s\n", err)
	}
}
