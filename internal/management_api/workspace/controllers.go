package workspace

import (
	"context"
	"switchboard/internal/common/db"
	"switchboard/internal/common/err_utils"
	"switchboard/internal/common/randomdata"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Workspace struct {
	ID        string    `json:"id" bson:"id,omitempty"`
	Name      string    `json:"name" bson:"name,omitempty"`
	ExpiresAt string    `json:"expiresAt,omitempty" bson:"expiresAt,omitempty"`
	CreatedBy string    `json:"createdBy" bson:"createdBy,omitempty"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt,omitempty"`
}

func CreateWorkspace(userId string, ws *CreateWorkspaceRequestBody) (*Workspace, *err_utils.DetailedError) {
	wsId := randomdata.GetShortId()
	currentTime := time.Now()
	newWs := &Workspace{
		ID:        wsId,
		Name:      ws.Name,
		ExpiresAt: ws.ExpiresAt,
		CreatedBy: userId,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
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
		return nil, err_utils.WrapAsDetailedError(findErr)
	}
	return &createdWs, nil
}

func FindWorkspaces(filter *bson.D) ([]Workspace, *err_utils.DetailedError) {
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
		return []Workspace{}, err_utils.WrapAsDetailedError(errFind)
	}
	result := make([]Workspace, 0)
	err := cursor.All(ctx, &result)
	if err != nil {
		return nil, err_utils.WrapAsDetailedError(err)
	}
	return result, nil
}

func GetWorkspaces() ([]Workspace, *err_utils.DetailedError) {
	return FindWorkspaces(&bson.D{})
}

func GetUserWorkspaces(userID string) ([]Workspace, *err_utils.DetailedError) {
	return FindWorkspaces(&bson.D{
		{Key: "createdBy", Value: userID},
	})
}

func IsWorkspaceOwner(userId, workspaceId string) (bool, *err_utils.DetailedError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	wsCol := db.Database.Collection(db.WORKSPACES_COLLECTION)
	count, err := wsCol.CountDocuments(ctx, bson.D{
		{Key: "id", Value: workspaceId},
		{Key: "createdBy", Value: userId},
	})
	if err != nil {
		return false, err_utils.WrapAsDetailedError(err)
	}
	return count > 0, nil
}

func DeleteWorkspace(userID, wsId string) (bool, *err_utils.DetailedError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	wsCol := db.Database.Collection(db.WORKSPACES_COLLECTION)
	result, errDel := wsCol.DeleteOne(ctx, bson.D{
		{Key: "id", Value: wsId},
		{Key: "createdBy", Value: userID},
	})
	if errDel != nil {
		return false, err_utils.WrapAsDetailedError(errDel)
	}
	return result.DeletedCount > 0, nil
}
