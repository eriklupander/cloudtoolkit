# Go Cloud Toolkit
#### WORK IN PROGRESS!!!!!
Strongly opinionated little toolkit for integrating a Go-based Microservice in a "typical" Spring Cloud / Netflix OSS landscape

This "cloud toolkit" is just a bit of glue around other 3rd party libraries for convenience for my own projects, but might be useful for others too. For more complete Go microservices toolkits, see go-kit.

### Supported integrations

* Hystrix through go-hystrix
* AMQP client through amqp (TODO link)
* Zipkin / Opentracing through go-opentracing (TODO link)
* Logging through Logrus (TODO link)
* Eureka (through my own little Eureka library) TODO, see https://github.com/eriklupander/eeureka
* OAuth token verification

### Usage

#### Setting up a config repository
See https://github.com/eriklupander/go-microservice-config for some sample configuration files.

#### Example
Loads configuration from Spring Cloud Config server deployed in your Swarm, inits zipkin, messaging and hystrix.
The sample below uses "from config property", e.g. the LoadSpringCloudConfig needs to load the following properties into viper (with sample value):

- amqp.connection.string: amqp://guest:guest@rabbitmq:5672
- zipkin.service.url: http://zipkin:9411


    import "github.com/eriklupander/cloudtoolkit"

    func main() {
        cloudtoolkit.LoadSpringCloudConfig("My application", "test", "http://configserver:8888")
        cloudtoolkit.InitTracingFromConfigProperty(appName)
        
        amqpClient = cloudtoolkit.InitMessagingClientFromConfigProperty()
        defer amqpClient.GetConn().Close()
        
        cloudtoolkit.ConfigureHystrix([]string{"get_account_secured"}, amqpClient)
    }
    
#### 

### 3rd party libraries
TODO

### LICENSE

MIT License.

See LICENSE.md