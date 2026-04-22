import React from 'react';
import { motion, type Variants } from 'motion/react';
import { cn } from '../../utils';

export type AnimationVariant = 'fadeIn' | 'slideUp' | 'scaleIn';

const animationVariants: Record<AnimationVariant, Variants> = {
  fadeIn: {
    initial: { opacity: 0 },
    animate: { opacity: 1 },
    exit: { opacity: 0 },
  },
  slideUp: {
    initial: { opacity: 0, y: 16 },
    animate: { opacity: 1, y: 0 },
    exit: { opacity: 0, y: -8 },
  },
  scaleIn: {
    initial: { opacity: 0, scale: 0.95 },
    animate: { opacity: 1, scale: 1 },
    exit: { opacity: 0, scale: 0.95 },
  },
};

interface GlassPanelProps {
  children: React.ReactNode;
  className?: string;
  padding?: 'none' | 'sm' | 'md' | 'lg';
  border?: boolean;
  hover?: boolean;
  animation?: AnimationVariant;
  loading?: boolean;
  skeleton?: React.ReactNode;
}

const paddingMap = {
  none: '',
  sm: 'p-4',
  md: 'p-6',
  lg: 'p-10',
};

export function GlassPanel({
  children,
  className,
  padding = 'md',
  border = true,
  hover = false,
  animation,
  loading = false,
  skeleton,
}: GlassPanelProps) {
  const Component = animation ? motion.div : 'div';
  const motionProps = animation
    ? {
        variants: animationVariants[animation],
        initial: 'initial' as const,
        animate: 'animate' as const,
        exit: 'exit' as const,
        transition: { duration: 0.3, ease: 'easeOut' },
      }
    : {};

  return (
    <Component
      className={cn(
        'glass-panel',
        border && 'border border-white/10',
        paddingMap[padding],
        hover && 'hover:bg-white/[0.04] transition-colors',
        className,
      )}
      {...motionProps}
    >
      {loading ? (
        skeleton || (
          <div className="space-y-3 animate-pulse">
            <div className="h-4 bg-white/10 w-2/3" />
            <div className="h-3 bg-white/5 w-full" />
            <div className="h-3 bg-white/5 w-4/5" />
          </div>
        )
      ) : (
        children
      )}
    </Component>
  );
}
