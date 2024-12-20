package prompt

type CreatePromptRequest struct {
	Value     any    `json:"value" form:"value"`
	ContextID string `json:"contextId" form:"contextId"`
	EntityID  string `json:"entityId" form:"entityId"`
}
