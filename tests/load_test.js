import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate } from 'k6/metrics';

const errorRate = new Rate('error_rate');

export const options = {
  stages: [
    { duration: '30s', target: 100 },  // Ramp up to 100 users
    { duration: '1m', target: 100 },   // Stay at 100 users
    { duration: '30s', target: 200 },  // Ramp up to 200 users
    { duration: '1m', target: 200 },   // Stay at 200 users
    { duration: '30s', target: 0 },    // Ramp down
  ],
};

export default function () {
  const baseUrl = __ENV.BASE_URL || 'http://localhost:8080';

  // Health check
  const healthRes = http.get(`${baseUrl}/health`);
  check(healthRes, {
    'health check is 200': (r) => r.status === 200,
  });
  errorRate.add(healthRes.status !== 200);

  // Register
  const registerRes = http.post(`${baseUrl}/api/v1/auth/register`, JSON.stringify({
    email: `test_${__VU}_${__IterationInTest}@example.com`,
    password: 'password123',
    name: `Test User ${__VU}`,
  }), {
    headers: { 'Content-Type': 'application/json' },
  });
  check(registerRes, {
    'register is 201 or 409': (r) => r.status === 201 || r.status === 409,
  });
  errorRate.add(registerRes.status >= 500);

  // Login
  const loginRes = http.post(`${baseUrl}/api/v1/auth/login`, JSON.stringify({
    email: 'test@example.com',
    password: 'password123',
  }), {
    headers: { 'Content-Type': 'application/json' },
  });
  check(loginRes, {
    'login is 200': (r) => r.status === 200,
  });
  errorRate.add(loginRes.status !== 200);

  const token = loginRes.json('token');

  if (token) {
    const headers = { Authorization: `Bearer ${token}` };

    // Get profile
    const meRes = http.get(`${baseUrl}/api/v1/users/me`, { headers });
    check(meRes, {
      'get me is 200': (r) => r.status === 200,
    });
    errorRate.add(meRes.status !== 200);

    // Create item
    const createRes = http.post(`${baseUrl}/api/v1/items`, JSON.stringify({
      name: 'Test Item',
      description: 'A test item for load testing',
    }), {
      headers: { ...headers, 'Content-Type': 'application/json' },
    });
    check(createRes, {
      'create item is 201': (r) => r.status === 201,
    });
    errorRate.add(createRes.status !== 201);

    // List items
    const listRes = http.get(`${baseUrl}/api/v1/items`, { headers });
    check(listRes, {
      'list items is 200': (r) => r.status === 200,
    });
    errorRate.add(listRes.status !== 200);
  }

  sleep(1);
}
