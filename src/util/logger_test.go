package util

import "testing"

func TestLog(b *testing.T) {
	log := NewLogger("./", "test")
	log.Debug("22222")
}
