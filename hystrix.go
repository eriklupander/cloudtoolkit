package cloudtoolkit

import (
	"encoding/json"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/spf13/viper"
	"net"
	"net/http"
	"time"
)

func ConfigureHystrix(commands []string, amqpClient *MessagingClient) {

	for _, command := range commands {
		hystrix.ConfigureCommand(command, hystrix.CommandConfig{
			Timeout:                resolveProperty(command, "Timeout"),
			MaxConcurrentRequests:  resolveProperty(command, "MaxConcurrentRequests"),
			ErrorPercentThreshold:  resolveProperty(command, "ErrorPercentThreshold"),
			RequestVolumeThreshold: resolveProperty(command, "RequestVolumeThreshold"),
			SleepWindow:            resolveProperty(command, "SleepWindow"),
		})
	}

	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()
	go http.ListenAndServe(net.JoinHostPort("", "8181"), hystrixStreamHandler)
	Log.Println("Launched hystrixStreamHandler at 8181")

	// Publish presence on RabbitMQ
	publishDiscoveryToken(amqpClient)
}

func publishDiscoveryToken(amqpClient *MessagingClient) {
	token := DiscoveryToken{
		State:   "UP",
		Address: GetLocalIP(),
	}
	json, _ := json.Marshal(token)
	go func() {
		for {
			amqpClient.SendMessage(string(json), "application/json", "discovery")
			time.Sleep(time.Second * 30)
		}
	}()
}

func resolveProperty(command string, prop string) int {
	if viper.IsSet("hystrix.command." + command + "." + prop) {
		return viper.GetInt("hystrix.command." + command + "." + prop)
	} else {
		return getDefaultHystrixConfigPropertyValue(prop)
	}
}
func getDefaultHystrixConfigPropertyValue(prop string) int {
	switch prop {
	case "Timeout":
		return hystrix.DefaultTimeout
	case "MaxConcurrentRequests":
		return hystrix.DefaultMaxConcurrent
	case "RequestVolumeThreshold":
		return hystrix.DefaultVolumeThreshold
	case "SleepWindow":
		return hystrix.DefaultSleepWindow
	case "ErrorPercentThreshold":
		return hystrix.DefaultErrorPercentThreshold
	}
	panic("Got unknown hystrix property: " + prop + ". Panicing!")
}
