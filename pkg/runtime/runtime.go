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

	"github.com/Sciebo-RDS/port-reva/pkg/server"
	"github.com/Sciebo-RDS/port-reva/pkg/service"
)

// Runtime implements the main program runtime.
type Runtime struct {
	log *zerolog.Logger

	conf Config

	webServer *server.WebServer
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

	svr, err := server.New(cfg.WebserverPort, rt.conf.Reva, log)
	if err != nil {
		return errors.Wrap(err, "unable to create the web server")
	}
	rt.webServer = svr
	log.Info().Uint16("port", cfg.WebserverPort).Msg("webserver started")

	// If all initialization went through, register the service with the token storage
	if err := rt.registerWithTokenStorage(); err != nil {
		log.Warn().Err(err).Msg("unable to register connector with token storage")
	}

	return nil
}

func (rt *Runtime) destroy() {

}

func (rt *Runtime) registerWithTokenStorage() error {
	endpoint, ok := os.LookupEnv("CENTRAL_SERVICE_TOKEN_STORAGE")
	if !ok {
		return errors.Errorf("token storage endpoint not set")
	}
	endpoint += "/service"

	if err := service.RegisterService(endpoint); err != nil {
		return errors.Wrap(err, "unable to register service")
	}

	return nil
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
		case c := <-stopSignal:
			rt.log.Info().Msgf("shutting down (%v)", c.String())
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
