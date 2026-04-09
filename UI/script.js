const API = "http://localhost:8080"

let tradetype = ""
let tradesymbol = ""
let watchlistData = [];
let watchSortDirection = {};
let currentWatchSortField = "";
let chartInstance = null
let stockData = [];     // store API data globally
let sortDirection = {}; // track asc/desc per column
let currentSortField = "";
let transactionData = [];
let transSortDirection = {};
let currentTransSortField = "";
console.log("js loaded");

// ================= TRADE MODAL =================
function opentrade(type, symbol, price){
    tradetype = type
    tradesymbol = symbol

    document.getElementById("trade-title").innerText = type + " " + symbol
    document.getElementById("trade-price").value = price
    document.getElementById("trade-quantity").value = ""

    document.getElementById("trademodal").style.display = "block"
}

function closeTrade(){
    document.getElementById("trademodal").style.display = "none"
}

function submitTrade(){

    const quantity = document.getElementById("trade-quantity").value
    const price = document.getElementById("trade-price").value
    const token = localStorage.getItem("token")

    if(!quantity || !price){
        alert("Enter quantity and price")
        return
    }

    let url = tradetype === "BUY"
        ? API + "/api/transactions/buy"
        : API + "/api/transactions/sell"

    fetch(url,{
        method:"POST",
        headers:{
            "Content-Type":"application/json",
            "Authorization":"Bearer " + token
        },
        body:JSON.stringify({
            symbol: tradesymbol,
            quantity: parseInt(quantity),
            price: parseFloat(price)
        })
    })
    .then(res=>res.json())
    .then(data=>{
        alert(data.message)
        closeTrade()
    })
}

// ================= AUTH =================
function login() {

    const email = document.getElementById("email").value
    const password = document.getElementById("password").value

    fetch(API + "/login", {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify({email, password})
    })
    .then(res => res.json())
    .then(data => {

        if(data.status === "success"){
            localStorage.setItem("token", data.data.token)
            window.location = "watchlist.html"
        } else {
            document.getElementById("msg").innerText = data.message
        }
    })
}

function register(){

    const name = document.getElementById("name").value
    const email = document.getElementById("email").value
    const phone = document.getElementById("phone").value
    const password = document.getElementById("password").value
    const address = document.getElementById("address").value
    const api_key = document.getElementById("api_key").value

    fetch(API + "/register",{
        method:"POST",
        headers:{"Content-Type":"application/json"},
        body:JSON.stringify({name,email,phone,password,address,api_key})
    })
    .then(res=>res.json())
    .then(data=>{
        alert(data.message)
        if(data.status === "success"){
            window.location = "login.html"
        }
    })
}

// ================= STOCK SEARCH =================
function searchStock(){

    const symbol = document.getElementById("symbol").value
    const token = localStorage.getItem("token")

    fetch(API + "/api/stocks/" + symbol,{
        headers:{
            "Authorization":"Bearer " + token
        }
    })
    .then(res=>res.json())
    .then(data=>{

        const stocks = data.data
        const table = document.getElementById("stock")
        table.innerHTML = ""

        if(!stocks || stocks.length === 0){
            table.innerHTML = "<tr><td colspan='3'>No stocks found</td></tr>"
            return
        }

        stocks.forEach(stock => {

            const tr = document.createElement("tr")

            tr.innerHTML = `
                <td>${stock.symbol}</td>
                <td>${stock.stock_name}</td>
                <td>
                    <button onclick="addWatch('${stock.symbol}')">Add</button>
                    <button onclick="opentrade('BUY','${stock.symbol}',0)">Buy</button>
                    <button onclick="opentrade('SELL','${stock.symbol}',0)">Sell</button>
                </td>
            `

            table.appendChild(tr)
        })
    })
}

// ================= WATCHLIST =================
function addWatch(symbol){

    const token = localStorage.getItem("token")

    fetch(API + "/api/watchlist?symbol=" + symbol,{
        method:"POST",
        headers:{
            "Authorization":"Bearer " + token
        }
    })
    .then(res=>res.json())
    .then(data=>{
        alert(data.message)
        loadWatchlist()
    })
}

function loadWatchlist(){

    const token = localStorage.getItem("token")

    fetch(API + "/api/watchlist",{
        headers:{
            "Authorization":"Bearer " + token
        }
    })
    .then(res => res.json())
    .then(response => {

        const data = response.data   // ✅ FIX

        const table = document.getElementById("watchlist")
        table.innerHTML = ""

        watchlistData = data;
renderWatchlist();
    })
}

