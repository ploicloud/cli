package output

import (
	"bytes"
	"strings"
	"testing"
)

func renderList(t *testing.T, items []any) string {
	t.Helper()
	var buf bytes.Buffer
	if err := printList(&buf, items); err != nil {
		t.Fatalf("printList: %v", err)
	}
	return buf.String()
}

func valkeyService() map[string]any {
	return map[string]any{
		"id":       float64(1326),
		"name":     "cache",
		"type":     "valkey",
		"status":   "active",
		"settings": map[string]any{"maxmemory_policy": "allkeys-lru", "memory_request": "512Mi", "volume_size": "4Gi", "password": "********"},
	}
}

func TestListShowsNestedSettings(t *testing.T) {
	out := renderList(t, []any{valkeyService()})

	if !strings.Contains(out, "SETTINGS") {
		t.Fatalf("no settings column:\n%s", out)
	}
	if !strings.Contains(out, "maxmemory_policy=allkeys-lru") {
		t.Fatalf("eviction policy not readable in list output:\n%s", out)
	}
	if !strings.Contains(out, "memory_request=512Mi") {
		t.Fatalf("sibling settings dropped:\n%s", out)
	}
}

func TestListKeepsScalarColumns(t *testing.T) {
	out := renderList(t, []any{valkeyService()})

	for _, want := range []string{"ID", "NAME", "TYPE", "STATUS", "1326", "cache", "valkey", "active"} {
		if !strings.Contains(out, want) {
			t.Fatalf("missing %q in:\n%s", want, out)
		}
	}
}

func TestListShowsSettingsWithMoreThanFourKeys(t *testing.T) {
	mysql := map[string]any{
		"id":   float64(1033),
		"name": "mysql",
		"type": "mysql",
		"settings": map[string]any{
			"database": "app_db", "username": "dbuser", "password": "********",
			"root_password": "********", "volume_size": "4Gi", "memory_request": "512Mi",
		},
	}
	out := renderList(t, []any{mysql})

	if !strings.Contains(out, "volume_size=4Gi") {
		t.Fatalf("six-key settings were dropped:\n%s", out)
	}
}

func TestListRowsWithoutSettingsAreUnaffected(t *testing.T) {
	app := map[string]any{
		"id":               float64(873),
		"name":             "laravel",
		"application_type": "laravel",
		"status":           "active",
		"domains":          []any{map[string]any{"domain": "example.com"}},
	}
	out := renderList(t, []any{app})

	if strings.Contains(out, "SETTINGS") {
		t.Fatalf("settings column appeared for a payload without settings:\n%s", out)
	}
	if strings.Contains(out, "DOMAINS") {
		t.Fatalf("list-valued column should still be skipped:\n%s", out)
	}
}

func TestEmptyAndDeeplyNestedSettingsFallBack(t *testing.T) {
	for name, settings := range map[string]any{
		"empty":  map[string]any{},
		"nested": map[string]any{"egress": map[string]any{"enabled": true}},
	} {
		out := renderList(t, []any{map[string]any{"id": float64(1), "name": "x", "settings": settings}})
		if strings.Contains(out, "map[") {
			t.Fatalf("%s: raw Go map leaked into output:\n%s", name, out)
		}
	}
}

func TestListMasksCredentials(t *testing.T) {
	out := renderList(t, []any{map[string]any{
		"id":   float64(1033),
		"name": "mysql",
		"settings": map[string]any{
			"username": "dbuser", "password": "PFTDMdA0QG9ThLbSUoG4Kudb",
			"root_password": "bPvGQZ39pKBZilRz1DTV4KEK", "memory_request": "512Mi",
		},
	}})

	for _, leaked := range []string{"PFTDMdA0QG9ThLbSUoG4Kudb", "bPvGQZ39pKBZilRz1DTV4KEK"} {
		if strings.Contains(out, leaked) {
			t.Fatalf("credential printed in plain output:\n%s", out)
		}
	}
	if !strings.Contains(out, "password=********") || !strings.Contains(out, "root_password=********") {
		t.Fatalf("credentials should still be listed as masked:\n%s", out)
	}
	if !strings.Contains(out, "username=dbuser") || !strings.Contains(out, "memory_request=512Mi") {
		t.Fatalf("non-secret settings must stay visible:\n%s", out)
	}
}

func TestObjectOutputMasksCredentials(t *testing.T) {
	var buf bytes.Buffer
	if err := printObject(&buf, map[string]any{
		"name":     "cache",
		"settings": map[string]any{"password": "zWn7edNfFsB9yOQbuiL8", "maxmemory_policy": "allkeys-lru"},
	}); err != nil {
		t.Fatal(err)
	}
	out := buf.String()

	if strings.Contains(out, "zWn7edNfFsB9yOQbuiL8") {
		t.Fatalf("credential printed by single-object output:\n%s", out)
	}
	if !strings.Contains(out, "allkeys-lru") {
		t.Fatalf("non-secret value lost:\n%s", out)
	}
}

func TestSecretKeyClassification(t *testing.T) {
	for k, want := range map[string]bool{
		"password": true, "root_password": true, "api_key": true, "access_token": true,
		"secret": true, "apikey": true,
		"maxmemory_policy": false, "memory_request": false, "username": false,
		"public_key": false, "key_id": false, "keyword": false, "volume_size": false,
	} {
		if got := isSecretKey(k); got != want {
			t.Errorf("isSecretKey(%q) = %v, want %v", k, got, want)
		}
	}
}
