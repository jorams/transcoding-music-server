# Transcoding music server

This is a simple, quick-and-dirty HTTP server that transcodes music files
before serving them. You can give it a directory containing music, and when a
music file is requested it is first converted to Opus, using FFmpeg. Other
files are served directly.

This is useful when playing remotely hosted music on a phone with a data plan,
because the transcoded Opus files are significantly smaller than the originals.
The file is converted before responding to the request, so there will be a
small delay while FFmpeg is doing its job. Converted files are saved, so the
next request will be near-instant.

Note that this means a HTTP request for a file called `example.flac` will
return an Opus file, with `Content-Type: audio/ogg`. The transcoded file will
be saved as `example.flac.opus`.

File extensions that will be transcoded are `.flac`, `.mp3` and `.m4a` (not
case sensitive).

## Usage

First build the program:

```sh
go build
```

Now consider a directory `music` that looks like this:

```
music
+ artist
  + example.flac
  + example.mp3
  + example.m4a
  + cover.jpg
```

Run the server:

```sh
transcoding-music-server --origin music --target transcoded-music --bind :8844
```

Note that it calls out to `ffmpeg` to do the transcoding.

Request some files:

```sh
curl -I -X GET http://localhost:8844/artist/example.flac
# HTTP/1.1 200 OK
# Content-Type: audio/ogg
# ...
curl -I -X GET http://localhost:8844/artist/example.mp3
# HTTP/1.1 200 OK
# Content-Type: audio/ogg
# ...
curl -I -X GET http://localhost:8844/artist/example.m4a
# HTTP/1.1 200 OK
# Content-Type: audio/ogg
# ...
curl -I -X GET http://localhost:8844/artist/cover.jpg
# HTTP/1.1 200 OK
# Content-Type: image/jpeg
# ...
```

The directory `transcoded-music` will now look like this:

```
transcoded-music
+ artist
  + example.flac.opus
  + example.mp3.opus
  + example.m4a.opus
```

The music files have been converted, but the cover image was served directly
from the origin directory.

## License

```
Copyright (c) 2020 Joram Schrijver

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```
