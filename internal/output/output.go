package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/ploicloud/cli/internal/client"
)

var JSON bool

func Print(resp *client.Response) error {
	v, err := resp.JSON()
	if err != nil {
		return err
	}
	w := os.Stdout
	if resp.Status >= 400 {
		fmt.Fprintf(os.Stderr, "Request failed: HTTP %d\n", resp.Status)
	}
	if JSON {
		return printJSON(w, v)
	}
	return printPretty(w, v)
}

func printJSON(w io.Writer, v any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func printPretty(w io.Writer, v any) error {
	m, ok := v.(map[string]any)
	if !ok {
		return printValue(w, v)
	}
	if msg, ok := m["message"].(string); ok && msg != "" {
		fmt.Fprintln(w, msg)
	}
	data, hasData := m["data"]
	if !hasData {
		filtered := map[string]any{}
		for k, val := range m {
			if k == "success" || k == "message" {
				continue
			}
			filtered[k] = val
		}
		if len(filtered) == 0 {
			if ok, _ := m["success"].(bool); ok {
				fmt.Fprintln(w, "OK")
			}
			return nil
		}
		return printValue(w, filtered)
	}
	if data == nil {
		return nil
	}
	return printValue(w, data)
}

func printValue(w io.Writer, v any) error {
	switch d := v.(type) {
	case nil:
		return nil
	case map[string]any:
		return printObject(w, d)
	case []any:
		return printList(w, d)
	default:
		fmt.Fprintln(w, formatScalar(v))
		return nil
	}
}

func printObject(w io.Writer, obj map[string]any) error {
	return printObjectIndented(w, obj, "")
}

func printObjectIndented(w io.Writer, obj map[string]any, indent string) error {
	keys := sortedKeys(obj)
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	for _, k := range keys {
		v := obj[k]
		switch vv := v.(type) {
		case nil:
			fmt.Fprintf(tw, "%s%s\t-\n", indent, k)
		case []any:
			if len(vv) == 0 {
				fmt.Fprintf(tw, "%s%s\t(empty)\n", indent, k)
			} else if isListOfScalars(vv) {
				fmt.Fprintf(tw, "%s%s\t%s\n", indent, k, joinScalars(vv))
			} else {
				fmt.Fprintf(tw, "%s%s\t(%d items)\n", indent, k, len(vv))
			}
		case map[string]any:
			fmt.Fprintf(tw, "%s%s\n", indent, k)
			tw.Flush()
			_ = printObjectIndented(w, vv, indent+"  ")
		default:
			fmt.Fprintf(tw, "%s%s\t%s\n", indent, k, formatScalar(v))
		}
	}
	return tw.Flush()
}

func isListOfScalars(items []any) bool {
	for _, it := range items {
		switch it.(type) {
		case map[string]any, []any:
			return false
		}
	}
	return true
}

func joinScalars(items []any) string {
	parts := make([]string, 0, len(items))
	for _, it := range items {
		parts = append(parts, formatScalar(it))
	}
	return strings.Join(parts, ", ")
}

func printList(w io.Writer, items []any) error {
	if len(items) == 0 {
		fmt.Fprintln(w, "No items.")
		return nil
	}
	first, ok := items[0].(map[string]any)
	if !ok {
		for _, it := range items {
			fmt.Fprintln(w, formatScalar(it))
		}
		return nil
	}
	cols := pickColumns(first)
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	headers := make([]string, len(cols))
	for i, c := range cols {
		headers[i] = strings.ToUpper(c)
	}
	fmt.Fprintln(tw, strings.Join(headers, "\t"))
	for _, it := range items {
		m, ok := it.(map[string]any)
		if !ok {
			continue
		}
		row := make([]string, len(cols))
		for i, c := range cols {
			row[i] = formatScalar(m[c])
		}
		fmt.Fprintln(tw, strings.Join(row, "\t"))
	}
	return tw.Flush()
}

func pickColumns(obj map[string]any) []string {
	preferred := []string{"id", "name", "slug", "type", "application_type", "status", "email", "host", "url", "created_at", "updated_at"}
	cols := []string{}
	seen := map[string]bool{}
	for _, p := range preferred {
		if v, ok := obj[p]; ok {
			if _, isMap := v.(map[string]any); isMap {
				continue
			}
			if _, isList := v.([]any); isList {
				continue
			}
			cols = append(cols, p)
			seen[p] = true
		}
	}
	if len(cols) < 4 {
		for _, k := range sortedKeys(obj) {
			if seen[k] {
				continue
			}
			v := obj[k]
			if _, isMap := v.(map[string]any); isMap {
				continue
			}
			if _, isList := v.([]any); isList {
				continue
			}
			cols = append(cols, k)
			if len(cols) >= 4 {
				break
			}
		}
	}
	if len(cols) == 0 {
		cols = sortedKeys(obj)
		if len(cols) > 4 {
			cols = cols[:4]
		}
	}
	return cols
}

func sortedKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func formatScalar(v any) string {
	switch x := v.(type) {
	case nil:
		return "-"
	case string:
		if x == "" {
			return "-"
		}
		return x
	case bool:
		if x {
			return "true"
		}
		return "false"
	case float64:
		if x == float64(int64(x)) {
			return fmt.Sprintf("%d", int64(x))
		}
		return fmt.Sprintf("%v", x)
	case json.Number:
		return string(x)
	default:
		b, err := json.Marshal(x)
		if err != nil {
			return fmt.Sprintf("%v", x)
		}
		return string(b)
	}
}

func formatNestedObject(m map[string]any) string {
	if len(m) == 0 || len(m) > 4 {
		return ""
	}
	parts := []string{}
	for _, k := range sortedKeys(m) {
		v := m[k]
		switch v.(type) {
		case map[string]any, []any:
			return ""
		}
		s := formatScalar(v)
		if s == "-" {
			continue
		}
		parts = append(parts, fmt.Sprintf("%s=%s", k, s))
	}
	return strings.Join(parts, "  ")
}
