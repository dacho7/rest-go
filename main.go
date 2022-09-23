package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type task struct {
	Id      int    `json:"Id"`
	Name    string `json:"Name"`
	Content string `json:"Content"`
}

type allTasks []task

var tasks = allTasks{
	{
		Id:      1,
		Name:    "lavar platos",
		Content: "No olvidar las ollas",
	},
}

func indexRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to my Api")
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func createTask(w http.ResponseWriter, r *http.Request) {
	var newTask task
	reqBody, er := ioutil.ReadAll(r.Body)
	if er != nil {
		fmt.Fprint(w, "Insert a valida task")
	}

	json.Unmarshal(reqBody, &newTask)

	newTask.Id = len(tasks) + 1
	tasks = append(tasks, newTask)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tasks)
}

func getTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	idTask, err := strconv.Atoi(vars["id"])

	if err != nil {
		fmt.Fprint(w, "Invalid Id")
		return
	}

	for index, task := range tasks {
		if task.Id == idTask {
			w.Header().Set("Content-Type", "application/json")
			fmt.Print(index)
			json.NewEncoder(w).Encode(task)
		}
	}

}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idTask, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Fprint(w, "Invalid id")
		return
	}

	for index, task := range tasks {
		if task.Id == idTask {
			tasks = append(tasks[:index], tasks[index+1:]...)
			json.NewEncoder(w).Encode(task)
		}
	}
}

func main() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", indexRoute)
	router.HandleFunc("/tasks", getTasks)
	router.HandleFunc("/createtask", createTask).Methods("POST")
	router.HandleFunc("/gettask/{id}", getTask).Methods("GET")
	router.HandleFunc("/deletetask/{id}", deleteTask).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":3010", router))
}
