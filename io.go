package main

import (
  "io"
  "os"
  "regexp"
  "strings"
  "errors"
)

func Parameterize(str string) (string, error) {
  // Replace unwanted chars into separator char
  r, err := regexp.Compile(`(?i)[^a-z0-9\-_]+`)
  if err != nil {
    return "", err
  }
  str = r.ReplaceAllString(str, "--")

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

  formats := []string{
    "video/mp4",
    "video/x-flv",
  }

  for _, v := range formats {
    if strings.Contains(format, v) {
      return v[len(v)-3:], nil
    }
  }

  return "", errors.New("Couldn't find a valid format")
}

func GetIOStream(yt *YouTube, format string) (io.WriteCloser, error) {
  path, err := os.Getwd()
  if err != nil {
    return nil, err
  }

  filename, err := Parameterize(yt.filename)
  if err != nil {
    return nil, err
  }

  ext, err := GetExtension(format)
  if err != nil {
    return nil, err
  }

  file := path + "/" + filename + "." + ext
  out, err := os.Create(file)
  if err != nil {
    return nil, err
  }

  return out, nil
}
