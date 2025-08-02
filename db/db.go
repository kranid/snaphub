package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql" // Важно: "_" для инициализации драйвера
)

type SnapInfo struct {
	JBId        string
	Name        string
	PackageName string
}

type SnapInfoStore struct {
	DB *sql.DB
}

// NewStore создает новый пул соединений с базой данных
func NewSnapInfoStore(user, password, dbName string) (*SnapInfoStore, error) {
	// DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(db:3306)/%s", user, password, dbName)

	// Настройка пула соединений
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Настраиваем пул соединений
	db.SetMaxOpenConns(100)  // Максимальное количество открытых соединений
	db.SetMaxIdleConns(10)   // Максимальное количество ожидающих соединений
	db.SetConnMaxLifetime(0) // Время жизни соединения (0 - не ограничено)

	// Пинг, чтобы проверить, что соединение было успешно установлено
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &SnapInfoStore{DB: db}, nil
}

func (store *SnapInfoStore) getSnapInfo(name string) (*SnapInfo, error) {
	q := "SELECT jbId, name, pkgName FROM snap_info WHERE name = ?"
	row := store.DB.QueryRow(q, name)
	snapInfo := &SnapInfo{}

	// Заполняем snapInfo данными из строки
	err := row.Scan(&snapInfo.JBId, &snapInfo.Name, &snapInfo.PackageName)
	if err != nil {
		if err == sql.ErrNoRows {
			// Если нет результатов, возвращаем nil и nil
			return nil, nil
		}
		// В случае ошибки возвращаем nil и саму ошибку
		return nil, err
	}
	return snapInfo, nil
}

func (store *SnapInfoStore) AddSnapInfo(info SnapInfo) (int64, error) { // Изменено: теперь возвращает int64, error
	q := "INSERT INTO snap_info (packagename, jbid, name) VALUES (?, ?, ?)"
	result, err := store.DB.Exec(q, info.PackageName, info.JBId, info.Name)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	log.Printf("DB: Added snap_info record with ID: %d, name: %s", id, info.Name)
	return id, nil // Возвращаем ID
}

func (store *SnapInfoStore) Get(name string) (SnapInfo, error) {
	q := "select jbid, name , packagename from snap_info where name =? order by id desc limit 1"
	row := store.DB.QueryRow(q, name)
	snapInfo := SnapInfo{}
	err := row.Scan(&snapInfo.JBId, &snapInfo.Name, &snapInfo.PackageName)
	log.Printf("DB: Fetched snap info for name '%s', found JBId '%s'", snapInfo.Name, snapInfo.JBId)
	return snapInfo, err
}

// CreateSnapshotRecord создает новую запись в таблице snapshots и возвращает ее ID.
func (store *SnapInfoStore) CreateSnapshotRecord(packageName, activityName string) (int64, error) {
	q := "INSERT INTO snapshots (package_name, activity_name) VALUES (?, ?)"
	result, err := store.DB.Exec(q, packageName, activityName)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	log.Printf("DB: Created snapshot record for package '%s' with ID: %d", packageName, id)
	return id, nil
}

// AddSnapshotJsonLink создает новую запись в таблице snapshot_json_links.
func (store *SnapInfoStore) AddSnapshotJsonLink(snapshotID, snapInfoID int64, dataType string) error {
	q := "INSERT INTO snapshot_json_links (snapshot_id, snap_info_id, data_type) VALUES (?, ?, ?)"
	_, err := store.DB.Exec(q, snapshotID, snapInfoID, dataType)
	if err != nil {
		log.Printf("ERROR: Failed to add snapshot JSON link: %v", err)
		return err
	}
	log.Printf("DB: Added snapshot JSON link: snapshotID=%d, snapInfoID=%d, dataType=%s", snapshotID, snapInfoID, dataType)
	return nil
}
