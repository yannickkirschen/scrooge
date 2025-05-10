package filter

import (
	"fmt"
	"slices"
	"time"

	"github.com/expr-lang/expr"
)

var BaseEnv = map[string]any{
	"RFC3339":          time.RFC3339,
	"DateTime":         time.DateTime,
	"DateOnly":         time.DateOnly,
	"TimeOnly":         time.TimeOnly,
	"YearOnly":         "2006",
	"MonthOnly":        "01",
	"DayOnly":          "02",
	"YearMonthOnly":    "2006-01",
	"HourOnly":         "15",
	"MinuteOnly":       "04",
	"SecondOnly":       "05",
	"HourMinuteOnly":   "15:04",
	"strSliceContains": func(a []string, e string) bool { return slices.Contains(a, e) },
}

func Matches(env map[string]any, ex string) (bool, error) {
	if ex == "" {
		return true, nil
	}

	program, err := expr.Compile(ex, expr.Env(env))
	if err != nil {
		return false, fmt.Errorf("compilation of expression failed: %s", err)
	}

	output, err := expr.Run(program, env)
	if err != nil {
		return false, fmt.Errorf("filter expression failed: %s", err)
	}

	matches, ok := output.(bool)
	if !ok {
		return false, fmt.Errorf("filter expression %s must return a boolean", ex)
	}

	return matches, nil
}
