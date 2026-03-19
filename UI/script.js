const API = "http://localhost:8080"

let tradetype = ""
let tradesymbol = ""

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

        const stock = data.data

        document.getElementById("stock").innerHTML =
        `
        <b>${stock.symbol}</b> - ${stock.stock_name} <br>
        Price: ${stock.last_price} <br>
        <button onclick="addWatch('${stock.symbol}')">Add to Watchlist</button>
        <button onclick="opentrade('BUY', '${stock.symbol}', ${stock.last_price})">Buy</button>
        <button onclick="opentrade('SELL', '${stock.symbol}', ${stock.last_price})">Sell</button>
        `
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
    .then(data => {   // ✅ FIX HERE

        const table = document.getElementById("transactions")
        table.innerHTML = ""

        if(!data || data.length === 0){
            table.innerHTML = "<tr><td colspan='6'>No transactions</td></tr>"
            return
        }

        const transactions = [...data].sort((a,b)=>
            new Date(b.created_at) - new Date(a.created_at)
        )

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

    const token = localStorage.getItem("token")

    if (!token) {
        alert("Please login again")
        window.location = "login.html"
        return
    }

    fetch(API + "/api/portfolio", {
        headers: {
            "Authorization": "Bearer " + token
        }
    })
    .then(res => {
        if (!res.ok) {
            throw new Error("Failed to fetch")
        }
        return res.json();
    })
    .then(data => {

        // ✅ Summary
        document.getElementById("totalInvestment").innerText = data.tot_investment || 0
        document.getElementById("currentValue").innerText = data.tot_cur_investment || 0
        document.getElementById("totalPL").innerText = data.total_profit_loss || 0

        // ✅ Table
        const table = document.getElementById("portfolio")
        table.innerHTML = ""

        if (!data.stocks || data.stocks.length === 0) {
            table.innerHTML = "<tr><td colspan='7'>No stocks found</td></tr>"
            return
        }

        // Optional: sort by highest profit
        const stocks = [...data.stocks].sort((a, b) => b.profit_loss - a.profit_loss)

        stocks.forEach(stock => {

            const tr = document.createElement("tr")

            tr.innerHTML = `
                <td>${stock.symbol || "N/A"}</td>
                <td>${stock.quantity || 0}</td>
                <td>${stock.avg_buy_price || 0}</td>
                <td>${stock.total_investment || 0}</td>
                <td>${stock.current_price || 0}</td>
                <td>${stock.current_value || 0}</td>
                <td style="color:${stock.profit_loss >= 0 ? 'green' : 'red'}">
                    ${stock.profit_loss || 0}
                </td>
            `

            table.appendChild(tr)
        })

    })
    .catch(err => {
        console.error(err)
        alert("Error loading portfolio")
    })
}

function goPortfolio(){
    window.location = "portfolio.html"
}