package grimoire

type Codex interface {
	Init() error
	GetAllRoute() []Route
	GetChannel() (signal chan struct{})
}
