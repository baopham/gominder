package models

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

type Reminder struct {
	ID          bson.ObjectId `json:"ID" bson:"_id,omitempty"`
	UserID      bson.ObjectId `json:"UserID" bson:"userid,omitempty"`
	Name        string
	Description string
	URL         string
	RemindAt    *time.Time
	RemindCount int
	Tags        []string
}

func FindReminder(id string, s *mgo.Session) (*Reminder, error) {
	session, c := reminderDB(s)
	defer session.Close()

	var reminder Reminder

	err := c.FindId(bson.ObjectIdHex(id)).One(&reminder)

	if err != nil {
		log.Println("failed to find: ", err)
		return nil, err
	}

	return &reminder, nil
}

func Remind(ids []string, s *mgo.Session) error {
	session, c := reminderDB(s)
	defer session.Close()

	docIDs := make([]bson.ObjectId, len(ids))
	for i, id := range ids {
		docIDs[i] = bson.ObjectIdHex(id)
	}

	var reminders []Reminder
	err := c.Find(bson.M{"_id": bson.M{"$in": docIDs}}).All(&reminders)

	if err != nil {
		log.Println("failed to find: ", err)
		return nil
	}

	for _, reminder := range reminders {
		reminder.RemindCount = reminder.RemindCount + 1
		err := reminder.Save(session)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetReminders(s *mgo.Session) ([]Reminder, error) {
	session, c := reminderDB(s)
	defer session.Close()

	var reminders []Reminder
	err := c.Find(nil).Limit(50).Sort("remindat", "-remindcount").All(&reminders)

	if err != nil {
		log.Println("failed to get reminders: ", err)
		return nil, err
	}

	return reminders, nil
}

func (r *Reminder) Remove(s *mgo.Session) error {
	session, c := reminderDB(s)
	defer session.Close()

	err := c.Remove(bson.M{"_id": r.ID})

	if err != nil {
		log.Println("failed to remove: ", err)
		return err
	}

	return nil
}

func (r *Reminder) Save(s *mgo.Session) error {
	session, c := reminderDB(s)
	defer session.Close()

	var err error

	if r.ID == "" {
		err = c.Insert(r)
	} else {
		_, err = c.Upsert(bson.M{"_id": r.ID}, r)
	}

	if err != nil {
		log.Println("failed to save: ", err)
		return err
	}

	return nil
}

func ReminderIndices() []mgo.Index {
	return []mgo.Index{
		{
			Key:    []string{"remindat"},
			Sparse: true,
		},
		{
			Key:    []string{"userid"},
			Sparse: true,
		},
		{
			Key:    []string{"tags"},
			Sparse: true,
		},
	}
}

func reminderDB(s *mgo.Session) (*mgo.Session, *mgo.Collection) {
	session := s.Copy()
	c := session.DB("gominder").C("reminders")
	return session, c
}
