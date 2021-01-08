/* **********************************
 * Date: 2021-01-07
 * *********************************/

package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/Sciebo-RDS/port-reva/pkg/reva"
)

// WebServer is used to handle all HTTP requests.
type WebServer struct {
	log *zerolog.Logger

	revaClient *reva.Client
}

type endpointHandler = func(*RequestData, http.ResponseWriter, *http.Request) ([]byte, error)
type endpointHandlers = map[string]endpointHandler

func (svr *WebServer) initialize(port uint16, revaClient *reva.Client, log *zerolog.Logger) error {
	if log == nil {
		return errors.Errorf("no logger specified")
	}
	svr.log = log

	if revaClient == nil {
		return errors.Errorf("no Reva client specified")
	}
	svr.revaClient = revaClient

	// Set up and start the HTTP server
	http.HandleFunc("/file", func(w http.ResponseWriter, r *http.Request) {
		svr.handleEndpoint(endpointHandlers{"GET": svr.handleFileGetRequest}, w, r)
	})
	http.HandleFunc("/folder", func(w http.ResponseWriter, r *http.Request) {
		svr.handleEndpoint(endpointHandlers{"GET": svr.handleFolderGetRequest}, w, r)
	})
	go http.ListenAndServe(fmt.Sprintf(":%v", port), nil)

	return nil
}

func (svr *WebServer) handleEndpoint(handlers endpointHandlers, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var respData []byte
	var err error

	if handler, ok := handlers[r.Method]; ok {
		if reqData, errUnmarshal := UnmarshalRequestData(r); errUnmarshal == nil {
			respData, err = handler(&reqData, w, r)
		} else {
			err = errUnmarshal
		}
	} else {
		err = errors.Errorf("unsupported method")
	}

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		errMsg := fmt.Sprintf("%v", err)
		respData = []byte(errMsg)

		svr.log.Warn().Str("method", r.Method).Str("path", r.URL.Path).Msg(errMsg)
	}

	if len(respData) > 0 {
		_, _ = w.Write(respData)
	}
}

func (svr *WebServer) handleFileGetRequest(reqData *RequestData, w http.ResponseWriter, r *http.Request) ([]byte, error) {
	svr.logRequest("file contents request", reqData, r.RemoteAddr)

	fileContents, err := svr.revaClient.DownloadFile(reqData.Metadata.FilePath)
	if err != nil {
		return nil, errors.Wrap(err, "error while retrieving file contents")
	}
	svr.log.Debug().Str("path", reqData.Metadata.FilePath).Int("size", len(fileContents)).Msg("retrieved file contents")

	return fileContents, nil
}

func (svr *WebServer) handleFolderGetRequest(reqData *RequestData, w http.ResponseWriter, r *http.Request) ([]byte, error) {
	svr.logRequest("folder contents request", reqData, r.RemoteAddr)

	folderContents, err := svr.revaClient.ListFolder(reqData.Metadata.FilePath)
	if err != nil {
		return nil, errors.Wrap(err, "error while retrieving folder contents")
	}
	svr.log.Debug().Str("path", reqData.Metadata.FilePath).Int("count", len(folderContents)).Msg("retrieved folder contents")

	reply := map[string][]string{"files": folderContents}
	jsonData, _ := json.Marshal(reply)
	return jsonData, nil
}

func (svr *WebServer) logRequest(msg string, reqData *RequestData, requester string) {
	svr.log.Info().
		Str("path", reqData.Metadata.FilePath).
		Str("userId", reqData.Metadata.UserID).
		Str("apiKey", reqData.Metadata.APIKey).
		Str("requester", requester).
		Msg(msg)
}

// New creates a new WebServer instance.
func New(port uint16, revaClient *reva.Client, log *zerolog.Logger) (*WebServer, error) {
	svr := &WebServer{}
	if err := svr.initialize(port, revaClient, log); err != nil {
		return nil, errors.Wrap(err, "unable to initialize the web server")
	}
	return svr, nil
}
