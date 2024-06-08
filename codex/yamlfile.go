package codex

import (
	"io/ioutil"
	"sync"

	g "github.com/ananrafs/grimoire"
	yaml "gopkg.in/yaml.v2"
)

type YamlCodex struct {
	filename string
	routes   []g.Route
	lock     sync.Mutex
}

func NewYamlCodex(filename string) *YamlCodex {
	return &YamlCodex{filename: filename}
}

func (yc *YamlCodex) Init() error {
	data, err := ioutil.ReadFile(yc.filename)
	if err != nil {
		return err
	}

	yc.lock.Lock()
	defer yc.lock.Unlock()

	config := struct {
		Routes []g.Route `json:"routes"`
	}{}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return err
	}
	yc.routes = config.Routes

	return nil
}

func (yc *YamlCodex) GetAllRoute() []g.Route {
	yc.lock.Lock()
	defer yc.lock.Unlock()
	return yc.routes
}

func (yc *YamlCodex) GetChannel() chan struct{} {
	signal := make(chan struct{})
	go watchChanges(yc.filename, func() {
		yc.Init()
		signal <- struct{}{}
	})
	return signal
}
