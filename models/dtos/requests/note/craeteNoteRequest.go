package note

type CreateNoteRequest struct {
	Header     string  `json:"header"`
	Payload    string  `json:"payload"`
	LanguageID string  `json:"languageId"`
	UserID     *string `json:"userId"`
	ContextID  string  `json:"contextId"`
}
