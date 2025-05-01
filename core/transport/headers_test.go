package transport

import (
	"net/http"
	"reflect"
	"testing"
)

func TestNewHeaders_NormalInput(t *testing.T) {
	input := []string{"Content-Type=application/json", "Accept=*/*"}
	want := http.Header{
		"Content-Type": {"application/json"},
		"Accept":       {"*/*"},
	}
	got := NewHeaders(input)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("NewHeaders() = %v, want %v", got, want)
	}
}

func TestNewHeaders_InvalidFormat(t *testing.T) {
	input := []string{"invalid-header", "AnotherOne"}
	want := http.Header{}
	got := NewHeaders(input)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("NewHeaders() = %v, want %v", got, want)
	}
}

func TestNewHeaders_DuplicateKeys(t *testing.T) {
	input := []string{"X-Test=1", "X-Test=2"}
	want := http.Header{
		"X-Test": {"1", "2"},
	}
	got := NewHeaders(input)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("NewHeaders() = %v, want %v", got, want)
	}
}

func TestNewHeaders_EmptyOrNilInput(t *testing.T) {
	var inputNil []string
	inputEmpty := []string{}
	want := http.Header{}

	gotNil := NewHeaders(inputNil)
	if !reflect.DeepEqual(gotNil, want) {
		t.Errorf("NewHeaders(nil) = %v, want %v", gotNil, want)
	}

	gotEmpty := NewHeaders(inputEmpty)
	if !reflect.DeepEqual(gotEmpty, want) {
		t.Errorf("NewHeaders(empty) = %v, want %v", gotEmpty, want)
	}
}
