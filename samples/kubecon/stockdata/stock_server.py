from flask import Flask, request, jsonify
import yfinance as yf
from datetime import datetime, timedelta
import threading
import time
import json
import os

app = Flask(__name__)

# Predefined list of stock symbols
STOCK_LIST = ["AAPL","NVDA","MSFT","AMZN","GOOGL","GOOG","META","TSLA","AVGO","BRK-B","WMT","LLY","JPM","V","MA","ORCL","XOM","UNH","COST","PG","HD","NFLX","JNJ","BAC","CRM","ABBV","KO","TMUS","CVX","MRK","WFC","CSCO","ACN","NOW","AXP","MCD","PEP","BX","IBM","DIS","LIN","TMO","MS","ABT","ADBE","AMD","PM","ISRG","PLTR","GE","INTU","GS","CAT","TXN","QCOM","VZ","BKNG","DHR","T","BLK","RTX","SPGI","PFE","HON","NEE","CMCSA","ANET","AMGN","PGR","LOW","SYK","UNP","TJX","KKR","SCHW","ETN","AMAT","BA","BSX","C","UBER","COP","PANW","ADP","DE","FI","BMY","LMT","GILD","NKE","CB","UPS","ADI","MMC","MDT","VRTX","MU","SBUX","PLD","GEV","LRCX","MO","SO","EQIX","CRWD","PYPL","SHW","ICE","CME","AMT","APH","ELV","TT","MCO","CMG","INTC","KLAC","ABNB","DUK","PH","CDNS","WM","DELL","MDLZ","MAR","MSI","WELL","AON","REGN","CI","HCA","PNC","ITW","SNPS","CTAS","CL","USB","FTNT","ZTS","MCK","GD","TDG","CEG","AJG","EMR","MMM","ORLY","NOC","COF","ECL","EOG","FDX","BDX","APD","WMB","SPG","ADSK","RCL","RSG","CARR","CSX","HLT","DLR","TGT","KMI","OKE","TFC","AFL","GM","BK","ROP","MET","CPRT","FCX","CVS","PCAR","SRE","AZO","TRV","NXPI","JCI","GWW","NSC","PSA","SLB","AMP","ALL","FICO","MNST","PAYX","CHTR","AEP","ROST","PWR","CMI","AXON","VST","URI","MSCI","LULU","O","PSX","AIG","FANG","D","HWM","DHI","KR","NDAQ","OXY","EW","COR","KDP","FIS","KMB","NEM","DFS","PCG","TEL","MPC","FAST","AME","PEG","PRU","KVUE","STZ","GLW","LHX","GRMN","BKR","CBRE","CTVA","HES","CCI","DAL","CTSH","F","VRSK","EA","ODFL","XEL","TRGP","A","IT","LVS","SYY","VLO","OTIS","LEN","EXC","IR","YUM","KHC","GEHC","IQV","GIS","CCL","RMD","VMC","HSY","ACGL","IDXX","WAB","ROK","MLM","EXR","DD","ETR","DECK","EFX","UAL","WTW","TTWO","HIG","RJF","AVB","MTB","DXCM","ED","EBAY","HPQ","IRM","EIX","LYV","VICI","CNC","WEC","MCHP","HUM","ANSS","BRO","CSGP","MPWR","GDDY","TSCO","STT","CAH","GPN","FITB","XYL","HPE","KEYS","DOW","EQR","ON","PPG","K","SW","NUE","EL","BR","WBD","TPL","CHD","MTD","DOV","TYL","FTV","TROW","VLTO","EQT","SYF","NVR","DTE","VTR","AWK","ADM","NTAP","WST","CPAY","PPL","LYB","AEE","EXPE","HBAN","CDW","FE","HUBB","HAL","ROL","PHM","CINF","PTC","WRB","DRI","FOXA","FOX","IFF","SBAC","WAT","ERIE","TDY","ATO","RF","BIIB","ZBH","CNP","MKC","ES","WDC","TSN","TER","STE","PKG","CLX","NTRS","ZBRA","DVN","CBOE","WY","LUV","ULTA","CMS","INVH","FSLR","BF-B","LDOS","CFG","LH","VRSN","IP","ESS","PODD","COO","SMCI","STX","MAA","FDS","NRG","BBY","SNA","L","PFG","STLD","TRMB","OMC","CTRA","HRL","ARE","BLDR","JBHT","GEN","DGX","KEY","NI","MOH","PNR","J","DG","BALL","NWS","NWSA","UDR","HOLX","JBL","GPC","IEX","MAS","KIM","ALGN","DLTR","EXPD","EG","MRNA","LNT","AVY","BAX","TPR","VTRS","CF","FFIV","DPZ","AKAM","RL","TXT","SWKS","EVRG","EPAM","DOC","APTV","RVTY","AMCR","REG","POOL","INCY","BXP","KMX","CAG","HST","JKHY","SWK","DVA","CPB","CHRW","JNPR","CPT","TAP","NDSN","PAYC","UHS","NCLH","DAY","SJM","TECH","SOLV","ALLE","BG","AIZ","IPG","BEN","EMN","ALB","MGM","AOS","WYNN","PNW","ENPH","LKQ","FRT","CRL","GNRC","AES","GL","LW","HSIC","MKTX","MTCH","TFX","WBA","HAS","IVZ","APA","MOS","PARA","MHK","CE","HII","CZR","BWA","QRVO","FMC","AMTM"]



