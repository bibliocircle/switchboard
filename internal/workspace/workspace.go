package workspace

import (
	"context"
	"switchboard/internal/common"
	"switchboard/internal/db"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateWorkspace(userId string, ws *CreateWorkspaceRequestBody) (*Workspace, *common.DetailedError) {
	wsId := common.GetShortId()
	currentTime := time.Now()
	var expiresAt *time.Time

	if ws.ExpiresAt != "" {
		var errParse error
		exp, errParse := time.Parse(time.RFC3339, ws.ExpiresAt)
		expiresAt = &exp
		if errParse != nil {
			return nil, common.NewDetailedError(common.ErrorInvalidInput, "could not parse expiresAt value")
		}
	}

	newWs := &Workspace{
		ID:        wsId,
		Name:      ws.Name,
		ExpiresAt: expiresAt,
		CreatedBy: userId,
		CreatedAt: &currentTime,
		UpdatedAt: &currentTime,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	wsCollection := db.Database.Collection(db.WORKSPACES_COLLECTION)
	_, insertErr := wsCollection.InsertOne(ctx, newWs)
	if insertErr != nil {
		return nil, db.GetDbError(insertErr)
	}
	var createdWs Workspace
	findErr := wsCollection.FindOne(ctx, bson.D{{
		Key:   "id",
		Value: wsId,
	}}).Decode(&createdWs)
	if findErr != nil {
		if findErr == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, db.GetDbError(findErr)
	}
	return &createdWs, nil
}

func FindWorkspaces(filter *bson.D) (*[]Workspace, *common.DetailedError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	wsCol := db.Database.Collection(db.WORKSPACES_COLLECTION)
	findOpts := &options.FindOptions{
		Sort: &map[string]int64{
			"createdAt": -1,
		},
	}

	cursor, errFind := wsCol.Find(ctx, filter, findOpts)
	if errFind != nil {
		return &[]Workspace{}, db.GetDbError(errFind)
	}
	result := make([]Workspace, 0)
	err := cursor.All(ctx, &result)
	if err != nil {
		return nil, db.GetDbError(err)
	}
	return &result, nil
}

func GetWorkspaces(userID string) (*[]Workspace, *common.DetailedError) {
	return FindWorkspaces(&bson.D{
		{Key: "createdBy", Value: userID},
	})
}

func GetUserWorkspaceByID(userID, workspaceID string) (*Workspace, *common.DetailedError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var ws Workspace
	wsCol := db.Database.Collection(db.WORKSPACES_COLLECTION)
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
		if findErr == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, db.GetDbError(findErr)
	}
	return &ws, nil
}

func IsWorkspaceOwner(userId, workspaceId string) (bool, *common.DetailedError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	wsCol := db.Database.Collection(db.WORKSPACES_COLLECTION)
	count, err := wsCol.CountDocuments(ctx, bson.D{
		{Key: "id", Value: workspaceId},
		{Key: "createdBy", Value: userId},
	})
	if err != nil {
		return false, db.GetDbError(err)
	}
	return count > 0, nil
}

func DeleteWorkspace(userID, wsId string) (bool, *common.DetailedError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	wsCol := db.Database.Collection(db.WORKSPACES_COLLECTION)
	result, errDel := wsCol.DeleteOne(ctx, bson.D{
		{Key: "id", Value: wsId},
		{Key: "createdBy", Value: userID},
	})
	if errDel != nil {
		return false, db.GetDbError(errDel)
	}
	return result.DeletedCount > 0, nil
}
