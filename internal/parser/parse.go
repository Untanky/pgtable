package parser

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"regexp"

	"github.com/untanky/pgtable"
)

var ErrInvalidFormat = errors.New("invalid format")

const columnSeparator = '|'

func Parse(reader io.Reader) (pgtable.Table, error) {
	scanner := bufio.NewScanner(reader)

	columns, err := parseColumns(scanner)
	if err != nil {
		return pgtable.Table{}, err
	}

	if err := parseSeparator(scanner); err != nil {
		return pgtable.Table{}, err
	}

	cells, err := parseCells(scanner)
	if err != nil {
		return pgtable.Table{}, err
	}

	return pgtable.Table{
		Columns: columns,
		Cells:   cells,
	}, nil
}

func readNextLine(scanner *bufio.Scanner) ([]byte, error) {
	ok := scanner.Scan()
	if !ok {
		return nil, scanner.Err()
	}

	return scanner.Bytes(), nil
}

func parseColumns(scanner *bufio.Scanner) ([]pgtable.Column, error) {
	line, err := readNextLine(scanner)
	if err != nil {
		return nil, err
	}

	columnBytes := bytes.Split(line, []byte{columnSeparator})
	columns := make([]pgtable.Column, len(columnBytes))

	for idx, column := range columnBytes {
		width := len(column)
		column = bytes.TrimSpace(column)
		name := make([]byte, len(column))

		copy(name, column)

		columns[idx] = pgtable.Column{
			Name:  string(name),
			Width: width,
		}
	}

	return columns, nil
}

func parseSeparator(scanner *bufio.Scanner) error {
	_, err := readNextLine(scanner)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidFormat, err)
	}

	return nil
}

func parseCells(scanner *bufio.Scanner) ([]pgtable.Cell, error) {
	cells := make([]pgtable.Cell, 0)

	regex := regexp.MustCompile("\\((\\d*) rows?\\)$")

	for {
		line, err := readNextLine(scanner)
		if err != nil {
			return nil, err
		}

		if regex.Match(line) {
			return cells, nil
		}

		columnBytes := bytes.SplitSeq(line, []byte{columnSeparator})

		for column := range columnBytes {
			column = bytes.TrimSpace(column)
			value := make([]byte, len(column))

			copy(value, column)

			cells = append(cells, pgtable.Cell{
				Value:  string(value),
				IsNull: len(value) == 0,
			})
		}

	}
}
