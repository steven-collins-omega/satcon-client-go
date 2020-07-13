package actions

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"text/template"
)

type QueryType string

const (
	QueryTypeQuery    QueryType = "query"
	QueryTypeMutation QueryType = "mutation"
)

const (
	QueryTemplate = `{"query":"{{.Type}} ({{buildArgsList .Args}}) {\n  {{.QueryName}}({{buildArgVarsList .Args}}{{print ") {"}}{{range .Returns}}{{printf "\\n    %s" .}}{{end}}\n  }\n}","variables":{{print "{"}}{{block "vars" .}}{{end}}{{print "}}"}}`
)

type GraphQLQuery struct {
	Type      QueryType
	QueryName string
	Args      map[string]string
	Returns   []string
}

func BuildArgsList(args map[string]string) string {
	argStrings := make([]string, 0)

	for k, v := range args {
		argStrings = append(argStrings, fmt.Sprintf("$%s: %s", k, v))
	}

	return strings.Join(argStrings, ", ")
}

func BuildArgVarsList(args map[string]string) string {
	argVarStrings := make([]string, 0)

	for k, _ := range args {
		argVarStrings = append(argVarStrings, fmt.Sprintf("%s: $%s", k, k))
	}

	return strings.Join(argVarStrings, ", ")
}

func BuildRequest(payload io.Reader, endpoint, token string) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, endpoint, payload)
	req.Header.Add("content-type", "application/json")
	req.Header.Add("authorization", fmt.Sprintf("Bearer %s", token))
	return req
}

// BuildRequestBody takes the request(query)-specific variables template in and uses it to overlay
// the common QueryTemplate using the supplied variables, then returns an io.Reader for
// use as a GraphQL http request body.  The user can supply additional functions for
// use during template processing, but for most queries it should be safe to pass nil
// and rely on the predefined template.FuncMap.
//
// NOTE: The supplied template *must* include an inlined template definition for "vars",
// e.g.:
// `{{define "vars"}}"var_1":"{{.Var1}},"var2":"{{.Var2}}"{{end}}`
func BuildRequestBody(requestTemplate string, vars interface{}, funcs template.FuncMap) (io.Reader, error) {
	defaultFuncs := template.FuncMap{
		"buildArgsList":    BuildArgsList,
		"buildArgVarsList": BuildArgVarsList,
	}

	// Merge in user-supplied functions
	if funcs != nil {
		for k, v := range funcs {
			defaultFuncs[k] = v
		}
	}

	master, err := template.New("master").Funcs(defaultFuncs).Parse(QueryTemplate)
	if err != nil {
		return nil, errors.New("Unable to parse master template")
	}

	final, err := template.Must(master.Clone()).Parse(requestTemplate)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse supplied template: %s", err)
	}

	buf := &bytes.Buffer{}
	err = final.Execute(buf, vars)
	if err != nil {
		return nil, fmt.Errorf("Unable to execute template: %s", err)
	}

	return buf, nil
}