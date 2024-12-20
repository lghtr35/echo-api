package language

type CreateLanguageRequest struct {
	Name       string `json:"name"`
	Alpha2Code string `json:"alpha2Code"`
	Alpha3Code string `json:"alpha3Code"`
	Icon       string `json:"icon"`
}
