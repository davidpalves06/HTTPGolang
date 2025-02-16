package easyhttp

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClientGet(t *testing.T) {
	client := NewHTTPClient()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("TestHeader", "Hello")
		w.WriteHeader(STATUS_OK)
		w.Write([]byte("Hello World!\n"))
	}))
	defer server.Close()

	request, err := NewRequest(server.URL + "/path")
	if err != nil {
		t.Fatal(err.Error())
	}

	response, err := client.GET(request)
	if err != nil || response == nil {
		t.Fatal(err.Error())
	}

	if response.StatusCode != STATUS_OK {
		t.Fatalf("Wrong Status Code: %d\n", response.StatusCode)
	}

	if !response.HasHeaderValue("TestHeader", "Hello") {
		t.Fatalf("Wrong Header\n")
	}
	bodyBuffer := make([]byte, 1024)
	bodyLength, _ := response.Read(bodyBuffer)
	if string(bodyBuffer[:bodyLength]) != "Hello World!\n" {
		t.Fatalf("Wrong Body")
	}
}

func TestClientPost(t *testing.T) {
	client := NewHTTPClient()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		body, _ := io.ReadAll(r.Body)
		w.Header().Set("TestHeader", "Hello")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(body))
	}))
	defer server.Close()

	request, err := NewRequestWithBody(server.URL+"/resource", []byte("Echo me this"))
	if err != nil {
		t.Fatal(err.Error())
	}
	response, err := client.POST(request)
	if err != nil {
		t.Fatal(err.Error())
	}

	if response.StatusCode != STATUS_OK {
		t.Fatalf("Wrong Status code")
	}
	if !response.HasHeaderValue("TestHeader", "Hello") {
		t.Fatalf("Wrong Header")
	}
	bodyBuffer := make([]byte, 1024)
	bodyLength, _ := response.Read(bodyBuffer)
	if string(bodyBuffer[:bodyLength]) != "Echo me this" {
		t.Fatalf("Wrong Body")
	}
}
