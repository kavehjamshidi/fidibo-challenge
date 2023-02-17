package fidibosearch

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
)

type httpMock struct {
	statusCode int
	response   interface{}
}

func (m *httpMock) mockHandler(w http.ResponseWriter, r *http.Request) {
	var resp []byte

	rt := reflect.TypeOf(m.response)
	if rt.Kind() == reflect.String {
		resp = []byte(m.response.(string))
	} else if rt.Kind() == reflect.Struct || rt.Kind() == reflect.Ptr {
		resp, _ = json.Marshal(m.response)
	} else {
		resp = []byte("{}")
	}

	w.WriteHeader(m.statusCode)
	w.Write(resp)
}

func newMockServer(pattern string, statusCode int, response interface{}) *httptest.Server {
	c := &httpMock{statusCode, response}

	handler := http.NewServeMux()
	handler.HandleFunc(pattern, c.mockHandler)

	return httptest.NewServer(handler)
}