function removeWatch(symbol){

    const token = localStorage.getItem("token")

    fetch(API + "/api/watchlist/" + symbol,{
        method:"DELETE",
        headers:{
            "Authorization":"Bearer " + token
        }
    })
    .then(res=>res.json())
    .then(data=>{
        alert(data.message)
        loadWatchlist()
    })
}

// ================= TRANSACTIONS =================
function loadTransactions(){

    const token = localStorage.getItem("token")

    if(!token){
        alert("Please login again")
        window.location = "login.html"
        return
    }

    fetch(API + "/api/transactions",{
        headers:{
            "Authorization":"Bearer " + token
        }
    })
    .then(res => {
        if(!res.ok){
            throw new Error("Failed to fetch")
        }
        return res.json()
    })
    .then(response => {   // ✅ FIX HERE

        const data = response.data
        const table = document.getElementById("transactions")
        table.innerHTML = ""

        if(!data || data.length === 0){
            table.innerHTML = "<tr><td colspan='6'>No transactions</td></tr>"
            return
        }

        transactionData = data ? [...data] : [];
renderTransactions();
    })
    .catch(err => {
        console.error(err)
        alert("Error loading transactions")
    })
}

// ================= HELPERS =================
function formatDate(dateString){
    return new Date(dateString).toLocaleString()
}

function goToTransactions(){
    window.location = "transactions.html"
}

function logout(){
    localStorage.removeItem("token")
    window.location = "login.html"
}

function goregister(){
    window.location = "register.html"
}

function loadPortfolio() {

    const token = localStorage.getItem("token")

    fetch(API + "/api/portfolio", {
        headers: {
            "Authorization": "Bearer " + token
        }
    })
    .then(res => res.json())
    .then(response => {

        const data = response.data

        console.log("PORTFOLIO:", data)

        document.getElementById("totalInvestment").innerText = (data.tot_investment || 0).toFixed(2)
        document.getElementById("currentValue").innerText = (data.tot_cur_investment || 0).toFixed(2)
        document.getElementById("totalPL").innerText = (data.total_profit_loss || 0).toFixed(2)

        if (!data || !data.stocks || data.stocks.length === 0) {
            table.innerHTML = "<tr><td colspan='6'>No portfolio data</td></tr>"
            return
        }

        stockData = data.stocks
        renderPortfolioTable()
    })
}

function goPortfolio(){
    window.location = "portfolio.html"
}

function loadProfile() {

    const token = localStorage.getItem("token")

    fetch(API + "/api/profile", {
        headers: {
            "Authorization": "Bearer " + token
        }
    })
    .then(res => res.json())
    .then(response => {

        const data = response.data
        console.log("profile", data)

        document.getElementById("profilename").value = data.name || ""
        document.getElementById("profilemail").value = data.email || ""
        document.getElementById("profilePhone").value = data.phone || ""
        document.getElementById("profileAddress").value = data.address || ""

    })
    .catch(err => {
        console.error(err)
        document.getElementById("profileMsg").innerText = "Failed to load profile"
    })
}

function updateProfile() {

    const token = localStorage.getItem("token")

    const name = document.getElementById("profilename").value
    const phone = document.getElementById("profilePhone").value
    const address = document.getElementById("profileAddress").value

    if (!name || !phone || !address) {
        document.getElementById("profileMsg").innerText = "All fields required"
        return
    }

    fetch(API + "/api/profile/update", {
        method: "PUT",
        headers: {
            "Content-Type": "application/json",
            "Authorization": "Bearer " + token
        },
        body: JSON.stringify({
            Name: name,
            phone: phone,
            Address: address
        })
    })
    .then(res => res.json())
    .then(data => {
        document.getElementById("profileMsg").innerText = data.message || data.error
    })
    .catch(err => {
        console.error(err)
        document.getElementById("profileMsg").innerText = "Update failed"
    })
}

function goProfile(){
    window.location = "profile.html"
}

// function downloaduserReport() {
//     const token = localStorage.getItem("token")

//     fetch(API + "/api/report", {
//         headers: {
//             "Authorization": "Bearer " + token
//         }
//     })
//     .then(res => res.blob())
//     .then(blob => {
//         const url = window.URL.createObjectURL(blob)
//         const a = document.createElement("a")
//         const today = new Date().toISOString().split("T")[0]
//         const filename = `portfolio_report_${today}.csv`
//         a.href = url
//         console.log("Downloading report from:", a)
//         a.download = filename
//         document.body.appendChild(a)
//         a.click()
//         window.URL.revokeObjectURL(url)
//     })
//     .catch(err => {
//         console.error(err)
//         alert("Error downloading report")
//     })
// }

