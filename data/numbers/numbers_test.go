package numbers

import (
	"testing"
)

func Test_getPrefix_3common9(t *testing.T) {
	from, to := getPrefix("123000", "123999")
	assertStringEqual(t, "123", from)
	assertStringEqual(t, "123", to)
}

func Test_getPrefix_1common9(t *testing.T) {
	from, to := getPrefix("123450", "123459")
	assertStringEqual(t, "12345", from)
	assertStringEqual(t, "12345", to)
}

func Test_getPrefix_nocommon9(t *testing.T) {
	from, to := getPrefix("123456", "123457")
	assertStringEqual(t, "123456", from)
	assertStringEqual(t, "123457", to)
}

func Test_getSingles_to_greater_from(t *testing.T) {
	expected := []string{"123", "124", "125"}
	actual, _ := getSingles("123", "125")
	for i, exp := range expected {
		if exp != actual[i] {
			t.Errorf("error in list at %d expected '%s', got '%s'", i, exp, actual[i])
		}
	}
}

func Test_getSingles_to_equal_from(t *testing.T) {
	expected := []string{"123"}
	actual, _ := getSingles("123", "123")

	for i, exp := range expected {
		if exp != actual[i] {
			t.Errorf("error in list at %d expected '%s', got '%s'", i, exp, actual[i])
		}
	}
}

func assertStringEqual(t *testing.T, expected, actual string) {
	if expected != actual {
		t.Errorf("expected '%s', got '%s'", expected, actual)
	}
}
