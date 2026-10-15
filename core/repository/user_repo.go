package repository

import (
	"context"
	"crypto/md5"
	"database/sql"
	"fmt"
	"strings"

	"shophoa/core/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Authenticate kiểm tra tài khoản đăng nhập trong SYS_USER bằng MD5 hash
func (r *UserRepository) Authenticate(ctx context.Context, username, password string) (*models.User, error) {
	query := `
		SELECT USERNAME, ISNULL(FULLNAME, USERNAME), ISNULL(EMAIL, ''),
		       ISNULL(DISABLED, 0), ISNULL(LOC, '01'), ISNULL(PASSWD, '')
		FROM SYS_USER
		WHERE UPPER(USERNAME) = UPPER(?)
	`
	var u models.User
	var storedHash string
	err := r.db.QueryRowContext(ctx, query, strings.TrimSpace(username)).Scan(
		&u.Username, &u.FullName, &u.Email, &u.Disabled, &u.LOC, &storedHash,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("tài khoản '%s' không tồn tại", username)
		}
		return nil, fmt.Errorf("lỗi kiểm tra người dùng: %w", err)
	}

	if u.Disabled {
		return nil, fmt.Errorf("tài khoản '%s' đang bị tạm khóa", username)
	}

	// Kiểm tra mật khẩu (MD5 hash hoặc chuỗi rỗng)
	inputHash := fmt.Sprintf("%x", md5.Sum([]byte(password)))
	storedHash = strings.TrimSpace(storedHash)

	// Các trường hợp hợp lệ:
	// 1. Mật khẩu trong DB để trống hoặc khớp hash MD5 của chuỗi rỗng
	// 2. Hash MD5 của input khớp với hash trong DB
	// 3. Khớp trực tiếp chuỗi mật khẩu
	match := false
	if (storedHash == "" || storedHash == "d41d8cd98f00b204e9800998ecf8427e") && password == "" {
		match = true
	} else if strings.EqualFold(storedHash, inputHash) {
		match = true
	} else if storedHash == password {
		match = true
	}

	if !match {
		return nil, fmt.Errorf("mật khẩu đăng nhập không chính xác")
	}

	return &u, nil
}

// GetWarehouses lấy danh sách kho bán lẻ từ DM_KHO
func (r *UserRepository) GetWarehouses(ctx context.Context) ([]models.Warehouse, error) {
	query := `
		SELECT MAKHO, ISNULL(TENKHO, MAKHO)
		FROM DM_KHO
		ORDER BY MAKHO
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query DM_KHO: %w", err)
	}
	defer rows.Close()

	var list []models.Warehouse
	for rows.Next() {
		var w models.Warehouse
		if err := rows.Scan(&w.MaKho, &w.TenKho); err != nil {
			return nil, err
		}
		list = append(list, w)
	}
	return list, nil
}

// GetWarehouseName lấy tên kho theo mã
func (r *UserRepository) GetWarehouseName(ctx context.Context, makho string) string {
	query := "SELECT ISNULL(TENKHO, '') FROM DM_KHO WHERE MAKHO = ?"
	var name string
	_ = r.db.QueryRowContext(ctx, query, makho).Scan(&name)
	if name == "" {
		return "Kho " + makho
	}
	return name
}
