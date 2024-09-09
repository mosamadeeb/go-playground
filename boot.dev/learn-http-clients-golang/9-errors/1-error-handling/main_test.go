package main

import (
	"fmt"
	"testing"
)

func TestFetchData(t *testing.T) {
	type testCase struct {
		inputUrl       string
		expectedError  error
		expectedStatus int
	}

	tests := []testCase{
		{
			inputUrl:       "https://api.boot.dev/v1/courses_rest_api/learn-http/items",
			expectedError:  nil,
			expectedStatus: 200,
		},
		{
			inputUrl:       "https://api.boot.dev/v1/wrong-path",
			expectedError:  fmt.Errorf("non-OK HTTP status: 404 Not Found"),
			expectedStatus: 404,
		},
	}

	if withSubmit {
		tests = append(tests, []testCase{
			{
				inputUrl:       "://example.com",
				expectedError:  fmt.Errorf("network error: parse \"://example.com\": missing protocol scheme"),
				expectedStatus: 0,
			},
			{
				inputUrl:       "https://api.boot.dev/v1/courses_rest_api/learn-http/locations",
				expectedError:  nil,
				expectedStatus: 200,
			},
		}...)
	}

	passCount := 0
	failCount := 0

	for _, test := range tests {
		actualStatus, actualError := fetchData(test.inputUrl)
		if (actualError == nil && test.expectedError != nil) || (actualError != nil && test.expectedError == nil) || (actualError != nil && test.expectedError != nil && actualError.Error() != test.expectedError.Error()) || actualStatus != test.expectedStatus {
			failCount++
			t.Errorf(`---------------------------------
URL:		%v
Expecting:  Error: %v - Status code: %v
Actual:     Error: %v - Status code: %v
Fail
`, test.inputUrl, test.expectedError, test.expectedStatus, actualError, actualStatus)

		} else {
			passCount++
			fmt.Printf(`---------------------------------
URL:		%v
Expecting:  Error: %v - Status code: %v
Actual:     Error: %v - Status code: %v
Pass
`, test.inputUrl, test.expectedError, test.expectedStatus, actualError, actualStatus)
		}
	}

	fmt.Println("---------------------------------")
	fmt.Printf("%d passed, %d failed\n", passCount, failCount)
}

var withSubmit = true
