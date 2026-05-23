package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/ploicloud/cli/internal/client"
)

var matchFields = []string{"name", "host", "slug"}

func resolvePathID(ctx context.Context, c *client.Client, pathTmpl string, paramName string, value string, resolved map[string]string) (string, error) {
	if _, err := strconv.Atoi(value); err == nil {
		return value, nil
	}

	marker := "/{" + paramName + "}"
	idx := strings.Index(pathTmpl, marker)
	if idx < 0 {
		return value, nil
	}
	listPath := pathTmpl[:idx]
	if listPath == "" {
		return value, nil
	}

	resp, err := c.Do(ctx, client.Request{
		Method:     "GET",
		PathTmpl:   listPath,
		PathParams: resolved,
		Query:      map[string]string{"per_page": "100"},
	})
	if err != nil {
		return "", err
	}
	if resp.Status >= 400 {
		return value, nil
	}

	var payload struct {
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return value, nil
	}

	items, err := extractItems(payload.Data)
	if err != nil || len(items) == 0 {
		return value, nil
	}

	matches := []map[string]any{}
	for _, it := range items {
		for _, field := range matchFields {
			if v, ok := it[field].(string); ok && v == value {
				matches = append(matches, it)
				break
			}
		}
	}

	switch len(matches) {
	case 0:
		return "", fmt.Errorf("no %s matches %q (searched %s)", paramName, value, url.QueryEscape(listPath))
	case 1:
		id, ok := matches[0]["id"]
		if !ok {
			return value, nil
		}
		idStr := scalarToString(id)
		fmt.Fprintf(os.Stderr, "(resolved %s %q → #%s)\n", paramName, value, idStr)
		return idStr, nil
	default:
		ids := []string{}
		for _, m := range matches {
			ids = append(ids, scalarToString(m["id"]))
		}
		return "", fmt.Errorf("%q is ambiguous, matches %s ids: %s", value, paramName, strings.Join(ids, ", "))
	}
}

func extractItems(raw json.RawMessage) ([]map[string]any, error) {
	if len(raw) == 0 {
		return nil, nil
	}
	if raw[0] == '[' {
		var arr []map[string]any
		if err := json.Unmarshal(raw, &arr); err != nil {
			return nil, err
		}
		return arr, nil
	}
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(raw, &obj); err != nil {
		return nil, err
	}
	for _, key := range []string{"data", "items"} {
		if v, ok := obj[key]; ok && len(v) > 0 && v[0] == '[' {
			var arr []map[string]any
			if err := json.Unmarshal(v, &arr); err == nil {
				return arr, nil
			}
		}
	}
	return nil, nil
}

func scalarToString(v any) string {
	switch x := v.(type) {
	case string:
		return x
	case float64:
		return strconv.FormatInt(int64(x), 10)
	case json.Number:
		return string(x)
	default:
		b, _ := json.Marshal(x)
		return string(b)
	}
}
