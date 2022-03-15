package server

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"go.uber.org/zap"

	"github.com/xdorro/golang-grpc-base-project/pkg/crypto"
)

type MiddlewarePayload struct {
	Data string `json:"data"`
}

// ResponseWriter wrap http response to CryptoService response body
type ResponseWriter struct {
	http.ResponseWriter
	buf bytes.Buffer
}

func (rw *ResponseWriter) Write(p []byte) (int, error) {
	return rw.buf.Write(p)
}

func (rw *ResponseWriter) String() string {
	return rw.buf.String()
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

		// middleware that decrypts the request body.
		rw := &ResponseWriter{ResponseWriter: w}
		parseRequestBody(s.log, rw, r)

		// Wrap response write by ResponseWriter
		s.httpServer.ServeHTTP(rw, r)

		encryptResponseBody(s.log, rw, r)
		return
	})
}

func writeResponse(w http.ResponseWriter, status int, message interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(message)
}

func parseRequestBody(log *zap.Logger, w *ResponseWriter, r *http.Request) {
	payload := &MiddlewarePayload{}
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		log.Error("json decode error", zap.Error(err))
		return
	}

	requestKeyRaw := r.Header.Get("requestKey")
	if requestKeyRaw == "" {
		log.Error("requestKey is empty")
		writeResponse(w, http.StatusBadRequest, "requestKey is empty")
		return
	}

	// requestKeyRaw := "9b1deb4d-3b7d-4bad-9bdd-2b0d7b3dcb6d"
	decodedQuery, err := crypto.DecryptAES(payload.Data, requestKeyRaw)
	if err != nil {
		log.Panic("decrypt error", zap.Error(err))
	}

	log.Info("request", zap.String("data", decodedQuery))
	r.Body = ioutil.NopCloser(strings.NewReader(decodedQuery))
}

func encryptResponseBody(log *zap.Logger, rw *ResponseWriter, r *http.Request) {
	rw.ResponseWriter.Header().Set("Content-Type", "application/json")

	requestKeyRaw := r.Header.Get("requestKey")
	if requestKeyRaw == "" {
		log.Error("requestKey is empty")
		writeResponse(rw.ResponseWriter, http.StatusBadRequest, "requestKey is empty")
		return
	}

	log.Info("response", zap.String("data", rw.String()))

	encodedQuery, err := crypto.EncryptAES(rw.String(), requestKeyRaw)
	if err != nil {
		log.Panic("decrypt error", zap.Error(err))
	}

	log.Info("encodedQuery", zap.String("data", encodedQuery))

	payload := &MiddlewarePayload{
		Data: encodedQuery,
	}
	writeResponse(rw.ResponseWriter, http.StatusOK, payload)
}
