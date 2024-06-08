package main

import (
	"github.com/ananrafs/grimoire"
	"github.com/ananrafs/grimoire/caster"
	"github.com/ananrafs/grimoire/codex"
	"net/http"
	"time"
)

func main() {
	var server grimoire.Server

	//exampleCodex, onQuit := exampleHttpCodex()
	//defer onQuit()

	server = grimoire.NewServer(
		//exampleHttpCodex,
		exampleJsonCodex(),
		//exampleYamlCodex(),
		exampleDummyHandler(),
		grimoire.WithThrottle(10*time.Second))
	quit := server.Serve("8188")
	defer quit()

}

func exampleJsonCodex() grimoire.Codex {
	return codex.NewJsonCodex("routes.json")
}

func exampleYamlCodex() grimoire.Codex {
	return codex.NewJsonCodex("routes.yaml")
}

func exampleHttpCodex() (grimoire.Codex, func()) {
	return codex.NewHttpCodex(8190)
}

func exampleDummyHandler() grimoire.Caster[http.Request] {
	return &caster.DummyCaster{}
}
