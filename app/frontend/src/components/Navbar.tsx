import React from 'react';
import { MapPin, Car, Calendar, Activity, CreditCard, Wrench, LogIn, LogOut, UserCheck, ShieldCheck } from 'lucide-react';
import type { User } from '../types';

interface NavbarProps {
  activeTab: string;
  setActiveTab: (tab: string) => void;
  currentUser: User | null;
  onOpenAuthModal: () => void;
  onLogout: () => void;
  hasActiveSession: boolean;
}

export const Navbar: React.FC<NavbarProps> = ({
  activeTab,
  setActiveTab,
  currentUser,
  onOpenAuthModal,
  onLogout,
  hasActiveSession,
}) => {
  // Strict Role-Based Navigation Items
  let navItems: { id: string; label: string; icon: any; badge?: string }[] = [];

  if (currentUser?.role === 'ADMIN') {
    // Admin Menu: Portal Manajemen Pusat SPKLU, Laman Operator, & Cari SPKLU
    navItems = [
      { id: 'admin', label: 'Portal Manajemen Pusat SPKLU', icon: ShieldCheck },
      { id: 'operator', label: 'Laman Operator & Monitoring Realtime', icon: Wrench },
      { id: 'finder', label: 'Cari SPKLU', icon: MapPin },
    ];
  } else if (currentUser?.role === 'OPERATOR') {
    // Operator Menu: Monitoring Realtime & Control Center
    navItems = [
      { id: 'operator', label: 'Laman Operator & Monitoring Realtime', icon: Wrench },
      { id: 'finder', label: 'Cari SPKLU', icon: MapPin },
    ];
  } else {
    // Driver User & Guest Menu Access
    navItems = [
      { id: 'finder', label: 'Cari SPKLU', icon: MapPin },
      { id: 'garage', label: 'Garasi EV', icon: Car },
      { id: 'bookings', label: 'Booking Saya', icon: Calendar },
      { 
        id: 'session', 
        label: 'Live Telemetry', 
        icon: Activity,
        badge: hasActiveSession ? 'CHARGING' : undefined 
      },
      { id: 'billing', label: 'Invoices & Payment', icon: CreditCard },
    ];
  }

  return (
    <header className="sticky top-0 z-50 bg-white border-b border-[#e6e6e6] shadow-xs">
      {/* M-Stripe Top Accent Line */}
      <div className="m-stripe-bar" />

      <div className="max-w-[1440px] mx-auto px-4 lg:px-8 h-16 flex items-center justify-between">
        {/* Brand Logo */}
        <div 
          className="flex items-center gap-3 cursor-pointer group"
          onClick={() => setActiveTab((currentUser?.role === 'ADMIN' || currentUser?.role === 'OPERATOR') ? 'operator' : 'finder')}
        >
          <div>
            <div className="flex items-center gap-2">
              <span className="font-extrabold text-xl tracking-tight text-[#262626]">
                VOLT<span className="text-[#1c69d4]">HUB</span>
              </span>
              <span className="px-2 py-0.5 text-[10px] font-bold bg-[#fafafa] text-[#262626] border border-[#e6e6e6] uppercase tracking-wider">
                CORPORATE SPKLU
              </span>
            </div>
            <p className="text-[11px] text-[#6b6b6b] font-light hidden sm:block">
              Sistem Booking Charger EV Indonesia
            </p>
          </div>
        </div>

        {/* Horizontal Nav Links (Role Scoped) */}
        <nav className="hidden md:flex items-center space-x-1">
          {navItems.map((item) => {
            const Icon = item.icon;
            const isActive = activeTab === item.id;
            return (
              <button
                key={item.id}
                onClick={() => setActiveTab(item.id)}
                className={`h-16 px-4 flex items-center gap-2 text-sm font-semibold transition-all relative border-b-2 ${
                  isActive
                    ? 'border-[#1c69d4] text-[#1c69d4] bg-[#fafafa]'
                    : 'border-transparent text-[#262626] hover:text-[#1c69d4] hover:bg-[#fafafa]'
                }`}
              >
                <Icon className={`w-4 h-4 ${isActive ? 'text-[#1c69d4]' : 'text-[#6b6b6b]'}`} />
                <span>{item.label}</span>
                {item.badge && (
                  <span className="w-2 h-2 bg-[#22c55e] absolute top-3 right-3 animate-ping" />
                )}
              </button>
            );
          })}
        </nav>

        {/* User Profile & Role Switcher */}
        <div className="flex items-center gap-3">
          {currentUser ? (
            <>
              <div className="flex items-center gap-2 px-3 py-1.5 bg-[#fafafa] border border-[#cccccc] text-xs font-bold text-[#262626]">
                <UserCheck className="w-4 h-4 text-[#1c69d4]" />
                <span className="hidden sm:inline">Role:</span>
                <span className={`px-2 py-0.5 text-[10px] font-extrabold uppercase ${
                  currentUser.role === 'ADMIN'
                    ? 'bg-[#1a2129] text-white'
                    : currentUser.role === 'OPERATOR'
                    ? 'bg-[#eab308] text-white'
                    : 'bg-[#1c69d4] text-white'
                }`}>
                  {currentUser.role}
                </span>
              </div>

              <div className="flex items-center gap-2.5 pl-3 border-l border-[#e6e6e6]">
                <div className="w-9 h-9 bg-[#1a2129] text-white flex items-center justify-center font-bold text-sm">
                  {currentUser.name.charAt(0)}
                </div>
                <div className="hidden lg:block text-left">
                  <div className="text-xs font-bold text-[#262626]">{currentUser.name}</div>
                  <div className="text-[11px] text-[#6b6b6b] font-light">{currentUser.email}</div>
                </div>

                <button
                  onClick={onLogout}
                  className="p-2 text-[#9a9a9a] hover:text-[#dc2626] transition-colors ml-1"
                  title="Keluar (Logout)"
                >
                  <LogOut className="w-4 h-4" />
                </button>
              </div>
            </>
          ) : (
            <button
              onClick={onOpenAuthModal}
              className="bmw-btn-primary text-xs shrink-0"
            >
              <LogIn className="w-4 h-4" />
              <span>LOGIN / REGISTRASI</span>
            </button>
          )}
        </div>
      </div>
    </header>
  );
};
