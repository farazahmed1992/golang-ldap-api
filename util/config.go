package util

import (
	"os"

	"github.com/spf13/viper"
)

// Config local app.env
type Config struct {
	APITOKEN     string `mapstructure:"API_TOKEN"`
	BindAddress  string `mapstructure:"BIND_ADDRESS"`
	LdapServer   string `mapstructure:"LDAP_SERVER"`
	LdapServer2  string `mapstructure:"LDAP_SERVER2"`
	LdapUser     string `mapstructure:"LDAP_USER"`
	LdapPassword string `mapstructure:"LDAP_PASSWORD"`
	DomainFirst  string `mapstructure:"Domain_First"`
	DomainLast   string `mapstructure:"Domain_Last"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			config.APITOKEN = os.Getenv("API_TOKEN")
			config.BindAddress = os.Getenv("BIND_ADDRESS")
			config.LdapServer = os.Getenv("LDAP_SERVER")
			config.LdapServer2 = os.Getenv("LDAP_SERVER2")
			config.LdapUser = os.Getenv("LDAP_USER")
			config.LdapPassword = os.Getenv("LDAP_PASSWORD")
			config.DomainFirst = os.Getenv("Domain_First")
			config.DomainLast = os.Getenv("Domain_Last")
		}
	}

	err = viper.Unmarshal(&config)

	return

}
