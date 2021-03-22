package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/tnyie/tny-api/models"
)

func createUsers() {
	for i := 0; i < 10; i++ {
		created_at := rand.Int63()
		userAuth := &models.UserAuth{
			UID:      "",
			Username: fmt.Sprint("username" + strconv.Itoa(i)),
			Email:    fmt.Sprint("test", strconv.Itoa(i), "@test.com"),
			Enabled:  true,
		}
		userAuth.Create("password")
		userAuth.Get()
		user := &models.User{
			UID:       userAuth.UID,
			Username:  userAuth.Username,
			Email:     userAuth.Email,
			CreatedAt: created_at,
			UpdatedAt: created_at + 10000,
		}
		err := user.Create()
		createLinks(userAuth.UID)
		if err != nil {
			fmt.Println("ERROR CREATING USER\n", err)
		}
	}
}

func createLinks(user_id string) {
	for i := 0; i < 3; i++ {
		link := &models.Link{
			OwnerID: user_id,
			Slug:    fmt.Sprint(user_id, "-", strconv.Itoa(i)),
			URL:     "https://netsoc.co/rk",
		}
		link.Create()
		link.Get()
		createVisits(link.ID)
	}
}

func createVisits(link_id string) {
	curr_time := time.Now()
	for i := 0; i < 10; i++ {
		curr_time = curr_time.Add(time.Hour * -24)
		for j := 0; i < 10; j++ {
			visit := &models.Visit{
				LinkID:    link_id,
				CreatedAt: curr_time,
			}
			visit.Create()
		}
	}
}
