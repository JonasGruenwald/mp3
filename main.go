package main

import (
	"embed"
	"github.com/JonasGruenwald/mp3/cmd"
)

//go:embed templates/default-service.tmpl
var templateFs embed.FS

func main() {
	cmd.Execute(templateFs)
}
