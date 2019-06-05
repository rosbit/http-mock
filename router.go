package main

import (
	"github.com/urfave/negroni"
	"github.com/gernest/alien"
	"fmt"
	"net/http"
	"strings"
)

func StartService() error {
	api := negroni.New()
	api.Use(negroni.NewRecovery())
	api.Use(negroni.NewLogger())
	api.Use(negroni.HandlerFunc(SetCorsHeader))

	router := alien.New()

	if ServiceConf.Root != "" && ServiceConf.Alias != "" {
		aliasPath := ServiceConf.Alias
		if !strings.HasSuffix(aliasPath, "/") {
			aliasPath = fmt.Sprintf("%s/", aliasPath)
		}

		alias := fmt.Sprintf("%s*", aliasPath)
		router.Get(alias, http.StripPrefix(aliasPath, http.FileServer(http.Dir(ServiceConf.Root))).ServeHTTP)
	}

	endpoints := CreateEndpoints()
	for _, e := range endpoints {
		uri := e.GetUri()
		handler := e.GetHandler()
		methods := e.GetMethods()
		for _, m := range methods {
			router.AddRoute(m, uri, handler)
		}
	}

	api.UseHandler(router)
	listenParam := fmt.Sprintf(":%d", ServiceConf.Port)
	fmt.Printf("I am listening at %s...\n", listenParam)
	fmt.Printf("%v\n", http.ListenAndServe(listenParam, api))
	return nil
}
