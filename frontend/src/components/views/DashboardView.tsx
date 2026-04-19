import React from 'react';
import { useAppStore } from '../../stores/appStore';

const DashboardView: React.FC = () => {
  const {
    appInfo,
    healthStatus,
    isInitialized,
    isExecuting,
    executionProgress,
    executionLogs,
    agentStatus,
    startExecution,
    stopExecution,
    addExecutionLog
  } = useAppStore();

  const handleStartExecution = () => {
    startExecution();
    addExecutionLog('info', 'Execution started');
    
    // 模拟执行进度
    let progress = 0;
    const interval = setInterval(() => {
      progress += 10;
      useAppStore.getState().setExecutionProgress(progress);
      
      if (progress >= 100) {
        clearInterval(interval);
        useAppStore.getState().stopExecution();
        addExecutionLog('success', 'Execution completed successfully');
      }
    }, 500);
  };

  const handleStopExecution = () => {
    stopExecution();
    addExecutionLog('warning', 'Execution stopped by user');
  };

  return (
    <div className="dashboard">
      <div className="dashboard-header">
        <h2>Dashboard</h2>
        <div className="dashboard-actions">
          <button 
            className={`btn ${isExecuting ? 'btn-stop' : 'btn-start'}`}
            onClick={isExecuting ? handleStopExecution : handleStartExecution}
            disabled={!isInitialized}
          >
            {isExecuting ? '⏹️ Stop' : '▶️ Start'} Execution
          </button>
        </div>
      </div>

      {/* 状态卡片 */}
      <div className="status-cards">
        <div className="status-card">
          <h3>Application Status</h3>
          {appInfo ? (
            <div className="status-card-content">
              <p><strong>Name:</strong> {appInfo.name}</p>
              <p><strong>Version:</strong> {appInfo.version}</p>
              <p><strong>Status:</strong> 
                <span className={`status-badge status-${appInfo.status}`}>
                  {appInfo.status}
                </span>
              </p>
            </div>
          ) : (
            <p>Loading app info...</p>
          )}
        </div>

        <div className="status-card">
          <h3>Service Health</h3>
          {healthStatus ? (
            <div className="status-card-content">
              <p><strong>Status:</strong> 
                <span className={`status-badge health-${healthStatus.status}`}>
                  {healthStatus.status}
                </span>
              </p>
              <p><strong>Service:</strong> {healthStatus.service}</p>
              <p><strong>Last Check:</strong> 
                {new Date(healthStatus.timestamp * 1000).toLocaleTimeString()}
              </p>
            </div>
          ) : (
            <p>Checking health...</p>
          )}
        </div>

        <div className="status-card">
          <h3>Execution Progress</h3>
          <div className="status-card-content">
            <div className="progress-container">
              <div 
                className="progress-bar" 
                style={{ width: `${executionProgress}%` }}
              >
                {executionProgress}%
              </div>
            </div>
            <p><strong>Status:</strong> {isExecuting ? 'Running' : 'Idle'}</p>
            <p><strong>Progress:</strong> {executionProgress}%</p>
          </div>
        </div>

        <div className="status-card">
          <h3>System Info</h3>
          <div className="status-card-content">
            <p><strong>Initialization:</strong> 
              <span className={isInitialized ? 'text-success' : 'text-warning'}>
                {isInitialized ? '✅ Complete' : '⏳ In Progress'}
              </span>
            </p>
            <p><strong>Active Agents:</strong> {Object.keys(agentStatus).length}</p>
            <p><strong>Log Entries:</strong> {executionLogs.length}</p>
          </div>
        </div>
      </div>

      {/* 执行日志 */}
      <div className="execution-logs">
        <h3>Execution Logs</h3>
        <div className="logs-container">
          {executionLogs.length === 0 ? (
            <p className="no-logs">No logs yet. Start an execution to see logs.</p>
          ) : (
            <div className="logs-list">
              {executionLogs.slice().reverse().map((log) => (
                <div key={log.id} className={`log-entry log-${log.level}`}>
                  <span className="log-time">
                    {new Date(log.timestamp).toLocaleTimeString()}
                  </span>
                  <span className="log-level">[{log.level}]</span>
                  <span className="log-message">{log.message}</span>
                </div>
              ))}
            </div>
          )}
        </div>
        {executionLogs.length > 0 && (
          <button 
            className="btn btn-clear"
            onClick={() => useAppStore.getState().clearExecutionLogs()}
          >
            Clear Logs
          </button>
        )}
      </div>

      {/* 快速操作 */}
      <div className="quick-actions">
        <h3>Quick Actions</h3>
        <div className="action-buttons">
          <button className="btn btn-action" onClick={() => {/* TODO */}}>
            📁 New Project
          </button>
          <button className="btn btn-action" onClick={() => {/* TODO */}}>
            🤖 Configure Agents
          </button>
          <button className="btn btn-action" onClick={() => {/* TODO */}}>
            ⚙️ Settings
          </button>
          <button className="btn btn-action" onClick={() => {/* TODO */}}>
            📊 View Metrics
          </button>
        </div>
      </div>
    </div>
  );
};

export default DashboardView;