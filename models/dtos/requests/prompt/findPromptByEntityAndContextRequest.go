package prompt

type FindPromptByEntityAndContextRequest struct {
	EntityID  string `json:"entityId" form:"entityId"`
	ContextID string `json:"contextId" form:"contextId"`
}
