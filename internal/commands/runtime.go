package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/ploicloud/cli/internal/client"
	"github.com/ploicloud/cli/internal/output"
	"github.com/spf13/cobra"
)

type paramDef struct {
	Name     string
	Type     string
	Required bool
	Desc     string
}

type opSpec struct {
	ID          string
	Method      string
	PathTmpl    string
	Use         string
	Short       string
	Long        string
	PathParams  []paramDef
	QueryParams []paramDef
	BodyParams  []paramDef
}

func buildCmd(c *client.Client, op opSpec) *cobra.Command {
	use := op.Use
	for _, p := range op.PathParams {
		use += " <" + kebab(p.Name) + ">"
	}

	cmd := &cobra.Command{
		Use:   use,
		Short: op.Short,
		Long:  op.Long,
		Args:  cobra.ExactArgs(len(op.PathParams)),
	}

	flagVals := map[string]any{}
	for _, p := range op.QueryParams {
		bindFlag(cmd, p, flagVals)
	}
	for _, p := range op.BodyParams {
		bindFlag(cmd, p, flagVals)
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		pathParams := map[string]string{}
		for i, p := range op.PathParams {
			value := args[i]
			if p.Type == "integer" {
				resolved, err := resolvePathID(ctx, c, op.PathTmpl, p.Name, value, pathParams)
				if err != nil {
					return err
				}
				value = resolved
			}
			pathParams[p.Name] = value
		}
		query := map[string]string{}
		for _, p := range op.QueryParams {
			if cmd.Flags().Changed(kebab(p.Name)) {
				query[p.Name] = stringifyFlag(p, flagVals[p.Name])
			}
		}
		var body map[string]any
		for _, p := range op.BodyParams {
			if cmd.Flags().Changed(kebab(p.Name)) {
				if body == nil {
					body = map[string]any{}
				}
				body[p.Name] = unwrapFlag(p, flagVals[p.Name])
			}
		}

		req := client.Request{
			Method:     op.Method,
			PathTmpl:   op.PathTmpl,
			PathParams: pathParams,
			Query:      query,
		}
		if body != nil {
			req.Body = body
		}

		resp, err := c.Do(ctx, req)
		if err != nil {
			return err
		}
		return printResponse(resp)
	}
	return cmd
}

func bindFlag(cmd *cobra.Command, p paramDef, vals map[string]any) {
	name := kebab(p.Name)
	switch {
	case p.Type == "integer":
		v := new(int)
		cmd.Flags().IntVar(v, name, 0, p.Desc)
		vals[p.Name] = v
	case p.Type == "number":
		v := new(float64)
		cmd.Flags().Float64Var(v, name, 0, p.Desc)
		vals[p.Name] = v
	case p.Type == "boolean":
		v := new(bool)
		cmd.Flags().BoolVar(v, name, false, p.Desc)
		vals[p.Name] = v
	case strings.HasPrefix(p.Type, "array:"):
		v := new([]string)
		cmd.Flags().StringSliceVar(v, name, nil, p.Desc)
		vals[p.Name] = v
	default:
		v := new(string)
		cmd.Flags().StringVar(v, name, "", p.Desc)
		vals[p.Name] = v
	}
}

func stringifyFlag(p paramDef, v any) string {
	switch x := v.(type) {
	case *string:
		return *x
	case *int:
		return strconv.Itoa(*x)
	case *float64:
		return strconv.FormatFloat(*x, 'f', -1, 64)
	case *bool:
		return strconv.FormatBool(*x)
	case *[]string:
		return strings.Join(*x, ",")
	}
	return fmt.Sprintf("%v", v)
}

func unwrapFlag(p paramDef, v any) any {
	switch x := v.(type) {
	case *string:
		s := *x
		if p.Type == "object" && s != "" {
			var parsed any
			if err := json.Unmarshal([]byte(s), &parsed); err == nil {
				return parsed
			}
		}
		return s
	case *int:
		return *x
	case *float64:
		return *x
	case *bool:
		return *x
	case *[]string:
		return *x
	}
	return v
}

func kebab(s string) string {
	return strings.ReplaceAll(s, "_", "-")
}

func printResponse(resp *client.Response) error {
	return output.Print(resp)
}

func Register(root *cobra.Command, c *client.Client) {
	registerGenerated(root, c)
}
