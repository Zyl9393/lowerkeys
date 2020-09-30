package lowerkeys_test

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/Zyl9393/lowerkeys"
)

func TestLowerKeysFrom(t *testing.T) {
	tests := []struct {
		header         http.Header
		expectedHeader lowerkeys.Header
	}{
		{nil, lowerkeys.Header{}},
		{http.Header{}, lowerkeys.Header{}},
		{http.Header{"content-type": {"foo"}}, lowerkeys.Header{"content-type": {"foo"}}},
		{http.Header{"Content-Type": {"foo"}}, lowerkeys.Header{"content-type": {"foo"}}},
		{http.Header{"Content-Type": {"foo"}, "Content-Length": {"3"}}, lowerkeys.Header{"content-type": {"foo"}, "content-length": {"3"}}},
	}
	for i, test := range tests {
		header := lowerkeys.From(test.header)
		if !reflect.DeepEqual(header, test.expectedHeader) {
			t.Errorf("Test #%d: header = %v. Expected %v.", i+1, header, test.expectedHeader)
		}
	}
}

func TestLowerKeysFromIsCopy(t *testing.T) {
	header := http.Header{
		"Content-Type":   {"foo"},
		"Content-Length": {"3"},
	}
	lowerKeysHeader := lowerkeys.From(header)
	header["Mark"] = []string{"Mutation"}
	if _, ok := lowerKeysHeader["Mark"]; ok {
		t.Fatalf("change in http.Header caused change in lowerkeys.Header created with lowerkeys.From()")
	}
}

func TestLowerKeysUsingIsSameMap(t *testing.T) {
	header := http.Header{
		"Content-Type":   {"foo"},
		"Content-Length": {"3"},
	}
	lowerKeysHeader := lowerkeys.Using(header)
	header["Mark"] = []string{"Mutation"}
	if _, ok := lowerKeysHeader["Mark"]; !ok {
		t.Fatalf("change in http.Header did not cause change in lowerkeys.Header created with lowerkeys.Using()")
	}
}

func TestLowerKeysHeaderAdd(t *testing.T) {
	header := lowerkeys.From(http.Header{
		"Content-Type":   {"foo"},
		"Content-Length": {"3"},
	})
	header.Add("Stop-Writing-Like-This", "Please")
	if _, ok := header["Stop-Writing-Like-This"]; ok {
		t.Errorf("key was added in non-lower original case")
	}
	if _, ok := header["stop-writing-like-this"]; !ok {
		t.Errorf("key was not added in lower-case")
	}
	if header.Get("sToP-wRiTiNg-LikE-tHiS") != "Please" {
		t.Errorf("added key not found")
	}
}

func TestLowerKeysHeaderDel(t *testing.T) {
	header := lowerkeys.From(http.Header{
		"Content-Type":   {"foo"},
		"Content-Length": {"3"},
	})
	header.Del("conTENT-TYpe")
	if _, ok := header["Content-Type"]; ok {
		t.Errorf("key was not deleted")
	}
	if header.Get("ConTent-tyPE") != "" {
		t.Errorf("deleted key is still found")
	}
}

func TestLowerKeysHeaderGet(t *testing.T) {
	header := lowerkeys.From(http.Header{
		"Content-Type":   {"foo"},
		"Content-Length": {"3"},
	})
	if header.Get("Content-Type") != "foo" {
		t.Errorf(`header.Get("Content-Type") != "foo"`)
	}
	if header.Get("Content-Length") != "3" {
		t.Errorf(`header.Get("Content-Length") != "3"`)
	}
	if header.Get("content-type") != "foo" {
		t.Errorf(`header.Get("content-type") != "foo"`)
	}
	if header.Get("content-length") != "3" {
		t.Errorf(`header.Get("content-length") != "3"`)
	}
}

func TestLowerKeysHeaderSet(t *testing.T) {
	header := lowerkeys.From(http.Header{
		"Content-Type":   {"foo"},
		"Content-Length": {"3"},
	})
	header.Set("contenT-TYPE", "BAR")
	header.Set("oTHerKEY", "value")
	if _, ok := header["content-type"]; !ok {
		t.Errorf("key replaced but not found")
	}
	val := header.Values("conTENT-type")
	if len(val) != 1 {
		t.Errorf("replaced one value but length was not 1")
	}
	if val[0] != "BAR" {
		t.Errorf("key found but wrong value")
	}
	val = header.Values("othErkey")
	if len(val) != 1 {
		t.Errorf("set one value but length was not 1")
	}
	if val[0] != "value" {
		t.Errorf("key found but wrong value")
	}
	if len(header) != 3 {
		t.Errorf("unexpected entry count in map")
	}
}

func TestLowerKeysHeaderValues(t *testing.T) {
	header := lowerkeys.From(http.Header{
		"Content-Type": {"foo", "bar"},
		"Nonsense":     {"3", "4"},
	})
	tests := []struct {
		headerName     string
		expectedValues map[string]int
	}{
		{"cOnTeNt-TyPe", map[string]int{"foo": 1, "bar": 1}},
		{"NonSenSe", map[string]int{"3": 1, "4": 1}},
	}
	for _, test := range tests {
		for _, value := range header.Values(test.headerName) {
			_, ok := test.expectedValues[value]
			if !ok {
				t.Errorf("unexpected or duplicate value '%s' for header name '%s'", value, test.headerName)
			} else {
				delete(test.expectedValues, value)
			}
		}
		if len(test.expectedValues) > 0 {
			t.Errorf("missing values for header name '%s': %v", test.headerName, test.expectedValues)
		}
	}
}

func TestLowerKeysFromPanic(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatalf("expected this test to panic, but it didn't.")
		}
	}()
	lowerkeys.From(http.Header{
		"Content-Type": {"foo"},
		"cOnTenT-tYpE": {"foo"},
	})
}

func TestLowerKeysUsingPanic(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatalf("expected this test to panic, but it didn't.")
		}
	}()
	lowerkeys.Using(http.Header{
		"Content-Type": {"foo"},
		"cOnTenT-tYpE": {"foo"},
	})
}
