package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/xuri/excelize/v2"
)

func processExcelFile(inputFile string) error {
	f, err := excelize.OpenFile(inputFile)
	if err != nil {
		return err
	}
	defer f.Close()

	// Get all rows in the first sheet
	rows, err := f.GetRows(f.GetSheetName(0))
	if err != nil {
		return err
	}

	currentInitial := ""
	currentRows := make(map[string][][]string)

	for i, row := range rows {
		if i > 0 && len(row) > 0 { // Check if the row is not empty and skip the first row
			cell := fmt.Sprintf("A%d", i+1) // Construct the cell reference for the first column
			if isBold(f, cell) {
				// If we encounter a new sales rep initials
				currentInitial = row[0]
				currentRows[currentInitial] = nil // Initialize the slice for this initial
			}
			// If the current initial is not empty, add the row to the corresponding sales rep
			if currentInitial != "" {
				currentRows[currentInitial] = append(currentRows[currentInitial], row)
			}
		}
	}

	// Write each sales rep's data to a new Excel file
	for initial, rows := range currentRows {
		if err := writeToExcelFile(initial, rows); err != nil {
			return err
		}
	}

	return nil
}

func isBold(f *excelize.File, cell string) bool {
	// Get the style ID for the cell
	styleID, err := f.GetCellStyle("Sheet1", cell)
	if err != nil {
		fmt.Printf("Error getting style for cell %s: %v\n", cell, err)
		return false // Return false if there's an error retrieving the style
	}
	if styleID == 0 {
		fmt.Printf("No style found for cell %s\n", cell)
		return false // Handle case where no style is found
	}

	// Get the style details
	style, err := f.GetStyle(styleID)
	if err != nil {
		fmt.Printf("Error getting style details for cell %s: %v\n", cell, err)
		return false // Return false if there's an error retrieving the style
	}

	return style.Font.Bold
}

func writeToExcelFile(initial string, rows [][]string) error {
	// Create output directory if it doesn't exist
	outputDir := "output"
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return err
	}

	f := excelize.NewFile()

	// Get the current date and calculate the previous month's first and last days
	now := time.Now()
	firstDayOfPrevMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).AddDate(0, -1, 0)
	lastDayOfPrevMonth := firstDayOfPrevMonth.AddDate(0, 0, -1)

	// Define the title and date range
	title := "Commission Report for Biely & Shoaf"
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
		return err
	}
	f.SetCellValue("Sheet1", "A1", title)
	f.SetCellStyle("Sheet1", "A1", "A1", style)

	f.SetCellValue("Sheet1", "A2", dateRange)

	// Define the header
	header := []string{
		"SPers No", "Name", "Cust No", "BillToName", "BillToAdd1", "BillToAdd2",
		"BillToCity", "BillToState", "BillToZip", "ShipToName", "ShipToAdd1",
		"ShipToAdd2", "ShipToCity", "ShipToState", "ShipToZip", "Code No.", "PO Number",
		"Number", "Date", "to Comm.", "Payable",
	}

	// Set the header in the third row
	for j, colCell := range header {
		cell, _ := excelize.CoordinatesToCellName(j+1, 3) // Row 3 for header
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
		return err
	}
	fmt.Printf("Created file: %s\n", outputFile)
	return nil
}
