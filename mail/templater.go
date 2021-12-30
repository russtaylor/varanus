package mail

import (
	"bytes"
	"html/template"
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
	tmpl, err := template.ParseFiles("mail/templates/failure_email.html")
	if err != nil {
		return "", err
	}
	buffer := new(bytes.Buffer)
	err = tmpl.Execute(buffer, templateVars)
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}
