package gen_rest

type MethodContent struct {
	Content struct {
		Json struct {
			Schema *Model `yaml:"schema"`
		} `yaml:"application/json"`
	} `yaml:"content"`
}

type Method struct {
	AuthRoles   []string                 `yaml:"x-roles"`
	OperationId string                   `yaml:"operationId"`
	Summary     string                   `yaml:"summary"`
	RequestBody MethodContent            `yaml:"requestBody"`
	Responses   map[string]MethodContent `yaml:"responses"`
}

type EndPoint map[string]interface{}

type Model struct {
	Title      string            `yaml:"title"`
	Type       string            `yaml:"type"`
	Enum       []string          `yaml:"enum"`
	Items      *Model            `yaml:"items"`
	Properties map[string]Model  `yaml:"properties"`
	Required   []string          `yaml:"required"`
	Ref        string            `yaml:"$ref"`
	ErrorDescs map[string]string `yaml:"x-error-descriptions"`
	AnyOf      []interface{}     `yaml:"anyOf"`
}

type Schemas map[string]Model

type Rest struct {
	Paths      map[string]EndPoint `yaml:"paths"`
	Components map[string]Schemas  `yaml:"components"`
}

type ModelGenerator struct {
	EnumData string
}
