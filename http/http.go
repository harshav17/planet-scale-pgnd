package http

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	planetscale "github.com/harshav17/planet_scale"
)

type ContentType string

const (
	ContentTypeJson ContentType = "application/json"
)

func ReceiveJson(w http.ResponseWriter, r *http.Request, thing any) error {
	err := MustBeContentType(r, ContentTypeJson)
	if err != nil {
		return planetscale.Errorf(planetscale.EINVALID, "invalid content-type")
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, thing)
	if err != nil {
		return err
	}

	return nil
}

func RespondJson(w http.ResponseWriter, r *http.Request, statusCode int, thing any) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-type", "application/json")
	err := json.NewEncoder(w).Encode(thing)
	if err != nil {
		Error(w, r, err)
	}
}

// Error prints & optionally logs an error message.
func Error(w http.ResponseWriter, r *http.Request, err error) {
	// Extract error code & message.
	code, message := planetscale.ErrorCode(err), planetscale.ErrorMessage(err)

	// Log & report internal errors.
	if code == planetscale.EINTERNAL {
		planetscale.ReportError(r.Context(), err, r)
		LogError(r, err)
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(ErrorStatusCode(code))
	json.NewEncoder(w).Encode(&ErrorResponse{Error: message})
}

// ErrorResponse represents a JSON structure for error output.
type ErrorResponse struct {
	Error string `json:"error"`
}

// LogError logs an error with the HTTP route information.
func LogError(r *http.Request, err error) {
	slog.Error("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
}

func MustBeContentType(r *http.Request, contentType ContentType) error {
	if !IsContentType(r, contentType) {
		return planetscale.Errorf(planetscale.EINVALID, "invalid content type")
	}
	return nil
}

func IsContentType(r *http.Request, contentType ContentType) bool {
	requestContentType := r.Header.Get("Content-Type")
	switch contentType {
	case ContentTypeJson:
		return requestContentType == "application/json"
	// ...
	default:
		return false
	}
}

// lookup of application error codes to HTTP status codes.
var codes = map[string]int{
	planetscale.ECONFLICT:       http.StatusConflict,
	planetscale.EINVALID:        http.StatusBadRequest,
	planetscale.ENOTFOUND:       http.StatusNotFound,
	planetscale.ENOTIMPLEMENTED: http.StatusNotImplemented,
	planetscale.EUNAUTHORIZED:   http.StatusUnauthorized,
	planetscale.EINTERNAL:       http.StatusInternalServerError,
}

// ErrorStatusCode returns the associated HTTP status code for a WTF error code.
func ErrorStatusCode(code string) int {
	if v, ok := codes[code]; ok {
		return v
	}
	return http.StatusInternalServerError
}
