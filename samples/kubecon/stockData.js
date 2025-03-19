// stockData.js
const accessToken = 'eyJ4NXQiOiJNVE15TkRSak1ETmlZemxrTjJVM1pUYzVZbVF3WW1WbFlqTTJaalJoWmpSaE4yTTVZak5oWmpZNVlqaGlNbVU0TXpZM1l6YzFNemd3TlRKaU1USTBNdyIsImtpZCI6Ik1UTXlORFJqTUROaVl6bGtOMlUzWlRjNVltUXdZbVZsWWpNMlpqUmhaalJoTjJNNVlqTmhaalk1WWpoaU1tVTRNelkzWXpjMU16Z3dOVEppTVRJME13X1JTMjU2IiwidHlwIjoiYXQrand0IiwiYWxnIjoiUlMyNTYifQ.eyJzdWIiOiIxZmYxMGFkZi02ZTU4LTRjYTktOTBiNy0yMTEyMjM5ZjU0YTMiLCJhdXQiOiJBUFBMSUNBVElPTiIsImF1ZCI6IkV1dVN5WjNfRVNJMk1UdFc5aHdscnMxTEtIY2EiLCJuYmYiOjE3NDIzNTc4OTMsImF6cCI6IkV1dVN5WjNfRVNJMk1UdFc5aHdscnMxTEtIY2EiLCJzY29wZSI6ImRlZmF1bHQiLCJpc3MiOiJodHRwczovL2FtLndzbzIuY29tOjQ0My9vYXV0aDIvdG9rZW4iLCJleHAiOjM2MTc0MjM1Nzg5MywiaWF0IjoxNzQyMzU3ODkzLCJqdGkiOiJhYjc2ZmE0ZS00ZGQwLTQ5ZjgtYTMxOS04YTMxMWYzZmY4ZTYiLCJjbGllbnRfaWQiOiJFdXVTeVozX0VTSTJNVHRXOWh3bHJzMUxLSGNhIn0.TEucKDOjtMsZQ4XF25JescujKJ6BnUWCKJHoLG0n6D0yv17oTzzaX0w7nO4RHmKB5stAwQuCz8U8E0llvuk8qUeN3F5oewcZ41mtPKxQZ370MrkxRzYsR8-FKbImhb0ofBFWpFN29AyVBUXl_YIfY0q8yY5z_E4kRealpW4_XbCX_Kfp34vSJrL_-BBjHtyLRtKrd0npn00VU5BI9ECK-HABgvzwf9zYRtwXu8M72vuVx4xwZWJ-4vjtBUKQZ9SkDkeOj-jra1GlPxqNHjgKSyhtINncxvWlrAbsr8PGXxFzkDmTxyzBiRN0L7tINvbMBQF38CymLsQY4U1VP7d5dw'; // Replace with your actual token
    
function extractStockData(response) {
    return {
        symbol: response.symbol,
        open: parseFloat(response.open),
        high: parseFloat(response.high),
        low: parseFloat(response.low),
        price: parseFloat(response.close), // Map "close" to "price" for consistency
        volume: parseInt(response.volume),
        latestTradingDay: response.date, // Map "date" to "latestTradingDay"
        previousClose: parseFloat(response.close) - parseFloat(response.absolute_change), // Calculate previous close
        change: parseFloat(response.absolute_change),
        changePercent: parseFloat(response.percentage_change)
    };
}

