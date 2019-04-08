package main

import (
	"fmt"
	"github.com/peterh/liner"
	"github.com/sstark/knxbaosip"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	wordSep = " "
)

var (
	cmds          = []string{"quit", "read", "write", "group", "list"}
	selectedGroup string
)

func makePrompt() (s string) {
	if selectedGroup != "" {
		s = fmt.Sprintf("[%s] ", selectedGroup)
	} else {
		s = ">>> "
	}
	return
}

func groupList(g map[string][]int) (l []string) {
	for k, _ := range g {
		l = append(l, k)
	}
	return
}

func complete(l interface{}, word string) (c []string) {
	var rl []string
	switch tl := l.(type) {
	case []string:
		rl = tl
	case []int:
		for _, elem := range tl {
			rl = append(rl, fmt.Sprintf("%d", elem))
		}
	default:
	}
	for _, n := range rl {
		if strings.HasPrefix(n, word) {
			c = append(c, strings.TrimPrefix(n, word))
		}
	}
	return
}

func prompt(knx *knxbaosip.Client, groups GroupMap) {
	// word completer for liner
	var wc = func(line string, pos int) (head string, c []string, tail string) {
		var wordPos int
		var tl int
		head = line[:pos]
		tail = line[pos:]
		words := strings.Split(line, wordSep)
		// set wordPos to the word where the cursor sits currently
		for i, w := range words {
			tl += len(w) + len(wordSep)
			if tl > pos {
				wordPos = i
				break
			}
		}
		switch wordPos {
		case 0:
			c = complete(cmds, words[0])
		case 1:
			switch words[0] {
			case "group":
				c = complete(groupList(groups), words[1])
			case "read":
				c = complete(groups[selectedGroup], words[1])
			}
		default:
		}
		return
	}
	line := liner.NewLiner()
	defer line.Close()
	line.SetCtrlCAborts(true)
	line.SetWordCompleter(wc)
	line.SetTabCompletionStyle(liner.TabCircular)

	for {
		if name, err := line.Prompt(makePrompt()); err == nil {
			words := strings.Split(name, wordSep)
			switch words[0] {
			case "":
				continue
			case "quit":
				os.Exit(0)
				line.AppendHistory(name)
			case "group":
				if len(words) > 1 {
					_, ok := groups[words[1]]
					if ok {
						selectedGroup = words[1]
					} else {
						log.Println("group not found")
					}
				}
				line.AppendHistory(name)
			case "list":
				listGroups(groups)
				line.AppendHistory(name)
			case "read":
				var dps []int
				dps = groups[selectedGroup]
				if len(words) > 1 && words[1] != "" {
					val, convErr := strconv.ParseInt(words[1], 0, 32)
					if convErr != nil {
						log.Printf("not a valid data point: %s\n", words[1])
						break
					}
					dps = []int{int(val)}
				}
				readDatapoints(knx, dps)
				line.AppendHistory(name)
			default:
				log.Println("command not found")
			}
		} else if err == liner.ErrPromptAborted {
			log.Print("Aborted")
		} else if err == io.EOF {
			log.Print("EOF")
			os.Exit(0)
		} else {
			log.Print("Error reading line: ", err)
		}
	}
}