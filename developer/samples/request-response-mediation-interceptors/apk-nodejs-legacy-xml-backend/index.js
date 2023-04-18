const http = require('http');
const url = require('url');

const server = http.createServer((req, res) => {
    const reqUrl = url.parse(req.url, true);

    if (reqUrl.pathname === '/books' && req.method === 'POST') {
        console.log('Backend service is called');

        const xUserHeader = req.headers['x-user'];
        if (xUserHeader !== 'admin') {
            res.setHeader('Content-Type', 'application/xml');
            res.statusCode = 401;
            res.end('<response>Error</response>');
            return;
        }

        let body = '';
        req.on('data', chunk => {
            body += chunk;
        });

        req.on('end', () => {
            console.log(`Received payload ${body}`);
            res.setHeader('Content-Type', 'text/plain');
            res.statusCode = 200;
            res.end('created');
        });
    } else {
        res.statusCode = 404;
        res.end();
    }
});

server.listen(9082, () => {
    console.log('Server listening on port 9082');
});
