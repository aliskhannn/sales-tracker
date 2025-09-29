const API_BASE = "/api"; // Относительный путь через nginx

// ---------- Items ----------
async function loadItems() {
    try {
        const res = await fetch(`${API_BASE}/items`);
        if (!res.ok) throw new Error(`Failed to fetch items: ${res.status}`);
        const json = await res.json();

        const tbody = document.querySelector("#items-table tbody");
        tbody.innerHTML = "";

        (json.items || []).forEach(item => {
            const tr = document.createElement("tr");
            tr.innerHTML = `
                <td>${item.id}</td>
                <td>${item.kind}</td>
                <td>${item.title}</td>
                <td>${item.amount}</td>
                <td>${item.currency}</td>
                <td>${new Date(item.occurred_at).toLocaleString()}</td>
                <td><pre>${item.metadata ? JSON.stringify(item.metadata, null, 2) : "{}"}</pre></td>
                <td><button onclick="deleteItem('${item.id}')">Delete</button></td>
            `;
            tbody.appendChild(tr);
        });
    } catch (err) {
        console.error(err);
    }
}

document.getElementById("item-form").addEventListener("submit", async e => {
    e.preventDefault();
    const formData = new FormData(e.target);
    const data = Object.fromEntries(formData.entries());
    data.amount = parseFloat(data.amount);
    data.occurred_at = new Date(data.occurred_at).toISOString();

    // обработка metadata
    if (data.metadata) {
        try {
            data.metadata = JSON.parse(data.metadata);
        } catch {
            alert("Metadata must be valid JSON");
            return;
        }
    } else {
        data.metadata = {};
    }

    try {
        const res = await fetch(`${API_BASE}/items`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(data)
        });

        if (!res.ok) throw new Error(`Failed to create item: ${res.status}`);
        e.target.reset();
        loadItems();
    } catch (err) {
        console.error(err);
    }
});

async function deleteItem(id) {
    try {
        const res = await fetch(`${API_BASE}/items/${id}`, { method: "DELETE" });
        if (!res.ok) throw new Error(`Failed to delete item: ${res.status}`);
        loadItems();
    } catch (err) {
        console.error(err);
    }
}

// ---------- Categories ----------
async function loadCategories() {
    try {
        const res = await fetch(`${API_BASE}/categories`);
        if (!res.ok) throw new Error(`Failed to fetch categories: ${res.status}`);
        const json = await res.json();

        const list = document.getElementById("categories-list");
        list.innerHTML = "";

        (json.categories || []).forEach(cat => {
            const li = document.createElement("li");
            li.textContent = `${cat.id}: ${cat.name}`;
            list.appendChild(li);
        });
    } catch (err) {
        console.error(err);
    }
}

document.getElementById("category-form").addEventListener("submit", async e => {
    e.preventDefault();
    const formData = new FormData(e.target);
    const data = Object.fromEntries(formData.entries());

    try {
        const res = await fetch(`${API_BASE}/categories`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(data)
        });

        if (!res.ok) throw new Error(`Failed to create category: ${res.status}`);
        e.target.reset();
        loadCategories();
    } catch (err) {
        console.error(err);
    }
});

// ---------- Analytics ----------
document.getElementById("analytics-form").addEventListener("submit", async e => {
    e.preventDefault();
    const formData = new FormData(e.target);
    const params = new URLSearchParams(formData);

    try {
        const [sumRes, avgRes, countRes, medianRes, percentileRes] = await Promise.all([
            fetch(`${API_BASE}/analytics/sum?${params}`),
            fetch(`${API_BASE}/analytics/avg?${params}`),
            fetch(`${API_BASE}/analytics/count?${params}`),
            fetch(`${API_BASE}/analytics/median?${params}`),
            fetch(`${API_BASE}/analytics/percentile?${params}`)
        ]);

        const [sum, avg, count, median, percentile] = await Promise.all([
            sumRes.json(),
            avgRes.json(),
            countRes.json(),
            medianRes.json(),
            percentileRes.json()
        ]);

        document.getElementById("analytics-result").innerHTML = `
            <p>Sum: ${sum.result.sum ?? sum.result.percentile ?? sum.result.avg ?? sum.result.count ?? sum.result.median ?? "[unknown]"}</p>
            <p>Avg: ${avg.result.avg}</p>
            <p>Count: ${count.result.count}</p>
            <p>Median: ${median.result.median}</p>
            <p>Percentile: ${percentile.result.percentile}</p>
        `;
    } catch (err) {
        console.error(err);
    }
});

// ---------- Init ----------
loadItems();
loadCategories();
