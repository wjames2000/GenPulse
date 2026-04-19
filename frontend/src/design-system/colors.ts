// Genpulse 设计系统颜色配置
// 基于 Material Design 3 和设计规范

export const colors = {
  // 表面颜色
  surface: '#131318',
  'surface-dim': '#131318',
  'surface-bright': '#39383e',
  'surface-container-lowest': '#0e0e13',
  'surface-container-low': '#1b1b20',
  'surface-container': '#1f1f25',
  'surface-container-high': '#2a292f',
  'surface-container-highest': '#35343a',
  'surface-variant': '#35343a',
  
  // 主色
  primary: '#c0c1ff',
  'primary-container': '#5b5fff',
  'primary-fixed': '#e1e0ff',
  'primary-fixed-dim': '#c0c1ff',
  'inverse-primary': '#4345e8',
  'surface-tint': '#c0c1ff',
  
  // 次要色
  secondary: '#c0c1ff',
  'secondary-container': '#3f4287',
  'secondary-fixed': '#e1e0ff',
  'secondary-fixed-dim': '#c0c1ff',
  
  // 第三色
  tertiary: '#ffb68f',
  'tertiary-container': '#c05500',
  'tertiary-fixed': '#ffdbca',
  'tertiary-fixed-dim': '#ffb68f',
  
  // 错误色
  error: '#ffb4ab',
  'error-container': '#93000a',
  
  // 轮廓色
  outline: '#908fa1',
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
  background: '#131318',
  
  // 系统颜色
  'sys-success': '#4ADE80',
} as const;

export type ColorKey = keyof typeof colors;