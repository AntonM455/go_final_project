package api

import (
	"encoding/json"
	"go_final_project/pkg/db"
	"net/http"
	"time"
)

// taskHandler - router using HTTP methods
func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// addTaskHandler — handling POST /api/task
func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	// deserialize the request body into a Task structure
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		// return status 400
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": "incorrect data in JSON"})
		return
	}

	// check for mandatory title
	if task.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		// return status 400
		writeJSON(w, map[string]string{"error": "task title not specified"})
		return
	}

	// checking and correcting the date using checkDate
	if err := checkDate(&task); err != nil {
		// return status 400
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	// add to DB
	id, err := db.AddTask(&task)
	if err != nil {
		// return status 500
		w.WriteHeader(http.StatusInternalServerError)
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, map[string]any{"id": id})
}

// getTaskHandler — handling GET /api/task
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		// return status 400
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": "ID not specified"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		// return status 404
		w.WriteHeader(http.StatusNotFound)
		writeJSON(w, map[string]string{"error": "task not found"})
		return
	}

	writeJSON(w, task)
}

// updateTaskHandler — handling PUT api/task
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		// return status 400
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": "invalid JSON"})
		return
	}

	// check for mandatory title
	if task.Title == "" {
		// return status 400
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": "task title not specified"})
		return
	}

	// checking and correcting the date using checkDate
	if err := checkDate(&task); err != nil {
		// return status 400
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	if err := db.UpdateTask(&task); err != nil {
		// return status 500
		w.WriteHeader(http.StatusInternalServerError)
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, map[string]string{})
}

// checkDate checks and adjusts the task date in the Date field.
func checkDate(task *db.Task) error {
	now := time.Now()

	// if the date is not specified, we substitute today's date
	if task.Date == "" {
		task.Date = now.Format("20060102")
	}
	// parse date from string
	t, err := time.Parse("20060102", task.Date)
	if err != nil {
		return err
	}

	var next string

	// if a repeat rule is specified, we check and calculate next
	if task.Repeat != "" {
		next, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}

	// date comparison and adjustment
	if afterNow(now, t) {
		if task.Repeat == "" {
			task.Date = now.Format("20060102")
		} else {
			task.Date = next
		}
	}

	return nil
}

// writeJSON serializes the response and sets the headers
func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(data)
}

// doneTaskHandler — handling POST /api/task
func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "ID not specified"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, map[string]string{"error": "task not found"})
		return
	}

	// if there is no repetition rule, then delete the task
	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			writeJSON(w, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, map[string]string{})
		return
	}

	// if the task is periodic, we calculate a new date
	now := time.Now()
	next, err := NextDate(now, task.Date, task.Repeat)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	if err := db.UpdateDate(next, id); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, map[string]string{})
}

// deleteTaskHandler — handling DELETE /api/task
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "ID not specified"})
		return
	}

	if err := db.DeleteTask(id); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, map[string]string{})
}
