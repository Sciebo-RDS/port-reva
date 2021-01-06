/* **********************************
 * Date: 2021-01-06
 * *********************************/

package logger

import (
	"os"

	"github.com/rs/zerolog"
)

// New creates a new default console logger.
func New() *zerolog.Logger {
	zl := zerolog.New(os.Stdout).With().Timestamp().Logger()
	zl = zl.Level(zerolog.DebugLevel)
	zl = zl.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "2006-01-02 15:04:05.999"})

	return &zl
}
