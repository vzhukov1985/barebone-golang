package gen_ws

import (
	gen_rest "core-generators/internal/gen-rest"
	"core-generators/internal/utils"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func GenHandlers(api *Api) error {
	filePath := "../core-websocket/internal/api/ws"

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(path.Join(wd, filePath, "handlers"), 0777); err != nil {
		return err
	}

	files, err := filepath.Glob(path.Join(wd, filePath+"/handlers", "*_handler.go"))
	if err != nil {
		return err
	}
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			return err
		}
		if strings.HasPrefix(string(data), "// Schema generated model - DO NOT EDIT!!!") {
			if err := os.Remove(f); err != nil {
				return err
			}
		}
	}

	if utils.FileExists(path.Join(wd, filePath, "handlers.go")) {
		if err := os.Remove(path.Join(wd, filePath, "handlers.go")); err != nil {
			return err
		}
	}

	if err := GenController(api); err != nil {
		return err
	}

	if err := GenHandlerFiles(api); err != nil {
		return err
	}

	return nil
}

func GenController(api *Api) error {
	var outData string
	outData = `// Schema generated model - DO NOT EDIT!!!

package ws

import (
	"core-websocket/internal/api/ws/handlers"
	"core-websocket/internal/models/ws"
	"github.com/{{.orgName}}/{{.pkgRepoName}}/models"
)

var handlerFuncs map[string]ws.ChannelHandler

func configureHandlers(h *handlers.Handlers) {
	handlerFuncs = map[string]ws.ChannelHandler{
`
	for _, v := range api.Handlers {
		roles := "[]models.UserRole{}"
		if v.Roles != nil && len(v.Roles) > 0 {
			roles = "[]models.UserRole{"
			rArr := make([]string, 0)
			for _, r := range v.Roles {
				rArr = append(rArr, "models.UserRole"+utils.UpperCaseFirstLetter(string(r)))
			}
			roles += strings.Join(rArr, ",") + "}"
		}
		outData += fmt.Sprintf(`"%s": {
			Handler: h.%sHandler,
			Roles:   %s,
		},
`, v.Channel, v.Name, roles)
	}

	outData += `}
}
`
	if err := os.WriteFile("../core-websocket/internal/api/ws/handlers.go", []byte(outData), 0666); err != nil {
		return err
	}
	utils.FormatFile("../core-websocket/internal/api/ws/handlers.go")

	return nil
}

