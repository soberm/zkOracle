package zkOracle

import (
	"github.com/rs/zerolog"
	"io"
	"os"
)

var logger zerolog.Logger

func init() {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "15:04:05"}
	logger = zerolog.New(output).With().Timestamp().Logger()
}

func SetOutput(w io.Writer) {
	logger = logger.Output(w)
}

func Set(l zerolog.Logger) {
	logger = l
}

func Disable() {
	logger = zerolog.Nop()
}

func Logger() zerolog.Logger {
	return logger
}
