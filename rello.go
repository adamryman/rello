package rello

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc"

	pb "github.com/adamryman/rello/rello-service"

	ambition "github.com/adamryman/ambition-model/ambition-service"
	ambitiongrpc "github.com/adamryman/ambition-model/ambition-service/generated/client/grpc"
)

var aService ambition.AmbitionServiceServer

func init() {
	//dbLocation := os.Getenv("SQLITE3")
	//InitDatabase(dbLocation)

	ambitionGRPCAddr := os.Getenv("ambition_grpc")
	fmt.Println(ambitionGRPCAddr)
	fmt.Printf("% x\n", ambitionGRPCAddr)
	conn, err := grpc.Dial(ambitionGRPCAddr, grpc.WithInsecure(), grpc.WithTimeout(time.Second*10))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while dialing grpc connection: %v", err)
		os.Exit(1)
	}
	aService, err = ambitiongrpc.New(conn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while creating ambition service: %v", err)
		os.Exit(1)
	}
}

const trelloTime = "2006-01-02T15:04:05.999Z07:00"

func HandleWebhook(w http.ResponseWriter, r *http.Request) {
	var update pb.ChecklistUpdate
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&update)
	if err != nil {
		fmt.Println(err)
		return
	}
	HandleUpdate(&update)
}

func HandleUpdate(u *pb.ChecklistUpdate) {
	cItem := u.Action.Data.CheckItem
	switch u.Action.Type {
	case "createCheckItem":
		// TODO: CreateAction
		// TODO: Map CheckItem.Id to Action.Id
		// TODO: Map u.Action.MemberCreator.Id to Action.UserId
		// UserId hardcoded to 1
		_ = u.Action.MemberCreator.Id

		user := readUser(u.Action.MemberCreator.Id)
		action, err := user.CreateAction(cItem)
		if err != nil {
			fmt.Println(errors.Wrap(err, "cannot create action"))
			return
		}
		fmt.Printf("%s action created\n", action.ActionName)
	case "updateCheckItemStateOnCard":
		if cItem.State == "incomplete" {
			fmt.Printf("%q is being unchecked\n", cItem.Name)
			return
		}
		dateString := &u.Action.Date
		date, err := time.Parse(trelloTime, *dateString)
		_ = date
		if err != nil {
			fmt.Println(errors.Wrap(err, "time parsing failed"))
			return
		}
		var req ambition.Action
		req.TrelloId = cItem.Id
		actionResponse, err := aService.ReadAction(context.Background(), &req)
		if err != nil || actionResponse.Action == nil {
			fmt.Println(errors.Wrap(err, "cannot read action"))
			fmt.Println(actionResponse.Error)
			return
		}
		var o ambition.Occurrence
		o.ActionId = actionResponse.Action.ActionId
		o.Datetime = date.String()
		resp, err := aService.CreateOccurrence(context.Background(), &o)
		if err != nil || resp.Occurrence == nil {
			fmt.Println(errors.Wrap(err, "cannot create occurrence"))
			fmt.Println(resp.Error)
			return
		}

		// TODO: Create Occurrence with the ActionId from CheckItem.Id
		// TODO: Add date to occurrence
		// Maps to action id
		_ = cItem.Id
		fmt.Printf("%s occurrence created of action %s\n", resp.Occurrence.Datetime, actionResponse.Action.ActionName)

	default:
		fmt.Println(u.Action.Type)
	}
}

type user int

func readUser(trelloMemberId string) user {
	return 1
}

func (u user) CreateAction(c *pb.CheckItem) (*ambition.Action, error) {
	var req ambition.Action
	req.ActionName = c.Name
	req.UserId = 1
	req.TrelloId = c.Id
	actionResponse, err := aService.CreateAction(context.Background(), &req)
	if err != nil {
		return nil, err
	}

	return actionResponse.Action, err
}

func (u user) CreateOccurrence() {

}
