package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/alexbrainman/odbc"
)

// DBProfile thông tin cấu hình kết nối CSDL SQL Server
type DBProfile struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Server   string `json:"server"`
	Database string `json:"database"`
	User     string `json:"user"`
	Password string `json:"password"`
	Driver   string `json:"driver"`
	MaKho    string `json:"makho"`
	Quay     string `json:"quay"`
}

// ProfilesConfig cấu trúc lưu trữ danh sách các cấu hình và cấu hình đang chọn gần nhất
type ProfilesConfig struct {
	ActiveProfileID string      `json:"active_profile_id"`
	Profiles        []DBProfile `json:"profiles"`
}

const profilesFileName = "db_profiles.json"

// getProfilesFilePath trả về đường dẫn file lưu cấu hình
func getProfilesFilePath() string {
	exePath, err := os.Executable()
	if err == nil {
		dir := filepath.Dir(exePath)
		return filepath.Join(dir, profilesFileName)
	}
	return profilesFileName
}

// LoadProfilesConfig tải danh sách kết nối từ file, khởi tạo mặc định nếu chưa có
func LoadProfilesConfig() *ProfilesConfig {
	path := getProfilesFilePath()
	data, err := os.ReadFile(path)
	if err == nil {
		var cfg ProfilesConfig
		if err := json.Unmarshal(data, &cfg); err == nil && len(cfg.Profiles) > 0 {
			return &cfg
		}
	}

	// Cấu hình mặc định ban đầu
	defaultProfile := DBProfile{
		ID:       "default",
		Name:     "Cửa hàng Q2 (.\\sqlexpress - Mặc định)",
		Server:   ".\\sqlexpress",
		Database: "SOFTZ_SHOPHOA",
		User:     "sa",
		Password: "123@456",
		Driver:   "ODBC Driver 17 for SQL Server",
		MaKho:    "01",
		Quay:     "01",
	}

	cfg := &ProfilesConfig{
		ActiveProfileID: defaultProfile.ID,
		Profiles:        []DBProfile{defaultProfile},
	}

	_ = SaveProfilesConfig(cfg)
	return cfg
}

// SaveProfilesConfig ghi cấu hình ra file
func SaveProfilesConfig(cfg *ProfilesConfig) error {
	path := getProfilesFilePath()
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// ConnectProfile tạo kết nối CSDL từ profile chỉ định
func ConnectProfile(p DBProfile) (*sql.DB, error) {
	if p.Driver == "" {
		p.Driver = "ODBC Driver 17 for SQL Server"
	}
	connStr := fmt.Sprintf("Driver={%s};Server=%s;Database=%s;Uid=%s;Pwd=%s;",
		p.Driver, p.Server, p.Database, p.User, p.Password)

	db, err := sql.Open("odbc", connStr)
	if err != nil {
		return nil, fmt.Errorf("không thể mở kết nối ODBC: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("kết nối thất bại tới %s [%s]: %w", p.Server, p.Database, err)
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	return db, nil
}
