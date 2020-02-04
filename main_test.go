package main

import "testing"

func TestGetJobs(t *testing.T) {
	got := getJobs(-1)
	if got != 1 {
		t.Errorf("Abs(-1) = %d; want 1", got)
	}
}
