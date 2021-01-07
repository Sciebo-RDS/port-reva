/* **********************************
 * Date: 2021-01-07
 * *********************************/

package server

import (
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
	http.HandleFunc("/file", svr.handleFileEndpoint)
	http.HandleFunc("/folder", svr.handleFolderEndpoint)
	go http.ListenAndServe(fmt.Sprintf(":%v", port), nil)

	return nil
}

func (svr *WebServer) handleFileEndpoint(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		svr.handleFileGetRequest(w, r)

	default:
		svr.log.Warn().Str("method", r.Method).Str("path", r.URL.Path).Msg("unsupported method")
	}
}

func (svr *WebServer) handleFolderEndpoint(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		svr.handleFolderGetRequest(w, r)

	default:
		svr.log.Warn().Str("method", r.Method).Str("path", r.URL.Path).Msg("unsupported method")
	}
}

func (svr *WebServer) handleFileGetRequest(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("file"))
}

func (svr *WebServer) handleFolderGetRequest(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("folder"))
}

// New creates a new WebServer instance.
func New(port uint16, revaClient *reva.Client, log *zerolog.Logger) (*WebServer, error) {
	svr := &WebServer{}
	if err := svr.initialize(port, revaClient, log); err != nil {
		return nil, errors.Wrap(err, "unable to initialize the web server")
	}
	return svr, nil
}
