package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/fsnotify/fsnotify"
)

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			fmt.Println("Waiting for events")
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
				}
				if event.Op&fsnotify.Create == fsnotify.Create {
					log.Println("created file:", event.Name)
					// Do something here when a file is created
					if strings.HasSuffix(event.Name, ".csv") {
						// Assign the result of removeTopNLines to a variable
						stockCsvPath, _ := removeTopNLines(event.Name, "Name")
						readStockAloneCsv(stockCsvPath)
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	// TODO: get from flags
	err = watcher.Add("./tmp")
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

type Investment struct {
	CurrentValue string
	TotalReturns string
}

type Stock struct {
	Name         string
	Ticker       string
	CurrentPrice string
	AvgBuyPrice  string
	Returns      string
	Weightage    string
	Shares       string
}

func readFile(filepath string) {
	// Read the file
	file, err := os.Open(filepath)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Create arrays to store the structs
	var investments []Investment
	var stocks []Stock

	// Read the file line by line
	scanner := bufio.NewScanner(file)
	stock_info := false
	for scanner.Scan() {
		line := scanner.Text()

		// // Check if the line starts with "Investment"
		// if strings.HasPrefix(line, "Investment") {
		// 	// Skip the header line
		// 	continue
		// }

		// Check if the line starts with "Name"
		if strings.HasPrefix(line, "Name") {
			// Skip the header line
			stock_info = true
			continue
		}

		if stock_info == false {
			continue
		}

		// Split the line by comma
		fields := strings.Split(line, ",")

		// Create an Investment struct and populate the fields
		// investment := Investment{
		// 	CurrentValue: fields[0],
		// 	TotalReturns: fields[1],
		// }

		// Create a Stock struct and populate the fields
		stock := Stock{
			Name:         fields[0],
			Ticker:       fields[1],
			CurrentPrice: fields[2],
			AvgBuyPrice:  fields[3],
			Returns:      fields[4],
			Weightage:    fields[5],
			Shares:       fields[6],
		}

		// Append the structs to the arrays
		// investments = append(investments, investment)
		fmt.Println(stock)
		stocks = append(stocks, stock)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Print the arrays
	fmt.Println("Investments:")
	for _, investment := range investments {
		fmt.Println(investment)
	}

	fmt.Println("Stocks:")
	for _, stock := range stocks {
		fmt.Println(stock)
	}
	//
	// Investment Overview,,,,,,Returns Breakdown

	// Current Value,,,,,,Total Returns
	// 21991.85,,,,,,29.20
	// Current Investment,Money Put In,,,Current Returns,Realized Returns,Dividends
	// 21962.65,21962.65,,,29.20,0.00,0.00

	// Name,Ticker,Current Price (Rs.),Avg Buy Price (Rs.),Returns (%),Weightage,Shares
	// Union Bank of India Ltd,UNIONBANK,149.00,148.65,0.23,9.48,14
	// Tata Power Company Ltd,TATAPOWER,392.20,392.25,-0.01,19.61,11
	// Nava Limited,NAVA,483.35,483.25,0.02,15.38,7
	// Jindal Stainless Ltd,JSL,694.90,693.97,0.13,15.79,5
	// JK Tyre & Industries Ltd,JKTYRE,416.40,415.16,0.29,9.46,5
	// Jindal SAW Ltd,JINDALSAW,437.45,437.10,0.08,15.91,8
	// Jammu and Kashmir Bank Ltd,J&KBANK,131.35,130.90,0.34,14.33,24

}

func removeTopNLines(filepath string, lineSuffix string) (string, error) {
	// Open the file
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	outputFilename := "results/" + filepath + "_stock.csv"

	// Create a new file
	newFile, err := os.Create(outputFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer newFile.Close()

	// Create a scanner to read the file
	scanner := bufio.NewScanner(file)
	writer := bufio.NewWriter(newFile)

	stock_info := false
	for scanner.Scan() {
		line := scanner.Text()

		// // Check if the line starts with "Investment"
		// if strings.HasPrefix(line, "Investment") {
		// 	// Skip the header line
		// 	continue
		// }

		// Check if the line starts with "Name"
		if strings.HasPrefix(line, "Name") {
			// Skip the header line
			stock_info = true
		}

		if stock_info == false {
			continue
		}

		writer.WriteString(scanner.Text() + "\n")
	}

	// Check for errors in the scanner
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Flush the writer
	writer.Flush()

	return outputFilename, nil
}

type StockDetailed struct {
	Name         string
	Ticker       string
	CurrentPrice float64
	AvgBuyPrice  float64
	Returns      float64
	Weightage    float64
	Shares       int
}

func readStockAloneCsv(filename string) ([]StockDetailed, error) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	var stocks []StockDetailed
	for _, line := range lines[1:] { // Skip the header line
		currentPrice, _ := strconv.ParseFloat(line[2], 64)
		avgBuyPrice, _ := strconv.ParseFloat(line[3], 64)
		returns, _ := strconv.ParseFloat(line[4], 64)
		weightage, _ := strconv.ParseFloat(line[5], 64)
		shares, _ := strconv.Atoi(line[6])

		stocks = append(stocks, StockDetailed{
			Name:         line[0],
			Ticker:       line[1],
			CurrentPrice: currentPrice,
			AvgBuyPrice:  avgBuyPrice,
			Returns:      returns,
			Weightage:    weightage,
			Shares:       shares,
		})
	}

	for _, stock := range stocks {
		log.Println(stock)
	}

	return stocks, nil
}
