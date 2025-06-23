// ==========================
// Zen Space - Main Script
// ==========================

// ==== Helper ==== //
function fetchJSON(url) {
  return fetch(url).then((res) => {
    if (!res.ok) {
      throw new Error(`HTTP error! Status: ${res.status}`);
    }
    return res.json();
  });
}

// ==== Hiển thị danh sách bài viết ==== //
function hienThiDanhSachBaiDang(dsBaiDang) {
  const danhSach = document.getElementById("danhSachBaiDang");
  danhSach.innerHTML = "";

  if (!Array.isArray(dsBaiDang) || dsBaiDang.length === 0) {
    danhSach.innerHTML = `
      <div style="padding: 20px; color: #999; text-align: center;">
        Không tìm thấy bài viết phù hợp.
      </div>`;
    danhSach.style.justifyContent = "center";
    return;
  }

  dsBaiDang.forEach((bai) => {
    const container = document.createElement("div");
    container.className = "bai-dang";
    container.style.marginBottom = "20px";
    container.style.padding = "15px";
    container.style.border = "1px solid #eee";
    container.style.borderRadius = "8px";
    container.style.backgroundColor = "#fff";
    container.style.boxShadow = "0 2px 6px rgba(0,0,0,0.05)";

    const title = document.createElement("h3");
    title.textContent = bai.tieuDe;
    title.style.cursor = "pointer";
    title.style.marginBottom = "8px";
    title.onclick = () => {
      window.location.href = `BaiViet.html?id=${bai.id}`;
    };

    const meta = document.createElement("div");
    meta.textContent = `Đăng bởi: ${bai.tenDangNhap || "Ẩn danh"} - ${
      bai.ngayDang || ""
    }`;
    meta.style.color = "#888";
    meta.style.fontSize = "13px";
    meta.style.marginBottom = "10px";

    const contentDiv = document.createElement("div");
    contentDiv.innerHTML = bai.noiDung;
    contentDiv.className = "noi-dung";

    container.appendChild(title);
    container.appendChild(meta);
    container.appendChild(contentDiv);

    // Thêm nút xóa nếu bạn muốn (tuỳ điều kiện)
    const userId = parseInt(localStorage.getItem("userId"));
    if (bai.idTaiKhoan === userId || localStorage.getItem("role") === "admin") {
      const deleteBtn = document.createElement("button");
      deleteBtn.textContent = "🗑 Xóa";
      deleteBtn.style.marginTop = "10px";
      deleteBtn.style.backgroundColor = "#ff4d4f";
      deleteBtn.style.color = "#fff";
      deleteBtn.style.border = "none";
      deleteBtn.style.padding = "6px 12px";
      deleteBtn.style.borderRadius = "4px";
      deleteBtn.style.cursor = "pointer";
      deleteBtn.onclick = () => xoaBaiViet(bai.id);
      container.appendChild(deleteBtn);
    }

    danhSach.appendChild(container);
  });
}

