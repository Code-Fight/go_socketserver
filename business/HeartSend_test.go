package business

import "testing"

func TestHeardEvent(t *testing.T) {
	testChan :=make(chan struct{})
	HeardEvent(nil,1,testChan)
}