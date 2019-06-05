package main

import (
	"net/http"
)

func SetCorsHeader(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	h := w.Header()
	cors := &ServiceConf.CorsConf

	if cors.AllowOrigin != "" {
		h.Set("Access-Control-Allow-Origin", cors.AllowOrigin)
	}
	if cors.AllowHeaders != "" {
		h.Set("Access-Control-Allow-Headers", cors.AllowHeaders)
	}
	if cors.AllowMethods != "" {
		h.Set("Access-Control-Allow-Methods", cors.AllowMethods)
	}
	if cors.ExposeHeaders != "" {
		h.Set("Access-Control-Expose-Headers", cors.ExposeHeaders)
	}
	if cors.AllowCredentials != "" {
		h.Set("Access-Control-Allow-Credentials", cors.AllowCredentials)
	}

	next(w, r)
}
