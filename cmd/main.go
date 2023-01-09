package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var voiceNameRegex *regexp.Regexp = regexp.MustCompile("\\s+")

const CONFIG_FILE = "./config.txt"

const RATE_DEF = "45"
const RATE_TIME = "300"

func getChineseVoiceName() string {
	out, err := exec.Command("say", "-v", "?").Output()
	if err != nil {
		log.Fatalf("Failed to retrieve voice names\n")
		return ""
	}
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		sections := voiceNameRegex.Split(line, -1)
		if len(sections) < 2 {
			continue
		}
		if strings.HasPrefix(sections[1], "zh_CN") {
			return sections[0]
		}
	}
	log.Fatalf("Failed to get Chinese voice name\n")
	return ""
}

func speak(str string, chineseVoiceName string, rate string) {
	cmd := exec.Command("say", "-v", chineseVoiceName, "-r", rate, str)
	_ = cmd.Run()
}

func readConfig() (time.Time, time.Time) {
	b, err := os.ReadFile(CONFIG_FILE)
	if err != nil {
		log.Fatalf("Failed to read file: %s\n", CONFIG_FILE)
	}
	str := string(b)
	str = strings.TrimSuffix(str, "\n")
	str = strings.TrimSuffix(str, "\r")
	arr := strings.Split(str, "-")

	if len(arr) < 2 {
		log.Fatalf("Failed to parse config file: %s due to it has no startDateStr or endDateStr \n", CONFIG_FILE)
	}

	startTimeStr := arr[0]
	endTimeStr := arr[1]

	now := time.Now()
	startDateTimeStr := fmt.Sprintf("%d-%02d-%02dT%s:00+08:00", now.Year(), now.Month(), now.Day(), startTimeStr)
	endDateTimeStr := fmt.Sprintf("%d-%02d-%02dT%s:00+08:00", now.Year(), now.Month(), now.Day(), endTimeStr)

	startDateTime, err := time.Parse(time.RFC3339, startDateTimeStr)
	if err != nil {
		log.Fatalf("Failed to parser startDateTimeStr: %v\n", startDateTimeStr)
	}

	endDateTime, err := time.Parse(time.RFC3339, endDateTimeStr)
	if err != nil {
		log.Fatalf("Failed to parser endDateTimeStr: %v\n", endDateTimeStr)
	}

	if endDateTime.Before(startDateTime) {
		endDateTime = endDateTime.Add(24 * time.Hour)
	}

	return startDateTime, endDateTime
}

func main() {
	chineseVoiceName := getChineseVoiceName()
	startDateTime, endDateTime := readConfig()
	fmt.Printf("The time duration is [%s - %s]\n", startDateTime.Format(time.RFC3339), endDateTime.Format(time.RFC3339))

	showTimeStr := ""
	for {
		now := time.Now()
		if now.After(endDateTime) {
			speak("时间到, 程序结束", chineseVoiceName, RATE_DEF)
			break
		}

		if now.Before(startDateTime) {
			time.Sleep(1 * time.Second)
			continue
		}

		nowStr := now.Format(time.RFC3339)
		if showTimeStr != nowStr {
			fmt.Printf("\033[2K\rNow is: %v", now.Format(time.RFC3339))
			sec := now.Second()
			str := strconv.Itoa(sec)
			speak(str, chineseVoiceName, RATE_TIME)
			showTimeStr = nowStr
		}
	}
}
