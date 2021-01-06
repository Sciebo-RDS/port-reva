/* **********************************
 * Date: 2021-01-06
 * *********************************/

package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"

	"github.com/Sciebo-RDS/port-reva/cmd/logger"
	"github.com/Sciebo-RDS/port-reva/cmd/version"
	"github.com/Sciebo-RDS/port-reva/pkg/runtime"
)

var (
	portFlag = flag.Uint("port", 80, "the webserver port")
)

func main() {
	flag.Parse()
	printWelcome()

	run(logger.New())
}

func run(log *zerolog.Logger) {
	host, _ := os.Hostname()
	log.Info().Msgf("hostname: %s", host)
	log.Info().Msgf("webserver port: %v", *portFlag)

	rt, err := runtime.New(log)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create the runtime")
	}

	// The stopSignal is used to intercept interruption signals to gracefully terminate the program
	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, syscall.SIGINT, syscall.SIGTERM)

	if err := rt.Run(stopSignal); err != nil {
		log.Fatal().Err(err).Msg("fatal error during execution")
	}
}

func printWelcome() {
	fmt.Println("Sciebo RDS <-> Reva connector -- V" + version.GetString())
	fmt.Println("------------------------------------------------------------")
}
