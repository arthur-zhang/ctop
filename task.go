package main

import (
	"fmt"
	linuxproc "github.com/c9s/goprocinfo/linux"
	"os"
	"strconv"
)

type TaskCount struct {
	Total           int
	Running         int
	Sleeping        int
	Stopped         int
	Zombie          int
	Uninterruptible int
}

func GetTaskCount() (c TaskCount) {
	f, err := os.Open("/proc")
	if err != nil {
		fmt.Println(err)
		return
	}
	files, err := f.Readdir(0)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, v := range files {
		fileName := v.Name()
		if fileName[0] < '0' || fileName[0] > '9' {
			continue
		}
		pid, err := strconv.Atoi(v.Name())
		if err != nil {
			return
		}

		taskDirFile, err := os.Open(fmt.Sprintf("/proc/%d/task", pid))
		if err != nil {
			continue
		}

		taskFiles, err := taskDirFile.ReadDir(0)
		if err != nil {
			continue
		}

		for _, taskFile := range taskFiles {
			tid, err := strconv.Atoi(taskFile.Name())
			if err != nil {
				continue
			}
			p, err := linuxproc.ReadProcessStat(fmt.Sprintf("/proc/%d/task/%d/stat", pid, tid))
			if err != nil {
				continue
			}
			switch p.State {
			case "R":
				c.Running++
			case "D":
				c.Uninterruptible++
			case "t", "T":
				c.Stopped++
			case "Z":
				c.Zombie++
			default:
				c.Sleeping++
			}
			c.Total++
		}
	}

	return
}
