package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// Struct cho người dùng
type NguoiDung struct {
	ID        int    `json:"id"`
	HoTen     string `json:"tenDangNhap"`
	Email     string `json:"email"`
	Avatar    string `json:"avatar"`
	PhanQuyen string `json:"phanQuyen"`
}

// Struct cho bài viết
type BaiViet struct {
	ID       int    `json:"id"`
	TieuDe   string `json:"tieuDe"`
	NoiDung  string `json:"noiDung"`
	NgayDang string `json:"ngayDang"` // ISO8601, ví dụ: 2023-06-21T15:04:05Z
}

// Hàm tiện ích set header CORS chuẩn
func setCorsHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*") // bạn có thể thay bằng domain cụ thể
	w.Header().Set("Access-Control-Allow-Methods", "GET, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

// Hàm tiện ích trả JSON chuẩn
func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		log.Println("Lỗi encode JSON:", err)
	}
}
func TrangCaNhanHandler(w http.ResponseWriter, r *http.Request) {
    setCorsHeaders(w)

    if r.Method == http.MethodOptions {
        w.WriteHeader(http.StatusOK)
        return
    }
    if r.Method != http.MethodGet {
        writeJSON(w, http.StatusMethodNotAllowed, map[string]interface{}{
            "success": false,
            "message": "Chỉ chấp nhận GET",
        })
        return
    }

    // Lấy userID từ URL
    userIDStr := mux.Vars(r)["id"]
    userID, err := strconv.Atoi(userIDStr)
    if err != nil {
        writeJSON(w, http.StatusBadRequest, map[string]interface{}{
            "success": false,
            "message": "ID không hợp lệ",
        })
        return
    }

    // Truy vấn thông tin người dùng
    var nguoiDung NguoiDung
    var avatar sql.NullString
    err = DB.QueryRow(`
        SELECT id, tenDangNhap, email, avatar, phanQuyen
        FROM taikhoan
        WHERE id = ?`, userID).
        Scan(&nguoiDung.ID, &nguoiDung.HoTen, &nguoiDung.Email, &avatar, &nguoiDung.PhanQuyen)
    if err != nil {
        if err == sql.ErrNoRows {
            writeJSON(w, http.StatusNotFound, map[string]interface{}{
                "success": false,
                "message": "Không tìm thấy người dùng",
            })
        } else {
            log.Println("❌ Lỗi query user:", err)
            writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
                "success": false,
                "message": "Lỗi khi truy vấn dữ liệu người dùng",
            })
        }
        return
    }
    nguoiDung.Avatar = "IMG/ZenUser.png"
    if avatar.Valid && avatar.String != "" {
        nguoiDung.Avatar = avatar.String
    }

    // Truy vấn các bài viết
    rows, err := DB.Query(`
        SELECT id, tieuDe, noiDung, ngayDang
        FROM baidang
        WHERE idTaiKhoan = ?
        ORDER BY ngayDang DESC`, userID)
    if err != nil {
        log.Println("❌ Lỗi query bài viết:", err)
        writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
            "success": false,
            "message": "Lỗi khi truy vấn bài viết",
        })
        return
    }
    defer rows.Close()

    var baiVietList []BaiViet
    for rows.Next() {
        var b BaiViet
        var ngay time.Time
        if err := rows.Scan(&b.ID, &b.TieuDe, &b.NoiDung, &ngay); err != nil {
            log.Println("❌ Lỗi scan bài viết:", err)
            continue
        }
        b.NgayDang = ngay.Format(time.RFC3339)
        baiVietList = append(baiVietList, b)
    }

    if err := rows.Err(); err != nil {
        log.Println("❌ Lỗi rows bài viết:", err)
    }

    // Trả về JSON
    writeJSON(w, http.StatusOK, map[string]interface{}{
        "success": true,
        "message": "Lấy trang cá nhân thành công",
        "data": map[string]interface{}{
            "user":     nguoiDung,
            "baiViets": baiVietList,
        },
    })
}


