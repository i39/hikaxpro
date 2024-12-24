package hikaxprogo

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetSessionParams(t *testing.T) {
	hik := &HikISAPI{
		host:     "http://localhost",
		port:     "8080",
		username: "admin",
		password: "password",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<sessionCapabilities xmlns="http://www.isapi.org/ver20/XMLSchema">
			<sessionID>12345</sessionID>
			<challenge>challenge</challenge>
			<salt>salt</salt>
			<salt2>salt2</salt2>
			<isIrreversible>true</isIrreversible>
			<iterations>1000</iterations>
		</sessionCapabilities>`))
	}))
	defer server.Close()

	hik.host = server.URL

	capabilities, err := hik.getSessionParams()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if capabilities.SessionID != "12345" {
		t.Errorf("Expected sessionID to be '12345', got %s", capabilities.SessionID)
	}
}

func TestEncodePassword(t *testing.T) {
	hik := &HikISAPI{
		username: "admin",
		password: "password",
	}

	cap := sessionCapabilities{
		Challenge:      "challenge",
		Salt:           "salt",
		Salt2:          "salt2",
		IsIrreversible: "true",
		Iterations:     1000,
	}

	encodedPassword := hik.encodePassword(cap)
	if encodedPassword == "" {
		t.Error("Expected encoded password, got empty string")
	}
}
