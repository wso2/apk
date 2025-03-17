from flask import Flask, request, jsonify
import yfinance as yf
from datetime import datetime, timedelta

app = Flask(__name__)

def get_latest_trading_day():
    """Get the most recent trading day (assuming weekends/holidays skipped)."""
    today = datetime.now()
    if today.weekday() == 0:  # Monday
        latest_day = today - timedelta(days=3)
    elif today.weekday() == 6:  # Sunday
        latest_day = today - timedelta(days=2)
    else:
        latest_day = today - timedelta(days=1)
    return latest_day.strftime("%Y-%m-%d")

@app.route('/stock/<symbol>', methods=['GET'])
def get_stock_data(symbol):
    print("Hello", flush=True)
    try:
        # Use the latest trading day as the start and end date (1-day range)
        latest_day = get_latest_trading_day()
        end_day = (datetime.strptime(latest_day, "%Y-%m-%d") + timedelta(days=1)).strftime("%Y-%m-%d")

        # Fetch data from yfinance
        data = yf.download(symbol, start=latest_day, end=end_day, interval="1d")
        print("__________________")
        print("data", data, flush=True)
        print("__________________")
        if data.empty:
            return jsonify({"error": f"No data available for {symbol} on {latest_day}"}), 404

        # Extract scalar values from Series
        close = data["Close"].iloc[0]
        high = data["High"].iloc[0]
        low = data["Low"].iloc[0]
        open_price = data["Open"].iloc[0]
        volume = data["Volume"].iloc[0]
        abs_change = close - open_price
        pct_change = (close - open_price) / open_price * 100

        # Prepare response with scalar values
        stock_data = {
            "symbol": symbol.upper(),
            "date": latest_day,
            "close": round(float(close), 2),  # Ensure float for JSON
            "high": round(float(high), 2),
            "low": round(float(low), 2),
            "open": round(float(open_price), 2),
            "volume": int(volume),  # Convert to int directly
            "absolute_change": round(float(abs_change), 2),
            "percentage_change": round(float(pct_change), 2)
        }

        return jsonify(stock_data), 200

    except Exception as e:
        return jsonify({"error": str(e)}), 500

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5001, debug=True)  # Changed port to 5001