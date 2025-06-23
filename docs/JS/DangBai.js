let currentUserId = null;
let currentRole = null;
let token = null;

document.addEventListener("DOMContentLoaded", function () {
  refreshUserInfo(); // 🛠 nạp token ngay lập tức sau khi DOM sẵn sàng
  const form = document.querySelector("#popupForm form");
  const popupOverlay = document.getElementById("popupOverlay");
  const closeBtn = document.getElementById("closePopupBtn");
  function refreshUserInfo() {
    token = localStorage.getItem("token");
    if (token) {
      const decoded = parseJwt(token);
      currentUserId = decoded?.id || null;
      currentRole = decoded?.phanQuyen || null;
    }
  }

  function parseJwt(token) {
    try {
      return JSON.parse(atob(token.split(".")[1]));
    } catch (e) {
      return null;
    }
  }
  console.log(parseJwt(localStorage.getItem("token")));
  async function loadBaiDang() {
    refreshUserInfo();
    try {
      const res = await fetch("http://localhost:8080/LayBaiDang");
      const data = await res.json();

      const container = document.getElementById("danhSachBaiDang");
      container.innerHTML = "";

      data.forEach(async (bai) => {
        const baiVietDiv = document.createElement("div");
        baiVietDiv.className = "bai-viet";
        baiVietDiv.style.marginBottom = "15px";
        baiVietDiv.style.borderBottom = "1px solid #eee";
        baiVietDiv.style.paddingBottom = "15px";

        let camXucHTML = "";
        try {
          const resCamXuc = await fetch(
            `http://localhost:8080/ThongKeCamXuc/${bai.id}`
          );
          const camXucData = await resCamXuc.json();
          const emojiMap = {
            like: "👍",
            love: "❤️",
            haha: "😂",
            wow: "😮",
            sad: "😢",
            angry: "😠",
          };
          for (const [k, v] of Object.entries(camXucData)) {
            camXucHTML += `<span style="margin-right:8px;">${
              emojiMap[k] || k
            } ${v}</span>`;
          }
        } catch (e) {
          camXucHTML = "<i>Không tải được cảm xúc</i>";
        }

        let html = `
  <div style="display: flex; justify-content: space-between; align-items: flex-start;">
    <div class="noi-dung-bai" style="flex-grow: 1; cursor: pointer;">
      <h3 style="margin: 0 0 5px 0;">${bai.tieuDe}</h3>
      <p style="margin: 0 0 5px 0; color: #555;">
        ${bai.noiDung.substring(0, 100)}...
      </p>
      <small style="color: #888;">
        Đăng bởi: ${bai.tenDangNhap} - ${bai.ngayDang}
      </small>
      <div style="margin-top: 10px; font-size: 18px;">
        ${camXucHTML}
      </div>
    </div>`;

        if (
          currentRole === "admin" ||
          parseInt(bai.idTaiKhoan) === currentUserId
        ) {
          html += `
    <div>
      ${
        parseInt(bai.idTaiKhoan) === currentUserId
          ? `<button class="btn-sua" data-id="${bai.id}" 
               style="margin-left: 10px; padding: 5px 10px; background-color: #ffc107; color: black; border: none; border-radius: 4px; cursor: pointer;">
               Sửa
             </button>`
          : ""
      }
      <button class="btn-xoa" data-id="${bai.id}" 
        style="margin-left: 10px; padding: 5px 10px; background-color: #ff4444; color: white; border: none; border-radius: 4px; cursor: pointer;">
        Xóa
      </button>
    </div>`;
        }

        html += `
  </div>
<div class="emoji-container">
  <div class="react-wrapper">
    <button class="react-toggle">😊</button>
    <div class="emoji-popover">
      <button class="emoji" data-type="like">👍</button>
      <button class="emoji" data-type="love">❤️</button>
      <button class="emoji" data-type="haha">😂</button>
      <button class="emoji" data-type="wow">😮</button>
      <button class="emoji" data-type="sad">😢</button>
      <button class="emoji" data-type="angry">😠</button>
    </div>
  </div>
</div>


`;
        baiVietDiv.innerHTML = html;
        // 👉 Gắn sự kiện cho các emoji trong popover
        baiVietDiv.querySelectorAll(".emoji").forEach((btn) => {
          btn.addEventListener("click", async (e) => {
            e.stopPropagation();
            const loaiCamXuc = btn.dataset.type;
            const idBaiDang = bai.id;

            try {
              const response = await fetch(
                `http://localhost:8080/ThemCamXuc/${idBaiDang}`,
                {
                  method: "POST",
                  headers: {
                    "Content-Type": "application/json",
                    Authorization: "Bearer " + token,
                  },
                  body: JSON.stringify({ loaiCamXuc }),
                }
              );

              if (response.ok) {
                loadBaiDang(); // hoặc chỉ cập nhật lại cảm xúc nếu muốn nhanh
              } else {
                const errText = await response.text();
                alert("Lỗi: " + errText);
              }
            } catch (err) {
              alert("Không thể gửi cảm xúc");
            }
          });
        });

        // 👉 Xem chi tiết
        baiVietDiv
          .querySelector(".noi-dung-bai")
          .addEventListener("click", (e) => {
            if (e.target.closest(".btn-sua, .btn-xoa, .btn-react")) return; // 👈 nếu nhấn nút thì không chuyển trang
            window.location.href = `baiviet.html?id=${bai.id}`;
          });

        // 👉 Xóa
        const btnXoa = baiVietDiv.querySelector(".btn-xoa");

        if (btnXoa) {
          btnXoa.addEventListener("click", async (e) => {
            e.stopPropagation();
            if (!confirm("Bạn có chắc chắn muốn xóa bài viết này?")) return;
            refreshUserInfo(); // 👈 thêm dòng này trước khi gọi fetch xoá!
            console.log("🔐 Token đang dùng để xoá:", token);
            console.log("🎟️ Token:", token);
            console.log("👤 Role:", currentRole);
            console.log("🪪 User ID:", currentUserId);
            console.log(
              "📤 Gửi đến:",
              `http://localhost:8080/XoaBaiDang/${bai.id}`
            );

            try {
              const response = await fetch(
                `http://localhost:8080/XoaBaiDang/${bai.id}`,
                {
                  method: "DELETE",
                  headers: {
                    "Content-Type": "application/json",
                    Authorization: "Bearer " + token,
                  },
                }
              );

              // 👉 Phân biệt rõ response.ok và response.status
              if (!response.ok) {
                const text = await response.text(); // đọc dưới dạng chuỗi
                throw new Error(`Lỗi ${response.status}: ${text}`);
              }

              const data = await response.json();
              if (data.success) {
                alert("✅ Đã xóa bài thành công!");
                baiVietDiv.remove();
              } else {
                alert("❌ Xóa thất bại: " + data.message);
              }
            } catch (err) {
              console.error("Lỗi khi gửi yêu cầu xóa:", err);
            }
          });
        }

        // 👉 Cảm xúc
        baiVietDiv.querySelectorAll(".btn-react").forEach((btn) => {
          btn.addEventListener("click", async (e) => {
            e.stopPropagation();
            const loaiCamXuc = btn.dataset.type;
            const idBaiDang = btn.dataset.id;

            try {
              const response = await fetch(
                `http://localhost:8080/ThemCamXuc/${idBaiDang}`,
                {
                  method: "POST",
                  headers: {
                    "Content-Type": "application/json",
                    Authorization: "Bearer " + token,
                  },
                  body: JSON.stringify({ loaiCamXuc }),
                }
              );

              if (response.ok) {
                loadBaiDang();
              } else {
                const errText = await response.text();
                alert("Lỗi: " + errText);
              }
            } catch (err) {
              alert("Không thể gửi cảm xúc");
            }
          });
        });
        const btnSua = baiVietDiv.querySelector(".btn-sua");
        if (btnSua) {
          btnSua.addEventListener("click", (e) => {
            e.stopPropagation(); // 👈 ngăn sự kiện "lan lên" vùng cha
            e.preventDefault();
            const popup = document.getElementById("popupSuaBai");
            const inputTieuDe = document.getElementById("popupEditTitle");
            const inputNoiDung = document.getElementById("popupEditContent");
            const btnLuu = document.getElementById("btnLuuSua");
            const btnHuy = document.getElementById("btnHuySua");

            inputTieuDe.value = bai.tieuDe;
            inputNoiDung.value = bai.noiDung;
            popup.style.display = "flex"; // hoặc "block" nếu không cần căn giữa
            popup.style.visibility = "visible";
            popup.style.opacity = "1";
            popup.style.pointerEvents = "auto";

            btnLuu.onclick = async () => {
              const tieuDeMoi = inputTieuDe.value.trim();
              const noiDungMoi = inputNoiDung.value.trim();
              if (!tieuDeMoi || !noiDungMoi) {
                alert("Không được để trống!");
                return;
              }

              try {
                const res = await fetch(
                  `http://localhost:8080/SuaBaiDang/${bai.id}`,
                  {
                    method: "PUT",
                    headers: {
                      "Content-Type": "application/json",
                      Authorization: "Bearer " + token,
                    },
                    body: JSON.stringify({
                      tieuDe: tieuDeMoi,
                      noiDung: noiDungMoi,
                    }),
                  }
                );

                const result = await res.json();
                alert(result.message || "Đã sửa bài viết!");
                popup.style.display = "none";
                loadBaiDang();
              } catch (err) {
                alert("Không thể gửi yêu cầu sửa.");
              }
            };

            btnHuy.onclick = () => {
              popup.style.display = "none";
            };
          });
        }

        container.appendChild(baiVietDiv);
      });
    } catch (err) {
      console.error("Lỗi khi tải bài viết:", err);
      document.getElementById("danhSachBaiDang").textContent =
        "Không thể tải danh sách bài viết.";
    }
  }

  loadBaiDang();

  // 👉 Submit đăng bài
  form.addEventListener("submit", async function (e) {
    e.preventDefault();
    const TieuDe = document.getElementById("title").value;
    const NoiDung = document.getElementById("content").value;

    try {
      const response = await fetch("http://localhost:8080/DangBai", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: "Bearer " + token,
        },
        body: JSON.stringify({ TieuDe, NoiDung }),
      });

      if (response.ok) {
        alert("Đăng bài thành công!");
        form.reset();
        popupOverlay.style.display = "none";
        loadBaiDang();
      } else {
        const error = await response.text();
        alert("Lỗi: " + error);
      }
    } catch (err) {
      alert("Không thể kết nối tới server");
    }
  });

  // 👉 Đóng popup
  closeBtn.addEventListener("click", function () {
    popupOverlay.style.display = "none";
  });
});
