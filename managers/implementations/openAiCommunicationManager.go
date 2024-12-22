package implementations

import (
	"echo-api/util"
	"errors"
)

// TODO Implement
type OpenAiCommunicationManager struct {
	configuration *util.Configuration
}

func NewOpenAiCommunicationManager(c *util.Configuration) OpenAiCommunicationManager {
	return OpenAiCommunicationManager{configuration: c}
}

func (cm OpenAiCommunicationManager) SendPrompt(contextID string, msg string) (string, error) {
	if cm.configuration.IsAiAssistantEnabled {
		return "", errors.ErrUnsupported
	}
	return msg, nil
}
func (cm OpenAiCommunicationManager) ResetContext(contextID string, isSoftReset bool) error {
	if cm.configuration.IsAiAssistantEnabled {
		return errors.ErrUnsupported
	}
	return nil
}

func (cm OpenAiCommunicationManager) DeleteContext(string, bool) error {
	if cm.configuration.IsAiAssistantEnabled {
		return errors.ErrUnsupported
	}
	return nil
}
func (cm OpenAiCommunicationManager) CreateContext(string, bool) error {
	if cm.configuration.IsAiAssistantEnabled {
		return errors.ErrUnsupported
	}
	return nil
}
