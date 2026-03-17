const API = "http://localhost:8080"
let tradetype = ""
let tradesymbol = ""

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

    let url = ""

    if(tradetype === "BUY"){
        url = API + "/api/transactions/buy"
    } else {
        url = API + "/api/transactions/sell"
    }

    fetch(url,{
        method:"POST",
        headers:{
            "Content-Type":"application/json",
            "Authorization":"Bearer " + token
        },
        body:JSON.stringify({symbol: tradesymbol, quantity: parseInt(quantity), price: parseFloat(price)})
}
    )
    .then(res=>res.json())
    .then(data=>{
        alert(data.message)
        closeTrade()
    })
}

function login() {

    const email = document.getElementById("email").value
    const password = document.getElementById("password").value

    fetch(API + "/login", {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify({email: email, password: password})
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

        window.location = "login.html"
    })
}
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
.then(data => {

const table = document.getElementById("watchlist")
table.innerHTML = ""

// sort by symbol
const stocks = data.data.sort((a,b)=>{
    return a.symbol.localeCompare(b.symbol)
})

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

fetch("http://localhost:8080/api/watchlist/" + symbol,{
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

function goToTransactions(){
    window.location = "transactions.html"
}

function logout(){
localStorage.removeItem("token")
window.location = "login.html"
}

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
    .then(data => {

        console.log(data) // your array

        const table = document.getElementById("transactions")
        table.innerHTML = ""

        if(data.length === 0){
            table.innerHTML = "<tr><td colspan='6'>No transactions</td></tr>"
            return
        }

        // Already sorted DESC from backend, but safe:
        const transactions = [...data].sort((a,b)=>
            new Date(b.created_at) - new Date(a.created_at)
        )

        transactions.forEach(tx => {

            const tr = document.createElement("tr")

            tr.innerHTML = `
                <td style="color:${tx.type === 'BUY' ? 'green' : 'red'}">
                    ${tx.type}
                </td>
                <td>${tx.stock_id}</td>
                <td>${tx.quantity}</td>
                <td>₹ ${tx.price}</td>
                <td>₹ ${tx.total_amount}</td>
                <td>${formatDate(tx.created_at)}</td>
            `

            table.appendChild(tr)
        })

    })
    .catch(err => {
        console.error(err)
        alert("Error loading transactions")
    })
}

function formatDate(dateString){
    const d = new Date(dateString)
    return d.toLocaleString()
}


