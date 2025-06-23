package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	redisv9 "github.com/redis/go-redis/v9"

	// các import sẵn có của bạn ở trên...

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/microcosm-cc/bluemonday"
)
var ctx = context.Background()

var redisClient = redisv9.NewClient(&redisv9.Options{
	Addr:     "localhost:6379", // Nếu Redis chạy local
	Password: "",               // Nếu Redis không có mật khẩu
	DB:       0,                // DB mặc định là 0
})
type BaiDang struct {
	TieuDe     string `json:"tieuDe"`
	NoiDung    string `json:"noiDung"`
	IDTaiKhoan int    `json:"idTaiKhoan"`
}


func renderMarkdown(input string) string {
	// Chuyển từ markdown sang HTML
	output := markdown.ToHTML([]byte(input), nil, html.NewRenderer(html.RendererOptions{
		Flags: html.CommonFlags | html.HrefTargetBlank,
	}))

	// Làm sạch HTML để tránh XSS
	policy := bluemonday.UGCPolicy()
	safe := policy.SanitizeBytes(output)

	return string(safe)
}

// ========== Đăng bài ==========
func DangBai(w http.ResponseWriter, r *http.Request) {
	// CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
	w.Header().Set("Content-Type", "application/json")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Chỉ hỗ trợ POST", http.StatusMethodNotAllowed)
		return
	}

	// ✅ Lấy user từ token JWT
	token, err := validateToken(r)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Unauthorized: Invalid claims", http.StatusUnauthorized)
		return
	}

	idTaiKhoanFloat, ok := claims["id"].(float64)
	if !ok {
		http.Error(w, "Unauthorized: Invalid user ID in token", http.StatusUnauthorized)
		return
	}
	idTaiKhoan := int(idTaiKhoanFloat)

	// ✅ Decode bài đăng nhưng KHÔNG nhận IDTaiKhoan từ client
	var bai struct {
		TieuDe  string `json:"tieuDe"`
		NoiDung string `json:"noiDung"`
	}
	if err := json.NewDecoder(r.Body).Decode(&bai); err != nil {
		http.Error(w, "Dữ liệu không hợp lệ", http.StatusBadRequest)
		return
	}

	// ✅ Lưu bài đăng với idTaiKhoan từ token
	stmt := `INSERT INTO baidang (tieuDe, noiDung, idTaiKhoan, ngayDang) VALUES (?, ?, ?, ?)`
	if _, err := DB.ExecContext(ctx, stmt, bai.TieuDe, bai.NoiDung, idTaiKhoan, time.Now()); err != nil {
		log.Println("Lỗi ghi DB:", err)
		http.Error(w, "Không thể lưu bài viết", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Đăng bài thành công",
	})
}


// ========== Lấy tất cả bài đăng ==========
func LayBaiDang(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
	w.Header().Set("Content-Type", "application/json")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Chỉ hỗ trợ GET", http.StatusMethodNotAllowed)
		return
	}

	rows, err := DB.QueryContext(ctx, `
		SELECT b.id, b.tieuDe, b.noiDung, b.idTaiKhoan, t.tenDangNhap, b.ngayDang
		FROM baidang b
		JOIN taikhoan t ON b.idTaiKhoan = t.id
		ORDER BY b.ngayDang DESC`)
	if err != nil {
		http.Error(w, "Không thể lấy bài viết", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var danhSach []map[string]interface{}
	for rows.Next() {
		var id, idTaiKhoan int
		var tieuDe, noiDung, tenDangNhap string
		var ngayDang time.Time

		if err := rows.Scan(&id, &tieuDe, &noiDung, &idTaiKhoan, &tenDangNhap, &ngayDang); err != nil {
			log.Println("Lỗi đọc hàng:", err)
			continue
		}

		danhSach = append(danhSach, map[string]interface{}{
			"id":           id,
			"tieuDe":       tieuDe,
			"noiDung":      renderMarkdown(noiDung),
			"idTaiKhoan":   idTaiKhoan,
			"tenDangNhap":  tenDangNhap,
			"ngayDang":     ngayDang.Format("02/01/2006 15:04"),
		})
	}

	json.NewEncoder(w).Encode(danhSach)
}

// ========== Lấy chi tiết bài đăng ==========
func LayChiTietBaiDang(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
	w.Header().Set("Content-Type", "application/json")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID không hợp lệ", http.StatusBadRequest)
		return
	}

	var tieuDe, noiDung string
	var idTaiKhoan int
	var ngayDang time.Time

	err = DB.QueryRowContext(ctx,
		`SELECT tieuDe, noiDung, idTaiKhoan, ngayDang FROM baidang WHERE id = ?`, id).
		Scan(&tieuDe, &noiDung, &idTaiKhoan, &ngayDang)
	if err != nil {
		http.Error(w, "Không tìm thấy bài viết", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":         id,
		"tieuDe":     tieuDe,
		"noiDung":    renderMarkdown(noiDung),
		"idTaiKhoan": idTaiKhoan,
		"ngayDang":   ngayDang.Format("02/01/2006 15:04"),
	})
}


