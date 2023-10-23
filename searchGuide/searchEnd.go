package searchGuide

import (
	"encoding/json"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"io"
	"main/corn"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type SearchResult struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Result []struct {
			ResultType string `json:"result_type"`
			Data       []struct {
				Type        string `json:"type"`
				Aid         int    `json:"aid"`
				Title       string `json:"title"`
				Description string `json:"description"`
			} `json:"data"`
		} `json:"result"`
	} `json:"data"`
}

func SearchEnd(app *tview.Application, inputStringChan *chan string) *tview.List {
	list := tview.NewList()
	go func() {
		for text := range *inputStringChan {
			whowmany := list.GetItemCount()
			for i := 0; i < whowmany; i++ {
				list.RemoveItem(i)
			}
			url := "https://api.bilibili.com/x/web-interface/search/all/v2?page=1&keyword=" + text
			var client http.Client
			req, _ := http.NewRequest("GET", url, nil)
			req.Header.Set("Cookie", "_uuid=85931211-63AF-6DE7-1AD3-9FD35103BEFA557562infoc; buvid4=7686E374-7A17-4F44-F11D-0DF892C8DAA758234-023095520-y7FzQHcGAvqnFVXmecmhpw%3D%3D; DedeUserID=699793285; home_feed_column=4; SESSDATA=f5c5ea78%2C1001461506%2C06537%2A92CjDnfSkAVGsW6h2mJnxm1RaA60sKU64rnpw5czIvr0qDSMOueSbZdNoUCEmLyxYWxCQSVjJ5VTlLczRSTGxHU2VORzE1YXgzd0RNQXBOazZySjJWMjBEMTQtQUpwX1hCdHVBUkdsZnV0M2tZQWx1eW5JNzA1NkMwOVFjX2NWbDRJRUpBdUtQWjlBIIEC; bili_jct=f7dcfa1c7802f7424c69619e8ce71a84; sid=4r0nphpo; b_lsid=E10EB41022_18AE1085C06")
			req.Header.Set("User-Agent", " Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36 Edg/117.0.2045.40")
			res, _ := client.Do(req)

			var searchResult SearchResult
			body, _ := io.ReadAll(res.Body)
			//fmt.Println(string(body))
			json.Unmarshal(body, &searchResult)
			for _, result := range searchResult.Data.Result {
				if result.ResultType == "video" {
					for _, datum := range result.Data {
						//fmt.Printf("title:%s", datum.Title)
						//time.Sleep(1 * time.Second)
						app.QueueUpdateDraw(func() {
							//fmt.Println(datum.Title)
							time.Sleep(500 * time.Millisecond)
							title := datum.Title + "~Aid~=" + strconv.Itoa(datum.Aid)
							list.AddItem(title, datum.Description, '👻', func() {
								aid := strings.Split(title, "~Aid~=")
								cid, title := corn.GetCid(aid[1])
								m4sURL := corn.GetSongm4sURL(aid[1], cid)
								cmd := exec.Command("/bin/bash", "-c", "wget '"+m4sURL+"' --referer 'https://www.bilibili.com' -O '"+aid[1]+".m4s'")
								cmd.Run()
								cm2 := fmt.Sprintf(`ffmpeg -i %s.m4s "%s.mp3"`, aid[1], title)
								//fmt.Println(cm2)
								cmd2 := exec.Command("/bin/bash", "-c", cm2)
								cmd2.Run()
								cm3 := fmt.Sprintf(`mv "%s.mp3" ~/Music/bili_music/`, title)
								cmd3 := exec.Command("/bin/bash", "-c", cm3)
								cmd3.Run()
								//fmt.Println(string(out))
								//fmt.Println("finish")
								//fmt.Printf("aid=%s,cid=%s,url=%s", aid[1], cid, m4sURL)
							})
						})
					}
				}
			}
		}

	}()
	list.SetBorder(true).
		SetBorderColor(tcell.ColorBlue).
		SetTitle("搜索结果：").
		SetTitleAlign(tview.AlignLeft).
		SetTitleColor(tcell.ColorDeepSkyBlue)
	return list
}
