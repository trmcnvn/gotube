package main

import (
  "os"
  //"os/exec"
  "fmt"
  "flag"
)

func main() {
  output := flag.String("output", "", "specify the output file name")
  play := flag.Bool("play", true, "open the downloaded file in default media player")
  audio := flag.Bool("audio", false, "output audio file")
  flag.Parse()

  // Exit out if there isn't a tailing argument
  if flag.NArg() == 0 {
    fmt.Fprintln(os.Stderr, "You must specify a YouTube URL")
    return
  }

  yt := &YouTube{
      *output,
      *play,
      *audio,
      "",
  }

  err := yt.ParseURL(flag.Arg(0))
  if err != nil {
    fmt.Fprintln(os.Stderr, err)
    return
  }

  fmt.Println("Downloading...")

  // Download the YouTube flv
  err = yt.Download()
  if err != nil {
    fmt.Fprintln(os.Stderr, err)
    return
  }

  fmt.Println("Finished Downloading!")

  // Play the file if set to
  if yt.play {
    // todo: determine if windows
  }
}
