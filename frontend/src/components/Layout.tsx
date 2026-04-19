import React from 'react';
import { useAppStore, selectSidebarOpen, selectDarkMode } from '../stores/appStore';
import Sidebar from './Sidebar';
import Header from './Header';

interface LayoutProps {
  children: React.ReactNode;
}

const Layout: React.FC<LayoutProps> = ({ children }) => {
  const sidebarOpen = useAppStore(selectSidebarOpen);
  const darkMode = useAppStore(selectDarkMode);

  return (
    <div className={`layout ${darkMode ? 'dark' : 'light'}`}>
      <Header />
      <div className="layout-content">
        {sidebarOpen && <Sidebar />}
        <main className={`main-content ${sidebarOpen ? 'with-sidebar' : 'full-width'}`}>
          {children}
        </main>
      </div>
    </div>
  );
};

export default Layout;