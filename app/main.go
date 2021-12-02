package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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

	db, err := gorm.Open(mysql.Open("user:pass@tcp(localhost:3306)/mydb?charset=utf8mb4&parseTime=True&loc=Local"))
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect to database")
	}

	var people []Person
	db.Find(&people)
	fmt.Println(people)

	js, err := json.Marshal(people)
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

	db, err := gorm.Open(mysql.Open("user:pass@tcp(localhost:3306)/mydb?charset=utf8mb4&parseTime=True&loc=Local"))
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect to database")
	}

	vars := mux.Vars(r)

	firstName := vars["FirstName"]
	surname := vars["SurName"]
	age, err := strconv.ParseInt(vars["Age"], 10, 64)
	if err != nil {
		fmt.Println(err.Error())
	}
	salary, err := strconv.ParseInt(vars["Salary"], 10, 64)
	if err != nil {
		fmt.Println(err.Error())
	}

	db.Create(&Person{FirstName: firstName, SurName: surname, Age: int16(age), Salary: salary})

	decoder := json.NewDecoder(r.Body)
	var staffMember Person
	err = decoder.Decode(&staffMember)
	if err != nil {
		panic(err)
	}

	fmt.Println(decoder)
	fmt.Println(vars["FirstName"])
	// fmt.Println(firstName)
	// fmt.Println(surname)
	// fmt.Println(salary)
	// fmt.Println(age)

	// err := json.NewDecoder(r.Body).Decode(&staffMember)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }

	// people = append(people, staffMember)

	// fmt.Fprintf(w, "Staff Member Added: %+v", staffMember)
	// fmt.Fprintf(w, "Staff Members: %+v", people)

	fmt.Println("Post made")
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("Hit the home route"))
	fmt.Fprintf(w, "Hi")

	// insert, err := db.Query("INSERT INTO User VALUES ('Joseph')")
	// if err != nil {
	// 	panic(err)
	// }

	// defer insert.Close()

}

func initialMigration() {
	db, err := gorm.Open(mysql.Open("user:pass@tcp(localhost:3306)/mydb?charset=utf8mb4&parseTime=True&loc=Local"))
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect to database")
	}

	db.AutoMigrate(&Person{})
}

func main() {

	initialMigration()

	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler).Methods("GET")
	r.HandleFunc("/staff", GetAllStaff).Methods("GET")
	r.HandleFunc("/staff/{name}", GetStaffByName).Methods("GET")
	r.HandleFunc("/staff", AddStaff).Methods("POST")

	fmt.Println("Server started")

	log.Fatal(http.ListenAndServe(":8000", r))

}
