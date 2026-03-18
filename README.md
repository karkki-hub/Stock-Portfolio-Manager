# 📈 Stock Portfolio Management System

A backend application built using **Go (Golang)** that allows users to manage stock investments, track portfolios, maintain watchlists, and generate daily reports using real-time stock data.

---

## 🚀 Features

### 🔐 Authentication
- User Registration (Email & Password)
- Password hashing using bcrypt
- JWT-based Login Authentication
- Protected API routes

### 📊 Stock Search
- Search stocks by symbol (e.g., AAPL, GOOGL)
- Fetch real-time stock data from Alpha Vantage API
- View current price, open, high, low, and volume

### ⭐ Watchlist
- Add stocks to watchlist
- Remove stocks from watchlist
- View watchlist with latest prices

### 💼 Portfolio Management
- Buy stocks (quantity & price)
- Sell stocks with validation
- Automatic calculations:
  - Average Buy Price
  - Total Investment
  - Current Value
  - Profit / Loss

### 🧾 Transactions
- Records every buy/sell transaction
- Stores quantity, price, total amount, and timestamp
- Linked to user via foreign key

### ⏰ Daily Cron Job
- Automatically updates stock prices daily
- Keeps portfolio values current
- Handles API errors gracefully

### 📄 Daily CSV Report
- Generates daily portfolio summary
- Stored in `/reports/YYYY-MM-DD/`
- Downloadable via API

---

## 🛠 Tech Stack

- **Language:** Go (Golang)
- **Framework:** Echo
- **Authentication:** JWT
- **Database:** MySQL
- **Scheduler:** robfig/cron
- **External API:** Alpha Vantage API

---
