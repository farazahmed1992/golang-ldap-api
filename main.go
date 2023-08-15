// Package main ldap active-directory API.
//
// the purpose of this application is to provide an application
// that is using plain go code to perform ldap actions over an API
//
// This documents all the possible endpoints, parameters, security and usage.
//
// Terms Of Service:
//
// there are no TOS at this moment, use at your own risk we take no responsibility
//
//	Schemes: https
//	Host: Ldap@personaldomain.com
//	BasePath: /
//	Version: 1.0.7
//	License: MIT http://opensource.org/licenses/MIT
//	Contact: Faraz Ahmed Khan<farazahmed1992@gmail.com>
//
//	Consumes:
//	- application/json
//
//
//	Produces:
//	- application/json
//
//
//	Security:
//	- api_key:
//
//	SecurityDefinitions:
//	api_key:
//	     type: apiKey
//	     name: X-KEY
//	     in: header
//
// swagger:meta
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	//"github.com/hashicorp/vault/api"

	myhanglers "go-ldap-api/handlers"
	"go-ldap-api/util"

	"github.com/go-openapi/runtime/middleware"
	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
)

var Config, _ = util.LoadConfig(".")

// Find is to check if the key is matched in the slice of APIKeys.
func Find(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

var ApiKeys []string

func apiKeyMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("X-KEY")
		if !Find(ApiKeys, key) {
			w.WriteHeader(401)
			//log.Println(ApiKeys)
			err := json.NewEncoder(w).Encode(map[string]string{"status": "failed", "data": "Auth failed"})
			if err != nil {
				return
			}
		} else {
			// Call the next handler, which can be another middleware in the chain, or the final handler.
			next.ServeHTTP(w, r)
		}

	})
}

func init() {

	h := Config.APITOKEN
	keys := strings.Split(h, ",")

	for _, item := range keys {
		ApiKeys = append(
			ApiKeys,
			item,
		)
	}

}

func main() {
	l := hclog.Default()
	config, err := util.LoadConfig(".")
	if err != nil {
		l.Info("cannot load config:", err)
	}

	// create a new serve mux and register the handlers
	sm := mux.NewRouter()
	protectedRoutes := sm.PathPrefix("/api").Subrouter()
	protectedRoutes.Use(apiKeyMiddleware)

	// CORS
	ch := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"*"}))

	// DOCs
	ops := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(ops, nil)

	sm.Handle("/docs", sh)
	sm.Handle("/swagger.yaml", http.FileServer(http.Dir("./"))).Methods("GET")

	// Handlers
	sm.HandleFunc("/", myhanglers.NewHome)
	protectedRoutes.HandleFunc("/add/usertogroup", myhanglers.AddToGroup).Methods("POST")
	protectedRoutes.HandleFunc("/check/usertogroup", myhanglers.CheckToGroup).Methods("POST")
	protectedRoutes.HandleFunc("/add/bulkusertogroup", myhanglers.AddBulkToGroup).Methods("POST")
	protectedRoutes.HandleFunc("/remove/userfromgroup", myhanglers.RemoveFromGroup).Methods("POST")
	protectedRoutes.HandleFunc("/add/user", myhanglers.AddUserToOU).Methods("POST")
	protectedRoutes.HandleFunc("/remove/bulkusersfromgroup", myhanglers.RemoveBulkFromGroup).Methods("POST")
	protectedRoutes.HandleFunc("/add/group", myhanglers.AddGroup).Methods("POST")
	protectedRoutes.HandleFunc("/remove/group", myhanglers.RemoveGroup).Methods("POST")

	// create a new server
	s := http.Server{
		Addr:         config.BindAddress,                               // configure the bind address
		Handler:      ch(sm),                                           // set the default handler
		ErrorLog:     l.StandardLogger(&hclog.StandardLoggerOptions{}), // set the logger for the server
		ReadTimeout:  160 * time.Second,                                // max time to read request from the client
		WriteTimeout: 240 * time.Second,                                // max time to write response to the client
		IdleTimeout:  240 * time.Second,                                // max time for connections using TCP Keep-Alive
	}

	// start the server
	go func() {
		l.Info(fmt.Sprintf("Starting stage server on %s", s.Addr))
		log.Printf("Start logging")
		err := s.ListenAndServe()
		if err != nil {
			l.Error("Error starting server", "error", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interrupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	log.Println("Got signal:", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	err = s.Shutdown(ctx)
	if err != nil {
		return
	}
	defer cancel()
}
