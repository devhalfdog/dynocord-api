package twitch

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/devhalfdog/dynocord-api/utils"
	"github.com/grafov/m3u8"
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/mowshon/moviego"
	"github.com/wmw64/twitchpl"
)

const (
	USER_AGENT = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/115.0"
)

var (
	TWITCH_CLIENT = utils.Environment("STREAM_CLIENT")
	TWITCH_TOKEN  = utils.Environment("STREAM_TOKEN")
	client        = http.Client{
		Timeout: 1 * time.Second,
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     1 * time.Second,
		},
	}
)

func GetStreamScreenShot(channel string) (string, error) {
	pl, err := twitchpl.Get(context.Background(), channel, false)
	if err != nil {
		return "", err
	}

	url := pl.Best().AsURL()

	now, err := timeIn(time.Now(), "Local")
	if err != nil {
		log.Println(err)
	}

	path, err := createChannelDir(channel, now.Format("15-04-05"))
	if err != nil {
		return "", err
	}

	dlc := make(chan *Segment, 1)
	defer close(dlc)

	go getPlaylist(channel, url, dlc)
	go downloadSegment(path, dlc)

	time.Sleep(1 * time.Second)

	video, err := moviego.Load(path)
	if err != nil {
		return "", err
	}

	screenshot, err := video.Screenshot(1, fmt.Sprintf("./static/%s.png", channel))
	if err != nil {
		return "", err
	}

	// file remove
	err = os.Remove(path)
	if err != nil {
		log.Println("file remove error :", err)
		return screenshot, nil
	}

	return screenshot, nil
}

func getPlaylist(channel string, urlStr string, dlc chan *Segment) {
	var recDuration time.Duration = 0
	var req *http.Request

	cache, _ := lru.New[any, any](1)
	defer cache.Purge()

	plUrl, err := url.Parse(urlStr)
	if err != nil {
		log.Println(err)
	}

	req, err = http.NewRequest("GET", urlStr, nil)
	if err != nil {
		log.Println()
	}

	res, err := doRequest(req)
	if err != nil {
		log.Println(err)
		return
	}

	playlist, listType, err := m3u8.DecodeFrom(res.Body, true)
	if err != nil {
		log.Println(err)
		return
	}
	res.Body.Close()

	if listType == m3u8.MEDIA {
		uri := playlist.(*m3u8.MediaPlaylist).Segments[0].URI
		var msURI string

		if strings.HasPrefix(uri, "http") {
			msURI, err = url.QueryUnescape(uri)
			if err != nil {
				log.Println(err)
			}
		} else {
			msUrl, err := plUrl.Parse(uri)
			if err != nil {
				log.Println(err)
			}

			msURI, err = url.QueryUnescape(msUrl.String())
			if err != nil {
				log.Println(err)
			}
		}
		_, hit := cache.Get(msURI)
		if !hit {
			cache.Add(msURI, nil)
			dlc <- &Segment{msURI, recDuration}
		}
	}
}

func downloadSegment(filepath string, dlc chan *Segment) {
	var req *http.Request

	out, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		log.Println(err)
	}
	defer out.Close()

	segment := <-dlc
	req, err = http.NewRequest("GET", segment.URI, nil)
	if err != nil {
		log.Println(err)
	}

	res, err := doRequest(req)
	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		res.Body.Close()
	}

	_, err = io.Copy(out, res.Body)
	if err != nil {
		log.Println(err)
	}
}

func doRequest(req *http.Request) (*http.Response, error) {
	req.Header.Add("User-Agent", USER_AGENT)
	req.Header.Add("Client-ID", TWITCH_CLIENT)
	req.Header.Add("Authorization", fmt.Sprintf("OAuth %s", TWITCH_TOKEN))

	res, err := client.Do(req)

	return res, err
}

func createChannelDir(channel string, time string) (string, error) {
	sep := string(os.PathSeparator)
	fName := fmt.Sprintf("%v_%v.ts", channel, time)
	path, _ := os.Getwd()
	cDir := path + sep + "stream"

	if err := os.MkdirAll(cDir, 0777); err != nil {
		return "", err
	}

	filepath := cDir + sep + fName

	return filepath, nil
}

func timeIn(t time.Time, name string) (time.Time, error) {
	loc, err := time.LoadLocation(name)
	if err == nil {
		t = t.In(loc)
	}

	return t, err
}
