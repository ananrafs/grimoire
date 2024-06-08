package grimoire

type Collector interface {
	Init() error
	GetAllRoute() []Route
	GetChannel() (signal chan struct{})
}
