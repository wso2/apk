function addLoading() {
    const chatBox = document.querySelector(".chat-box");
    
    // Add the loading class
    chatBox.classList.add("loading");

}

function removeLoading() {
    const chatBox = document.querySelector(".chat-box");
    
    // Add the loading class
    chatBox.classList.remove("loading");

}

async function sendQuestion() {
    const url = 'https://default.gw.wso2.com:9095/t/forbizlogic.com/chat/1.0.0/chat/completions';
    const apiKey = 'eyJ4NXQjUzI1NiI6Ik1XVXhZbU5rWldJd1lUUTJORGN5TVRVd1l6VTFOVFF5WVRsall6QXlaak01TkRneFpUVmtaREZsTm1WaE5Ea3pZemd5WWpBeU0yTmlaVEF6WWpRMFl3PT0iLCJraWQiOiJnYXRld2F5X2NlcnRpZmljYXRlX2FsaWFzIiwidHlwIjoiSldUIiwiYWxnIjoiUlMyNTYifQ==.eyJzdWIiOiJhZG1pbkBjYXJib24uc3VwZXIiLCJhcHBsaWNhdGlvbiI6eyJpZCI6NiwidXVpZCI6ImRmNWFhOWMwLWJhNDUtNGQ1Zi1iNjM1LWNjNGViN2U1YjRhNiJ9LCJpc3MiOiJodHRwczpcL1wvYW0ud3NvMi5jb206NDQzXC9vYXV0aDJcL3Rva2VuIiwia2V5dHlwZSI6IlBST0RVQ1RJT04iLCJwZXJtaXR0ZWRSZWZlcmVyIjoiIiwidG9rZW5fdHlwZSI6ImFwaUtleSIsInBlcm1pdHRlZElQIjoiIiwiaWF0IjoxNzQyMjg0NDIwLCJqdGkiOiI0M2VjMDhjMy0zN2Q2LTRjNmMtYjYyNC0wYjI5YzUxOThjZjkifQ==.TyjnLOubvdtTBW6kQd41tMjhmcanWbiComrzKQBZskT5HzmnL7-UD7kmMBKeD53DXA79FA3qd0UoSACJbTcCkjcvREGR9mItFd5U5Xp1aiEj2ZXMHUI_eGaFZHDjzxd9wgHLzooGzhfxUGGZ3zs2XG3h4-50fdbQGkXw216McOQAp5IvFDjA91N5-cye63F0WxbhTSfRSK8Z84uP6HzBuhhn_dQo4tb3ehMfCMbShmS83uzqmowfLFhOlCPVcUbS9jUib3aznsTSq-S48Dw8BgNUD4LgXRZFRKsLixPlCG2IR4ygf7jsHnOCVjCDnMsywswRR-J-nHGkaoFlpX89Yw==';
    const atoken = 'eyJ4NXQiOiJNVE15TkRSak1ETmlZemxrTjJVM1pUYzVZbVF3WW1WbFlqTTJaalJoWmpSaE4yTTVZak5oWmpZNVlqaGlNbVU0TXpZM1l6YzFNemd3TlRKaU1USTBNdyIsImtpZCI6Ik1UTXlORFJqTUROaVl6bGtOMlUzWlRjNVltUXdZbVZsWWpNMlpqUmhaalJoTjJNNVlqTmhaalk1WWpoaU1tVTRNelkzWXpjMU16Z3dOVEppTVRJME13X1JTMjU2IiwidHlwIjoiYXQrand0IiwiYWxnIjoiUlMyNTYifQ.eyJzdWIiOiIxZmYxMGFkZi02ZTU4LTRjYTktOTBiNy0yMTEyMjM5ZjU0YTMiLCJhdXQiOiJBUFBMSUNBVElPTiIsImF1ZCI6Im8zOUdmMDFwZktROUFtbTZXMjJ2cjFSZXZ4Z2EiLCJuYmYiOjE3NDIzNTc4MzcsImF6cCI6Im8zOUdmMDFwZktROUFtbTZXMjJ2cjFSZXZ4Z2EiLCJzY29wZSI6ImRlZmF1bHQiLCJpc3MiOiJodHRwczovL2FtLndzbzIuY29tOjQ0My9vYXV0aDIvdG9rZW4iLCJleHAiOjM2MTc0MjM1NzgzNywiaWF0IjoxNzQyMzU3ODM3LCJqdGkiOiJmYTRhNGM1Mi0wMDM1LTQyMmMtYjg5Mi1iZmY0YzEzYTQxMzYiLCJjbGllbnRfaWQiOiJvMzlHZjAxcGZLUTlBbW02VzIydnIxUmV2eGdhIn0.NJdqk3dCe2r6wcUYvLFXJ0N2vu0g7QWLXs5pjKqbC6HNUpqlSGGQOHWnZmjUtDZHS1mCTeNRYBLOzR3sDFgYm2LgGz9BJD_ulZ7mDFY2jhRnZTuzE1qmJ7O9QhtwKbY7O2We86ZHxWRaotm09wRQaRU08vSqtilDgB8ICLaGhgHnU2F4LI4kJok54u7uLHdVPya9Ffp1celGk7B1VkZYQb664xRb1gvMSFxyiTf5yV2hXzF93pdhOmYMYkH6sCVgYDROIE8PofDOYRDJMaOxiO9J9sC3uwqwoxctIk0Xa1hRiujSGmUrQ0TkEgohN8cIluqgePhzryBNUWwExK2GqA'
    const requestBody = {
        model: "gpt-4o",
        store: true,
        messages: [
            {
                role: "system",
                content: "You are a Stock market guru. You will not tell about the inaccebility of the realtime data. you will confidently suggest some tickers. Whenever you are mentioning about ticker you will wrap it within -- --. you will answer in less or equal to three sentence."
            }
        ]
    };
    questionsAndAnswers.forEach((qa, index) => {
        requestBody.messages.push({
            role: index % 2 === 0 ? "user" : "assistant",
            content: qa
        });
    });
    
    console.log("Request body:", requestBody);
    try {
        const response = await fetch(url, {
            method: 'POST',
            headers: {
                'accept': 'application/json',
                'Content-Type': 'application/json',
                // 'ApiKey': apiKey,
                'Authorization': `Bearer ${atoken}`
            },
            body: JSON.stringify(requestBody)
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();
        console.log("Response ::::: ",data);
        return data;
    } catch (error) {
        console.error('Error:', error);
        throw error;
    }
}

const questionsAndAnswers = []
const answers = []

// Example usage:
// getStockRecommendation().then(response => console.log(response));
// Basic send message functionality with chat and dummy response
async function sendMessage() {
    const input = document.getElementById('chatInput');
    const chatMessages = document.getElementById('chatMessages');
    if (input.value.trim()) {
        console.log("Question sent:", input.value);
        questionsAndAnswers.push(input.value);
        const welcomeMessage = document.getElementById("welcomeMessage");
        if (welcomeMessage) {
            welcomeMessage.remove();
        }
        // Add question to chat area
        const question = document.createElement('div');
        question.className = 'chat-message q';
        question.textContent = `Q: ${input.value}`;
        chatMessages.appendChild(question);
        addLoading();
        // Add dummy answer to chat area
        const answer = document.createElement('div');
        answer.className = 'chat-message a';

        // Scroll to the latest message if chat area exceeds half screen height
        if (chatMessages.scrollHeight > window.innerHeight * 0.5) {
            chatMessages.scrollTop = chatMessages.scrollHeight;
        }

        await sendQuestion().then(response => {
            console.log("Response ::::: ",response);
            if (response.choices && response.choices.length > 0) {
                questionsAndAnswers.push(response.choices[0].message.content);
                answer.textContent = `A: ${response.choices[0].message.content}`;
            } else {
                answer.textContent = `A: No valid response received.`;
            }
            chatMessages.appendChild(answer);
            chatMessages.appendChild(answer);
        })
        
        

        // Scroll to the latest message if chat area exceeds half screen height
        if (chatMessages.scrollHeight > window.innerHeight * 0.5) {
            chatMessages.scrollTop = chatMessages.scrollHeight;
        }

        input.value = ''; // Clear the input after sending

        // Generate and display dummy response table
        const responseTableBody = document.getElementById('responseTableBody');
        responseTableBody.innerHTML = ''; // Clear previous data
        // Extract symbols wrapped in -- -- from the answer text
        const symbolPattern = /--(.*?)--/g;
        const symbols = [];
        let match;
        while ((match = symbolPattern.exec(answer.textContent)) !== null) {
            symbols.push(match[1]);
        }
        console.log("Extracted symbols:", symbols);

        botSymbolData = []
        removeLoading();

        const promises = symbols.map(async symbol => {
            try {
                const data = await getStockQuote(symbol, accessToken);
                const result = data ? extractStockData(data) : null;
                botSymbolData.push(result);
                console.log("botSymbolData", botSymbolData);
                return result;
            } catch (error) {
                console.log(`Failed to fetch data for ${symbol}`);
                return null;
            }
        });
        
        await Promise.all(promises);
        console.log("****botSymbolData", botSymbolData);

        botSymbolData.forEach(row => {
            const tr = document.createElement('tr');
            tr.innerHTML = `
                <td>${row.symbol}</td>
                <td>${row.price}</td>
                <td style="color: ${row.change < 0 ? 'red' : 'green'}">${row.change}</td>
                <td>${row.volume}</td>
            `;
            responseTableBody.appendChild(tr);
        });

        // Show the response table
        document.getElementById('responseTable').style.display = 'block';
    }
}

// Allow sending message with Enter key
document.getElementById('chatInput').addEventListener('keypress', function(e) {
    if (e.key === 'Enter') {
        sendMessage();
    }
});