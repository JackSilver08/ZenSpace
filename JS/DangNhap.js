// Lấy các phần tử cần thiết từ DOM
const loginForm = document.getElementById("login-form");
const registerForm = document.getElementById("register-form");
const switchToRegister = document.getElementById("switch-to-register");
const switchToLogin = document.getElementById("switch-to-login");
const loginMessage = document.getElementById("login-message");

// Chuyển form đăng nhập -> đăng ký
switchToRegister.addEventListener("click", (e) => {
  e.preventDefault();
  loginForm.style.display = "none";
  registerForm.style.display = "block";
  loginMessage.textContent = "";
});

// Chuyển form đăng ký -> đăng nhập
switchToLogin.addEventListener("click", (e) => {
  e.preventDefault();
  registerForm.style.display = "none";
  loginForm.style.display = "block";
  loginMessage.textContent = "";
});

// Xử lý đăng nhập khi submit form
loginForm.addEventListener("submit", (e) => {
  e.preventDefault(); // Ngăn reload trang

  const username = document.getElementById("login-username").value.trim();
  const password = document.getElementById("login-password").value;

  if (!username || !password) {
    loginMessage.textContent = "Vui lòng nhập đầy đủ thông tin.";
    return;
  }

  dangNhap(username, password);
});

// Hàm gửi yêu cầu đăng nhập đến server
async function dangNhap(username, password) {
  try {
    const response = await fetch("http://localhost:8080/DangNhap", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ TenDangNhap: username, MatKhau: password }), // ✅ Viết hoa đúng
    });

    if (!response.ok) {
      throw new Error("Lỗi server: " + response.status);
    }

    const data = await response.json();
    console.log("🧪 Token trả về từ server:", data.token);

    if (data.success) {
      console.log("✅ Token nhận được:", data.token);
      const user = data.user || {};
      const avatarPath =
        user.avatar && user.avatar.trim() !== ""
          ? user.avatar
          : "IMG/ZenUSer.png";

      // ✅ THÊM DÒNG NÀY
      localStorage.setItem("token", data.token);

      localStorage.setItem("idTaiKhoan", user.id);
      localStorage.setItem("username", user.hoten || username);
      localStorage.setItem("avatarUrl", avatarPath);
      localStorage.setItem("role", user.phanQuyen); // ✅ THÊM DÒNG NÀY
      localStorage.setItem("isLoggedIn", "true");

      window.location.href = "index.html";
    } else {
      loginMessage.textContent =
        data.message || "Tên đăng nhập hoặc mật khẩu không đúng.";
    }
  } catch (error) {
    console.error("Lỗi khi kết nối:", error);
    loginMessage.textContent =
      "Không thể kết nối đến server. Vui lòng thử lại sau.";
  }
}

// Xử lý đăng ký khi submit form
registerForm.addEventListener("submit", async (e) => {
  e.preventDefault();

  const username = document.getElementById("register-username").value.trim();
  const email = document.getElementById("register-email").value.trim();
  const password = document.getElementById("register-password").value;
  const confirmPassword = document.getElementById(
    "register-confirm-password"
  ).value;

  if (!username || !email || !password || !confirmPassword) {
    loginMessage.textContent = "Vui lòng nhập đầy đủ thông tin đăng ký.";
    return;
  }

  if (password !== confirmPassword) {
    loginMessage.textContent = "Mật khẩu xác nhận không khớp.";
    return;
  }

  try {
    const response = await fetch("http://localhost:8080/DangKy", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        username: username,
        email: email,
        password: password,
        confirmPassword: confirmPassword,
      }),
    });

    if (!response.ok) {
      throw new Error("Lỗi server: " + response.status);
    }

    const data = await response.json();

    if (data.success) {
      loginMessage.style.color = "green";
      loginMessage.textContent =
        data.message || "Đăng ký thành công! Vui lòng đăng nhập.";
      // Tự động chuyển về form đăng nhập
      registerForm.style.display = "none";
      loginForm.style.display = "block";
    } else {
      loginMessage.style.color = "red";
      loginMessage.textContent = data.message || "Đăng ký thất bại.";
    }
  } catch (error) {
    console.error("Lỗi khi kết nối:", error);
    loginMessage.style.color = "red";
    loginMessage.textContent =
      "Không thể kết nối đến server. Vui lòng thử lại sau.";
  }
});
