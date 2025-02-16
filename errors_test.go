package easyhttp

import (
	"log"
	"testing"
	"time"
)

func TestInvalidLength(t *testing.T) {
	tearDown := setupServer(t)
	defer tearDown(t)
	client := NewHTTPClient()

	request, err := NewRequest("http://localhost:1234/path")
	if err != nil {
		t.Fatal(err.Error())
	}
	request.version = "1.0"
	request.SetHeader("Content-Length", "avc")
	response, err := client.GET(request)
	if err != nil {
		t.Fatal(err.Error())
	}
	if response.version != "1.0" {
		t.Fatalf("HTTP VERSION IS WRONG")
	}

	if response.StatusCode != STATUS_LENGTH_REQUIRED {
		t.FailNow()
	}

}

func TestInvalidMethod(t *testing.T) {
	tearDown := setupServer(t)
	defer tearDown(t)
	client := NewHTTPClient()

	request, err := NewRequest("http://localhost:1234/path")
	if err != nil {
		t.Fatal(err.Error())
	}
	request.version = "1.0"
	request.method = "AVC"
	response, err := client.sendRequest(request)
	if err != nil {
		t.Fatal(err.Error())
	}
	if response.version != "1.0" {
		t.Fatalf("HTTP VERSION IS WRONG")
	}

	if response.StatusCode != STATUS_METHOD_NOT_ALLOWED {
		t.FailNow()
	}

}

func TestUnsupportedVersion(t *testing.T) {
	tearDown := setupServer(t)
	defer tearDown(t)
	client := NewHTTPClient()

	request, err := NewRequest("http://localhost:1234/path")
	if err != nil {
		t.Fatal(err.Error())
	}

	request.version = "2.0"
	response, err := client.GET(request)
	if err != nil {
		t.Fatal(err.Error())
	}

	if response.StatusCode != STATUS_HTTP_VERSION_NOT_SUPPORTED {
		t.FailNow()
	}

}

func TestPanicOnHandler(t *testing.T) {
	tearDown := setupServer(t)
	defer tearDown(t)
	client := NewHTTPClient()

	request, err := NewRequest("http://localhost:1234/panic")
	if err != nil {
		t.Fatal(err.Error())
	}

	request.CloseConnection()
	response, err := client.GET(request)
	if err != nil {
		t.Fatal(err.Error())
	}

	if response.StatusCode != STATUS_INTERNAL_ERROR {
		log.Println(response.StatusCode)
		t.FailNow()
	}

}

func TestTimeoutOnHandler(t *testing.T) {
	tearDown := setupServer(t)
	defer tearDown(t)
	client := NewHTTPClient()

	request, err := NewRequest("http://localhost:1234/timeout")
	if err != nil {
		t.Fatal(err.Error())
	}
	request.CloseConnection()
	response, err := client.GET(request)
	if err != nil {
		t.Fatal(err.Error())
	}

	if response.StatusCode != STATUS_REQUEST_TIMEOUT {
		t.FailNow()
	}

}

func TestTimeoutOnClient(t *testing.T) {
	tearDown := setupServer(t)
	defer tearDown(t)
	client := NewHTTPClient()

	request, err := NewRequest("http://localhost:1234/timeout")
	if err != nil {
		t.Fatal(err.Error())
	}
	request.CloseConnection()
	request.SetTimeout(time.Duration(2000) * time.Millisecond)
	_, err = client.GET(request)
	if err != ErrClientTimeout {
		t.Fatalf("Got wrong error %s\n", err.Error())
	}

}

func TestNotFound(t *testing.T) {
	tearDown := setupServer(t)
	defer tearDown(t)
	client := NewHTTPClient()

	request, err := NewRequest("http://localhost:1234/notfound")
	if err != nil {
		t.Fatal(err.Error())
	}
	response, err := client.GET(request)
	if err != nil {
		t.Fatal(err.Error())
	}
	if response.StatusCode != STATUS_NOT_FOUND {
		t.Fatalf("Got wrong STATUS %d\n", response.StatusCode)
	}

}

func TestMethodNotAllowed(t *testing.T) {
	tearDown := setupServer(t)
	defer tearDown(t)
	client := NewHTTPClient()

	request, err := NewRequest("http://localhost:1234/")
	if err != nil {
		t.Fatal(err.Error())
	}
	response, err := client.PUT(request)
	if err != nil {
		t.Fatal(err.Error())
	}
	if response.StatusCode != STATUS_METHOD_NOT_ALLOWED {
		t.Fatalf("Got wrong STATUS %d\n", response.StatusCode)
	}

	if !response.HasHeaderValue("Allow", "GET") {
		t.Fatalf("Missing allow header")
	}

}
