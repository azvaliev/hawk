package httpfs

import (
	"bytes"
	"embed"
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

// printDirEntryInfo Information necessary for printing a directory entry
type printDirEntryInfo struct {
	Name       string
	Path       string
	ModifiedAt string
	IsDir      bool
	Size       string
}

type directoryTemplateParams struct {
	Title   string
	Entries []*printDirEntryInfo
}

//go:embed directory.gohtml
var dirTemplateEmbed embed.FS

func (h HTTPFileServer) handleDirRequest(requestedFilePath string) FileServerResponse {
	res := FileServerResponse{Status: http.StatusOK}
	dirEntries, err := fs.ReadDir(h.fs, requestedFilePath)
	if err == nil {
		res.Content, res.Err = h.getDirContent(requestedFilePath, dirEntries)
	} else {
		res.Status = 500
		res.Err = err
	}

	return res
}

func (h HTTPFileServer) getDirContent(dirPath string, dirEntries []fs.DirEntry) (content []byte, err error) {
	entriesWithMeta := make([]*printDirEntryInfo, len(dirEntries))

	// Get info for each entry
	for idx, dirEntry := range dirEntries {
		info, err := dirEntry.Info()
		if err != nil {
			return nil, err
		}

		if strings.ToLower(dirEntry.Name()) == "index.html" {
			file, err := fs.ReadFile(h.fs, filepath.Join(dirPath, dirEntry.Name()))
			return file, err
		}
		entriesWithMeta[idx] = &printDirEntryInfo{
			Name:       dirEntry.Name(),
			Path:       filepath.Join(dirPath, dirEntry.Name()),
			ModifiedAt: info.ModTime().Format(time.RFC3339),
			IsDir:      info.IsDir(),
			Size:       fileSizeFormat(info.Size()),
		}
	}

	// Fill out the template
	t, err := template.ParseFS(dirTemplateEmbed, "directory.gohtml")
	if err != nil {
		return nil, err
	}

	var contentBuffer bytes.Buffer
	err = t.Execute(&contentBuffer, directoryTemplateParams{
		Title:   dirPath,
		Entries: entriesWithMeta,
	})
	if err != nil {
		return nil, err
	}

	return contentBuffer.Bytes(), nil
}
