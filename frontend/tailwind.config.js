/** @type {import('tailwindcss').Config} */
export default {
  darkMode: 'class',
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        // 表面颜色
        'surface': '#131318',
        'surface-dim': '#131318',
        'surface-bright': '#39383e',
        'surface-container-lowest': '#0e0e13',
        'surface-container-low': '#1b1b20',
        'surface-container': '#1f1f25',
        'surface-container-high': '#2a292f',
        'surface-container-highest': '#35343a',
        'surface-variant': '#35343a',
        
        // 主色
        'primary': '#c0c1ff',
        'primary-container': '#5b5fff',
        'primary-fixed': '#e1e0ff',
        'primary-fixed-dim': '#c0c1ff',
        'inverse-primary': '#4345e8',
        'surface-tint': '#c0c1ff',
        
        // 次要色
        'secondary': '#c0c1ff',
        'secondary-container': '#3f4287',
        'secondary-fixed': '#e1e0ff',
        'secondary-fixed-dim': '#c0c1ff',
        
        // 第三色
        'tertiary': '#ffb68f',
        'tertiary-container': '#c05500',
        'tertiary-fixed': '#ffdbca',
        'tertiary-fixed-dim': '#ffb68f',
        
        // 错误色
        'error': '#ffb4ab',
        'error-container': '#93000a',
        
        // 轮廓色
        'outline': '#908fa1',
        'outline-variant': '#454555',
        
        // 文字颜色
        'on-surface': '#e4e1e9',
        'on-surface-variant': '#c6c4d8',
        'on-background': '#e4e1e9',
        'on-primary': '#0d00aa',
        'on-primary-container': '#fffcff',
        'on-primary-fixed': '#05006c',
        'on-primary-fixed-variant': '#2623d1',
        'on-secondary': '#26286c',
        'on-secondary-container': '#afb2ff',
        'on-secondary-fixed': '#0e0f58',
        'on-secondary-fixed-variant': '#3d3f84',
        'on-tertiary': '#542100',
        'on-tertiary-container': '#fffcff',
        'on-tertiary-fixed': '#331100',
        'on-tertiary-fixed-variant': '#773200',
        'on-error': '#690005',
        'on-error-container': '#ffdad6',
        'inverse-surface': '#e4e1e9',
        'inverse-on-surface': '#303036',
        
        // 背景色
        'background': '#131318',
        
        // 系统颜色
        'sys-success': '#4ADE80',
        
        // 功能色（来自设计规范 7.1.3）
        'success': '#4ADE80',
        'warning': '#FBBF24',
        'error': '#F87171',
        'info': '#60A5FA',
        'running': '#3B82F6',
        'pending': '#A1A1AA',
        
        // Agent 角色专属色（来自设计规范 7.1.4）
        'agent-orchestrator': '#A78BFA',
        'agent-product-manager': '#2DD4BF',
        'agent-architect': '#FB923C',
        'agent-frontend': '#38BDF8',
        'agent-backend': '#4ADE80',
        'agent-qa': '#F472B6',
        'agent-devops': '#FBBF24',
        'agent-reviewer': '#818CF8',
      },
      borderRadius: {
        'none': '0',
        'sm': '0.125rem',
        'DEFAULT': '0.125rem',
        'md': '0.25rem',
        'lg': '0.25rem',
        'xl': '0.5rem',
        '2xl': '0.75rem',
        '3xl': '1rem',
        'full': '0.75rem',
        
        // 组件特定圆角
        'card': '0.5rem',
        'card-lg': '0.75rem',
        'button': '0.25rem',
        'button-lg': '0.5rem',
        'input': '0.25rem',
        'input-lg': '0.5rem',
        'badge': '0.75rem',
        'avatar': '0.75rem',
      },
      fontFamily: {
        // 系统界面字体（来自设计规范 7.2.1）
        'sans': ['"Inter"', '-apple-system', 'BlinkMacSystemFont', '"Segoe UI"', 'sans-serif'],
        'sans-chinese': ['"PingFang SC"', '"Microsoft YaHei"', '"Source Han Sans CN"', 'sans-serif'],
        'headline': ['Inter', 'sans-serif'],
        'body': ['Inter', 'sans-serif'],
        'label': ['Inter', 'sans-serif'],
        // 代码/终端字体
        'mono': ['"JetBrains Mono"', '"Fira Code"', '"Cascadia Code"', '"Source Code Pro"', 'monospace']
      },
      fontSize: {
        'display-large': ['3.5625rem', { lineHeight: '4rem', fontWeight: '400' }],
        'display-medium': ['2.8125rem', { lineHeight: '3.25rem', fontWeight: '400' }],
        'display-small': ['2.25rem', { lineHeight: '2.75rem', fontWeight: '400' }],
        
        'headline-large': ['2rem', { lineHeight: '2.5rem', fontWeight: '400' }],
        'headline-medium': ['1.75rem', { lineHeight: '2.25rem', fontWeight: '400' }],
        'headline-small': ['1.5rem', { lineHeight: '2rem', fontWeight: '400' }],
        
        'title-large': ['1.375rem', { lineHeight: '1.75rem', fontWeight: '400' }],
        'title-medium': ['1rem', { lineHeight: '1.5rem', fontWeight: '500' }],
        'title-small': ['0.875rem', { lineHeight: '1.25rem', fontWeight: '500' }],
        
        'label-large': ['0.875rem', { lineHeight: '1.25rem', fontWeight: '500' }],
        'label-medium': ['0.75rem', { lineHeight: '1rem', fontWeight: '500' }],
        'label-small': ['0.6875rem', { lineHeight: '1rem', fontWeight: '500' }],
        
        'body-large': ['1rem', { lineHeight: '1.5rem', fontWeight: '400' }],
        'body-medium': ['0.875rem', { lineHeight: '1.25rem', fontWeight: '400' }],
        'body-small': ['0.75rem', { lineHeight: '1rem', fontWeight: '400' }],
      },
      fontWeight: {
        'regular': '400',
        'medium': '500',
        'semibold': '600',
        'bold': '700',
        'black': '900',
      },
      letterSpacing: {
        'tight': '-0.02em',
        'normal': '0',
        'wide': '0.05em',
      },
      spacing: {
        // 8px基础网格系统（来自设计规范 4.2）
        'xs': '4px',
        'sm': '8px',
        'md': '16px',
        'lg': '24px',
        'xl': '32px',
        '2xl': '48px',
      },
       keyframes: {
        slideUpFade: {
          'from': { opacity: '0', transform: 'translateY(10px)' },
          'to': { opacity: '1', transform: 'translateY(0)' }
        },
        pulse: {
          '0%, 100%': { opacity: '1' },
          '50%': { opacity: '0.5' }
        },
        fadeIn: {
          'from': { opacity: '0' },
          'to': { opacity: '1' }
        },
        slideInRight: {
          'from': { transform: 'translateX(100%)' },
          'to': { transform: 'translateX(0)' }
        },
        slideInLeft: {
          'from': { transform: 'translateX(-100%)' },
          'to': { transform: 'translateX(0)' }
        },
        scaleIn: {
          'from': { transform: 'scale(0.95)', opacity: '0' },
          'to': { transform: 'scale(1)', opacity: '1' }
        },
        // Agent 运行脉冲动画
        'pulse-border': {
          '0%, 100%': { 'box-shadow': '0 0 0 0 rgba(59, 130, 246, 0.7)' },
          '50%': { 'box-shadow': '0 0 0 4px rgba(59, 130, 246, 0)' }
        }
      },
      animation: {
        'slideUpFade': 'slideUpFade 0.5s ease forwards',
        'pulse': 'pulse 2s infinite',
        'fadeIn': 'fadeIn 0.3s ease-in',
        'slideInRight': 'slideInRight 0.3s ease-out',
        'slideInLeft': 'slideInLeft 0.3s ease-out',
        'scaleIn': 'scaleIn 0.2s ease-out',
        // 动效令牌（来自设计规范 7.4.3）
        'instant': '100ms ease',
        'fast': '150ms ease-out',
        'base': '200ms ease-out',
        'slow': '300ms ease-out',
        'pulse-border': 'pulse-border 1.5s ease-in-out infinite',
      },
       boxShadow: {
        // 阴影定义（深色模式）- 来自设计规范 7.5.1
        'sm': '0 1px 2px 0 rgba(0,0,0,0.4)',
        'md': '0 4px 6px -2px rgba(0,0,0,0.4), 0 2px 4px -1px rgba(0,0,0,0.3)',
        'lg': '0 10px 15px -3px rgba(0,0,0,0.5), 0 4px 6px -2px rgba(0,0,0,0.3)',
        'xl': '0 20px 25px -5px rgba(0,0,0,0.5), 0 8px 10px -3px rgba(0,0,0,0.4)',
        '2xl': '0 25px 50px -12px rgba(0, 0, 0, 0.6)',
        'inner': 'inset 0 2px 4px 0 rgba(0,0,0,0.5)',
        // 液态玻璃效果（Liquid Glass）- 来自设计规范 7.5.2
        'glass': '0 8px 32px rgba(0, 0, 0, 0.4), inset 0 1px 0 rgba(255, 255, 255, 0.1)',
        // 旧版兼容
        'card': '0 4px 24px rgba(0, 0, 0, 0.4), 0 1px 0 rgba(255, 255, 255, 0.05)',
        'float': '0 20px 40px rgba(0, 0, 0, 0.5), 0 1px 0 rgba(255, 255, 255, 0.1)',
      }
    }
  },
  plugins: [],
}