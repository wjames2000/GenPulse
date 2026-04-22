import React from 'react';
import { cn } from '../../utils';

type StatusColor = 'healthy' | 'warning' | 'error' | 'idle';
type StatusSize = 'sm' | 'md' | 'lg';

interface StatusDotProps {
  status: StatusColor;
  size?: StatusSize;
  label?: string;
  className?: string;
  animate?: boolean;
}

const colorMap: Record<StatusColor, string> = {
  healthy: 'bg-primary',
  warning: 'bg-yellow-500',
  error: 'bg-red-500',
  idle: 'bg-white/20',
};

const pulseMap: Record<StatusColor, string> = {
  healthy: 'shadow-[0_0_8px_rgba(251,223,36,0.5)]',
  warning: 'shadow-[0_0_8px_rgba(234,179,8,0.5)]',
  error: 'shadow-[0_0_8px_rgba(239,68,68,0.5)]',
  idle: '',
};

const sizeMap: Record<StatusSize, string> = {
  sm: 'w-1.5 h-1.5',
  md: 'w-2 h-2',
  lg: 'w-3 h-3',
};

export function StatusDot({
  status,
  size = 'md',
  label,
  className,
  animate = true,
}: StatusDotProps) {
  return (
    <div className={cn('flex items-center gap-2', className)}>
      <span className="relative inline-flex shrink-0">
        <span
          className={cn(
            'rounded-full',
            colorMap[status],
            sizeMap[size],
            animate && status !== 'idle' && 'animate-pulse',
            animate && status !== 'idle' && pulseMap[status],
          )}
        />
        {animate && status === 'healthy' && (
          <span
            className={cn(
              'absolute inset-0 rounded-full pulse-dot',
              colorMap[status],
            )}
          />
        )}
      </span>
      {label && (
        <span className="text-[9px] font-black uppercase tracking-widest text-white/40">
          {label}
        </span>
      )}
    </div>
  );
}
