/* **********************************
 * Date: 2021-01-06
 * *********************************/

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/rs/zerolog"

	"github.com/Sciebo-RDS/port-reva/cmd/logger"
	"github.com/Sciebo-RDS/port-reva/cmd/version"
	"github.com/Sciebo-RDS/port-reva/pkg/runtime"
)

var (
	portFlag = flag.Uint("port", 80, "the webserver port")
	hostFlag = flag.String("host", "", "the Reva host (<host>:<port>)")
	userFlag = flag.String("user", "", "the user name to log in to Reva")
	passFlag = flag.String("pass", "", "the user password to log in to Reva")
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

	cfg := getRuntimeConfig()
	verifyRuntimeConfig(&cfg, log)

	rt, err := runtime.New(cfg, log)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create the runtime")
	}

	// Start the runtime, which will enter its own loop
	if err := rt.Run(); err != nil {
		log.Fatal().Err(err).Msg("fatal error during execution")
	}
}

func getRuntimeConfig() runtime.Config {
	cfg := runtime.Config{}
	cfg.WebserverPort = uint16(*portFlag)
	cfg.Reva.Host = *hostFlag
	cfg.Reva.User = *userFlag
	cfg.Reva.Password = *passFlag
	return cfg
}

func verifyRuntimeConfig(cfg *runtime.Config, log *zerolog.Logger) {
	if cfg.Reva.Host == "" {
		log.Fatal().Msg("no Reva host specified")
	}

	if cfg.Reva.User == "" {
		log.Fatal().Msg("no Reva user specified")
	}

	if cfg.Reva.Password == "" {
		log.Fatal().Msg("no Reva password specified")
	}
}

func printWelcome() {
	fmt.Println("Sciebo RDS <-> Reva connector -- V" + version.GetString())
	fmt.Println("------------------------------------------------------------")
}
