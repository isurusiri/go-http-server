// this is used to demonstrate how to minimize
// duplication in http handlers
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	// go get github.com/justinas/alice
	// alice is a package that helps to chain
	// handlers. This allows to create a
	// common list of handlers and reuse
	// those handlers in each route.
	"github.com/justinas/alice"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome!")
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "About!")
}

// request logging is a common function that
// can be reuse in each route.
// ServeHTTP of http.Handler interface initiates
// request to given http handler. We create a
// function that does our resuable work
// and call chanining handlers.
// this method is returning a handler which can be
// used to register routes and handlers.
func loggingHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		next.ServeHTTP(w, r)
		t2 := time.Now()
		log.Printf("[%s] %q %v\n", r.Method, r.URL.String(), t2.Sub(t1))
	}

	return http.HandlerFunc(fn)
}

// a function to log internal server errors.
// a differed function is used to recover in such
// situations.
func recoverHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("error: %+v", err)
				http.Error(w, http.StatusText(500), 500)
			}
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func main() {
	commonHandlers := alice.New(loggingHandler, recoverHandler)
	http.Handle("/", commonHandlers.ThenFunc(indexHandler))
	http.Handle("/about", commonHandlers.ThenFunc(aboutHandler))
	http.ListenAndServe(":8085", nil)
}
