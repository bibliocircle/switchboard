package db

import (
	"context"
	"switchboard/internal/common"
	"switchboard/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateWorkspace(userId string, ws *models.CreateWorkspaceRequestBody) (*models.Workspace, *common.DetailedError) {
	wsId := common.GetShortId()
	currentTime := time.Now()
	newWs := &models.Workspace{
		ID:        wsId,
		Name:      ws.Name,
		ExpiresAt: ws.ExpiresAt,
		CreatedBy: userId,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	wsCollection := Database.Collection(WORKSPACES_COLLECTION)
	_, insertErr := wsCollection.InsertOne(ctx, newWs)
	if insertErr != nil {
		return nil, GetDbError(insertErr)
	}
	var createdWs models.Workspace
	findErr := wsCollection.FindOne(ctx, bson.D{{
		Key:   "id",
		Value: wsId,
	}}).Decode(&createdWs)
	if findErr != nil {
		return nil, common.WrapAsDetailedError(findErr)
	}
	return &createdWs, nil
}

func FindWorkspaces(filter *bson.D) (*[]models.Workspace, *common.DetailedError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	wsCol := Database.Collection(WORKSPACES_COLLECTION)
	findOpts := &options.FindOptions{
		Sort: &map[string]int64{
			"createdAt": -1,
		},
	}

	cursor, errFind := wsCol.Find(ctx, filter, findOpts)
	if errFind != nil {
		return &[]models.Workspace{}, common.WrapAsDetailedError(errFind)
	}
	result := make([]models.Workspace, 0)
	err := cursor.All(ctx, &result)
	if err != nil {
		return nil, common.WrapAsDetailedError(err)
	}
	return &result, nil
}

func GetWorkspaces() (*[]models.Workspace, *common.DetailedError) {
	return FindWorkspaces(&bson.D{})
}

func GetUserWorkspaces(userID string) (*[]models.Workspace, *common.DetailedError) {
	return FindWorkspaces(&bson.D{
		{Key: "createdBy", Value: userID},
	})
}

func GetUserWorkspaceByID(userID, workspaceID string) (*models.Workspace, *common.DetailedError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var ws models.Workspace
	wsCol := Database.Collection(WORKSPACES_COLLECTION)
	findErr := wsCol.FindOne(ctx, bson.D{
		{
			Key:   "id",
			Value: workspaceID,
		},
		{
			Key:   "createdBy",
			Value: userID,
		},
	}).Decode(&ws)
	if findErr != nil {
		return nil, common.WrapAsDetailedError(findErr)
	}
	return &ws, nil
}

func IsWorkspaceOwner(userId, workspaceId string) (bool, *common.DetailedError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	wsCol := Database.Collection(WORKSPACES_COLLECTION)
	count, err := wsCol.CountDocuments(ctx, bson.D{
		{Key: "id", Value: workspaceId},
		{Key: "createdBy", Value: userId},
	})
	if err != nil {
		return false, common.WrapAsDetailedError(err)
	}
	return count > 0, nil
}

func DeleteWorkspace(userID, wsId string) (bool, *common.DetailedError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	wsCol := Database.Collection(WORKSPACES_COLLECTION)
	result, errDel := wsCol.DeleteOne(ctx, bson.D{
		{Key: "id", Value: wsId},
		{Key: "createdBy", Value: userID},
	})
	if errDel != nil {
		return false, common.WrapAsDetailedError(errDel)
	}
	return result.DeletedCount > 0, nil
}
