package resolvepath

import (
	"os"
	"path/filepath"
)

type GetDirAbsolutePathError struct {
	Context string
	Err     error
}

func GetDirAbsolutePath(path string) (string, *GetDirAbsolutePathError) {
	var absolutePath string
	if filepath.IsAbs(path) {
		absolutePath = path
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			return "", &GetDirAbsolutePathError{
				Context: "Failed to get CWD",
				Err:     err,
			}
		}
		absolutePath = filepath.Join(cwd, path)
	}

	pathType, err := GetPathType(os.Stat(absolutePath))

	switch pathType {
	case PathIsFile:
		return "", &GetDirAbsolutePathError{
			Context: "Path is not a directory",
			Err:     err,
		}
	case PathDoesNotExist:
		return "", &GetDirAbsolutePathError{
			Context: "Path does not exist",
			Err:     err,
		}
	case PathForbidden:
		return "", &GetDirAbsolutePathError{
			Context: "Insufficient permission to access the path",
			Err:     err,
		}
	case PathUnknown:
		return "", &GetDirAbsolutePathError{
			Context: "Unknown issue with accessing path",
			Err:     err,
		}
	default:
		return absolutePath, nil
	}

}
