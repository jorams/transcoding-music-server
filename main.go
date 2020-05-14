package main

import (
	"errors"
	"flag"
	"github.com/gorilla/handlers"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

type TranscodedDirectory struct {
	Origin string
	Target string
}

func (d TranscodedDirectory) Open(name string) (http.File, error) {
	if filepath.Separator != '/' && strings.ContainsRune(name, filepath.Separator) {
		return nil, errors.New("invalid character in file path")
	}

	originName := filepath.Join(d.Origin, filepath.FromSlash(path.Clean("/"+name)))
	lowerOriginName := strings.ToLower(originName)

	if strings.HasSuffix(lowerOriginName, ".flac") || strings.HasSuffix(lowerOriginName, ".mp3") || strings.HasSuffix(lowerOriginName, ".m4a") {
		targetName := filepath.Join(d.Target, filepath.FromSlash(path.Clean("/"+name))) + ".opus"
		targetFile, targetErr := os.Open(targetName)
		if targetErr == nil {
			return targetFile, nil
		}

		if _, originErr := os.Stat(originName); originErr != nil {
			return nil, originErr
		}

		baseName := filepath.Dir(targetName)
		if dirErr := os.MkdirAll(baseName, 0755); dirErr != nil {
			return nil, dirErr
		}

		ffmpeg := exec.Command("ffmpeg", "-i", originName, "-b:a", "80K", targetName)
		ffmpegErr := ffmpeg.Run()
		if ffmpegErr != nil {
			return nil, ffmpegErr
		}

		targetFile, targetErr = os.Open(targetName)
		if targetErr != nil {
			return nil, targetErr
		}
		return targetFile, nil
	}

	originFile, originErr := os.Open(originName)
	return originFile, originErr
}

func main() {
	var bind string
	var origin string
	var target string

	flag.StringVar(&bind, "bind", ":8844", "port to serve files on")
	flag.StringVar(&origin, "origin", "origin", "origin directory")
	flag.StringVar(&target, "target", "target", "target directory to store transcoded files in")
	flag.Parse()

	fileServer := http.FileServer(TranscodedDirectory{
		Origin: origin,
		Target: target,
	})

	server := http.Server{
		Addr:    bind,
		Handler: handlers.LoggingHandler(os.Stdout, fileServer),
	}

	log.Printf("Starting transcoding server on %s", bind)
	log.Fatal(server.ListenAndServe())
}
