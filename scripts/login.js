document.getElementById("loginForm").addEventListener("submit", async function (event) {
    event.preventDefault();

    const email = document.getElementById("email").value;
    const password = document.getElementById("password").value;

    const response = await fetch("http://localhost:8080/todo/login", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({ email, password })
    });

    if (!response.ok) {
        const errorMessage = await response.text();
        alert("Login failed: " + errorMessage);
        return;
    }

    const data = await response.json();
    localStorage.setItem("token", data.token);

    window.location.href = "index.html";  // ไปที่หน้า products
});

async function loginUser(event) {
    event.preventDefault();
    const email = document.getElementById("email").value;
    const password = document.getElementById("password").value;

    const response = await fetch("http://localhost:8080/todo/login", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({
            email: email,
            password: password
        })
    });

    if (response.ok) {
        const data = await response.json();  // แปลงคำตอบเป็น JSON
        alert(data.message);  // แสดงข้อความจาก JSON ที่ได้

        // บันทึก token ใน localStorage
        localStorage.setItem("token", data.token);
        // ย้ายไปยังหน้ารายการผลิตภัณฑ์
        window.location.href = "/index.html";
    } else {
        const errorData = await response.json();  // แปลงคำตอบเป็น JSON หากเกิดข้อผิดพลาด
        alert(errorData.message);
    }
}
