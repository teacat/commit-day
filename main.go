package main

import (
	"fmt"
	. "fmt"
	"os"
	"os/exec"
	"regexp"
	"time"
)

func main() {
	var c []commitBox

	app := "git"

	arg0 := "log"
	arg1 := "--branches"
	arg2 := "--not"
	arg3 := "--remotes"

	cmd := exec.Command(app, arg0, arg1, arg2, arg3)
	stdout, err := cmd.Output()

	if err != nil {
		Println(err.Error())
		return
	}

	//Print()

	commitRegexp := regexp.MustCompile(`commit (.*)\nAuthor: (.*)\nDate:   (.*)\n\n    (.*)`)
	commits := commitRegexp.FindAllStringSubmatch(string(stdout), -1)

	for _, commit := range commits {
		id := commit[1]
		//author := commit[2]
		date := commit[3]
		message := commit[4]

		app := "git"

		arg0 := "show"
		arg1 := id

		cmd := exec.Command(app, arg0, arg1)
		stdout, err := cmd.Output()

		if err != nil {
			Println(err.Error())
			return
		}

		diffRegexp := regexp.MustCompile(`diff --git a\/(.*) b\/(.*)`)
		diffs := diffRegexp.FindAllStringSubmatch(string(stdout), -1)

		var srcs []string
		var lastTime time.Time
		for _, diff := range diffs {
			src := diff[1]

			file, err := os.Stat(src)

			if err != nil {
				fmt.Println(err)
			}

			//Printf("%s, %s, %s, %s\n", id, src, file.ModTime(), message)

			if lastTime.Before(file.ModTime()) {
				lastTime = file.ModTime()
			}
			srcs = append(srcs, src)
		}

		c = append(c, commitBox{
			files:   srcs,
			message: message,
			oldDate: date,
			newDate: lastTime.Format("Mon Jan 2 15:04:05 2006 -0700"),
		})
	}

	app = "git"

	arg0 = "reset"
	arg1 = "origin/master"

	cmd = exec.Command(app, arg0, arg1)
	stdout, err = cmd.Output()

	if err != nil {
		Println(err.Error())
		return
	}

	//Println(string(stdout))

	for _, v := range c {

		app := "git"
		argAdd := []string{"add"}

		for _, k := range v.files {
			argAdd = append(argAdd, k)
		}

		cmd := exec.Command(app, argAdd...)
		_, err := cmd.Output()

		if err != nil {
			Println(err.Error())
			return
		}
		//Println(string(stdout))

		app = "git"
		arg0 := "commit"
		arg1 := "-m"
		arg2 := Sprintf(`%s`, v.message)
		arg3 := "--date"
		arg4 := Sprintf(`"%s"`, v.newDate)

		cmd = exec.Command(app, arg0, arg1, arg2, arg3, arg4)
		stdout, err = cmd.Output()

		if err != nil {
			Println(err.Error())
			return
		}
		//Println(string(stdout))

		Printf("%s, %s, %s -> %s\n", v.files, v.message, v.oldDate, v.newDate)
	}

}

type commitBox struct {
	files   []string
	message string
	oldDate string
	newDate string
}
