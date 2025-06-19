package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

type BinhLuan struct {
	ID          int    `json:"id,omitempty"`
	IDBaiDang   int    `json:"idBaiDang"`
	IDTaiKhoan  int    `json:"idTaiKhoan"`
	NoiDung     string `json:"noiDung"`
	NgayBinhLuan string `json:"ngayBinhLuan,omitempty"`
	TenDangNhap string `json:"tenDangNhap,omitempty"`
}

// Lấy danh sách bình luận theo id bài viết
func LayBinhLuanTheoBaiDang(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idBaiDang := vars["id"]

	rows, err := DB.Query(`
		SELECT b.noiDung, t.tenDangNhap, b.ngayBinhLuan
		FROM binhluan b
		JOIN taikhoan t ON b.idTaiKhoan = t.id
		WHERE b.idBaiDang = ?
		ORDER BY b.ngayBinhLuan ASC
	`, idBaiDang)
	if err != nil {
		http.Error(w, "Lỗi truy vấn bình luận", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var binhLuans []BinhLuan
	for rows.Next() {
		var bl BinhLuan
		if err := rows.Scan(&bl.NoiDung, &bl.TenDangNhap, &bl.NgayBinhLuan); err != nil {
			http.Error(w, "Lỗi đọc dữ liệu", http.StatusInternalServerError)
			return
		}
		binhLuans = append(binhLuans, bl)
	}
if binhLuans == nil {
	binhLuans = []BinhLuan{}
}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(binhLuans)
}

// Thêm bình luận mới
func ThemBinhLuan(w http.ResponseWriter, r *http.Request) {
    var bl BinhLuan

    bodyBytes, _ := io.ReadAll(r.Body)
    fmt.Println("📦 Raw body nhận được từ frontend:", string(bodyBytes))
    r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

    if err := json.NewDecoder(r.Body).Decode(&bl); err != nil {
        http.Error(w, "Dữ liệu gửi lên không hợp lệ: "+err.Error(), http.StatusBadRequest)
        return
    }

    fmt.Printf("✅ Decode JSON thành công: %+v\n", bl)

    if bl.IDBaiDang == 0 || bl.IDTaiKhoan == 0 || bl.NoiDung == "" {
        http.Error(w, fmt.Sprintf("Thiếu thông tin: idBaiDang=%d, idTaiKhoan=%d, noiDung='%s'",
            bl.IDBaiDang, bl.IDTaiKhoan, bl.NoiDung), http.StatusBadRequest)
        return
    }

    _, err := DB.Exec(`
        INSERT INTO binhluan (idBaiDang, idTaiKhoan, noiDung)
        VALUES (?, ?, ?)
    `, bl.IDBaiDang, bl.IDTaiKhoan, bl.NoiDung)
    if err != nil {
        http.Error(w, "Lỗi thêm bình luận vào database", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"message": "Bình luận thành công"})
}

