package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/yanzay/tbot/v2"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
		log.Print(err)
	}
}

type stateInfo struct {
	state string
	phone string
}

var mapStates = map[string]stateInfo{
	"IN": {"India", "+91-11-23978046"},
	"MH": {"Maharashtra", "020-26127394"},
	"KL": {"Kerala", "0471-2552056"},
	"TN": {"Tamil Nadu", "044-29510500"},
	"UP": {"Uttar Pradesh", "18001805145"},
	"KA": {"Karnataka", "104"},
	"DL": {"Delhi", "011-22307145"},
	"RJ": {"Rajasthan", "0141-2225624"},
	"TG": {"Telangana", "104"},
	"GJ": {"Gujarat", "104"},
	"MP": {"Madhya Pradesh", "104"},
	"JK": {"Jammu and Kashmir", "0194-2440283"},
	"HR": {"Haryana", "8558893911"},
	"PB": {"Punjab", "104"},
	"AP": {"Andhra Pradesh", "0866-2410978"},
	"WB": {"West Bengal", "1800313444222"},
	"BR": {"Bihar", "104"},
	"LK": {"Ladakh", "01982256462"},
	"CH": {"Chandigarh", "9779558282"},
	"AN": {"Andaman and Nicobar Islands", "03192-232102"},
	"CT": {"Chhattisgarh", "104"},
	"UT": {"Uttarakhand", "104"},
	"GA": {"Goa", "104"},
	"HP": {"Himachal Pradesh", "104"},
	"OR": {"Odisha", "9439994859"},
	"MN": {"Manipur", "3852411668"},
	"MZ": {"Mizoram", "102"},
	"PY": {"Puducherry", "104"},
	"AS": {"Assam", "6913347770"},
	"JH": {"Jharkhand", "104"},
	"AR": {"Arunachal Pradesh", "9436055743"},
	"DN": {"Dadra and Nagar Haveli", "104"},
	"DD": {"Daman and Diu", "104"},
	"LD": {"Lakshadweep", "104"},
	"ML": {"Meghalaya", "108"},
	"NL": {"Nagaland", "7005539653"},
	"SK": {"Sikkim", "104"},
	"TR": {"Tripura", "0381-2315879"},
}

//StatsDel function
type StatsDel struct {
	ActiveDel    string `json:"delta{active}"`
	ConfirmedDel string `json:"delta{confirmed}"`
	DeathsDel    string `json:"delta{deaths}"`
	RecoveredDel string `json:"delta{recovered}"`
}

//StatsAll function
type StatsAll struct {
	ActiveNow    string `json:"active"`
	ConfirmedNow string `json:"confirmed"`
	DeathsNow    string `json:"deaths"`
	RecoveredNow string `json:"recovered"`
	State        string `json:"state"`
	LastUpdTime  string `json:"lastupdatedtime"`
	// ListDel      []StatsDel `json:"delta"`
}

//StateWise function
type StateWise struct {
	ListAll []StatsAll `json:"statewise"`
}

func callTheAPI() (StateWise, error) {
	var err error = nil
	res, err2 := http.Get("https://api.covid19india.org/data.json")
	if err2 != nil {
		err = errors.New("Couldn't connect to API")
	}
	// fmt.Prstringln(reflect.TypeOf(res))
	text, err1 := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err1 != nil {
		err = errors.New("Couldn't Read Body")
	}

	theStats := StateWise{}
	jsonErr := json.Unmarshal(text, &theStats)
	if jsonErr != nil {
		err = errors.New("Couldn't Unmarshall the JSON")
	}
	return theStats, err
}

func getThemStats(state string) (string, string) {
	var msg string
	var phnMsg string
	var i int
	theStats, err1 := callTheAPI()
	if err1 != nil {
		msg = err1.Error()
		return msg, msg
	}
	var initMsg string
	if state == "IN" {
		initMsg = "Across the country,"
	} else {
		initMsg = "In the state of " + mapStates[state].state + ","
	}
	for i = 0; i < len(theStats.ListAll); i++ {
		// msg = msg + theStats.ListAll[i].State + " " + mapStates[state] + "\n"
		if theStats.ListAll[i].State == mapStates[state].state {
			phnMsg = "For any emergency, contact: " + mapStates[state].phone
			res := strings.Split(theStats.ListAll[i].LastUpdTime, " ")
			date := res[0]
			time := res[1]
			msg = fmt.Sprint(initMsg, " as of ", time, ", on ", date, ", there have been ", theStats.ListAll[i].ActiveNow, " cases reported, of which ", theStats.ListAll[i].DeathsNow, " people have been reported to be dead, with ", theStats.ListAll[i].RecoveredNow, " having recoved from the Novel COVID-19.")
		}
	}
	if i > 38 {
		msg += "Invalid Command Entered. " + state + initMsg
	}
	return msg, phnMsg
}

func main() {
	var codeTxt string = ""
	for key, value := range mapStates {
		codeTxt += key + " - " + value.state + "\n"
	}

	token, _ := os.LookupEnv("TOKEN")
	bot := tbot.New((token), tbot.WithWebhook("https://covidindia-bot.herokuapp.com", ":"+os.Getenv("PORT")))
	cli := bot.Client()
	bot.HandleMessage("/stats .+", func(m *tbot.Message) {
		text := strings.TrimPrefix(m.Text, "/stats ")
		theText, emerMsg := getThemStats(text)
		fmt.Println(m.Text)
		cli.SendMessage(m.Chat.ID, theText)
		cli.SendMessage(m.Chat.ID, emerMsg)
	})
	bot.HandleMessage("/about", func(m *tbot.Message) {
		aboutText := "Made with <3 by Gopikrishnan K \n\n Catch me here : github.com/SuperGops7 \n\n Powered by the wonderful people at https://api.covid19india.org/"
		cli.SendMessage(m.Chat.ID, aboutText)
	})
	bot.HandleMessage("/help", func(m *tbot.Message) {
		cli.SendMessage(m.Chat.ID, "The states/UT and their codes are as follows:")
		cli.SendMessage(m.Chat.ID, codeTxt)
	})
	log.Fatal(bot.Start())

	// fmt.Printf("%s", theStats.ListAll[6].State)
}
