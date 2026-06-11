// All the code in this file belongs to the package named golangCrud.
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

const (
	dbDriver = "mysql"
	dbUser   = "root"
	dbPass   = "1234"
	dbName   = "gocrud"
)

func main() {
	//create new router
	r := mux.NewRouter()

	// define endpoints
	r.HandleFunc("/user", createUserHandler).Methods("POST")
	r.HandleFunc("/user/{id}", getUserHandler).Methods("GET")
	r.HandleFunc("/user/{id}", updateUserHandler).Methods("PUT")
	r.HandleFunc("/user/{id}", deleteUserHandler).Methods("DELETE")

	// Start the HTTP server on port 8090
	log.Println("server listening on port 8090")
	log.Fatal(http.ListenAndServe(":8090", r))
}
func createUserHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// Parse JSON data from the request body
	var user User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Failed to parse user data", http.StatusBadRequest)
		return
	}

	err = createUser(db, user.Name, user.Email)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "User created successfully")
}

// Create user actual crud function
func createUser(db *sql.DB, name, email string) error {
	query := "INSERT INTO users(name,email) VALUES(?,?)"
	_, err := db.Exec(query, name, email)
	if err != nil {
		fmt.Println("DATABASE ERROR:", err)
		return err
	}
	fmt.Println("User created successfully")
	return nil
}

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// Get the 'id' parameter from the URL
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Convert id to int
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid 'id' parameter", http.StatusBadRequest)
		return
	}

	user, err := getUser(db, userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Convert the user object to JSON and send it in the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// Get user actual crud function
func getUser(db *sql.DB, id int) (*User, error) {
	query := "SELECT * FROM users WHERE id=?"
	row := db.QueryRow(query, id)

	user := &User{} // this is for creating a new user object to hold the data retrieved from the database.
	err := row.Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// Get the 'id' parameter from the URL
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Convert 'id' to an integer
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid 'id' parameter", http.StatusBadRequest)
		return
	}

	var user User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Failed to parse user data", http.StatusBadRequest)
		return
	}

	err = updateUser(db, userID, user.Name, user.Email)
	if err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "User updated successfully")
}

// Update user actual crud function
func updateUser(db *sql.DB, id int, name, email string) error {
	query := "UPDATE users SET name = ?, email = ? WHERE id = ?"
	_, err := db.Exec(query, name, email, id)
	if err != nil {
		return err
	}
	return nil
}

// delete data
func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	// Get the 'id' parameter from the URL
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Convert 'id' to an integer
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid 'id' parameter", http.StatusBadRequest)
		return
	}

	err = deleteUser(db, userID)
	if err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "User deleted successfully")
}

// this function is for deleting the user from the database
func deleteUser(db *sql.DB, id int) error {
	query := "DELETE FROM users WHERE id = ?"
	_, err := db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
