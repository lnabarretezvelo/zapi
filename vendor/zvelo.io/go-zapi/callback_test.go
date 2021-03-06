package zapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"

	"zvelo.io/msg"
)

func callbackHandler(m **msg.QueryResult) Handler {
	return HandlerFunc(func(in *msg.QueryResult) {
		*m = in
	})
}

func TestCallbackHandler(t *testing.T) {
	var buf bytes.Buffer
	var m *msg.QueryResult

	srv := httptest.NewServer(CallbackHandler(callbackHandler(&m), WithDebug(&buf)))

	r := msg.QueryResult{
		ResponseDataset: &msg.DataSet{
			Categorization: &msg.DataSet_Categorization{
				Value: []msg.Category{
					msg.BLOG_4,
					msg.NEWS_4,
				},
			},
		},
		QueryStatus: &msg.QueryStatus{
			Complete:  true,
			FetchCode: http.StatusOK,
		},
	}

	body, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}

	if _, err = http.Post(srv.URL, "application/json", bytes.NewReader(body)); err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(&r, m) {
		t.Log(cmp.Diff(&r, m))
		t.Error("got unexpected result")
	}
}
