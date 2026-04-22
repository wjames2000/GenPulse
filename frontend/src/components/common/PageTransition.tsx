import React from 'react';
import { motion, AnimatePresence } from 'motion/react';
import { cn } from '../../utils';

interface PageTransitionProps {
  children: React.ReactNode;
  uniqueKey: string;
  className?: string;
  duration?: number;
  direction?: 'up' | 'down' | 'left' | 'right';
}

const directionVariants = {
  up: { initial: { opacity: 0, y: 16 }, exit: { opacity: 0, y: -16 } },
  down: { initial: { opacity: 0, y: -16 }, exit: { opacity: 0, y: 16 } },
  left: { initial: { opacity: 0, x: 16 }, exit: { opacity: 0, x: -16 } },
  right: { initial: { opacity: 0, x: -16 }, exit: { opacity: 0, x: 16 } },
};

export function PageTransition({
  children,
  uniqueKey,
  className,
  duration = 0.25,
  direction = 'up',
}: PageTransitionProps) {
  const variant = directionVariants[direction];

  return (
    <AnimatePresence mode="wait">
      <motion.div
        key={uniqueKey}
        initial={variant.initial}
        animate={{ opacity: 1, y: 0, x: 0 }}
        exit={variant.exit}
        transition={{ duration, ease: 'easeOut' }}
        className={cn('h-full', className)}
      >
        {children}
      </motion.div>
    </AnimatePresence>
  );
}
