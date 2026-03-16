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

function logout(){
localStorage.removeItem("token")
window.location = "login.html"
}