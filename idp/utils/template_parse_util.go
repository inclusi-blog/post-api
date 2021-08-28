package util

import (
	"bytes"
	"context"
	"github.com/gola-glitch/gola-utils/logging"
	"html/template"
)

func ParseTemplate(goContext context.Context, templatePath string, data interface{}) (string, error) {
	logger := logging.GetLogger(goContext)
	content := new(bytes.Buffer)
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		logger.Error("Could not find email content template : ", err)
		return "", err
	}
	err = t.Execute(content, data)
	if err != nil {
		logger.Error("Could not generate email content from template : ", err)
		return "", err
	}
	return content.String(), nil
}
