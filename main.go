package main

import (
	"log"
	"net/http"
	"net/http/cookiejar"
	"regexp"
	"strings"
	"time"
)

var (
	CLEAR                    map[string]func()
	CLIENT                   http.Client
	CONFIG                   Config
	MODE, QARGS, QAMT, DEBUG = "currentFaults", "", "", false
)

func main() {
	jar, _ := cookiejar.New(nil)
	CLIENT.Jar, CLIENT.Timeout = jar, time.Second*25
	GetConfig()    // load config file
	Cookie(CLIENT) // get inital cookie from SCADA
	Login(CLIENT)  // login as guest
	RunMain(CLIENT)
}

func RunMain(client http.Client) {
	go BackGroundUpdate(client) // refresh every minute when MODE==currentFaults
	for {
		switch MODE {
		case "historySearch", "histText":
			HistorySearch(client)
			break
		case "currentFaults":
			CurrentFaults(client)
			break
		case "help":
			Help()
		}
	}
}

func BackGroundUpdate(client http.Client) {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	// if in currentFaults refresh
	for _ = range ticker.C {
		if MODE == "currentFaults" {
			go RunMain(client)
		}
	}
}

func Refresh() {
	tmpArgs, tmpAmt := QARGS, QAMT
	QARGS, QAMT = "", ""

	if MODE == "help" {
		_ = GetUserInput("Press Enter to return")
	} else {
		_ = GetUserInput("Enter to refresh, change modes, or enter new search: ")
	}

	switch QARGS {
	case "":
		if MODE == "help" {
			MODE = "currentFaults"
		}
		QARGS = tmpArgs
		QAMT = tmpAmt
		break
	case "h":
		MODE = "help"
		PrintHeader()
		break
	case "t":
		QAMT = ""
		QARGS = ""
		MODE = "histText"
		PrintHeader()
		break
	case "m":
		QAMT = ""
		QARGS = ""
		MODE = "historySearch"
		PrintHeader()
		break
	case "f":
		MODE = "currentFaults"
		PrintHeader()
		QAMT = ""
		break
	case " ebug": // lazy fix for switching in and out of debug mode
	case "debug":
		DEBUG = !DEBUG
		MODE = "currentFaults"
		PrintHeader()
		QAMT = ""
	}
	return
}

func BuildURL(_strings []string) string {
	var queryURL, url string
	// Make a Regex to say we only want letters and numbers
	reg, err := regexp.Compile("[^0-9]+") // regex to strip chars and numbers from querty amt
	if err != nil {
		log.Fatal(err)
	}
	QAMT = reg.ReplaceAllString(QAMT, "")

	switch MODE {
	case "historySearch":
		for idx, str := range _strings {
			if idx == len(_strings)-1 {
				_strings[idx] = "%27" + str + "-PHYSICAL%27"
			} else {
				_strings[idx] = "%27" + str + "-PHYSICAL%27,"
			}
			queryURL += _strings[idx]
		}
		url = CONFIG.Historysearch.URL
		url = strings.Replace(url, "\"PLACEHOLDER1\"", queryURL, 1)
		url = strings.Replace(url, "\"PLACEHOLDER2\"", QAMT, 1)
		break
	case "histText":
		for idx, str := range _strings {
			s := strings.Split(str, "&&")
			for i, v := range s {
				s[i] = "+AND+($%3C%3Ctext%3E%3E+LIKE+%27%25" + v + "%25%27)"
				queryURL += s[i]
			}
			_strings[idx] = str
		}
		//strip out wildcards since they are included in the loop above
		queryURL = strings.ReplaceAll(queryURL, "*", "")
		url = CONFIG.Historytext.URL
		url = strings.Replace(url, "\"PLACEHOLDER1\"", queryURL, 1)
		url = strings.Replace(url, "\"PLACEHOLDER2\"", QAMT, 1)
	}
	QAMT = "" //reset after each query
	return url
}

func FindMatch(str string) []string {
	var _prefix, _suffix, _contains bool = false, false, false

	str = strings.ToUpper(str)
	if _contains = strings.HasPrefix(str, "*") && strings.HasSuffix(str, "*"); _contains == true {
		str = str[1 : len(str)-1]
	} else if _suffix = strings.HasPrefix(str, "*"); _suffix == true {
		str = str[1:]
	} else if _prefix = strings.HasSuffix(str, "*"); _prefix == true {
		str = str[:len(str)-1]
	}

	var match []string
	for _, val := range CONFIG.MachineNames {
		txt := strings.Split(val, " ")
		if _contains {
			if strings.Contains(txt[0], str) {
				match = append(match, txt[0])
			}
		} else if _prefix {
			if strings.HasPrefix(txt[0], str) {
				match = append(match, txt[0])
			}
		} else if _suffix {
			if strings.HasSuffix(txt[0], str) {
				match = append(match, txt[0])
			}
		}
	}

	if _contains == false && _prefix == false && _suffix == false {
		match = append(match, str)
		return match
	} else {
		return match
	}
}

func faultExists(hist []string, str string) bool {
	for _, val := range hist {
		if val == str {
			return true
		}
	}
	return false
}
