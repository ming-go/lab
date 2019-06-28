package hraft

import (
	"testing"
	"time"
)

func TestHRaft(t *testing.T) {
	go New("123", "127.0.0.1:8787", false)
	//go New("456", "127.0.0.1:8787", true)

	<-time.After(10 * time.Minute)
}
