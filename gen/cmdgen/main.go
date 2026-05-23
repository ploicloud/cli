package main

import (
	"encoding/json"
	"fmt"
	"go/format"
	"os"
	"regexp"
	"sort"
	"strings"
)

type spec struct {
	Paths      map[string]map[string]*operation `json:"paths"`
	Components *components                      `json:"components"`
}

type components struct {
	Schemas map[string]*paramSchema `json:"schemas"`
}

type operation struct {
	OperationID string       `json:"operationId"`
	Summary     string       `json:"summary"`
	Description string       `json:"description"`
	Parameters  []parameter  `json:"parameters"`
	RequestBody *requestBody `json:"requestBody"`
}

type parameter struct {
	Name     string        `json:"name"`
	In       string        `json:"in"`
	Required bool          `json:"required"`
	Schema   *paramSchema  `json:"schema"`
	Desc     string        `json:"description"`
}

type requestBody struct {
	Required bool                 `json:"required"`
	Content  map[string]mediaType `json:"content"`
}

type mediaType struct {
	Schema *paramSchema `json:"schema"`
}

type paramSchema struct {
	Ref        string                  `json:"$ref"`
	Type       any                     `json:"type"`
	Properties map[string]*paramSchema `json:"properties"`
	Required   []string                `json:"required"`
	Enum       []any                   `json:"enum"`
	Items      *paramSchema            `json:"items"`
	Format     string                  `json:"format"`
}

func (s *paramSchema) primary() string {
	if s == nil || s.Type == nil {
		return "string"
	}
	switch v := s.Type.(type) {
	case string:
		return v
	case []any:
		for _, e := range v {
			if str, ok := e.(string); ok && str != "null" {
				return str
			}
		}
	}
	return "string"
}

