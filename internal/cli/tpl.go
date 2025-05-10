package cli

import (
	"fmt"
	"text/template"
	"time"
)

var TplBaseFunctions = template.FuncMap{
	"now": func() string {
		return time.Now().Format(time.RFC3339)
	},
	"before": func(a, b any) (bool, error) {
		t1, err := time.Parse(time.RFC3339, fmt.Sprintf("%v", a))
		if err != nil {
			return false, fmt.Errorf("error parsing first argument %s as ISO-8601 time: %s", a, err)
		}

		t2, err := time.Parse(time.RFC3339, fmt.Sprintf("%v", b))
		if err != nil {
			return false, fmt.Errorf("error parsing first argument %s as ISO-8601 time: %s", a, err)
		}

		return t1.Before(t2), nil
	},
}
