package parser

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/untanky/pgtable"
)

var ErrInvalidFormat = errors.New("invalid format")

const columnSeparator = '|'

func Parse(reader io.Reader) (pgtable.Table, error) {
	bufferedReader := bufio.NewReader(reader)

	columns, err := parseColumns(bufferedReader)
	if err != nil {
		return pgtable.Table{}, err
	}

	if err := parseSeparator(bufferedReader); err != nil {
		return pgtable.Table{}, err
	}

	cells, err := parseCells(bufferedReader)
	if err != nil {
		return pgtable.Table{}, err
	}

	return pgtable.Table{
		Columns: columns,
		Cells:   cells,
	}, nil
}

func readNextLine(reader *bufio.Reader) ([]byte, error) {
	line, err := reader.ReadBytes('\n')
	if err != nil {
		return nil, fmt.Errorf("reading line: %w", err)
	}

	line = bytes.TrimSuffix(line, []byte{'\n'})

	return line, nil
}

func parseColumns(bufferedReader *bufio.Reader) ([]pgtable.Column, error) {
	line, err := readNextLine(bufferedReader)
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

func parseSeparator(reader *bufio.Reader) error {
	_, err := readNextLine(reader)
	if err == nil {
		return nil
	}

	if errors.Is(err, io.EOF) {
		return ErrInvalidFormat
	}

	return fmt.Errorf("%w: %w", ErrInvalidFormat, err)
}

func parseCells(reader *bufio.Reader) ([]pgtable.Cell, error) {
	cells := make([]pgtable.Cell, 0)

	for {
		line, err := readNextLine(reader)
		if err != nil {
			if errors.Is(err, io.EOF) {
				if err := parseFooter(line); err != nil {
					return nil, err
				}

				return cells, nil
			}

			return nil, err
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

func parseFooter(line []byte) error {
	return nil
}
