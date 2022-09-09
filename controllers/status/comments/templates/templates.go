package templates

import "embed"

//go:embed *.md
var templatesFs embed.FS

func GetTemplate(name string) (string, error) {
	data, err := templatesFs.ReadFile(name)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func MustGetTemplate(name string) string {
	template, err := GetTemplate(name)
	if err != nil {
		panic(err)
	}
	return template
}
