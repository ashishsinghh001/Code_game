package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/gorilla/mux"
)

type Response struct {
	Name      string      `json:"name"`
	Character []Character `json:"character"`
}

type Character struct {
	Name     string `json:"name"`
	MaxPower int    `json:"max_power"`
}

type Power struct {
	Name  string `json:"Name"`
	Power int    `json:"Power"`
}

var Powers []Power

var a []string

/*
Challenge 0 to build an api:
http://localhost:10000/ReturnPower    */

func ReturnPower(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "******Welcome in Marvels Game!******")
	fmt.Fprintln(w, "=================================================")
	fmt.Fprintln(w, "Please find all characters : ")
	fmt.Println("Endpoint Hit: ReturnPower")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Powers)
}

func CharacterName(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Gets params
	// Loop through Powers and find one with the Name from the params
	for _, item := range Powers {
		if item.Name == params["Name"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Power{})
}

func deleteChar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range Powers {
		if item.Name == params["Name"] {
			Powers = append(Powers[:index], Powers[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(Powers)
}

func handleRequests() {
	myRouter := mux.NewRouter()
	myRouter.HandleFunc("/ReturnPower", ReturnPower).Methods("GET")
	myRouter.HandleFunc("/ReturnPower/{Name}", CharacterName).Methods("GET")
	myRouter.HandleFunc("/deleteChar/{Name}", deleteChar).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {

	ticker := time.NewTicker(10 * time.Second)
	done := make(chan bool)

	var response1Data []byte
	var response2Data []byte
	var response3Data []byte

	go func() {
		for {
			response1, err := http.Get("http://www.mocky.io/v2/5ecfd5dc3200006200e3d64b")
			response2, err := http.Get("http://www.mocky.io/v2/5ed39e2b340000e46801f3c6")
			response3, err := http.Get("http://www.mocky.io/v2/5ecfd6473200009dc1e3d64e")
			if err != nil {
				fmt.Print(err.Error())
				os.Exit(1)
			}
			select {
			case <-done:
				return
			case t := <-ticker.C:
				fmt.Println("Tick at", t)
				response1Data, err = ioutil.ReadAll(response1.Body) // this line will read given Avengers API
				response2Data, err = ioutil.ReadAll(response2.Body) // this line will read given Anti heroes API
				response3Data, err = ioutil.ReadAll(response3.Body) // this line will read given Mutant API
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}()
	time.Sleep(10 * time.Second)
	done <- true

	mapdata := make(map[string]int) // API data is stored in map data structure

	var response1Object Response
	json.Unmarshal(response1Data, &response1Object)

	var response2Object Response
	json.Unmarshal(response2Data, &response2Object)

	var response3Object Response
	json.Unmarshal(response3Data, &response3Object)

	for i := 0; i < len(response1Object.Character); i++ { //this loop will help to store Avengers data in map
		mapdata[response1Object.Character[i].Name] = response1Object.Character[i].MaxPower // <Key, Value >(Character name, Power)
		Powers = append(Powers, Power{Name: response1Object.Character[i].Name, Power: response1Object.Character[i].MaxPower})
	}

	for i := 0; i < len(response2Object.Character); i++ { //this loop will help to store Anti heroes data in map
		mapdata[response2Object.Character[i].Name] = response2Object.Character[i].MaxPower
		Powers = append(Powers, Power{Name: response2Object.Character[i].Name, Power: response2Object.Character[i].MaxPower})
	}

	for i := 0; i < len(response3Object.Character); i++ { //this loop will help to store Mutants data in map
		mapdata[response3Object.Character[i].Name] = response3Object.Character[i].MaxPower
		Powers = append(Powers, Power{Name: response3Object.Character[i].Name, Power: response3Object.Character[i].MaxPower})
	}

	/*
		Challenge 2:Drop lowest powered characters to make space.
		Steps:1- Sorted the mapdata by value and saved sorted value in new map hack.
			  2- Appending key of hack into "a" array
			  3- Total characters are 18. Therefore deleting 3 minimum power characters to maintain 15 size of map
	*/
	if len(mapdata) > 15 {
		// used to switch key and value
		hack := map[int]string{}
		hackkeys := []int{}
		for key, val := range mapdata {
			hack[val] = key
			hackkeys = append(hackkeys, val)
		}
		sort.Ints(hackkeys)

		for _, val := range hackkeys {
			a = append(a, hack[val])
		}

		for i := 0; i < 3; i++ {
			fmt.Println("Deleting minimum power character: " + a[i])
			delete(mapdata, a[i])
		}
		fmt.Println("Done...")
	}

	fmt.Println("Please visit this API http://localhost:10000/ReturnPower for list of characters")
	fmt.Println("==========================================================")
	fmt.Println("Enter character name after API link to know maximum power level of your character.  ")
	fmt.Println("Example: http://localhost:10000/ReturnPower/Thanos  ")
	fmt.Println("==========================================================")
	fmt.Println("If you want to delete any element then visit http://localhost:10000/deleteChar/{Name} followed by charater name")

	ticker.Stop()
	fmt.Println("Ticker stopped")

	handleRequests() //calling handleRequest() function
}
