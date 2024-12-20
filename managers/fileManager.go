package managers

import "os"

type FileManager interface {
	SaveFile(string, string, []byte, FileOpeningOptions) (int, error)
	GetFile(string, string, FileOpeningOptions) (*os.File, error)
	ListFiles(string) ([]string, error)
	DeleteFile(string, string) error
	GetFullPath(string, string) string
}

type FileOpeningOptions struct {
	StartPoint DocumentStartPoint
	Offset     uint64
}

func DefaultFileOpeningOptions() FileOpeningOptions {
	return FileOpeningOptions{StartPoint: BEGINNING, Offset: 0}
}

type DocumentStartPoint uint

const (
	BEGINNING DocumentStartPoint = iota
	END
	CUSTOM
)
