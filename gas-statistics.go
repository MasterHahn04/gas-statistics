package main

import (
	"encoding/json"
	"fmt"
	"gas-statistics/pkg/request"
	"gas-statistics/pkg/storage"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"
)

func main() {
	//TODO Request every 5 Minutes the Prices of the Stations which are given via there ID
	//TODO Add random seconds of delay to the request to avoid the blocking of the API
	//TODO Ensure that the request not made at a full hour to avoid the blocking of the API
	//TODO Ensure that the request not made at a full second to avoid the blocking of the API e.g. XX:X4

	//TODO One Time per Day request the list of the Stations to hold the stations Table up to date

	//TODO Save the Prices in a Database
	//TODO Ensure to not store the same price twice beyond each other

	//Features for the Future (not in the first Version)
	//TODO Calculate the Gas-Statistics
	//TODO Calculate the three cheapest times to refuel for each station

	//--------------------------------

	//Implement the interfaces
	//TODO Create the Functions that it can easily be switch to proceed the Data of the MTK-S
	//TODO Create the Functions that it can easily be switch to proceed the Data of the Tankerkoenig
	//^^^ not in the first Version ^^^

	//Tankerkoenig Functions
	//TODO Create the Function to request the Prices of the Stations via Tankerkoenig by the prices.php
	//TODO Create the Function to request the List of the Stations via Tankerkoenig by the list.php
	//For the Tankerkoenig Api must be the API-Key set in the Environment Variables
	//For the Tankerkoenig Api must be the ID, the lng and the lat of max 10 stations set in the Database

	//Database Functions
	//TODO Create the Function to store the Prices in the Database
	//TODO Create the Function to store the Stations Details in the Database
	//TODO Create the Function to store the whole Respond in the Database

	//TODO Load the Environment Variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	databaseType, used := os.LookupEnv("DATABASE_TYPE")
	if used != true {
		fmt.Println("The Environment Variable DATABASE_TYPE is not set\n")
	}
	fmt.Printf("The Database Type is: %s\n", databaseType)
	var s storage.Storage
	//Switch Case to choose the Database
	switch strings.ToLower(databaseType) {
	case "mariadb", "mysql":
		var db storage.DatabaseMariaDB
		db.Host, used = os.LookupEnv("DB_HOST")
		if used != true {
			fmt.Println("The Environment Variable DB_HOST is not set\n")
		}
		db.Port, used = os.LookupEnv("DB_PORT")
		if used != true {
			fmt.Println("The Environment Variable DB_PORT is not set\n")
		}
		db.User, used = os.LookupEnv("DB_USER")
		if used != true {
			fmt.Println("The Environment Variable DB_USER is not set\n")
		}
		db.Password, used = os.LookupEnv("DB_PASSWORD")
		if used != true {
			fmt.Println("The Environment Variable DB_PASSWORD is not set\n")
		}
		db.Database, used = os.LookupEnv("DB_DATABASE")
		if used != true {
			fmt.Println("The Environment Variable DB_DATABASE is not set\n")
		}
		s = &db
	}
	err = s.Connect()
	if err != nil {
		fmt.Printf("The Connection to the Database failed\n Error: %s\n", err)
	}

	//--------------------------------

	var r request.Request
	//Switch Case to choose the Request
	source, used := os.LookupEnv("SOURCE")
	if used != true {
		fmt.Println("The Environment Variable SOURCE is not set\n")
	}
	fmt.Printf("The Source is: %s\n", source)
	switch strings.ToLower(source) {
	case "tankerkoenig":
		var src request.RequestTankerkoenig
		//Key for the access to the free Tankerkoenig-Gas-Price-API
		//For own Key please register here https://creativecommons.tankerkoenig.de
		src.ApiKey, used = os.LookupEnv("TANKERKOENIG_API_KEY")
		if used != true {
			fmt.Println("The Environment Variable TANKERKOENIG_API_KEY is not set\n")
		}
		fmt.Printf("The Tankerkoenig API-Key is: %s\n", src.ApiKey)
		r = &src
	}
	//Test
	log.Println(r)
	/*
		ids, err := s.GetIDs()
		if err != nil {
			fmt.Printf("The Request to the Database failed\n Error: %s\n", err)
		}
	*/
	/*
		err = r.MakeRequest(ids)
		if err != nil {
			fmt.Printf("The Request to the Source failed\n Error: %s\n", err)
		}
	*/
	test := new(request.PricesRespond)
	//TODO Undo the Test
	err = json.NewDecoder(strings.NewReader(rohdaten)).Decode(test)
	if err != nil {
		fmt.Printf("The Decoding of the Respond failed\n Error: %s\n", err)
	}
	fmt.Printf("The Respond is: %v\n", test)
}
