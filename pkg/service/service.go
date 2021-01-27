/* **********************************
 * Date: 2021-01-11
 * *********************************/

package service

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

// RegisterService registers this service with the provided token storage given by endpoint.
func RegisterService(endpoint string) error {
	// Send service registration request via POST
	jsonData := `{
		"servicename": "reva",
		"implements": ["fileStorage"],
		"fileTransferMode": 0,
		"fileTransferArchive": 0,
		"credentials": {
			"userId": true,
			"password": true
		}
	}`
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer([]byte(jsonData)))
	if err != nil {
		return errors.Wrap(err, "unable to create HTTP POST request")
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "unable to send HTTP POST request")
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	bodyStr := string(body)

	// Check registration response
	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("unable to register service with token storage (status=%v): %v", resp.StatusCode, bodyStr)
	}

	objects := make(map[string]interface{})
	if err := json.Unmarshal(body, &objects); err != nil {
		return errors.Wrapf(err, "invalid JSON response: %v", bodyStr)
	}

	if _, ok := objects["success"]; !ok {
		return errors.Errorf("unable to register service with token storage: %v", bodyStr)
	}

	return nil
}
