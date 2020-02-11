package main

import (
	"fmt"
	"time"
	"os/exec"
	"os"
	"io/ioutil"
	"io"
	"bytes"
	"strings"
	"errors"
)

func check(e error) {
    if e != nil {
        fmt.Println(e)
	os.Exit(1)
    }
}

func diarysearch(diarydir string) {
	//this generates a list of entries
	list, err := exec.Command("ls","-t",diarydir).Output()
	check(err)
	if string(list) == "" {
		check(errors.New("Your diary directory is empty, write an entry first"))
	}
	//command to pipe to
	cmd := exec.Command("fzf")
	//setup stdin pipe
	stdin, err := cmd.StdinPipe()
	check(err)
	io.WriteString(stdin, string(list))
	stdin.Close()
	//setup stderr
	cmd.Stderr = os.Stderr
	//setup stdout capture to get the chosen entry
	var entry bytes.Buffer
	cmd.Stdout = &entry
	//run the command
	cmd.Run()

	editor := os.Getenv("EDITOR");
	file := fmt.Sprintf("%s%s", diarydir, strings.TrimSuffix(entry.String(),"\n"))

	if file == diarydir {
		check(errors.New("Make sure to choose something"))
	}
	if _, err := os.Stat(file); os.IsNotExist(err) {
		check(errors.New("Choose a real file from the list"))
	}

	cmd = exec.Command(editor,file)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Run()
}


func diarytoday(diarydir string) {
	editor := os.Getenv("EDITOR");

	t := time.Now()
	entry := t.Format("02-Jan-2006")
	file := fmt.Sprintf("%s%s",diarydir,entry)

	fancydate := t.Format("January 02, 2006")
	content := []byte(fmt.Sprintf("%s\n\nDear Diary,\n\n\n", fancydate))

	cmd := exec.Command("")
	if _, err := os.Stat(file); os.IsNotExist(err) {
		ioutil.WriteFile(file, content, 0755)
		cmd = exec.Command(editor,"+ normal G $","+startinsert",file)
	} else {
		cmd = exec.Command(editor,"+ normal G $",file)
	}

	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func main() {
	home := os.Getenv("HOME");
	diarydir := os.Getenv("DIARY_DIR");
	if diarydir == "" {
		diarydir = fmt.Sprintf("%s/dox/diary/", home)
	}
	if _, err := os.Stat(diarydir); os.IsNotExist(err) {
		os.Mkdir(diarydir, 0755)
	}t 

	if len(os.Args) > 1 {
		if os.Args[1] == "search" {
			diarysearch(diarydir)
		} else {
			fmt.Println("Either type [ diary search ] to search entries or just [ diary ] to start or add to today's one");
		}
	} else {
		diarytoday(diarydir)
	}
}
