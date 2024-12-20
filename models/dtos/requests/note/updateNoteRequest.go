package note

type UpdateNoteRequest struct {
	ID         string  `json:"id"`
	Header     *string `json:"header"`
	Payload    *string `json:"payload"`
	UserID     *string `json:"userId"`
	LanguageID *string `json:"languageId"`
}
