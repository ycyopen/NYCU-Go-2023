package main

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// TODO: Please create a struct to include the information of a video

type Video struct {
	Items []struct {
		ID      string `json:"id"`
		Snippet struct {
			PublishedAt  time.Time `json:"publishedAt"`
			Title        string    `json:"title"`
			ChannelTitle string    `json:"channelTitle"`
		} `json:"snippet"`
		Statistics struct {
			ViewCount    string `json:"viewCount"`
			LikeCount    string `json:"likeCount"`
			CommentCount string `json:"commentCount"`
		} `json:"statistics"`
	} `json:"items"`
	PageInfo struct {
		TotalResults int `json:"totalResults"`
	} `json:"pageInfo"`
}

type VideoHTML struct {
	Id           string
	Title        string
	ChannelTitle string
	PublishedAt  string
	ViewCount    string
	LikeCount    string
	CommentCount string
}

func AddComma(s string) string {
	n := len(s)
	if n <= 3 {
		return s
	}
	return AddComma(s[:n-3]) + "," + s[n-3:]
}

func YouTubePage(w http.ResponseWriter, r *http.Request) {
	// TODO: Get API token from .env file
	err_env := godotenv.Load()
	api_key := os.Getenv("YOUTUBE_API_KEY")
	t, err_t := template.ParseFiles("index.html")
	errpage, err_e := template.ParseFiles("error.html")
	if err_env != nil || err_t != nil || err_e != nil {
		errpage.Execute(w, nil)
		return
	}
	// TODO: Get video ID from URL query `v`
	video_id := r.URL.Query().Get("v")
	// TODO: Get video information from YouTube API
	url := "https://www.googleapis.com/youtube/v3/videos?&part=statistics&part=snippet&id=" + video_id + "&key=" + api_key
	resp, err_r := http.Get(url)
	if err_r != nil {
		errpage.Execute(w, nil)
		return
	}
	defer resp.Body.Close()
	body, err_b := io.ReadAll(resp.Body)
	if err_b != nil {
		errpage.Execute(w, nil)
		return
	}
	// TODO: Parse the JSON response and store the information into a struct
	video := Video{}
	err_m := json.Unmarshal([]byte(body), &video)
	if err_m != nil || video.PageInfo.TotalResults == 0 {
		errpage.Execute(w, nil)
		return
	}
	video_output := VideoHTML{
		Id:           video.Items[0].ID,
		Title:        video.Items[0].Snippet.Title,
		ChannelTitle: video.Items[0].Snippet.ChannelTitle,
		PublishedAt:  video.Items[0].Snippet.PublishedAt.Format("2006年01月02日"),
		ViewCount:    AddComma(video.Items[0].Statistics.ViewCount),
		LikeCount:    AddComma(video.Items[0].Statistics.LikeCount),
		CommentCount: AddComma(video.Items[0].Statistics.CommentCount),
	}
	if video_output.Id != video_id {
		errpage.Execute(w, nil)
		return
	}
	// TODO: Display the information in an HTML page through `template`
	t.Execute(w, video_output)
}

func main() {
	http.HandleFunc("/", YouTubePage)
	log.Fatal(http.ListenAndServe(":8085", nil))
}
