import http from 'k6/http';

export const options = {
  scenarios: {
    contacts: {
      executor: 'ramping-arrival-rate',
      preAllocatedVUs: 1000,
      timeUnit: '1s',
      startRate: 500,
      stages: [
        { target: 800, duration: '5s' },
        { target: 1000, duration: '0' },
        { target: 1000, duration: '10s' },
      ],
    },
  },
};

export default function () {
  const res = http.get('http://localhost:8080/?search=GTA6');
}