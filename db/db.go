package db

import (
	"database/sql"
	"fmt"

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
	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s", user, password, dbName)

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

func (store *SnapInfoStore) AddSnapInfo(info SnapInfo) error {
	q := "INSERT INTO snap_info (packagename, jbid, name) VALUES (?, ?, ?)"
	_, err := store.DB.Exec(q, info.PackageName, info.JBId, info.Name)
	return err
}

func (store *SnapInfoStore) Get(name string) (SnapInfo, error) {
	q := "select jbid, name , packagename from snap_info where name =? order by id desc limit 1"
	row := store.DB.QueryRow(q, name)
	snapInfo := SnapInfo{}
	err := row.Scan(&snapInfo.JBId, &snapInfo.Name, &snapInfo.PackageName)
	fmt.Printf("result of query name - %s, jbid - %s", snapInfo.Name, snapInfo.JBId)
	return snapInfo, err
}
