package worker

import (
	"dev-helper/internal/services"
	"dev-helper/internal/shell"
)

func Start() {

	for _, v := range services.Slice {
		go shell.ExecuteBackendService(v)
	}
}

func Stop() {

}
