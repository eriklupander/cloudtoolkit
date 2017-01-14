package cloudtoolkit

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"testing"
	"time"
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

func TestSessionCache(t *testing.T) {
	Convey("Given that the sessionCache is empty", t, func() {
		sessionCache := SessionCache{}
		instant := time.Now()
		token := "token-123"

		Convey("When a new Token is inserted", func() {
			sessionCache.Put(token, instant)

			Convey("Assert that the token has expiry + 1 hour", func() {
				So(sessionCache.IsValid(token), ShouldBeTrue)
				So(sessionCache.Get(token), ShouldResemble, instant.Add(time.Hour*1))
			})

			Convey("Assert that another token is not valid", func() {
				So(sessionCache.IsValid("other-token"), ShouldBeFalse)
			})
		})

		Convey("When a faked expired Token is inserted", func() {
			sessionCache.Put("expired-token", instant.Add(time.Hour*-3))

			Convey("Assert that this is NOT a valid token", func() {
				So(sessionCache.IsValid("expired-token"), ShouldBeFalse)
			})
		})
	})
}
