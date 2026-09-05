package main

import (
	"encoding/json"
	"strings"
	"testing"
)

func schemaFrom(t *testing.T, raw string) *paramSchema {
	t.Helper()
	s := &paramSchema{}
	if err := json.Unmarshal([]byte(raw), s); err != nil {
		t.Fatalf("unmarshal schema: %v", err)
	}
	return s
}

func find(params []bodyParam, name string) (bodyParam, bool) {
	for _, p := range params {
		if p.Name == name {
			return p, true
		}
	}
	return bodyParam{}, false
}

func names(params []bodyParam) []string {
	out := make([]string, 0, len(params))
	for _, p := range params {
		out = append(out, p.Name)
	}
	return out
}

func TestBodyParamsFlattensObjectProperties(t *testing.T) {
	schema := schemaFrom(t, `{
		"type": "object",
		"required": ["type", "name"],
		"properties": {
			"type": {"type": "string"},
			"name": {"type": "string"},
			"settings": {
				"type": "object",
				"properties": {
					"memory_request": {"type": ["string", "null"]},
					"replicas": {"type": "integer"}
				}
			}
		}
	}`)

	params := bodyParams(schema)

	for _, want := range []string{"type", "name", "settings", "settings.memory_request", "settings.replicas"} {
		if _, ok := find(params, want); !ok {
			t.Errorf("missing param %q, got %v", want, names(params))
		}
	}

	parent, _ := find(params, "settings")
	if parent.Type != "object" {
		t.Errorf("parent settings type = %q, want object", parent.Type)
	}

	replicas, _ := find(params, "settings.replicas")
	if replicas.Type != "integer" {
		t.Errorf("settings.replicas type = %q, want integer", replicas.Type)
	}

	memory, _ := find(params, "settings.memory_request")
	if memory.Type != "string" {
		t.Errorf("settings.memory_request type = %q, want string", memory.Type)
	}

	typeParam, _ := find(params, "type")
	if !typeParam.Required {
		t.Error("top-level required property lost its Required flag")
	}
}

func TestBodyParamsDescribesEnum(t *testing.T) {
	schema := schemaFrom(t, `{
		"type": "object",
		"properties": {
			"settings": {
				"type": "object",
				"properties": {
					"maxmemory_policy": {
						"type": "string",
						"enum": ["noeviction", "allkeys-lru", "volatile-ttl"]
					}
				}
			}
		}
	}`)

	policy, ok := find(bodyParams(schema), "settings.maxmemory_policy")
	if !ok {
		t.Fatal("settings.maxmemory_policy not generated")
	}
	want := "one of: noeviction, allkeys-lru, volatile-ttl"
	if policy.Desc != want {
		t.Errorf("Desc = %q, want %q", policy.Desc, want)
	}
}

func TestBodyParamsDescriptionAndEnumTogether(t *testing.T) {
	schema := schemaFrom(t, `{
		"type": "object",
		"properties": {
			"region": {
				"type": "string",
				"description": "Deployment region",
				"enum": ["ams1", "fra1"]
			},
			"label": {"type": "string", "description": "Display name"}
		}
	}`)

	params := bodyParams(schema)

	region, _ := find(params, "region")
	if region.Desc != "one of: ams1, fra1; Deployment region" {
		t.Errorf("region Desc = %q", region.Desc)
	}

	label, _ := find(params, "label")
	if label.Desc != "Display name" {
		t.Errorf("label Desc = %q", label.Desc)
	}
}

func TestBodyParamsEmptyDescriptionIsNotFatal(t *testing.T) {
	schema := schemaFrom(t, `{
		"type": "object",
		"properties": {"version": {"type": "string"}}
	}`)

	version, ok := find(bodyParams(schema), "version")
	if !ok {
		t.Fatal("version not generated")
	}
	if version.Desc != "" {
		t.Errorf("Desc = %q, want empty", version.Desc)
	}
}

func TestBodyParamsArrayLeafKeepsItemType(t *testing.T) {
	schema := schemaFrom(t, `{
		"type": "object",
		"properties": {
			"settings": {
				"type": "object",
				"properties": {
					"extensions": {"type": "array", "items": {"type": "string"}}
				}
			}
		}
	}`)

	ext, ok := find(bodyParams(schema), "settings.extensions")
	if !ok {
		t.Fatal("settings.extensions not generated")
	}
	if ext.Type != "array:string" {
		t.Errorf("Type = %q, want array:string", ext.Type)
	}
}

