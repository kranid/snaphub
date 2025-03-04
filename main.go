package main

import (
	"errors"

	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql" // Важно: "_" для инициализации драйвера
	"github.com/kranid/snaphub/config"
	"github.com/kranid/snaphub/db"
	"github.com/kranid/snaphub/jsonbin"
)

func addHandler(snapHub *SnapHub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Запущен handler")
		info := db.SnapInfo{
			Name:        r.Header.Get("name"),
			PackageName: r.Header.Get("packagename"),
		}
		if info.Name == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err := snapHub.Add(r.Body, info)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func getHandler(sh *SnapHub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")
		if name == "" {
			fmt.Println("name is not provided")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		body, err := sh.Get(name)
		if err != nil {
			fmt.Println(err)
			var NotFound *jsonbin.NotFoundError
			if errors.Is(err, sql.ErrNoRows) || errors.As(err, &NotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	}

}

func main() {
	config := config.InitConfig()
	infostore, err := db.NewSnapInfoStore(config.DBUser, config.DBPass, config.DBName)
	if err != nil {
		panic(err)
	}
	defer infostore.DB.Close()
	snapStore := jsonbin.JsonBinNew(config.ApiKey, config.AccessKey)
	snapHub, err := SnapHubNew(infostore, snapStore)
	if err != nil {
		panic(err)
	}
	m := http.NewServeMux()
	m.Handle("POST /snaphub/add", addHandler(snapHub))
	m.Handle("GET /snaphub/get/{name}", getHandler(snapHub))
	err = http.ListenAndServe(":8080", m)
	if err != nil {
		fmt.Println("Не удалось запустить сервер")
		return
	}

}
