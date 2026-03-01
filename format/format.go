package format

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
)

// Printer handles output formatting (table or JSON).
type Printer struct {
	JSON bool
}

// Table writes tabular data to stdout.
func (p *Printer) Table(headers []string, rows [][]string) {
	if p.JSON {
		p.tableJSON(headers, rows)
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
	fmt.Fprintln(w, strings.Join(headers, "\t"))
	for _, row := range rows {
		fmt.Fprintln(w, strings.Join(row, "\t"))
	}
	w.Flush()
}

func (p *Printer) tableJSON(headers []string, rows [][]string) {
	var result []map[string]string
	for _, row := range rows {
		m := make(map[string]string)
		for i, h := range headers {
			if i < len(row) {
				m[strings.ToLower(h)] = row[i]
			}
		}
		result = append(result, m)
	}
	printJSON(result)
}

// Detail writes key-value detail to stdout.
func (p *Printer) Detail(pairs [][]string) {
	if p.JSON {
		m := make(map[string]string)
		for _, pair := range pairs {
			m[strings.ToLower(pair[0])] = pair[1]
		}
		printJSON(m)
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 2, 2, ' ', 0)
	for _, pair := range pairs {
		fmt.Fprintf(w, "%s:\t%s\n", pair[0], pair[1])
	}
	w.Flush()
}

// Block writes a labeled block of text (e.g., notes).
func (p *Printer) Block(label, text string, w io.Writer) {
	if text == "" {
		return
	}
	if w == nil {
		w = os.Stdout
	}
	fmt.Fprintf(w, "\n%s:\n", label)
	for _, line := range strings.Split(text, "\n") {
		fmt.Fprintf(w, "  %s\n", line)
	}
}

func printJSON(v any) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(v)
}
