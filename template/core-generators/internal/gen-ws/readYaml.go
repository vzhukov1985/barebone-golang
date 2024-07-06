package gen_ws

import (
	gen_rest "core-generators/internal/gen-rest"
	"core-generators/internal/utils"
	"gopkg.in/yaml.v3"
	"os"
	"sort"
	"strings"
)

func ReadYaml() (*Api, error) {
	data, err := os.ReadFile("../{{.prjName}}-schema/websocket/ws-asyncapi.yaml")
	if err != nil {
		return nil, err
	}

	var f map[string]interface{}
	var m AsyncApi
	if err := yaml.Unmarshal(data, &f); err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, err
	}

	res := Api{Handlers: make([]Handler, 0), Models: make([]gen_rest.Model, 0)}

	for k, v := range m.Channels {
		var req, resp *Payload

		if v.Publish != nil && v.Publish.Message != nil {
			pub := *v.Publish
			req = &Payload{
				Model:             &pub.Message.Payload,
				CommonModelFolder: "",
				CommonModelName:   "",
			}
		}

		if v.Subscribe != nil && v.Subscribe.Message != nil {
			sub := *v.Subscribe
			resp = &Payload{
				Model:             &sub.Message.Payload,
				CommonModelFolder: "",
				CommonModelName:   "",
			}
		}

		/*
			ref := v.Publish.Message.Payload.Ref

			if ref == "" {
				req = &Payload{
					Model:             &v.Publish.Message.Payload,
					CommonModelFolder: "",
					CommonModelName:   "",
				}
			} else {
				msgName := filepath.Base(ref)
				msg, ok := m.Components.Messages[msgName]
				if ok {
					pRef := msg.Payload.Ref
					if pRef == "" {
						req = &Payload{
							Model:             &msg.Payload,
							CommonModelFolder: "",
							CommonModelName:   "",
						}
					} else {
						modName := filepath.Base(pRef)
						if strings.HasPrefix(pRef, "#") {
							model, ok := m.Components.Schemas[modName]
							if ok {
								req = &Payload{
									Model:             &model,
									CommonModelFolder: "",
									CommonModelName:   "",
								}
							}
						} else {
							var modFolder string
							pathParts := strings.Split(pRef, "/")
							modFolder = pathParts[len(pathParts)-2]
							if modFolder == "models" {
								modFolder = ""
							}
							req = &Payload{
								CommonModelFolder: modFolder,
								CommonModelName:   strings.TrimSuffix(filepath.Base(modName), filepath.Ext(modName)),
							}
						}
					}
				}
			}

			ref = v.Subscribe.Message.Ref
			if ref != "" {
				msgName := filepath.Base(ref)
				msg, ok := m.Components.Messages[msgName]
				if ok {
					pRef := msg.Payload.Ref
					if pRef != "" {
						modName := filepath.Base(pRef)
						if strings.HasPrefix(pRef, "#") {
							model, ok := m.Components.Schemas[modName]
							if ok {
								resp = &Payload{
									Model:             &model,
									CommonModelFolder: "",
									CommonModelName:   "",
								}
							}
						} else {
							var modFolder string
							pathParts := strings.Split(pRef, "/")
							modFolder = pathParts[len(pathParts)-2]
							if modFolder == "models" {
								modFolder = ""
							}
							resp = &Payload{
								CommonModelFolder: modFolder,
								CommonModelName:   strings.TrimSuffix(filepath.Base(modName), filepath.Ext(modName)),
							}
						}
					}
				}
			}

		*/

		res.Handlers = append(res.Handlers, Handler{
			Channel:  k,
			Name:     utils.NameFromChannel(k),
			Roles:    v.Roles,
			Request:  req,
			Response: resp,
		})
	}

	sort.Slice(res.Handlers, func(i, j int) bool {
		return res.Handlers[i].Name < res.Handlers[j].Name
	})

	for k, v := range m.Components.Schemas {
		if !strings.HasSuffix(k, "Response") && !strings.HasSuffix(k, "Request") {
			v.Title = utils.UpperCaseFirstLetter(k)
			res.Models = append(res.Models, v)
		}
	}

	return &res, nil
}
