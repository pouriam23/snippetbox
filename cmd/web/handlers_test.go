package main

import (
	"net/http"
	"snippetbox.alexedwards.net/internal/assert"
	"testing"
)

// unit testing
/*func TestPing(t *testing.T) {
	// Create a new ResponseRecorder
	rr := httptest.NewRecorder()


	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err) // Stop the test if the request can't be created
	}


	ping(rr, r)


	rs := rr.Result()


	assert.Equal(t, rs.StatusCode, http.StatusOK)


	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err) // Stop the test if reading body fails
	}

	bytes.TrimSpace(body)

	assert.Equal(t, string(body), "OK")
}*/

// E2E testing
func TestPing(t *testing.T) {
	// create a fake version of our application
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())

	code, _, body := ts.get(t, "/ping")
	assert.Equal(t, code, http.StatusOK)
	assert.Equal(t, string(body), "ok")
}

func TestSnippetView(t *testing.T) {
	// Create app instance with mocks
	app := newTestApplication(t)

	// Start a test server using our routes
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	// Define all the test cases we want to run
	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody string
	}{
		{
			name:     "Valid ID",
			urlPath:  "/snippet/view/1",
			wantCode: http.StatusOK,
			wantBody: "An old silent pond...", // Check for this text
		},
		{
			name:     "Non-existent ID",
			urlPath:  "/snippet/view/2",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Negative ID",
			urlPath:  "/snippet/view/-1",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Decimal ID",
			urlPath:  "/snippet/view/1.23",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "String ID",
			urlPath:  "/snippet/view/foo",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Empty ID",
			urlPath:  "/snippet/view/",
			wantCode: http.StatusNotFound,
		},
	}

	// Run each test case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.get(t, tt.urlPath)

			// Check response status code
			assert.Equal(t, code, tt.wantCode)

			// If we expect content, check it's in the body
			if tt.wantBody != "" {
				assert.StringContains(t, body, tt.wantBody)
			}
		})
	}
}
