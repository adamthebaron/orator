package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"

	yaml "gopkg.in/yaml.v2"

	"github.com/adamthebaron/orator/config"
	"github.com/adamthebaron/orator/gen"
	"github.com/adamthebaron/orator/util"
)

var Fm *util.FrontMatter
var Layouts map[string]gen.Layout
var RootTemplate *template.Template
var SiteConfig *config.SiteConfig
var gendir string
var configFilePath = "config.yaml"
var layoutDir      = "layouts"
var contentDir     = "content"
var staticDir      = "static"

func Init() {
	Fm = util.NewFrontMatter("---")
	Layouts = make(map[string]gen.Layout)
	RootTemplate = template.New("root")
}

func usage() {
	fmt.Print(
		`usage: orator [-h] [-s] [-g gendir]

options:
	-h - print this message
	-s - scaffold a new project into the current directory
	-g - directory to place generated html

usage:
	invoke orator to generate the site in the gen directory in the current working directory.
`,
	)
}

func main() {
	var showUsage, doScaffold bool
	flag.BoolVar(&showUsage, "h", false, "show help")
	flag.BoolVar(&doScaffold, "s", false, "make the required directory structure in this directory")
	flag.StringVar(&gendir, "g", "docs", "directory to place generated html")
	flag.Parse()
	if showUsage {
		usage()
		os.Exit(0)
	}

	if doScaffold {
		scaffold()
		os.Exit(0)
	}

	Init()
	log.Print("init done")
	SiteConfig = new(config.SiteConfig)
	SiteConfig.ReadConfig(configFilePath)
	gen.LoadLayouts(layoutDir, Layouts, RootTemplate, Fm, SiteConfig)
	log.Print("loaded layout")
	err := gen.GenerateSite(contentDir, gendir, staticDir, Fm, Layouts, RootTemplate, SiteConfig)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("job's done.")
}

func scaffold() {
	conf := config.SiteConfig{}
	f, err := os.Create(configFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	out, err := yaml.Marshal(conf)
	f.Write(out)
	os.Mkdir(layoutDir, os.ModePerm)
	os.Mkdir(contentDir, os.ModePerm)
	os.Mkdir(gendir, os.ModePerm)
	os.Mkdir(staticDir, os.ModePerm)
}
