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

	out, err := GetIOStream(yt, streams[0]["type"][0])
	if err != nil {
		return err
	}

	_, err = io.Copy(out, res.Body)
	if err != nil {
		return err
	}

	// convert to wav if audio is true
	if yt.audio {
		wav_path := yt.output[:len(yt.output)-3]
		wav_path = wav_path + "wav"

		err = FFmpeg(yt.output, wav_path)
		if err != nil {
			return err
		}

		// delete the video file
		vid_path := yt.output
		defer func() {
			// Close Writer otherwise Windows won't delete the file
			out.Close()
			os.Remove(vid_path)
		}()

		// update output for .wav file
		yt.output = wav_path
	} else {
		out.Close()
	}

	return nil
}