type op struct {
	id     string
	method string
	path   string
	op     *operation
}

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "usage: cmdgen <spec.json> <out.go>")
		os.Exit(2)
	}
	specPath, outPath := os.Args[1], os.Args[2]

	raw, err := os.ReadFile(specPath)
	must(err)
	s := &spec{}
	must(json.Unmarshal(raw, s))

	var ops []op
	for path, methods := range s.Paths {
		for method, o := range methods {
			if o == nil || o.OperationID == "" {
				continue
			}
			m := strings.ToUpper(method)
			if m != "GET" && m != "POST" && m != "PUT" && m != "PATCH" && m != "DELETE" {
				continue
			}
			if o.RequestBody != nil {
				for mt, content := range o.RequestBody.Content {
					content.Schema = resolveRef(content.Schema, s.Components, map[string]bool{})
					o.RequestBody.Content[mt] = content
				}
			}
			displayID := strings.TrimPrefix(o.OperationID, "v1.")
			ops = append(ops, op{id: displayID, method: m, path: path, op: o})
		}
	}
	sort.Slice(ops, func(i, j int) bool { return ops[i].id < ops[j].id })

	out := newWriter()
	out.line("package commands")
	out.line("")
	out.line("import (")
	out.line(`	"context"`)
	out.line(`	"fmt"`)
	out.line(`	"strconv"`)
	out.line("")
	out.line(`	"github.com/spf13/cobra"`)
	out.line(`	"github.com/ploicloud/cli/internal/client"`)
	out.line(")")
	out.line("")
	out.line("var _ = strconv.Itoa")
	out.line("var _ = fmt.Sprintf")
	out.line("var _ context.Context")
	out.line("")

	groups := map[string]string{}
	out.line("func registerGenerated(root *cobra.Command, c *client.Client) {")

	for _, o := range ops {
		segments := strings.Split(o.id, ".")
		if len(segments) < 2 {
			continue
		}
		groupSegs := segments[:len(segments)-1]
		verb := mapVerb(segments[len(segments)-1])

		parent := "root"
		groupKey := ""
		for _, seg := range groupSegs {
			if groupKey == "" {
				groupKey = seg
			} else {
				groupKey = groupKey + "." + seg
			}
			varName, exists := groups[groupKey]
			if !exists {
				varName = "g" + sanitize(groupKey)
				groups[groupKey] = varName
				out.line(fmt.Sprintf("\t%s := &cobra.Command{Use: %q, Short: %q}", varName, seg, "Manage "+seg))
				out.line(fmt.Sprintf("\t%s.AddCommand(%s)", parent, varName))
			}
			parent = varName
		}

		cmdVar := "c" + sanitize(o.id)
		out.line(fmt.Sprintf("\t%s := buildCmd(c, opSpec{", cmdVar))
		out.line(fmt.Sprintf("\t\tID:       %q,", o.id))
		out.line(fmt.Sprintf("\t\tMethod:   %q,", o.method))
		out.line(fmt.Sprintf("\t\tPathTmpl: %q,", o.path))
		out.line(fmt.Sprintf("\t\tUse:      %q,", verb))
		out.line(fmt.Sprintf("\t\tShort:    %q,", o.op.Summary))
		if o.op.Description != "" {
			out.line(fmt.Sprintf("\t\tLong:     %q,", o.op.Description))
		}

		pathParams := []parameter{}
		queryParams := []parameter{}
		for _, p := range o.op.Parameters {
			switch p.In {
			case "path":
				pathParams = append(pathParams, p)
			case "query":
				queryParams = append(queryParams, p)
			}
		}

		out.line("\t\tPathParams: []paramDef{")
		for _, p := range pathParams {
			out.line(fmt.Sprintf("\t\t\t{Name: %q, Type: %q, Required: true, Desc: %q},",
				p.Name, p.Schema.primary(), p.Desc))
		}
		out.line("\t\t},")

		out.line("\t\tQueryParams: []paramDef{")
		for _, p := range queryParams {
			out.line(fmt.Sprintf("\t\t\t{Name: %q, Type: %q, Required: %v, Desc: %q},",
				p.Name, p.Schema.primary(), p.Required, p.Desc))
		}
		out.line("\t\t},")

		out.line("\t\tBodyParams: []paramDef{")
		if o.op.RequestBody != nil {
			if mt, ok := o.op.RequestBody.Content["application/json"]; ok && mt.Schema != nil {
				reqSet := map[string]bool{}
				for _, r := range mt.Schema.Required {
					reqSet[r] = true
				}
				keys := make([]string, 0, len(mt.Schema.Properties))
				for k := range mt.Schema.Properties {
					keys = append(keys, k)
				}
				sort.Strings(keys)
				for _, name := range keys {
					ps := mt.Schema.Properties[name]
					if ps == nil {
						continue
					}
					t := ps.primary()
					if t == "array" && ps.Items != nil {
						t = "array:" + ps.Items.primary()
					}
					out.line(fmt.Sprintf("\t\t\t{Name: %q, Type: %q, Required: %v},",
						name, t, reqSet[name]))
				}
			}
		}
		out.line("\t\t},")
		out.line("\t})")
		out.line(fmt.Sprintf("\t%s.AddCommand(%s)", parent, cmdVar))
		out.line("")
	}
	out.line("}")

	formatted, err := format.Source([]byte(out.String()))
	if err != nil {
		_ = os.WriteFile(outPath, []byte(out.String()), 0o644)
		fmt.Fprintf(os.Stderr, "warning: gofmt failed: %v\n", err)
		os.Exit(1)
	}
	must(os.WriteFile(outPath, formatted, 0o644))
	fmt.Fprintf(os.Stderr, "wrote %s with %d operations\n", outPath, len(ops))
}

func resolveRef(s *paramSchema, comps *components, seen map[string]bool) *paramSchema {
	if s == nil {
		return nil
	}
	if s.Ref != "" {
		prefix := "#/components/schemas/"
		if strings.HasPrefix(s.Ref, prefix) && comps != nil {
			name := strings.TrimPrefix(s.Ref, prefix)
			if seen[name] {
				return s
			}
			if target, ok := comps.Schemas[name]; ok {
				seen[name] = true
				return resolveRef(target, comps, seen)
			}
		}
		return s
	}
	if s.Items != nil {
		s.Items = resolveRef(s.Items, comps, seen)
	}
	for k, v := range s.Properties {
		s.Properties[k] = resolveRef(v, comps, seen)
	}
	return s
}

func mapVerb(v string) string {
	switch v {
	case "index":
		return "list"
	case "store":
		return "create"
	case "show":
		return "get"
	case "destroy":
		return "delete"
	}
	return v
}

var sanitizeRE = regexp.MustCompile(`[^a-zA-Z0-9]+`)

func sanitize(s string) string {
	parts := sanitizeRE.Split(s, -1)
	out := ""
	for _, p := range parts {
		if p == "" {
			continue
		}
		out += strings.ToUpper(p[:1]) + p[1:]
	}
	return out
}

type writer struct{ b strings.Builder }

func newWriter() *writer { return &writer{} }
func (w *writer) line(s string) {
	w.b.WriteString(s)
	w.b.WriteByte('\n')
}
func (w *writer) String() string { return w.b.String() }

func must(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
