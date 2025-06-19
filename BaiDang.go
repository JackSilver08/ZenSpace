package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type BaiDang struct {
	TieuDe     string `json:"tieuDe"`
	NoiDung    string `json:"noiDung"`
	IDTaiKhoan int    `json:"idTaiKhoan"`
}



func DangBai(w http.ResponseWriter, r *http.Request) {
	// ======= BỔ SUNG CORS CHO ROUTE NÀY =========
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	// ===========================================

	if r.Method != http.MethodPost {
		http.Error(w, "Chỉ hỗ trợ POST", http.StatusMethodNotAllowed)
		return
	}

	var bai BaiDang
	err := json.NewDecoder(r.Body).Decode(&bai)
	if err != nil {
		log.Println("Lỗi JSON:", err)
		http.Error(w, "Lỗi dữ liệu gửi lên", http.StatusBadRequest)
		return
	}

	stmt := `INSERT INTO baidang (tieuDe, noiDung, idTaiKhoan, ngayDang) VALUES (?, ?, ?, ?)`
	_, err = DB.Exec(stmt, bai.TieuDe, bai.NoiDung, bai.IDTaiKhoan, time.Now())
	if err != nil {
		log.Println("Lỗi ghi DB:", err)
		http.Error(w, "Lỗi khi lưu bài viết", http.StatusInternalServerError)
		return
	}

	log.Println("✅ Đăng bài thành công:", bai.TieuDe)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Đăng bài thành công"))
}

func LayBaiDang(w http.ResponseWriter, r *http.Request) {
	// Bổ sung CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "Chỉ hỗ trợ GET", http.StatusMethodNotAllowed)
		return
	}

	rows, err := DB.Query(`
	SELECT b.id, b.tieuDe, b.noiDung, b.idTaiKhoan, t.tenDangNhap, b.ngayDang
	FROM baidang b
	JOIN taikhoan t ON b.idTaiKhoan = t.id
	ORDER BY b.ngayDang DESC`)

	if err != nil {
		log.Println("Lỗi truy vấn:", err)
		http.Error(w, "Không thể lấy bài viết", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var ds []map[string]interface{}
	for rows.Next() {
		var (
	id         int
	tieuDe     string
	noiDung    string
	idTaiKhoan int
	tenDangNhap string
	ngayDang   time.Time
)

if err := rows.Scan(&id, &tieuDe, &noiDung, &idTaiKhoan, &tenDangNhap, &ngayDang); err != nil {
	log.Println("Lỗi đọc dữ liệu:", err)
	continue
}


	ds = append(ds, map[string]interface{}{
			"id":         id,
			"tenDangNhap": tenDangNhap,
			"tieuDe":     tieuDe,
			"noiDung":    noiDung,
			"idTaiKhoan": idTaiKhoan,
			"ngayDang":   ngayDang.Format("02/01/2006 15:04"),
		})
	}

	json.NewEncoder(w).Encode(ds)
}

// Lấy chi tiết 1 bài đăng theo ID
func LayChiTietBaiDang(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	idStr := params["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID không hợp lệ", http.StatusBadRequest)
		return
	}

	var (
		tieuDe     string
		noiDung    string
		idTaiKhoan int
		ngayDang   time.Time
	)

	err = DB.QueryRow("SELECT tieuDe, noiDung, idTaiKhoan, ngayDang FROM baidang WHERE id = ?", id).Scan(&tieuDe, &noiDung, &idTaiKhoan, &ngayDang)
	if err != nil {
		log.Println("Không tìm thấy bài viết:", err)
		http.Error(w, "Không tìm thấy bài viết", http.StatusNotFound)
		return
	}

	chiTiet := map[string]interface{}{
		"id":         id,
		"tieuDe":     tieuDe,
		"noiDung":    noiDung,
		"idTaiKhoan": idTaiKhoan,
		"ngayDang":   ngayDang.Format("02/01/2006 15:04"),
	}

	json.NewEncoder(w).Encode(chiTiet)
}
func XoaBaiDang(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodDelete {
		http.Error(w, "Chỉ hỗ trợ DELETE", http.StatusMethodNotAllowed)
		return
	}

	params := mux.Vars(r)
	idStr := params["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID không hợp lệ", http.StatusBadRequest)
		return
	}

	stmt := `DELETE FROM baidang WHERE id = ?`
	result, err := DB.Exec(stmt, id)
	if err != nil {
		log.Println("Lỗi khi xóa bài viết:", err)
		http.Error(w, "Không thể xóa bài viết", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Không tìm thấy bài viết để xóa", http.StatusNotFound)
		return
	}

	log.Printf("✅ Đã xóa bài viết ID %d\n", id)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Đã xóa bài viết thành công"))
}
