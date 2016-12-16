package cloudtoolkit

import (
        "flag"
        "net/http"
        "io/ioutil"
        "encoding/json"
        "github.com/spf13/viper"
)

func loadConfig(appName string, configServerUrl string) {

        // Get environment from flag, otherwise default to dev.
        envProfilePtr := flag.String("profile", "dev", "environment profile, somewhat like spring.profiles.active in java")
        flag.Parse()
        envProfile := *envProfilePtr
        Log.Println("Running as profile: " + envProfile)
        if envProfile != "dev" {
                configServerUrl = "http://configserver:8888"
        }

        // Try to load from spring cloud config...
        Log.Println("Loading configuration from Spring Cloud Config server at: " + configServerUrl)

        resp, err := http.Get(configServerUrl + "/" + appName + "-" + envProfile + "/" + envProfile)
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
