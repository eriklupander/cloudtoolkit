package cloudtoolkit

import (
        "net/http"
        "strings"
        "crypto/tls"
)

func OAuth2Handler(inner http.Handler, name string) http.Handler {

        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                if checkAuth(w, r) {
                        inner.ServeHTTP(w, r)
                        return
                }

                w.WriteHeader(401)
                w.Write([]byte("401 Unauthorized\n"))
        })
}

func extractAuthorizationFromHeader(r *http.Request) string {
        s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
        if len(s) != 2 { return "" }

        return s[1]
}

func checkAuth(w http.ResponseWriter, r *http.Request) bool {
        // try to find authorization header
        token := extractAuthorizationFromHeader(r)
        req, _ := http.NewRequest("GET", "https://192.168.99.100:9999/uaa/user?access_token=" + token, nil)
        var DefaultTransport http.RoundTripper = &http.Transport{
                TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        }
        resp, err := DefaultTransport.RoundTrip(req)

        if err != nil {
                panic("Could not contact OAuth server: " + err.Error())
        }
        if resp.StatusCode == 200 {
               return true
        } else {
                return false
        }
}

