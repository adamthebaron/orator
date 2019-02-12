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

var fm *util.FrontMatter
var layouts map[string]gen.Layout
var roottemplate *template.Template
var siteconf *config.SiteConf
var gendir string
var sitedir string
var confpath = "config.yaml"
var layoutDir = "layouts"
var contentDir = "content"
var staticDir = "static"

func Init() {
	fm = util.NewFrontMatter("---")
	layouts = make(map[string]gen.Layout)
	roottemplate = template.New("root")
}

func usage() {
	fmt.Print(
		`usage: orator [-h] [-s] [-g gendir] [-d sitedir]

options:
	-h - print this message
	-s - scaffold a new project into the current directory
	-g - directory to place generated html
	-s - directory containing stuff to generate (default is cwd)
`,
	)
}

func main() {
	var showUsage, doScaffold bool
	flag.BoolVar(&showUsage, "h", false, "show help")
	flag.BoolVar(&doScaffold, "s", false, "make the required directory structure in this directory")
	flag.StringVar(&gendir, "g", "docs", "directory to place generated html")
	flag.StringVar(&sitedir, "s", ".", "directory containing stuff to generate")
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
	siteconf = new(config.SiteConf)
	siteconf.ReadConfig(confpath)
	gen.Loadlayouts(layoutDir, layouts, roottemplate, fm, siteconf)
	log.Print("loaded layout")
	err := gen.GenerateSite(contentDir, gendir, staticDir, fm, layouts, roottemplate, siteconf)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("job's done.")
}

func scaffold() {
	conf := config.SiteConf{}
	f, err := os.Create(confpath)
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
