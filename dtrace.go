package cloudtoolkit

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
	"github.com/spf13/viper"
	"net/http"
)

var ZIPKIN_SERVICE_URL string = "zipkin.service.url"

var Tracer opentracing.Tracer

var zipkinUrl string

// Typical URL: http://192.168.99.100:9411
func InitTracingFromConfigProperty(serviceName string) {
	if viper.IsSet(ZIPKIN_SERVICE_URL) {
		zipkinUrl = viper.GetString(ZIPKIN_SERVICE_URL)
		initTracing(serviceName)
	} else {
		panic("Config property " + ZIPKIN_SERVICE_URL + " not set, panicing...")
	}
}

// Typical URL: http://192.168.99.100:9411
func InitTracingUsingUrl(serviceName string, zipkinHost string) {
	zipkinUrl = zipkinHost
	initTracing(serviceName)
}

func initTracing(serviceName string) {

	collector, err := zipkin.NewHTTPCollector(
		fmt.Sprintf("%s/api/v1/spans", zipkinUrl))
	if err != nil {
		Log.Errorln("Error connecting to zipkin server at " +
			fmt.Sprintf("%s/api/v1/spans", zipkinUrl) + ". Error: " + err.Error())
	}
	Tracer, err = zipkin.NewTracer(
		zipkin.NewRecorder(collector, false, "127.0.0.1:0", serviceName))
	if err != nil {
		Log.Errorln("Error starting new zipkin tracer. Error: " + err.Error())
	}

}

// Loads tracing information from an INCOMING HTTP request.
func StartHTTPTrace(r *http.Request, opName string) opentracing.Span {
	carrier := opentracing.HTTPHeadersCarrier(r.Header)
	clientContext, err := Tracer.Extract(opentracing.HTTPHeaders, carrier)
	var span opentracing.Span
	if err == nil {
		span = Tracer.StartSpan(
			opName, ext.RPCServerOption(clientContext))
	} else {
		span = Tracer.StartSpan(opName)
	}
	return span
}

// Adds tracing information to an OUTGOING HTTP request
func AddTracingToReq(req *http.Request, span opentracing.Span) {
	carrier := opentracing.HTTPHeadersCarrier(req.Header)
	err := Tracer.Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		carrier)
	if err != nil {
		panic("Unable to inject tracing context: " + err.Error())
	}
}