function downloaduserReport() {
    const token = localStorage.getItem("token")

    fetch(API + "/api/report", {
        headers: {
            "Authorization": "Bearer " + token
        }
    })
    .then(res => {
        if (!res.ok) {
            throw new Error("Failed to fetch report")
        }
        return res.blob()
    })
    .then(async (blob) => {
        const text = await blob.text()

        // 🧠 Extract from first line
        const firstLine = text.split("\n")[0]
        const parts = firstLine.split(",")

        // Expected: ["NAME:", "jon", "", "DATE:", "2026-04-09"]
        let username = parts[1] || "user"
        let date = parts[4] || new Date().toISOString().split("T")[0]

        // sanitize username
        const safeUsername = username.replace(/[^a-z0-9]/gi, "_").toLowerCase()

        const filename = `report_${safeUsername}_${date}.csv`

        // 🔽 Download
        const url = window.URL.createObjectURL(blob)
        const a = document.createElement("a")
        a.href = url
        a.download = filename

        document.body.appendChild(a)
        a.click()
        document.body.removeChild(a)

        window.URL.revokeObjectURL(url)
    })
    .catch(err => {
        console.error(err)
        alert("Error downloading report")
    })
}

function showGraph(symbol){

    const token = localStorage.getItem("token")

    fetch(API + "/api/watchlist/history/" + symbol, {
        headers: {
            "Authorization": "Bearer " + token
        }
    })
    .then(res => res.json())
    .then(response => {

        const data = response.data   // ✅ correct

        if(!data || data.length === 0){
            alert("No history available")
            return
        }

        // ✅ Format dates properly
        const labels = data.map(d => 
            new Date(d.date).toLocaleDateString()
        )

        const prices = data.map(d => d.price)

        document.getElementById("chartTitle").innerText = symbol + " Price History"
        document.getElementById("chartModal").style.display = "block"

        const ctx = document.getElementById("stockChart").getContext("2d")

        if(chartInstance){
            chartInstance.destroy()
        }

        chartInstance = new Chart(ctx, {
            type: "line",
            data: {
                labels: labels,
                datasets: [{
                    label: "Price",
                    data: prices,
                    tension: 0.3   // smooth curve
                }]
            },
            options: {
                responsive: true,
                plugins: {
                    legend: {
                        display: true
                    }
                }
            }
        })
    })
    .catch(err => {
        console.error(err)
        alert("Failed to load graph")
    })
}

function closeChart(){
    document.getElementById("chartModal").style.display = "none"
    if(chartInstance){
        chartInstance.destroy()
        chartInstance = null
    }
}

function openReset(){
    document.getElementById("resetModal").style.display = "block"
}

function closeReset(){
    document.getElementById("resetModal").style.display = "none"
}

function submitReset(){

    const token = localStorage.getItem("token")

    const oldpassword = document.getElementById("oldPass").value
    const newpassword = document.getElementById("newPass").value
    const reenterpassword = document.getElementById("rePass").value

    if(!oldpassword || !newpassword || !reenterpassword){
        document.getElementById("resetMsg").innerText = "All fields required"
        return
    }

    if(newpassword !== reenterpassword){
        document.getElementById("resetMsg").innerText = "Passwords do not match"
        return
    }

    fetch(API + "/api/profile/reset_pswd", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            "Authorization": "Bearer " + token
        },
        body: JSON.stringify({
            oldpassword,
            newpassword,
            reenterpassword
        })
    })
    .then(res => res.json())
    .then(data => {

        document.getElementById("resetMsg").innerText = data.message

        if(data.status === "success"){
            setTimeout(() => {
                closeReset()
            }, 1500)
        }
    })
    .catch(err => {
        console.error(err)
        document.getElementById("resetMsg").innerText = "Reset failed"
    })
}

function sortStocks(field) {
    // toggle direction
    sortDirection[field] = !sortDirection[field];
    currentSortField = field;

    // SORT DATA
    stockData.sort((a, b) => {
        let valA = a[field] ?? 0;
        let valB = b[field] ?? 0;

        if (typeof valA === "string") {
                return sortDirection[field]
                    ? valA.localeCompare(valB)
                    : valB.localeCompare(valA);
            }

            return sortDirection[field]
                ? valA - valB
                : valB - valA;
        });

    updateArrows();   // 🔥 update arrows
    renderPortfolioTable();
}

