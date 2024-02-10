package resolvepath

import (
	"errors"
	"io/fs"
	"os"
)

type PathType int

const (
	PathIsDir PathType = iota
	PathIsFile
	PathDoesNotExist
	PathForbidden
	PathUnknown
)

// GetPathType Map a path to a PathType, forward any errors encountered in the process.
// Call os.Stat or resolve-path.Stat for parameters
func GetPathType(info os.FileInfo, err error) (PathType, error) {
	if err == nil {
		if info.IsDir() {
			return PathIsDir, nil
		}
		return PathIsFile, nil
	}
	if errors.Is(err, fs.ErrNotExist) {
		return PathDoesNotExist, err
	}
	if errors.Is(err, fs.ErrPermission) {
		return PathForbidden, err
	}

	return PathUnknown, err
}
