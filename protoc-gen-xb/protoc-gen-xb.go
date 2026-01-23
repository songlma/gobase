package main

import (
	"flag"
	"fmt"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
)

func main() {
	var (
		flags        flag.FlagSet
		plugins      = flags.String("plugins", "", "list of plugins to enable (supported values: grpc)")
		importPrefix = flags.String("import_prefix", "", "prefix to prepend to import paths")
	)
	importRewriteFunc := func(importPath protogen.GoImportPath) protogen.GoImportPath {
		switch importPath {
		case "context", "fmt", "math":
			return importPath
		}
		if *importPrefix != "" {
			return protogen.GoImportPath(*importPrefix) + importPath
		}
		return importPath
	}
	protogen.Options{
		ParamFunc:         flags.Set,
		ImportRewriteFunc: importRewriteFunc,
	}.Run(func(gen *protogen.Plugin) error {
		for _, plugin := range strings.Split(*plugins, ",") {
			switch plugin {
			case "grpc":
			case "":
			default:
				return fmt.Errorf("protoc-gen-go: unknown plugin %q", plugin)
			}
		}
		//println(fmt.Sprintf("gen_Files:%+v", gen))
		println("===========================================================")
		for _, f := range gen.Files {
			//println(fmt.Sprintf("f:%+v", f))
			//println(fmt.Sprintf("f.Desc:%+v", f.Desc))
			println(fmt.Sprintf("Services:%+v", f.Desc.Services()))
			println("===========================================================")
			println(fmt.Sprintf("fMessages:%+v", f.Desc.Messages()))
			println("===========================================================")
			println(fmt.Sprintf("FullName:%+v", f.Desc.FullName()))
			println("===========================================================")
			if !f.Generate {
				continue
			}

		}
		return nil
	})
}
