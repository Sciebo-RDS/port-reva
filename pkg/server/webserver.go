/* **********************************
 * Date: 2021-01-07
 * *********************************/

package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"

	"github.com/Sciebo-RDS/port-reva/pkg/reva"
)

// WebServer is used to handle all HTTP requests.
type WebServer struct {
	log *zerolog.Logger

	revaConfig reva.Config
}

type endpointHandler = func(*RequestData, *reva.Client, http.ResponseWriter, *http.Request) ([]byte, error)
type endpointHandlers = map[string]endpointHandler

func (svr *WebServer) initialize(port uint16, revaConfig reva.Config, log *zerolog.Logger) error {
	if log == nil {
		return errors.Errorf("no logger specified")
	}
	svr.log = log

	svr.revaConfig = revaConfig

	// Set up and start the HTTP server
	http.HandleFunc("/storage/file", func(w http.ResponseWriter, r *http.Request) {
		svr.handleEndpoint(endpointHandlers{"GET": svr.handleFileGetRequest}, w, r)
	})
	http.HandleFunc("/storage/folder", func(w http.ResponseWriter, r *http.Request) {
		svr.handleEndpoint(endpointHandlers{"GET": svr.handleFolderGetRequest}, w, r)
	})

	// Also serve Prometheus metrics for health checking etc.
	http.Handle("/metrics", promhttp.Handler())

	go http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
	return nil
}

func (svr *WebServer) handleEndpoint(handlers endpointHandlers, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var respData []byte
	var err error

	if handler, ok := handlers[r.Method]; ok {
		if reqData, errUnmarshal := UnmarshalRequestData(r); errUnmarshal == nil {
			// Found a handler, so create a Reva client used to handle the request
			if client, errClient := svr.createRevaClient(); err == nil {
				respData, err = handler(&reqData, client, w, r)
			} else {
				err = errClient
			}
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

func (svr *WebServer) handleFileGetRequest(reqData *RequestData, revaClient *reva.Client, w http.ResponseWriter, r *http.Request) ([]byte, error) {
	svr.logRequest("file contents request", reqData, r.RemoteAddr)

	fileContents, err := revaClient.DownloadFile(reqData.Metadata.FilePath)
	if err != nil {
		return nil, errors.Wrap(err, "error while retrieving file contents")
	}
	svr.log.Debug().Str("path", reqData.Metadata.FilePath).Int("size", len(fileContents)).Msg("retrieved file contents")

	return fileContents, nil
}

func (svr *WebServer) handleFolderGetRequest(reqData *RequestData, revaClient *reva.Client, w http.ResponseWriter, r *http.Request) ([]byte, error) {
	svr.logRequest("folder contents request", reqData, r.RemoteAddr)

	folderContents, err := revaClient.ListFolder(reqData.Metadata.FilePath)
	if err != nil {
		return nil, errors.Wrap(err, "error while retrieving folder contents")
	}
	svr.log.Debug().Str("path", reqData.Metadata.FilePath).Int("count", len(folderContents)).Msg("retrieved folder contents")

	reply := map[string][]string{"files": folderContents}
	jsonData, _ := json.Marshal(reply)
	return jsonData, nil
}

func (svr *WebServer) createRevaClient() (*reva.Client, error) {
	client, err := reva.New(svr.revaConfig.Host, svr.revaConfig.User, svr.revaConfig.Password, svr.log)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create the Reva client")
	}
	svr.log.Debug().Str("host", svr.revaConfig.Host).Str("user", svr.revaConfig.User).Msg("established Reva session")
	return client, nil

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
func New(port uint16, revaConfig reva.Config, log *zerolog.Logger) (*WebServer, error) {
	svr := &WebServer{}
	if err := svr.initialize(port, revaConfig, log); err != nil {
		return nil, errors.Wrap(err, "unable to initialize the web server")
	}
	return svr, nil
}
