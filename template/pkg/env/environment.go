package env

type Environment string

const (
	Development Environment = "development"
	Stage       Environment = "stage"
	Production  Environment = "production"
)

var EnvironmentSlice = []Environment{
	Development,
	Stage,
	Production,
}

func IsDevelopment() bool {
	return environment == Development
}

func IsStage() bool {
	return environment == Stage
}

func IsProduction() bool {
	return environment == Production
}

func Current() Environment {
	return environment
}
