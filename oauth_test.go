package cloudtoolkit

import (
        "testing"
        "net/http"
)

var sampleToken = "my-cool-token"

func TestExtract(t *testing.T) {

        httpReq := http.Request{
                Header: map[string][]string{
                        "Authorization": {"Bearer: " + sampleToken},
                },
        }
        token := extractAuthorizationFromHeader(&httpReq)
        assertEquals(token, sampleToken, t)

}

func assertEquals(s1 string, s2 string, t *testing.T) {
        if s1 != s2 {
                t.Error("Expected " + s1 + ", got " + s2)
        }
}
