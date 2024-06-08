package grimoire

type Meta struct {
	Request  map[string]interface{} `json:"request"`
	Response map[string]interface{} `json:"response"`
	Type     string                 `json:"__type"`
	Param    map[string]interface{} `json:"param"`
}

type Caster[T any] interface {
	Cast(meta []Meta, request T) (map[string]interface{}, error)
}
