package main

import (
	"time"
	"strconv"
	"encoding/json"
	"html/template"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/dgrijalva/jwt-go"
	"github.com/qiniu/qpaint/paintdom"
)

// M represents a general dictionary type.
type M map[string]interface{}

// ---------------------------------------------------------------------------

type errorRet struct {
	Err string `json:"error"`
}

// Reply replies a http response.
func Reply(w http.ResponseWriter, code int, data interface{}) {

	msg, err := json.Marshal(data)
	if err != nil {
		Reply(w, 500, &errorRet{Err: err.Error()})
		return
	}

	h := w.Header()
	h.Set("Content-Length", strconv.Itoa(len(msg)))
	h.Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(msg)
}

// ReplyErr replies an error as a http response.
func ReplyErr(w http.ResponseWriter, code int, err string) {

	Reply(w, code, &errorRet{Err: err})
}

// ----------------------------------------------------------

type sessionToAuth struct {
	base http.Handler
}

func newSessionToAuth(base http.Handler) http.Handler {
	return &sessionToAuth{base: base}
}

func (p *sessionToAuth) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("session")
	if err != nil || isExpired(cookie) {
		ReplyErr(w, 401, "bad token")
		return
	}
	id, ok := parseToken(cookie.Value)
	if !ok {
		ReplyErr(w, 401, "bad token")
		return
	}
	req = cloneRequest(req)
	req.Header.Del("Cookie")
	req.Header.Set("Authorization", "QPaintStub " + id)
	p.base.ServeHTTP(w, req)
}

func isExpired(cookie *http.Cookie) bool {
	return false
}

func cloneRequest(r *http.Request) *http.Request {
	// shallow copy of the struct
	r2 := *r
	// deep copy of the Header
	r2.Header = make(http.Header, len(r.Header))
	for k, s := range r.Header {
		r2.Header[k] = append([]string(nil), s...)
	}
	return &r2
}

// ----------------------------------------------------------

var loginTmpl = template.Must(template.New("login").Parse(`<html>
  <body>
    <form action="/login" method="post">
       <p>
         Username: <input type="text" name="username">
       </p>
       <p>
		 Password: <input type="text" name="password">
		 <input type="hidden" name="return" value="{{ .ReturnURL }}">
       </p>
       <input type="submit" value="Login">
    </form>
  </body>
</html>`))

func handleLogin(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET", "HEAD":
		loginTmpl.Execute(w, M{"ReturnURL": req.FormValue("return")})
	case "POST":
		userName := req.PostFormValue("username")
		password := req.PostFormValue("password")
		returnURL := req.PostFormValue("return")
		token, ok := login(userName, password)
		if !ok {
			ReplyErr(w, 401, "bad token")
			return
		}
		http.SetCookie(w, &http.Cookie{Name: "session", Value: token, MaxAge: 86400})
		http.Redirect(w, req, returnURL, 301)
	}
}

var (
	signKey = []byte("please protect your key")
)

// login 用来检查到数据库 userName+password 是否正确
// 如果正确，返回对应的 accessToken。
func login(userName, password string) (token string, ok bool) {
	// 正常应该到数据库查，但是我们这里假设 userName, password 都是 uid
	if userName != password {
		return
	}
	_, err := strconv.Atoi(userName)
	if err != nil {
		return
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(86400 * time.Second).Unix(),
		Id: userName,
	})
	if token, err = t.SignedString(signKey); err != nil {
		return
	}
	return token, true
}

func parseToken(token string) (id string, ok bool) {
	t, err := jwt.ParseWithClaims(token, new(jwt.StandardClaims), func(token *jwt.Token) (interface{}, error) {
        return signKey, nil
	})
	if err != nil {
		return
	}
	claims := t.Claims
	if claims.Valid() != nil {
		return
	}
	return claims.(*jwt.StandardClaims).Id, true
}

// ----------------------------------------------------------

func newReverseProxy(baseURL string) *httputil.ReverseProxy {
	rpURL, _ := url.Parse(baseURL)
	return httputil.NewSingleHostReverseProxy(rpURL)
}

func handleDefault(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/" {
		http.ServeFile(w, req, "www/index.htm")
		return
	}
	req.URL.RawQuery = "" // skip "?params"
	wwwServer.ServeHTTP(w, req)
}

var (
	apiReverseProxy = newReverseProxy("http://localhost:9999")
	wwwServer       = http.FileServer(http.Dir("www"))
)

func main() {
	go paintdom.Main()
	http.Handle("/api/", http.StripPrefix("/api/", newSessionToAuth(apiReverseProxy)))
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/", handleDefault)
	http.ListenAndServe(":8888", nil)
}

// ----------------------------------------------------------
