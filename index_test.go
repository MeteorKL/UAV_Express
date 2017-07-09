package main

import (
	"testing"
	"log"
)

func TestReBuildIndex(t *testing.T) {
	userTable = []DB_User{
		{User_id: 1, User_name: "S1Meteor",},
		{User_id: 2, User_name: "S2Meteor"},
	}
	ReBuildIndex()
	iter, _ := user_id_index.tree.GetIterator(2)
	log.Println(iter.Value().(UserRecord).DB_User)
	iter.Prev()
	log.Println(iter.Value().(UserRecord).DB_User)
}
