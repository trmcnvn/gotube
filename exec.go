package main

import (
  "os/exec"
  "runtime"
)

func IsWindows() bool {
  return runtime.GOOS == "windows"
}

func OpenFile(path string) error {
  var cmd *exec.Cmd
  if IsWindows() {
    cmd = exec.Command("start", path)
  } else {
    cmd = exec.Command("open", path)
  }

  err := cmd.Run()
  if err != nil {
    return err
  }

  return nil
}

func FFmpeg(src string, dest string) error {
  // check if ffmpeg is installed
  _, err := exec.LookPath("ffmpeg")
  if err != nil {
    return err
  }

  cmd := exec.Command("ffmpeg", "-i", src, dest, "-y")
  err = cmd.Run()
  if err != nil {
    return err
  }

  return nil
}
