package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
)

// Creates a workshop.lua file with commands to download the workshop collection
func main() {

	apiKey := ""
	collectionId := 0
	out := ""

	flag.StringVar(&apiKey, "api-key", "", "Your steam api key")
	flag.IntVar(&collectionId, "collection", 0, "The workshop collection id")
	flag.StringVar(&out, "output", "", "Your steam api key")

	flag.Parse()

	if apiKey == "" {
		fmt.Println("Invalid api key")
		return
	}

	if collectionId <= 0 {
		fmt.Println("Invalid collection id")
		return
	}

	if out == "" {
		fmt.Println("Invalid output file")
		return
	}

	files, err := download(apiKey, collectionId)
	if err != nil {
		panic(err)
	}

	file, err := openLua(out)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = file.Close()
	}()

	writer := bufio.NewWriter(file)
	writeLua(writer, collectionId, files)

	_ = writer.Flush()
}

func download(apiKey string, collectionId int) ([]WorkshopFile, error) {
	fileIds, err := requestCollection(apiKey, collectionId)
	if err != nil {
		return nil, err
	}

	files, err := requestFileDetails(apiKey, fileIds)
	if err != nil {
		return nil, err
	}

	return files, nil
}

func writeLua(s *bufio.Writer, collectionId int, files []WorkshopFile) {
	_, _ = s.WriteString(
		"-- Created with workshop-lua (https://github.com/LukWebsForge/LuaWorkshopGen)\n" +
			"-- List based on collection https://steamcommunity.com/sharedfiles/filedetails/?id=" + strconv.Itoa(collectionId) +
			"\n" +
			"\n")

	for _, file := range files {
		_, _ = s.WriteString("resource.AddWorkshop(\"" + strconv.Itoa(file.Id) + "\") -- " + file.Name + "\n")
	}
}

func openLua(filename string) (*os.File, error) {
	file, err := os.Create(filename)
	if err != nil {
		return nil, err
	}

	return file, nil
}
