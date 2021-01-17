package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/pcrandall/figlet4go"
)

func PrintHeader() {

	CallClear()

	ascii := figlet4go.NewAsciiRender()
	colors := [...]color.Attribute{
		color.FgWhite,
	}
	options := figlet4go.NewRenderOptions()
	options.FontColor = make([]color.Attribute, len("LIGHTHOUSE"))
	for i := range options.FontColor {
		options.FontColor[i] = colors[i%len(colors)]
	}
	renderStr, _ := ascii.RenderOpts("LIGHTHOUSE", options)

	// remove the last three blank rows, all uppercase chars have a height of 8, the font height for default font is 11
	fmt.Println(renderStr[:len(renderStr)-len(renderStr)/11*3-1])
	banner := color.New(color.FgBlack, color.BgWhite).SprintFunc()

	switch MODE {
	case "historySearch":
		fmt.Printf("(f)Current Faults    %s    (t)Search Text    (h)Help          pcrandall '20\n\n", banner("(m)Search Mechanical Name"))
		break
	case "histText":
		fmt.Printf("(f)Current Faults    (m)Search Mechanical Name    %s    (h)Help          pcrandall '20\n\n", banner("(t)Search Text"))
		break
	case "currentFaults":
		fmt.Printf("%s    (m)Search Mechanical Name    (t)Search Text    (h)Help          pcrandall '20\n\n", banner("(f)Current Faults"))
		break
	case "help":
		fmt.Printf("(f)Current Faults    (m)Search Mechanical Name    (t)Search Text    %s         pcrandall '20\n\n", banner("(h)Help"))
		break
	}
}

func init() {
	CLEAR = make(map[string]func()) //Initialize it
	CLEAR["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	CLEAR["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func CallClear() {
	value, ok := CLEAR[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	//if we defined a clear func for that platform:
	if ok {
		value() //we execute it
	} else {
		panic("Your platform is unsupported! I can't clear terminal screen :(") //unsupported platform
	}
}

func Help() {
	banner := color.New(color.FgBlack, color.BgYellow).SprintFunc()

	PrintHeader()

	color.Set(color.FgBlack, color.BgWhite)
	fmt.Println("")
	fmt.Println("CURRENT FAULTS")
	fmt.Println("")
	color.Unset()
	fmt.Printf("\tWill refresh automatically every %s or when %s is pressed\n", banner("60 seconds"), banner("Enter"))
	color.Set(color.FgWhite, color.BgRed)
	fmt.Println("\tFaults in the Matrix and C05 area will be highlighted in red.")
	// fmt.Println("You can keep terminal open, will refresh every  or when Enter is pressed.")

	color.Set(color.FgBlack, color.BgWhite)
	fmt.Println("")
	fmt.Println("SEARCH MECHANICAL NAME")
	fmt.Println("")
	color.Set(color.FgBlack, color.BgYellow)
	fmt.Println("\tUSAGE: <Search><Single Space><Amount(optional)>")
	color.Unset()
	fmt.Println("\tSearches are not case sensitive, wildcards are supported. Default number of values returned are 15")
	fmt.Println("\tIf wildcards(*) are omitted only exact matches will be returned")
	fmt.Println("\tExamples:")
	fmt.Println("\t\t15 most recent faults for \"N1111\": \"n1111\"")
	fmt.Println("\t\t3 most recent faults for \"NL3085\": \"nl3085 3\"")
	fmt.Println("\t\t15 most recent faults for \"N1111\": \"n1111\"")
	fmt.Println("\t\t10 most recent Navette faults: \"n* 10\"")
	fmt.Println("\t\t5 most recent Navette Lift faults: \"nl* 5\"")
	fmt.Println("\t\t100 most recent faults in \"C05\": \"c05* 100\"")
	fmt.Println("\t\t15 most recent faults for \"Sorter 1\": \"*so21*\"")
	fmt.Println("\t\t10 most recent faults for all lifts ending in \"82\": \"*82 10\"")

	color.Set(color.FgBlack, color.BgWhite)
	fmt.Println("")
	fmt.Println("SEARCH TEXT")
	fmt.Println("")
	color.Set(color.FgBlack, color.BgYellow)
	fmt.Println("\tUSAGE: <Search><Single Space><Amount(optional)>")
	fmt.Println("\tYou can refine text search using multiple keywords separated by \"&&\" for example \"C05&&PEC 12\"")
	color.Set(color.FgBlack, color.BgWhite)
	fmt.Println("\t\tThis will return 12 results with the text field containing \"C05 and PEC\"")
	color.Unset() // Don't forget to unset
	fmt.Println("\tSearches ARE CASE sensitive. Wildcards are included by default. Default number of rows returned are 15.")
	fmt.Println("")
	fmt.Println("\t\tExamples:")
	fmt.Println("\t\t5 most recent Safety Gate Open faults: \"afet 5\"")
	fmt.Println("\t\t7 most recent PEC faults: \"PEC 7\"")
	fmt.Println("\t\t15 most recent clearance sensor faults: \"clearance 7\"")
	fmt.Println("")
	Refresh()
}

func GetUserInput(prompt string) string {
	fmt.Printf(prompt)
	r := bufio.NewReader(os.Stdin)
	str, err := r.ReadString('\n')
	if err != nil {
		panic(err)
	}
	args := strings.Fields(str)
	if len(args) > 1 {
		QAMT = args[1]
	}
	if len(args) > 0 {
		QARGS = args[0]
	}
	return strings.TrimSpace(str)
}
