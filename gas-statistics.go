package main

import (
	"fmt"
	"gas-statistics/pkg/request"
	"gas-statistics/pkg/storage"
	"github.com/joho/godotenv"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	//TODO Request every 5 Minutes the Prices of the Stations which are given via there ID
	//TODO Add random seconds of delay to the request to avoid the blocking of the API
	//TODO Ensure that the request not made at a full hour to avoid the blocking of the API
	//TODO Ensure that the request not made at a full second to avoid the blocking of the API e.g. XX:X4

	//TODO One Time per Day request the list of the Stations to hold the stations Table up to date

	//TODO Ensure to not store the same price twice beyond each other

	//Features for the Future (not in the first Version)
	//TODO Calculate the Gas-Statistics
	//TODO Calculate the three cheapest times to refuel for each station

	//--------------------------------

	//Implement the interfaces
	//TODO Create the Functions that it can easily be switch to proceed the Data of the MTK-S
	//^^^ not in the first Version ^^^

	//Tankerkoenig Functions
	//TODO Create the Function to request the List of the Stations via Tankerkoenig by the list.php
	//For the Tankerkoenig Api must be the API-Key set in the Environment Variables
	//For the Tankerkoenig Api must be the ID, the lng and the lat of max 10 stations set in the Database

	//Database Functions
	//TODO Create the Function to store the Stations Details in the Database

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

	//Circle the following Code every 5 Minutes and add a random delay of 1-50 seconds to the 5 Minutes that it will never run on the full minute
	for {
		//Connect to the Database
		err = s.Connect()
		if err != nil {
			fmt.Printf("The Connection to the Database failed\n Error: %s\n", err)
		}

		//--------------------------------

		var r request.Request

		//Lookup what the Content of the Environment Variable SOURCE is
		source, used := os.LookupEnv("SOURCE")
		if used != true {
			fmt.Println("The Environment Variable SOURCE is not set\n")
		}

		//Print which Source is used
		fmt.Printf("The Source is: %s\n", source)

		//Switch Case to choose the Request
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
		ids, err := s.GetIDs()
		if err != nil {
			fmt.Printf("The Request to the Database failed\n Error: %s\n", err)
		}

		//Make the Request to the Source, get the Prices and give it back
		prices, body, err := r.MakeRequest(ids)
		if err != nil {
			fmt.Printf("The Request to the Source failed\n Error: %s\n", err)
		}

		//Store the Response of the Request in the Database
		err = s.StoreResponse(body)
		if err != nil {
			fmt.Printf("The Storage of the Response failed\n Error: %s\n", err)
		}

		//Store the Prices of the Stations in the Database
		for index, value := range prices.Prices {
			err := s.StorePrices(index, value.Status, value.E5, value.E10, value.Diesel)
			if err != nil {
				fmt.Printf("The Storage of the Prices failed\n Error: >%s<\n", err)
			}
		}

		//Close the Connection to the Database
		err = s.Close()
		if err != nil {
			fmt.Printf("To close the Connection to the Database failed\n Error: %s\n", err)
		}

		//Print the Time when the for loop is finished
		fmt.Printf("The for loop is finished at: %s\n", time.Now().Format("2006-01-02 15:04:05"))

		//Wait 5 Minutes and add a random delay of 1-50 seconds to the 5 Minutes
		time.Sleep(time.Minute*5 + time.Second*time.Duration(rand.Intn(50)))

		//Check if the Time is a full Minute and add a random delay of 1-50 seconds to the 5 Minutes
		//If the Time is a full Minute, the Request will be blocked by the API
		if time.Now().Second() == 0 {
			time.Sleep(time.Second * 3)
		}
		if time.Now().Minute() == 0 {
			time.Sleep(time.Second * 23)
		}

		//Print the Time that the Loop will now start again
		fmt.Printf("The for loop will now start again at: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	}

}
