package rello

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/davecgh/go-spew/spew"
)

const trelloTime = "2006-01-02T15:04:05.999Z07:00"

type Check struct {
	Id          int64
	CheckItemId int64
	Date        time.Time
}

func HandleWebhook(w http.ResponseWriter, r *http.Request) {
	var update ChecklistUpdate
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&update)
	if err != nil {
		fmt.Println(err)
		return
	}
	handleUpdate(update)
}

func handleUpdate(u ChecklistUpdate) {
	var err error
	ci := u.Action.Data.CheckItem
	switch u.Action.Type {
	case "updateCheckItemStateOnCard":
		if ci.State == "incomplete" {
			fmt.Printf("%q is being unchecked\n", ci.Name)
			return
		}
		var c Check
		c.Date, err = time.Parse(trelloTime, u.Action.Date)
		if err != nil {
			fmt.Printf("time parsing failed: %v\n", err)
			return
		}
		spew.Dump(c)
	case "createCheckItem":
		ci.UserId = u.Action.MemberCreator.Id
		id, err := CreateCheckItem(*ci)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("%s created with Id %d\n", ci.Name, id)
	default:
		fmt.Println(u.Action.Type)
	}

}
