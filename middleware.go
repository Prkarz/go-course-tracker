package main

import "net/http"

func JWT_Middleware(net http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		net.
			next(w, r)
	}
}
