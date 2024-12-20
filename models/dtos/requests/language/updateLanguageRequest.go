package language

type UpdateLanguageRequest struct {
	ID         string  `json:"id"`
	Name       *string `json:"name"`
	Alpha2Code *string `json:"alpha2code"`
	Alpha3Code *string `json:"alpha3code"`
}
