document.addEventListener("DOMContentLoaded", function () {
  // ==== Chat popup ====
  const chatIcon = document.getElementById("chat-icon");
  const chatPopup = document.getElementById("chat-popup");
  const chatClose = document.getElementById("chat-close");
  const chatBody = document.querySelector(".chat-body");

  function checkEmptyChat() {
    if (chatBody && chatBody.children.length === 0) {
      const emptyMessage = document.createElement("div");
      emptyMessage.className = "empty-chat-message";
      emptyMessage.innerHTML = `
        <div class="empty-chat-content">
          <img src="IMG/empty-chat.png" alt="Empty chat" />
          <p>Bạn chưa có ai để nói chuyện</p>
        </div>
      `;
      chatBody.appendChild(emptyMessage);
    }
  }

  if (chatIcon && chatPopup && chatClose) {
    chatIcon.addEventListener("click", function (e) {
      e.stopPropagation();
      chatPopup.classList.toggle("hidden");
      checkEmptyChat();
    });

    chatClose.addEventListener("click", function () {
      chatPopup.classList.add("hidden");
    });

    document.addEventListener("click", function (e) {
      if (!chatPopup.contains(e.target) && e.target !== chatIcon) {
        chatPopup.classList.add("hidden");
      }
    });
  }

  // ==== User avatar & login ====
  const username = localStorage.getItem("username");
  const avatarUrl = localStorage.getItem("avatarUrl") || "IMG/ZenUser.png";

  if (username) {
    const loginLink = document.getElementById("loginLink");
    if (username && loginLink) {
      loginLink.style.display = "none";
    }

    const userAvatarDiv = document.getElementById("userAvatar");
    if (userAvatarDiv) userAvatarDiv.style.display = "block";

    const userAvatarImg = document.getElementById("userAvatarImg");

    function preloadImage(url, onSuccess, onError) {
      const img = new Image();
      img.onload = () => onSuccess(url);
      img.onerror = () => onError();
      userAvatarImg.src = "IMG/ZenUser.png";
    }

    if (userAvatarImg) {
      preloadImage(
        avatarUrl,
        (validUrl) => {
          userAvatarImg.src = validUrl;
        },
        () => {
          userAvatarImg.src = "IMG/ZenUser.png";
        }
      );
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

  // ==== Popup đăng bài ====
  const openBtn = document.getElementById("openPopupBtn");
  const popupOverlay = document.getElementById("popupOverlay");
  const closeBtn = document.getElementById("closePopupBtn");

  if (openBtn && popupOverlay && closeBtn) {
    openBtn.addEventListener("click", () => {
      popupOverlay.classList.add("active");
    });

    closeBtn.addEventListener("click", () => {
      popupOverlay.classList.remove("active");
    });

    popupOverlay.addEventListener("click", (e) => {
      if (e.target === popupOverlay) {
        popupOverlay.classList.remove("active");
      }
    });
  }
});

function hienThiDanhSachBaiDang(dsBaiDang) {
  const danhSach = document.getElementById("danhSachBaiDang");
  danhSach.innerHTML = "";

  if (dsBaiDang.length === 0) {
    danhSach.innerHTML = `
      <div style="padding: 20px; color: #999; text-align: center;">
        Không tìm thấy bài viết phù hợp.
      </div>
    `;
    danhSach.style.justifyContent = "center"; // nếu đang dùng flex
    return;
  }

  dsBaiDang.forEach((bai) => {
    const button = document.createElement("button");
    button.textContent = bai.tieuDe;
    button.className = "btn-tieu-de";
    button.onclick = () => {
      window.location.href = `baiviet.html?id=${bai.id}`;
    };
    danhSach.appendChild(button);
  });
}

// Gọi API lấy danh sách bài viết
fetch("http://localhost:8080/LayBaiDang")
  .then((res) => res.json())
  .then((data) => hienThiDanhSachBaiDang(data))
  .catch((err) => console.error("Lỗi tải bài viết:", err));

function xoaBaiViet(id) {
  if (!confirm("Bạn có chắc muốn xóa bài viết này?")) return;

  fetch(`http://localhost:8080/XoaBaiDang/${id}`, {
    method: "DELETE",
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

document.addEventListener("DOMContentLoaded", () => {
  let danhSachBaiViet = [];
  const danhSach = document.getElementById("danhSachBaiDang");
  danhSach.style.minHeight = "300px";

  fetch("http://localhost:8080/LayBaiDang")
    .then((res) => res.json())
    .then((data) => {
      danhSachBaiViet = data;
      hienThiDanhSachBaiDang(data);
    });

  const searchBox = document.getElementById("search-box");
  searchBox.addEventListener("input", () => {
    const tuKhoa = searchBox.value.toLowerCase();
    const ketQua = danhSachBaiViet.filter((bai) =>
      bai.tieuDe.toLowerCase().includes(tuKhoa)
    );
    hienThiDanhSachBaiDang(ketQua);
  });
});

document.addEventListener("DOMContentLoaded", () => {
  const thongBaoBtn = document.getElementById("thongBaoBtn");
  const soThongBao = document.getElementById("soThongBao");

  // Hàm kiểm tra bài viết mới định kỳ
  function kiemTraBaiMoi() {
    fetch("http://localhost:8080/LayBaiDang")
      .then((res) => res.json())
      .then((data) => {
        if (!Array.isArray(data) || data.length === 0) return;

        const latestPostId = data[0].id; // Giả sử bài viết mới nhất có id lớn nhất
        const lastSeenId =
          parseInt(localStorage.getItem("lastSeenPostId")) || 0;

        if (latestPostId > lastSeenId) {
          const soMoi = latestPostId - lastSeenId;
          soThongBao.textContent = `(${soMoi})`; // Gắn số vào nút thông báo
        } else {
          soThongBao.textContent = ""; // Không có bài mới thì không hiển thị
        }
      })
      .catch((err) => {
        console.error("Lỗi kiểm tra bài viết mới:", err);
      });
  }

  // Gọi khi bấm nút thông báo: cập nhật đã xem
  thongBaoBtn.addEventListener("click", () => {
    fetch("http://localhost:8080/LayBaiDang")
      .then((res) => res.json())
      .then((data) => {
        if (data && data.length > 0) {
          localStorage.setItem("lastSeenPostId", data[0].id);
          soThongBao.textContent = "";
          alert("Đã xem thông báo bài viết mới!");
        }
      });
  });

  // Kiểm tra tự động mỗi 15 giây
  setInterval(kiemTraBaiMoi, 15000); // 15 giây/lần
  kiemTraBaiMoi(); // kiểm tra ngay khi load
});
