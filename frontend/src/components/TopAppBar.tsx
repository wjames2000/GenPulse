import React, { useState } from 'react';

const TopAppBar: React.FC = () => {
  const [searchQuery, setSearchQuery] = useState('');

  return (
    <header className="flex justify-between items-center w-full px-6 py-3 bg-[#131318]/80 backdrop-blur-xl shadow-2xl shadow-black/40 sticky top-0 z-50">
      {/* Left: Brand / Search context */}
      <div className="flex items-center gap-6 flex-1">
        <span className="text-xl font-bold tracking-tight text-[#C0C1FF] font-['Inter'] antialiased hidden lg:block">
          Genpulse AI
        </span>
        <div className="relative w-64">
          <span className="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-[#908FA1] text-sm">
            search
          </span>
          <input
            className="w-full bg-surface-container-lowest text-on-surface text-sm rounded-lg pl-9 pr-3 py-1.5 border-none focus:ring-0 focus:border-b-2 focus:border-b-primary transition-all placeholder:text-outline outline-none"
            placeholder="Search settings..."
            type="text"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
          />
        </div>
      </div>

      {/* Right: Actions & Profile */}
      <div className="flex items-center gap-4">
        <div className="flex items-center gap-1 text-[#C0C1FF]">
          <button className="p-2 rounded-full hover:bg-[#35343a]/50 transition-colors duration-200 active:scale-95 transition-transform flex items-center justify-center">
            <span className="material-symbols-outlined">notifications</span>
          </button>
          <button className="p-2 rounded-full hover:bg-[#35343a]/50 transition-colors duration-200 active:scale-95 transition-transform flex items-center justify-center">
            <span className="material-symbols-outlined">settings</span>
          </button>
          <button className="p-2 rounded-full hover:bg-[#35343a]/50 transition-colors duration-200 active:scale-95 transition-transform flex items-center justify-center">
            <span className="material-symbols-outlined">help</span>
          </button>
        </div>
        <div className="w-8 h-8 rounded-full bg-surface-container-high overflow-hidden border border-outline-variant/15 flex-shrink-0 cursor-pointer hover:ring-2 hover:ring-primary/50 transition-all">
          <img
            alt="User avatar"
            className="w-full h-full object-cover"
            src="https://lh3.googleusercontent.com/aida-public/AB6AXuCeTVh53Ij19u5mGo-i-arSIOndiU1YLxujdHTfHh_E0IvrGptF61hoB6uTlPX0U7Kci35PJc8aYaBA7c-N-wYmEG4nUoFyqmlG1O7pz5I1UzmZcXvg-tkSNhnGZJZTTWmPwR0DSl5-90Lx4Alk-QghY-btO7bqFELE5T2GgVBfnqLOcVL1MeGE_rpKhE7enTWkesGebKs1v61NuZ123hVQuTmnBoUldayUnT5UsI7pLw0WUuYTfVJzpLvm2JGFvN_jxegXP_YW5J13"
          />
        </div>
      </div>
    </header>
  );
};

export default TopAppBar;