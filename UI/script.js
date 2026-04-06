const API = "http://localhost:8080"

let tradetype = ""
let tradesymbol = ""

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

        const stocks = data.sort((a,b)=>
            a.symbol.localeCompare(b.symbol)
        )

        stocks.forEach(stock => {

            const tr = document.createElement("tr")

            tr.innerHTML =
            `
            <td>${stock.symbol}</td>
            <td>${stock.stock_name}</td>
            <td>${stock.last_price}</td>
            <td>
                <button onclick="opentrade('BUY','${stock.symbol}',${stock.last_price})">Buy</button>
                <button onclick="opentrade('SELL','${stock.symbol}',${stock.last_price})">Sell</button>
                <button onclick="removeWatch('${stock.symbol}')">Remove</button>
            </td>
            `

            table.appendChild(tr)
        })
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

        const transactions = [...data].sort((a,b)=>
            new Date(b.created_at) - new Date(a.created_at)
        )

        console.log("Transactions:", transactions)

        transactions.forEach(tx => {

            const tr = document.createElement("tr")

        
            tr.innerHTML = `
                <td style="color:${tx.type === 'BUY' ? 'green' : 'red'}">
                    ${tx.type || "N/A"}
                </td>
                <td>${tx.symbol || tx.stock_id || "N/A"}</td>
                <td>${tx.quantity || 0}</td>
                <td>${tx.price || 0}</td>
                <td>${tx.total_amount || 0}</td>
                <td>${tx.created_at ? formatDate(tx.created_at) : "N/A"}</td>
            `

            table.appendChild(tr)
        })
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

    console.log("Function called")

    const token = localStorage.getItem("token")

    fetch(API + "/api/portfolio", {
        headers: {
            "Authorization": "Bearer " + token
        }
    })
    .then(res => res.json())
    .then(response => {
        const data = response.data

        console.log("DATA:", data)

        document.getElementById("totalInvestment").innerText = (data.tot_investment || 0).toFixed(2)
        document.getElementById("currentValue").innerText = (data.tot_cur_investment || 0).toFixed(2)
        document.getElementById("totalPL").innerText = (data.total_profit_loss || 0).toFixed(2)

        const table = document.getElementById("stockTable") // ✅ FIXED
        table.innerHTML = "";

        data.stocks.forEach(stock => {
            const tr = document.createElement("tr")

            tr.innerHTML = `
                <td>${stock.symbol}</td>
                <td>${stock.quantity}</td>
                <td>${(stock.avg_buy_price || 0).toFixed(2)}</td>
                <td>${(stock.current_price || 0).toFixed(2)}</td>
                <td>${(stock.profit_loss || 0).toFixed(2)}</td>
                <td style="color:${stock.profit_loss >= 0 ? 'green' : 'red'}">${stock.profit_loss >= 0 ? 'Profit' : 'Loss'}</td>
            `

            table.appendChild(tr)
        })

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
    .then(data => {
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

    const email = document.getElementById("profileEmail").value
    const phone = document.getElementById("profilePhone").value
    const address = document.getElementById("profileAddress").value

    if (!email || !phone || !address) {
        document.getElementById("profileMsg").innerText = "All fields required"
        return
    }

    fetch(API + "/api/profile", {
        method: "PUT",
        headers: {
            "Content-Type": "application/json",
            "Authorization": "Bearer " + token
        },
        body: JSON.stringify({
            email: email,
            phone: phone,
            address: address
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