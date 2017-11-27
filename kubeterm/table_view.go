package kubeterm

import (
	"fmt"
)

type TableView struct {
	Format string
	Header []string
	Rows   [][]string
}

func NewTableView(format string) *TableView {
	return &TableView{
		Header: []string{},
		Rows:   [][]string{},
		Format: format,
	}
}

func (vt *TableView) AddHeader(header []string) {
	vt.Header = header
}

func (vt *TableView) AddRow(row []string) {
	vt.Rows = append(vt.Rows, row)
}

func (vt *TableView) Lines() []string {
	lines := []string{}

	for _, h := range vt.Header {
		lines = append(lines, fmt.Sprintf(vt.Format, h))
	}

	for _, row := range vt.Rows {
		lines = append(lines, fmt.Sprintf(vt.Format, row))
	}

	return lines
}

func (vt *TableView) Height() int {
	return len(vt.Rows) + len(vt.Header)
}
