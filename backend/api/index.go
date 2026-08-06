// Package api is the Vercel serverless entry point. Vercel's Go runtime
// looks for an exported `Handler` symbol of type `http.Handler` inside
// `api/index.go` (or any `api/*.go` file) and invokes it for every
// incoming HTTP request.
//
// The actual Echo application is initialized in `internal/server` and
// built once via `sync.Once` so the database connection pool and
// registered routes are reused across serverless invocations.
package api

import (
	"net/http"

	"spotsync-backend/internal/server"
)

// Handler is what Vercel's runtime invokes. Echo's `*Echo` implements
// `http.Handler`, so we expose it directly. We call `server.Get()`
// here so the very first cold start performs initialization.
var Handler http.Handler = server.Get()
