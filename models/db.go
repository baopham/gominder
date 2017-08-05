package models

import (
	"gopkg.in/mgo.v2"
)

func NewSession() (*mgo.Session, error) {
	session, err := mgo.Dial("localhost")

	if err != nil {
		return nil, err
	}

	session.SetMode(mgo.Monotonic, true)

	err = ensureIndex(session)

	if err != nil {
		return nil, err
	}

	return session, nil
}

func ensureIndex(s *mgo.Session) error {
	session := s.Copy()
	defer session.Close()

	c := session.DB("gominder").C("reminders")

	for _, index := range ReminderIndices() {
		err := c.EnsureIndex(index)

		if err != nil {
			return err
		}
	}

	return nil
}