// ==== Thông báo bài viết mới nâng cao ==== //
document.addEventListener("DOMContentLoaded", () => {
  const thongBaoBtn = document.getElementById("thongBaoBtn");
  const soThongBao = document.getElementById("soThongBao");

  // Tạo container cho thông báo
  const thongBaoContainer = document.createElement("div");
  thongBaoContainer.id = "thongBaoPopup";
  thongBaoContainer.style.position = "fixed";
  thongBaoContainer.style.top = "60px";
  thongBaoContainer.style.right = "20px";
  thongBaoContainer.style.backgroundColor = "#fff";
  thongBaoContainer.style.border = "1px solid #eee";
  thongBaoContainer.style.borderRadius = "8px";
  thongBaoContainer.style.boxShadow = "0 2px 10px rgba(0,0,0,0.1)";
  thongBaoContainer.style.padding = "15px";
  thongBaoContainer.style.maxWidth = "350px";
  thongBaoContainer.style.zIndex = "1000";
  thongBaoContainer.style.display = "none";
  thongBaoContainer.style.transition = "all 0.3s ease";
  document.body.appendChild(thongBaoContainer);

  // Biến lưu trạng thái thông báo
  let isShowingNotification = false;
  let newPosts = [];

  // Hàm kiểm tra bài mới
  async function kiemTraBaiMoi() {
    try {
      const data = await fetchJSON("http://localhost:8080/LayBaiDang");
      if (!Array.isArray(data) || data.length === 0) return;

      const latestPost = data[0];
      const lastSeenId = parseInt(localStorage.getItem("lastSeenPostId")) || 0;

      if (latestPost.id > lastSeenId) {
        // Lấy danh sách bài viết mới
        newPosts = data.filter((post) => post.id > lastSeenId);
        const newPostsCount = newPosts.length;

        // Hiển thị số thông báo
        soThongBao.textContent = `(${newPostsCount})`;

        // Chỉ hiển thị popup nếu không đang hiển thị
        if (!isShowingNotification) {
          hienThiThongBaoPopup();
        }
      }
    } catch (err) {
      console.error("Lỗi khi kiểm tra bài mới:", err.message);
    }
  }

  // Hiển thị popup thông báo
  function hienThiThongBaoPopup() {
    if (newPosts.length === 0) return;

    isShowingNotification = true;

    // Hiển thị bài mới nhất đầu tiên
    const latestPost = newPosts[0];

    thongBaoContainer.innerHTML = `
      <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 10px;">
        <h3 style="margin: 0; color: #d9534f;">Bài viết mới</h3>
        <span style="font-size: 12px; color: #888;">${new Date().toLocaleTimeString()}</span>
      </div>
      <div style="border-bottom: 1px solid #eee; padding-bottom: 10px; margin-bottom: 10px;">
        <h4 style="margin: 0 0 5px 0; cursor: pointer;" onclick="window.location.href='BaiViet.html?id=${
          latestPost.id
        }'">
          ${latestPost.tieuDe || "Bài viết mới"}
        </h4>
        <p style="margin: 5px 0; color: #666; font-size: 14px;">
          ${
            latestPost.noiDung
              ? latestPost.noiDung.substring(0, 80) +
                (latestPost.noiDung.length > 80 ? "..." : "")
              : ""
          }
        </p>
        <div style="display: flex; justify-content: space-between; align-items: center;">
          <small style="color: #888;">${
            latestPost.tenDangNhap || "Ẩn danh"
          } • ${new Date(latestPost.ngayDang).toLocaleString() || ""}</small>
          <span style="font-size: 12px; background: #f0f0f0; padding: 2px 5px; border-radius: 3px;">
            ${newPosts.length > 1 ? `+${newPosts.length - 1} bài khác` : ""}
          </span>
        </div>
      </div>
      <div style="display: flex; justify-content: space-between;">
        <button id="xemTatCaBtn" style="background: #5bc0de; color: white; border: none; padding: 5px 10px; border-radius: 3px; cursor: pointer; font-size: 13px;">
          Xem tất cả (${newPosts.length})
        </button>
        <button id="dongThongBaoBtn" style="background: none; border: none; color: #888; cursor: pointer; font-size: 13px;">
          Đóng
        </button>
      </div>
    `;

    thongBaoContainer.style.display = "block";

    // Xử lý sự kiện cho nút trong popup
    document.getElementById("xemTatCaBtn")?.addEventListener("click", () => {
      localStorage.setItem("lastSeenPostId", newPosts[0].id);
      soThongBao.textContent = "";
      thongBaoContainer.style.display = "none";
      isShowingNotification = false;
      window.location.href = "index.html#new-posts"; // Có thể thêm hash để highlight bài mới
    });

    document
      .getElementById("dongThongBaoBtn")
      ?.addEventListener("click", () => {
        thongBaoContainer.style.display = "none";
        isShowingNotification = false;
      });

    // Tự động ẩn sau 8 giây
    setTimeout(() => {
      if (isShowingNotification) {
        thongBaoContainer.style.display = "none";
        isShowingNotification = false;
      }
    }, 8000);
  }

  // Xử lý sự kiện click nút thông báo
  if (thongBaoBtn) {
    thongBaoBtn.addEventListener("click", async () => {
      const data = await fetchJSON("http://localhost:8080/LayBaiDang");
      if (data && data.length > 0) {
        localStorage.setItem("lastSeenPostId", data[0].id);
        soThongBao.textContent = "";
        thongBaoContainer.style.display = "none";
        isShowingNotification = false;

        // Highlight các bài viết mới trên trang
        const newPosts = data.filter(
          (post) =>
            post.id > (parseInt(localStorage.getItem("lastSeenPostId")) || 0)
        );
        if (newPosts.length > 0) {
          hienThiDanhSachBaiDang(data); // Sử dụng hàm hiển thị có sẵn
          window.scrollTo(0, 0); // Cuộn lên đầu trang
        }
      }
    });
  }

  // Kiểm tra bài mới mỗi 15 giây
  setInterval(kiemTraBaiMoi, 15000);
  kiemTraBaiMoi();
});
// ==== Giao diện người dùng ==== //
document.addEventListener("DOMContentLoaded", () => {
  const username = localStorage.getItem("username");
  const avatarUrl = "IMG/ZenUser.png"; // hoặc "./IMG/ZenUser.png"

  if (username) {
    const loginLink = document.getElementById("loginLink");
    if (loginLink) loginLink.style.display = "none";

    const userAvatarDiv = document.getElementById("userAvatar");
    if (userAvatarDiv) userAvatarDiv.style.display = "block";

    const userAvatarImg = document.getElementById("userAvatarImg");
    if (userAvatarImg) {
      if (avatarUrl && avatarUrl.trim() !== "") {
        const img = new Image();
        img.onload = () => {
          userAvatarImg.src = avatarUrl;
        };
        img.onerror = () => {
          console.warn("Ảnh avatar không tải được, dùng ảnh mặc định.");
          userAvatarImg.src = "IMG/ZenUser.png"; // ✅ đúng với cấu trúc thư mục
        };
        console.log("Avatar URL đang được set là:", avatarUrl);
        img.src = avatarUrl;
      } else {
        userAvatarImg.src = "IMG/ZenUser.png"; // ✅ đúng với cấu trúc thư mục
        // fallback mặc định
      }
    }

    const welcomeText = document.getElementById("welcomeText");
    if (welcomeText) welcomeText.textContent = `Xin chào, ${username}`;

    const logoutBtn = document.getElementById("logoutBtn");
    if (logoutBtn) {
      logoutBtn.addEventListener("click", () => {
        localStorage.removeItem("username");
        localStorage.removeItem("avatarUrl");
        location.reload();
      });
    }
  }
});

