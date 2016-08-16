package device

import (
	"bytes"
	"fmt"
	kitlog "github.com/go-kit/kit/log"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"golang.org/x/net/context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var (
	db  *sqlx.DB
	err error
	ds  DeviceRepository
	svc Service
)

func setup() {
	db, err = sqlx.Open("postgres", "user=devicestore password=devicestore dbname=devicestore sslmode=disable")
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	db.Exec("DELETE FROM devices")

	ds = NewRepository(db)
	svc = NewService(ds)
}

func teardown() {
	db.Close()
}

// a face io.ReadCloser for constructing request Body
type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

func TestServiceHandler(t *testing.T) {
	setup()
	defer teardown()

	reqBody := []byte(`{
		"data": {
			"type": "devices",
			"attributes": {
				"name": "Test Device",
				"udid": "00000000-1111-2222-3333-444455556666",
				"serial_number": "ABCDEFG0000111",
				"os_version": "9.3.3"
			}
		}
	}`)

	req, err := http.NewRequest("POST", "/devices", &nopCloser{bytes.NewBuffer(reqBody)})
	req.Header.Set("Content-Type", "application/vnd.api+json")
	req.Header.Set("Accept", "application/vnd.api+json")

	fmt.Printf("%v\n", req)

	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	ctx := context.Background()
	logger := kitlog.NewLogfmtLogger(os.Stdout)

	handler := MakeHTTPHandler(ctx, svc, logger)
	handler.ServeHTTP(w, req)

	if w.Code != 201 {
		t.Errorf("expected http status: 201 created, got: %v", w.Code)
	}

	if w.Header().Get("Location") == "" {
		t.Error("expected Location header")
	}

	fmt.Println(w.Body)

}