function renderPortfolioTable() {
    const table = document.getElementById("stockTable");
    table.innerHTML = "";

    stockData.forEach(stock => {
        const tr = document.createElement("tr");

        tr.innerHTML = `
            <td>${stock.symbol}</td>
            <td>${stock.quantity}</td>
            <td>${(stock.avg_buy_price || 0).toFixed(2)}</td>
            <td>${(stock.current_price || 0).toFixed(2)}</td>
            <td>${(stock.profit_loss || 0).toFixed(2)}</td>
            <td style="color:${stock.profit_loss >= 0 ? 'green' : 'red'}">
                ${stock.profit_loss >= 0 ? 'Profit' : 'Loss'}
            </td>
        `;

        table.appendChild(tr);
    });
}

function updateArrows() {
    // clear all arrows
    document.querySelectorAll("th span").forEach(span => {
        span.innerText = "";
    });

    // set arrow for active column
    const arrow = document.getElementById("arrow-" + currentSortField);

    if (arrow) {
        arrow.innerText = sortDirection[currentSortField] ? " ↑" : " ↓";
    }
}

function sortWatchlist(field) {

    watchSortDirection[field] = !watchSortDirection[field];
    currentWatchSortField = field;

    watchlistData.sort((a, b) => {
        let valA = a[field] ?? "";
        let valB = b[field] ?? "";

        // convert numbers if needed
        if (!isNaN(valA)) valA = Number(valA);
        if (!isNaN(valB)) valB = Number(valB);

        if (typeof valA === "string") {
            return watchSortDirection[field]
                ? valA.localeCompare(valB)
                : valB.localeCompare(valA);
        }

        return watchSortDirection[field]
            ? valA - valB
            : valB - valA;
    });

    updateWatchArrows();
    renderWatchlist();
}

function renderWatchlist() {

    const table = document.getElementById("watchlist");
    table.innerHTML = "";

    watchlistData.forEach(stock => {

        const tr = document.createElement("tr");

        tr.innerHTML = `
            <td>${stock.symbol}</td>
            <td>${stock.stock_name}</td>
            <td>${stock.last_price}</td>
            <td>
                <button onclick="opentrade('BUY','${stock.symbol}',${stock.last_price})">Buy</button>
                <button onclick="opentrade('SELL','${stock.symbol}',${stock.last_price})">Sell</button>
                <button onclick="removeWatch('${stock.symbol}')">Remove</button>
                <button onclick="showGraph('${stock.symbol}')">Graph</button>
            </td>
        `;

        table.appendChild(tr);
    });
}

function updateWatchArrows() {

    document.querySelectorAll("th span").forEach(span => {
        span.innerText = "";
    });

    const arrow = document.getElementById("arrow-" + currentWatchSortField);

    if (arrow) {
        arrow.innerText = watchSortDirection[currentWatchSortField] ? " ↑" : " ↓";
    }
}

function sortTransactions(field) {
    transSortDirection[field] = !transSortDirection[field];
    currentTransSortField = field;

    transactionData.sort((a, b) => {
        let valA = a[field] ?? "";
        let valB = b[field] ?? "";

        // Parse dates properly
        if (field === "created_at") {
            valA = new Date(valA);
            valB = new Date(valB);
        }

        // Convert to numbers if numeric
        if (!isNaN(valA) && field !== "created_at") valA = Number(valA);
        if (!isNaN(valB) && field !== "created_at") valB = Number(valB);

        if (typeof valA === "string") {
            return transSortDirection[field]
                ? valA.localeCompare(valB)
                : valB.localeCompare(valA);
        }

        return transSortDirection[field]
            ? valA - valB
            : valB - valA;
    });

    updateTransactionArrows();
    renderTransactions();
}

function renderTransactions() {
    const table = document.getElementById("transactions");
    table.innerHTML = "";

    if (!transactionData.length) {
        table.innerHTML = "<tr><td colspan='6'>No transactions</td></tr>";
        return;
    }

    transactionData.forEach(tx => {
        const tr = document.createElement("tr");

        tr.innerHTML = `
            <td style="color:${tx.type === 'BUY' ? 'green' : 'red'}">
                ${tx.type || "N/A"}
            </td>
            <td>${tx.symbol || tx.stock_id || "N/A"}</td>
            <td>${tx.quantity || 0}</td>
            <td>${tx.price || 0}</td>
            <td>${tx.total_amount || 0}</td>
            <td>${tx.created_at ? formatDate(tx.created_at) : "N/A"}</td>
        `;

        table.appendChild(tr);
    });
}

