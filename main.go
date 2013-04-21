// Copyright 2013 marpie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"github.com/marpie/winfsnotify"
	"strings"
)

func main() {
	w, err := winfsnotify.NewWatcher()
	if err != nil {
		return
	}
	defer w.Close()

	event := flag.String("event", "all", "possible Events: all, access, close, create, delete, modify, move.")
	one_shot := flag.Bool("one-shot", false, "only execute once (true/false)")
	format := flag.String("format", "[[Filename]]", "[[Cookie]]=unique cookie of the event, [[Name]]=File/Directory name, [[Info]]=Debug information.")
	//recursive := flag.Bool("recursive", true, "only fire on directory events  (true/false)")
	flag.Parse()

	command := flag.Arg(0)

	if len(command) < 1 {
		println("Command to execute is missing!")
		return
	}

	/*
		output-format:
			[[Cookie]]
			[[Filename]]
			[[Info]]
	*/

	var mode uint32 = winfsnotify.FS_ALL_EVENTS // listen to all events
	switch *event {
	case "all":
		mode = winfsnotify.FS_ALL_EVENTS
	case "access":
		mode = winfsnotify.FS_ACCESS
	case "close":
		mode = winfsnotify.FS_CLOSE
	case "create":
		mode = winfsnotify.FS_CREATE
	case "delete":
		mode = winfsnotify.FS_DELETE
	case "modify":
		mode = winfsnotify.FS_MODIFY
	case "move":
		mode = winfsnotify.FS_MOVE
	default:
		println("Unknown mode! Listening to all events!")
	}

	if *one_shot {
		mode |= winfsnotify.FS_ONESHOT
	}

	if err := w.AddWatch(".", mode); err != nil {
		fmt.Printf("Error adding watch: %v", err)
		return
	}

	for event := range w.Event {
		output := strings.Replace(*format, "[[Cookie]]", fmt.Sprintf("%d", event.Cookie), -1)
		output = strings.Replace(output, "[[Filename]]", event.Name, -1)
		output = strings.Replace(output, "[[Info]]", event.String(), -1)

		// DEBUG
		println(output)
		// DEBUG END

		if *one_shot {
			break
		}
	}

	println("Done.")
}
