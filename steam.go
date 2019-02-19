package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

const CollectionUrl string = "https://api.steampowered.com/ISteamRemoteStorage/GetCollectionDetails/v1/"
const FileDetailsUrl string = "https://api.steampowered.com/ISteamRemoteStorage/GetPublishedFileDetails/v1/"

type WorkshopFile struct {
	Id   int
	Name string
}

func requestSteam(data url.Values, url string) (*http.Response, error) {
	data.Set("format", "json")

	response, err := http.PostForm(url, data)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, fmt.Errorf("invalid http status: %v", response.StatusCode)
		} else {
			return nil, fmt.Errorf("invalid http status: %v\n body: %v", response.StatusCode, string(body))
		}
	}

	return response, nil
}

func requestCollection(apiKey string, collectionId int) ([]int, error) {
	data := url.Values{}
	data.Set("key", apiKey)
	data.Set("collectioncount", "1")
	data.Set("publishedfileids[0]", strconv.Itoa(collectionId))

	response, err := requestSteam(data, CollectionUrl)
	if err != nil {
		return nil, err
	}

	return parseCollections(&response.Body, collectionId)
}

func requestFileDetails(apiKey string, fileIds []int) ([]WorkshopFile, error) {
	data := url.Values{}
	data.Set("key", apiKey)
	data.Set("itemcount", strconv.Itoa(len(fileIds)))
	for key, fileId := range fileIds {
		data.Set("publishedfileids[" + strconv.Itoa(key) + "]", strconv.Itoa(fileId))
	}

	response, err := requestSteam(data, FileDetailsUrl)
	if err != nil {
		return nil, err
	}

	return parseFileDetails(&response.Body)
}

func parseSteamResp(stream *io.ReadCloser, category string) ([]interface{}, error) {
	jsonDec := json.NewDecoder(*stream)
	jsonCon := new(map[string]interface{})

	err := jsonDec.Decode(&jsonCon)
	if err != nil {
		return nil, err
	}

	response := (*jsonCon)["response"].(map[string]interface{})
	resultCount := response["resultcount"].(float64)

	if resultCount <= 0 {
		return nil, fmt.Errorf("zero collection count")
	}

	details := response[category+"details"].([]interface{})
	return details, nil
}

func parseCollections(stream *io.ReadCloser, wantedCollectionId int) ([]int, error) {
	collections, err := parseSteamResp(stream, "collection")
	if err != nil {
		return nil, err
	}

	collection := collections[0].(map[string]interface{})

	collectionId, err := strconv.Atoi(collection["publishedfileid"].(string))
	if err != nil {
		return nil, fmt.Errorf("collection id isn't a int: %v", collection["publishedfileid"])
	} else if collectionId != wantedCollectionId {
		return nil, fmt.Errorf("wanted and given collections ids don't match: %v / %v", wantedCollectionId, collectionId)
	}

	children := collection["children"].([]interface{})
	fileIds := make([]int, len(children))

	for key, child := range children {
		childMap := child.(map[string]interface{})
		fileIdStr := childMap["publishedfileid"].(string)

		i, err := strconv.Atoi(fileIdStr)
		if err != nil {
			fmt.Printf("error while reading: %v is not a int", fileIdStr)
			i = -1
		}

		fileIds[key] = i
	}

	return fileIds, nil
}

func parseFileDetails(stream *io.ReadCloser) ([]WorkshopFile, error) {
	publishedFiles, err := parseSteamResp(stream, "publishedfile")
	if err != nil {
		return nil, err
	}

	workFiles := make([]WorkshopFile, len(publishedFiles))

	for key, pubFile := range publishedFiles {
		pubMap := pubFile.(map[string]interface{})
		fileIdStr := pubMap["publishedfileid"].(string)
		title := pubMap["title"].(string)

		fid, err := strconv.Atoi(fileIdStr)
		if err != nil {
			fmt.Printf("error while reading: %v is not a int", fileIdStr)
			fid = -1
		}

		workFiles[key] = WorkshopFile{
			Id:fid,
			Name:title,
		}
	}

	return workFiles, nil
}