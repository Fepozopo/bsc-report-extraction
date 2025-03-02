package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2/app"
)

func main() {
	startTime := time.Now()

	logger, logFile, err := CreateLogger("main", "", "", "ERROR")
	if err != nil {
		logger.Printf("failed to create log file: %v", err)
		return
	}
	defer logFile.Close()

	myApp := app.New()
	defer myApp.Quit()

	inputExcelFile := selectFiles(myApp)

	// If no file is selected, exit the program
	if inputExcelFile == "" {
		logger.Printf("No file selected")
		return
	}

	// Process the Excel file
	err = processExcelFile(inputExcelFile)
	if err != nil {
		logger.Printf("Error processing Excel file: %v", err)
		return
	}

	fmt.Printf("Done!\nElapsed time: %v\n", time.Since(startTime))

	for i := 0; i < 3; i++ {
		fmt.Printf("Quitting in %d seconds...\n", 3-i)
		time.Sleep(1 * time.Second)
	}
}
