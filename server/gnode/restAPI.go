package gnode

import (
	"net/http"
)

const baseURL = "/api/v1"

func (g *GNode) startRESTAPI() {
	logf.info("Start REST API server on port %s\n", config.restPort)
	go func() {
		http.HandleFunc(baseURL+"/health", g.health)
		http.HandleFunc(baseURL+"/ready", g.checkReady)
		http.ListenAndServe(":"+config.restPort, nil)
	}()
}

func (g *GNode) checkReady(resp http.ResponseWriter, req *http.Request) {
	if g.grpcReady {
		resp.WriteHeader(200)
	} else {
		logf.debug("execute /ready: return not reday")
		resp.WriteHeader(400)
	}
}

func (g *GNode) health(resp http.ResponseWriter, req *http.Request) {
	if g.healthy && g.ready {
		resp.WriteHeader(200)
	} else {
		logf.debug("execute /health: return not healthy")
		resp.WriteHeader(400)
	}
}
