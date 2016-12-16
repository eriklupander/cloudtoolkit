package cloudtoolkit

import (
	"encoding/json"
	"flag"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
)

const SPRING_CLOUD_CONFIG_SERVER_URL string = "spring.cloud.config.url"

func ResolveProfile() string {
	envProfilePtr := flag.String("profile", "dev", "environment profile, somewhat like spring.profiles.active in java")
	flag.Parse()
	return *envProfilePtr
}

// Tries to resolve the URL to the spring cloud config server using a standardized value injected into Viper.
func LoadSpringCloudConfigFromConfigProperty(appName string, profile string) {
	if viper.IsSet(SPRING_CLOUD_CONFIG_SERVER_URL) {
		url := viper.GetString(SPRING_CLOUD_CONFIG_SERVER_URL)
		LoadSpringCloudConfig(appName, profile, url)
	} else {
		panic("Config property " + SPRING_CLOUD_CONFIG_SERVER_URL + " not set, panicing...")
	}
}

func LoadSpringCloudConfig(appName string, profile string, configServerUrl string) {

	// Try to load from spring cloud config...
	Log.Info("Loading configuration from Spring Cloud Config server at: " + configServerUrl)

	resp, err := http.Get(configServerUrl + "/" + appName + "-" + profile + "/" + profile)
	if err != nil {
		panic("Failed to load configuration: " + err.Error())
	}

	body, err := ioutil.ReadAll(resp.Body)

	var cloudConfig SpringCloudConfig
	json.Unmarshal(body, &cloudConfig)

	for key, value := range cloudConfig.PropertySources[0].Source {
		viper.Set(key, value)
	}
	viper.SetConfigType("json")
}
