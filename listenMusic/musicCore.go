package listenMusic

import (
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"log"
	"math"
	"os"
	"strconv"
	"time"
)

func PlayMusic(path string, pausedchbool *chan bool, musicTimeChan *chan int, nextMusic chan bool) {
	var startStopTime int
	var endStopTime int
	var stopTimeTest int
	var stopAllTime int //一首歌暂停的总时间
	//因为这个函数是使用时间来return,所以上面这几个变量用来控制音乐暂停时，延长return的时间的
	//fmt.Println("start")
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	// 获取音频长度（以秒为单位）
	lengthInSeconds := float64(streamer.Len()) / float64(format.SampleRate)
	var MusicTime int = int(math.Round(lengthInSeconds))

	*musicTimeChan <- MusicTime
	timerStart := int(int64(int(time.Now().Unix())))
	ctrl := &beep.Ctrl{Streamer: beep.Loop(1, streamer), Paused: false}
	speaker.Play(ctrl)
	timer := time.After(time.Duration(MusicTime) * time.Second)
	for {
		select {
		case paused := <-*pausedchbool:
			if paused { //获取开始暂停时的时间
				startStopTime = int(int64(int(time.Now().Unix())))
			} else { //获取结束暂停时的时间
				endStopTime = int(time.Now().Unix())
			}
			stopTimeTest = endStopTime - startStopTime

			if stopTimeTest >= 0 {
				//stopTime = append(stopTime, stopTimeTest)
				//	useTime := startStopTime - timerStart
				//alwaysListen := startStopTime - timerStart
				//stopTime <- stopTimeTest

				stopAllTime = stopAllTime + stopTimeTest

				//useTi:=endStopTime-timerStart
				timer = time.After(time.Duration(timerStart+MusicTime+stopAllTime-endStopTime) * time.Second) //这里相当于重新设置音乐return时间
			} else { //如果stopTimeTest<0，说明只按了暂停，没有按开始,所以一直暂停
				timer = time.After(time.Duration(MusicTime+10000) * time.Second)
			}

			file, _ := os.OpenFile("hello", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			//fmt.Println(file, error)
			file.WriteString(strconv.FormatBool(paused) + "\n")
			speaker.Lock()
			ctrl.Paused = paused
			speaker.Unlock()
		case <-nextMusic:
			file, _ := os.OpenFile("hello", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			//fmt.Println(file, error)
			file.WriteString("下一首开始\n")
			return
		case <-timer:
			file, _ := os.OpenFile("hello", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			//fmt.Println(file, error)
			file.WriteString(path + "   \t(finish)\n")
			return
		}
	}
}