// ========== Xóa bài đăng ==========
func XoaBaiDang(w http.ResponseWriter, r *http.Request) {
	// --- Cấu hình CORS và header ---
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
	w.Header().Set("Content-Type", "application/json")

	// --- Xử lý preflight request ---
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// --- Chỉ cho phép DELETE ---
	if r.Method != http.MethodDelete {
		http.Error(w, "Chỉ hỗ trợ DELETE", http.StatusMethodNotAllowed)
		return
	}

	// --- Xác thực token ---
	token, err := validateToken(r)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Unauthorized: Invalid token claims", http.StatusUnauthorized)
		return
	}

	// --- Lấy thông tin người dùng từ token ---
	idFloat, okID := claims["id"].(float64)
	username, _ := claims["tenDangNhap"].(string)

	// Kiểm tra quyền hợp lệ
	quyenRaw, exists := claims["phanQuyen"]
	if !okID || !exists || quyenRaw == nil {
		http.Error(w, "Unauthorized: Thiếu thông tin trong token", http.StatusUnauthorized)
		return
	}
	quyen, ok := quyenRaw.(string)
	if !ok {
		http.Error(w, "Unauthorized: Quyền không hợp lệ", http.StatusUnauthorized)
		return
	}

	idNguoiDung := int(idFloat)
	fmt.Printf("🧾 Token người dùng: ID=%d | Tên=%s | Quyền=%s\n", idNguoiDung, username, quyen)

	// --- Lấy ID bài đăng từ URL ---
	idStr := mux.Vars(r)["id"]
	idBaiDang, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID bài đăng không hợp lệ", http.StatusBadRequest)
		return
	}

	// --- Truy vấn ID chủ bài viết từ CSDL ---
	var idChuBai int
	err = DB.QueryRowContext(ctx, "SELECT idTaiKhoan FROM baidang WHERE id = ?", idBaiDang).Scan(&idChuBai)
	if err == sql.ErrNoRows {
		http.Error(w, "Bài viết không tồn tại", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Lỗi truy vấn bài viết", http.StatusInternalServerError)
		return
	}

	// --- Kiểm tra quyền xóa ---
	if idNguoiDung != idChuBai && quyen != "admin" {
		http.Error(w, "Bạn không có quyền xóa bài viết này", http.StatusForbidden)
		return
	}

	// --- Bắt đầu transaction ---
	tx, err := DB.BeginTx(ctx, nil)
	if err != nil {
		http.Error(w, "Không thể khởi tạo giao dịch", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// --- Xóa dữ liệu liên quan ---
	if _, err := tx.ExecContext(ctx, "DELETE FROM binhluan WHERE idBaiDang = ?", idBaiDang); err != nil {
		http.Error(w, "Lỗi khi xóa bình luận", http.StatusInternalServerError)
		return
	}
	if _, err := tx.ExecContext(ctx, "DELETE FROM cam_xuc WHERE idBaiDang = ?", idBaiDang); err != nil {
		http.Error(w, "Lỗi khi xóa cảm xúc", http.StatusInternalServerError)
		return
	}

	// --- Xóa bài đăng ---
	res, err := tx.ExecContext(ctx, "DELETE FROM baidang WHERE id = ?", idBaiDang)
	if err != nil {
		http.Error(w, "Lỗi khi xóa bài viết", http.StatusInternalServerError)
		return
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Không tìm thấy bài viết để xóa", http.StatusNotFound)
		return
	}

	// --- Xác nhận transaction ---
	if err := tx.Commit(); err != nil {
		http.Error(w, "Lỗi khi xác nhận xóa bài", http.StatusInternalServerError)
		return
	}

	// --- Phản hồi thành công ---
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Đã xóa bài viết và dữ liệu liên quan",
		"id":      idBaiDang,
	})
}



