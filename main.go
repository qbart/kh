package main

import (
	"bufio"
	"flag"
	"fmt"
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
	for _, t := range trusted {
		found := false
		for _, a := range active {
			if t == a {
				found = true
				break
			}
			if !found {
				temp = append(temp, t)
			}
		}
	}

	if len(temp) > 0 {
		fmt.Print(ctc.ForegroundGreen, "No new hosts", ctc.Reset)
	} else {
		fmt.Println("New hosts:")
		for _, s := range temp {
			fmt.Println(s)
		}
	}
}

func readKnownHosts(suffix string) []string {
	res := make([]string, 0)
	if f, err := os.Open(fmt.Sprint(os.Getenv("HOME"), "/.ssh/known_hosts", suffix)); err == nil {
		defer f.Close()

		reader := bufio.NewReader(f)

		for {
			line, err := reader.ReadString('\n')
			line = strings.TrimSuffix(line, "\n")
			res = append(res, line)

			if err != nil {
				break
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