func GenHandlerFiles(api *Api) error {
	filesToFormat := make([]string, 0)
	for _, v := range api.Handlers {
		var outData string
		pkgs := make([]gen_rest.PkgImport, 0)
		var modData string
		if v.Request != nil {
			if v.Request.CommonModelName == "" {
				oData, eData, pData := gen_rest.GenType(v.Name+"Request", *v.Request.Model, true, true, false)
				pkgs = append(pkgs, pData...)
				modData += oData + "\n"
				modData += eData
			} else {
				var pkgPrefix string
				pkgPath := "github.com/{{.orgName}}/{{.pkgRepoName}}/models"
				if v.Request.CommonModelFolder == "" {
					pkgPrefix = "common"
				} else {
					pkgPrefix = v.Request.CommonModelFolder
					pkgPath += "/" + v.Request.CommonModelFolder
				}
				pkgs = append(pkgs, gen_rest.PkgImport{
					Name: fmt.Sprintf("%sModels", pkgPrefix),
					Path: pkgPath,
				})
				modData += fmt.Sprintf("type %sRequest %sModels.%s\n\n", v.Name, pkgPrefix, v.Request.CommonModelName)
			}
		}

		if v.Response != nil {
			if v.Response.CommonModelName == "" {
				oData, eData, pData := gen_rest.GenType(v.Name+"Response", *v.Response.Model, true, true, false)
				pkgs = append(pkgs, pData...)
				modData += oData
				modData += eData
			} else {
				var pkgPrefix string
				pkgPath := "github.com/{{.orgName}}/{{.pkgRepoName}}/models"
				if v.Response.CommonModelFolder == "" {
					pkgPrefix = "common"
				} else {
					pkgPrefix = v.Response.CommonModelFolder
					pkgPath += "/" + v.Response.CommonModelFolder
				}
				pkgs = append(pkgs, gen_rest.PkgImport{
					Name: fmt.Sprintf("%sModels", pkgPrefix),
					Path: pkgPath,
				})
				modData += fmt.Sprintf("type %sResponse %sModels.%s\n\n", v.Name, pkgPrefix, v.Response.CommonModelName)
			}
		}

		outData = `// Schema generated model - DO NOT EDIT!!!

package handlers

import (
	"context"
	"core-websocket/internal/models/ws"
	"github.com/goccy/go-json"
	"github.com/{{.orgName}}/{{.pkgRepoName}}/log"
	"github.com/{{.orgName}}/{{.pkgRepoName}}/models"
`

		for _, p := range pkgs {
			outData += fmt.Sprintf("\t%s \"%s\"\n", p.Name, p.Path)
		}

		outData += `
)

`
		outData += fmt.Sprintf(`const %sChannel = "%s"

`, v.Name, v.Channel)

		outData += modData

		outData = strings.ReplaceAll(outData, "restModels", "wsModels")

		if !strings.Contains(outData, "wsModels") {
			outData = strings.ReplaceAll(outData, "\"core-websocket/internal/models/ws-models\"", "")
		}

		outData += fmt.Sprintf("func (h *Handlers) %sHandler(ctx context.Context, w *ws.Ws, payload []byte) {\n", v.Name)
		if v.Request != nil {
			outData += fmt.Sprintf(`var p %sRequest
	if err := json.Unmarshal(payload, &p); err != nil {
		log.Errorw("Failed to unmarshal subscription payload", "error", err, "channel", %sChannel)
		w.SendError(%sChannel, err)
		return
	}
`, v.Name, v.Name, v.Name)
		}

		outData += fmt.Sprintf(`log.Debugw("User subscribed to channel", "clientKey", w.C.Headers("Sec-Websocket-Key"), "channel", %sChannel)

`, v.Name)
		var p string
		if v.Request != nil {
			p = `, &p`
		}
		outData += fmt.Sprintf("h.%s(ctx, w", v.Name) + p + ")\n\n"

		outData += fmt.Sprintf(`log.Debugw("User unsubscribed from channel", "clientKey", w.C.Headers("Sec-Websocket-Key"), "channel", %sChannel)
`, v.Name)

		outData += "}"

		outFileHandler := path.Join("../core-websocket/internal/api/ws/handlers", utils.LowerCaseFirstLetter(v.Name)+"_handler.go")
		if err := os.WriteFile(outFileHandler, []byte(outData), 0666); err != nil {
			return err
		}

		filesToFormat = append(filesToFormat, outFileHandler)

		if !utils.FileExists(path.Join("../core-websocket/internal/api/ws/handlers", utils.LowerCaseFirstLetter(v.Name)+".go")) {
			outData = `package handlers

import (
	"context"
	"core-websocket/internal/models/ws"
)
`
			p = ""
			if v.Request != nil {
				p = fmt.Sprintf(", r *%sRequest", v.Name)
			}
			outData += fmt.Sprintf("func (h *Handlers) %s(ctx context.Context, c *ws.Ws%s) {\n\n", v.Name, p)
			outData += "<-ctx.Done()\n\n"
			outData += "}"
			outFileProcess := path.Join("../core-websocket/internal/api/ws/handlers", utils.LowerCaseFirstLetter(v.Name)+".go")
			if err := os.WriteFile(outFileProcess, []byte(outData), 0666); err != nil {
				return err
			}
			filesToFormat = append(filesToFormat, outFileProcess)
		}
	}

	for _, v := range filesToFormat {
		utils.FormatFile(v)
	}

	return nil
}
