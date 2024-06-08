package codex

import (
	"encoding/json"
	g "github.com/ananrafs/grimoire"
	"io/ioutil"
	"sync"
)

type JsonFileCodex struct {
	filename string
	routes   []g.Route
	lock     sync.Mutex
}

func NewJsonCodex(filename string) *JsonFileCodex {
	return &JsonFileCodex{filename: filename}
}

func (jc *JsonFileCodex) Init() error {
	data, err := ioutil.ReadFile(jc.filename)
	if err != nil {
		return err
	}

	jc.lock.Lock()
	defer jc.lock.Unlock()

	config := struct {
		Routes []g.Route `json:"routes"`
	}{}

	if err := json.Unmarshal(data, &config); err != nil {
		return err
	}

	jc.routes = config.Routes

	return nil
}

func (jc *JsonFileCodex) GetAllRoute() []g.Route {
	jc.lock.Lock()
	defer jc.lock.Unlock()
	return jc.routes
}

func (jc *JsonFileCodex) GetChannel() <-chan struct{} {
	signal := make(chan struct{})
	go watchChanges(jc.filename, func() {
		jc.Init()
		signal <- struct{}{}
	})
	return signal
}
