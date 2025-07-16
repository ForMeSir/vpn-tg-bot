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
	Id      int
	Country string
	Type    string
	Online  string
	Key     string
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
	s := gocron.NewScheduler(time.UTC)
	// добавляем одну задачу на каждые 10 минут
	s.Cron("*/10 * * * *").Do(task)
	// запускаем планировщик с блокировкой текущего потока
	s.StartBlocking()
	//task()
}

// func testTask(){
// fmt.Println("Wow")
// }\n\xE2\x8F\xB1Действует: 24ч
func task() {
	fmt.Println("Начало")
	keys := GetKeys()
	c := telegram.New(BotToken)
	if len(keys) != 0 {
		for i := len(keys) - 1; i >= 0; i-- {
			if keys[i].Key != "" {
				err := c.SendMessage("Доступен новый ключ \xF0\x9F\x94\x91 \n\n\xF0\x9F\x93\x8D Локация: "+keys[i].Country+"\n\xF0\x9F\x8C\x90 Трафик: \xE2\x99\xBE \n\n\xF0\x9F\x94\x92 Тип: "+keys[i].Type+"  ```\n"+keys[i].Key+"```", int64(-1002343650923))
				if err != nil {
					fmt.Println("Cлишком много ключей")
					time.Sleep(2 * time.Minute)
				}
				time.Sleep(1 * time.Second)
			}
		}
	}
	fmt.Println("Конец")
}

func GetKeys() []Object {
	var lastID int

	var keys []Object

	file, err := os.Open("cmd/fer.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	data := make([]byte, 64)

	n, err := file.Read(data)
	if err == io.EOF {
		fmt.Println("Ошибка при чтении файла")
	}
	lastID, _ = strconv.Atoi(string(data[:n]))
	if err != nil {
		fmt.Println("Ошибка c получением последнего ID")
	}

	file.Close()

	var Continue bool = true
	for i := 1; len(keys) == 0 || lastID < keys[len(keys)-1].Id; i++ {
		c := colly.NewCollector()
		// set a valid User-Agent header
		c.UserAgent = "Mozilla/5.0 (compatible; MSIE 9.0; Windows; U; Windows NT 6.1; Win64; x64 Trident/5.0)"

		c.OnHTML(".col", func(e *colly.HTMLElement) {
			var key Object
			var err error

			count := e.ChildText(".d-flex.justify-content-between.align-items-center h3")
			IdAndCount := strings.Split(count, " #")
			key.Id, err = strconv.Atoi(IdAndCount[1])
			if err != nil {
				fmt.Println("Ошибка при получении ID")
			}
			key.Country = IdAndCount[0]

			key.Type = e.ChildText(".me-2")

			key.Online = e.ChildText(".text-success")
			if key.Online == "" {
				key.Online = "Offline"
			}
			if key.Id <= lastID && Continue {

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
				if len(keys) == 0 {
					file.WriteString(strconv.Itoa(key.Id))
				} else {
					file.WriteString(strconv.Itoa(keys[0].Id))
				}
				Continue = false
			}
			if Continue {
				keys = append(keys, key)
			}
		})
		num := strconv.Itoa(i)
		c.Visit("https://outlinekeys.com/?page=" + num)
		if !Continue {
			break
		}
	}

	var SendKeys []Object

	for _, value := range keys {
		if value.Country != "Russia" && value.Type != "Vless" && value.Online != "Offline" {
			SendKeys = append(SendKeys, value)
		}
	}
	for num, value := range SendKeys {
		href := "/key/" + strconv.Itoa(value.Id) + "/"
		vpnkey := GetVpnKey(href)
		if len(vpnkey) > 0 {
			SendKeys[num].Key = vpnkey[0]
			fmt.Println(SendKeys[num])
		}
	}
	return SendKeys
}

func GetVpnKey(url string) (surs []string) {
	time.Sleep(1 * time.Second)
	g := colly.NewCollector()
	g.UserAgent = "Mozilla/5.0 (compatible; MSIE 9.0; Windows; U; Windows NT 6.1; Win64; x64 Trident/5.0)"
	g.OnHTML(".form-control", func(e *colly.HTMLElement) {
		sur := e.Text
		surs = append(surs, sur)

	})
	g.Visit("https://outlinekeys.com" + url)
	return surs
}
