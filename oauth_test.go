package cloudtoolkit

import (
	"net/http"
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

var sampleToken = "my-oauth-token"

func TestOAuthSpec(t *testing.T) {

        Convey("Given a HTTP request with a valid sample token", t, func() {
                httpReq := http.Request{
                        Header: map[string][]string{
                                "Authorization": {"Bearer: " + sampleToken},
                        },
                }

                Convey("When the auth header is extracted", func() {
                        token := extractAuthorizationFromHeader(&httpReq)

                        Convey("The auth header extracted should match the sample one", func() {
                                So(token, ShouldNotBeNil)
                                So(token, ShouldEqual, sampleToken)
                        })
                })
        })

        Convey("Given a HTTP request missing a token", t, func() {
                httpReq := http.Request{
                        Header: map[string][]string{
                                "Authorization": {""},
                        },
                }

                Convey("When the auth header is extracted", func() {
                        token := extractAuthorizationFromHeader(&httpReq)

                        Convey("The auth header extracted should be empty", func() {
                                So(token, ShouldNotBeNil)
                                So(token, ShouldBeEmpty)
                        })
                })
        })
}