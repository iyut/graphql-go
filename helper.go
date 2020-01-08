package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func openJSONFile() Settings {

	jsonFile, err := os.Open(getMainPath() + "settings.json")

	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var settings Settings

	json.Unmarshal(byteValue, &settings)

	/*
		if err := json.Unmarshal([]byte(settings), &val); err != nil {
			panic(err)
		}
	*/

	return settings
}

func getMainPath() string {

	pathLocal := "/Users/luthfi/go/bin/"
	pathServer := "/root/go/bin/"

	if _, err := os.Stat(pathLocal); !os.IsNotExist(err) {
		return pathLocal
	}

	return pathServer
}
