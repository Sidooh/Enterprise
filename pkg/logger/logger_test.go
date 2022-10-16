package logger

import (
	"os"
	"testing"
)

func TestLoggerInit(t *testing.T) {

	if ClientLog.Out != nil {
		t.Errorf("Init() = %v; want nil", ClientLog.Out)
	}

	if ClientLog.Out != nil {
		t.Errorf("Init() = %v; want nil", ClientLog.Out)
	}

	err := os.Setenv("APP_ENV", "TEST")
	if err != nil {
		return
	}

	Init()

	if ClientLog.Out == nil {
		t.Errorf("Init() = %v; want value", ClientLog.Out)
	}

	if ClientLog.Level.String() != "info" {
		t.Errorf("Init() = %v; want info", ClientLog.Level)
	}

}
