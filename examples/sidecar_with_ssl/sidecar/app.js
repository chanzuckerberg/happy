import express from 'express';
import fetch from 'node-fetch'
import https from 'https'
import fs from 'fs'

const app = express()
const port = 8443

app.get('/', (req, res) => {
  res.send('Hello World!');
});

app.get('/stacklist', async (req, res) => {
  const result = await (await fetch(process.env.API_URL)).json();
  res.send(result);
});

const httpsServer = https.createServer({
  key: fs.readFileSync(process.env.KEY_FILE),
  cert: fs.readFileSync(process.env.CERT_FILE),
  ca: fs.readFileSync(process.env.CA_CERT_FILE),
}, app);

httpsServer.listen(port, () => {
  console.log(`HTTP Server running on port ${port}`);
});
