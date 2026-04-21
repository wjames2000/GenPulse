import { api } from './services/api';

async function testApi() {
  console.log('Testing API integration...\n');

  try {
    // 1. Test app info
    console.log('1. Testing GetAppInfo...');
    const appInfo = await api.getAppInfo();
    console.log('App Info:', appInfo ? 'Success' : 'Failed');

    // 2. Test health check
    console.log('\n2. Testing HealthCheck...');
    const health = await api.healthCheck();
    console.log('Health Check:', health);

    // 3. Test list agents
    console.log('\n3. Testing ListAgents...');
    const agents = await api.listAgents();
    console.log('Agents:', agents.length > 0 ? `Found ${agents.length} agents` : 'No agents found');

    // 4. Test get all agents status
    console.log('\n4. Testing GetAllAgentsStatus...');
    const agentsStatus = await api.getAllAgentsStatus();
    console.log('Agents Status:', Object.keys(agentsStatus).length > 0 ? 'Success' : 'No status available');

    // 5. Test log message
    console.log('\n5. Testing LogMessage...');
    await api.logMessage('info', 'Test log message from frontend');
    console.log('Log message sent');

    // 6. Test get logs
    console.log('\n6. Testing GetLogs...');
    const logs = await api.getLogs();
    console.log('Logs:', Array.isArray(logs) ? `Found ${logs.length} logs` : 'No logs available');

    console.log('\n✅ All API tests completed!');
  } catch (error) {
    console.error('\n❌ API test failed:', error);
  }
}

// Run test if this file is executed directly
if (typeof window !== 'undefined') {
  console.log('Running API tests in browser...');
  testApi();
}

export { testApi };