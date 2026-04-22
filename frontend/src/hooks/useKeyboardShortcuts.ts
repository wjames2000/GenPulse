import { useEffect, useCallback } from 'react';

export interface ShortcutMap {
  [key: string]: () => void;
}

const DEFAULT_SHORTCUTS: Record<string, { key: string; meta?: boolean; ctrl?: boolean; shift?: boolean; alt?: boolean }> = {
  'search': { key: 'k', meta: true },
  'toggleSidebar': { key: 'b', meta: true },
  'escape': { key: 'Escape' },
  'nav1': { key: '1', meta: true },
  'nav2': { key: '2', meta: true },
  'nav3': { key: '3', meta: true },
  'nav4': { key: '4', meta: true },
  'nav5': { key: '5', meta: true },
  'nav6': { key: '6', meta: true },
  'nav7': { key: '7', meta: true },
  'nav8': { key: '8', meta: true },
  'nav9': { key: '9', meta: true },
};

function matchesShortcut(e: KeyboardEvent, binding: { key: string; meta?: boolean; ctrl?: boolean; shift?: boolean; alt?: boolean }): boolean {
  if (e.key.toLowerCase() !== binding.key.toLowerCase()) return false;
  if (binding.meta && !e.metaKey) return false;
  if (binding.ctrl && !e.ctrlKey) return false;
  if (binding.shift && !e.shiftKey) return false;
  if (binding.alt && !e.altKey) return false;
  if (!binding.meta && !binding.ctrl && !binding.shift && !binding.alt) {
    if (e.metaKey || e.ctrlKey || e.shiftKey || e.altKey) return false;
  }
  return true;
}

export function useKeyboardShortcuts(handlers: ShortcutMap) {
  const handleKeyDown = useCallback(
    (e: KeyboardEvent) => {
      const isInputFocused = ['INPUT', 'TEXTAREA', 'SELECT'].includes(
        (e.target as HTMLElement)?.tagName
      );

      for (const [name, handler] of Object.entries(handlers)) {
        const binding = DEFAULT_SHORTCUTS[name];
        if (!binding) continue;

        if (matchesShortcut(e, binding)) {
          if (isInputFocused && !binding.meta && !binding.ctrl) continue;
          e.preventDefault();
          e.stopPropagation();
          handler();
          return;
        }
      }
    },
    [handlers]
  );

  useEffect(() => {
    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [handleKeyDown]);
}
