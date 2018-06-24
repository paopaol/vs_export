package main

import (
	//"awesomeProject/src/github.com/gin-gonic/gin/json"
    "encoding/json"
	"flag"
	"fmt"
	"os"
	"sln"
)

func main() {
	path := flag.String("s", "", "sln file path")
	configuration := flag.String("c", "Debug|Win32",
		"Configuration, [configuration|platform], default Debug|Win32")
	flag.Parse()

	if *path == "" {
		usage()
		os.Exit(1)
	}

	solution, err := sln.NewSln(*path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	cmdList, err := solution.CompileCommandsJson(*configuration)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	js, err := json.Marshal(cmdList)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("%s\n", js[:])
}


func usage() {
	echo := `usage:sln_export_compile_commands options
			 -s   path                        sln filename
           -c   configuration               project configuration,eg Debug|Win32.
                                            default Debug|Win32
	`
	fmt.Println(echo)
}