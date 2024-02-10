package cli

import (
	"flag"
	"fmt"
	"github.com/go-playground/validator/v10"
	"os"
	"strings"
)

type Args struct {
	Port uint   `validate:"gte=0,lte=65535"`
	Dir  string `validate:"required"`
	Help bool
}

// defaultPort 0 will just have http.ListenAndServe pick one randomly
const defaultPort = 0
const dirUsage = "directory to serve files from (required)"
const portUsage = "specify port number for the file server (default to a random available port)"
const helpUsage = "view command usage"

const helpMessage = "" +
	"A simple tool for running a quick static file server.\n" +
	"Not intended for production use cases.\n"

var validate *validator.Validate

func MustParseCLIArgs() Args {
	args, errors := ParseCliArgs()
	if args.Help {
		fmt.Println(helpMessage)
		flag.Usage()
		os.Exit(0)
	}
	if errors != nil {
		var badArgs []string
		for _, err := range errors {
			badArgs = append(badArgs, err.Field())
		}

		fmt.Printf("Invalid arguments: %s\n\n", strings.Join(badArgs, ", "))
		flag.Usage()
		os.Exit(1)
	}

	return args
}

func ParseCliArgs() (Args, validator.ValidationErrors) {
	args := Args{}

	flag.BoolVar(&args.Help, "h", false, helpUsage)
	flag.StringVar(&args.Dir, "d", "", dirUsage)
	flag.UintVar(&args.Port, "p", defaultPort, portUsage)

	flag.Parse()
	validate = validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(args)

	if err == nil {
		return args, nil
	}
	return args, err.(validator.ValidationErrors)
}
