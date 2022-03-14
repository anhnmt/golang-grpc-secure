package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

type MiddlewarePayload struct {
	Data string `json:"data"`
}

// httpGrpcRouter is http grpc router.
func (s *Server) httpGrpcRouter() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			s.grpcServer.ServeHTTP(w, r)
			return
		}

		// middleware that adds CORS headers to the response.
		h := w.Header()
		h.Set("Access-Control-Allow-Origin", "http://localhost:3000")
		h.Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			h.Set("Access-Control-Methods", "POST, PUT, PATCH, DELETE")
			h.Set("Access-Control-Allow-Headers", "Access-Control-Allow-Origin,Content-Type")
			h.Set("Access-Control-Max-Age", "86400")
			w.WriteHeader(http.StatusNoContent)
			return
		}

		payload := &MiddlewarePayload{}
		err := json.NewDecoder(r.Body).Decode(payload)
		if err != nil {
			s.log.Error("json decode error", zap.Error(err))
			return
		}

		s.log.Info("request", zap.Any("data", payload))
		r.Body = ioutil.NopCloser(strings.NewReader("ahihi"))

		s.httpServer.ServeHTTP(w, r)

		_, err = w.Write([]byte("ahihi"))
		if err != nil {
			return
		}
		return
	})
}
