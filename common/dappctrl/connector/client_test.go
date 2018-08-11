package connector

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func testServer(t *testing.T, args, result interface{}) *httptest.Server {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rawResult, err := json.Marshal(result)
			if err != nil {
				t.Fatal(err)
			}

			data := &Response{Result: rawResult}

			w.Header().Set("Content-Type",
				"application/json; charset=utf-8")
			w.WriteHeader(http.StatusOK)

			if err := json.NewEncoder(w).Encode(data); err != nil {
				t.Fatal(err)
			}
		}))
	return ts
}

func TestPostRequest(t *testing.T) {
	expResult := map[string]interface{}{"A": "123",
		"B": float64(456), "C": true}
	args := map[string]interface{}{"arg1": "val1", "arg2": 2}

	ts := testServer(t, args, expResult)
	defer ts.Close()

	data, err := json.Marshal(args)
	if err != nil {
		t.Fatal(err)
	}

	var resp map[string]interface{}
	if err := post(httpClient(DefaultConfig()),
		ts.URL, "admin", "pass", data, &resp); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expResult, resp) {
		t.Fatalf("expect result %v, got %v", expResult, resp)
	}
}
