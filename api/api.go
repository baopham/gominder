package api

import (
	"encoding/json"
	"github.com/baopham/gominder/models"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io"
	"log"
	"net/http"
)

func Remind(w http.ResponseWriter, r *http.Request) {
	db, ok := r.Context().Value("database").(*mgo.Session)

	if !ok {
		badRequest(w)
		return
	}

	var ids []string
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&ids)

	if err != nil {
		badRequest(w)
		return
	}

	err = models.Remind(ids, db)

	if err != nil {
		badRequest(w)
		return
	}

	respond("ok", w)
}

func GetReminders(w http.ResponseWriter, r *http.Request) {
	db, ok := r.Context().Value("database").(*mgo.Session)

	if !ok {
		badRequest(w)
		return
	}

	reminders, err := models.GetReminders(db)

	if err != nil {
		badRequest(w)
		return
	}

	respond(reminders, w)
}

func GetReminder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	db, ok := r.Context().Value("database").(*mgo.Session)

	if !ok {
		badRequest(w)
		return
	}

	reminder, err := models.FindReminder(id, db)

	if err != nil {
		badRequest(w)
		return
	}

	respond(reminder, w)
}

func AddReminder(w http.ResponseWriter, r *http.Request) {
	err := validateReminderPayload(r)

	if err != nil {
		badRequest(w)
		return
	}

	db, ok := r.Context().Value("database").(*mgo.Session)

	if !ok {
		badRequest(w)
		return
	}

	reminder, err := parseToReminder(r.Body)
	reminder.ID = ""
	reminder.RemindCount = 0

	if err != nil {
		log.Println("AddReminder:", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	err = reminder.Save(db)

	if err != nil {
		badRequest(w)
		return
	}

	respond("ok", w)
}

func validateReminderPayload(r *http.Request) error {
	// TODO
	return nil
}

func badRequest(w http.ResponseWriter) {
	http.Error(w, "Invalid request", http.StatusBadRequest)
}

func parseToReminder(r io.Reader) (models.Reminder, error) {
	var reminder models.Reminder
	decoder := json.NewDecoder(r)
	err := decoder.Decode(&reminder)
	// TODO: embed user id by the current user
	reminder.UserID = bson.ObjectId("19729de860ea")
	return reminder, err
}

func setContentType(w *http.ResponseWriter) {
	(*w).Header().Set("Content-Type", "application/json")
}

func respond(payload interface{}, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(payload)

	if err != nil {
		log.Println("failed to encode to JSON: ", err)
	}
}
