package cloudtoolkit

import (
        "crypto/tls"
        "github.com/spf13/viper"
        "net/http"
        "strings"
        "time"
)

const AUTH_SERVER_USER_URL = "auth.server.user.url"

var authServerUserUrl string

var sessionCache SessionCache

func InitOAuth2Handler() {
        authServerUserUrl = viper.GetString(AUTH_SERVER_USER_URL)
}

func InitOAuth2HandlerUsingUrl(url string) {
        authServerUserUrl = url
}

func OAuth2Handler(inner http.Handler) http.Handler {
        // Init the cache
        sessionCache = SessionCache{}

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
        if len(s) != 2 {
                return ""
        }

        return s[1]
}

func checkAuth(w http.ResponseWriter, r *http.Request) bool {
        if authServerUserUrl == "" {
                panic("authServerUserUrl is not specified, should be something like https://192.168.99.100:9999/uaa/user")
        }
        // try to find authorization header
        token := extractAuthorizationFromHeader(r)

        if token == "" {
                return false
        }

        // Check session cache
        if sessionCache.IsValid(token) {
                return true
        }

        req, _ := http.NewRequest("GET", authServerUserUrl + "?access_token=" + token, nil) //"https://192.168.99.100:9999/uaa/user?access_token="+token, nil)
        var DefaultTransport http.RoundTripper = &http.Transport{
                TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        }
        resp, err := DefaultTransport.RoundTrip(req)

        if err != nil {
                panic("Could not contact OAuth server: " + err.Error())
        }
        if resp.StatusCode == 200 {
                sessionCache.Put(token, time.Now())
                return true
        } else {
                return false
        }
}

type SessionCache struct {
        Store map[string]time.Time
}

func (s *SessionCache) Put(token string, instant time.Time) {
        if s.Store == nil {
                s.Store = make(map[string]time.Time)
        }
        s.Store[token] = instant.Add(time.Hour * 1)
}

func (s *SessionCache) Get(token string) time.Time {
        return s.Store[token]
}

func (s *SessionCache) IsValid(token string) bool {
        if !s.Store[token].IsZero() {
                if time.Now().After(s.Store[token]) {
                        // Expired
                        delete(s.Store, token)
                        return false
                }
                return true
        }
        return false
}