function updateTransactionArrows() {
    document.querySelectorAll("th span").forEach(span => {
        span.innerText = "";
    });

    const arrow = document.getElementById("arrow-" + currentTransSortField);

    if (arrow) {
        arrow.innerText = transSortDirection[currentTransSortField] ? " ↑" : " ↓";
    }
}

function goReports(){
    window.location = "reports.html"
}

function downloadReport(filename) {
    const token = localStorage.getItem("token")

    console.log("Downloading:", filename)
    console.log("Token:", token)

    fetch(API + "/api/reports/" + filename, {
        headers: {
            "Authorization": "Bearer " + token
        }
    })
    .then(res => {
        console.log("Status:", res.status)

        if (!res.ok) {
            return res.text().then(text => {
                console.error("Server error:", text)
                throw new Error(text)
            })
        }

        return res.blob()
    })
    .then(blob => {
        const url = window.URL.createObjectURL(blob)

        const a = document.createElement("a")
        a.href = url
        a.download = filename
        document.body.appendChild(a)
        a.click()
        a.remove()

        window.URL.revokeObjectURL(url)
    })
    .catch(err => {
        console.error("Download error:", err)
        alert("Download failed: " + err.message)
    })
}

function goWatchlist(){
    window.location = "watchlist.html"
}

let currentSort = { column: null, asc: true }

let allReports = []

document.addEventListener("DOMContentLoaded", () => {

    console.log("DOM loaded");

    // ✅ REPORTS PAGE
    if (document.getElementById("reportsTable")) {

        const filterBtn = document.getElementById("filter-btn")
        const resetBtn = document.getElementById("reset-btn")

        if (filterBtn) {
            filterBtn.addEventListener("click", applyDateFilter)
        }

        if (resetBtn) {
            resetBtn.addEventListener("click", resetDateFilter)
        }

        loadReports()
    }

})
function loadReports() {
    const token = localStorage.getItem("token")
    if (!token) return alert("Please login again")

    fetch(API + "/api/reports", {
        headers: { "Authorization": "Bearer " + token }
    })
    .then(res => res.json())
    .then(response => {
        allReports = response.data
        renderTable(allReports)
    })
    .catch(err => {
        console.error(err)
        alert("Failed to load reports")
    })
}

function renderTable(files) {
    const table = document.getElementById("reportsTable")
    table.innerHTML = ""

    // Header row with clickable sorting
    const header = document.createElement("tr")
    const columns = ["File Name", "Date", "Action"]
    columns.forEach((col, index) => {
        const th = document.createElement("th")
        th.textContent = col
        if (col !== "Action") { // make sortable
            th.style.cursor = "pointer"
            th.addEventListener("click", () => {
                sortTable(files, index)
                renderTable(files)
            })

            // add arrow indicator
            if (currentSort.column === index) {
                th.textContent += currentSort.asc ? " ▲" : " ▼"
            }
        }
        header.appendChild(th)
    })
    table.appendChild(header)

    // Data rows
    files.forEach(file => {
        const tr = document.createElement("tr")
        tr.innerHTML = `
            <td>${file[0]}</td>
            <td>${file[1]}</td>
            <td><button>Download</button></td>
        `
        const btn = tr.querySelector("button")
        btn.addEventListener("click", () => downloadReport(file[0]))
        table.appendChild(tr)
    })
}

function sortTable(files, columnIndex) {
    if (currentSort.column === columnIndex) {
        currentSort.asc = !currentSort.asc // toggle asc/desc
    } else {
        currentSort.column = columnIndex
        currentSort.asc = true
    }

    files.sort((a, b) => {
        let valA = a[columnIndex]
        let valB = b[columnIndex]

        // Detect if column is date (column 1)
        if (columnIndex === 1) {
            return currentSort.asc
                ? new Date(valA) - new Date(valB)
                : new Date(valB) - new Date(valA)
        }

        // Default: string comparison
        return currentSort.asc
            ? valA.localeCompare(valB)
            : valB.localeCompare(valA)
    })
}

function applyDateFilter() {
    const from = document.getElementById("date-from").value
    const to = document.getElementById("date-to").value

    const filtered = allReports.filter(file => {
        const reportDate = new Date(file[1])

        if (from && reportDate < new Date(from)) return false
        if (to && reportDate > new Date(to)) return false

        return true
    })

    renderTable(filtered)
}

// Reset filter
function resetDateFilter() {
    document.getElementById("date-from").value = ""
    document.getElementById("date-to").value = ""
    renderTable(allReports)
}

