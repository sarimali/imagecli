package main

import (
	"encoding/csv"
	"fmt"
	"github.com/corona10/goimagehash"
	"github.com/urfave/cli/v2"
	"image/jpeg"
	"log"
	"os"
	"sort"
	"strconv"
	"time"
)

func main() {
	app := cli.NewApp()
	app.Name = "imagecli"
	app.Usage = "Image CLI"
	app.Version = "2020.02.02"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:     "path",
			Value:    "~/example.csv",
			Usage:    "Required path to csv file",
			Required: true,
		},
	}

	app.Commands = []*cli.Command{
		{
			Name:    "compare",
			Aliases: []string{"c"},
			Usage:   "Use this to compare the files listed in the csv",
			Action:  compare,
		},
		{
			Name:    "list",
			Aliases: []string{"l"},
			Action:  listFilesToCompare,
			Usage:   "List files that will be compared",
			// TODO: Reverse csv list flag
			//Flags: []cli.Flag{
			//	&cli.StringFlag{
			//		Name:  "reverse",
			//		Usage: "reverse list",
			//	},
			//},
		},
	}
	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))
	err := app.Run(os.Args)
	if err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}

func compare(c *cli.Context) error {
	file := c.String("path")

	rows := verifyFile(file)
	image1 := make([]string, len(rows)-1)
	image2 := make([]string, len(rows)-1)

	for i, line := range rows {
		if i == 0 {
			rows[0] = append(rows[0], "Similar", "Elapsed")
			continue
		}
		image1[i-1] = line[0]
		image2[i-1] = line[1]

		result := compareImages(image1[i-1], image2[i-1])
		distance := strconv.FormatFloat(result[0], 'f', -1, 64)
		elapsed := strconv.FormatFloat(result[1], 'f', -1, 64)
		rows[i] = append(rows[i], distance)
		rows[i] = append(rows[i], elapsed)
	}

	writeResults(file+"results.csv", rows)

	return nil

}



func compareImages(image1 string, image2 string) []float64 {
	// Start timer
	start := time.Now()

	// handle opening and closing of the images
	file1, err := os.Open(image1)
	if err != nil {
		log.Fatalf("Cannot open image '%s': %s\n", file1, err.Error())
	}
	file2, err := os.Open(image2)
	if err != nil {
		log.Fatalf("Cannot open image '%s': %s\n", file2, err.Error())
	}
	defer file1.Close()
	defer file2.Close()

	// decode
	img1, _ := jpeg.Decode(file1)
	if err != nil {
		log.Fatalf("Cannot decode '%s': %s\n", img1, err.Error())
	}
	img2, _ := jpeg.Decode(file2)
	if err != nil {
		log.Fatalf("Cannot decode '%s': %s\n", img2, err.Error())
	}

	// calculate average hash
	avg1, _ := goimagehash.AverageHash(img1)
	avg2, _ := goimagehash.AverageHash(img2)
	distance1, _ := avg1.Distance(avg2)

	// calculate difference hash
	diff1, _ := goimagehash.DifferenceHash(img1)
	diff2, _ := goimagehash.DifferenceHash(img2)
	distance2, _ := diff1.Distance(diff2)

	// pick the greater one
	distance := 0

	if distance1 > distance2{
		distance = distance1
	} else {
		distance = distance2
	}
	// set it
	elapsed := float64(time.Since(start) / time.Millisecond)
	values := make([]float64, 2)
	values[0] = float64(distance)/100
	values[1] = elapsed
	return values
}

// IMAGE HELPERS



//CSV FILE HELPERS
// `verifyFile` takes a filename and returns a two-dimensional list of spreadsheet cells
// Additionally has checks to manage the file

func verifyFile(file string) [][]string {
	f, err := os.Open(file)
	//handle error here since its just a small cli tool
	//alternatively return
	if err != nil {
		log.Fatalf("Cannot open '%s': %s\n", file, err.Error())
	}
	//Successfully opened so we want to ensure its closed
	defer f.Close()

	r := csv.NewReader(f)
	// set the delimiter to be spaces like the assignment denotes
	// TODO: Ensure they are using spaces else reject the file
	r.Comma = ' '
	// Read the whole file at once
	// TODO: Handle large files

	rows, err := r.ReadAll()
	// Again, we check for any error,
	if err != nil {
		log.Fatalln("Cannot read CSV data:", err.Error())
	}
	// Successfully loaded file in
	// and finally we can return the rows.
	return rows
}

// Function to read the csv and show files being compared
func listFilesToCompare(c *cli.Context) error {

	file := c.String("path")

	rows := verifyFile(file)

	lines := rows
	image1 := make([]string, len(lines)-1)
	image2 := make([]string, len(lines)-1)

	for i, line := range lines {
		if i == 0 {
			// skip header line
			continue
		}
		image1[i-1] = line[0]
		image2[i-1] = line[1]
	}

	for i, _ := range image1 {
		fmt.Println(image1[i], "comparing with", image2[i])
	}

	return nil
}

func writeResults(name string, rows [][]string) {

	f, err := os.Create(name)
	if err != nil {
		log.Fatalf("Cannot create '%s': %s\n", name, err.Error())
	}

	defer func() {
		e := f.Close()
		if e != nil {
			log.Fatalf("Cannot close created '%s': %s\n", name, e.Error())
		}
	}()

	w := csv.NewWriter(f)
	err = w.WriteAll(rows)
}