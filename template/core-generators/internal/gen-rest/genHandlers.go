package gen_rest

import (
	"core-generators/internal/utils"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

type ControllerHandler struct {
	Method    string
	EndPoint  string
	F         string
	AuthRoles []string
}

func GenHandlers(serviceName string, paths map[string]EndPoint) error {
	filePath := fmt.Sprintf("../core-%s/internal/api/rest", serviceName)

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(path.Join(wd, filePath, "handlers"), 0777); err != nil {
		return err
	}

	if !utils.FileExists(path.Join(wd, filePath, "handlers", "handlers.go")) {
		if err := GenSpareHandlersService(path.Join(wd, filePath, "handlers")); err != nil {
			return err
		}
	}

	files, err := filepath.Glob(path.Join(wd, filePath+"/handlers", "*_handler.go"))
	if err != nil {
		return err
	}
	for _, f := range files {
		if err := os.Remove(f); err != nil {
			return err
		}
	}

	filesToFormat := make([]string, 0)
	controllerHandlers := make([]ControllerHandler, 0)

	for ep, p := range paths {
		for mType, mm := range p {
			var m Method
			mp, ok := mm.(EndPoint)
			if ok {
				marsh, err := yaml.Marshal(mp)
				if err != nil {
					return err
				}
				if err := yaml.Unmarshal(marsh, &m); err != nil {
					return err
				}
			} else {
				continue
			}
			controllerHandlers = append(controllerHandlers, ControllerHandler{
				Method:    strings.ToUpper(mType[:1]) + mType[1:],
				EndPoint:  ep,
				F:         m.OperationId,
				AuthRoles: m.AuthRoles,
			})
			procFilePath := path.Join(filePath, "handlers", fmt.Sprintf("%s.go", m.OperationId))
			if !utils.FileExists(procFilePath) {
				if err := genProcessor(filePath, m); err != nil {
					return nil
				}
				filesToFormat = append(filesToFormat, procFilePath)
			}

			if err := genHandler(serviceName, filePath, m); err != nil {
				return nil
			}
			filesToFormat = append(filesToFormat, path.Join(filePath, "handlers", fmt.Sprintf("%s_handler.go", m.OperationId)))
		}
	}

	if err := genController(serviceName, filePath, controllerHandlers); err != nil {
		return err
	}

	filesToFormat = append(filesToFormat, path.Join(filePath, "endpoints.go"))

	for _, f := range filesToFormat {
		utils.FormatFile(f)
	}

	return nil
}

func GenSpareHandlersService(outPath string) error {
	var outData string
	outData = "package handlers\n\n"

	outData += `type Handlers struct {
}

func Create() *Handlers {
	return &Handlers{}
}
`
	fileName := "handlers.go"
	if err := os.WriteFile(path.Join(outPath, fileName), []byte(outData), 0666); err != nil {
		return err
	}

	return nil
}

func genProcessor(outPath string, m Method) error {
	var outData string
	outData = "package handlers\n\n"

	imports := `import (
	"github.com/gofiber/fiber/v2"
	"github.com/{{.orgName}}/{{.pkgRepoName}}/errors"
)

`
	outData += imports

	hasRequestBody := false
	hasResponseBody := false
	if m.RequestBody.Content.Json.Schema != nil {
		hasRequestBody = true
	}

	_, ok := m.Responses["200"]
	if ok {
		if m.Responses["200"].Content.Json.Schema != nil {
			hasResponseBody = true
		}
	}

	capOperationId := strings.ToUpper(m.OperationId[:1]) + m.OperationId[1:]

	var reqParam string
	if hasRequestBody {
		reqParam += fmt.Sprintf(", req *%sRequest", capOperationId)
	}

	outData += fmt.Sprintf(`func (h *Handlers) %sValidate(c *fiber.Ctx%s) error {

	return nil
}
`, capOperationId, reqParam)

	var respParamA, respParamB string
	if hasResponseBody {
		respParamA = fmt.Sprintf("(*%sResponse, ", capOperationId)
		respParamB = ")"
	}

	outData += fmt.Sprintf(`func (h *Handlers) %sProcess(c *fiber.Ctx%s) %serror%s {

`, capOperationId, reqParam, respParamA, respParamB)

	if hasResponseBody {
		outData += fmt.Sprintf("return &%sResponse{}, nil\n", capOperationId)
	} else {
		outData += "return nil\n"
	}

	outData += "}"

	fileName := fmt.Sprintf("%s.go", m.OperationId)
	if err := os.WriteFile(path.Join(outPath, "handlers", fileName), []byte(outData), 0666); err != nil {
		return err
	}

	return nil
}

func genHandler(serviceName, outPath string, m Method) error {
	var outData string

	outData += utils.GenModelHeader

	outData += "package handlers\n\n"

	hasRequestBody := false
	hasResponseBody := false
	if m.RequestBody.Content.Json.Schema != nil {
		hasRequestBody = true
	}

	_, ok := m.Responses["200"]
	if ok {
		if m.Responses["200"].Content.Json.Schema != nil {
			hasResponseBody = true
		}
	}

	capOperationId := strings.ToUpper(m.OperationId[:1]) + m.OperationId[1:]

	var oReqData, eReqData, oRespData, eRespData string
	refImports := make([]PkgImport, 0)

	if hasRequestBody {
		oData, eData, pData := GenType(capOperationId+"Request", *m.RequestBody.Content.Json.Schema, true, true, false)
		oReqData = oData
		eReqData = eData
		refImports = append(refImports, pData...)
	}

	if hasResponseBody {
		oData, eData, pData := GenType(capOperationId+"Response", *m.Responses["200"].Content.Json.Schema, true, true, false)
		oRespData = oData
		eRespData = eData
		refImports = append(refImports, pData...)
	}

	imports := fmt.Sprintf(`import (
	restModels "core-%s/internal/models/rest"
	"github.com/gofiber/fiber/v2"
	"github.com/{{.orgName}}/{{.pkgRepoName}}/errors"
`, serviceName)
	outData += imports

	for _, i := range refImports {
		outData += fmt.Sprintf("\t%s \"%s\"\n", i.Name, i.Path)
	}

	outData += ")\n"
	outData += oReqData
	outData += eReqData
	outData += oRespData
	outData += eRespData

	outData += fmt.Sprintf("func (h *Handlers) %s(c *fiber.Ctx) error {\n", capOperationId)

	if hasRequestBody {
		outData += fmt.Sprintf("req := new(%sRequest)\n", capOperationId)
		outData += `	err := c.BodyParser(req)
	if err != nil {
		return errors.ParseRequest.WithInfo(err.Error())
	}

`
	}

	var reqParam string
	if hasRequestBody {
		reqParam += ", req"
	}

	outData += fmt.Sprintf(`	if err := h.%sValidate(c%s); err != nil {
		return err
	}

`, capOperationId, reqParam)

	hasAnyBody := hasRequestBody || hasResponseBody

	if hasResponseBody {
		outData += fmt.Sprintf(`res, err := h.%sProcess(c%s)
	if err != nil {
		return err
	}
	c.JSON(res)

`, capOperationId, reqParam)
	} else {
		hb := ""
		if !hasAnyBody {
			hb = ":"
		}
		outData += fmt.Sprintf(`err %s= h.%sProcess(c%s)
	if err != nil {
		return err
	}

`, hb, capOperationId, reqParam)
	}

	outData += "return nil\n"
	outData += "}"

	fileName := fmt.Sprintf("%s_handler.go", m.OperationId)
	if err := os.WriteFile(path.Join(outPath, "handlers", fileName), []byte(outData), 0666); err != nil {
		return err
	}

	return nil
}

func genController(serviceName string, outPath string, hs []ControllerHandler) error {
	var outData string
	outData = "package rest\n\n"

	outData += fmt.Sprintf(`import (
	"core-%s/internal/api/rest/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/{{.orgName}}/{{.pkgRepoName}}/models"
	"github.com/{{.orgName}}/{{.pkgRepoName}}/rest/middleware"
)

func configureEndPoints(c *fiber.App, h *handlers.Handlers) {
`, serviceName)

	sort.Slice(hs, func(i, j int) bool {
		return hs[i].EndPoint < hs[j].EndPoint
	})

	for _, ch := range hs {
		authMiddleware := " middleware.Authorize("
		if ch.AuthRoles != nil && len(ch.AuthRoles) > 0 {
			roles := make([]string, 0)
			for _, r := range ch.AuthRoles {
				roles = append(roles, fmt.Sprintf("models.UserRole%s", utils.UpperCaseFirstLetter(string(r))))
			}
			authMiddleware += strings.Join(roles, ", ")
		}
		authMiddleware += "), "
		outData += fmt.Sprintf("c.%s(\"%s\",%s h.%s)\n", ch.Method, ch.EndPoint, authMiddleware, strings.ToUpper(ch.F[:1])+ch.F[1:])
	}

	outData += "}"
	fileName := "endpoints.go"
	if err := os.WriteFile(path.Join(outPath, fileName), []byte(outData), 0666); err != nil {
		return err
	}

	return nil
}
