package cloudtoolkit

import (
        "github.com/opentracing/opentracing-go"
        "net/http"
        "github.com/spf13/viper"
        "github.com/opentracing/opentracing-go/ext"
         zipkin "github.com/openzipkin/zipkin-go-opentracing"
        "fmt"
)

var Tracer opentracing.Tracer

func InitTracing() {
        var zipkinHost = "192.168.99.100"
        if viper.GetString("profile") != "dev" {
                zipkinHost = "zipkin"
        }
        collector, err := zipkin.NewHTTPCollector(
                fmt.Sprintf("http://%s:9411/api/v1/spans", zipkinHost))
        if err != nil {
                Log.Errorln("Error connecting to zipkin server at " +
                        fmt.Sprintf("http://%s:9411/api/v1/spans", zipkinHost) + ". Error: " + err.Error())
        }
        Tracer, err = zipkin.NewTracer(
                zipkin.NewRecorder(collector, false, "127.0.0.1:0", "compservice"))
        if err != nil {
                Log.Errorln("Error starting new zipkin tracer. Error: " + err.Error())
        }

}

func StartTracing(r *http.Request, opName string) opentracing.Span {
        carrier := opentracing.HTTPHeadersCarrier(r.Header)
        clientContext, err := Tracer.Extract(opentracing.HTTPHeaders, carrier)
        var span opentracing.Span
        if err == nil {
                Log.Println("Compservice could not find an existing Span to attach to in header from caller, creating new span...")
                span = Tracer.StartSpan(
                        opName, ext.RPCServerOption(clientContext))
        } else {
                span = Tracer.StartSpan(opName)
        }
        return span
}
