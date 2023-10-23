package steam

import (
	"github.com/gdamore/tcell/v2"
	"github.com/navidys/tvxwidgets"
	"github.com/rivo/tview"
	"os"
	"strconv"
	"sync"
	"time"
)

var zan bool

// 进度条
func JinDuTiao(app *tview.Application, musicTime *chan int, pausedchboolJinDu *chan bool, nextMusic chan bool) *tvxwidgets.PercentageModeGauge {
	//app := tview.NewApplication()
	gauge := tvxwidgets.NewPercentageModeGauge()
	gauge.SetTitle("进度条🥵")
	gauge.SetRect(10, 4, 50, 3)
	gauge.SetBorder(true).
		SetTitleAlign(tview.AlignLeft).
		SetBorderColor(tcell.Color51).
		SetTitleColor(tcell.Color49)

	value := 0
	gauge.SetMaxValue(100)

	update := func() {
		for {
		aa:
			var wg sync.Mutex
			var s int

			s = <-*musicTime
			wg.Lock()
			file, _ := os.OpenFile("hello", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			//	fmt.Println(file, error)
			file.WriteString("开始" + strconv.Itoa(s) + "\n")
			wg.Unlock()
			//fmt.Println(s)
			//a := s / 50
			//fmt.Println(s)
			var tick *time.Ticker
			if s/100 > 1 {
				tick = time.NewTicker(time.Duration(s/100) * time.Second)
			} else {
				tick = time.NewTicker(time.Duration(s*1000/100) * time.Millisecond)
			}

			/*	for {
				select {
				case <-tick.C:
					if value > gauge.GetMaxValue() {
						value = 0
					} else {
						value = value + 1
					}
					gauge.SetValue(value)
					app.Draw()
				}
			}*/

			for _ = range tick.C {
				if value > gauge.GetMaxValue() {
					value = 0
					gauge.SetValue(value)
					app.Draw()
					//fmt.Println("ssssssssssssssssssssssss")
					break
				} else {
					//用来控制用了下一首时的进度条
					select {
					case <-nextMusic:
						value = 0
						gauge.SetValue(value)
						app.Draw()
						goto aa
					default:
					}
					//用来控制音乐播放时的进度条（暂停和播放）
					go func() {
						/*for b := range *pausedchboolJinDu {
							wg.Lock()
							file, _ := os.OpenFile("hello", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
							//	fmt.Println(file, error)
							file.WriteString("拿到啦是\t" + strconv.FormatBool(b) + "\n")
							wg.Unlock()
							if !b {
								*zan = false
							} else {
								*zan = true
							}
						}*/

						select {
						case b := <-*pausedchboolJinDu:

							wg.Lock()
							file, _ := os.OpenFile("hello", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
							//	fmt.Println(file, error)
							file.WriteString("拿到啦是\t" + strconv.FormatBool(b) + "\n")
							wg.Unlock()
							if b {
								zan = true
							} else {
								zan = false
							}

						}
					}()

					if !zan {
						wg.Lock()
						/*file, _ := os.OpenFile("hello", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
						//	fmt.Println(file, error)
						file.WriteString("运行时拿到啦是\t" + strconv.FormatBool(zan) + "\n")*/
						time.Sleep(2 * time.Millisecond)
						wg.Unlock()
						value = value + 1
					} else {
						wg.Lock()
						/*file, _ := os.OpenFile("hello", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
						//	fmt.Println(file, error)
						file.WriteString("运行时拿到啦是\t" + strconv.FormatBool(zan) + "\n")*/
						time.Sleep(2 * time.Millisecond)
						wg.Unlock()
					}
					//fmt.Println(value)
				}
				gauge.SetValue(value)
				app.Draw()

			}

		}
	}
	go update()

	return gauge
}
