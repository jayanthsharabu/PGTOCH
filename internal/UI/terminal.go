package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func PrintLogo() {
	fmt.Println(LogoStyle.Render(logo))
}

func PrintTitle(title string) {
	fmt.Println(TitleStyle.Render(title))
}

func PrintSubtitle(subtitle string) {
	fmt.Println(SubtitleStyle.Render(subtitle))
}

func PrintSuccess(message string) {
	fmt.Println(SuccessStyle.Render("✓ " + message))
}

func PrintError(message string) {
	fmt.Println(ErrorStyle.Render("✗ " + message))
}

func PrintWarning(message string) {
	fmt.Println(WarningStyle.Render("! " + message))
}

func PrintInfo(message string) {
	fmt.Println(InfoStyle.Render(message))
}

func PrintHighlight(message string) {
	fmt.Println(HighlightStyle.Render(message))
}

func PrintBox(title string, content string) {
	titleText := HighlightStyle.Render(title)
	contentText := InfoStyle.Render(content)
	boxContent := lipgloss.JoinVertical(lipgloss.Left, titleText, contentText)
	fmt.Println(BoxStyle.Render(boxContent))
}

func ExitWithError(message string) {
	PrintError(message)
	os.Exit(1)
}

func displayTable(headers []string, rows [][]string) {
	colWidths := make([]int, len(headers))
	for i, header := range headers {
		colWidths[i] = len(header)
	}

	for _, row := range rows {
		for i, cell := range row {
			if i < len(colWidths) && len(cell) > colWidths[i] {
				colWidths[i] = len(cell)
			}
		}
	}

	headerCells := make([]string, len(headers))
	for i, header := range headers {
		headerCells[i] = TableHeaderStyle.Render(
			lipgloss.PlaceHorizontal(
				colWidths[i]+2,
				lipgloss.Left,
				header,
			),
		)
	}

	headerRow := lipgloss.JoinHorizontal(lipgloss.Top, headerCells...)
	fmt.Println(headerRow)

	separator := make([]string, len(headers))
	for i, width := range colWidths {
		separator[i] = strings.Repeat("─", width+2)
	}
	separatorRow := lipgloss.JoinHorizontal(lipgloss.Top, separator...)
	fmt.Println(HighlightStyle.Render(separatorRow))

	for _, row := range rows {
		rowCells := make([]string, len(row))
		for i, cell := range row {
			if i < len(colWidths) {
				rowCells[i] = TableCellStyle.Render(
					lipgloss.PlaceHorizontal(
						colWidths[i]+2,
						lipgloss.Left,
						cell,
					),
				)
			}
		}
		fmt.Println(lipgloss.JoinHorizontal(lipgloss.Top, rowCells...))
	}
}
