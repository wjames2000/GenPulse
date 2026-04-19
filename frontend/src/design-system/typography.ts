// Genpulse 设计系统字体配置
// 基于设计规范

export const typography = {
  fontFamily: {
    headline: ['Inter', 'sans-serif'],
    body: ['Inter', 'sans-serif'],
    label: ['Inter', 'sans-serif'],
    mono: ['JetBrains Mono', 'monospace'],
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
    regular: '400',
    medium: '500',
    semibold: '600',
    bold: '700',
    black: '900',
  },
  
  letterSpacing: {
    tight: '-0.02em',
    normal: '0',
    wide: '0.05em',
  },
} as const;