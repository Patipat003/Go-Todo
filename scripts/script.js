const API_URL = "http://localhost:8080/todo/";

// ตรวจสอบว่า user ได้เข้าสู่ระบบแล้วหรือยัง
const token = localStorage.getItem("token");

if (!token) {
  // ถ้าไม่มี token ให้เปลี่ยนเส้นทางไปที่หน้า login.html
  window.location.href = "login.html";
}

// ฟังก์ชันสำหรับการดึง todos จาก API
async function fetchTodos() {
  const response = await fetch(API_URL, {
    headers: {
      "Authorization": `Bearer ${token}`
    }
  });

  if (!response.ok) {
    if (response.status === 401) {
      alert("Session expired. Please login again.");
      localStorage.removeItem("token");
      window.location.href = "login.html";
    } else {
      alert("Failed to fetch todos");
    }
    return;
  }

  const todos = await response.json();
  const todoList = document.getElementById("todos");
  todoList.innerHTML = "";
  todos.forEach(todo => {
    const todoItem = document.createElement("div");
    todoItem.className = "todo";
    todoItem.draggable = true;
    todoItem.setAttribute("data-id", todo.id);

    todoItem.innerHTML = `
      <label>
        <input type="checkbox" ${todo.complete ? "checked" : ""} onchange="toggleComplete(${todo.id}, this.checked)">
        ${todo.text}
      </label>
      <button onclick="editTodo(${todo.id})">Edit</button>
      <button onclick="deleteTodo(${todo.id})">Delete</button>
    `;

    todoItem.addEventListener("dragstart", handleDragStart);
    todoItem.addEventListener("dragover", handleDragOver);
    todoItem.addEventListener("drop", handleDrop);
    todoList.appendChild(todoItem);
  });
}

// ฟังก์ชันสำหรับการเพิ่ม Todo ใหม่
async function addTodo() {
  const text = document.getElementById("newTodoText").value;
  if (!text) return alert("Please enter a todo text");

  const response = await fetch(API_URL, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "Authorization": `Bearer ${token}`
    },
    body: JSON.stringify({ text, complete: false })
  });

  if (!response.ok) {
    alert("Failed to add todo");
    return;
  }

  document.getElementById("newTodoText").value = "";
  fetchTodos();
}

document.getElementById("newTodoText").addEventListener("keydown", function (event) {
  if (event.key === "Enter") {
    event.preventDefault();
    addTodo();
  }
});

async function toggleComplete(id, complete) {
  const response = await fetch(API_URL + id, {
    headers: {
      "Authorization": `Bearer ${token}` // ส่ง token ไปใน header
    },
  });

  if (!response.ok) {
    alert("Failed to fetch todo details. Please check your token or server.");
    return;
  }

  const todo = await response.json();
  await fetch(API_URL + id, {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
      "Authorization": `Bearer ${token}` // ส่ง token ไปใน header
    },
    body: JSON.stringify({
      text: todo.text,
      complete,
    }),
  });
  fetchTodos();
}

async function deleteTodo(id) {
  const response = await fetch(API_URL + id, {
    method: "DELETE",
    headers: {
      "Authorization": `Bearer ${token}`
    },
  });

  if (!response.ok) {
    alert("Failed to delete todo");
    return;
  }

  fetchTodos();
}

async function editTodo(id) {
  const response = await fetch(API_URL + id, {
    headers: {
      "Authorization": `Bearer ${token}`
    }
  });

  if (!response.ok) {
    alert("Failed to fetch todo details");
    return;
  }

  const todo = await response.json();
  const newText = prompt("Edit Todo:", todo.text);
  if (newText !== null && newText.trim() !== "") {
    const updateResponse = await fetch(API_URL + id, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        "Authorization": `Bearer ${token}`
      },
      body: JSON.stringify({ text: newText.trim(), complete: todo.complete })
    });

    if (!updateResponse.ok) {
      alert("Failed to edit todo");
      return;
    }

    fetchTodos();
  }
}

// ฟังก์ชันสำหรับจัดการการลากและวาง
function handleDragStart(event) {
  event.dataTransfer.setData("text/plain", event.target.getAttribute("data-id"));
}

function handleDragOver(event) {
  event.preventDefault();
}

function handleDrop(event) {
  event.preventDefault();
  const draggedId = event.dataTransfer.getData("text/plain");
  const targetElement = event.target.closest(".todo");

  if (targetElement) {
    const targetId = targetElement.getAttribute("data-id");

    if (draggedId !== targetId) {
      // ทำการย้ายเฉพาะใน DOM เท่านั้น
      const draggedElement = document.querySelector(`[data-id='${draggedId}']`);
      const targetParent = targetElement.parentNode;

      // ลบ draggedElement ออกจากตำแหน่งเดิม
      targetParent.removeChild(draggedElement);

      // ใส่ draggedElement ลงในตำแหน่งใหม่
      targetParent.insertBefore(draggedElement, targetElement);
    }
  }
}


document.getElementById("logoutButton").addEventListener("click", function () {
  localStorage.removeItem("token"); // ลบ token ออกจาก localStorage
  window.location.href = "login.html"; // เปลี่ยนเส้นทางไปยังหน้า login
});

fetchTodos();
