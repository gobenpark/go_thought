package log

import "testing"

func TestZapLogger_Info(t *testing.T) {
	l := NewZapLogger(false)
	l.Info("test", "data", "1")
}
