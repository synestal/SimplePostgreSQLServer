package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Person struct {
	Id         string `json:"id"`
	Name       string `json:"nme"`
	Birthday   string `json:"birthday"`
	Occupation string `json:"occupation"`
}

type AuthStruct struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Dbname   string `json:"dbname"`
	Sslmode  string `json:"sslmode"`
}

func getPersonHandler(w http.ResponseWriter, r *http.Request) {
	personList, err := store.GetPerson()
	if err != nil {
		w.Write([]byte(""))
		fmt.Println(fmt.Errorf("Error: %v", err))
		return
	}
	// Convert the `personList` variable to JSON
	personListBytes, err := json.Marshal(personList)
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(personListBytes)
	if err != nil {
		return
	}
}

func createPersonHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	person := Person{}
	person.Name = r.Form.Get("nme")
	person.Birthday = r.Form.Get("birthday")
	person.Occupation = r.Form.Get("occupation")

	err = store.CreatePerson(&person)
	if err != nil {
		fmt.Println(err)
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func deletePersonHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	id, err_ := strconv.Atoi(r.Form.Get("id"))
	if err_ != nil {
		fmt.Println(err_)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	err = store.DeletePerson(id)
	if err != nil {
		fmt.Println(err)
	}

	http.Redirect(w, r, "/", http.StatusFound)

}

func Authorise(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	/*
		Connect to PostgreSQL
		user - a member from whom do you make a work
		password - password for user to log in PostgreSQL Database
		host - default localhost
		port - port for PostgreSQL Database, setted in settings, default - 5432
		dbname - name of Dtabase to which you are connecting
		sslmode - level of protection, every iteration within disable requires additional code and Database setup. Default - disable
	*/

	client := &AuthStruct{}
	client.User = r.Form.Get("user")
	client.Password = r.Form.Get("password")
	client.Host = r.Form.Get("host")
	client.Port = r.Form.Get("port")
	client.Dbname = r.Form.Get("dbname")
	client.Sslmode = r.Form.Get("sslmode")
	portAtoi, _ := strconv.Atoi(client.Port)
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", client.Host, portAtoi, client.User, client.Password, client.Dbname, client.Sslmode)

	db, err_ := sql.Open("postgres", psqlconn)
	if err_ != nil {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		panic(err_)
		return
	}
	err = db.Ping()
	if err != nil {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		panic(err)
		return
	}
	store = &dbStore{db: db}

	http.Redirect(w, r, "/", http.StatusFound)
	return
}
