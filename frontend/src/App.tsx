import {useEffect} from 'react';
import './App.css';
import {Greet} from "../wailsjs/go/main/App";
import baseService from './services/baseService';
import { useAppStore, selectCurrentView } from './stores/appStore';
import LayoutGenpulse from './components/LayoutGenpulse';
import DashboardViewGenpulseThree from './components/views/DashboardViewGenpulseThree';
import AgentViewNew from './components/views/AgentViewNew';
import ProjectsViewGenpulseThree from './components/views/ProjectsViewGenpulseThree';
import SkillsViewGenpulseThree from './components/views/SkillsViewGenpulseThree';
import MemoryViewGenpulseThree from './components/views/MemoryViewGenpulseThree';
import KanbanViewGenpulseThree from './components/views/KanbanViewGenpulseThree';
import TerminalViewGenpulseThree from './components/views/TerminalViewGenpulseThree';
import AgentViewGenpulseThree from './components/views/AgentViewGenpulseThree';
import SettingsViewDesignThree from './components/views/SettingsViewDesignThree';

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

  const currentView = useAppStore(selectCurrentView);

  const renderView = () => {
    switch (currentView) {
      case 'dashboard':
        return <DashboardViewGenpulseThree />;
      case 'projects':
        return <ProjectsViewGenpulseThree />;
      case 'agents':
        return <AgentViewGenpulseThree />;
      case 'skills':
        return <SkillsViewGenpulseThree />;
      case 'memory':
        return <MemoryViewGenpulseThree />;
      case 'kanban':
        return <KanbanViewGenpulseThree />;
      case 'terminal':
        return <TerminalViewGenpulseThree />;
      case 'settings':
        return <SettingsViewDesignThree />;
      default:
        return <DashboardViewGenpulseThree />;
    }
  };

  return (
    <LayoutGenpulse>
      {renderView()}
    </LayoutGenpulse>
  );
}

export default App
