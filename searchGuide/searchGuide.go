package searchGuide

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Tag struct {
	Value string `json:"value"`
}

type Result struct {
	Tag []Tag `json:"tag"`
}

type Response struct {
	Result Result `json:"result"`
}

// 一个输入框，加上了引导输入，按下enter执行的操作在SetDoneFunc设置
func SearchInput(app *tview.Application, input *chan string) *tview.InputField {
	inputField := tview.NewInputField()
	inputField.SetLabelColor(tcell.ColorHotPink).
		SetFieldBackgroundColor(tcell.ColorBlack).
		SetFieldTextColor(tcell.ColorTomato)
	inputField.SetLabel("🔍").
		SetFieldWidth(30)
	inputField.SetBorder(true).SetTitle("搜索框🐕").
		SetBorderColor(tcell.ColorHotPink).SetTitleColor(tcell.ColorLightPink)
	inputField.SetTitleAlign(tview.AlignLeft)
	inputField.SetDoneFunc(func(key tcell.Key) {
		*input <- inputField.GetText()
		//app.Stop()
	})
	// Set up autocomplete function.
	var mutex sync.Mutex
	prefixMap := make(map[string][]string)
	inputField.SetAutocompleteFunc(func(currentText string) []string {
		// Ignore empty text.
		prefix := strings.TrimSpace(strings.ToLower(currentText))
		if prefix == "" {
			return nil
		}

		// Do we have entries for this text already?
		mutex.Lock()
		defer mutex.Unlock()
		entries, ok := prefixMap[prefix]
		if ok {
			return entries
		}

		// No entries yet. Issue a request to the API in a goroutine.
		go func() {
			// Ignore errors in this demo.
			url := "https://s.search.bilibili.com/main/suggest?term=" + url.QueryEscape(prefix)
			res, err := http.Get(url)
			if err != nil {
				return
			}

			// Store the result in the prefix map.
			var companies Response
			dec := json.NewDecoder(res.Body)
			if err := dec.Decode(&companies); err != nil {
				return
			}
			entries := make([]string, 0, len(companies.Result.Tag))
			for _, c := range companies.Result.Tag {
				entries = append(entries, c.Value)
				/*file, _ := os.OpenFile("hello", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
				//	fmt.Println(file, error)
				//a, _ := io.ReadAll(res.Body)
				file.WriteString(c.Value + "   \t(qiao test)\n")*/
			}
			mutex.Lock()
			prefixMap[prefix] = entries
			mutex.Unlock()

			// Trigger an update to the input field.
			inputField.Autocomplete()

			// Also redraw the screen.
			app.Draw()
		}()

		return nil
	})

	return inputField
}