func ThemCamXuc(w http.ResponseWriter, r *http.Request) {
	// CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Chỉ hỗ trợ phương thức POST", http.StatusMethodNotAllowed)
		return
	}

	// --- Xác thực token ---
	token, err := validateToken(r)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Unauthorized: Invalid claims", http.StatusUnauthorized)
		return
	}

	idTaiKhoanFloat, ok := claims["id"].(float64)
	if !ok {
		http.Error(w, "Unauthorized: Invalid user ID in token", http.StatusUnauthorized)
		return
	}
	idTaiKhoan := int(idTaiKhoanFloat)

	// --- Parse request body ---
	var request struct {
		LoaiCamXuc string `json:"loaiCamXuc"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Dữ liệu không hợp lệ", http.StatusBadRequest)
		return
	}

	// Lấy idBaiDang từ URL params
	params := mux.Vars(r)
	idBaiDang, err := strconv.Atoi(params["idBaiDang"])
	if err != nil {
		http.Error(w, "ID bài viết không hợp lệ", http.StatusBadRequest)
		return
	}

	// Kiểm tra loại cảm xúc hợp lệ
	validEmotions := map[string]bool{
		"like": true, "dislike": true, "love": true,
		"haha": true, "wow": true, "sad": true, "angry": true,
	}
	if !validEmotions[request.LoaiCamXuc] {
		http.Error(w, "Loại cảm xúc không hợp lệ", http.StatusBadRequest)
		return
	}


	tx, err := DB.Begin()
	if err != nil {
		log.Println("Lỗi khởi tạo transaction:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	var exists bool
	err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM baidang WHERE id = ?)", idBaiDang).Scan(&exists)
	if err != nil || !exists {
		http.Error(w, "Bài viết không tồn tại", http.StatusNotFound)
		return
	}

	_, err = tx.Exec(`
		INSERT INTO cam_xuc (idBaiDang, idTaiKhoan, loaiCamXuc)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE loaiCamXuc = VALUES(loaiCamXuc)
	`, idBaiDang, idTaiKhoan, request.LoaiCamXuc)

	if err != nil {
		log.Println("Lỗi DB:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Println("Lỗi commit:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// ✅ Xoá cache cũ
	cacheKey := fmt.Sprintf("reactions:%d", idBaiDang)
	redisClient.Del(ctx, cacheKey)

	// ✅ Trả về JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"action":  "updated",
		"emotion": request.LoaiCamXuc,
	})
}


func ThongKeCamXuc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	idBaiDang, err := strconv.Atoi(params["idBaiDang"])
	if err != nil {
		http.Error(w, "ID bài viết không hợp lệ", http.StatusBadRequest)
		return
	}

	cacheKey := fmt.Sprintf("reactions:%d", idBaiDang)
	if cached, err := redisClient.Get(ctx, cacheKey).Result(); err == nil {
		w.Write([]byte(cached))
		return
	}

	rows, err := DB.Query(`
		SELECT loaiCamXuc, COUNT(*) as count 
		FROM cam_xuc 
		WHERE idBaiDang = ?
		GROUP BY loaiCamXuc
	`, idBaiDang)
	if err != nil {
		log.Println("Lỗi truy vấn:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	result := make(map[string]int)
	for rows.Next() {
		var loai string
		var count int
		if err := rows.Scan(&loai, &count); err != nil {
			log.Println("Lỗi đọc dữ liệu:", err)
			continue
		}
		result[loai] = count
	}

	// ✅ Lưu cache trong 30 phút
	jsonData, _ := json.Marshal(result)
	redisClient.Set(ctx, cacheKey, string(jsonData), 30*time.Minute)

	json.NewEncoder(w).Encode(result)
}

func SuaBaiDang(w http.ResponseWriter, r *http.Request) {
 vars := mux.Vars(r)
    idStr := vars["id"]
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "ID không hợp lệ", http.StatusBadRequest)
        return
    }

    var data struct {
        TieuDe  string `json:"tieuDe"`
        NoiDung string `json:"noiDung"`
    }

    if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
        http.Error(w, "Dữ liệu không hợp lệ", http.StatusBadRequest)
        return
    }

    _, err = DB.Exec(`UPDATE baidang SET tieuDe = ?, noiDung = ? WHERE id = ?`, data.TieuDe, data.NoiDung, id)
    if err != nil {
        http.Error(w, "Không thể cập nhật bài viết", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(map[string]string{
        "message": "Cập nhật bài viết thành công!",
    })
}

func XoaBinhLuan(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    idStr := vars["id"]
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "ID không hợp lệ", http.StatusBadRequest)
        return
    }

    _, err = DB.Exec(`DELETE FROM binhluan WHERE id = ?`, id)
    if err != nil {
        http.Error(w, "Không thể xóa bình luận", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(map[string]string{
        "message": "Đã xóa bình luận thành công!",
    })
}

