package things

import (
	"fmt"
	"net/url"
	"os/exec"
	"strings"
)

// ExecURL opens a things:/// URL via macOS `open` command.
func ExecURL(thingsURL string) error {
	cmd := exec.Command("open", thingsURL)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("open URL: %w\n%s", err, output)
	}
	return nil
}

// AddTodo builds and executes a things:///add URL.
func AddTodo(params map[string]string) error {
	u := buildURL("add", params)
	return ExecURL(u)
}

// UpdateTask builds and executes a things:///update URL.
func UpdateTask(params map[string]string) error {
	u := buildURL("update", params)
	return ExecURL(u)
}

func buildURL(command string, params map[string]string) string {
	var b strings.Builder
	b.WriteString("things:///")
	b.WriteString(command)
	b.WriteString("?")

	first := true
	for k, v := range params {
		if !first {
			b.WriteString("&")
		}
		b.WriteString(url.QueryEscape(k))
		b.WriteString("=")
		b.WriteString(url.QueryEscape(v))
		first = false
	}
	return b.String()
}
