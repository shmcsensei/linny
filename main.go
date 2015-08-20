package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/simonchong/linny/config"

	"github.com/zenazn/goji"
)

var conf config.Config

func init() {

	conf = config.Config{}
	conf.Init()
}

func main() {

	goji.Get("/*", contentController)

	goji.Serve()
}

func getFile(baseDir string, name string) (string, error) {

	// fmt.Println(conf.ContentRoot, baseDir, name)

	absBaseDir, e1 := filepath.Abs(baseDir)
	if e1 != nil {
		return "", e1
	}
	absBaseDir = filepath.Clean(absBaseDir)
	absFile, e2 := filepath.Abs(absBaseDir + "/" + name)
	if e2 != nil {
		return "", e2
	}
	absFile = filepath.Clean(absFile)

	if strings.HasPrefix(absFile, absBaseDir) {
		content, err := ioutil.ReadFile(absFile)
		return string(content), err
	}
	return "", errors.New("Invalid Path :" + absFile + " | " + absBaseDir)
}

func getAdFile(name string) (string, error) {
	adsDir := conf.ContentRoot + "/ad/"
	fmt.Println("Get Ad File: ", adsDir)
	return getFile(adsDir, name)
}

func getResource(name string) (string, error) {
	return getFile(conf.ContentRoot, name)
}

func getWrappedContent(name string) (string, error) {

	content, err0 := getAdFile(name)
	if err0 != nil {
		return "", err0
	}
	header, err1 := getResource("header.frag")
	if err1 != nil {
		return "", err1
	}
	footer, err2 := getResource("footer.frag")
	if err2 != nil {
		return "", err2
	}

	rtn := header + content + footer

	return rtn, nil
}

func contentController(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Requested: ", r.URL.Path[1:])

	contentRequested := r.URL.Path[1:]

	content, err := getWrappedContent(contentRequested)

	if err != nil {
		fmt.Println("Content Controller Error: ", err)
		http.NotFound(w, r)
		return
	}

	w.Header().Set(
		"Content-Type",
		"text/html",
	)
	fmt.Fprint(w, content)
}
