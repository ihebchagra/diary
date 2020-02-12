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
	"path/filepath"
)

func check(e error) {
    if e != nil {
        fmt.Println(e)
	os.Exit(1)
    }
}

func deleteHTML(diarydir string) {
	files, err := filepath.Glob(diarydir + "*.html")
	check(err)
	for _, f := range files {
		err = os.Remove(f)
		check(err)
	}
}



func formatDate(t time.Time) string {
    suffix := "th"
    switch t.Day() {
    case 1, 21, 31:
        suffix = "st"
    case 2, 22:
        suffix = "nd"
    case 3, 23:
        suffix = "rd"
    }
    return t.Format("January 2" + suffix + " 2006")
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
	base := fmt.Sprintf("%s%s",diarydir,entry)
	file := fmt.Sprintf("%s.md",base)

	fancydate := formatDate(t)
	content := []byte(fmt.Sprintf("%% %s\n\n*Dear Diary*,\n\n\n", fancydate))

	isnewentry := false
	cmd := exec.Command("")
	if _, err := os.Stat(file); os.IsNotExist(err) {
		isnewentry = true
		ioutil.WriteFile(file, content, 0755)
		cmd = exec.Command(editor,"+ normal G $","+startinsert",file)
	} else {
		cmd = exec.Command(editor,"+ normal G $",file)
	}

	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Run()

	if isnewentry {
		newcontent, err := ioutil.ReadFile(file)
		check(err)
		if strings.TrimSpace(string(newcontent)) == strings.TrimSpace(string(content)) {
			err = os.Remove(file)
			check(err)
		}
	}
}

func main() {
	home := os.Getenv("HOME");
	diarydir := os.Getenv("DIARY_DIR");
	if diarydir == "" {
		diarydir = fmt.Sprintf("%s/dox/diary/", home)
	}
	if _, err := os.Stat(diarydir); os.IsNotExist(err) {
		os.Mkdir(diarydir, 0755)
	}

	deleteHTML(diarydir)

	if len(os.Args) > 1 {
		if os.Args[1] == "search" {
			diarysearch(diarydir)
		}
	} else {
		diarytoday(diarydir)
	}


	deleteHTML(diarydir)
}
