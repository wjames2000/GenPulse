// Genpulse 设计系统圆角配置
// 基于设计规范

export const borderRadius = {
  // 基础圆角值
  none: '0',
  sm: '0.125rem', // 2px
  DEFAULT: '0.125rem', // 2px
  md: '0.25rem',  // 4px
  lg: '0.25rem',  // 4px
  xl: '0.5rem',   // 8px
  '2xl': '0.75rem', // 12px
  '3xl': '1rem',  // 16px
  full: '0.75rem', // 12px
  
  // 组件特定圆角
  card: '0.5rem',     // 8px
  'card-lg': '0.75rem', // 12px
  button: '0.25rem',  // 4px
  'button-lg': '0.5rem', // 8px
  input: '0.25rem',   // 4px
  'input-lg': '0.5rem', // 8px
  badge: '0.75rem',   // 12px
  avatar: '0.75rem',  // 12px
} as const;