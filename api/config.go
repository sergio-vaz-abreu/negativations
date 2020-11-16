package api

type ApplicationConfig struct {
	SymmetricKeyConfig SymmetricKeyConfig `json:"symmetricKeyConfig"`
	ApiConfig          ApiConfig          `json:"apiConfig"`
	ArangoConfig       ArangoConfig       `json:"arangoConfig"`
	LegacyConfig       LegacyConfig       `json:"legacyConfig"`
}

type SymmetricKeyConfig struct {
	SymmetricKey      string `json:"symmetricKey"`
	EncryptionContext string `json:"encryptionContext"`
}

type ApiConfig struct {
	Port int `json:"port"`
}

type ArangoConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type LegacyConfig struct {
	Url string `json:"url"`
}
