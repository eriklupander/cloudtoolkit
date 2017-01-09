package cloudtoolkit

import (
	"net/http"
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

var sampleToken = "my-cool-token"

func TestOAuthSpec(t *testing.T) {

        Convey("Given a HTTP request with a sample token", t, func() {
                httpReq := http.Request{
                        Header: map[string][]string{
                                "Authorization": {"Bearer: " + sampleToken},
                        },
                }
                Convey("When the auth header is extracted", func() {
                        token := extractAuthorizationFromHeader(&httpReq)

                        Convey("The auth header extracted should match the sample one", func() {
                                So(token, ShouldEqual, sampleToken)
                                So(token, ShouldEqual, sampleToken)
                        })
                })
        })
}