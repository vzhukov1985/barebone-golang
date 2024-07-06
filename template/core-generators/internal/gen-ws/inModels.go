package gen_ws

import (
	gen_rest "core-generators/internal/gen-rest"
)

type AsyncApi struct {
	Channels   map[string]Channel `yaml:"channels"`
	Components Components         `yaml:"components"`
}

type Channel struct {
	Publish   *ChannelOp `yaml:"publish"`
	Subscribe *ChannelOp `yaml:"subscribe"`
	Roles     []string   `yaml:"x-roles"`
}

type ChannelOp struct {
	Message *Message `yaml:"message"`
}

type Components struct {
	Messages map[string]Message        `yaml:"messages"`
	Schemas  map[string]gen_rest.Model `yaml:"schemas"`
}

type Message struct {
	Payload gen_rest.Model `yaml:"payload"`
}
