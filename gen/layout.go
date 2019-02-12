package gen

import (
	"html/template"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/adamthebaron/orator/config"
	"github.com/adamthebaron/orator/util"
)

// layout metadata
type layout struct {
	name        string
	frontmatter map[string]interface{}
}

// load layouts from the layouts directory
func LoadLayouts(dirpath string, layouts map[string]Layout, rootTemplate *template.Template, fm *util.FrontMatter, conf *config.SiteConfig) {
	files, err := ioutil.ReadDir(dirpath)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if !file.IsDir() {
			// relative pathname from the project root
			fileName := dirpath + "/" + file.Name()
			// template's name is file's base name without suffix
			templateName := strings.TrimSuffix(file.Name(), filepath.Ext(fileName))
			buf, err := ioutil.ReadFile(fileName)
			if err != nil {
				log.Fatal(err)
			}
			fileContents := string(buf)
			front, body, err := fm.Parse(strings.NewReader(fileContents))
			if err != nil {
				log.Fatal(err)
			}
			lfm := make(map[string]interface{})
			yaml.Unmarshal([]byte(front), &lfm)
			_, err = rootTemplate.Parse(body)
			if err != nil {
				log.Fatal("error generating template '%s' '%s'", tempateName, err)
			}
			layouts[templateName] = Layout{templateName, lfm}
			log.Printf("loaded layout '%s'.", templateName)
		}
	}
}
