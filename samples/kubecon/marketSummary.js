// marketSummary.js
document.addEventListener('DOMContentLoaded', function () {
    const marketSummaryData = {
        gainers: [
            { symbol: "PTGX", name: "Protagonist Ther...", price: 55.95, change: "+17.60 (+45.89%)" },
            { symbol: "VRN", name: "Veren Inc.", price: 6.70, change: "+0.79 (+16.09%)" },
            { symbol: "MWA", name: "Mueller Water Pr...", price: 26.50, change: "+1.85 (+7.39%)" },
            { symbol: "KTOS", name: "Kratos Defense...", price: 29.18, change: "+1.89 (+6.93%)" },
            { symbol: "BECN", name: "Beacon Roofing...", price: 118.52, change: "+7.42 (+6.68%)" }
        ],
        losers: [
            { symbol: "RDDT", name: "Reddit, Inc.", price: 107.29, change: "-26.69 (-19.92%)" },
            { symbol: "HOOD", name: "Robinhood Mark...", price: 36.50, change: "-8.79 (-19.45%)" },
            { symbol: "VICR", name: "Vicor Corporation", price: 48.94, change: "-10.61 (-17.82%)" },
            { symbol: "COIN", name: "Coinbase Global...", price: 179.23, change: "-38.22 (-17.58%)" },
            { symbol: "MSTR", name: "MicroStrategy In...", price: 239.27, change: "-47.91 (-16.68%)" }
        ],
        mostActive: [
            { symbol: "NVDA", name: "NVIDIA Corporat...", price: 106.98, change: "-5.71 (-5.07%)" },
            { symbol: "TSLA", name: "Tesla, Inc.", price: 222.15, change: "-40.52 (-15.43%)" },
            { symbol: "F", name: "Ford Motor Comp...", price: 9.96, change: "+0.06 (+0.61%)" }
        ]
    };

    // const tbody = document.getElementById('marketSummaryBody');

    // // Populate Top Gainers
    // tbody.innerHTML = '<tr><th colspan="2">TOP GAINERS</th></tr>';
    // marketSummaryData.gainers.forEach(stock => {
    //     const row = document.createElement('tr');
    //     row.innerHTML = `
    //         <td>${stock.symbol}</td>
    //         <td>${stock.price} ${stock.change}</td>
    //     `;
    //     tbody.appendChild(row);
    // });

    // // Populate Top Losers
    // tbody.innerHTML += '<tr><th colspan="2">TOP LOSERS</th></tr>';
    // marketSummaryData.losers.forEach(stock => {
    //     const row = document.createElement('tr');
    //     row.innerHTML = `
    //         <td>${stock.symbol}</td>
    //         <td>${stock.price} ${stock.change}</td>
    //     `;
    //     tbody.appendChild(row);
    // });

    // // Populate Most Active
    // tbody.innerHTML += '<tr><th colspan="2">MOST ACTIVE</th></tr>';
    // marketSummaryData.mostActive.forEach(stock => {
    //     const row = document.createElement('tr');
    //     row.innerHTML = `
    //         <td>${stock.symbol}</td>
    //         <td>${stock.price} ${stock.change}</td>
    //     `;
    //     tbody.appendChild(row);
    // });
});