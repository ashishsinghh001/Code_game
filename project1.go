package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
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
var s []string
var a []string
var count int = 0
var z []string

/*
Challenge 0 to build an api:
Below code will print maximum power of character on API
http://localhost:10000/ReturnPower    */

func ReturnPower(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "******Welcome in Marvels Game!******")
	fmt.Fprintln(w, "=================================================")
	fmt.Fprintln(w, "Please find character maximum power: ")
	fmt.Println("Endpoint Hit: ReturnPower")
	json.NewEncoder(w).Encode(Powers)
}

func handleRequests() {
	http.HandleFunc("/ReturnPower", ReturnPower)
	log.Fatal(http.ListenAndServe(":10000", nil))
}

/*
Power level changes at every 10s interval. I have tried using time.Tick(10*time.Second) in loop to refresh data
but faced few issues and error. Therefore not using it.
*/

func main() {
	response1, err := http.Get("http://www.mocky.io/v2/5ecfd5dc3200006200e3d64b")
	response2, err := http.Get("http://www.mocky.io/v2/5ed39e2b340000e46801f3c6")
	response3, err := http.Get("http://www.mocky.io/v2/5ecfd6473200009dc1e3d64e")
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	mapdata := make(map[string]int) // API data is stored in map data structure
	newmap := make(map[string]int)  // this map will help in finding least used character

	response1Data, err := ioutil.ReadAll(response1.Body) // this line will read given Avengers API
	response2Data, err := ioutil.ReadAll(response2.Body) // this line will read given Anti heroes API
	response3Data, err := ioutil.ReadAll(response3.Body) // this line will read given Mutant API
	if err != nil {
		log.Fatal(err)
	}

	var response1Object Response
	json.Unmarshal(response1Data, &response1Object)

	var response2Object Response
	json.Unmarshal(response2Data, &response2Object)

	var response3Object Response
	json.Unmarshal(response3Data, &response3Object)

	for i := 0; i < len(response1Object.Character); i++ { //this loop will help to store Avengers data in map
		mapdata[response1Object.Character[i].Name] = response1Object.Character[i].MaxPower // <Key, Value >(Character name, Power)
		newmap[response1Object.Character[i].Name] = count                                  //setting value of newmap as 0 initially
		s = append(s, response1Object.Character[i].Name)                                   //this line is for appending all Charaters name in array
	}

	for i := 0; i < len(response2Object.Character); i++ { //this loop will help to store Anti heroes data in map
		mapdata[response2Object.Character[i].Name] = response2Object.Character[i].MaxPower
		newmap[response2Object.Character[i].Name] = count
		s = append(s, response2Object.Character[i].Name)
	}

	for i := 0; i < len(response3Object.Character); i++ { //this loop will help to store Mutants data in map
		mapdata[response3Object.Character[i].Name] = response3Object.Character[i].MaxPower
		newmap[response3Object.Character[i].Name] = count
		s = append(s, response3Object.Character[i].Name)
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

	/*Challenge 1: If there more characters then we should remove the least used characters.
	Steps: 1- I have used newmap to store all character names as key and values as "0"
		   2- When user enters its character then that character value will be incremented by 1
		   3- This way least used character will have smaller value
		   4- sort the newmap using value
		   5- Append all key in "z" array and delete the element at index "0"
	*/
	if len(mapdata) > 15 {
		// used to switch key and value
		hack := map[int]string{}
		hackkeys := []int{}
		for key, val := range newmap {
			hack[val] = key
			hackkeys = append(hackkeys, val)
		}
		sort.Ints(hackkeys)

		for _, val := range hackkeys {
			z = append(z, hack[val])
		}

		fmt.Println("Deleting least used character: " + z[0])
		//delete(mapdata, z[0])

		fmt.Println("Done...")
	}

	fmt.Println("Please choose the character from below given list:")
	fmt.Println("==========================================================")
	fmt.Println("List of Characters fetched from given APIs:  ")
	fmt.Println(s) // this array has all character names before deleting least used and min power
	fmt.Println("==========================================================")
	fmt.Println("Enter your character name to know its maximum power except deleted character : ")

	scanner := bufio.NewScanner(os.Stdin)
	var n string
	if scanner.Scan() {
		n = scanner.Text() //taking user input from console
	}

	value, okay := newmap[n]
	if okay { //this line will increment value of newmap
		newmap[n] = count + 1
	}

	value, ok := mapdata[n]
	if ok {
		fmt.Println("Please visit this API http://localhost:10000/ReturnPower to know maximum power of " + n)
	} else {
		fmt.Println("Character not found")
	}

	Powers = []Power{
		Power{Name: n, Power: value},
	}

	handleRequests() //calling handleRequest() function

}
