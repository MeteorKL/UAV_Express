package main

import (
	"testing"
	"log"
)

func Test_getUserById(t *testing.T) {
	userTable = []DB_User{{User_id: 1, User_name:"S1Meteor"},
		{User_id: 2, User_name: "S2Meteor"},
		{User_id: 3, User_name: "S3Meteor"}}
	ReBuildIndex()
	log.Println(getUserById(2))
}

func Test_getRecentPayments(t *testing.T) {
	userTable = []DB_User{{User_id: 1, User_name:"S1Meteor"},
						  {User_id: 2, User_name: "S2Meteor"},
						  {User_id: 3, User_name: "S3Meteor"}}
	uavTable = []DB_UAV{{UAV_id: 1, UAV_name:"No0"},
						{UAV_id:2, UAV_name: "No1"}}
	paymentTable = []DB_Payment{{Payment_id:1, Payment_user_id:2, Payment_time:15646463},
		{Payment_id:3, Payment_user_id:3, Payment_time:183456486}}
	ReBuildIndex()
	user := getUserById(2)
	log.Println(user.getRecentPayments(2))
}

func Test_getUAVList(t *testing.T) {
	uavTable = []DB_UAV{{UAV_id: 1, UAV_name:"No0"},
		{UAV_id:2, UAV_name: "No1"}}
	ReBuildIndex()
	log.Println(getUAVList(1,4))
}