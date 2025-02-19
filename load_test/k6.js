import http from 'k6/http';
import { sleep } from 'k6';

function randomFloat(min, max) {
    return Math.random() * (max - min) + min;
  }

export let options = {
    scenarios: {
        load_increase: {
            executor: 'ramping-arrival-rate',
            startRate: 100,
            timeUnit: '1m',
            preAllocatedVUs: 50,
            maxVUs: 10000,
            stages: [
                { duration: '2m', target: 4000 },
                { duration: '2m', target: 8000 },
                { duration: '2m', target: 16000 },
            ],
            startTime: '2m',
        }
    }
};

export default function () {
    const urls = [
        'http://localhost:5001/predict',
        'http://localhost:5002/predict',
        'http://localhost:5003/predict'
    ];

    urls.forEach(url => {
        const payload = JSON.stringify({
            sepal_length: randomFloat(4.3, 7.9),
            sepal_width: randomFloat(2, 4.4),
            petal_length: randomFloat(1, 6.9),
            petal_width: randomFloat(0.1, 2.5),
        });

        const params = {
            headers: { 'Content-Type': 'application/json' }
        };

        http.post(url, payload, params);
    });

    sleep(1);
}

