package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

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
	 var err error

    DB, err = sql.Open("mysql", "root:Tuan@1234@tcp(127.0.0.1:3306)/ZenSpaceDB?parseTime=true")
    if err != nil {
        log.Fatal("❌ Không thể mở kết nối MySQL:", err)
    }

    err = DB.Ping()
    if err != nil {
        log.Fatal("❌ Không thể ping tới MySQL:", err)
    }

    log.Println("✅ Kết nối MySQL thành công!")
	
    router := mux.NewRouter()
	router.Use(enableCORS)
    router.HandleFunc("/hello", Hello).Methods("GET", "OPTIONS")
    router.HandleFunc("/", handler)
    router.HandleFunc("/DangNhap", DangNhap)
    router.HandleFunc("/DangKy", DangKy)
    router.HandleFunc("/DangBai", DangBai)
    router.HandleFunc("/LayBaiDang", LayBaiDang)
    router.HandleFunc("/api/baiviet/{id}", LayChiTietBaiDang).Methods("GET", "OPTIONS")
	router.HandleFunc("/XoaBaiDang/{id}", XoaBaiDang).Methods("DELETE", "OPTIONS") // thay thế LayBaiDang nếu tên handler khác
    router.HandleFunc("/api/binhluan", ThemBinhLuan).Methods("POST", "OPTIONS")
    router.HandleFunc("/api/binhluan/{id}", LayBinhLuanTheoBaiDang).Methods("GET", "OPTIONS")
    log.Println("🚀 Server chạy tại cổng 8080...")
    err = http.ListenAndServe(":8080", router)
    if err != nil {
        log.Fatal("❌ Không thể chạy server:", err)
    }
}