async function getStockQuote(symbol, accessToken) {
    const url = `https://default.gw.wso2.com:9095/t/forbizlogic.com/stock/1.0.0/stock/${symbol}`;
    // const url = `http://127.0.0.1:5001/stock/${symbol}`;
    
    try {
        const response = await fetch(url, {
            method: 'GET',
            headers: {
                'Accept': '*/*',
                'Authorization': `Bearer ${accessToken}`
            }
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();
        // console.log("data", data);
        return data;
    } catch (error) {
        console.error(`Error fetching stock quote for ${symbol}:`, error);
        return null; // Return null for failed requests
    }
}

document.addEventListener('DOMContentLoaded', function () {
    var symbols = ["AAPL","NVDA","MSFT","AMZN","GOOGL","GOOG","META","TSLA","AVGO","BRK-B","WMT","LLY","JPM","V","MA","ORCL","XOM","UNH","COST","PG","HD","NFLX","JNJ","BAC","CRM","ABBV","KO","TMUS","CVX","MRK","WFC","CSCO","ACN","NOW","AXP","MCD","PEP","BX","IBM","DIS","LIN","TMO","MS","ABT","ADBE","AMD","PM","ISRG","PLTR","GE","INTU","GS","CAT","TXN","QCOM","VZ","BKNG","DHR","T","BLK","RTX","SPGI","PFE","HON","NEE","CMCSA","ANET","AMGN","PGR","LOW","SYK","UNP","TJX","KKR","SCHW","ETN","AMAT","BA","BSX","C","UBER","COP","PANW","ADP","DE","FI","BMY","LMT","GILD","NKE","CB","UPS","ADI","MMC","MDT","VRTX","MU","SBUX","PLD","GEV","LRCX","MO","SO","EQIX","CRWD","PYPL","SHW","ICE","CME","AMT","APH","ELV","TT","MCO","CMG","INTC","KLAC","ABNB","DUK","PH","CDNS","WM","DELL","MDLZ","MAR","MSI","WELL","AON","REGN","CI","HCA","PNC","ITW","SNPS","CTAS","CL","USB","FTNT","ZTS","MCK","GD","TDG","CEG","AJG","EMR","MMM","ORLY","NOC","COF","ECL","EOG","FDX","BDX","APD","WMB","SPG","ADSK","RCL","RSG","CARR","CSX","HLT","DLR","TGT","KMI","OKE","TFC","AFL","GM","BK","ROP","MET","CPRT","FCX","CVS","PCAR","SRE","AZO","TRV","NXPI","JCI","GWW","NSC","PSA","SLB","AMP","ALL","FICO","MNST","PAYX","CHTR","AEP","ROST","PWR","CMI","AXON","VST","URI","MSCI","LULU","O","PSX","AIG","FANG","D","HWM","DHI","KR","NDAQ","OXY","EW","COR","KDP","FIS","KMB","NEM","DFS","PCG","TEL","MPC","FAST","AME","PEG","PRU","KVUE","STZ","GLW","LHX","GRMN","BKR","CBRE","CTVA","HES","CCI","DAL","CTSH","F","VRSK","EA","ODFL","XEL","TRGP","A","IT","LVS","SYY","VLO","OTIS","LEN","EXC","IR","YUM","KHC","GEHC","IQV","GIS","CCL","RMD","VMC","HSY","ACGL","IDXX","WAB","ROK","MLM","EXR","DD","ETR","DECK","EFX","UAL","WTW","TTWO","HIG","RJF","AVB","MTB","DXCM","ED","EBAY","HPQ","IRM","EIX","LYV","VICI","CNC","WEC","MCHP","HUM","ANSS","BRO","CSGP","MPWR","GDDY","TSCO","STT","CAH","GPN","FITB","XYL","HPE","KEYS","DOW","EQR","ON","PPG","K","SW","NUE","EL","BR","WBD","TPL","CHD","MTD","DOV","TYL","FTV","TROW","VLTO","EQT","SYF","NVR","DTE","VTR","AWK","ADM","NTAP","WST","CPAY","PPL","LYB","AEE","EXPE","HBAN","CDW","FE","HUBB","HAL","ROL","PHM","CINF","PTC","WRB","DRI","FOXA","FOX","IFF","SBAC","WAT","ERIE","TDY","ATO","RF","BIIB","ZBH","CNP","MKC","ES","WDC","TSN","TER","STE","PKG","CLX","NTRS","ZBRA","DVN","CBOE","WY","LUV","ULTA","CMS","INVH","FSLR","BF-B","LDOS","CFG","LH","VRSN","IP","ESS","PODD","COO","SMCI","STX","MAA","FDS","NRG","BBY","SNA","L","PFG","STLD","TRMB","OMC","CTRA","HRL","ARE","BLDR","JBHT","GEN","DGX","KEY","NI","MOH","PNR","J","DG","BALL","NWS","NWSA","UDR","HOLX","JBL","GPC","IEX","MAS","KIM","ALGN","DLTR","EXPD","EG","MRNA","LNT","AVY","BAX","TPR","VTRS","CF","FFIV","DPZ","AKAM","RL","TXT","SWKS","EVRG","EPAM","DOC","APTV","RVTY","AMCR","REG","POOL","INCY","BXP","KMX","CAG","HST","JKHY","SWK","DVA","CPB","CHRW","JNPR","CPT","TAP","NDSN","PAYC","UHS","NCLH","DAY","SJM","TECH","SOLV","ALLE","BG","AIZ","IPG","BEN","EMN","ALB","MGM","AOS","WYNN","PNW","ENPH","LKQ","FRT","CRL","GNRC","AES","GL","LW","HSIC","MKTX","MTCH","TFX","WBA","HAS","IVZ","APA","MOS","PARA","MHK","CE","HII","CZR","BWA","QRVO","FMC","AMTM"];
    symbols = ["AAPL","NVDA","MSFT","AMZN","GOOGL","GOOG","META","TSLA","AVGO","BRK-B","WMT","LLY","JPM","V","MA","ORCL","XOM","UNH","COST","PG","HD","NFLX","JNJ","BAC","CRM","ABBV","KO","TMUS","CVX","MRK","WFC","CSCO","ACN","NOW","AXP","MCD","PEP","BX","IBM","DIS","LIN","TMO","MS","ABT","ADBE","AMD","PM","ISRG"]
    const tbody = document.getElementById('stockTableBody');
    const batchSize = 20; // Number of symbols to fetch at a time
    let currentIndex = 0; // Starting index for round-robin

    // Function to populate the table
    function populateTable(data) {
        // console.log("data", data);
        tbody.innerHTML = ''; // Clear existing rows
        data.forEach(stock => {
            // console.log("stock", stock);
            if (stock) { // Only add valid stock data
                const row = document.createElement('tr');
                row.innerHTML = `
                    <td>${stock.symbol}</td>
                    <td>${stock.price.toFixed(2)}</td>
                    <td>${stock.open.toFixed(2)}</td>
                    <td>${stock.high.toFixed(2)}</td>
                    <td>${stock.low.toFixed(2)}</td>
                    <td style="color: ${stock.change < 0 ? 'red' : 'green'}">${stock.change.toFixed(2)}</td>
                    <td style="color: ${stock.change < 0 ? 'red' : 'green'}">${stock.changePercent.toFixed(2)}%</td>
                    <td>${stock.volume.toLocaleString()}</td>
                `;
                tbody.appendChild(row);
            }
        });
    }

    // Function to fetch data for a batch of symbols
    async function fetchStockBatch() {
        // Get the next 20 symbols in round-robin fashion
        const batchSymbols = [];
        for (let i = 0; i < batchSize; i++) {
            const index = (currentIndex + i) % symbols.length;
            batchSymbols.push(symbols[index]);
        }

        // Fetch data for all symbols in the batch concurrently
        const promises = batchSymbols.map(symbol => 
            getStockQuote(symbol, accessToken)
                .then(data => {
                    // console.log(`Raw data for ${symbol}:`, data); // Log the raw response
                    const result = data ? extractStockData(data) : null;
                    // console.log(`Processed data for ${symbol}:`, result); // Log the processed result
                    return result; // Ensure the transformed data is returned
                })
                .catch(() => {
                    console.log(`Failed to fetch data for ${symbol}`); // Log the failure
                    return null; // Handle individual failures gracefully
                })
        );

        const stockData = await Promise.all(promises);
        // console.log("stockData", stockData);
        populateTable(stockData.filter(data => data !== null)); // Filter out failed requests

        // Update the index for the next batch
        currentIndex = (currentIndex + batchSize) % symbols.length;
    }

    // Initial fetch
    fetchStockBatch();

    // Update every second
    setInterval(fetchStockBatch, 15000);
});