package template

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
	"math"
	"net/http"
	"reflect"
	"strings"
)

//go:embed templates/*.html
var templateFiles embed.FS

type TemplateData map[string]any

type Manager interface {
	Render(w http.ResponseWriter, name string, data TemplateData) error
	RenderToString(name string, data TemplateData) (string, error)
}

type TemplateManager struct {
	templates *template.Template
}

func NewManager() (*TemplateManager, error) {
	log.Println("Initializing template manager")
	templateFS, err := fs.Sub(templateFiles, "templates")
	if err != nil {
		log.Printf("Error getting template subdirectory: %v", err)
		return nil, err
	}

	tmpl := template.New("").Funcs(template.FuncMap{
		"add": func(a, b any) float64 {
			return toFloat(a) + toFloat(b)
		},
		"sub": func(a, b any) float64 {
			return toFloat(a) - toFloat(b)
		},
		"mul": func(a, b any) float64 {
			return toFloat(a) * toFloat(b)
		},
		"div": func(a, b any) float64 {
			bVal := toFloat(b)
			if bVal == 0 {
				return 0
			}
			return toFloat(a) / bVal
		},
		"round": func(num any) float64 {
			return math.Round(toFloat(num)*100) / 100
		},
		"replace": func(input, old, new string) string {
			return strings.Replace(input, old, new, -1)
		},
	})

	log.Println("Parsing templates from root directory")
	tmpl, err = tmpl.ParseFS(templateFS, "*.html")
	if err != nil {
		log.Printf("Error parsing templates: %v", err)
		return nil, err
	}

	definedTemplates := tmpl.DefinedTemplates()
	log.Printf("Defined templates: %s", definedTemplates)

	return &TemplateManager{
		templates: tmpl,
	}, nil
}

func (tm *TemplateManager) Render(w http.ResponseWriter, name string, data TemplateData) error {
	log.Printf("Rendering template: %s", name)
	err := tm.templates.ExecuteTemplate(w, name, data)
	if err != nil {
		log.Printf("Error rendering template %s: %v", name, err)
	}
	return err
}

func (tm *TemplateManager) RenderToString(name string, data TemplateData) (string, error) {
	var buf strings.Builder
	log.Printf("Rendering template to string: %s", name)
	err := tm.templates.ExecuteTemplate(&buf, name, data)
	if err != nil {
		log.Printf("Error rendering template to string %s: %v", name, err)
		return "", err
	}
	return buf.String(), nil
}

func toFloat(val any) float64 {
	switch v := val.(type) {
	case int:
		return float64(v)
	case float64:
		return v
	case int32, int64, uint, uint32, uint64:
		return float64(reflect.ValueOf(v).Int())
	case float32:
		return float64(v)
	}
	return 0
}
