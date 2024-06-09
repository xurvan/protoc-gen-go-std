package main

import (
	"log"
	"os/exec"
	"strings"
	"text/template"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"

	"protoc-gen-go-std/assets"
)

const Version = "0.1.0"

type Param struct {
	PackageName     string
	PluginVersion   string
	CompilerVersion string
	FileName        string
	Imports         []string
	Name            string
	ProjectName     string
}

func generateFile(plugin *protogen.Plugin, file *protogen.File) {
	filename := file.GeneratedFilenamePrefix + "_std.pb.go"
	generatedFile := plugin.NewGeneratedFile(filename, file.GoImportPath)

	firstMessageName := (*file.Messages[0]).Desc.Name()
	projectName := strings.Split(string(file.GoImportPath), "/")[0]

	cmd := exec.Command("go", "version")
	out, err := cmd.Output()
	if err != nil {
		log.Fatal("Error getting Go version:", err)
		return
	}
	compilerVersion := string(out)

	data := Param{
		PackageName:     string(file.GoPackageName),
		PluginVersion:   Version,
		CompilerVersion: compilerVersion,
		FileName:        file.Proto.GetName(),
		Name:            string(firstMessageName),
		ProjectName:     projectName,
	}
	tmplFile := "templates/endpoint.tmpl"
	tmpl, err := template.ParseFS(assets.TemplatesFS, tmplFile)
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(generatedFile, data)
	if err != nil {
		panic(err)
	}
}

func pluginHandler(plugin *protogen.Plugin) error {
	plugin.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

	for _, f := range plugin.Files {
		if f.Generate == false {
			continue
		}
		generateFile(plugin, f)
	}
	return nil
}

func main() {
	opts := protogen.Options{}
	opts.Run(pluginHandler)
	log.Println("All the standard methods has been generated!")
}