func XoaTaiKhoanHandler(w http.ResponseWriter, r *http.Request) {
    setCorsHeaders(w)

    if r.Method == http.MethodOptions {
        w.WriteHeader(http.StatusOK)
        return
    }
    if r.Method != http.MethodDelete {
        writeJSON(w, http.StatusMethodNotAllowed, map[string]interface{}{
            "success": false,
            "message": "Chỉ chấp nhận DELETE",
        })
        return
    }

    // Lấy ID từ URL
    userIDStr := mux.Vars(r)["id"]
    userID, err := strconv.Atoi(userIDStr)
    if err != nil {
        writeJSON(w, http.StatusBadRequest, map[string]interface{}{
            "success": false,
            "message": "ID không hợp lệ",
        })
        return
    }

    // Kiểm tra tồn tại
    var exists bool
    if err := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM taikhoan WHERE id = ?)", userID).Scan(&exists); err != nil || !exists {
        log.Printf("❌ Tài khoản %d không tồn tại hoặc lỗi truy vấn: %v\n", userID, err)
        writeJSON(w, http.StatusNotFound, map[string]interface{}{
            "success": false,
            "message": "Không tìm thấy tài khoản",
        })
        return
    }

    // Xoá dữ liệu phụ thuộc theo thứ tự khoá ngoại
    relatedTables := []string{"cam_xuc", "binhluan", "baidang"}
    for _, tbl := range relatedTables {
        if _, err := DB.Exec("DELETE FROM "+tbl+" WHERE idTaiKhoan = ?", userID); err != nil {
            log.Printf("❌ Không thể xoá bảng %s của user %d: %v\n", tbl, userID, err)
            writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
                "success": false,
                "message": "Không thể xoá dữ liệu " + tbl,
            })
            return
        }
    }
// Xoá tin nhắn mà user là người gửi hoặc người nhận
if _, err := DB.Exec(`
  DELETE FROM tin_nhan 
  WHERE nguoi_gui_id = ? OR nguoi_nhan_id = ?`, userID, userID); err != nil {
    log.Printf("❌ Không thể xoá tin nhắn của user %d: %v\n", userID, err)
    writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
        "success": false,
        "message": "Không thể xoá tin nhắn của người dùng",
    })
    return
}

    // Sau cùng: xoá tài khoản
    if _, err := DB.Exec("DELETE FROM taikhoan WHERE id = ?", userID); err != nil {
        log.Printf("❌ Không thể xoá tài khoản %d: %v\n", userID, err)
        writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
            "success": false,
            "message": "Không thể xoá tài khoản",
        })
        return
    }

    writeJSON(w, http.StatusOK, map[string]interface{}{
        "success": true,
        "message": "Tài khoản và dữ liệu liên quan đã được xoá.",
    })
}


