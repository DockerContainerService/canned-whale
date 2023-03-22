package main

import (
	"github.com/DockerContainerService/canned-whale/cmd"
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&nested.Formatter{
		HideKeys:       false,
		NoColors:       false,
		NoFieldsColors: false,
		ShowFullLevel:  false,
		TrimMessages:   false,
		CallerFirst:    true,
		FieldsOrder:    []string{"component", "category"},
	})
	logrus.SetLevel(logrus.InfoLevel)
}

func main() {
	cmd.Execute()
}
