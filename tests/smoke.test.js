const { test } = require("node:test");
const assert = require("node:assert/strict");

const BASE = process.env.BASE || "http://localhost:8086";

// 1. Tes Kesehatan Server
test("GET /health membalas ok", async () => {
  const res = await fetch(`${BASE}/health`);
  assert.equal(res.status, 200);
});

// 2. Tes Login Pengguna yang Sah
test("POST /auth/login yang sah → 200", async () => {
  const res = await fetch(`${BASE}/auth/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      email: "ferdi@gmail.com",
      password: "user123"
    }),
  });
  assert.equal(res.status, 200);
});

// 3. Tes Validasi Error (Login tanpa password merespons 401)
test("POST /auth/login tanpa password → 401", async () => {
  const res = await fetch(`${BASE}/auth/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      email: "ferdi@gmail.com"
    }),
  });
  assert.equal(res.status, 401);
});