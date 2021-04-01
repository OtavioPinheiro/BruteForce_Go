package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/yeka/zip"
)

func main() {

	zipPath := flag.String("zip", "", "Path to the zip file")
	passwordListFile := flag.String("password-list", "", "Path to the password list file")
	threads := flag.Int("threads", 1, "Number of threads")
	linesPerThread := flag.Int("lines-thread", 1000, "Number of lines to be read by each thread")
	timeout := flag.Int("timeout", 60, "Timeout for the proccess - in seconds")
	flag.Parse()

	unzipFile(*zipPath, *passwordListFile, *threads, *linesPerThread, *timeout)
}

func unzipFile(zipPath string, passwordListPath string, threads int, linesPerThread int, timeout int) {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	zipFile := r.File[0]

	if zipFile.IsEncrypted() {
		fmt.Printf("The file is protected by password! Let's crack it down...\n")
		ch := make(chan string)
		pwdList := getPasswordList(passwordListPath)
		initialIndex := 0
		start := time.Now()
		for i := 0; i < threads; i++ {
			finalIndex := linesPerThread * (i + 1)

			outputColor := RandomOutputColor(i)
			output := color.New(outputColor)
			output.Printf("Starting thread %d reading from line %d till line %d\n", i+1, initialIndex, finalIndex)
			go bruteForce(zipPath, pwdList[initialIndex:finalIndex], outputColor, ch)

			initialIndex = finalIndex + 1
		}
		fmt.Println("----------------------------------------------------")

		color.Yellow("Cracking the password...")
		fmt.Println()

		select {
		case pwd := <-ch:
			fmt.Printf("\nThe password is:\"%v\"\n", pwd)
			fmt.Printf("\nThe cracking process took %v \n", time.Since(start))
		case <-time.After(time.Duration(timeout) * time.Second):
			fmt.Printf("Timeout after: %d seconds \n", timeout)
			fmt.Printf("Password not found :( \n")
		}

	} else {
		fmt.Printf("No protection...\n")
	}
}

func bruteForce(zipPath string, passwordList []string, outputColor color.Attribute, ch chan<- string) {
	output := color.New(outputColor)
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	zipFile := r.File[0]

	for _, value := range passwordList {
		//	output.Printf("Trying to crack the file with password: %v \n", string(value))

		zipFile.SetPassword(string(value))
		_, err := zipFile.Open()

		if err == nil {
			output.Printf("Ulala, we found the password!\n")

			zipReader, err := zipFile.Open()

			if err != nil {
				log.Fatal(err)
			}

			buf, err := ioutil.ReadAll(zipReader)
			if err != nil {
				log.Fatal(err)
			}

			defer zipReader.Close()

			output.Printf("Size of %v: %v byte(s)\n", zipFile.Name, len(buf))

			ch <- string(value)
			break
		}
	}
}

func RandomOutputColor(threadId int) color.Attribute {

	switch threadId {
	case 0:
		return color.FgHiBlue
	case 1:
		return color.FgRed
	case 2:
		return color.FgCyan
	case 3:
		return color.FgGreen
	case 4:
		return color.FgHiYellow
	case 5:
		return color.FgMagenta
	}

	return color.FgBlue
}

func getPasswordList(passwordListFile string) []string {
	file, err := os.Open(passwordListFile)

	if err != nil {
		log.Fatalf("failed to open")
	}

	scanner := bufio.NewScanner(file)

	scanner.Split(bufio.ScanLines)
	var text []string

	for scanner.Scan() {
		text = append(text, scanner.Text())
	}

	file.Close()

	wordListStat, _ := os.Stat(passwordListFile)
	fmt.Print("-----------------------------\n", "List file: ", wordListStat.Name(), "\nFile size: ", wordListStat.Size()/(1024), " KB\n")
	fmt.Println("Total passwords: ", len(text))
	fmt.Println("-----------------------------")

	return text
}
