package kubeterm

import (
	"fmt"
)

type TableView struct {
	Format  string
	Headers [][]interface{}
	Rows    [][]interface{}
}

func NewTableView(format string) *TableView {
	return &TableView{
		Headers: [][]interface{}{},
		Rows:    [][]interface{}{},
		Format:  format,
	}
}

func (vt *TableView) AddHeader(header ...interface{}) {
	vt.Headers = append(vt.Headers, header)
}

func (vt *TableView) AddRow(row ...interface{}) {
	vt.Rows = append(vt.Rows, row)
}

func (vt *TableView) Lines() []string {
	lines := []string{}

	for _, header := range vt.Headers {
		lines = append(lines, fmt.Sprintf(vt.Format, header...))
	}

	for _, row := range vt.Rows {
		lines = append(lines, fmt.Sprintf(vt.Format, row...))
	}

	return lines
}

func (vt *TableView) Reset() {
	vt.Rows = [][]interface{}{}
}

func (vt *TableView) Height() int {
	return len(vt.Rows) + len(vt.Headers)
}
