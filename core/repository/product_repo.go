package repository

import (
	"context"
	"database/sql"
	"fmt"

	"shophoa/core/models"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// SearchProducts tìm kiếm nhanh sản phẩm theo mã, tên hoặc barcode
func (r *ProductRepository) SearchProducts(ctx context.Context, keyword string, limit int) ([]models.ProductSearchResult, error) {
	if limit <= 0 {
		limit = 20
	}

	query := `
		SELECT TOP (?) 
			MAVT, ISNULL(TENVT, ''), ISNULL(DVT, ''), 
			ISNULL(GIABAN, 0), ISNULL(LOAITHUE, 'T10'), ISNULL(VAT_RA, 10)
		FROM DM_HH
		WHERE (CHECKED = 1 OR CHECKED IS NULL)
		  AND (MAVT LIKE ? OR TENVT LIKE ? OR PLU = ?)
		ORDER BY MAVT
	`

	kw := "%" + keyword + "%"
	rows, err := r.db.QueryContext(ctx, query, limit, kw, kw, keyword)
	if err != nil {
		return nil, fmt.Errorf("failed to search products: %w", err)
	}
	defer rows.Close()

	var list []models.ProductSearchResult
	for rows.Next() {
		var p models.ProductSearchResult
		if err := rows.Scan(&p.MaVT, &p.TenVT, &p.DVT, &p.GiaBan, &p.LoaiThue, &p.VatRa); err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		list = append(list, p)
	}

	return list, nil
}

// GetByMaVT lấy thông tin chi tiết một mặt hàng theo mã MAVT
func (r *ProductRepository) GetByMaVT(ctx context.Context, mavt string) (*models.Product, error) {
	query := `
		SELECT MAVT, ISNULL(TENVT, ''), ISNULL(DVT, ''), ISNULL(GIABAN, 0),
		       ISNULL(GIABAN1, 0), ISNULL(GIASI, 0), ISNULL(GIAVON, 0),
		       ISNULL(LOAITHUE, 'T10'), ISNULL(VAT_RA, 10), ISNULL(MANHOM, ''),
		       ISNULL(MANGANH, ''), ISNULL(TINHTRANG, ''), ISNULL(PLU, '')
		FROM DM_HH
		WHERE MAVT = ?
	`

	var p models.Product
	err := r.db.QueryRowContext(ctx, query, mavt).Scan(
		&p.MaVT, &p.TenVT, &p.DVT, &p.GiaBan,
		&p.GiaBan1, &p.GiaSi, &p.GiaVon,
		&p.LoaiThue, &p.VatRa, &p.MaNhom,
		&p.MaNganh, &p.TinhTrang, &p.PLU,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get product by mavt %s: %w", mavt, err)
	}

	return &p, nil
}

// GetTopCategories lấy danh sách các nhóm hàng có nhiều sản phẩm nhất
func (r *ProductRepository) GetTopCategories(ctx context.Context) ([]models.Category, error) {
	query := `
		SELECT TOP 8 n.MANHOM, ISNULL(n.TENNHOM, n.MANHOM), COUNT(*) AS Total
		FROM DM_HH h
		JOIN DM_NHOM n ON h.MANHOM = n.MANHOM
		WHERE (h.CHECKED = 1 OR h.CHECKED IS NULL)
		GROUP BY n.MANHOM, n.TENNHOM
		ORDER BY Total DESC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get top categories: %w", err)
	}
	defer rows.Close()

	var list []models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.MaNhom, &c.TenNhom, &c.Count); err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, nil
}

// SearchProductsWithCategory tìm kiếm sản phẩm theo từ khóa và nhóm hàng
func (r *ProductRepository) SearchProductsWithCategory(ctx context.Context, keyword string, manhom string, limit int) ([]models.ProductSearchResult, error) {
	if limit <= 0 {
		limit = 30
	}

	whereClause := "WHERE (CHECKED = 1 OR CHECKED IS NULL)"
	var args []interface{}

	if manhom != "" {
		whereClause += " AND MANHOM = ?"
		args = append(args, manhom)
	}

	if keyword != "" {
		whereClause += " AND (MAVT LIKE ? OR TENVT LIKE ? OR PLU = ?)"
		kw := "%" + keyword + "%"
		args = append(args, kw, kw, keyword)
	}

	query := fmt.Sprintf(`
		SELECT TOP (%d) 
			MAVT, ISNULL(TENVT, ''), ISNULL(DVT, ''), 
			ISNULL(GIABAN, 0), ISNULL(LOAITHUE, 'T10'), ISNULL(VAT_RA, 10)
		FROM DM_HH
		%s
		ORDER BY MAVT
	`, limit, whereClause)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search products: %w", err)
	}
	defer rows.Close()

	var list []models.ProductSearchResult
	for rows.Next() {
		var p models.ProductSearchResult
		if err := rows.Scan(&p.MaVT, &p.TenVT, &p.DVT, &p.GiaBan, &p.LoaiThue, &p.VatRa); err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		list = append(list, p)
	}
	return list, nil
}
