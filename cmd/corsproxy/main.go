package main

import (
	"fmt"
	"log"
	"net/http"
	"slices"
	"strings"

	"github.com/rs/cors"
	"github.com/zachcheung/eflag"

	"github.com/zachcheung/corsproxy"
)

var (
	allowedOrigins            eflag.StringList
	allowedMethods            eflag.StringList
	allowedHeaders            eflag.StringList
	exposedHeaders            eflag.StringList
	maxAge                    int
	allowCredentials          bool
	allowPrivateNetwork       bool
	passthrough               bool
	successStatus             int
	debug                     bool
	allowedTargets            eflag.StringList
	allowPrivateNetworkTarget bool
	addr                      string
	defaultAllowedMethods     = []string{http.MethodGet, http.MethodHead, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete}
	defaultAllowedHeaders     = []string{"accept", "content-type", "x-requested-with", "authorization"}
	normalizedAllowedTargets  []string
)

func main() {
	eflag.Var(&allowedOrigins, "allowedOrigins", "", "a list of origins a cross-domain request can be executed from", "")
	eflag.Var(&allowedMethods, "allowedMethods", strings.Join(defaultAllowedMethods, ", "), "a list of methods the client is allowed to use with cross-domain requests", "")
	eflag.Var(&allowedHeaders, "allowedHeaders", strings.Join(defaultAllowedHeaders, ", "), "a list of headers the client is allowed to use with cross-domain requests", "")
	eflag.Var(&exposedHeaders, "exposedHeaders", "", "indicates which headers are safe to expose to the API of a CORS API specification", "")
	eflag.Var(&maxAge, "maxAge", 0, "indicates how long (in seconds) the results of a preflight request can be cached", "")
	eflag.Var(&allowCredentials, "allowCredentials", false, "indicates whether the request can include user credentials like cookies, HTTP authentication or client side SSL certificates", "")
	eflag.Var(&allowPrivateNetwork, "allowPrivateNetwork", false, "indicates whether to accept cross-origin requests over a private network", "")
	eflag.Var(&passthrough, "passthrough", false, "instructs preflight to let other potential next handlers to process the OPTIONS method", "")
	eflag.Var(&successStatus, "successStatus", 204, "provides a status code to use for successful OPTIONS requests", "")
	eflag.Var(&debug, "debug", false, "adds additional output to debug server side CORS issues", "")
	eflag.Var(&allowedTargets, "allowedTargets", "", "AllowedTargets is a list of targets a cross-domain request can reach", "")
	eflag.Var(&allowPrivateNetworkTarget, "allowPrivateNetworkTarget", false, "indicates whether to accept private network targets", "")
	eflag.Var(&addr, "addr", ":8000", "bind address", "")
	eflag.Parse()

	opt := cors.Options{
		MaxAge:               maxAge,
		AllowCredentials:     allowCredentials,
		AllowPrivateNetwork:  allowPrivateNetwork,
		OptionsPassthrough:   passthrough,
		OptionsSuccessStatus: successStatus,
		Debug:                debug,
	}

	if v := allowedOrigins.Value(); len(v) > 0 {
		opt.AllowedOrigins = v
	}
	if v := allowedMethods.Value(); len(v) > 0 {
		opt.AllowedMethods = v
	}
	if v := allowedHeaders.Value(); len(v) > 0 {
		opt.AllowedHeaders = v
	}
	if v := exposedHeaders.Value(); len(v) > 0 {
		opt.ExposedHeaders = v
	}
	if v := allowedTargets.Value(); len(v) > 0 {
		for _, target := range v {
			target, err := corsproxy.StripURLQuery(target)
			if err != nil {
				panic(fmt.Sprintf("Invalid target %s: %v", target, err))
			}
			normalizedAllowedTargets = append(normalizedAllowedTargets, target)
		}
	}

	if debug {
		log.Print("[DEBUG] Debug mode")
		log.Printf("[DEBUG] Options:\nallowedOrigins: %v\nallowedMethods: %v\nallowedHeaders: %v\nexposedHeaders: %v\nmaxAge: %v\nallowCredentials: %v\nallowPrivateNetwork: %v\npassthrough: %v\nsuccessStatus: %v\nallowedTargets: %v\nnormalizedAllowedTargets: %v\nallowPrivateNetworkTarget: %v", allowedOrigins.Value(), allowedMethods.Value(), allowedHeaders.Value(), exposedHeaders.Value(), maxAge, allowCredentials, allowPrivateNetwork, passthrough, successStatus, allowedTargets.Value(), normalizedAllowedTargets, allowPrivateNetworkTarget)
	}

	cp := corsproxy.New(corsproxy.Options{
		Options:                   opt,
		AllowedTargets:            allowedTargets.Value(),
		AllowPrivateNetworkTarget: allowPrivateNetworkTarget,
	})

	if (len(allowedTargets.Value()) == 0 || slices.Contains(allowedTargets.Value(), "*")) && allowPrivateNetworkTarget {
		log.Print("[WARN] Private network targets have been allowed without any configured allowedTargets rule!")
	}

	log.Printf("Listen %s", addr)
	log.Fatal(http.ListenAndServe(addr, cp.Handler()))
}
