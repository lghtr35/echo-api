package implementations

import (
	"echo-api/managers"
	"errors"
	"fmt"
	"os"
)

type OnServerFileManager struct {
	basePath    string
	locationMap map[string]string
}

func NewOnServerFileManager(basePath string, locations []string) *OnServerFileManager {
	locs := make(map[string]string, 0)
	for _, i := range locations {
		locs[i] = fmt.Sprintf("%s/%s", basePath, i)
	}
	return &OnServerFileManager{basePath: basePath, locationMap: locs}

}

func (m *OnServerFileManager) SaveFile(location string, filename string, buffer []byte, options managers.FileOpeningOptions) (int, error) {
	fullpath := m.GetFullPath(location, filename)
	f, err := m.openFileToWrite(fullpath, options)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	count, err := f.Write(buffer)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (m *OnServerFileManager) GetFile(location string, filename string, options managers.FileOpeningOptions) (*os.File, error) {
	fullpath := m.GetFullPath(location, filename)
	f, err := m.openFileToRead(fullpath, options)
	if err != nil {
		return nil, err
	}
	return f, nil
}
func (m *OnServerFileManager) ListFiles(fileType string) ([]string, error) {
	fullpath := m.GetFullPath(fileType, "")
	entries, err := os.ReadDir(fullpath)
	if err != nil {
		return make([]string, 0), err
	}
	res := make([]string, len(entries))
	for i, v := range entries {
		res[i] = v.Name()
	}
	return res, err
}
func (m *OnServerFileManager) DeleteFile(location string, filename string) error {
	fullpath := m.GetFullPath(location, filename)
	err := os.Remove(fullpath)
	if err != nil {
		return err
	}
	return nil
}

func (m *OnServerFileManager) GetFullPath(location, filename string) string {
	return fmt.Sprintf("%s/%s/%s", m.basePath, m.locationMap[location], filename)
}

func (m *OnServerFileManager) openFileToWrite(fullpath string, options managers.FileOpeningOptions) (*os.File, error) {
	point := options.StartPoint
	if point == managers.END {
		f, err := os.OpenFile(fullpath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
		return f, nil
	} else if point == managers.BEGINNING {
		f, err := os.OpenFile(fullpath, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
		return f, nil
	} else if point == managers.CUSTOM {
		f, err := os.OpenFile(fullpath, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
		f.Seek(int64(options.Offset), 0)
		return f, nil
	}
	return nil, errors.New("argumentErrorUnknownStartPoint")
}

func (m *OnServerFileManager) openFileToRead(fullpath string, options managers.FileOpeningOptions) (*os.File, error) {
	point := options.StartPoint
	if point == managers.END {
		return nil, errors.ErrUnsupported
	} else if point == managers.BEGINNING {
		f, err := os.OpenFile(fullpath, os.O_RDONLY, 0644)
		if err != nil {
			return nil, err
		}
		return f, nil
	} else if point == managers.CUSTOM {
		f, err := os.OpenFile(fullpath, os.O_RDONLY, 0644)
		if err != nil {
			return nil, err
		}
		f.Seek(int64(options.Offset), 0)
		return f, nil
	}
	return nil, errors.New("argumentErrorUnknownStartPoint")
}
