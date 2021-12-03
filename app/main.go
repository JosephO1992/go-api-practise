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
	ID        uint `gorm:"primaryKey;auto_increment;not_null"`
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
	db, err := gorm.Open(mysql.Open("user:pass@tcp(localhost:3306)/mydb?charset=utf8mb4&parseTime=True&loc=Local"))
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect to database")
	}

	// Get the id from the vars passed in.
	vars := mux.Vars(r)
	staffToFoundID := vars["id"]
	id, err := strconv.ParseInt(staffToFoundID, 10, 64)
	if err != nil {
		fmt.Println(err.Error())
		panic("Couldn't find in database")
	}

	var staffMember Person
	db.First(&staffMember, id)

	js, err := json.Marshal(staffMember)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

	fmt.Println("Staff member found")

}

func AddStaff(w http.ResponseWriter, r *http.Request) {

	db, err := gorm.Open(mysql.Open("user:pass@tcp(localhost:3306)/mydb?charset=utf8mb4&parseTime=True&loc=Local"))
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect to database")
	}

	decoder := json.NewDecoder(r.Body)
	var staffMember Person
	err = decoder.Decode(&staffMember)
	if err != nil {
		panic(err)
	}

	db.Create(&Person{
		FirstName: staffMember.FirstName,
		SurName:   staffMember.SurName,
		Age:       staffMember.Age,
		Salary:    staffMember.Salary})

	fmt.Println("Post made")
	fmt.Fprintf(w, "Staff Member Added: %+v", staffMember)
}

func DeleteStaff(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(mysql.Open("user:pass@tcp(localhost:3306)/mydb?charset=utf8mb4&parseTime=True&loc=Local"))
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect to database")
	}

	vars := mux.Vars(r)
	staffToDeleteId := vars["id"]
	id, err := strconv.ParseInt(staffToDeleteId, 10, 64)
	if err != nil {
		fmt.Println(err.Error())
		panic("Couldn't delete from database")
	}

	db.Delete(&Person{}, "id LIKE ?", id)
	fmt.Println("Staff member deleted")

}

func HomeHandler(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("Hit the home route"))
	fmt.Fprintf(w, "Hi")

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
	r.HandleFunc("/staff/{id}", GetStaffByName).Methods("GET")
	r.HandleFunc("/staff", AddStaff).Methods("POST")
	r.HandleFunc("/staff/{id}", DeleteStaff).Methods("DELETE")

	fmt.Println("Server started")

	log.Fatal(http.ListenAndServe(":8000", r))

}
