package main

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/kranid/snaphub/db"
	"github.com/kranid/snaphub/jsonbin"
)

type SnapHub struct {
	InfoStore *db.SnapInfoStore
	SnapStore jsonbin.SnapStore
}

func SnapHubNew(infoStore *db.SnapInfoStore, snapStore jsonbin.SnapStore) (*SnapHub, error) {
	if infoStore == nil {
		return nil, fmt.Errorf("db.SnapInfoStore can not be nil")
	}
	return &SnapHub{
		InfoStore: infoStore,
		SnapStore: snapStore,
	}, nil
}

func (sh *SnapHub) Add(body io.Reader, info db.SnapInfo) error {
	fmt.Printf("name = %s, packagename =%s \n", info.Name, info.PackageName)
	id, err := sh.SnapStore.Create(body, info.Name)
	if err != nil {
		return err
	}
	info.JBId = id
	err = sh.InfoStore.AddSnapInfo(info)
	return err
}

func (sh *SnapHub) Get(name string) (json.RawMessage, error) {
	info, err := sh.InfoStore.Get(name)
	if err != nil {

		return nil, err
	}
	body, err := sh.SnapStore.Get(info.JBId)
	if err != nil {
		return nil, err
	}
	return body, nil
}
