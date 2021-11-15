package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Person struct {
	FirstName string
	SurName   string
	Age       int16
	Salary    int64
}

var people = []Person{
	{
		FirstName: "Joe",
		SurName:   "O'Reilly",
		Age:       29,
		Salary:    35000,
	},
	{
		FirstName: "Chloe",
		SurName:   "Stanley",
		Age:       23,
		Salary:    40000,
	},
	{
		FirstName: "James",
		SurName:   "Hutchinson",
		Age:       34,
		Salary:    45000,
	},
}

func GetAllStaff(w http.ResponseWriter, r *http.Request) {
	staff := people

	js, err := json.Marshal(staff)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func GetStaffByName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	nameQuery := vars["name"]

	for i := 0; i < len(people); i++ {
		if people[i].FirstName == nameQuery {
			res, err := json.Marshal(people[i])
			if err != nil {
				panic(err.Error())
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(res)
		}
	}

}

func AddStaff(w http.ResponseWriter, r *http.Request) {

	var staffMember Person

	err := json.NewDecoder(r.Body).Decode(&staffMember)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	people = append(people, staffMember)

	fmt.Fprintf(w, "Staff Member Added: %+v", staffMember)
	fmt.Fprintf(w, "Staff Members: %+v", people)

	fmt.Println("Post made")
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hit the home route"))
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler).Methods("GET")
	r.HandleFunc("/staff", GetAllStaff).Methods("GET")
	r.HandleFunc("/staff/{name}", GetStaffByName).Methods("GET")
	r.HandleFunc("/staff", AddStaff).Methods("POST")
	fmt.Println("Server started")

	log.Fatal(http.ListenAndServe(":8000", r))
}
