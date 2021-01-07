/* **********************************
 * Date: 2021-01-06
 * *********************************/

package runtime

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/Sciebo-RDS/port-reva/pkg/reva"
	"github.com/Sciebo-RDS/port-reva/pkg/server"
)

// Runtime implements the main program runtime.
type Runtime struct {
	log *zerolog.Logger

	conf Config

	revaClient *reva.Client
	webServer  *server.WebServer
}

const (
	runLoopSleeptime = time.Millisecond * 100
)

func (rt *Runtime) initialize(cfg Config, log *zerolog.Logger) error {
	if log == nil {
		return errors.Errorf("no logger specified")
	}
	rt.log = log

	rt.conf = cfg

	client, err := reva.New(cfg.Reva.Host, cfg.Reva.User, cfg.Reva.Password, log)
	if err != nil {
		return errors.Wrap(err, "failed to create the Reva client")
	}
	rt.revaClient = client
	log.Info().Str("host", cfg.Reva.Host).Str("user", cfg.Reva.User).Msg("established Reva session")

	svr, err := server.New(cfg.WebserverPort, rt.revaClient, log)
	if err != nil {
		return errors.Wrap(err, "unable to create the web server")
	}
	rt.webServer = svr
	log.Info().Uint16("port", cfg.WebserverPort).Msg("webserver started")

	return nil
}

func (rt *Runtime) destroy() {

}

// Run starts the runtime's main loop.
func (rt *Runtime) Run() error {
	defer rt.destroy()

	// The stopSignal is used to intercept interruption signals to gracefully terminate the program
	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, syscall.SIGINT, syscall.SIGTERM)

loop:
	for {
		// Poll the stopSignal channel; if a signal was received, break the loop
		select {
		case <-stopSignal:
			rt.log.Info().Msg("shutting down")
			break loop

		default:
		}

		rt.tick()
		time.Sleep(runLoopSleeptime)
	}

	return nil
}

func (rt *Runtime) tick() {

}

// New creates a new runtime object.
func New(cfg Config, log *zerolog.Logger) (*Runtime, error) {
	rt := &Runtime{}
	if err := rt.initialize(cfg, log); err != nil {
		return nil, errors.Wrap(err, "unable to initialize the runtime")
	}
	return rt, nil
}
