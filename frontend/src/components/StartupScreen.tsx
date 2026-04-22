import React, { useEffect, useState, useCallback } from 'react';
import { motion, AnimatePresence } from 'motion/react';
import { Cpu } from 'lucide-react';
import { EventsOn, EventsOff } from '../../wailsjs/runtime/runtime';

interface StartupPhase {
  phase: number;
  name: string;
  elapsed: number;
}

export default function StartupScreen() {
  const [visible, setVisible] = useState(true);
  const [phase, setPhase] = useState(0);
  const [phases, setPhases] = useState<StartupPhase[]>([]);
  const [fadeOut, setFadeOut] = useState(false);

  const phaseLabels = ['', 'Loading Core Services...', 'Initializing Engine...', 'Preparing Environment...'];

  const handlePhaseComplete = useCallback((phaseNum: number, data: { phase: number; elapsed: number }) => {
    setPhase(phaseNum + 1);
    setPhases(prev => [...prev, { phase: data.phase, name: phaseLabels[data.phase] || '', elapsed: data.elapsed }]);
    if (phaseNum === 0) {
      setTimeout(() => setFadeOut(true), 600);
    }
  }, [phaseLabels]);

  useEffect(() => {
    const onPhase1 = (data: { phase: number; elapsed: number }) => handlePhaseComplete(0, data);
    const onPhase2 = (data: { phase: number; elapsed: number }) => handlePhaseComplete(1, data);
    const onPhase3 = (data: { phase: number; elapsed: number }) => handlePhaseComplete(2, data);

    EventsOn('startup:phase1_complete', onPhase1);
    EventsOn('startup:phase2_complete', onPhase2);
    EventsOn('startup:phase3_complete', onPhase3);

    return () => {
      EventsOff('startup:phase1_complete');
      EventsOff('startup:phase2_complete');
      EventsOff('startup:phase3_complete');
    };
  }, [handlePhaseComplete]);

  useEffect(() => {
    if (fadeOut) {
      const timer = setTimeout(() => setVisible(false), 500);
      return () => clearTimeout(timer);
    }
  }, [fadeOut]);

  return (
    <AnimatePresence>
      {visible && (
        <motion.div
          initial={{ opacity: 1 }}
          animate={{ opacity: fadeOut ? 0 : 1 }}
          exit={{ opacity: 0 }}
          transition={{ duration: 0.5, ease: 'easeInOut' }}
          className="fixed inset-0 z-[9999] flex flex-col items-center justify-center bg-background"
        >
          <motion.div
            initial={{ scale: 0.9, opacity: 0 }}
            animate={{ scale: 1, opacity: 1 }}
            transition={{ duration: 0.6, ease: 'easeOut' }}
            className="flex flex-col items-center"
          >
            <div className="relative mb-8">
              <div className="w-20 h-20 border-2 border-white/10 rounded-full flex items-center justify-center">
                <Cpu size={36} className="text-primary" />
              </div>
              <motion.div
                className="absolute inset-0 rounded-full border-2 border-primary/30"
                animate={{ scale: [1, 1.3, 1], opacity: [0.3, 0, 0.3] }}
                transition={{ duration: 2, repeat: Infinity, ease: 'easeInOut' }}
              />
            </div>

            <h1 className="text-3xl font-black tracking-tighter text-primary mb-2">
              GenPulse
            </h1>
            <p className="text-[10px] uppercase font-bold tracking-[0.3em] text-white/40 mb-12">
              Cognitive Lab / V2.4
            </p>

            <div className="w-64 mb-8">
              <div className="h-[2px] bg-white/5 rounded-full overflow-hidden">
                <motion.div
                  className="h-full bg-primary"
                  initial={{ width: '0%' }}
                  animate={{ width: `${Math.min((phase / 3) * 100, 100)}%` }}
                  transition={{ duration: 0.5, ease: 'easeOut' }}
                />
              </div>
            </div>

            <div className="text-center">
              <p className="text-xs font-medium text-white/60 mb-4 font-mono">
                {phaseLabels[phase] || phaseLabels[0]}
              </p>
              <div className="flex flex-col items-center space-y-2">
                {[1, 2, 3].map(p => (
                  <div
                    key={p}
                    className={`flex items-center space-x-3 transition-all duration-500 ${
                      phase >= p ? 'text-white/80' : 'text-white/20'
                    }`}
                  >
                    <div className={`w-2 h-2 rounded-full ${
                      phase > p ? 'bg-primary' :
                      phase === p ? 'bg-primary animate-pulse' :
                      'bg-white/10'
                    }`} />
                    <span className="text-[10px] font-bold uppercase tracking-widest">
                      {phaseLabels[p]}
                    </span>
                    {phase > p && (
                      <span className="text-[9px] font-mono text-white/30">
                        OK
                      </span>
                    )}
                  </div>
                ))}
              </div>
            </div>
          </motion.div>
        </motion.div>
      )}
    </AnimatePresence>
  );
}