func TestBodyParamsStopsAtOneLevel(t *testing.T) {
	schema := schemaFrom(t, `{
		"type": "object",
		"properties": {
			"security": {
				"type": "object",
				"properties": {
					"headers": {
						"type": "object",
						"properties": {"X-Frame-Options": {"type": "string"}}
					}
				}
			}
		}
	}`)

	params := bodyParams(schema)

	if _, ok := find(params, "security.headers"); !ok {
		t.Error("depth-1 child security.headers should be generated")
	}
	if _, ok := find(params, "security.headers.X-Frame-Options"); ok {
		t.Errorf("depth-2 grandchild should not be generated, got %v", names(params))
	}

	headers, _ := find(params, "security.headers")
	if headers.Type != "object" {
		t.Errorf("depth-1 object child Type = %q, want object", headers.Type)
	}
}

func TestBodyParamsIgnoresNestedRequired(t *testing.T) {
	schema := schemaFrom(t, `{
		"type": "object",
		"properties": {
			"settings": {
				"type": "object",
				"required": ["command"],
				"properties": {"command": {"type": "string"}}
			}
		}
	}`)

	command, ok := find(bodyParams(schema), "settings.command")
	if !ok {
		t.Fatal("settings.command not generated")
	}
	if command.Required {
		t.Error("leaf marked required from a nested required list; conditional requirements belong to the API")
	}
}

func TestBodyParamsResolvedRefSchema(t *testing.T) {
	raw := `{
		"paths": {
			"/services": {
				"post": {
					"operationId": "services.store",
					"requestBody": {
						"content": {
							"application/json": {
								"schema": {"$ref": "#/components/schemas/StoreServiceRequest"}
							}
						}
					}
				}
			}
		},
		"components": {
			"schemas": {
				"StoreServiceRequest": {
					"type": "object",
					"properties": {
						"settings": {
							"type": "object",
							"properties": {"volume_size": {"type": "string"}}
						}
					}
				}
			}
		}
	}`

	s := &spec{}
	if err := json.Unmarshal([]byte(raw), s); err != nil {
		t.Fatalf("unmarshal spec: %v", err)
	}
	mt := s.Paths["/services"]["post"].RequestBody.Content["application/json"]

	if got := bodyParams(mt.Schema); len(got) != 0 {
		t.Fatalf("expected an unresolved $ref to yield nothing, got %v", names(got))
	}

	mt.Schema = resolveRef(mt.Schema, s.Components, map[string]bool{})
	params := bodyParams(mt.Schema)
	if _, ok := find(params, "settings.volume_size"); !ok {
		t.Errorf("referenced schema not walked after resolution, got %v", names(params))
	}
}

func TestBodyParamsNilAndEmpty(t *testing.T) {
	if got := bodyParams(nil); got != nil {
		t.Errorf("nil schema returned %v", got)
	}
	if got := bodyParams(schemaFrom(t, `{"type": "object"}`)); got != nil {
		t.Errorf("property-less schema returned %v", got)
	}
}

func TestCheckCollisions(t *testing.T) {
	cases := []struct {
		name    string
		query   []parameter
		body    []bodyParam
		wantErr bool
	}{
		{
			name: "no collision",
			body: []bodyParam{{Name: "settings"}, {Name: "settings.command"}},
		},
		{
			name:    "leaf collides with sibling top-level property",
			body:    []bodyParam{{Name: "settings_command"}, {Name: "settings.command"}},
			wantErr: true,
		},
		{
			name:    "body collides with query param",
			query:   []parameter{{Name: "per_page"}},
			body:    []bodyParam{{Name: "per.page"}},
			wantErr: true,
		},
		{
			name:  "distinct query and body",
			query: []parameter{{Name: "dry_run"}},
			body:  []bodyParam{{Name: "settings.command"}},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := checkCollisions("applications.services.update", tc.query, tc.body)
			if tc.wantErr && err == nil {
				t.Fatal("expected a collision error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tc.wantErr {
				msg := err.Error()
				if !strings.Contains(msg, "applications.services.update") {
					t.Errorf("error does not name the operation: %q", msg)
				}
				if !strings.Contains(msg, "--settings-command") && !strings.Contains(msg, "--per-page") {
					t.Errorf("error does not name the flag: %q", msg)
				}
			}
		})
	}
}

func TestFlagName(t *testing.T) {
	cases := map[string]string{
		"settings":                  "settings",
		"settings.maxmemory_policy": "settings-maxmemory-policy",
		"php_settings":              "php-settings",
	}
	for in, want := range cases {
		if got := flagName(in); got != want {
			t.Errorf("flagName(%q) = %q, want %q", in, got, want)
		}
	}
}
