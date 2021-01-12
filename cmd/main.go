/* **********************************
 * Date: 2021-01-06
 * *********************************/

package main

import (
	"flag"
	"os"

	"github.com/rs/zerolog"

	"github.com/Sciebo-RDS/port-reva/cmd/cmdline"
	"github.com/Sciebo-RDS/port-reva/cmd/logger"
	"github.com/Sciebo-RDS/port-reva/cmd/version"
	"github.com/Sciebo-RDS/port-reva/pkg/runtime"
)

const (
	portFlagName     = "port"
	hostFlagName     = "host"
	userFlagName     = "user"
	passwordFlagName = "password"

	hostEnvName     = "RDS_REVA_HOST"
	userEnvName     = "RDS_REVA_USER"
	passwordEnvName = "RDS_REVA_PASSWORD"
)

var (
	portFlag     = flag.Uint(portFlagName, 80, "the webserver port")
	hostFlag     = flag.String(hostFlagName, "", "the Reva host (<host>:<port>)")
	userFlag     = flag.String(userFlagName, "", "the user name to log in to Reva")
	passwordFlag = flag.String(passwordFlagName, "", "the user password to log in to Reva")
)

func main() {
	flag.Parse()

	log := logger.New()
	printWelcome(log)

	run(log)
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

	// Read settings from command-line
	cfg.WebserverPort = uint16(*portFlag)
	cfg.Reva.Host = *hostFlag
	cfg.Reva.User = *userFlag
	cfg.Reva.Password = *passwordFlag

	// Read (missing) settings from environment variables
	setFromEnvironment := func(value *string, flagName string, envName string) {
		if !cmdline.IsFlagSet(flagName) {
			if val, ok := os.LookupEnv(envName); ok {
				*value = val
			}
		}
	}

	setFromEnvironment(&cfg.Reva.Host, hostFlagName, hostEnvName)
	setFromEnvironment(&cfg.Reva.User, userFlagName, userEnvName)
	setFromEnvironment(&cfg.Reva.Password, passwordFlagName, passwordEnvName)

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

func printWelcome(log *zerolog.Logger) {
	log.Info().Str("version", version.GetString()).Msg("Sciebo RDS to Reva connector")
}
