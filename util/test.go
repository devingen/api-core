package util

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/devingen/api-core/database"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func InsertDataFromFile(db *database.Database, databaseName string, collectionName string) {

	log.Println("Preparing collection:", collectionName)
	var list []interface{}
	ReadFile("./data/"+collectionName+".json", &list)

	for _, item := range list {
		InsertData(db, databaseName, collectionName, item)
	}
	return
}

func ReadFile(name string, list interface{}) {
	jsonFile, err := os.Open(name)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &list)
}

func InsertData(db *database.Database, databaseName string, collectionName string, data interface{}) {
	collection, err := db.ConnectToCollection(databaseName, collectionName)
	if err != nil {
		return
	}

	collection.Drop(context.TODO())
	_, err = collection.InsertOne(context.TODO(), data)
	if err != nil {
		return
	}
}

func SaveResultFile(name string, data interface{}) {
	filePrefix, _ := filepath.Abs("./result")
	f, err := os.Create(filePrefix + "/" + name + ".json")
	if err != nil {
		fmt.Println(1, err.Error())
		return
	}

	json, _ := json.MarshalIndent(data, "", "  ")
	_, err = f.Write(json)
	if err != nil {
		fmt.Println(2, err.Error())
		return
	}
	//err := ioutil.WriteFile("./result/"+name+".json", file, 0644)
}
