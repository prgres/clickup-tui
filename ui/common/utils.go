package common

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

func OpenUrlInWebBrowser(url string) error {
	switch runtime.GOOS {
	case "linux":
		return exec.Command("xdg-open", url).Start()
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		return exec.Command("open", url).Start()
	default:
		return fmt.Errorf("unsupported platform")
	}
}

type EditorFinishedMsg struct {
	Id   string
	Data interface{}
	Err  error
}

func OpenEditor(id string, data string) tea.Cmd {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}

	tmpfileName, err := createTempFileWithData(data)
	if err != nil {
		return ErrCmd(err)
	}

	execCmd := exec.Command(editor, tmpfileName)
	return tea.ExecProcess(execCmd, func(err error) tea.Msg {
		defer os.Remove(tmpfileName) // Clean up the file afterwards

		if err != nil {
			return EditorFinishedMsg{
				Err: err,
			}
		}

		updatedData, err := os.ReadFile(tmpfileName)

		return EditorFinishedMsg{
			Id:   id,
			Data: string(updatedData),
			Err:  err,
		}
	})
}

func createTempFileWithData(data string) (string, error) {
	tmpfile, err := os.CreateTemp("", "clickup-tui.*.txt")
	if err != nil {
		return "", fmt.Errorf("error creating temporary file: %w", err)
	}

	defer tmpfile.Close()

	if _, err := tmpfile.Write([]byte(data)); err != nil {
		return "", fmt.Errorf("error writing to temporary file: %w", err)
	}

	return tmpfile.Name(), nil
}

type ResourceType string

var ResourceTypeRegistry = struct {
	VIEW      ResourceType
	WIDGET    ResourceType
	COMPONENT ResourceType
}{
	VIEW:      "view",
	WIDGET:    "widget",
	COMPONENT: "component",
}

func NewLogger(l *log.Logger, resourceType ResourceType, id string) *log.Logger {
	return l.WithPrefix(l.GetPrefix() + string(resourceType) + id)
}
