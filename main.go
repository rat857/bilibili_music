// Demo code for the Flex primitive.
package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"main/getYourMusic"
	"main/listenMusic"
	"main/searchGuide"
	"main/steam"
	"os"
	"time"
)

func main() {
	file, _ := os.OpenFile("hello", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	//	fmt.Println(file, error)
	file.WriteString("application well run \n")
	file.Close()

	//获取音乐路径
	folderPath := getYourMusic.GetYourMusicPath()
	//这个管道用来传输输入的搜索的值，用来 SearchInput和SearchEnd之间的值传输
	var searchsring = make(chan string, 2)

	//这个管道用来控制音乐的暂停和播放
	var pausedchbool = make(chan bool, 2)

	//这个管道用来控制进度条的暂停
	var JinDuPausedchbool = make(chan bool, 2)

	//这个管道用来获取歌曲的时长
	var musicTimeInt = make(chan int, 2)

	//这个管道用来控制下一曲
	var nextMusic = make(chan bool, 2)
	//这个管道用来控制下一曲进度条
	var nextMusicJinDu = make(chan bool, 2)
	/*	//这个管道用来获取歌曲的绝对路径
		var filePathStr = make(chan string, 2)*/

	app := tview.NewApplication()

	musicTree, filePathStr := getYourMusic.ShowYourMusic(app, folderPath)
	time.Sleep(2 * time.Second)

	go func() {
		for {
			for s := range filePathStr {

				if len(s) != 0 {
					file, _ := os.OpenFile("hello", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
					//	fmt.Println(file, error)
					file.WriteString(s + "   \t(start)\n")
					JinDuPausedchbool <- false
					listenMusic.PlayMusic(s, &pausedchbool, &musicTimeInt, nextMusic)

				}

			}
		}
	}()
	flex := tview.NewFlex().
		AddItem(musicTree, 0, 1, true).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(steam.JinDuTiao(app, &musicTimeInt, &JinDuPausedchbool, nextMusicJinDu), 3, 1, false).
			AddItem(steam.MusicSteam(app), 0, 3, false).
			AddItem(searchGuide.SearchInput(app, &searchsring), 3, 1, false),
			0, 2, false).
		AddItem(searchGuide.SearchEnd(app, &searchsring), 0, 1, false)
	//控制音乐的暂停和播放
	var paused = false
	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 's' || event.Rune() == 'S' { //暂停
			// 在这里执行你的函数
			//fmt.Println("Space key pressed!")
			paused = !paused
			pausedchbool <- paused
			JinDuPausedchbool <- paused
			//JinDuPausedchbool = append(JinDuPausedchbool, paused)
		} else if event.Rune() == 'd' || event.Rune() == 'D' { //控制下一首
			nextMusic <- true
			nextMusicJinDu <- true
		}
		return event
	})
	if err := app.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}

	file2, _ := os.OpenFile("hello", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	//	fmt.Println(file, error)
	file2.WriteString("\033[1;31;40mRed.\033[0m\n" + "exit, length of the music list well be zero \n")
	file2.Close()
}
