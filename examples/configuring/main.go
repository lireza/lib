package main

import (
	"github.com/lireza/lib/configuring"
	"github.com/lireza/lib/logging"
)

func main() {
	logger := logging.NewLogger(logging.INFO)

	// Lets suppose we want to load logger.level and db.postgres.user configuration values.
	configs := configuring.New()

	// Here if LOGGER_LEVEL is defined as environment variable, config tries to load its value;
	// If not, and if logger_level is defined as command line argument, config tries to load its value;
	// If not, the default value provided as replacement is assigned to level.
	level := configs.Get("logger.level").StringOrElse("DEBUG")
	logger.Info("Loaded log level as %v", level)

	// Same as above, however instead of using a default value, this function tries to return an error if no key found.
	user, e := configs.Get("db.postgres.user").String()
	if e != nil {
		logger.Error(e.Error())
	} else {
		logger.Info("Loaded postgres user as %v", user)
	}

	// Config instances can be used to load configuration from a JSON file.
	configs, e = configuring.New().LoadJSON("config.json")
	if e != nil {
		logger.Error(e.Error())
	} else {
		// Here if LOGGER_LEVEL is defined as environment variable, config tries to load its value;
		// If not, and if logger_level is defined as command line argument, config tries to load its value;
		// If not, and if it is defined in provided JSON configuration file, it tries to load the value.
		// If not, the default value provided as replacement is assigned to level.
		level = configs.Get("logger.level").StringOrElse("DEBUG")
		logger.Info("Loaded log level as %v", level)

		// Same as above, however instead of using a default value, this function tries to return an error if no key found.
		user, e = configs.Get("db.postgres.user").String()
		if e != nil {
			logger.Error(e.Error())
		} else {
			logger.Info("Loaded postgres user as %v", user)
		}
	}
}
