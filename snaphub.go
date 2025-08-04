package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"

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

func (sh *SnapHub) Add(body io.Reader, info db.SnapInfo) (int64, error) {
	log.Printf("SnapHub: Adding snap with name = %s, packagename =%s", info.Name, info.PackageName)
	jsonBinID, err := sh.SnapStore.Create(body, info.Name)
	if err != nil {
		log.Printf("WARN: SnapHub failed to create JSONBin for %s: %v. Proceeding without JSONBin ID.", info.Name, err)
		info.JBId = "" // Ensure JBId is empty if creation failed
	} else {
		info.JBId = jsonBinID
		log.Printf("SnapHub: Successfully added snap with JSONBin ID: %s", jsonBinID)
	}

	snapInfoID, err := sh.InfoStore.AddSnapInfo(info) // Получаем ID из AddSnapInfo
	if err != nil {
		log.Printf("ERROR: SnapHub failed to add snap info to DB: %v", err)
		return 0, err
	}
	log.Printf("SnapHub: Successfully added snap info to DB with ID: %d", snapInfoID)
	return snapInfoID, nil // Возвращаем ID из DB
}

func (sh *SnapHub) Get(name string) (json.RawMessage, error) {
	info, err := sh.InfoStore.Get(name)
	if err != nil {
		log.Printf("ERROR: SnapHub failed to get snap info from DB for name %s: %v", name, err)
		return nil, err
	}
	body, err := sh.SnapStore.Get(info.JBId)
	if err != nil {
		log.Printf("ERROR: SnapHub failed to get snap from JSONBin for ID %s: %v", info.JBId, err)
		return nil, err
	}
	return body, nil
}