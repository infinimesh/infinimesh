package config

type OrgAccess struct {
	Namespace string `yaml:"ns"`
	Level     int32  `yaml:"level"`
}

type Config struct {
	ClientId            string               `yaml:"client_id"`
	ClientSecret        string               `yaml:"client_secret"`
	RedirectUrl         string               `yaml:"redirect_url"`
	Scopes              []string             `yaml:"scopes"`
	State               string               `yaml:"state"`
	ApiUrl              string               `yaml:"api_url"`
	AuthUrl             string               `yaml:"auth_url"`
	TokenUrl            string               `yaml:"token_url"`
	OrganizationMapping map[string]OrgAccess `yaml:"organization_mapping"`
}
