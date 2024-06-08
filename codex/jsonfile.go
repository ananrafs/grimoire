package collector

import (
	"encoding/json"
	g "github.com/ananrafs/grimoire"
	"io/ioutil"
	"sync"
)

type JsonCollector struct {
	filename string
	routes   []g.Route
	lock     sync.Mutex
}

func NewJsonCollector(filename string) *JsonCollector {
	return &JsonCollector{filename: filename}
}

func (jc *JsonCollector) Init() error {
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

func (jc *JsonCollector) GetAllRoute() []g.Route {
	jc.lock.Lock()
	defer jc.lock.Unlock()
	return jc.routes
}

func (jc *JsonCollector) GetChannel() chan struct{} {
	signal := make(chan struct{})
	go watchChanges(jc.filename, func() {
		jc.Init()
		signal <- struct{}{}
	})
	return signal
}
