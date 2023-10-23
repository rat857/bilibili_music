package corn

import (
	"encoding/json"
	"io"
	"net/http"
)

type M4SURL struct {
	Data struct {
		Dash struct {
			Audio []struct {
				BaseUrl string `json:"baseurl"`
			} `json:"audio"`
		} `json:"dash"`
	} `json:"data"`
}

func GetSongm4sURL(avid string, cid string) string {
	url := "https://api.bilibili.com/x/player/playurl?avid=" + avid + "&cid=" + cid + "&qn=0&fnval=80&fnver=0&fourk=1"
	var client http.Client
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Cookie", "_uuid=85931211-63AF-6DE7-1AD3-9FD35103BEFA557562infoc; buvid4=7686E374-7A17-4F44-F11D-0DF892C8DAA758234-023095520-y7FzQHcGAvqnFVXmecmhpw%3D%3D; DedeUserID=699793285; home_feed_column=4; SESSDATA=f5c5ea78%2C1001461506%2C06537%2A92CjDnfSkAVGsW6h2mJnxm1RaA60sKU64rnpw5czIvr0qDSMOueSbZdNoUCEmLyxYWxCQSVjJ5VTlLczRSTGxHU2VORzE1YXgzd0RNQXBOazZySjJWMjBEMTQtQUpwX1hCdHVBUkdsZnV0M2tZQWx1eW5JNzA1NkMwOVFjX2NWbDRJRUpBdUtQWjlBIIEC; bili_jct=f7dcfa1c7802f7424c69619e8ce71a84; sid=4r0nphpo; b_lsid=E10EB41022_18AE1085C06")
	req.Header.Set("User-Agent", " Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36 Edg/117.0.2045.40")
	res, _ := client.Do(req)

	var m4sUrl M4SURL
	body, _ := io.ReadAll(res.Body)
	json.Unmarshal(body, &m4sUrl)
	return m4sUrl.Data.Dash.Audio[0].BaseUrl
}
