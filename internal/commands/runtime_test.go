package commands

import (
	"reflect"
	"testing"

	"github.com/spf13/cobra"
)

var serviceUpdateParams = []paramDef{
	{Name: "version", Type: "string"},
	{Name: "settings", Type: "object"},
	{Name: "settings.command", Type: "string"},
	{Name: "settings.extensions", Type: "array:string"},
	{Name: "settings.maxmemory_policy", Type: "string"},
	{Name: "settings.volume_size", Type: "integer"},
}

func bodyFor(t *testing.T, params []paramDef, args ...string) map[string]any {
	t.Helper()

	cmd := &cobra.Command{Use: "update"}
	vals := map[string]any{}
	for _, p := range params {
		bindFlag(cmd, p, vals)
	}
	if err := cmd.Flags().Parse(args); err != nil {
		t.Fatalf("parse %v: %v", args, err)
	}
	return assembleBody(cmd, params, vals)
}

func TestBodyNestsDottedFlag(t *testing.T) {
	got := bodyFor(t, serviceUpdateParams, "--settings-maxmemory-policy", "allkeys-lru")

	want := map[string]any{
		"settings": map[string]any{"maxmemory_policy": "allkeys-lru"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("body = %#v, want %#v", got, want)
	}
}

func TestBodyOmitsUnsetFlags(t *testing.T) {
	got := bodyFor(t, serviceUpdateParams)

	if got != nil {
		t.Fatalf("body = %#v, want nil; sending untouched keys overwrites stored settings", got)
	}
}

func TestBodyOmitsSiblingsOfASetFlag(t *testing.T) {
	got := bodyFor(t, serviceUpdateParams, "--settings-command", "php artisan queue:work")

	settings, ok := got["settings"].(map[string]any)
	if !ok {
		t.Fatalf("settings missing or not an object: %#v", got)
	}
	if len(settings) != 1 {
		t.Errorf("settings = %#v, want only the flag that was set", settings)
	}
	if _, exists := got["version"]; exists {
		t.Error("unset top-level flag leaked into the body")
	}
}

func TestBodyPreservesDeclaredTypes(t *testing.T) {
	got := bodyFor(t, serviceUpdateParams,
		"--settings-volume-size", "20",
		"--settings-extensions", "postgis,uuid-ossp",
	)

	settings := got["settings"].(map[string]any)

	if size, ok := settings["volume_size"].(int); !ok || size != 20 {
		t.Errorf("volume_size = %#v, want int 20", settings["volume_size"])
	}

	want := []string{"postgis", "uuid-ossp"}
	if ext, ok := settings["extensions"].([]string); !ok || !reflect.DeepEqual(ext, want) {
		t.Errorf("extensions = %#v, want %v", settings["extensions"], want)
	}
}

func TestBodyMergesJSONFlagWithLeafFlags(t *testing.T) {
	got := bodyFor(t, serviceUpdateParams,
		"--settings", `{"memory_request":"512Mi"}`,
		"--settings-command", "php artisan horizon",
	)

	settings, ok := got["settings"].(map[string]any)
	if !ok {
		t.Fatalf("settings missing or not an object: %#v", got)
	}
	if settings["memory_request"] != "512Mi" {
		t.Errorf("key present only in the JSON document was lost: %#v", settings)
	}
	if settings["command"] != "php artisan horizon" {
		t.Errorf("command = %#v, want the leaf flag value", settings["command"])
	}
}

func TestLeafFlagOverridesJSONFlag(t *testing.T) {
	got := bodyFor(t, serviceUpdateParams,
		"--settings", `{"maxmemory_policy":"volatile-ttl"}`,
		"--settings-maxmemory-policy", "allkeys-lfu",
	)

	settings := got["settings"].(map[string]any)
	if settings["maxmemory_policy"] != "allkeys-lfu" {
		t.Errorf("maxmemory_policy = %#v, want the more specific flag to win", settings["maxmemory_policy"])
	}
}

func TestJSONFlagAloneStillWorks(t *testing.T) {
	got := bodyFor(t, serviceUpdateParams, "--settings", `{"maxmemory_policy":"volatile-ttl"}`)

	want := map[string]any{
		"settings": map[string]any{"maxmemory_policy": "volatile-ttl"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("body = %#v, want %#v", got, want)
	}
}

func TestMalformedJSONFlagIsNotSilentlyReplaced(t *testing.T) {
	got := bodyFor(t, serviceUpdateParams,
		"--settings", `{not json`,
		"--settings-command", "php artisan horizon",
	)

	if _, isObject := got["settings"].(map[string]any); isObject {
		t.Fatalf("undecodable parent value was replaced, hiding the mistake: %#v", got)
	}
	if got["settings"] != `{not json` {
		t.Errorf("settings = %#v, want the raw value passed through for the API to reject", got["settings"])
	}
}

func TestTopLevelFlagIsNotNested(t *testing.T) {
	got := bodyFor(t, serviceUpdateParams, "--version", "8.1")

	want := map[string]any{"version": "8.1"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("body = %#v, want %#v", got, want)
	}
}

func TestKebabRendersDotsAndUnderscores(t *testing.T) {
	cases := map[string]string{
		"settings":                  "settings",
		"settings.maxmemory_policy": "settings-maxmemory-policy",
		"php_settings":              "php-settings",
	}
	for in, want := range cases {
		if got := kebab(in); got != want {
			t.Errorf("kebab(%q) = %q, want %q", in, got, want)
		}
	}
}