// ==== Đăng bài Popup ==== //
document.addEventListener("DOMContentLoaded", () => {
  const openBtn = document.getElementById("openPopupBtn");
  const popupOverlay = document.getElementById("popupOverlay");
  const closeBtn = document.getElementById("closePopupBtn");

  if (openBtn && popupOverlay && closeBtn) {
    openBtn.addEventListener("click", () =>
      popupOverlay.classList.add("active")
    );
    closeBtn.addEventListener("click", () =>
      popupOverlay.classList.remove("active")
    );
    popupOverlay.addEventListener("click", (e) => {
      if (e.target === popupOverlay) popupOverlay.classList.remove("active");
    });
  }
});
// Tìm kiếm user và post
document.addEventListener("DOMContentLoaded", () => {
  const searchBox = document.getElementById("search-box");
  const danhSach = document.getElementById("danhSachBaiDang");

  // Thêm khối để hiển thị user (ở trên danh sách bài viết)
  const userResultDiv = document.createElement("div");
  userResultDiv.id = "ketQuaUser";
  userResultDiv.style.marginBottom = "20px";
  danhSach.parentNode.insertBefore(userResultDiv, danhSach);

  let danhSachBaiViet = [];

  // Hiển thị bài viết
  function hienThiBaiViet(dsBaiDang) {
    danhSach.innerHTML = "";

    if (!Array.isArray(dsBaiDang) || dsBaiDang.length === 0) {
      danhSach.innerHTML = `
        <div style="padding: 20px; color: #999; text-align: center;">
          Không tìm thấy bài viết phù hợp.
        </div>`;
      return;
    }

    dsBaiDang.forEach((bai) => {
      const container = document.createElement("div");
      container.style.marginBottom = "15px";
      container.style.padding = "12px";
      container.style.borderBottom = "1px solid #eee";

      const title = document.createElement("h3");
      title.textContent = bai.tieuDe;
      title.style.cursor = "pointer";
      title.style.margin = "0";
      title.onclick = () => {
        window.location.href = `BaiViet.html?id=${bai.id}`;
      };

      const meta = document.createElement("div");
      meta.textContent = `Đăng bởi ${bai.tenDangNhap || "Ẩn danh"} - ${
        bai.ngayDang || ""
      }`;
      meta.style.color = "#888";
      meta.style.fontSize = "13px";

      container.appendChild(title);
      container.appendChild(meta);
      danhSach.appendChild(container);
    });
  }

  // Hiển thị người dùng
  function hienThiNguoiDung(users) {
    userResultDiv.innerHTML = "";

    if (!Array.isArray(users) || users.length === 0) return;

    const title = document.createElement("h4");
    title.textContent = "Người dùng khớp:";
    title.style.marginBottom = "10px";
    userResultDiv.appendChild(title);

    users.forEach((user) => {
      const userDiv = document.createElement("div");
      userDiv.style.padding = "8px";
      userDiv.style.cursor = "pointer";
      userDiv.style.border = "1px solid #ccc";
      userDiv.style.borderRadius = "6px";
      userDiv.style.marginBottom = "6px";
      userDiv.style.background = "#f9f9f9";
      userDiv.textContent = user.tenDangNhap;
      userDiv.onclick = () => {
        window.location.href = `TaiKhoan.html?id=${user.id}`;
      };
      userResultDiv.appendChild(userDiv);
    });
  }

  // Hàm fetch JSON tiện lợi
  function fetchJSON(url) {
    return fetch(url).then((res) => {
      if (!res.ok) throw new Error("Lỗi khi gọi API");
      return res.json();
    });
  }

  // Tải bài viết ban đầu
  fetchJSON("http://localhost:8080/LayBaiDang")
    .then((data) => {
      danhSachBaiViet = data;
      hienThiBaiViet(danhSachBaiViet);
    })
    .catch((err) => {
      console.error("Lỗi khi tải bài viết:", err);
    });

  // Xử lý tìm kiếm bài viết + người dùng
  if (searchBox) {
    searchBox.addEventListener("input", function () {
      const tuKhoa = this.value.trim().toLowerCase();

      // Nếu trống, hiển thị lại như ban đầu
      if (!tuKhoa) {
        hienThiBaiViet(danhSachBaiViet);
        userResultDiv.innerHTML = "";
        return;
      }

      // Tìm trong tiêu đề bài viết
      const ketQuaBai = danhSachBaiViet.filter((bai) =>
        bai.tieuDe.toLowerCase().includes(tuKhoa)
      );
      hienThiBaiViet(ketQuaBai);

      // Gọi API tìm người dùng
      fetchJSON(
        `http://localhost:8080/api/timkiemnguoidung?key=${encodeURIComponent(
          tuKhoa
        )}`
      )
        .then((users) => {
          hienThiNguoiDung(users);
        })
        .catch((err) => {
          console.error("Lỗi tìm người dùng:", err);
          userResultDiv.innerHTML = "";
        });
    });
  }
});

