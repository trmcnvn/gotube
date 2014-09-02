# gotube

gotube downloads, and opens YouTube media.

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
| -output   | YouTube title | specify the output file name |
| -play     | false         | open the downloaded file in default media player |
| -audio    | false         | output audio file |

## Examples

```

```

## FFmpeg

ffmpeg is required if you wish to download audio only versions of YouTube media.

#### Mac

`brew install ffmpeg`

#### Windows

Check out [http://ffmpeg.zeranoe.com/builds/](http://ffmpeg.zeranoe.com/builds/)

## License

[The MIT License](https://github.com/vevix/gotube/blob/master/LICENSE.md)
