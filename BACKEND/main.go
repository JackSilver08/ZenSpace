package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)


var DB *sql.DB  // ✅ Chỉ giữ ở đây

func validateToken(r *http.Request) (*jwt.Token, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, errors.New("authorization header missing")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, errors.New("authorization header format must be Bearer {token}")
	}

	tokenString := parts[1]

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token, nil
}


func init() {
 
}

func enableCORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization") // ✅ thêm Authorization

        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusOK)
            return
        }

        next.ServeHTTP(w, r)
    })
}

func Hello(w http.ResponseWriter, r *http.Request) {
    // CORS headers
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

    // Nếu là preflight request (OPTIONS)
    if r.Method == http.MethodOptions {
        w.WriteHeader(http.StatusOK)
        return
    }

    // Nội dung trả về
    w.Header().Set("Content-Type", "text/plain")
    fmt.Fprint(w, "Server đang hoạt động tại ZenSpace!")
}


func handler(w http.ResponseWriter, r *http.Request) {
	// Cho phép gọi từ bất kỳ nơi đâu (chỉ dùng khi phát triển)
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Cho phép các phương thức bạn cần
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	// Cho phép header bạn cần (nếu có)
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// Nếu đây là preflight OPTIONS request thì trả về ngay
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	fmt.Fprintf(w, "Server đang chạy tại cổng 8080, %s", r.URL.Path[1:])
}





func main() {
 host := os.Getenv("DB_HOST")
if host == "" { host = "hopper.proxy.rlwy.net" }

port := os.Getenv("DB_PORT")
if port == "" { port = "53455" }

user := os.Getenv("DB_USER")
if user == "" { user = "root" }

pass := os.Getenv("DB_PASS")
if pass == "" { pass = "dieShqEVeInoVrIswYLVLkwkyzYBcCHo" }

dbName := os.Getenv("DB_NAME")
if dbName == "" { dbName = "railway" }



    
var err error

   dsn := "root:Tuan@1234@tcp(127.0.0.1:3306)/ZenSpaceDB?parseTime=true"



    DB, err = sql.Open("mysql", dsn)
    if err != nil {
        log.Fatal("❌ Không thể mở kết nối MySQL:", err)
    }

    if err = DB.Ping(); err != nil {
        log.Fatal("❌ Không thể ping tới MySQL:", err)
    }

    log.Println("✅ Kết nối MySQL thành công!")

    // ✅ Khởi tạo router
    router := mux.NewRouter()
    router.Use(enableCORS)

    // Các endpoint handler
    router.HandleFunc("/", handler)
router.HandleFunc("/hello", Hello).Methods("GET", "OPTIONS")
// Tài khoản
router.HandleFunc("/api/nguoidung/{id}", TrangCaNhanHandler)
router.HandleFunc("/api/xoataikhoan/{id}", XoaTaiKhoanHandler).Methods("DELETE", "OPTIONS")
router.HandleFunc("/api/nguoidung/{id}/avatar", DoiAvatarHandler).Methods("PUT", "OPTIONS")
router.HandleFunc("/api/timkiemnguoidung", TimKiemNguoiDung).Methods("GET", "OPTIONS")
router.HandleFunc("/api/chat/gui", GuiTinNhanHandler).Methods("POST", "OPTIONS")
router.HandleFunc("/api/chat/lichsu/{id1}/{id2}", LayLichSuTinNhanHandler).Methods("GET", "OPTIONS")

// Đăng nhập / Đăng ký
router.HandleFunc("/DangNhap", DangNhap).Methods("POST", "OPTIONS")
router.HandleFunc("/DangKy", DangKy).Methods("POST", "OPTIONS")

// Bài đăng
router.HandleFunc("/DangBai", DangBai).Methods("POST", "OPTIONS")
router.HandleFunc("/LayBaiDang", LayBaiDang).Methods("GET", "OPTIONS")
router.HandleFunc("/api/baiviet/{id}", LayChiTietBaiDang).Methods("GET", "OPTIONS")
router.HandleFunc("/SuaBaiDang/{id}", SuaBaiDang).Methods("PUT", "OPTIONS")
router.HandleFunc("/XoaBaiDang/{id}", XoaBaiDang).Methods("DELETE", "OPTIONS")

// Bình luận
router.HandleFunc("/api/binhluan", ThemBinhLuan).Methods("POST", "OPTIONS")
router.HandleFunc("/api/binhluan/{id}", LayBinhLuanTheoBaiDang).Methods("GET", "OPTIONS")
router.HandleFunc("/XoaBinhLuan/{id}", XoaBinhLuan).Methods("DELETE", "OPTIONS") // nếu có

// Cảm xúc
router.HandleFunc("/ThemCamXuc/{idBaiDang}", ThemCamXuc).Methods("POST", "OPTIONS")
router.HandleFunc("/ThongKeCamXuc/{idBaiDang}", ThongKeCamXuc).Methods("GET", "OPTIONS")



    // ✅ Cổng chạy server
    serverPort := os.Getenv("PORT")
    if serverPort == "" {
        serverPort = "8080"
    }
    log.Printf("🚀 Server chạy tại cổng %s...\n", serverPort)

    if err := http.ListenAndServe(":"+serverPort, router); err != nil {
        log.Fatal("❌ Không thể chạy server:", err)
    }
}