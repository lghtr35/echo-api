package util

import (
	"encoding/json"
	"errors"
	"os"
)

type Configuration struct {
	Version              string   `json:"version"`
	DbConnectionString   string   `json:"dbConnectionString"`
	SwaggerUrl           string   `json:"swaggerUrl"`
	Title                string   `json:"title"`
	Salt                 string   `json:"passwordSalt"`
	AcceptedExtensions   []string `json:"acceptedExtensions"`
	SaveLocations        []string `json:"saveLocations"`
	IsAiAssistantEnabled bool     `json:"isAiAssistantEnabled"`
	secretKey            string
}

func NewConfiguration(logger *Logger) (*Configuration, error) {
	config := new(Configuration)
	//Start filling config with reads
	config, err := ReadConfigFromJSON(config, logger)
	if err != nil {
		logger.Error().Err(err).Msg("Error occurred while reading configuration from local")
	}
	//Other reads may be done in the future
	config = ReadConfigFromEnv(config, logger)

	//Check if config is missing values
	if !isConfigFilled(config) {
		return nil, errors.New("configNotLoadedProperly")
	}

	config.secretKey = config.Title + "_V" + config.Version
	return config, nil
}

func (c *Configuration) GetSecretKey() string {
	return c.secretKey
}

func ReadConfigFromEnv(config *Configuration, logger *Logger) *Configuration {
	c := new(Configuration)
	c.DbConnectionString = os.Getenv("APP_DB_CONN_STR")
	c.Version = os.Getenv("APP_VERSION")
	c.Salt = os.Getenv("APP_PASSWORD_SALT")
	config = copyConfigVals(config, c)
	return config
}

func ReadConfigFromJSON(config *Configuration, logger *Logger) (*Configuration, error) {
	file, err := os.ReadFile("./config.json")
	if err != nil {
		logger.Error().Err(err).Msg("Error reading local JSON config file")
		return config, err
	}

	// Unmarshal the JSON data
	var c Configuration
	err = json.Unmarshal(file, &c)
	if err != nil {
		logger.Error().Err(err).Msg("Error unmarshalling config JSON")
		return config, err
	}
	config = copyConfigVals(config, &c)
	return config, nil
}

func isConfigFilled(c *Configuration) bool {
	return c.DbConnectionString != "" && c.Version != "" && c.SwaggerUrl != "" && c.Title != "" && c.Salt != ""
}

func copyConfigVals(c1 *Configuration, c2 *Configuration) *Configuration {
	if c2.DbConnectionString != "" {
		c1.DbConnectionString = c2.DbConnectionString
	}
	if c2.Version != "" {
		c1.Version = c2.Version
	}
	if c2.SwaggerUrl != "" {
		c1.SwaggerUrl = c2.SwaggerUrl
	}
	if c2.Title != "" {
		c1.Title = c2.Title
	}
	if c2.Salt != "" {
		c1.Salt = c2.Salt
	}

	return c1
}
