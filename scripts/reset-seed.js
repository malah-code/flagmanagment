const { request, seedData } = require('./seed');

async function resetAndSeed() {
  console.log('🧹 Starting Reseed (Cleanup + Fresh Seed)...\n');

  // 1. Authenticate
  const loginRes = await request('/auth/login', 'POST', {
    email: 'admin@example.com',
    password: 'admin123'
  });

  if (loginRes.status !== 200 || !loginRes.data.token) {
    console.error('❌ Authentication failed. Is FlagManagment backend running?', loginRes.status);
    process.exit(1);
  }

  const token = loginRes.data.token;
  console.log('✅ Authenticated as admin@example.com');

  // 2. Fetch existing projects
  const projectsRes = await request('/projects', 'GET', null, token);
  const projects = projectsRes.data?.data || [];

  console.log(`\n🗑️  Found ${projects.length} existing project(s). Deleting to start fresh...`);

  for (const proj of projects) {
    const delRes = await request(`/projects/${proj.id}`, 'DELETE', null, token);
    if (delRes.status === 204 || delRes.status === 200) {
      console.log(`   ❌ Deleted project: "${proj.name}" (${proj.id})`);
    } else {
      console.warn(`   ⚠️ Failed to delete project "${proj.name}":`, delRes.status, delRes.raw);
    }
  }

  console.log('\n✨ Database cleaned up cleanly!\n');

  // 3. Seed fresh data
  await seedData(token);
}

resetAndSeed();
