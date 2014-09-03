# gotube

gotube allows you to download YouTube media in either video, or audio format.

Current limitations: Will only download highest quality version available.

## Install

```
go get https://github.com/vevix/gotube
```

## Usage

```
gotube <youtube-url>
```

#### Options

| Parameter | Default Value | Description |
|-----------|---------------|-------------|
| -output   | YouTube title | specify the output file name (extension is added by gotube) |
| -play     | false         | open the downloaded file in default media player |
| -audio    | false         | output audio file |

## Examples

```
gotube -audio=true -play=true https://www.youtube.com/watch?v=YE7VzlLtp-4
// => big-buck-bunny.wav

gotube youtu.be/YE7VzlLtp-4
// => big-buck-bunny.mp4
```

## FFmpeg

ffmpeg is required if you wish to download audio only versions of YouTube media.

#### Mac

`brew install ffmpeg`

## License

[The MIT License](https://github.com/vevix/gotube/blob/master/LICENSE.md)
