document.addEventListener("DOMContentLoaded", function () {
  const form = document.querySelector("#popupForm form");
  const popupOverlay = document.getElementById("popupOverlay");
  const closeBtn = document.getElementById("closePopupBtn");

  // Hàm load bài viết
  async function loadBaiDang() {
    try {
      const res = await fetch("http://localhost:8080/LayBaiDang");
      const data = await res.json();

      const container = document.getElementById("danhSachBaiDang");
      container.innerHTML = ""; // Xóa cũ

      data.forEach((bai) => {
        const div = document.createElement("div");
        div.className = "bai-viet";
        div.innerHTML = `
        <button class="btn-tieu-de" style="display:block;width:100%;text-align:left;margin-bottom:10px;padding:10px">
          <strong>${bai.tieuDe}</strong><br/>
          <small>Người đăng: ${bai.tenDangNhap} - ${bai.ngayDang}</small>
        </button>
      `;

        div.querySelector("button").addEventListener("click", () => {
          window.location.href = `baiviet.html?id=${bai.id}`;
        });

        container.appendChild(div);
      });
    } catch (err) {
      console.error("Lỗi khi tải bài:", err);
      document.getElementById("danhSachBaiDang").textContent =
        "Không thể tải danh sách bài viết.";
    }
  }

  // Gọi khi tải trang
  loadBaiDang();

  // Submit form
  form.addEventListener("submit", async function (e) {
    e.preventDefault();

    const TieuDe = document.getElementById("title").value;
    const NoiDung = document.getElementById("content").value;
    const IDTaiKhoan = 1;

    const data = { TieuDe, NoiDung, IDTaiKhoan };

    try {
      const response = await fetch("http://localhost:8080/DangBai", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      });

      if (response.ok) {
        alert("Đăng bài thành công!");
        form.reset();
        popupOverlay.style.display = "none";
        loadBaiDang(); // Tải lại danh sách bài viết
      } else {
        const error = await response.text();
        alert("Lỗi: " + error);
      }
    } catch (err) {
      alert("Không thể kết nối tới server");
      console.error(err);
    }
  });

  closeBtn.addEventListener("click", function () {
    popupOverlay.style.display = "none";
  });
});
