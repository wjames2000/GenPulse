import { useState, useCallback, useMemo, useRef, useEffect, type RefObject, type UIEvent } from 'react';

interface UseVirtualScrollOptions {
  totalItems: number;
  itemHeight: number;
  containerRef: RefObject<HTMLElement | null>;
  overscan?: number;
  getItemHeight?: (index: number) => number;
}

interface UseVirtualScrollReturn {
  visibleItems: { index: number; offsetY: number }[];
  startIndex: number;
  endIndex: number;
  totalHeight: number;
  offsetY: number;
  isScrolledToBottom: boolean;
  scrollToBottom: () => void;
  onScroll: (e: UIEvent<HTMLElement>) => void;
}

export function useVirtualScroll({
  totalItems,
  itemHeight,
  containerRef,
  overscan = 5,
  getItemHeight,
}: UseVirtualScrollOptions): UseVirtualScrollReturn {
  const [scrollTop, setScrollTop] = useState(0);
  const [containerHeight, setContainerHeight] = useState(0);
  const rafRef = useRef<number>(0);
  const isScrolledToBottomRef = useRef(true);

  useEffect(() => {
    const container = containerRef.current;
    if (!container) return;

    const observer = new ResizeObserver((entries) => {
      for (const entry of entries) {
        setContainerHeight(entry.contentRect.height);
      }
    });
    observer.observe(container);
    return () => observer.disconnect();
  }, [containerRef]);

  const heights = useMemo(() => {
    if (getItemHeight) {
      const h = new Array(totalItems);
      for (let i = 0; i < totalItems; i++) {
        h[i] = getItemHeight(i);
      }
      return h;
    }
    return null;
  }, [totalItems, getItemHeight, itemHeight]);

  const getItemHeightAtIndex = useCallback(
    (index: number): number => {
      if (heights) return heights[index];
      return itemHeight;
    },
    [heights, itemHeight]
  );

  const cumulativeHeights = useMemo(() => {
    if (!heights) return null;
    const cum = new Array(totalItems);
    let sum = 0;
    for (let i = 0; i < totalItems; i++) {
      sum += heights[i];
      cum[i] = sum;
    }
    return cum;
  }, [heights, totalItems]);

  const totalHeight = useMemo(() => {
    if (cumulativeHeights && totalItems > 0) {
      return cumulativeHeights[totalItems - 1];
    }
    return totalItems * itemHeight;
  }, [cumulativeHeights, totalItems, itemHeight]);

  const findStartIndex = useCallback(
    (offset: number): number => {
      if (!cumulativeHeights) {
        return Math.max(0, Math.floor(offset / itemHeight));
      }
      let lo = 0;
      let hi = cumulativeHeights.length - 1;
      while (lo < hi) {
        const mid = Math.floor((lo + hi) / 2);
        if (cumulativeHeights[mid] <= offset) {
          lo = mid + 1;
        } else {
          hi = mid;
        }
      }
      return lo;
    },
    [cumulativeHeights, itemHeight]
  );

  const startIndex = useMemo(
    () => Math.max(0, findStartIndex(scrollTop) - overscan),
    [findStartIndex, scrollTop, overscan]
  );

  const endIndex = useMemo(() => {
    const maxVisible = findStartIndex(scrollTop + containerHeight) + overscan;
    return Math.min(totalItems - 1, maxVisible);
  }, [findStartIndex, scrollTop, containerHeight, overscan, totalItems]);

  const offsetY = useMemo(() => {
    if (cumulativeHeights && startIndex > 0) {
      return cumulativeHeights[startIndex - 1];
    }
    return startIndex * itemHeight;
  }, [cumulativeHeights, startIndex, itemHeight]);

  const visibleItems = useMemo(() => {
    const items: { index: number; offsetY: number }[] = [];
    for (let i = startIndex; i <= endIndex; i++) {
      const yOff = cumulativeHeights
        ? (i > 0 ? cumulativeHeights[i - 1] : 0)
        : i * itemHeight;
      items.push({ index: i, offsetY: yOff });
    }
    return items;
  }, [startIndex, endIndex, cumulativeHeights, itemHeight]);

  const isScrolledToBottom = useMemo(() => {
    const container = containerRef.current;
    if (!container) return true;
    const threshold = 50;
    return scrollTop + containerHeight >= totalHeight - threshold;
  }, [scrollTop, containerHeight, totalHeight, containerRef]);

  const scrollToBottom = useCallback(() => {
    const container = containerRef.current;
    if (container) {
      container.scrollTop = totalHeight;
    }
  }, [containerRef, totalHeight]);

  const onScroll = useCallback(
    (e: UIEvent<HTMLElement>) => {
      if (rafRef.current) {
        cancelAnimationFrame(rafRef.current);
      }
      rafRef.current = requestAnimationFrame(() => {
        const target = e.currentTarget;
        setScrollTop(target.scrollTop);
      });
    },
    []
  );

  useEffect(() => {
    return () => {
      if (rafRef.current) {
        cancelAnimationFrame(rafRef.current);
      }
    };
  }, []);

  return {
    visibleItems,
    startIndex,
    endIndex,
    totalHeight,
    offsetY,
    isScrolledToBottom,
    scrollToBottom,
    onScroll,
  };
}
