// Supabase Configuration
// Update SUPABASE_URL to your project URL after creating the project.
const SUPABASE_URL = localStorage.getItem('supabase_url') || 'https://rmekfsdhglyiralxvkwc.supabase.co';
const SUPABASE_ANON_KEY = localStorage.getItem('supabase_anon_key') || '';

async function supabaseFetch(endpoint, options = {}) {
  const url = `${SUPABASE_URL}/rest/v1/${endpoint}`;
  const headers = {
    'apikey': SUPABASE_ANON_KEY,
    'Authorization': `Bearer ${SUPABASE_ANON_KEY}`,
    'Content-Type': 'application/json',
    'Prefer': 'return=representation',
    ...options.headers
  };
  const res = await fetch(url, { ...options, headers });
  if (!res.ok) throw new Error(`Supabase ${endpoint}: ${res.status} ${res.statusText}`);
  return res.json();
}

async function supabaseUpsert(table, rows) {
  return supabaseFetch(table, {
    method: 'POST',
    headers: { 'Prefer': 'resolution=merge-duplicates' },
    body: JSON.stringify(rows)
  });
}

async function supabaseSelect(table, query = '') {
  return supabaseFetch(`${table}${query ? '?' + query : ''}`);
}

async function supabaseDelete(table, query) {
  return supabaseFetch(`${table}?${query}`, { method: 'DELETE' });
}

async function supabaseRpc(procedure, params = {}) {
  const url = `${SUPABASE_URL}/rest/v1/rpc/${procedure}`;
  const headers = {
    'apikey': SUPABASE_ANON_KEY,
    'Authorization': `Bearer ${SUPABASE_ANON_KEY}`,
    'Content-Type': 'application/json'
  };
  const res = await fetch(url, { method: 'POST', headers, body: JSON.stringify(params) });
  if (!res.ok) throw new Error(`RPC ${procedure}: ${res.status}`);
  return res.json();
}