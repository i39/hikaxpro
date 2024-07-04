package hikaxprogo

import (
	"crypto/sha256"
	"net/http"
	"time"

	"bytes"
	"encoding/hex"
	json "encoding/json"
	xml "encoding/xml"
	"errors"
	"io"
	"strconv"
)

type HikISAPI struct {
	host     string
	port     string
	username string
	password string
	session  http.Cookie // Session cookie
}
type sessionCapabilities struct {
	XMLNS          string `xml:"xmlns,attr"`
	SessionID      string `xml:"sessionID"`
	Challenge      string `xml:"challenge"`
	Salt           string `xml:"salt"`
	Salt2          string `xml:"salt2"`
	IsIrreversible string `xml:"isIrreversible"`
	Iterations     int    `xml:"iterations"`
}

type SessionLogin struct {
	SessionID        string `xml:"sessionID"`
	Password         string `xml:"password"`
	UserName         string `xml:"userName"`
	SessionIDVersion string `xml:"sessionIDVersion"`
}

func (hik *HikISAPI) getSessionParams() (sessionCapabilities, error) {

	capabilities := sessionCapabilities{}

	resp, err := hik.makeRequest("GET", hik.host+":"+hik.port+Session_Capabilities+hik.username, "")
	if err != nil {
		return capabilities, err
	}
	defer resp.Body.Close()
	// parse xmlns from response

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// Unmarshal the XML from the response into the struct
	err = xml.Unmarshal(body, &capabilities)
	if err != nil {
		return capabilities, err
	}

	return capabilities, nil
}

// encodePassword encodes the password using the session capabilities
func (hik *HikISAPI) encodePassword(cap sessionCapabilities) string {
	var result [32]byte

	// encode hiwatch password
	//
	if cap.IsIrreversible == "true" {
		result = sha256.Sum256([]byte(hik.username + cap.Salt + hik.password))
		result = sha256.Sum256([]byte(hik.username + cap.Salt2 + hex.EncodeToString(result[:])))
		result = sha256.Sum256([]byte(hex.EncodeToString(result[:]) + cap.Challenge))

		for i := 2; i < cap.Iterations; i++ {
			result = sha256.Sum256([]byte(hex.EncodeToString(result[:])))
		}

	} else {
		result = sha256.Sum256([]byte(hik.password + cap.Challenge))
		for i := 1; i < cap.Iterations; i++ {
			result = sha256.Sum256([]byte(hex.EncodeToString(result[:])))
		}
	}

	return hex.EncodeToString(result[:])

}

func (hik *HikISAPI) Login() error {

	cap, err := hik.getSessionParams()
	if err != nil {
		return err
	}
	encodedPassword := hik.encodePassword(cap)
	// Build the login request
	loginRequest := SessionLogin{
		UserName:         hik.username,
		SessionID:        cap.SessionID,
		Password:         string(encodedPassword[:]),
		SessionIDVersion: "2.1",
	}
	// Make the login request
	xmlLoginRequest, err := xml.Marshal(loginRequest)
	if err != nil {
		return err
	}
	//get POSIX time
	dt := time.Now().Unix()
	sessionLoginUrl := hik.host + ":" + hik.port + Session_Login + "?timeStamp=" + strconv.FormatInt(dt, 10)

	strLoginRequest := string((xmlLoginRequest[:]))
	resp, err := hik.makeRequest("POST", sessionLoginUrl, strLoginRequest)

	if err != nil {
		return err
	}
	if resp.StatusCode == 200 {
		hik.session = *resp.Cookies()[0]
		return nil
	}

	return errors.New("Login failed")

}

func (hik *HikISAPI) makeRequest(method string, url string, body string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(body)))
	if err != nil {
		return nil, err
	}
	// req.Header.Add("Accept", "application/xml")
	// req.Header.Add("Content-Type", "application/xml")
	req.Header.Add("Cookie", hik.session.String())
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == 401 {
		// Session expired, try to login again
		err := hik.Login()
		if err != nil {
			return nil, err
		}
		return hik.makeRequest(method, url, body)
	}
	return resp, nil
}

func (hik *HikISAPI) ZoneStatus() (ZoneList, error) {

	z := ZoneList{}

	resp, err := hik.makeRequest("GET", hik.host+":"+hik.port+ZoneStatus, "")
	if err != nil {
		return z, err
	}
	defer resp.Body.Close()
	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// Unmarshal the JSON from the response into the struct
	err = json.Unmarshal(body, &z)
	if err != nil {
		return z, err
	}

	return z, nil
}

func (hik *HikISAPI) ExDevData() (ExDevData, error) {
	e := ExDevData{}

	resp, err := hik.makeRequest("GET", hik.host+":"+hik.port+PeripheralsStatus, "")
	if err != nil {
		return e, err
	}

	defer resp.Body.Close()
	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// Unmarshal the JSON from the response into the struct
	err = json.Unmarshal(body, &e)
	if err != nil {
		return e, err
	}

	return e, nil
}

func New(host string, port string, username string, password string) *HikISAPI {
	hik := new(HikISAPI)
	// check if host starts with http:// or https://
	if host[0:4] == "http" {
		hik.host = host
	} else {
		hik.host = "http://" + host
	}
	hik.port = port
	hik.username = username
	hik.password = password
	hik.session = http.Cookie{}
	return hik
}