function xoaBaiViet(id) {
  if (!confirm("Bạn có chắc muốn xóa bài viết này?")) return;

  fetch(`http://localhost:8080/XoaBaiDang/${id}`, {
    method: "DELETE",
    headers: {
      Authorization: `Bearer ${localStorage.getItem("token")}`, // nếu bạn cần xác thực token
    },
  })
    .then((res) => {
      if (res.ok) {
        alert("Xóa thành công!");
        location.reload();
      } else {
        alert("Không thể xóa.");
      }
    })
    .catch((err) => {
      console.error("Lỗi khi xóa:", err);
      alert("Đã xảy ra lỗi khi xóa bài viết.");
    });
}
function hienThiTimKiemDonGian(dsBaiDang) {
  const danhSach = document.getElementById("danhSachBaiDang");
  danhSach.innerHTML = "";

  if (!Array.isArray(dsBaiDang) || dsBaiDang.length === 0) {
    danhSach.innerHTML = `
      <div style="padding: 20px; color: #999; text-align: center;">
        Không tìm thấy bài viết phù hợp.
      </div>`;
    return;
  }

  dsBaiDang.forEach((bai) => {
    const container = document.createElement("div");
    container.style.marginBottom = "15px";
    container.style.padding = "12px";
    container.style.borderBottom = "1px solid #eee";

    const title = document.createElement("h3");
    title.textContent = bai.tieuDe;
    title.style.cursor = "pointer";
    title.style.margin = "0";
    title.onclick = () => {
      window.location.href = `BaiViet.html?id=${bai.id}`;
    };

    const meta = document.createElement("div");
    meta.textContent = `Đăng bởi ${bai.tenDangNhap || "Ẩn danh"} - ${
      bai.ngayDang || ""
    }`;
    meta.style.color = "#888";
    meta.style.fontSize = "13px";

    container.appendChild(title);
    container.appendChild(meta);
    danhSach.appendChild(container);
  });
}

const avatarToggle = document.getElementById("userAvatar");
const dropdownContainer = document.getElementById("userAvatarDropdown");

avatarToggle.addEventListener("click", () => {
  dropdownContainer.classList.toggle("active");
});

// Đóng dropdown nếu click ra ngoài
document.addEventListener("click", function (event) {
  if (!dropdownContainer.contains(event.target)) {
    dropdownContainer.classList.remove("active");
  }
});
