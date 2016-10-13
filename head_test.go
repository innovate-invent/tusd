package tusd_test

import (
	"net/http"
	"os"
	"testing"

	. "github.com/tus/tusd"
)

func TestHead(t *testing.T) {
	SubTest(t, "Successful HEAD request", func(t *testing.T, store *MockFullDataStore) {
		store.EXPECT().GetInfo("yes").Return(FileInfo{
			Offset: 11,
			Size:   44,
			MetaData: map[string]string{
				"name": "lunrjs.png",
				"type": "image/png",
			},
		}, nil)

		handler, _ := NewHandler(Config{
			DataStore: store,
		})

		res := (&httpTest{
			Name:   "Successful request",
			Method: "HEAD",
			URL:    "yes",
			ReqHeader: map[string]string{
				"Tus-Resumable": "1.0.0",
			},
			Code: http.StatusOK,
			ResHeader: map[string]string{
				"Upload-Offset": "11",
				"Upload-Length": "44",
				"Cache-Control": "no-store",
			},
		}).Run(handler, t)

		// Since the order of a map is not guaranteed in Go, we need to be prepared
		// for the case, that the order of the metadata may have been changed
		if v := res.Header().Get("Upload-Metadata"); v != "name bHVucmpzLnBuZw==,type aW1hZ2UvcG5n" &&
			v != "type aW1hZ2UvcG5n,name bHVucmpzLnBuZw==" {
			t.Errorf("Expected valid metadata (got '%s')", v)
		}
	})

	SubTest(t, "Non-existing file", func(t *testing.T, store *MockFullDataStore) {
		store.EXPECT().GetInfo("no").Return(FileInfo{}, os.ErrNotExist)

		handler, _ := NewHandler(Config{
			DataStore: store,
		})

		res := (&httpTest{
			Method: "HEAD",
			URL:    "no",
			ReqHeader: map[string]string{
				"Tus-Resumable": "1.0.0",
			},
			Code: http.StatusNotFound,
			ResHeader: map[string]string{
				"Content-Length": "0",
			},
		}).Run(handler, t)

		if string(res.Body.Bytes()) != "" {
			t.Errorf("Expected empty body for failed HEAD request")
		}
	})
}
