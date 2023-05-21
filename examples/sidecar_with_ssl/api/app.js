import express from 'express';
import fs from 'fs';

const app = express()
const port = 3000

const stacklist = fs.readFileSync(process.env.STACK_LIST_FILE).toString();

app.get('/', (req, res) => {
  res.send('Hello World!');
});

app.get('/stacklist', (req, res) => {
  res.send(stacklist);
});


app.listen(port, () => {
  console.log(`Example app listening on port ${port}`);
});