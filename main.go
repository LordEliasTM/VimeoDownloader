package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/dlclark/regexp2"
)

var wg sync.WaitGroup

func main() {
	args := os.Args[1:]

	url := args[0]
	url = strings.Replace(url, "base64_init=1", "base64_init=0", -1)
	baseUrlMatch, err := regexp2.MustCompile(`^.+(?=sep\/)`, regexp2.None).FindStringMatch(url)
	if err != nil {
		panic(err)
	}
	baseUrl := baseUrlMatch.String()

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var master Master
	json.Unmarshal(body, &master)

	var v1080p Video
	for _, video := range master.Video {
		if video.Height == 1080 {
			v1080p = video
			break
		}
	}
	if v1080p.Duration == 0.0 {
		panic("Fail Video")
	}

	var aMP4 Audio
	for _, audio := range master.Audio {
		if strings.HasPrefix(audio.Codecs, "mp4a") {
			aMP4 = audio
			break
		}
	}
	if v1080p.Duration == 0.0 {
		panic("Fail Audio")
	}

	wg.Add(2)
	var vFileName string
	var aFileName string
	fmt.Println("Downloading video...")
	fmt.Println("Downloading audio...")
	go DownloadMediaSegments(v1080p.Media, baseUrl, &vFileName)
	go DownloadMediaSegments(aMP4.Media, baseUrl, &aFileName)
	wg.Wait()

	fmt.Println("Combining video and audio...")

	exec.Command(filepath.FromSlash("ffmpeg/bin/ffmpeg"),
		"-i", vFileName,
		"-i", aFileName,
		"-c", "copy",
		"-map", "0:v:0",
		"-map", "1:a:0",
		master.ClipId+".mp4",
	).Run()

	os.Remove(vFileName)
	os.Remove(aFileName)

	fmt.Println("Done! Output: ", master.ClipId+".mp4")
}

func DownloadMediaSegments(media Media, baseUrl string, returnVal *string) {
	defer wg.Done()

	videoBaseUrl := baseUrl + "parcel/" + media.BaseUrl
	filename := media.Id + ".mp4"

	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// write init segment
	err = doStuff(file, videoBaseUrl+media.InitSegment)
	if err != nil {
		panic(err)
	}

	// download index segment
	err = doStuff(file, videoBaseUrl+media.IndexSegment)
	if err != nil {
		panic(err)
	}

	for _, segment := range media.Segments {
		err = doStuff(file, videoBaseUrl+segment.Url)
		if err != nil {
			panic(err)
		}
	}

	*returnVal = filename
}

func doStuff(file *os.File, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	file.Write(body)

	return nil
}
