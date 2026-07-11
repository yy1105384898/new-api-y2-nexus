package middleware

import "testing"

func TestImageInFlightCounterPerUserLimit(t *testing.T) {
	limiter := imageInFlightCounter{users: make(map[int]int)}
	if !limiter.tryAcquire(1, 2, 10) || !limiter.tryAcquire(1, 2, 10) {
		t.Fatal("first two requests should be admitted")
	}
	if limiter.tryAcquire(1, 2, 10) {
		t.Fatal("third request for the same user should be rejected")
	}
	if !limiter.tryAcquire(2, 2, 10) {
		t.Fatal("another user should still be admitted")
	}
	limiter.release(1)
	if !limiter.tryAcquire(1, 2, 10) {
		t.Fatal("released capacity should be reusable")
	}
}

func TestImageInFlightCounterGlobalLimit(t *testing.T) {
	limiter := imageInFlightCounter{users: make(map[int]int)}
	if !limiter.tryAcquire(1, 10, 2) || !limiter.tryAcquire(2, 10, 2) {
		t.Fatal("requests up to the global limit should be admitted")
	}
	if limiter.tryAcquire(3, 10, 2) {
		t.Fatal("request above the global limit should be rejected")
	}
	limiter.release(1)
	if !limiter.tryAcquire(3, 10, 2) {
		t.Fatal("released global capacity should be reusable")
	}
}
