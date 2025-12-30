const API_BASE = "/api";

// Create a new order
function createOrder() {
    const restaurant = document.getElementById("restaurant").value;
    const deadline = document.getElementById("deadline").value;
    const password = document.getElementById("password").value;

    if (!restaurant || !password) {
        alert("Please enter restaurant and password.");
        return;
    }

    fetch(`${API_BASE}/orders`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ restaurant, deadline, password })
    })
    .then(res => res.json())
    .then(data => {
        document.getElementById("create-result").innerText = "Order ID: " + data.id;
        document.getElementById("order-id").value = data.id;

        document.getElementById("restaurant").value = "";
        document.getElementById("deadline").value = "";
        document.getElementById("password").value = "";
    })
    .catch(err => alert("Error creating order"));
}

// Fetch an order and its items
function fetchOrder() {
    const id = document.getElementById("order-id").value;
    const password = document.getElementById("fetch-password").value;

    if (!id || !password) {
        alert("Please enter order ID and password.");
        return;
    }

    fetch(`${API_BASE}/orders/${id}?password=${encodeURIComponent(password)}`)
    .then(res => res.json())
    .then(data => {
        if (data.error) {
            alert(data.error);
            return;
        }

        let html = `<h4>${data.restaurant}</h4>
                    <p>Deadline: ${data.deadline ? new Date(data.deadline).toLocaleString() : "Not set"}</p>
                    <p>Items:</p>`;

        if (!data.items || data.items.length === 0) {
            html += "<p>No items yet</p>";
        } else {
            html += "<ul>";
            data.items.forEach(item => {
                html += `<li>${item.user_name} - ${item.item_name} : ""} ${item.notes ? '('+item.notes+')' : ''}</li>`;
            });
            html += "</ul>";
        }

        document.getElementById("order-data").innerHTML = html;
    })
    .catch(err => alert("Error fetching order"));
}