// PUT /api/nguoidung/{id}/avatar
func DoiAvatarHandler(w http.ResponseWriter, r *http.Request) {
	setCorsHeaders(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method != http.MethodPut {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]interface{}{
			"success": false,
			"message": "Chỉ chấp nhận PUT",
		})
		return
	}

	idStr := mux.Vars(r)["id"]
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "ID không hợp lệ",
		})
		return
	}

	var body struct {
		Avatar string `json:"avatar"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Dữ liệu không hợp lệ",
		})
		return
	}

	_, err = DB.Exec(`UPDATE taikhoan SET avatar = ? WHERE id = ?`, body.Avatar, userID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Không thể cập nhật avatar",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Đã cập nhật avatar",
	})
}

func TimKiemNguoiDung(w http.ResponseWriter, r *http.Request) {
	setCorsHeaders(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	key := r.URL.Query().Get("key")
	if key == "" {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Thiếu từ khóa",
		})
		return
	}

	rows, err := DB.Query(`SELECT id, tenDangNhap FROM taikhoan WHERE tenDangNhap LIKE ? LIMIT 10`, "%"+key+"%")
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Lỗi truy vấn DB",
		})
		return
	}
	defer rows.Close()

	var users []map[string]interface{}
	for rows.Next() {
		var id int
		var ten string
		rows.Scan(&id, &ten)
		users = append(users, map[string]interface{}{
			"id":          id,
			"tenDangNhap": ten,
		})
	}

	writeJSON(w, http.StatusOK, users)
}

func GuiTinNhanHandler(w http.ResponseWriter, r *http.Request) {
    setCorsHeaders(w)

    if r.Method == http.MethodOptions {
        w.WriteHeader(http.StatusOK)
        return
    }
    if r.Method != http.MethodPost {
        writeJSON(w, http.StatusMethodNotAllowed, map[string]interface{}{
            "success": false,
            "message": "Chỉ chấp nhận POST",
        })
        return
    }

    var data struct {
        NguoiGuiID   int    `json:"nguoiGuiID"`
        NguoiNhanID  int    `json:"nguoiNhanID"`
        NoiDung      string `json:"noiDung"`
    }

    if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
        writeJSON(w, http.StatusBadRequest, map[string]interface{}{
            "success": false,
            "message": "Dữ liệu gửi lên không hợp lệ",
        })
        return
    }

    if data.NoiDung == "" || data.NguoiGuiID == 0 || data.NguoiNhanID == 0 {
        writeJSON(w, http.StatusBadRequest, map[string]interface{}{
            "success": false,
            "message": "Thiếu thông tin cần thiết",
        })
        return
    }

    // Lưu vào DB
    _, err := DB.Exec(`
        INSERT INTO tin_nhan (nguoi_gui_id, nguoi_nhan_id, noi_dung)
        VALUES (?, ?, ?)`, data.NguoiGuiID, data.NguoiNhanID, data.NoiDung)

    if err != nil {
        log.Println("❌ Lỗi lưu tin nhắn:", err)
        writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
            "success": false,
            "message": "Lỗi khi lưu tin nhắn",
        })
        return
    }

    writeJSON(w, http.StatusOK, map[string]interface{}{
        "success": true,
        "message": "Gửi tin nhắn thành công",
    })
}

func LayLichSuTinNhanHandler(w http.ResponseWriter, r *http.Request) {
  setCorsHeaders(w)

  if r.Method == http.MethodOptions {
    w.WriteHeader(http.StatusOK)
    return
  }

  if r.Method != http.MethodGet {
    writeJSON(w, http.StatusMethodNotAllowed, map[string]interface{}{
      "success": false,
      "message": "Chỉ chấp nhận GET",
    })
    return
  }

  // 👉 Lấy id1 và id2 từ URL
  parts := strings.Split(r.URL.Path, "/")
  if len(parts) < 6 {
    writeJSON(w, http.StatusBadRequest, map[string]interface{}{
      "success": false,
      "message": "Thiếu ID người dùng trong URL",
    })
    return
  }

  id1, err1 := strconv.Atoi(parts[4])
  id2, err2 := strconv.Atoi(parts[5])
  if err1 != nil || err2 != nil {
    writeJSON(w, http.StatusBadRequest, map[string]interface{}{
      "success": false,
      "message": "ID không hợp lệ",
    })
    return
  }

  // 👉 Truy vấn 2 chiều tin nhắn giữa id1 và id2
  rows, err := DB.Query(`
    SELECT id, nguoi_gui_id, nguoi_nhan_id, noi_dung, thoi_gian
    FROM tin_nhan
    WHERE (nguoi_gui_id = ? AND nguoi_nhan_id = ?)
       OR (nguoi_gui_id = ? AND nguoi_nhan_id = ?)
    ORDER BY thoi_gian ASC
  `, id1, id2, id2, id1)

  if err != nil {
    log.Println("❌ Lỗi khi lấy tin nhắn:", err)
    writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
      "success": false,
      "message": "Lỗi khi truy vấn tin nhắn",
    })
    return
  }
  defer rows.Close()

  var tinNhan []map[string]interface{}
  for rows.Next() {
    var id, gui, nhan int
    var noiDung string
    var thoiGian time.Time

    if err := rows.Scan(&id, &gui, &nhan, &noiDung, &thoiGian); err == nil {
      tinNhan = append(tinNhan, map[string]interface{}{
        "id":          id,
        "nguoiGuiID":  gui,
        "nguoiNhanID": nhan,
        "noiDung":     noiDung,
        "thoiGian":    thoiGian.Format(time.RFC3339),
      })
    }
  }

  writeJSON(w, http.StatusOK, map[string]interface{}{
    "success": true,
    "tinNhan": tinNhan,
  })
}
