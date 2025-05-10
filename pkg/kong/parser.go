package kong_addon

import (
	"github.com/alecthomas/kong"
)

// ParseString constructs a new Kong parser and parses the given string array.
// This has been copied from the Kong source and modified to accept an arbitrary
// string array and not just command args. Call with `os.Args[1:]` to parse
// command line arguments. It will return an error instead of panicking so you
// can handle it the way you want.
func ParseString(cli any, args []string, options ...kong.Option) (*kong.Context, error) {
	parser, err := kong.New(cli, options...)
	if err != nil {
		panic(err)
	}

	ctx, err := parser.Parse(args)
	if err != nil {
		return nil, err
	}

	return ctx, nil
}
