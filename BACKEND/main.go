package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
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
    // ✅ Lấy thông tin DB từ biến môi trường, fallback cho local
    host := os.Getenv("DB_HOST")
    if host == "" {
        host = "127.0.0.1"
    }

    port := os.Getenv("DB_PORT")
    if port == "" {
        port = "3306"
    }

    user := os.Getenv("DB_USER")
    if user == "" {
        user = "root"
    }

    pass := os.Getenv("DB_PASS")
    if pass == "" {
        pass = "Tuan@1234"
    }

    dbName := os.Getenv("DB_NAME")
    if dbName == "" {
        dbName = "ZenSpaceDB"
    }

    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, pass, host, port, dbName)

    // ✅ Kết nối DB
    var err error
    DB, err = sql.Open("mysql", dsn)
    if err != nil {
        log.Fatal("❌ Không thể mở kết nối MySQL:", err)
    }

    if err = DB.Ping(); err != nil {
        log.Fatal("❌ Không thể ping tới MySQL:", err)
    }

    log.Println("✅ Kết nối MySQL thành công!")

    // ✅ Thiết lập router
    router := mux.NewRouter()
    router.Use(enableCORS)
    router.HandleFunc("/hello", Hello).Methods("GET", "OPTIONS")
    router.HandleFunc("/", handler)
    router.HandleFunc("/DangNhap", DangNhap)
    router.HandleFunc("/DangKy", DangKy)
    router.HandleFunc("/DangBai", DangBai)
    router.HandleFunc("/LayBaiDang", LayBaiDang)
    router.HandleFunc("/api/baiviet/{id}", LayChiTietBaiDang).Methods("GET", "OPTIONS")
    router.HandleFunc("/XoaBaiDang/{id}", XoaBaiDang).Methods("DELETE", "OPTIONS")
    router.HandleFunc("/api/binhluan", ThemBinhLuan).Methods("POST", "OPTIONS")
    router.HandleFunc("/api/binhluan/{id}", LayBinhLuanTheoBaiDang).Methods("GET", "OPTIONS")

    // ✅ Cấu hình cổng
    port = os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    log.Printf("🚀 Server chạy tại cổng %s...\n", port)

    // ✅ Khởi động server
    if err := http.ListenAndServe(":"+port, router); err != nil {
        log.Fatal("❌ Không thể chạy server:", err)
    }
}
