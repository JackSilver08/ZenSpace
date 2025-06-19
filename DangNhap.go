package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

// ĐĂNG KÝ TÀI KHOẢN
// Cấu trúc dữ liệu cho đăng ký
type YeuCauDangKy struct {  // JSON = Javascript Object Notation ( 1 định dạng phổ biến để truyền dữ liệu trong GO)
    TenDangNhap string `json:"username"` // Định nghĩa cấu trúc ánh xạ JSON
    Email       string `json:"email"`
    MatKhau     string `json:"password"`
    XacNhanMK   string `json:"confirmPassword"`
}
func DangKy(w http.ResponseWriter, r *http.Request) {
    // CORS headers : Ý nghĩa của CORS (Cross-Origin Resource Sharing) là một cơ chế bảo mật của trình duyệt cho phép (hoặc từ chối) một website (frontend) gửi yêu cầu HTTP đến một domain khác (backend)
    w.Header().Set("Access-Control-Allow-Origin", "*") // Cho phép mọi nguồn (domain) truy cập vào tài nguyên của bạn (API/server)
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type") // Cho phép phía client (trình duyệt hoặc ứng dụng) được gửi header Content-Type trong yêu cầu
    w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS") // Cho phép phía client gọi các HTTP method như: POST: Gửi dữ liệu, OPTIONS: Dùng cho preflight request, tức là trình duyệt kiểm tra trước xem có được gửi POST thật hay không
    w.Header().Set("Content-Type", "application/json") // Server trả về nội dung kiểu JSON
    // Đây là kiểm tra xem request có phải là phương thức OPTIONS không
    if r.Method == http.MethodOptions { 
        w.WriteHeader(http.StatusOK)
        return
    }
    // Đảm bảo rằng endpoint này chỉ dùng để xử lý POST, ví dụ như đăng nhập, đăng ký, gửi dữ liệu JSON
    if r.Method != http.MethodPost {
        http.Error(w, "Chỉ chấp nhận POST", http.StatusMethodNotAllowed)
        return
    }

    // Phần xử lý đăng ký vẫn giữ nguyên từ đây...
    var req YeuCauDangKy
    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        http.Error(w, "Lỗi khi đọc dữ liệu", http.StatusBadRequest)
        return
    }

    if req.MatKhau != req.XacNhanMK {
        json.NewEncoder(w).Encode(PhanHoi{ThanhCong: false, ThongBao: "Mật khẩu xác nhận không khớp"})
        return
    }

    db, err := sql.Open("mysql", "root:Tuan@1234@tcp(127.0.0.1:3306)/ZenSpaceDB")
    if err != nil {
        http.Error(w, "Không thể kết nối database", http.StatusInternalServerError)
        return
    }
    defer db.Close()

    var exists int
    err = db.QueryRow("SELECT COUNT(*) FROM TaiKhoan WHERE tenDangNhap = ?", req.TenDangNhap).Scan(&exists)
    if err != nil || exists > 0 {
        json.NewEncoder(w).Encode(PhanHoi{ThanhCong: false, ThongBao: "Tên đăng nhập đã tồn tại"})
        return
    }

    res, err := db.Exec("INSERT INTO TaiKhoan (tenDangNhap, email, matKhau) VALUES (?, ?, ?)", req.TenDangNhap, req.Email, req.MatKhau)
    if err != nil {
        json.NewEncoder(w).Encode(PhanHoi{ThanhCong: false, ThongBao: "Lỗi khi tạo tài khoản: " + err.Error()})
        return
    }

    rowsAffected, _ := res.RowsAffected()
    if rowsAffected == 0 {
        json.NewEncoder(w).Encode(PhanHoi{ThanhCong: false, ThongBao: "Không thể thêm tài khoản"})
        return
    }

    json.NewEncoder(w).Encode(PhanHoi{ThanhCong: true, ThongBao: "Đăng ký thành công"})
}

// ĐĂNG NHẬP
// Cấu trúc để nhận dữ liệu JSON từ client
type YeuCauDangNhap struct {
	TenDangNhap string `json:"TenDangNhap"`
	MatKhau     string `json:"MatKhau"`
}

type PhanHoi struct {
	ThanhCong bool        `json:"success"`
	ThongBao  string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
}


func DangNhap(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
    w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
    if r.Method == "OPTIONS" {
        w.WriteHeader(http.StatusOK)
        return
    }
    if r.Method != "POST" {
        http.Error(w, "Chỉ chấp nhận POST", http.StatusMethodNotAllowed)
        return
    }
    w.Header().Set("Content-Type", "application/json")

    var req YeuCauDangNhap
    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        http.Error(w, "Lỗi khi đọc dữ liệu", http.StatusBadRequest)
        return
    }

    db, err := sql.Open("mysql", "root:Tuan@1234@tcp(127.0.0.1:3306)/ZenSpaceDB")
    if err != nil {
        http.Error(w, "Không thể kết nối database", http.StatusInternalServerError)
        return
    }
    defer db.Close()

    // Lấy mật khẩu đã lưu trong db theo username
    var storedPassword string
    err = db.QueryRow("SELECT matKhau FROM TaiKhoan WHERE tenDangNhap = ?", req.TenDangNhap).Scan(&storedPassword)
    if err != nil {
        if err == sql.ErrNoRows {
            json.NewEncoder(w).Encode(PhanHoi{ThanhCong: false, ThongBao: "Tài khoản không tồn tại"})
        } else {
            json.NewEncoder(w).Encode(PhanHoi{ThanhCong: false, ThongBao: "Lỗi truy vấn database: " + err.Error()})
        }
        return
    }

    // So sánh mật khẩu
    if storedPassword != req.MatKhau {
        json.NewEncoder(w).Encode(PhanHoi{ThanhCong: false, ThongBao: "Mật khẩu không đúng"})
        return
    }

    // Lấy thêm thông tin email và id (hoặc bạn có thể lấy trong cùng query trên)
    var id int
    var email string
    err = db.QueryRow("SELECT id, email FROM TaiKhoan WHERE tenDangNhap = ?", req.TenDangNhap).Scan(&id, &email)
    if err != nil {
        json.NewEncoder(w).Encode(PhanHoi{ThanhCong: false, ThongBao: "Lỗi lấy thông tin tài khoản: " + err.Error()})
        return
    }

  nguoiDung := map[string]interface{}{
    "id":     id,
    "hoten":  req.TenDangNhap,
    "avatar": "../IMG/ZenUser.png", // hoặc để rỗng
}

json.NewEncoder(w).Encode(map[string]interface{}{
    "success": true,
    "user":    nguoiDung,
})


}
