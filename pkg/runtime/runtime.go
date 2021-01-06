/* **********************************
 * Date: 2021-01-06
 * *********************************/

package runtime

import (
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// Runtime implements the main program runtime.
type Runtime struct {
	log *zerolog.Logger
}

const (
	runLoopSleeptime = time.Millisecond * 250
)

func (rt *Runtime) initialize(log *zerolog.Logger) error {
	if log == nil {
		return errors.Errorf("no logger specified")
	}
	rt.log = log
	rt.log.Debug().Msg("initializing runtime")

	return nil
}

func (rt *Runtime) destroy() {
	rt.log.Debug().Msg("destroying runtime")
}

// Run starts the runtime's main loop.
func (rt *Runtime) Run(stopSignal <-chan os.Signal) error {
	rt.log.Info().Msg("starting runtime")

	defer rt.destroy()

loop:
	for {
		// Poll the stopSignal channel; if a signal was received, break the loop
		select {
		case <-stopSignal:
			rt.log.Info().Msg("shutting down runtime")
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
func New(log *zerolog.Logger) (*Runtime, error) {
	rt := &Runtime{}
	if err := rt.initialize(log); err != nil {
		return nil, errors.Wrap(err, "unable to initialize the runtime")
	}
	return rt, nil
}
