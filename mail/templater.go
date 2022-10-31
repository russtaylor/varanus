package mail

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"html/template"
	"os"
)

type TemplateVars struct {
	FormattedTime    string
	Url              string
	Hostname         string
	FailureType      string
	FailureMessage   string
	ShortFailureType string
}

func CompileTemplate(templateVars TemplateVars) (string, error) {
	// TODO: Look in ./serverless_function_source_code if we can't find the templates by the env var
	sourceDir, present := os.LookupEnv("SOURCE_DIR")
	templatePath := "mail/templates/failure_email.html"
	if present {
		templatePath = fmt.Sprintf("%s/mail/templates/failure_email.html", sourceDir)
	}
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", err
	}
	log.Tracef("Loaded the failure_email template from %s", templatePath)
	buffer := new(bytes.Buffer)
	err = tmpl.Execute(buffer, templateVars)
	if err != nil {
		return "", err
	}
	log.Tracef("Rendered template: %s", buffer.String())
	return buffer.String(), nil
}
