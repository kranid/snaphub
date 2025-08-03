package main

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	_ "github.com/go-sql-driver/mysql" // Важно: "_" для инициализации драйвера
	"github.com/kranid/snaphub/config"
	"github.com/kranid/snaphub/db"
	"github.com/kranid/snaphub/jsonbin"
)

func addHandler(snapHub *SnapHub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		info := db.SnapInfo{
			Name:        r.Header.Get("name"),
			PackageName: r.Header.Get("packagename"),
		}
		log.Printf("INFO: addHandler received request for name: %s, package: %s", info.Name, info.PackageName)
		if info.Name == "" {
			log.Println("WARN: addHandler missing name in header")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, err := snapHub.Add(r.Body, info)
		if err != nil {
			log.Printf("ERROR: addHandler failed to add snap: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func getHandler(sh *SnapHub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")
		log.Printf("INFO: getHandler received request for name: %s", name)
		if name == "" {
			log.Println("WARN: getHandler missing name in path")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		body, err := sh.Get(name)
		if err != nil {
			var NotFound *jsonbin.NotFoundError
			if errors.Is(err, sql.ErrNoRows) || errors.As(err, &NotFound) {
				log.Printf("INFO: getHandler snap not found for name: %s", name)
				w.WriteHeader(http.StatusNotFound)
				return
			}
			log.Printf("ERROR: getHandler failed to get snap: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	}
}

func addSnapshotHandler(sh *SnapHub) http.HandlerFunc { // Изменено: теперь принимает *SnapHub
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("INFO: addSnapshotHandler received request")
		// 1. Парсинг multipart/form-data
		err := r.ParseMultipartForm(32 << 20) // 32MB max memory
		if err != nil {
			log.Printf("ERROR: Failed to parse multipart form: %v", err)
			http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
			return
		}

		// 2. Получение текстовых полей
		packageName := r.FormValue("package_name")
		activityName := r.FormValue("activity_name")
		log.Printf("INFO: addSnapshotHandler params: package_name='%s', activity_name='%s'", packageName, activityName)
		if packageName == "" || activityName == "" {
			log.Println("WARN: addSnapshotHandler missing package_name or activity_name")
			http.Error(w, "package_name and activity_name are required", http.StatusBadRequest)
			return
		}

		// 3. Создание записи в БД (таблица snapshots)
		snapshotID, err := sh.InfoStore.CreateSnapshotRecord(packageName, activityName)
		if err != nil {
			log.Printf("ERROR: Failed to create snapshot record: %v", err)
			http.Error(w, "Failed to create snapshot record", http.StatusInternalServerError)
			return
		}
		log.Printf("INFO: Created new snapshot record with ID: %d", snapshotID)

		// 4. Создание папки для снапшотов на диске
		snapshotDir := filepath.Join("snapshots", packageName, strconv.FormatInt(snapshotID, 10))
		if err := os.MkdirAll(snapshotDir, os.ModePerm); err != nil {
			log.Printf("ERROR: Failed to create snapshot directory '%s': %v", snapshotDir, err)
			http.Error(w, "Failed to create snapshot directory", http.StatusInternalServerError)
			return
		}
		log.Printf("INFO: Created snapshot directory: %s", snapshotDir)

		// 5. Сохранение JSON-файлов в JSONBin и связывание в БД
		jsonFilesToUpload := map[string]string{
			"original_snapshot": "original",
			"expected_snapshot": "expected",
			"tech_report":       "technical_report",
			"human_report":      "human_report",
		}

		for formFileName, dataType := range jsonFilesToUpload {
			file, _, err := r.FormFile(formFileName)
			if err != nil {
				log.Printf("WARN: Failed to get JSON file %s from form: %v", formFileName, err)
				// Не возвращаем ошибку, если JSON-файл отсутствует, так как он может быть необязательным
				continue
			}
			defer file.Close()

			// Создаем уникальное имя для JSONBin, используя snapshotID
			jsonBinName := fmt.Sprintf("%d_%s", snapshotID, dataType)

			// Создаем SnapInfo для передачи в SnapHub.Add
			info := db.SnapInfo{
				Name:        jsonBinName,
				PackageName: packageName,
			}

			snapInfoID, err := sh.Add(file, info) // Используем sh.Add
			if err != nil {
				log.Printf("ERROR: Failed to add %s via SnapHub.Add: %v", formFileName, err)
				http.Error(w, fmt.Sprintf("Failed to add %s", formFileName), http.StatusInternalServerError)
				return
			}
			log.Printf("INFO: Successfully added %s to JSONBin and DB with snap_info ID: %d", formFileName, snapInfoID)

			// Связываем в snapshot_json_links
			err = sh.InfoStore.AddSnapshotJsonLink(snapshotID, snapInfoID, dataType)
			if err != nil {
				log.Printf("ERROR: Failed to link %s (snap_info ID: %d) to snapshot record: %v", formFileName, snapInfoID, err)
				http.Error(w, fmt.Sprintf("Failed to link %s to snapshot record", formFileName), http.StatusInternalServerError)
				return
			}
			log.Printf("INFO: Linked %s (snap_info ID: %d) to snapshot ID: %d", formFileName, snapInfoID, snapshotID)
		}

		// 6. Сохранение скриншота на диск
		screenshotFile, _, err := r.FormFile("screenshot")
		if err != nil {
			log.Printf("WARN: Failed to get screenshot file from form: %v", err)
			// Скриншот может быть необязательным, поэтому не возвращаем ошибку
		} else {
			defer screenshotFile.Close()
			dstPath := filepath.Join(snapshotDir, "screenshot.jpg")
			dst, err := os.Create(dstPath)
			if err != nil {
				log.Printf("ERROR: Failed to create screenshot file '%s': %v", dstPath, err)
				http.Error(w, "Failed to create screenshot file", http.StatusInternalServerError)
				return
			}
			defer dst.Close()

			if _, err := io.Copy(dst, screenshotFile); err != nil {
				log.Printf("ERROR: Failed to save screenshot file '%s': %v", dstPath, err)
				http.Error(w, "Failed to save screenshot file", http.StatusInternalServerError)
				return
			}
			log.Printf("INFO: Successfully saved screenshot: %s", dstPath)
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Snapshot saved successfully with ID: %d", snapshotID)
	}
}

func main() {
	log.Println("INFO: Starting snaphub application...")
	config := config.InitConfig()
	infostore, err := db.NewSnapInfoStore(config.DBUser, config.DBPass, config.DBName)
	if err != nil {
		log.Fatalf("FATAL: Failed to connect to database: %v", err)
	}
	defer infostore.DB.Close()
	log.Println("INFO: Database connection successful.")

	snapStore := jsonbin.JsonBinNew(config.MasterKey, config.AccessKey)
	snapHub, err := SnapHubNew(infostore, snapStore)
	if err != nil {
		log.Fatalf("FATAL: Failed to create SnapHub: %v", err)
	}

	m := http.NewServeMux()
	m.Handle("POST /snaphub/add", addHandler(snapHub))
	m.Handle("GET /snaphub/get/{name}", getHandler(snapHub))
	m.Handle("POST /snapshots/add", addSnapshotHandler(snapHub)) // Изменено: теперь передаем snapHub

	log.Println("INFO: Starting HTTP server on :8080")
	err = http.ListenAndServe(":8080", m)
	if err != nil {
		log.Fatalf("FATAL: Failed to start HTTP server: %v", err)
	}
}