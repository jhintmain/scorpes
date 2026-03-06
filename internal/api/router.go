package api

import "net/http"

type Middleware func(http.Handler) http.Handler

type Router struct {
	mux         *http.ServeMux
	middlewares []Middleware
	prefix      string
}

func NewRouter() *Router {
	return &Router{
		mux:         http.NewServeMux(),
		middlewares: make([]Middleware, 0),
		prefix:      "",
	}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

func (r *Router) Use(middleware Middleware) {
	r.middlewares = append(r.middlewares, middleware)
}

func (r *Router) Group(prefix string, fn func(r *Router)) {
	subRouter := &Router{
		mux:         r.mux,
		prefix:      r.prefix + prefix,
		middlewares: append([]Middleware{}, r.middlewares...),
	}

	fn(subRouter)
}

func (r *Router) Handle(method, path string, handler http.Handler) {
	fullPath := method + " " + r.prefix + path

	var finalHandler http.Handler = handler

	for _, m := range r.middlewares {
		finalHandler = m(finalHandler)
	}

	r.mux.Handle(fullPath, finalHandler)
}

func (r *Router) GET(path string, handler http.HandlerFunc) {
	r.Handle(http.MethodGet, path, handler)
}

func (r *Router) POST(path string, handler http.HandlerFunc) {
	r.Handle(http.MethodPost, path, handler)
}

func (r *Router) DELETE(path string, handler http.HandlerFunc) {
	r.Handle(http.MethodDelete, path, handler)
}

func (r *Router) PUT(path string, handler http.HandlerFunc) {
	r.Handle(http.MethodPut, path, handler)
}
