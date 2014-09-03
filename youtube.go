package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

type YouTube struct {
	output string
	play   bool
	audio  bool
	id     string
}

func (yt *YouTube) ParseURL(url string) error {
	r, err := regexp.Compile(`(?:https?:\/\/)?(?:www\.)?(?:youtube\.com|youtu\.be)\/(?:watch\?v=)?([\w-]+)`)
	if err != nil {
		return err
	}

	matched := r.MatchString(url)
	if matched == false {
		return errors.New("Couldn't parse YouTube URL")
	}

	yt.id = r.FindStringSubmatch(url)[1]
	return nil
}

func (yt *YouTube) GetStreams() ([]url.Values, error) {
	if yt.id == "" {
		return nil, errors.New("YouTube id isn't present")
	}

	res, err := http.Get("https://www.youtube.com/get_video_info?video_id=" + yt.id)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Recieved invalid HTTP status code: %d", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	data, err := url.ParseQuery(string(body))
	if err != nil {
		return nil, err
	}

	// check for good status
	status, found := data["status"]
	if found == false || status[0] != "ok" {
		return nil, errors.New("Server didn't respond with a valid status")
	}

	// set filename if needed
	if yt.output == "" {
		yt.output = data["title"][0]
	}

	stream_map, found := data["url_encoded_fmt_stream_map"]
	if found == false {
		return nil, errors.New("Server didn't respond with a stream map")
	}

	// all available streams for the video
	var ret []url.Values
	streams := strings.Split(stream_map[0], ",")
	for _, stream := range streams {
		data, err := url.ParseQuery(stream)
		if err != nil {
			return nil, err
		}
		ret = append(ret, data)
	}

	return ret, nil
}

func (yt *YouTube) Download() error {
	streams, err := yt.GetStreams()
	if err != nil {
		return err
	}

	// todo: validate and don't assume
	// we assume that the highest quality is the first index
	res, err := http.Get(streams[0]["url"][0])
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("Recieved invalid HTTP status code: %d", res.StatusCode)
	}

	out, err := CreateTmpFile()
	if err != nil {
		return err
	}
	defer func() {
		out.Close()
		os.Remove(out.Name())
	}()

	_, err = io.Copy(out, res.Body)
	if err != nil {
		return err
	}

	// parameterize the title of the file
	filename, err := Parameterize(yt.output)
	if err != nil {
		return err
	}

	workingDir, err := os.Getwd()
	if err != nil {
		return err
	}

	path := workingDir + "/" + filename

	// convert to wav if audio is true
	if yt.audio {
		path = path + ".wav"

		err = FFmpeg(out.Name(), path)
		if err != nil {
			return err
		}
	} else {
		// get the extension based on the stream format
		ext, err := GetExtension(streams[0]["type"][0])
		if err != nil {
			return err
		}
		path = path + "." + ext

		// copy the file from the tmp dir to wd
		f, err := CreateFile(path)
		if err != nil {
			return err
		}
		defer f.Close()

		// reset offset for read
		_, err = out.Seek(0, 0)
		if err != nil {
			return err
		}

		_, err = io.Copy(f, out)
		if err != nil {
			return err
		}
	}

	yt.output = path
	return nil
}
