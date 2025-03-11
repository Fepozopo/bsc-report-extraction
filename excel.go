package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

// processExcelFile processes an input Excel file by parsing its contents,
// identifying rows associated with different sales representatives based on
// bolded initials in the first column, and writing separate Excel files for
// each sales rep group. It logs the processing steps to a log file.
// The input is the path to the Excel file, and it returns an error if any
// operation fails during processing.
func processExcelFile(reportType, inputFile string) error {
	logger, logFile, err := CreateLogger("processExcelFile", "", "", "INFO")
	if err != nil {
		return fmt.Errorf("error creating log file: %w", err)
	}
	defer logFile.Close()

	// Open the input Excel file
	f, err := excelize.OpenFile(inputFile)
	if err != nil {
		return fmt.Errorf("error opening Excel file: %w", err)
	}
	defer f.Close()

	// Get all rows in the first sheet
	rows, err := f.GetRows(f.GetSheetName(0))
	if err != nil {
		return fmt.Errorf("error getting rows: %w", err)
	}

	currentInitial := ""
	currentRows := make(map[string][][]string)

	var column string
	var columnIndex int
	switch {
	case reportType == "commission":
		column = "A"
		columnIndex = 0
	case reportType == "royalty":
		column = "B"
		columnIndex = 1
	default:
		return fmt.Errorf("unknown report type: %s", reportType)
	}

	for i, row := range rows {
		if i > 1 && len(row) > 0 { // Ignore the first two rows
			cell := fmt.Sprintf("%s%d", column, i+1) // Construct the cell reference for the sales rep initials
			bold, err := isBold(f, cell)
			if err != nil {
				return fmt.Errorf("error checking bold status for cell %s: %w", cell, err)
			}
			if bold && !strings.Contains(row[columnIndex], "Totals") {
				// If we encounter a new sales rep initials
				currentInitial = row[columnIndex]
				currentRows[currentInitial] = nil // Initialize the slice for this initial
			}
			// If the current initial is not empty, add the row to the corresponding sales rep
			if currentInitial != "" {
				currentRows[currentInitial] = append(currentRows[currentInitial], row)
			}
			logger.Printf("Processing row %d: current initial: %s, bold: %v\n%v", i+1, currentInitial, bold, row)
		}
	}

	// Write each sales rep's data to a new Excel file
	for initial, rows := range currentRows {
		if err := writeToExcelFile(reportType, initial, rows); err != nil {
			return fmt.Errorf("error writing to Excel file: %w", err)
		}
	}

	return nil
}

// isBold checks if the font style of a given cell in an Excel sheet is bold.
// It returns true if the font is bold, false otherwise. If there's an error
// retrieving the style or its details, an error is returned.
func isBold(f *excelize.File, cell string) (bool, error) {
	// Get the style ID for the cell
	styleID, err := f.GetCellStyle("Sheet1", cell)
	if err != nil {
		// Return false if there's an error retrieving the style
		return false, fmt.Errorf("error getting style for cell %s: %w", cell, err)
	}

	// Get the style details
	style, err := f.GetStyle(styleID)
	if err != nil {
		// Return false if there's an error retrieving the style
		return false, fmt.Errorf("error getting style details for cell %s: %w", cell, err)
	}

	return style.Font.Bold, nil
}

// writeToExcelFile writes a given set of rows to a new Excel file with title, date range, header, and data.
// The file is saved in the output directory with the given initial as the file name.
// If there's an error writing the file, an error is returned.
func writeToExcelFile(reportType, initial string, rows [][]string) error {
	// Create output directory if it doesn't exist
	outputDir := "output"
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return fmt.Errorf("error creating output directory: %w", err)
	}

	f := excelize.NewFile()

	// Get the current date and calculate the previous month's first and last days
	now := time.Now()
	firstDayOfPrevMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).AddDate(0, -1, 0)
	lastDayOfPrevMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).AddDate(0, 0, -1)

	// Define the title and date range
	var title string
	switch reportType {
	case "commission":
		title = "Commission Report for Biely & Shoaf"
	case "royalty":
		title = "Royalty Report for Biely & Shoaf"
	default:
		return fmt.Errorf("unknown report type: %s", reportType)
	}
	dateRange := fmt.Sprintf("For Invoices Dated %s to %s", firstDayOfPrevMonth.Format("1/2/2006"), lastDayOfPrevMonth.Format("1/2/2006"))

	// Set the title and date range in the first two rows
	style, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:   true,
			Size:   16,
			Family: "Calibri",
		},
	})
	if err != nil {
		return fmt.Errorf("error creating style: %w", err)
	}
	f.SetCellValue("Sheet1", "A1", title)
	f.SetCellStyle("Sheet1", "A1", "A1", style)

	f.SetCellValue("Sheet1", "A2", dateRange)

	// Define the header
	var header []string
	switch reportType {
	case "commission":
		header = []string{
			"SPers No", "Name", "Cust No", "BillToName", "BillToAdd1", "BillToAdd2", "",
			"BillToCity", "BillToState", "BillToZip", "ShipToName", "", "ShipToAdd1",
			"ShipToAdd2", "ShipToCity", "ShipToState", "ShipToZip", "Code No.", "PO Number",
			"Number", "Date", "to Comm.", "Payable",
		}
	case "royalty":
		header = []string{
			"", "Artist", "", "", "QuantityShipped", "Net Sales Amount", "Royalty Percent", "", "",
			"Royalty Amount",
		}
	}

	// Set the header in the third row
	for j, colCell := range header {
		cell, _ := excelize.CoordinatesToCellName(j+1, 4) // Row 4 for header
		f.SetCellValue("Sheet1", cell, colCell)
	}

	// Set the data rows starting from the fourth row
	for i, row := range rows {
		for j, colCell := range row {
			cell, _ := excelize.CoordinatesToCellName(j+1, i+5) // Start from Row 5 for data
			f.SetCellValue("Sheet1", cell, colCell)
		}
	}

	// Save the file in the output directory
	outputFile := filepath.Join(outputDir, fmt.Sprintf("%s.xlsx", initial))
	if err := f.SaveAs(outputFile); err != nil {
		return fmt.Errorf("error saving Excel file: %w", err)
	}
	fmt.Printf("Created file: %s\n", outputFile)
	return nil
}
