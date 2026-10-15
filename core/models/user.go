package models

// User ánh xạ người dùng trong bảng SYS_USER
type User struct {
	Username string `json:"username"`
	FullName string `json:"fullname"`
	Email    string `json:"email"`
	Disabled bool   `json:"disabled"`
	LOC      string `json:"loc"`
}

// SessionInfo thông tin phiên làm việc sau khi đăng nhập thành công
type SessionInfo struct {
	Username  string `json:"username"`
	FullName  string `json:"fullname"`
	MaKho     string `json:"makho"`
	TenKho    string `json:"tenkho"`
	Quay      string `json:"quay"`
	Ca        string `json:"ca"`
	LOC       string `json:"loc"`
	ProfileID string `json:"profile_id"`
}
