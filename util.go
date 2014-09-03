package main

import (
	"errors"
	"io/ioutil"
	"regexp"
	"strings"
	"os"
)

var formats []string = []string{
	"video/mp4",
	"video/x-flv",
}

func Parameterize(str string) (string, error) {
	// Replace unwanted chars into separator char
	r, err := regexp.Compile(`(?i)[^a-z0-9\-_]+`)
	if err != nil {
		return "", err
	}
	str = r.ReplaceAllString(str, "-")

	// Replace any occurancs of more than one separator
	r, err = regexp.Compile(`\-{2,}`)
	if err != nil {
		return "", err
	}
	str = r.ReplaceAllString(str, "-")

	// Remove leading/trailing separator
	r, err = regexp.Compile(`(?i)^\-|\-$`)
	if err != nil {
		return "", err
	}
	str = r.ReplaceAllString(str, "")

	return strings.ToLower(str), nil
}

func GetExtension(format string) (string, error) {
	if format == "" {
		return "", errors.New("Format can't be blank")
	}

	for _, v := range formats {
		if strings.Contains(format, v) {
			return v[len(v)-3:], nil
		}
	}

	return "", errors.New("Couldn't find a valid format")
}

func CreateTmpFile() (*os.File, error) {
	f, err := ioutil.TempFile("", "gotube")
	if err != nil {
		return nil, err
	}

	return f, nil
}

func CreateFile(path string) (*os.File, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	return f, nil
}
