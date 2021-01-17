package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
)

func Cookie(client http.Client) {
	req, err := http.NewRequest(CONFIG.Cookie.Method, CONFIG.Cookie.URL, nil)
	// set the appropriate headers before making request
	SetHeaders(*req, "cookie")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	client.Jar.SetCookies(req.URL, resp.Cookies())
	defer resp.Body.Close()
}

func Login(client http.Client) {
	data := Payload{CONFIG.Login.Auth.Un, CONFIG.Login.Auth.Pw}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}

	body := bytes.NewReader(payloadBytes)
	req, err := http.NewRequest(CONFIG.Login.Method, CONFIG.Login.URL, body)
	if err != nil {
		fmt.Println(err)
	}

	// set the appropriate headers before making request
	SetHeaders(*req, "login")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	client.Jar.SetCookies(req.URL, resp.Cookies())
	defer resp.Body.Close()
}

func CurrentFaults(client http.Client) {
	PrintHeader()
	req, err := http.NewRequest(CONFIG.Currentfaults.Method, CONFIG.Currentfaults.URL, nil)
	if err != nil {
		fmt.Println(err)
	}
	// set the appropriate headers before making request
	SetHeaders(*req, "currentFaults")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		main()
	}
	defer resp.Body.Close()

	faults := new(lhJSON)
	// write the resp.body json to faults (lhJSON)
	err = json.NewDecoder(resp.Body).Decode(&faults)
	if err != nil {
		fmt.Println(err)
	}

	RenderTable(faults)
	Refresh()
}

func HistorySearch(client http.Client) {
	var URL string
	var unit []string

	PrintHeader()

	for QARGS == "" {
		str := GetUserInput("Enter new search or change modes: ")
		args := strings.Fields(str)
		if len(args) > 1 {
			QAMT = args[1]
		}
		if len(args) > 0 {
			QARGS = args[0]
			if QARGS == "f" {
				MODE = "currentFaults"
				QARGS = ""
				QAMT = ""
				return
			} else if QARGS == "h" {
				MODE = "help"
				QARGS = ""
				QAMT = ""
				return
			}
			break
		}
	}

	if MODE == "historySearch" {
		unit = FindMatch(QARGS)
		if QAMT == "" {
			QAMT = "15"
		}
		URL = BuildURL(unit)
	} else if MODE == "histText" {
		if QAMT == "" {
			QAMT = "15"
		}
		unit = append(unit, QARGS)
		URL = BuildURL(unit)
	} else {
		QARGS = ""
		QAMT = ""
		return
	}

	req, err := http.NewRequest(CONFIG.Historysearch.Method, URL, nil)
	if err != nil {
		fmt.Println(err)
	}

	// set the appropriate headers before making request
	SetHeaders(*req, "historySearch")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	if !(resp.StatusCode >= 200 && resp.StatusCode <= 299) {
		if DEBUG == true {
			fmt.Println("HTTP Response Status:", resp.StatusCode, http.StatusText(resp.StatusCode))
			fmt.Println("Press Enter to continue")
			fmt.Scanln()
		}
		// Somethings wrong, call main again to initiate new client
		main()
	}
	defer resp.Body.Close()

	faults := new(lhJSON)

	// write the resp.body json to faults (lhJSON)
	err = json.NewDecoder(resp.Body).Decode(&faults)
	if err != nil {
		fmt.Println(err)
	}

	RenderTable(faults)

	Refresh()
}

func RenderTable(faults *lhJSON) {
	table := tablewriter.NewWriter(os.Stdout)
	switch MODE {
	case "historySearch", "histText":
		table.SetHeader([]string{"Fault History", "Start Time", "End Time"})
		table.SetTablePadding("\t") // pad with tabs
		table.SetAutoWrapText(false)
		table.SetAutoFormatHeaders(true)
		table.SetRowLine(true)
		table.SetCaption(true, "Last Updated: "+time.Now().Format(time.Stamp))
		// iterate over each fault and print Mechanical Name, Start Time, Notification
		for _, val := range faults.Data {
			et := time.Unix((val.EndTime / 1000), 0)
			st := time.Unix((val.StartTime / 1000), 0)
			// Remove -PHYSICAL from Node it, st is the start time, val.Text is the actual fault
			var row = []string{strings.Replace(val.NodeID, "-PHYSICAL", " ", 1) + val.Text, st.Format("01/02/2006") + " " + st.Format(time.Kitchen), et.Format("01/02/2006") + " " + et.Format(time.Kitchen)}
			table.Append(row)
		}
		table.Render() // Send output
		break
	case "currentFaults":
		table.SetHeader([]string{"Mechanical Name", "Start Time", "Notification"})
		table.SetTablePadding("\t") // pad with tabs
		table.SetAutoWrapText(false)
		table.SetAutoFormatHeaders(true)
		table.SetRowLine(true)
		table.SetCaption(true, "Last Updated: "+time.Now().Format(time.Stamp))

		for _, val := range faults.Data {
			st := time.Unix((val.StartTime / 1000), 0)
			// Remove -PHYSICAL from Node it, st is the start time, val.Text is the actual fault
			var row = []string{strings.Replace(val.NodeID, "-PHYSICAL", "", 1), st.Format(time.Kitchen) + " " + st.Format("2006/01/02"), val.Text}
			if strings.HasPrefix(strings.Replace(val.NodeID, "-PHYSICAL", "", 1), "C05") {
				table.Rich(row, []tablewriter.Colors{tablewriter.Colors{tablewriter.BgRedColor, tablewriter.FgWhiteColor}, tablewriter.Colors{tablewriter.BgRedColor, tablewriter.FgWhiteColor}, tablewriter.Colors{tablewriter.BgRedColor, tablewriter.FgWhiteColor}})
				continue
			} else if strings.HasPrefix(strings.Replace(val.NodeID, "-PHYSICAL", "", 1), "N") {
				table.Rich(row, []tablewriter.Colors{tablewriter.Colors{tablewriter.BgRedColor, tablewriter.FgWhiteColor}, tablewriter.Colors{tablewriter.BgRedColor, tablewriter.FgWhiteColor}, tablewriter.Colors{tablewriter.BgRedColor, tablewriter.FgWhiteColor}})
				continue
			}
			table.Append(row)
		}
		table.Render() // Send output
		break
	default:
	}
}

func SetHeaders(_req http.Request, origin string) {
	switch origin {
	case "cookie":
		_req.Host = CONFIG.Cookie.Host
		for _, val := range CONFIG.Cookie.Headers {
			_req.Header.Set(val.Key, val.Val)
		}
		break
	case "login":
		for _, val := range CONFIG.Login.Headers {
			_req.Header.Set(val.Key, val.Val)
		}
		break
	case "historySearch":
		for _, val := range CONFIG.Historysearch.Headers {
			_req.Header.Set(val.Key, val.Val)
		}
		break
	case "currentFaults":
		for _, val := range CONFIG.Currentfaults.Headers {
			_req.Header.Set(val.Key, val.Val)
		}
	default:
	}
}
