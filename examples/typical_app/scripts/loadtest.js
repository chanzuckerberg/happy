import http from 'k6/http';
import { sleep } from 'k6';

export default function () {
  http.get('https://alokshin.k8s-test-app.rdev.si.czi.technology/proxy/?count=5000');
  sleep(1);
}