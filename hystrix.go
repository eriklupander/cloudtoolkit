package cloudtoolkit

import (
        "github.com/afex/hystrix-go/hystrix"
        "net/http"
        "net"
        "encoding/json"
        "time"
)

func configureHystrix(amqpClient IMessagingClient) {

        hystrix.ConfigureCommand("get_account_image", hystrix.CommandConfig{
                Timeout:               3000,
                MaxConcurrentRequests: 100,
                ErrorPercentThreshold: 25,
        })
        hystrix.ConfigureCommand("get_account", hystrix.CommandConfig{
                Timeout:               3000,
                MaxConcurrentRequests: 100,
                ErrorPercentThreshold: 25,
        })

        hystrixStreamHandler := hystrix.NewStreamHandler()
        hystrixStreamHandler.Start()
        go http.ListenAndServe(net.JoinHostPort("", "8181"), hystrixStreamHandler)
        Log.Println("Launched hystrixStreamHandler at 8181")

        // Publish presence on RabbitMQ
        publishDiscoveryToken(amqpClient)
}

func publishDiscoveryToken(amqpClient IMessagingClient) {
        token := DiscoveryToken{
                State:   "UP",
                Address: GetLocalIP(),
        }
        json, _ :=  json.Marshal(token)
        go func() {
                for {
                        amqpClient.SendMessage(string(json), "application/json", "discovery")
                        time.Sleep(time.Second * 30)
                }
        }()
}