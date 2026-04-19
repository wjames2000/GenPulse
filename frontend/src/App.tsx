import {useEffect} from 'react';
import './App.css';
import {Greet} from "../wailsjs/go/main/App";
import baseService from './services/baseService';
import { useAppStore } from './stores/appStore';
import Layout from './components/Layout';
import DashboardView from './components/views/DashboardView';

function App() {
  const {
    appInfo,
    healthStatus,
    isInitialized,
    setAppInfo,
    setHealthStatus,
    setInitialized,
    addExecutionLog
  } = useAppStore();

  // 初始化基础服务
  useEffect(() => {
    const init = async () => {
      try {
        const initialized = await baseService.initialize();
        setInitialized(initialized);
        
        if (initialized) {
          const info = await baseService.getAppInfo();
          const health = await baseService.healthCheck();
          setAppInfo(info);
          setHealthStatus(health);
          
          // 记录初始化日志
          await baseService.logMessage('info', 'Frontend application initialized');
          addExecutionLog('info', 'Application initialized successfully');
        } else {
          addExecutionLog('error', 'Failed to initialize base service');
        }
      } catch (error) {
        console.error('Initialization error:', error);
        addExecutionLog('error', `Initialization error: ${error}`);
        setInitialized(false);
      }
    };

    init();

    // 清理函数
    return () => {
      baseService.cleanupEventListeners();
    };
  }, [setAppInfo, setHealthStatus, setInitialized, addExecutionLog]);

  // 测试Greet函数
  const testGreet = async () => {
    try {
      const result = await Greet('GenPulse User');
      addExecutionLog('info', `Greet test: ${result}`);
    } catch (error) {
      addExecutionLog('error', `Greet test failed: ${error}`);
    }
  };

  // 在初始化后测试
  useEffect(() => {
    if (isInitialized) {
      testGreet();
    }
  }, [isInitialized]);

  return (
    <Layout>
      <DashboardView />
    </Layout>
  );
}

export default App
