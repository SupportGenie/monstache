package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rwynn/monstache/monstachemap"
	"github.com/rwynn/monstache/supportgenie"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Perform(input *monstachemap.MapperPluginInput) (output *monstachemap.MapperPluginOutput, err error) {
	client := input.MongoClient
	ctx := context.TODO()

	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})

	log.Info("inside prepare function")
	log.Info("getting ticket Id")
	var ticket map[string]interface{}
	ticket, ok := input.Document["ticketId"].(map[string]interface{})
	if !ok {
		log.Error("ticket not of type map[string]interface{}")
		return nil, fmt.Errorf("ticket not of type map[string]interface{}")
	}
	ticketId, ok := ticket["id"].(string)
	if !ok {
		log.Error("ticketId not of type string")
		return nil, fmt.Errorf("ticketId not of type string")
	}
	objId, err := primitive.ObjectIDFromHex(ticketId)
	if err != nil {
		log.Error(err)
		return nil, fmt.Errorf("error converting string to ObjectID")
	}
	logger := log.WithField("ticketId", objId.Hex())
	collection := client.Database("supportgenie_test").Collection("ticket")
	err = collection.FindOne(ctx, bson.D{{"ticketId", objId}}).Decode(&ticket)
	if err != nil || len(ticket) == 0 {
		logger.Error("ticket not found", err)
		return nil, fmt.Errorf("ticket not found")
	}
	logger.Info("got ticket data")

	logger.Info("converting companyId to objectId")
	company_id, ok := ticket["companyId"].(primitive.ObjectID)
	if !ok {
		logger.Error("could not convert companyId to objectId")
		return nil, fmt.Errorf("could not convert companyId to objectId")
	}
	logger.Info("converted companyId to objectId")

	logger.Info("converting userId to objectId")
	userId, ok := ticket["userId"].(primitive.ObjectID)
	if !ok {
		logger.Error("could not convert userId to objectId")
		return nil, fmt.Errorf("could not convert userId to objectId")
	}

	logger = logger.WithFields(log.Fields{
		"ticketId":  objId.Hex(),
		"companyId": company_id.Hex(),
		"userId":    userId.Hex(),
	})

	logger.Info("getting company data")
	var company_data supportgenie.Company
	collection = client.Database("supportgenie_test").Collection("company")

	err = collection.FindOne(ctx, bson.D{{"companyId", company_id}}).Decode(&company_data)
	if err != nil {
		logger.Error("get company error", err)
		return nil, fmt.Errorf("get company error")
	}
	logger.Info("got company data")

	logger.Info("getting user data")
	var user_data supportgenie.User
	collection = client.Database("supportgenie_test").Collection("user")
	err = collection.FindOne(ctx, bson.D{{"userId", userId}}).Decode(&user_data)
	if err != nil {
		logger.Error("get user error", err)
		return nil, fmt.Errorf("get user error")
	}
	logger.Info("got user data")

	//clean ticket data
	logger.Info("cleaning ticket data")
	fields_to_remove := []string{"CompanyId", "UserId", "extraFields", "primaryAgentId"}
	for _, field := range fields_to_remove {
		logger.Info("removing field from ticket ", field)
		delete(ticket, field)
	}
	logger.Info("cleaned ticket data")

	logger.Info("creating final ticket")

	ticket_data := supportgenie.Ticket{
		Company: company_data,
		User:    user_data,
		Ticket:  ticket,
	}
	final_data := make(map[string]interface{})
	final_marshal, err := json.Marshal(ticket_data)
	if err != nil {
		logger.Error("marshal error", err)
		return nil, fmt.Errorf("marshal error")
	}
	json.Unmarshal(final_marshal, &final_data)
	output = &monstachemap.MapperPluginOutput{Document: final_data}
	return
}
