package main

import (

	// import Colly

	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	telegram "web-scraper/tg-bot"

	"github.com/go-co-op/gocron"
	"github.com/gocolly/colly"
	"github.com/joho/godotenv"
)

type Object struct {
	Id        int
	Country   string
	CreatedAt string
	//Time      string
	Type   string
	Online string
	Key    string
}

var BotToken string

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	BotToken = os.Getenv("BOT_TOKEN")
	fmt.Println(BotToken)
	//инициализируем объект планировщика
	s := gocron.NewScheduler(time.UTC)
	// добавляем одну задачу на каждые 10 минут
	s.Cron("*/10 * * * *").Do(task)
	// запускаем планировщик с блокировкой текущего потока
	s.StartBlocking()
}

// func testTask(){
// fmt.Println("Wow")
// }\n\xE2\x8F\xB1Действует: 24ч
func task() {
	fmt.Println("Начало")
	keys := GetKeys()
	c := telegram.New(BotToken)
	for _, value := range keys {
		err := c.SendMessage("Доступен новый ключ \xF0\x9F\x94\x91 \n\n\xF0\x9F\x93\x8D Локация: "+value.Country+"\n\xF0\x9F\x8C\x90 Трафик: \xE2\x99\xBE \n\n\xF0\x9F\x94\x92 Тип: "+value.Type+"  ```\n"+value.Key+"```", int64(-1002343650923))
		if err != nil {
			fmt.Println(err)
		}
	}
	fmt.Println("Конец")
}

func GetKeys() []Object {
	var lastID int

	c := colly.NewCollector()
	// set a valid User-Agent header

	c.UserAgent = "Mozilla/5.0 (compatible; MSIE 9.0; Windows; U; Windows NT 6.1; Win64; x64 Trident/5.0)"
	var der Object
	var keys []Object

	file, err := os.Open("cmd/fer.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	data := make([]byte, 64)

	for {
		n, err := file.Read(data)
		if err == io.EOF {
			break
		}
		lastID, _ = strconv.Atoi(string(data[:n]))
		if err != nil {
			fmt.Println("Ошибка c получением последнего ID")
		}
	}

	file.Close()

	var jart bool = true
	der.Id = lastID + 1
	keys = append(keys, der)
	for i := 1; lastID < keys[len(keys)-1].Id; i++ {

		c.OnHTML(".col", func(e *colly.HTMLElement) {
			count := e.ChildText(".d-flex.justify-content-between.align-items-center h3")
			var key Object
			var err error
			IdAndCount := strings.Split(count, "#")
			key.Id, err = strconv.Atoi(IdAndCount[1])
			if err != nil {
				fmt.Println("Ошибка при получении ID")
			}
			key.Country = IdAndCount[0]

			// count = e.ChildText(".card-footer.text-muted.d-flex.justify-content-between.align-items-center .d-flex.align-items-center")
			// times := strings.Split(count, "ago")
			// key.Time = times[0]

			key.Type = e.ChildText(".me-2")
			key.Online = e.ChildText(".text-success")
			if key.Online == "" {
				key.Online = "Offline"
			}
			if key.Id <= lastID {
				if len(keys) > 1 {
					err := os.Remove("cmd/fer.txt")
					if err != nil {
						fmt.Println("Ошибка при удалении файла:", err)
					}

					file, err := os.Create("cmd/fer.txt")

					if err != nil {
						fmt.Println("Unable to create file:", err)
						os.Exit(1)
					}
					defer file.Close()

					file.WriteString(strconv.Itoa(keys[1].Id))
				}
				jart = false
			}
			if jart {
				keys = append(keys, key)
			}
		})
		num := strconv.Itoa(i)
		c.Visit("https://outlinekeys.com/?page=" + num)
		if lastID < keys[len(keys)-1].Id {
			break
		}
	}
	var SendKeys []Object
	for key, value := range keys {
		if value.Country == "Russia" || value.Type == "Vless" || value.Online == "Offline" {
			continue
		}
		if key == 0 {
			continue
		}
		SendKeys = append(SendKeys, value)
	}
	for num, value := range SendKeys {
		href := "/key/" + strconv.Itoa(value.Id) + "/"
		vpnkey := GetVpnKey(href)
		SendKeys[num].Key = vpnkey[0]
		fmt.Println(SendKeys[num])
	}
	return SendKeys
}

func GetVpnKey(url string) (surs []string) {
	time.Sleep(2 * time.Second)
	g := colly.NewCollector()
	g.UserAgent = "Mozilla/5.0 (compatible; MSIE 9.0; Windows; U; Windows NT 6.1; Win64; x64 Trident/5.0)"
	g.OnHTML(".form-control", func(e *colly.HTMLElement) {
		sur := e.Text
		surs = append(surs, sur)

	})
	g.Visit("https://outlinekeys.com" + url)
	return surs
}
