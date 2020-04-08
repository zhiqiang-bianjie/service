package service

import (
	"github.com/irismod/service/keeper"
	"github.com/irismod/service/types"
)

const (
	ModuleName        = types.ModuleName
	StoreKey          = types.StoreKey
	QuerierRoute      = types.QuerierRoute
	RouterKey         = types.RouterKey
	DefaultParamspace = types.DefaultParamspace
)

var (
	NewKeeper     = keeper.NewKeeper
	NewQuerier    = keeper.NewQuerier
	RegisterCodec = types.RegisterCodec
)

type (
	Keeper = keeper.Keeper
)
