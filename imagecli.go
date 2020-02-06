package main

import (
	"encoding/csv"
	"fmt"
	"github.com/corona10/goimagehash"
	"github.com/urfave/cli/v2"
	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
	"golang.org/x/image/webp"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
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
	// write to csv
	writeResults(file+"results.csv", rows)

	// satisfy cli action reqs
	return nil

}

func compareImages(image1 string, image2 string) []float64 {
	// Start timer
	start := time.Now()

	// calculate average hash
	avg1, _ := goimagehash.AverageHash(imageDecoded(image1))
	avg2, _ := goimagehash.AverageHash(imageDecoded(image2))
	distance1, _ := avg1.Distance(avg2)

	// calculate difference hash
	diff1, _ := goimagehash.DifferenceHash(imageDecoded(image1))
	diff2, _ := goimagehash.DifferenceHash(imageDecoded(image2))
	distance2, _ := diff1.Distance(diff2)

	// pick the greater one
	distance := 0

	if distance1 > distance2 {
		distance = distance1
	} else {
		distance = distance2
	}
	// set it
	elapsed := float64(time.Since(start) / time.Millisecond)
	values := make([]float64, 2)
	values[0] = float64(distance) / 100
	values[1] = elapsed
	return values
}

// IMAGE HELPERS
type FileType int

const (
	PNG FileType = iota
	JPG
	GIF
	WEBP
	BMP
	TIFF
	ERR
)

func getFileType(input string) FileType {
	switch input {
	case "jpg":
		fallthrough
	case "jpeg":
		return JPG
	case "png":
		return PNG
	case "gif":
		return GIF
	case "bmp":
		return BMP
	case "webp":
		return WEBP
	case "tiff":
		return TIFF
	default:
		return ERR
	}
}

func imageDecoded(image string) image.Image {

	// handle opening and closing of the images
	file, err := os.Open(image)
	if err != nil {
		log.Fatalf("Cannot open image '%s': %s\n", file, err.Error())
	}
	defer file.Close()

	// get extension type
	ext := strings.ToLower(filepath.Ext(image))
	// validate file type
	startType := getFileType(ext[1:])
	if startType == ERR {
		log.Fatalf("file input type not valid")
	}

	// decode
	if startType == JPG {
		img, _ := jpeg.Decode(file)
		if err != nil {
			log.Fatalf("Cannot decode '%s': %s\n", img, err.Error())
		}
		return img

	} else if startType == PNG {
		img, _ := png.Decode(file)
		if err != nil {
			log.Fatalf("Cannot decode '%s': %s\n", img, err.Error())
		}
		return img
	} else if startType == GIF {
		img, _ := gif.Decode(file)
		if err != nil {
			log.Fatalf("Cannot decode '%s': %s\n", img, err.Error())
		}
		return img
	} else if startType == BMP {
		img, _ := bmp.Decode(file)
		if err != nil {
			log.Fatalf("Cannot decode '%s': %s\n", img, err.Error())
		}
		return img
	} else if startType == WEBP {
		img, _ := webp.Decode(file)
		if err != nil {
			log.Fatalf("Cannot decode '%s': %s\n", img, err.Error())
		}
		return img
	} else if startType == TIFF {
		img, _ := tiff.Decode(file)
		if err != nil {
			log.Fatalf("Cannot decode '%s': %s\n", img, err.Error())
		}
		return img
	}

	return nil
}

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
