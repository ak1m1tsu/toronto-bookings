package types

import (
	"fmt"
	"testing"
)

type TestCase struct {
	Input any
	Want  any
}

func TestIsFieldValidWithEmailRegexPattern(t *testing.T) {
	testCases := []TestCase{
		{"roman@gmail.com", true},
		{"roman@gmail.co", true},
		{"roman@gmail.c", false},
		{"roman@gmailcom", false},
		{"romangmail.com", false},
		{"roman\\|/@gmail.com", false},
		{"@gmail.com", false},
		{"r@gmail.com", true},
		{"roman@", false},
		{"roman_+_roman@gmail.com", true},
		{"test.test.tes.tes.set.set.tes.t.setset..setset@gmail.com", true},
		{"+@gmail.com", true},
	}

	for _, tc := range testCases {
		actual := isFieldValid(tc.Input.(string), emailRegexPattern)
		if actual != tc.Want {
			t.Errorf("for %s validation failed. got: %v, want: %v", tc.Input, actual, tc.Want)
		}
	}
}

func TestIsFieldValidWithPhoneNumberRegexPattern(t *testing.T) {
	testCases := []TestCase{
		{"89501262318", true},
		{"8 950 126 23 18", true},
		{"8 (950) 126 23 18", true},
		{"8 (950) 126-23-18", true},
		{"8 (950)-126-23-18", true},
		{"8950 126 23 18", true},
		{"8(950) 126 23 18", true},
		{"8(950) 126-23-18", true},
		{"8(950)-126-23-18", true},
		{"+79501262318", true},
		{"+7 950 126 23 18", true},
		{"+7 (950) 126 23 18", true},
		{"+7 (950) 126-23-18", true},
		{"+7 (950)-126-23-18", true},
		{"+7950 126 23 18", true},
		{"+7(950) 126 23 18", true},
		{"+7(950) 126-23-18", true},
		{"+7(950)-126-23-18", true},
		{"+7(950241124)-126-23-18", false},
		{"+7(950)-126421424123-18", false},
		{"421532", false},
		{"329582309857023985702398", false},
		{"+7(950)-126+23+18", false},
	}

	for _, tc := range testCases {
		actual := isFieldValid(tc.Input.(string), phoneNumberRegexPattern)
		if actual != tc.Want {
			t.Errorf("for %s validation failed. got: %v, want: %v", tc.Input, actual, tc.Want)
		}
	}
}

func TestNormalizePhoneNumber(t *testing.T) {
	testCases := []TestCase{
		{"89501262318", "89501262318"},
		{"8 950 126 23 18", "89501262318"},
		{"8 (950) 126 23 18", "89501262318"},
		{"8 (950) 126-23-18", "89501262318"},
		{"8 (950)-126-23-18", "89501262318"},
		{"8950 126 23 18", "89501262318"},
		{"8(950) 126 23 18", "89501262318"},
		{"8(950) 126-23-18", "89501262318"},
		{"8(950)-126-23-18", "89501262318"},
		{"+79501262318", "79501262318"},
		{"+7 950 126 23 18", "79501262318"},
		{"+7 (950) 126 23 18", "79501262318"},
		{"+7 (950) 126-23-18", "79501262318"},
		{"+7 (950)-126-23-18", "79501262318"},
		{"+7950 126 23 18", "79501262318"},
		{"+7(950) 126 23 18", "79501262318"},
		{"+7(950) 126-23-18", "79501262318"},
		{"+7(950)-126-23-18", "79501262318"},
	}

	for _, tc := range testCases {
		t.Run(tc.Input.(string), func(t *testing.T) {
			actual := normalizePhoneNumber(tc.Input.(string))
			if actual != tc.Want {
				t.Errorf("phone number is not corrected. got: %s, want: %v", actual, tc.Want)
			}
		})
	}
}

func TestValidateCreateReservationRequest(t *testing.T) {
	testCases := []TestCase{
		{&CreateReservationRequest{FirstName: ""}, FirstNameValidationError},
		{&CreateReservationRequest{FirstName: "123", LastName: ""}, LastNameValidationError},
		{&CreateReservationRequest{FirstName: "123", LastName: "123", Email: "resa@"}, EmailAddressValidationError},
		{&CreateReservationRequest{FirstName: "123", LastName: "123", Email: "resa@asd.com", PhoneNumber: "521521"}, PhoneNumberValidationError},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("validate create reservation request test case #%d", i+1), func(t *testing.T) {
			err := ValidateCreateReservationRequest(tc.Input.(*CreateReservationRequest))
			if err != tc.Want {
				t.Errorf("validate failed. got: %s, want: %s", err.Error(), tc.Want)
			}
		})
	}
}
