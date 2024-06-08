package collector

import (
	"io/ioutil"
	"sync"

	g "github.com/ananrafs/grimoire"
	yaml "gopkg.in/yaml.v2"
)

type YamlCollector struct {
	filename string
	routes   []g.Route
	lock     sync.Mutex
}

func NewYamlCollector(filename string) *YamlCollector {
	return &YamlCollector{filename: filename}
}

func (yc *YamlCollector) Init() error {
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

func (yc *YamlCollector) GetAllRoute() []g.Route {
	yc.lock.Lock()
	defer yc.lock.Unlock()
	return yc.routes
}

func (yc *YamlCollector) GetChannel() chan struct{} {
	signal := make(chan struct{})
	go watchChanges(yc.filename, func() {
		yc.Init()
		signal <- struct{}{}
	})
	return signal
}
