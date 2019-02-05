package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/ONSdigital/go-ns/log"
)

// Request ...
type Request struct {
	Authorization *Authorization
	Client        *http.Client
}

// Authorization ...
type Authorization struct {
	Username string
	Password string
}

// Institution ...
type Institution struct {
	Name string `json:"Name"`
}

// GetInstitutionName ...
func (request *Request) GetInstitutionName(path string) (string, error) {
	logData := log.Data{"path": path}
	URL, err := url.Parse(path)
	if err != nil {
		log.ErrorC("failed to create url for api call", err, logData)
		return "", err
	}
	path = URL.String()
	logData["path"] = path

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		log.ErrorC("failed to create request for unistats api", err, logData)
		return "", err
	}

	req.SetBasicAuth(request.Authorization.Username, request.Authorization.Password)
	req.Header.Set("Content-Type", "application/json")

	resp, err := request.Client.Do(req)
	if err != nil {
		log.ErrorC("Failed to action unistats api", err, logData)
		return "", err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.ErrorC("failed to read body from unistats api response", err, logData)
		return "", err
	}

	var institution Institution
	if err = json.Unmarshal(b, &institution); err != nil {
		log.ErrorC("unable to unmarshal bytes into institution resource", err, logData)
		return "", err
	}

	return institution.Name, nil
}
