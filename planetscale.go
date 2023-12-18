package planetscale

import (
	"context"
	"net/http"
)

var ReportError = func(ctx context.Context, err error, args *http.Request) {}

var ReportPanic = func(err interface{}) {}
