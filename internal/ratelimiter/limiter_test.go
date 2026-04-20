package ratelimiter

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestLimiterBasicAllowBlock(t *testing.T) {

	limiter := NewLimiter(2, 2*time.Second)
	ip := "127.0.0.1"

	// First 2 → allowed
	if !limiter.Allow(ip) {
		t.Error("expected first request to pass")
	}

	if !limiter.Allow(ip) {
		t.Error("expected second request to pass")
	}

	// Third → blocked
	if limiter.Allow(ip) {
		t.Error("expected third request to be blocked")
	}
}

func TestLimiterWindowReset(t *testing.T) {

	limiter := NewLimiter(1, 1*time.Second)
	ip := "127.0.0.1"

	if !limiter.Allow(ip) {
		t.Error("expected first request")
	}

	if limiter.Allow(ip) {
		t.Error("expected second request to be blocked")
	}

	// wait for window reset
	time.Sleep(1100 * time.Millisecond)

	if !limiter.Allow(ip) {
		t.Error("expected request after window reset")
	}
}

func TestLimiterMultipleIPs(t *testing.T) {

	limiter := NewLimiter(1, 2*time.Second)

	ip1 := "127.0.0.1"
	ip2 := "192.168.1.1"

	if !limiter.Allow(ip1) {
		t.Error("ip1 should pass")
	}

	if !limiter.Allow(ip2) {
		t.Error("ip2 should pass independently")
	}
}

func TestLimiterConcurrency(t *testing.T) {

	limiter := NewLimiter(5, 2*time.Second)
	ip := "127.0.0.1"

	const totalRequests = 20
	var successCount int32

	done := make(chan bool)

	for i := 0; i < totalRequests; i++ {
		go func() {
			if limiter.Allow(ip) {
				atomic.AddInt32(&successCount, 1)
			}
			done <- true
		}()
	}

	for i := 0; i < totalRequests; i++ {
		<-done
	}

	if successCount > 5 {
		t.Errorf("expected max 5 allowed, got %d", successCount)
	}
}

func TestLimiterZeroLimit(t *testing.T) {

	limiter := NewLimiter(0, 1*time.Second)
	ip := "127.0.0.1"

	if limiter.Allow(ip) {
		t.Error("expected no requests to be allowed")
	}
}