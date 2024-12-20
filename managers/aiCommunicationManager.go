package managers

type AiCommunicationManager interface {
	SendPrompt(string, string) (string, error)
	ResetContext(string, bool) error
	DeleteContext(string, bool) error
	CreateContext(string, bool) error
}
