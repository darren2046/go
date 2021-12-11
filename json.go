package golanglibs

import (
	"bytes"
	"encoding/json"

	"github.com/ghodss/yaml"
	goyaml "gopkg.in/yaml.v2"
)

type jsonStruct struct {
	Dumps     func(v interface{}, pretty ...bool) string
	Loads     func(str string) map[string]interface{}
	Valid     func(j string) bool
	Yaml2json func(y string) string
	Json2yaml func(j string) string
	Format    func(js string) string
}

var Json jsonStruct

func init() {
	Json = jsonStruct{
		Dumps:     jsonDumps,
		Loads:     jsonLoads,
		Valid:     jsonValid,
		Yaml2json: yaml2json,
		Json2yaml: json2yaml,
		Format:    formatJson,
	}
}

func formatJson(js string) string {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(js), "", "    "); err != nil {
		panicerr(err)
	}
	return prettyJSON.String()
}

func yaml2json(y string) string {
	goyaml.FutureLineWrap()
	outBytes, err := yaml.YAMLToJSON([]byte(y))
	panicerr(err)
	return Str(outBytes)
}

func json2yaml(j string) string {
	if !Json.Valid(j) {
		panicerr("Not a json string")
	}
	goyaml.FutureLineWrap()
	outBytes, err := yaml.JSONToYAML([]byte(j))
	panicerr(err)
	return Str(outBytes)
}

func jsonValid(j string) bool {
	return json.Valid([]byte(j))
}

func jsonDumps(v interface{}, pretty ...bool) string {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	if len(pretty) != 0 {
		encoder.SetIndent(" ", " ")
	}
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(v)

	panicerr(err)
	return String(buffer.String()).Strip().Get()
}

func jsonLoads(str string) map[string]interface{} {
	var dat map[string]interface{}
	err := json.Unmarshal([]byte(str), &dat)
	panicerr(err)
	return dat
}
