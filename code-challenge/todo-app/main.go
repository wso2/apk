package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Todo struct {
	ID   int    `json:"id"`
	Task string `json:"task"`
	Done bool   `json:"done"`
}

var todos []Todo
var nextID = 1

// Get all to-dos
func getTodos(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(todos)
}

// Create a new to-do
func createTodo(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	_ = json.NewDecoder(r.Body).Decode(&todo)
	todo.ID = nextID
	nextID++
	todos = append(todos, todo)
	json.NewEncoder(w).Encode(todo)
}

// Get a single to-do by ID
func getTodoByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])
	for _, todo := range todos {
		if todo.ID == id {
			json.NewEncoder(w).Encode(todo)
			return
		}
	}
	http.Error(w, "To-do not found", http.StatusNotFound)
}

// Update a to-do by ID
func updateTodo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])
	for index, todo := range todos {
		if todo.ID == id {
			todos = append(todos[:index], todos[index+1:]...)
			var updatedTodo Todo
			_ = json.NewDecoder(r.Body).Decode(&updatedTodo)
			updatedTodo.ID = id
			todos = append(todos, updatedTodo)
			json.NewEncoder(w).Encode(updatedTodo)
			return
		}
	}
	http.Error(w, "To-do not found", http.StatusNotFound)
}

// Delete a to-do by ID
func deleteTodo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])
	for index, todo := range todos {
		if todo.ID == id {
			todos = append(todos[:index], todos[index+1:]...)
			json.NewEncoder(w).Encode(todos)
			return
		}
	}
	http.Error(w, "To-do not found", http.StatusNotFound)
}

// Register a user
func registerUser(w http.ResponseWriter, r *http.Request) {
	regHeader := r.Header.Get("user-reg")
	if regHeader == "" {
		http.Error(w, "Missing user-reg header", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User registered with header: " + regHeader))
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/todos", getTodos).Methods("GET")
	router.HandleFunc("/todos", createTodo).Methods("POST")
	router.HandleFunc("/todos/{id}", getTodoByID).Methods("GET")
	router.HandleFunc("/todos/{id}", updateTodo).Methods("PUT")
	router.HandleFunc("/todos/{id}", deleteTodo).Methods("DELETE")
	router.HandleFunc("/register", registerUser).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))
}
