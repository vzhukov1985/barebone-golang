package gen_ws

import (
	gen_rest "core-generators/internal/gen-rest"
)

type Handler struct {
	Channel  string
	Name     string
	Roles    []string
	Request  *Payload
	Response *Payload
}

type Api struct {
	Handlers []Handler
	Models   []gen_rest.Model
}

type Payload struct {
	Model             *gen_rest.Model
	CommonModelFolder string
	CommonModelName   string
}
