package implementations

import (
	"echo-api/managers"
	"echo-api/models/entities"
	"errors"
	"fmt"
	"io"
	"strings"
)

/* This implementation is for locally generating prompts
 * I am planning to create another go module to do this seperately then dev another implementation of promptManager that interacts with that API
 */
type LocalPromptGenManager struct {
	fileManager managers.FileManager
}

func NewLocalPromptGenManager(fm managers.FileManager) LocalPromptGenManager {
	return LocalPromptGenManager{fileManager: fm}
}

func (m LocalPromptGenManager) GeneratePrompt(val any) (string, error) {
	switch val := val.(type) {
	case entities.Note:
		return m.generatePromptForNote(val)
	case entities.Document:
		return m.generatePromptForDocument(val)
	case string:
		return m.promptizeString(managers.Remember, val), nil
	default:
		return "", errors.ErrUnsupported
	}
}
func (m LocalPromptGenManager) GeneratePromptWith(val any, action managers.PromptAction) (string, error) {
	switch val := val.(type) {
	case entities.Note:
		return m.generatePromptForNote(val)
	case entities.Document:
		return m.generatePromptForDocument(val)
	case string:
		return m.promptizeString(managers.Remember, val), nil
	default:
		return "", errors.ErrUnsupported
	}
}

func (m LocalPromptGenManager) GenerateMessage(val string) (string, error) {
	return m.messageizeString(val), nil
}

func (m LocalPromptGenManager) generatePromptForNote(val entities.Note) (string, error) {
	return m.GeneratePrompt(val.Payload)
}

func (m LocalPromptGenManager) generatePromptForDocument(val entities.Document) (string, error) {
	f, err := m.fileManager.GetFile(val.Location, val.Name, managers.DefaultFileOpeningOptions())
	if err != nil {
		return "", err
	}
	defer f.Close()
	const buffSize = 1024
	buffer := make([]byte, buffSize)
	var sb strings.Builder
	for {
		readTotal, err := f.Read(buffer)
		if err != nil {
			if err == io.EOF {
				sb.Write(buffer[:readTotal])
				break
			}
			return "", err
		}
		sb.Write(buffer[:readTotal])
	}

	return m.GeneratePrompt(sb.String())
}

func (m LocalPromptGenManager) messageizeString(s string) string {
	return fmt.Sprintf(managers.Message.String(), s)
}

func (m LocalPromptGenManager) promptizeString(act managers.PromptAction, s string) string {
	return fmt.Sprintf(managers.Prompt.String(), fmt.Sprintf(act.String(), s))
}