# File to store stock data
DATA_FILE = "stock_data.json"

# Lock for file access
file_lock = threading.Lock()

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

def fetch_stock_data(symbol):
    """Fetch stock data for a given symbol using yfinance."""
    try:
        latest_day = get_latest_trading_day()
        end_day = (datetime.strptime(latest_day, "%Y-%m-%d") + timedelta(days=1)).strftime("%Y-%m-%d")
        
        data = yf.download(symbol, start=latest_day, end=end_day, interval="1d", threads=True)
        
        if data.empty:
            return None

        close = data["Close"].iloc[0]
        high = data["High"].iloc[0]
        low = data["Low"].iloc[0]
        open_price = data["Open"].iloc[0]
        volume = data["Volume"].iloc[0]
        abs_change = close - open_price
        pct_change = (close - open_price) / open_price * 100

        stock_data = {
            "symbol": symbol.upper(),
            "date": latest_day,
            "close": round(float(close), 2),
            "high": round(float(high), 2),
            "low": round(float(low), 2),
            "open": round(float(open_price), 2),
            "volume": int(volume),
            "absolute_change": round(float(abs_change), 2),
            "percentage_change": round(float(pct_change), 2),
            "timestamp": datetime.now().isoformat()
        }
        return stock_data
    except Exception as e:
        print(f"Error fetching data for {symbol}: {str(e)}")
        return None

def update_stock_file():
    """Background thread to update stock data in round-robin fashion."""
    index = 0
    while True:
        symbol = STOCK_LIST[index % len(STOCK_LIST)]
        stock_data = fetch_stock_data(symbol)
        
        if stock_data:
            with file_lock:
                # Read existing data
                if os.path.exists(DATA_FILE):
                    with open(DATA_FILE, 'r') as f:
                        try:
                            data = json.load(f)
                        except json.JSONDecodeError:
                            data = {}
                else:
                    data = {}
                
                # Update with new data
                data[symbol.upper()] = stock_data
                
                # Write back to file
                with open(DATA_FILE, 'w') as f:
                    json.dump(data, f)
        
        index += 1
        time.sleep(5)  # Wait 5 seconds before next fetch

@app.route('/stock/<symbol>', methods=['GET'])
def get_stock_data(symbol):
    symbol = symbol.upper()
    
    # Try to get data from file first
    with file_lock:
        if os.path.exists(DATA_FILE):
            with open(DATA_FILE, 'r') as f:
                try:
                    data = json.load(f)
                    if symbol in data:
                        return jsonify(data[symbol]), 200
                except json.JSONDecodeError:
                    pass

    # If not in file or file doesn't exist, fetch fresh data
    stock_data = fetch_stock_data(symbol)
    
    if stock_data is None:
        return jsonify({"error": f"No data available for {symbol}"}), 404
    
    # Update file with fresh data
    with file_lock:
        if os.path.exists(DATA_FILE):
            with open(DATA_FILE, 'r') as f:
                try:
                    data = json.load(f)
                except json.JSONDecodeError:
                    data = {}
        else:
            data = {}
        
        data[symbol] = stock_data
        
        with open(DATA_FILE, 'w') as f:
            json.dump(data, f)
    
    return jsonify(stock_data), 200

if __name__ == '__main__':
    # Start the background thread
    thread = threading.Thread(target=update_stock_file, daemon=True)
    thread.start()
    
    # Start the Flask app
    app.run(host='0.0.0.0', port=5001, debug=True)