// stockData.js


function extractStockData(response) {
    const globalQuote = response["Global Quote"];
    
    return {
        symbol: globalQuote["01. symbol"],
        open: parseFloat(globalQuote["02. open"]),
        high: parseFloat(globalQuote["03. high"]),
        low: parseFloat(globalQuote["04. low"]),
        price: parseFloat(globalQuote["05. price"]),
        volume: parseInt(globalQuote["06. volume"]),
        latestTradingDay: globalQuote["07. latest trading day"],
        previousClose: parseFloat(globalQuote["08. previous close"]),
        change: parseFloat(globalQuote["09. change"]),
        changePercent: parseFloat(globalQuote["10. change percent"].replace('%', ''))
    };
}

async function getStockQuote(symbol, accessToken) {
    const url = `https://default.gw.wso2.com:9095/query/1.0.0/query?function=GLOBAL_QUOTE&symbol=${symbol}&apikey=demo`;
    
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
        return data;
    } catch (error) {
        console.error('Error fetching stock quote:', error);
        throw error;
    }
}


document.addEventListener('DOMContentLoaded', function () {
    // Original stock data with initial prices
    const symbols = [
        "MMM","AOS","ABT","ABBV","ACN","ADBE","ADI","ADM","ADP","AES","AFL","A","APD","ABNB","AKAM","ALB","ARE","ALGN","ALLE","LNT","ALL","GOOGL","GOOG","MO","AMZN","AMCR","AMD","AEE","AAL","AEP","AXP","AIG","AMT","AWK","AMP","AME","AMGN","APH","APC","ADI","ANSS","AON","APA","AAPL","AMAT","APTV","ACGL","ADM","ANET","AJG","AIZ","T","ATO","ADSK","AZO","AVB","AVY","AXON","BKR","BALL","BAC","BK","BBWI","BAX","BDX","BBY","BIO","TECH","BIIB","BLK","BX","BA","BKNG","BWA","BSX","BMY","AVGO","BR","BRO","BF.B","BLDR","BG","BXP","CDNS","CZR","CPT","CPB","COF","CAH","KMX","CCL","CARR","CTLT","CAT","CBOE","CBRE","CDW","CE","COR","CNC","CNP","CF","CHRW","CRL","SCHW","CHTR","CVX","CMG","CB","CHD","CI","CINF","CTAS","CSCO","C","CFG","CLX","CME","CMS","KO","CTSH","CL","CMCSA","CAG","COP","ED","STZ","CEG","COO","CPRT","GLW","CPAY","CTVA","CSGP","COST","CTRA","CRWD","CCI","CSX","CMI","CVS","DHI","DHR","DRI","DVA","DAY","DECK","DE","DAL","DVN","DXCM","FANG","DLR","DFS","DG","DLTR","D","DPZ","DOV","DOW","DTE","DUK","DD","EMN","ETN","EBAY","ECL","EIX","EW","EA","ELV","LLY","EMR","ENPH","ETR","EOG","EPAM","EQT","EFX","EQIX","EQR","ESS","EL","ETSY","EG","EVRG","ES","EXC","EXPE","EXPD","EXR","XOM","FFIV","FDS","FICO","FAST","FRT","FDX","FIS","FITB","FSLR","FE","FI","FMC","F","FTNT","FTV","FOXA","FOX","BEN","FCX","GRMN","IT","GE","GEHC","GEV","GEN","GNRC","GD","GIS","GM","GPC","GILD","GPN","GL","GDDY","GS","HAL","HIG","HAS","HCA","DOC","HSIC","HSY","HES","HPE","HLT","HOLX","HD","HON","HRL","HST","HWM","HPQ","HUBB","HUM","HBAN","HII","IBM","IEX","IDXX","ITW","INCY","IR","PODD","INTC","ICE","IFF","IP","IPG","INTU","ISRG","IVZ","INVH","IQV","IRM","JBHT","JBL","JKHY","J","JNJ","JCI","JPM","JNPR","K","KVUE","KDP","KEY","KEYS","KMB","KIM","KMI","KKR","KLAC","KHC","KR","LHX","LH","LRCX","LW","LVS","LDOS","LEN","LIN","LYV","LKQ","LMT","L","LOW","LULU","LYB","MTB","MRO","MPC","MKTX","MAR","MMC","MLM","MAS","MA","MTCH","MKC","MCD","MCK","MDT","MRK","META","MET","MTD","MGM","MCHP","MU","MSFT","MAA","MRNA","MHK","MOH","TAP","MDLZ","MPWR","MNST","MCO","MS","MOS","MSI","MSCI","NDAQ","NTAP","NFLX","NEM","NWSA","NWS","NEE","NKE","NI","NDSN","NSC","NTRS","NOC","NCLH","NRG","NUE","NVDA","NVR","NXPI","ORLY","OXY","ODFL","OMC","ON","OKE","ORCL","OTIS","PCAR","PKG","PANW","PARA","PH","PAYX","PAYC","PYPL","PNR","PEP","PFE","PCG","PM","PSX","PNW","PNC","POOL","PPG","PPL","PFG","PG","PGR","PLD","PRU","PEG","PTC","PSA","PHM","QRVO","PWR","QCOM","DGX","RL","RJF","RTX","O","REG","REGN","RF","RSG","RMD","RVTY","ROK","ROL","ROP","ROST","RCL","SPGI","CRM","SBAC","SLB","STX","SRE","NOW","SHW","SPG","SWKS","SJM","SNA","SOLV","SO","LUV","SWK","SBUX","STT","STLD","STE","SYK","SMCI","SYF","SNPS","SYY","TMUS","TROW","TTWO","TPR","TRGP","TGT","TEL","TDY","TFX","TER","TSLA","TXN","TXT","TMO","TJX","TSCO","TT","TDG","TRV","TRMB","TFC","TYL","TSN","USB","UBER","UDR","ULTA","UNP","UAL","UPS","URI","UNH","UHS","VLO","VTR","VLTO","VRSN","VRSK","VZ","VRTX","VTRS","VICI","V","VST","VMC","WRB","GWW","WAB","WBA","WMT","DIS","WBD","WM","WAT","WEC","WFC","WELL","WST","WDC","WRK","WY","WMB","WTW","WYNN","XEL","XYL","YUM","ZBRA","ZBH","ZTS"
    ];
    const originalStockData = [
        { symbol: "NVDA", name: "NVIDIA Corporation", price: 106.98, change: -5.71, changePercent: -5.07, volume: "360.13M", avgVol: "265.50M", marketCap: "2.61T", peRatio: 36.39 },
        { symbol: "TSLA", name: "Tesla, Inc.", price: 222.15, change: -40.52, changePercent: -15.43, volume: "184.90M", avgVol: "87.44M", marketCap: "714.55B", peRatio: 108.90 },
        { symbol: "F", name: "Ford Motor Company", price: 9.96, change: +0.06, changePercent: +0.61, volume: "161.41M", avgVol: "81.31M", marketCap: "39.47B", peRatio: 6.82 },
        { symbol: "PLTR", name: "Palantir Technologies Inc.", price: 70.38, change: -8.53, changePercent: -10.85, volume: "133.31M", avgVol: "97.94M", marketCap: "179.13B", peRatio: 402.00 },
        { symbol: "LCID", name: "Lucid Group, Inc.", price: 2.0800, change: -0.0700, changePercent: -3.26, volume: "127.25M", avgVol: "92.82M", marketCap: "6.30B", peRatio: "-" },
        { symbol: "HOOD", name: "Robinhood Markets, Inc.", price: 36.50, change: -8.79, changePercent: -19.45, volume: "86.11M", avgVol: "29.58M", marketCap: "31.54B", peRatio: 22.84 },
        { symbol: "AAL", name: "American Airlines Group Inc.", price: 12.50, change: -0.53, changePercent: -4.07, volume: "81.91M", avgVol: "34.45M", marketCap: "8.22B", peRatio: 10.08 },
        { symbol: "GRAB", name: "Grab Holdings Limited", price: 4.1200, change: -0.4800, changePercent: -10.42, volume: "78.85M", avgVol: "34.30M", marketCap: "16.74B", peRatio: "-" },
        { symbol: "ALTM", name: "Arcadium Lithium plc", price: 5.84, change: 0.00, changePercent: 0.00, volume: "22.39M", avgVol: "18.53M", marketCap: "6.29B", peRatio: 64.89 },
        { symbol: "AAPL", name: "Apple Inc.", price: 167.48, change: -11.59, changePercent: -6.48, volume: "71.37M", avgVol: "52.01M", marketCap: "3.41T", peRatio: 36.11 },
        { symbol: "SMCI", name: "Super Micro Computer, Inc.", price: 38.00, change: -1.34, changePercent: -3.40, volume: "72.69M", avgVol: "70.88M", marketCap: "21.99B", peRatio: 16.04 },
        { symbol: "INTC", name: "Intel Corporation", price: 19.93, change: -0.71, changePercent: -3.44, volume: "79.67M", avgVol: "90.84M", marketCap: "86.29B", peRatio: "-" },
        { symbol: "SOFI", name: "SoFi Technologies, Inc.", price: 11.18, change: -1.41, changePercent: -11.20, volume: "69.59M", avgVol: "45.87M", marketCap: "12.25B", peRatio: 28.67 },
        { symbol: "VRN", name: "Veren Inc.", price: 6.70, change: +0.79, changePercent: +16.09, volume: "65.76M", avgVol: "17.71M", marketCap: "3.48B", peRatio: 17.81 },
        { symbol: "NIO", name: "NIO Inc.", price: 4.4600, change: -0.0100, changePercent: -0.22, volume: "64.90M", avgVol: "48.09M", marketCap: "9.32B", peRatio: "-" },
        { symbol: "NU", name: "Nu Holdings Ltd.", price: 10.12, change: -0.69, changePercent: -6.38, volume: "63.67M", avgVol: "48.06M", marketCap: "48.64B", peRatio: 25.30 },
        { symbol: "HOOD", name: "Robinhood Markets, Inc.", price: 36.50, change: -8.79, changePercent: -19.45, volume: "86.11M", avgVol: "29.58M", marketCap: "31.54B", peRatio: 22.84 },
        { symbol: "AAL", name: "American Airlines Group Inc.", price: 12.50, change: -0.53, changePercent: -4.07, volume: "81.91M", avgVol: "34.45M", marketCap: "8.22B", peRatio: 10.08 },
        { symbol: "GRAB", name: "Grab Holdings Limited", price: 4.1200, change: -0.4800, changePercent: -10.42, volume: "78.85M", avgVol: "34.30M", marketCap: "16.74B", peRatio: "-" },
        { symbol: "ALTM", name: "Arcadium Lithium plc", price: 5.84, change: 0.00, changePercent: 0.00, volume: "22.39M", avgVol: "18.53M", marketCap: "6.29B", peRatio: 64.89 },
        { symbol: "AAPL", name: "Apple Inc.", price: 167.48, change: -11.59, changePercent: -6.48, volume: "71.37M", avgVol: "52.01M", marketCap: "3.41T", peRatio: 36.11 },
        { symbol: "SMCI", name: "Super Micro Computer, Inc.", price: 38.00, change: -1.34, changePercent: -3.40, volume: "72.69M", avgVol: "70.88M", marketCap: "21.99B", peRatio: 16.04 },
        { symbol: "INTC", name: "Intel Corporation", price: 19.93, change: -0.71, changePercent: -3.44, volume: "79.67M", avgVol: "90.84M", marketCap: "86.29B", peRatio: "-" },
        { symbol: "SOFI", name: "SoFi Technologies, Inc.", price: 11.18, change: -1.41, changePercent: -11.20, volume: "69.59M", avgVol: "45.87M", marketCap: "12.25B", peRatio: 28.67 },
        { symbol: "VRN", name: "Veren Inc.", price: 6.70, change: +0.79, changePercent: +16.09, volume: "65.76M", avgVol: "17.71M", marketCap: "3.48B", peRatio: 17.81 },
        { symbol: "NIO", name: "NIO Inc.", price: 4.4600, change: -0.0100, changePercent: -0.22, volume: "64.90M", avgVol: "48.09M", marketCap: "9.32B", peRatio: "-" },
        { symbol: "NU", name: "Nu Holdings Ltd.", price: 10.12, change: -0.69, changePercent: -6.38, volume: "63.67M", avgVol: "48.06M", marketCap: "48.64B", peRatio: 25.30 }
    ];

    const tbody = document.getElementById('stockTableBody');

    // Function to populate the table
    function populateTable(data) {
        tbody.innerHTML = ''; // Clear existing rows
        data.forEach(stock => {
            const row = document.createElement('tr');
            row.innerHTML = `
                <td>${stock.symbol}</td>
                <td>${stock.name}</td>
                <td>${stock.price.toFixed(2)}</td>
                <td style="color: ${stock.change < 0 ? 'red' : 'green'}">${stock.change.toFixed(2)}</td>
                <td style="color: ${stock.change < 0 ? 'red' : 'green'}">${stock.changePercent.toFixed(2)}%</td>
                <td>${stock.volume}</td>
                <td>${stock.avgVol}</td>
                <td>${stock.marketCap}</td>
                <td>${stock.peRatio}</td>
            `;
            tbody.appendChild(row);
        });
    }

    // Initial population
    populateTable(originalStockData);
    console.log('Initial stock data:', originalStockData)
    // Update prices every 1 second
    setInterval(() => {
        const updatedStockData = originalStockData.map(stock => {
            // Generate a random price change between -2.0 and +2.0
            const priceChange = (Math.random() * 4 - 2).toFixed(2); // Random value between -2 and +2
            const newPrice = parseFloat((stock.price + parseFloat(priceChange)).toFixed(2));

            // Calculate new change and change percent based on the original price
            const originalPrice = stock.price;
            const newChange = (newPrice - originalPrice).toFixed(2);
            const newChangePercent = ((newChange / originalPrice) * 100).toFixed(2);

            return {
                ...stock,
                price: newPrice,
                change: parseFloat(newChange),
                changePercent: parseFloat(newChangePercent)
            };
        });
        console.log('Updated stock data:', updatedStockData)
        // Repopulate the table with updated data
        populateTable(updatedStockData);
    }, 1000); // Run every 1 second (1000 milliseconds)
});