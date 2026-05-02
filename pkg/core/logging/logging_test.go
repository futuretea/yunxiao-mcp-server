package logging

import (
	"bytes"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func preserveLogger() func() {
	globalLevel := zerolog.GlobalLevel()
	logger := log.Logger
	timeFieldFormat := zerolog.TimeFieldFormat
	return func() {
		zerolog.SetGlobalLevel(globalLevel)
		log.Logger = logger
		zerolog.TimeFieldFormat = timeFieldFormat
	}
}

func TestInitializeSetsLevelAndOutput(t *testing.T) {
	restore := preserveLogger()
	defer restore()

	var buf bytes.Buffer
	if err := Initialize("debug", &buf); err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}

	if zerolog.GlobalLevel() != zerolog.DebugLevel {
		t.Fatalf("global level = %v, want debug", zerolog.GlobalLevel())
	}
	log.Info().Msg("hello")
	if !bytes.Contains(buf.Bytes(), []byte("hello")) {
		t.Fatalf("output = %q, missing hello", buf.String())
	}
}

func TestInitializeDefaultsToStderr(t *testing.T) {
	restore := preserveLogger()
	defer restore()

	if err := Initialize("info", nil); err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}
}

func TestInitializeRejectsInvalidLevel(t *testing.T) {
	restore := preserveLogger()
	defer restore()

	if err := Initialize("not-a-level", &bytes.Buffer{}); err == nil {
		t.Fatal("Initialize() expected invalid level error")
	}
}

func TestDisableSilencesLogging(t *testing.T) {
	restore := preserveLogger()
	defer restore()

	var buf bytes.Buffer
	if err := Initialize("info", &buf); err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}

	Disable()

	if zerolog.GlobalLevel() != zerolog.Disabled {
		t.Fatalf("global level = %v, want disabled", zerolog.GlobalLevel())
	}
}
