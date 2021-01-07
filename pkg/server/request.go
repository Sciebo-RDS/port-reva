/* **********************************
 * Date: 2021-01-07
 * *********************************/

package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

// RequestMetadata holds all basic metadata of a request.
type RequestMetadata struct {
	FilePath string
	UserID   string
	APIKey   string
}

// RequestData holds the data of a request.
type RequestData struct {
	Metadata RequestMetadata
	Data     map[string]interface{}
}

// UnmarshalRequestData extracts the request data from an HTTP request.
func UnmarshalRequestData(r *http.Request) (RequestData, error) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return RequestData{}, errors.Wrap(err, "unable to read the request body")
	}

	objects := make(map[string]interface{})
	if err := json.Unmarshal(data, &objects); err != nil {
		return RequestData{}, errors.Wrap(err, "unable to unmarshal the JSON data")
	}

	// Some hardcoded values are stored into the request metadata;
	// any other values are stored as-is in the request data map.
	reqdata := RequestData{Data: make(map[string]interface{})}
	for key, val := range objects {
		switch key {
		case "filepath":
			reqdata.Metadata.FilePath = val.(string)

		case "userId":
			reqdata.Metadata.UserID = val.(string)

		case "apiKey":
			reqdata.Metadata.APIKey = val.(string)

		default:
			reqdata.Data[key] = val
		}
	}

	return reqdata, nil
}
