package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type PageHandler struct {
	RootDir  string
	PageName string
	ICODir   string
	CSSDir   string
	JSDir    string

	page      []byte
	resources map[string][]byte
}

func BuildPageHandler(
	name string,
	root string,
	icoDir string,
	cssDir string,
	jsDir string,
) *PageHandler {
	toReturn := PageHandler{
		PageName:  name,
		RootDir:   root,
		ICODir:    icoDir,
		CSSDir:    cssDir,
		JSDir:     jsDir,
		resources: make(map[string][]byte),
	}
	return &toReturn
}

func (ph *PageHandler) Handle(
	w http.ResponseWriter,
	r *http.Request,
) {
	log.Printf(
		"PageHandler.Handle\n\tpath:%s",
		r.URL.Path,
	)
	name, typ := getResourceNameAndType(r)
	switch typ {
	case Unknown:
		w.WriteHeader(http.StatusNotFound)
		return
	case HTTP:
		page, err := ph.LoadPage()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(page)
		break
	case CSS:
		resource, err := ph.LoadResource(
			name,
			typ,
		)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("content-type", "text/css")
		w.Write(resource)
	default:
		resource, err := ph.LoadResource(
			name,
			typ,
		)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(resource)
		break
	}
}

func (ph *PageHandler) LoadPage() (
	[]byte,
	error,
) {
	if ph.page != nil {
		return ph.page, nil
	}
	var err error
	ph.page, err = LoadFile(
		ph.RootDir,
		ph.PageName,
		"html",
	)
	return ph.page, err
}

func (ph *PageHandler) LoadResource(
	name string,
	typ ResourceType,
) (
	[]byte,
	error,
) {
	if val, ok := ph.resources[name]; ok {
		return val, nil
	}
	var dir string
	switch typ {
	case ICO:
		dir = ph.ICODir
		break
	case CSS, CSSMap:
		dir = ph.CSSDir
		break
	case JS, JSMap:
		dir = ph.JSDir
		break
	}
	var err error
	if dir != "" {
		dir = fmt.Sprintf(
			"%s/%s",
			ph.RootDir,
			dir,
		)
	} else {
		dir = ph.RootDir
	}
	ph.resources[name], err = LoadFile(
		dir,
		name,
		"",
	)
	return ph.resources[name], err
}

func LoadFile(
	baseDir string,
	name string,
	extension string,
) (
	[]byte,
	error,
) {
	var fullPath string
	if extension == "" {
		fullPath = fmt.Sprintf(
			"%s/%s",
			baseDir,
			name,
		)
	} else {
		fullPath = fmt.Sprintf(
			"%s/%s.%s",
			baseDir,
			name,
			extension,
		)
	}
	toReturn, err := ioutil.ReadFile(fullPath)
	if err != nil {
		log.Printf(
			"Error Loading File\n\t:%s",
			err,
		)
		return nil, err
	}
	return toReturn, nil
}

type ResourceType int

const (
	Unknown = ResourceType(iota)
	HTTP
	ICO
	CSS
	CSSMap
	JS
	JSMap
)

func getResourceNameAndType(
	r *http.Request,
) (
	string,
	ResourceType,
) {
	urlSplit := strings.Split(
		r.URL.Path,
		"/",
	)
	resource := urlSplit[len(urlSplit)-1]
	if resource == "" {
		return "", HTTP
	}
	resourceSplit := strings.Split(
		resource,
		".",
	)
	if len(resourceSplit) <= 1 {
		log.Printf(
			"Invalid Resource\n\t:%s",
			resource,
		)
		return "", Unknown
	}
	finalExtension := resourceSplit[len(resourceSplit)-1]
	var typ ResourceType
	switch strings.ToLower(finalExtension) {
	case "js":
		typ = JS
		break
	case "css":
		typ = CSS
		break
	case "ico":
		typ = ICO
		break
	case "map":
		secondaryExtension := resourceSplit[len(resourceSplit)-2]
		if secondaryExtension == "js" {
			typ = JSMap
			break
		} else if secondaryExtension == "css" {
			typ = CSSMap
			break
		}
	default:
		log.Printf(
			"Invalid Resource\n\t:%s",
			resource,
		)
		return "", Unknown
	}
	return resource, typ

}
