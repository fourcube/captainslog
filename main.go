package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"sort"
	"strings"
	"text/template"
	"time"
)

const FileTemplate = `
## END ##
#########
{{ range . }}
# Date {{ .Time }}
# ===================================
{{ range .Lines }}{{ if . }}# {{ . }}
{{ end }}{{end}}
{{ end }}
`

const TimeLayout = "Jan 2, 2006 at 3:04pm (MST)"

func main() {
	// Get $EDITOR and $CAPTAINSLOG
	editor, logpath := settings()
	log.Printf("Captainslog: %s", logpath)

	tempLogFile := createTempFile()
	defer tempLogFile.Close()

	writeHeader(tempLogFile, logpath)
	startEditor(editor, tempLogFile.Name())

	text := getText(tempLogFile)

	if len(text) > 0 {
		fmt.Println("#####")
		fmt.Println(text)
		fmt.Println("#####")

		err := appendLog(text, logpath)

		if err != nil {
			log.Printf("Keeping temporary log file %s", tempLogFile.Name())
			return
		}

		tempLogFile.Close()
		os.Remove(tempLogFile.Name())
	} else {
		log.Printf("Nothing logged.")
	}

	tempLogFile.Close()
	os.Remove(tempLogFile.Name())
}

func appendLog(text, logFilePath string) error {
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0660)
	defer logFile.Close()

	if err != nil {
		log.Printf("Failed to open log file, %v", err)
		return err
	}

	// Add Header
	header := fmt.Sprintf("## %s\n\n", time.Now().Format(TimeLayout))
	logFile.WriteString(header)

	if err != nil {
		log.Printf("Failed to write to log file, %v", err)
		return err
	}

	_, err = logFile.WriteString(text + "\n")
	if err != nil {
		log.Printf("Failed to write to log file, %v", err)
		return err
	}

	return nil
}

func getText(tempLogFile *os.File) string {
	var s, text string
	var err error

	tempLogFile.Seek(0, 0)
	r := bufio.NewReader(tempLogFile)

	for err == nil {
		s, err = r.ReadString('\n')
		log.Printf(s)

		// Skip comments
		if strings.HasPrefix(s, "#") {
			continue
		}

		s = strings.TrimSpace(s)
		if len(s) > 0 {
			text += fmt.Sprintf("%s\n", s)
		}
	}

	return text
}

func startEditor(editor, tempLogFilePath string) {
	cmd := exec.Command(editor, tempLogFilePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

type sortableLogEntries []LogEntry

func (l sortableLogEntries) Len() int {
	return len(l)
}

func (l sortableLogEntries) Less(i, j int) bool {
	return l[i].Before(l[j].Time)
}

func (l sortableLogEntries) Swap(i, j int) {
	l[j], l[i] = l[i], l[j]
}

func writeHeader(tempLogFile *os.File, logpath string) {
	logs, err := ioutil.ReadFile(logpath)
	var res sortableLogEntries

	if err != nil {
		res = make(sortableLogEntries, 0)
	} else {
		res = sortableLogEntries(Parse(string(logs)))
		sort.Sort(sort.Reverse(res))
	}

	previewCount := 5
	if len(res) < 5 {
		previewCount = len(res)
	}

	data := res[:previewCount]

	t := template.New("Header")
	t, err = t.Parse(FileTemplate)
	if err != nil {
		log.Fatalf("Couldn't parse template %s, %v", FileTemplate, err)
	}

	err = t.ExecuteTemplate(tempLogFile, "Header", data)
	if err != nil {
		panic("Couldn't write header to temporary log file.")
	}
}

func createTempFile() (tempFile *os.File) {
	tempFile, err := ioutil.TempFile("", "captainslog")
	if err != nil {
		panic("Couldn't create temporary file to edit log entry")
	}

	return
}

func settings() (editor string, path string) {
	path = os.Getenv("CAPTAINSLOG")
	if len(path) < 1 {
		panic("$CAPTAINSLOG environment variable not set!")
	}

	editor = os.Getenv("EDITOR")
	if len(editor) < 1 {
		// Assume nano
		log.Printf("Defaulting to 'nano' editor")
		editor = "nano"
	}

	editor, err := exec.LookPath(editor)
	if err != nil {
		log.Fatalf("$EDITOR not found!", editor)
	}
	log.Printf("Using %s", editor)

	return
}
