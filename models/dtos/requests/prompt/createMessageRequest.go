package prompt

type CreateMessageRequest struct {
	Value     string `json:"value" form:"value"`
	ContextID string `json:"contextId" form:"contextId"`
}
