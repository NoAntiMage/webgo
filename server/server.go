package server

import (
	"errors"
	"fmt"
	account "goweb/app/account/router"
	info "goweb/app/info/router"
	"goweb/common"
	"goweb/common/logx"
	"goweb/engine"
	"net/http"
	"runtime/debug"
	"time"
)

func Start() {
	common.GlobalInit()
	defer func() {
		if p := recover(); p != nil {
			panicResult := fmt.Sprintln(p)
			stackBytes := debug.Stack()
			err := errors.New(string(panicResult) + string(stackBytes))
			logx.Loggerx.Error(err)
		}
	}()

	//CUSTOM: Add router here.
	initRouters := []engine.GroupRouterHandler{
		info.InitInfoRouter,
		account.InitAccountRouter,
	}

	rootRouter := engine.InitEngine(initRouters...)

	//TODO CUSTOM http config
	s := &http.Server{
		Addr:           ":8080",
		Handler:        rootRouter,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	logx.Loggerx.Info("server starting...")
	s.ListenAndServe()
}
