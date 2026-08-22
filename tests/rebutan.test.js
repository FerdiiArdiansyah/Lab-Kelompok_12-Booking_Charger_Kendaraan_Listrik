const { test } = require("node:test");
const assert = require("node:assert/strict");

const BASE = process.env.BASE || "http://localhost:8086";
const PENYERBU = 300;

async function serbuLogin() {
  const res = await fetch(`${BASE}/auth/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email: "ferdi@gmail.com", password: "user123" }),
  });
  return res.status;
}

test(`Uji rebutan login (${PENYERBU} koneksi serentak)`, async () => {
  const status = await Promise.all(Array.from({ length: PENYERBU }, serbuLogin));
  const sukses = status.filter((s) => s === 200).length;
  const gagal = status.filter((s) => s !== 200).length;

  console.log(`Sukses (200): ${sukses} | Gagal/Lainnya: ${gagal}`);
  assert.equal(sukses + gagal, PENYERBU, "Setiap permintaan harus mendapat respons");
});