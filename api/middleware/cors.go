package middleware

import (
	"github.com/rs/cors"
	"net/http"
)

func CorsEverywhere(mux http.Handler) http.Handler {
	return cors.Default().Handler(mux)
}
