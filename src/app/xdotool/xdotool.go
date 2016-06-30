package xdotool

import (
	"log"
	"os/exec"
	"strings"
)

func Check() {
	_, err := exec.LookPath("xdotool")
	if err != nil {
		log.Fatal("xdotool was not found")
	}
}

func GetWindowTitle() string {
	out, err := exec.Command("xdotool", "getwindowfocus", "getwindowname").Output()
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimSpace(string(out))
}
