package gen

import (
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/adamthebaron/orator/config"
	"github.com/adamthebaron/orator/util"
)

func GenerateSite(contentDir, outputDir, staticDir string, fm *util.FrontMatter, layouts map[string]Layout,
	rootTpl *template.Template, conf *config.SiteConfig) error {

	err := filepath.Walk(
		contentDir,
		func(path string, info os.FileInfo, err error) error {
			return makePage(path, info, contentDir, fm, layouts, rootTpl, conf, outputDir, err)
		},
	)
	log.Print("walked filepath")
	if err != nil {
		return err
	}

	err = copyToDir(staticDir, outputDir)
	log.Print("copied to dir")
	if err != nil {
		return err
	}

	log.Print("static files copied successfully")

	return nil
}

func makePage(path string, info os.FileInfo, prefix string, fm *util.FrontMatter,
			  layouts map[string]Layout, rootTpl *template.Template, conf *config.SiteConfig,
			  outputDir string, err error) error {

	if err != nil {
		log.Fatal(err)
	}

	noPrefix := strings.TrimPrefix(path, prefix)
	if info.IsDir() {
		err := os.MkdirAll(filepath.Join(outputDir, noPrefix), os.ModePerm)
		if err != nil {
			return err
		}
	} else {
		p := Page{}
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		generatedPage, err := p.GeneratePage(f, filepath.Ext(path)[1:], conf, rootTpl, fm, layouts)
		if err != nil {
			return err
		}
		pagePath := strings.TrimSuffix(noPrefix, filepath.Ext(noPrefix)) + ".html"
		ioutil.WriteFile(filepath.Join(outputDir, pagePath), []byte(generatedPage), os.ModePerm)

		log.Printf("generated page '%s'", noPrefix)
	}

	return nil
}

func copyToDir(sourceDir, destDir string) error {
	os.MkdirAll(destDir, os.ModePerm)
	os.MkdirAll(sourceDir, os.ModePerm)
	err := filepath.Walk(
		sourceDir,
		func(path string, info os.FileInfo, err error) error {
			dd := strings.TrimPrefix(path, sourceDir)
			dd = filepath.Join(destDir, dd)
			var er error
			if info.IsDir() {
				er = os.MkdirAll(dd, os.ModePerm)
			} else {
				er = copyFile(path, dd)
			}

			return er
		},
	)

	return err
}

func copyFile(source, dest string) error {
	srcFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}
	srcInfo, err := os.Stat(source)
	if err != nil {
		return err
	}
	err = os.Chmod(dest, srcInfo.Mode())
	if err != nil {
		return err
	}
	return nil
}
