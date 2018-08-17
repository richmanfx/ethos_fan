package main

import "testing"

func TestGetGpuQuantity(t *testing.T) {

	actualResult := getGpuQuantity()

	if actualResult < 0 {
		t.Errorf("GPU quantity was incorrect - less than zero: %d", actualResult)
	}

	if actualResult > 99 {
		t.Errorf("GPU quantity was incorrect - more than 99: %d", actualResult)
	}
}
