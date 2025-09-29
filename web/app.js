const API_URL = "http://localhost:8080"; // поменяй, если у тебя другой порт

// ---------- Items ----------
async function loadItems() {
	const res = await fetch(`${API_URL}/api/items`);
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
      <td><button onclick="deleteItem('${item.id}')">Delete</button></td>
    `;
		tbody.appendChild(tr);
	});
}

document.getElementById("item-form").addEventListener("submit", async (e) => {
	e.preventDefault();
	const formData = new FormData(e.target);
	const data = Object.fromEntries(formData.entries());
	data.amount = parseFloat(data.amount);
	data.occurred_at = new Date(data.occurred_at).toISOString();

	const res = await fetch(`${API_URL}/api/items`, {
		method: "POST",
		headers: { "Content-Type": "application/json" },
		body: JSON.stringify(data),
	});

	if (res.ok) {
		e.target.reset();
		loadItems();
	}
});

async function deleteItem(id) {
	await fetch(`${API_URL}/api/items/${id}`, { method: "DELETE" });
	loadItems();
}

// ---------- Categories ----------
async function loadCategories() {
	const res = await fetch(`${API_URL}/api/categories`);
	const json = await res.json();
	const list = document.getElementById("categories-list");
	list.innerHTML = "";
	(json.categories || []).forEach(cat => {
		const li = document.createElement("li");
		li.textContent = `${cat.id}: ${cat.name}`;
		list.appendChild(li);
	});
}

document.getElementById("category-form").addEventListener("submit", async (e) => {
	e.preventDefault();
	const formData = new FormData(e.target);
	const data = Object.fromEntries(formData.entries());

	const res = await fetch(`${API_URL}/api/categories`, {
		method: "POST",
		headers: { "Content-Type": "application/json" },
		body: JSON.stringify(data),
	});

	if (res.ok) {
		e.target.reset();
		loadCategories();
	}
});

// ---------- Analytics ----------
document.getElementById("analytics-form").addEventListener("submit", async (e) => {
	e.preventDefault();
	const formData = new FormData(e.target);
	const params = new URLSearchParams(formData);

	const [sumRes, avgRes, countRes, medianRes, percentileRes] = await Promise.all([
		fetch(`${API_URL}/api/analytics/sum?${params}`),
		fetch(`${API_URL}/api/analytics/avg?${params}`),
		fetch(`${API_URL}/api/analytics/count?${params}`),
		fetch(`${API_URL}/api/analytics/median?${params}`),
		fetch(`${API_URL}/api/analytics/percentile?${params}`),
	]);

	const [sum, avg, count, median, percentile] = await Promise.all([
		sumRes.json(),
		avgRes.json(),
		countRes.json(),
		medianRes.json(),
		percentileRes.json(),
	]);

	document.getElementById("analytics-result").innerHTML = `
    <p>Sum: ${sum.result}</p>
    <p>Avg: ${avg.result}</p>
    <p>Count: ${count.result}</p>
    <p>Median: ${median.result}</p>
    <p>Percentile: ${percentile.result}</p>
  `;
});

// ---------- Init ----------
loadItems();
loadCategories();
