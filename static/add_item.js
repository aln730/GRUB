const API_BASE = "/api";

function addItem() {
    const orderId = document.getElementById("order-id").value;
    const userName = document.getElementById("item-user").value;
    const itemName = document.getElementById("item-name").value;
    const notes = document.getElementById("item-notes").value;

    if (!orderId || !userName || !itemName) {
        alert("Please fill out all required fields.");
        return;
    }

    fetch(`${API_BASE}/orders/${orderId}/items`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ user_name: userName, item_name: itemName, notes })
    })
    .then(res => res.json())
    .then(data => {
        if (data.status === "success") {
            alert("Item added!");
            fetchOrder(orderId);

            // Clear form
            document.getElementById("item-user").value = "";
            document.getElementById("item-name").value = "";
            document.getElementById("item-notes").value = "";
        } else {
            alert("Failed to add item");
        }
    })
    .catch(err => alert("Error adding item"));
}

function fetchOrder(orderId) {
    fetch(`${API_BASE}/orders/${orderId}`)
    .then(res => res.json())
    .then(data => {
        let html = `<h4>${data.restaurant}</h4>
                    <p>Deadline: ${new Date(data.deadline).toLocaleString()}</p>
                    <p>Items:</p>`;

        if (data.items.length === 0) {
            html += "<p>No items yet</p>";
        } else {
            html += "<ul>";
            data.items.forEach(item => {
                html += `<li>${item.user_name} - ${item.item_name} ${item.notes ? '('+item.notes+')' : ''}</li>`;
            });
            html += "</ul>";
        }

        document.getElementById("order-data").innerHTML = html;
    });
}