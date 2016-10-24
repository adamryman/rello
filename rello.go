package rello

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
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
			fmt.Println(errors.Wrap(err, "time parsing failed"))
			return
		}
		ci, err = readCheckItemByTrelloId(ci.Id)
		if err != nil {
			fmt.Println(errors.Wrapf(err, "error reading %s from database", ci.Id))
			return
		}
		c.CheckItemId = ci.RelloId

		id, err := createCheck(c)
		if err != nil {
			fmt.Println(errors.Wrapf(err, "cannot create check %v", c))
			return
		}
		c.Id = id
		fmt.Printf("Check created! %v\n", c)
	case "createCheckItem":
		ci.UserId = u.Action.MemberCreator.Id
		id, err := createCheckItem(*ci)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("%s created with Id %d\n", ci.Name, id)
	default:
		fmt.Println(u.Action.Type)
	}

}
