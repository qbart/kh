package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/wzshiming/ctc"
)

func init() {
	flag.Parse()
}

func main() {
	trusted := readKnownHosts(".db")
	active := readKnownHosts("")

	if flag.Arg(0) == "reset" {
		writeKnownHosts("", trusted)
		fmt.Printf("%s%d hosts written%s\n", ctc.ForegroundGreen, len(trusted), ctc.Reset)
		return
	}

	// check diff
	temp := make([]string, 0)
	for _, a := range active {
		found := false
		for _, t := range trusted {
			if t == a {
				found = true
				break
			}
		}
		if !found {
			temp = append(temp, a)
		}
	}

	if len(temp) > 0 {
		fmt.Print(ctc.ForegroundYellow, "New hosts (", len(temp), ")", ctc.Reset, "\n")
		for _, s := range temp {
			fmt.Println(s)
		}
	} else {
		fmt.Print(ctc.ForegroundGreen, "No new hosts\n", ctc.Reset)
	}
}

func readKnownHosts(suffix string) []string {
	res := make([]string, 0)
	if f, err := os.Open(fmt.Sprint(os.Getenv("HOME"), "/.ssh/known_hosts", suffix)); err == nil {
		defer f.Close()

		b, err := ioutil.ReadAll(f)
		if err == nil {
			for _, s := range strings.Split(string(b), "\n") {
				if strings.TrimSpace(s) != "" {
					res = append(res, s)
				}
			}
		}
	}

	return res
}

func writeKnownHosts(suffix string, lines []string) {
	if f, err := os.OpenFile(fmt.Sprint(os.Getenv("HOME"), "/.ssh/known_hosts", suffix), os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
		defer f.Close()

		writer := bufio.NewWriter(f)

		for _, line := range lines {
			writer.WriteString(fmt.Sprint(line, "\n"))
		}
	}
}
