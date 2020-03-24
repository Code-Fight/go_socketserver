package units

import "testing"

func TestSafeCloseChan(t *testing.T) {
	ch :=make(chan struct{})
	close(ch)
	SafeCloseChan(ch)
}
