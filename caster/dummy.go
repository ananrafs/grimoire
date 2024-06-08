package caster

import (
	"github.com/ananrafs/grimoire"
	"net/http"
)

type DummyCaster struct{}

func (dc *DummyCaster) Cast(meta []grimoire.Meta, req http.Request) (map[string]interface{}, error) {
	// Dummy implementation,
	return map[string]interface{}{}, nil
}

// TODO: add more caster
