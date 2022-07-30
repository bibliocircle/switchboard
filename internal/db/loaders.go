package db

import "github.com/graph-gophers/dataloader"

const LoadersCtxKey = "loaders"

type Loaders struct {
	Scenarios *dataloader.Loader
	Endpoints *dataloader.Loader
}
