import React from 'react';
import { motion } from 'motion/react';
import { cn } from '../../utils';
import { LucideIcon } from 'lucide-react';

type EmptyStateSize = 'sm' | 'md' | 'lg';

interface EmptyStateProps {
  icon: LucideIcon;
  title: string;
  description?: string;
  action?: {
    label: string;
    onClick: () => void;
  };
  size?: EmptyStateSize;
  className?: string;
}

const iconSizes: Record<EmptyStateSize, string> = {
  sm: 'w-8 h-8',
  md: 'w-12 h-12',
  lg: 'w-16 h-16',
};

const titleSizes: Record<EmptyStateSize, string> = {
  sm: 'text-sm',
  md: 'text-lg',
  lg: 'text-2xl',
};

const descSizes: Record<EmptyStateSize, string> = {
  sm: 'text-[10px]',
  md: 'text-xs',
  lg: 'text-sm',
};

export function EmptyState({
  icon: Icon,
  title,
  description,
  action,
  size = 'md',
  className,
}: EmptyStateProps) {
  return (
    <motion.div
      initial={{ opacity: 0, y: 12 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.3, ease: 'easeOut' }}
      className={cn(
        'flex flex-col items-center justify-center text-center py-16 px-8',
        className,
      )}
    >
      <div className={cn(
        'mb-6 text-white/20 flex items-center justify-center',
        iconSizes[size],
      )}>
        <Icon size="100%" strokeWidth={1} />
      </div>
      <h3 className={cn(
        'font-black uppercase tracking-tighter text-white/60 mb-2',
        titleSizes[size],
      )}>
        {title}
      </h3>
      {description && (
        <p className={cn(
          'text-white/40 max-w-md leading-relaxed font-medium',
          descSizes[size],
        )}>
          {description}
        </p>
      )}
      {action && (
        <button
          onClick={action.onClick}
          className="mt-8 px-8 py-3 bg-primary text-black font-black uppercase text-[10px] tracking-widest hover:scale-105 transition-all"
        >
          {action.label}
        </button>
      )}
    </motion.div>
  );
}
