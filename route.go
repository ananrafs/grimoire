package grimoire

type Route struct {
	URL  string `json:"url" yaml:"url"`
	Meta []Meta `json:"meta" yaml:"meta"`
}
