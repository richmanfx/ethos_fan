package main

import (
	"testing"
)

func TestGetGpuQuantity(t *testing.T) {

	actualResult := getGpuQuantity()

	if actualResult < 0 {
		t.Errorf("GPU quantity was incorrect - less than zero: %d", actualResult)
	}

	if actualResult > 99 {
		t.Errorf("GPU quantity was incorrect - more than 99: %d", actualResult)
	}
}

func TestCheckValidInRange(t *testing.T) {

	type TestData struct {
		minimum        int
		value          int
		maximum        int
		expectedResult bool
	}

	var testDataSlice = make([]TestData, 0, 10)

	testDataSlice = append(testDataSlice,
		TestData{minimum: 0, value: 50, maximum: 100, expectedResult: true},
		TestData{minimum: 30, value: 31, maximum: 32, expectedResult: true},
		TestData{minimum: 30, value: 69, maximum: 70, expectedResult: true},
		TestData{minimum: 30, value: 90, maximum: 80, expectedResult: false},
		TestData{minimum: 30, value: 90, maximum: 10, expectedResult: false},
		TestData{minimum: 30, value: 90, maximum: 10, expectedResult: false},
		TestData{minimum: 30, value: 90, maximum: 29, expectedResult: false},
		TestData{minimum: 30, value: 91, maximum: 90, expectedResult: false},
		TestData{minimum: 30, value: 30, maximum: 90, expectedResult: true},
		TestData{minimum: 30, value: 90, maximum: 90, expectedResult: true},
	)

	// Пробежать по всем данным
	for _, testDataItem := range testDataSlice {

		actualResult := checkValidInRange(
			testDataItem.minimum,
			testDataItem.maximum,
			testDataItem.value)

		if actualResult != testDataItem.expectedResult {
			t.Errorf("Invalid check value '%d' in range (%d...%d) - result '%v'",
				testDataItem.value, testDataItem.minimum, testDataItem.maximum, actualResult)
		}
	}
}

func TestSetNewFanSpeedForAllGpu(t *testing.T) {

	type TestData struct {
		gpuQuantity  int
		initFanSpeed int
		lowTemp      int
		highTemp     int
		speedStep    int
		minFanSpeed  int
		expectedResult
	}

}
