package managers

type PromptGenManager interface {
	GeneratePrompt(any) (string, error)
	GeneratePromptWith(any, PromptAction) (string, error)
	GenerateMessage(string) (string, error)
}

type PromptAction string

const (
	Initial   PromptAction = "Hi, you are going to assist customers with their questions or any request within the context given to you. Rules are these:\n 1. There will be prompts where you will need to do according to the action in them. Syntax is \"Prompt(<action>)\"\n 2. You will answer messages within the context as an assistant when a message sent. Syntax is \"Message(<string>)\"\n 3. Actions might be remember, forget or forgetAll. You will do the action and if it is done successfully respond \"done\", if there is any error on your side please respond with \"failed. <error>\". Syntax is \"<action>(<string optional>)\"\n 4. Remember action is for you to keep a given message in mind for future interactions\n 5. Forget action is for you to forget and dont bring up a given info anymore\n 6. ForgetAll action is for you to forget all the previous Prompts given and start fresh.\nPlease, try to keep answers short and focused and thank you for assisting me and the customers. "
	Prompt    PromptAction = "Prompt(%s)"
	Message   PromptAction = "Message(%s)"
	Remember  PromptAction = "Remember(%s)"
	Forget    PromptAction = "Forget(%s)"
	ForgetAll PromptAction = "ForgetAll()"
)

func (pa PromptAction) String() string {
	return string(pa)
}
