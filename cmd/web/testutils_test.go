package main

import (
	"bytes"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"html"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"regexp"
	"snippetbox.alexedwards.net/internal/models/mocks"
	"testing"
	"time"
)

// Define a regular expression to capture the CSRF token value
var csrfTokenRX = regexp.MustCompile(`<input type='hidden' name='csrf_token' value='(.+)'>`)

func extractCSRFToken(t *testing.T, body string) string {
	// Try to find the CSRF token from the HTML body
	matches := csrfTokenRX.FindStringSubmatch(body)
	if len(matches) < 2 {
		t.Fatal("no csrf token found in body")
	}
	// Return the unescaped token string
	return html.UnescapeString(string(matches[1]))
}

func newTestApplication(t *testing.T) *application {
	templateCache, err := newTemplateCache()
	if err != nil {
		t.Fatal(err)
	}

	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	return &application{
		errorLog:       log.New(io.Discard, "", 0),
		infoLog:        log.New(io.Discard, "", 0),
		snippet:        &mocks.SnippetModel{},
		user:           &mocks.UserModel{},
		templateCache:  templateCache,
		sessionManager: sessionManager,
		formDecoder:    formDecoder,
	}

}

// testServer struct wraps our httptest.Server to make our code cleaner
type testServer struct {
	*httptest.Server
}

func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)

	// create new cookiejar to store cookies
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	// attach cookiejar to test client
	ts.Client().Jar = jar

	// stop the client from following redirect
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{ts}
}

func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, string) {
	res, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	bytes.TrimSpace(body)

	return res.StatusCode, res.Header, string(body)
}

// This method sends a POST request to the test server.
// The last parameter is a `url.Values` object that holds the form data.
func (ts *testServer) postForm(t *testing.T, urlPath string, form url.Values) (int, http.Header, string) {
	rs, err := ts.Client().PostForm(ts.URL+urlPath, form)
	if err != nil {
		t.Fatal(err)
	}

	// Read the response body from the test server.
	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	bytes.TrimSpace(body)

	// Return status code, headers, and body.
	return rs.StatusCode, rs.Header, string(body)
}
