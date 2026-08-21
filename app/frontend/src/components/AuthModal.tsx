import React, { useState } from 'react';
import { ShieldCheck, Mail, Lock, User, Phone, Info, AlertTriangle } from 'lucide-react';
import type { User as UserType } from '../types';

interface AuthModalProps {
  isOpen: boolean;
  onClose: () => void;
  onLoginSuccess: (user: UserType, token: string) => void;
  onRegister: (data: { name: string; email: string; password: string; phone: string; role: string }) => Promise<UserType | null>;
  onLogin: (email: string, password: string) => Promise<UserType | null>;
}

export const AuthModal: React.FC<AuthModalProps> = ({
  isOpen,
  onClose,
  onLoginSuccess,
  onRegister,
  onLogin,
}) => {
  const [mode, setMode] = useState<'LOGIN' | 'REGISTER'>('LOGIN');

  // Form states - pure manual input
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [phone, setPhone] = useState('');

  const [errorMessage, setErrorMessage] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  if (!isOpen) return null;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setErrorMessage(null);
    setIsLoading(true);

    try {
      if (mode === 'REGISTER') {
        // Public Registration is strictly USER (Driver EV)
        const user = await onRegister({ name, email, password, phone, role: 'USER' });
        if (user) {
          const token = localStorage.getItem('volt_token') || 'jwt-token';
          onLoginSuccess(user, token);
          onClose();
        }
      } else {
        const user = await onLogin(email, password);
        if (user) {
          const token = localStorage.getItem('volt_token') || 'jwt-token';
          onLoginSuccess(user, token);
          onClose();
        } else {
          setErrorMessage('Email atau password tidak terdaftar di database.');
        }
      }
    } catch (err: any) {
      let msg = err.message || 'Terjadi kesalahan pada sistem. Silakan coba beberapa saat lagi.';
      if (msg.toLowerCase().includes('already registered') || msg.toLowerCase().includes('already exists')) {
        msg = 'Alamat Email ini sudah terdaftar. Silakan gunakan email lain atau masuk melalui menu Login Akun.';
      } else if (msg.toLowerCase().includes('failed to fetch')) {
        msg = 'Gagal terhubung ke Server Authentication. Pastikan layanan backend aktif.';
      }
      setErrorMessage(msg);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="fixed inset-0 z-50 bg-black/75 backdrop-blur-xs flex items-center justify-center p-4 overflow-y-auto">
      <div className="bg-white border border-[#262626] max-w-md w-full p-6 sm:p-8 space-y-6 relative my-8 shadow-2xl">
        {/* Dark Header Banner with Bright Contrast Text */}
        <div className="bg-[#1a2129] p-4 -mx-6 -mt-6 sm:-mx-8 sm:-mt-8 flex items-center justify-between border-b border-[#333333]">
          <div className="flex items-center gap-3">
            <div className="p-2 bg-[#1c69d4] shadow-md">
              <ShieldCheck className="w-6 h-6 text-white" />
            </div>
            <div>
              <span className="text-[10px] font-extrabold text-white bg-white/20 px-2 py-0.5 uppercase tracking-wider inline-block">
                VOLTHUB AUTHENTICATION SERVICE
              </span>
              <h2 className="text-xl font-black tracking-tight text-white mt-1">
                {mode === 'LOGIN' ? 'MASUK AKUN SPKLU' : 'DAFTAR AKUN PENGEMBUNG (DRIVER EV)'}
              </h2>
            </div>
          </div>
          <button
            onClick={onClose}
            className="text-[#9a9a9a] hover:text-white text-xl font-bold p-1 transition-colors"
          >
            ✕
          </button>
        </div>

        {/* Mode Switcher */}
        <div className="flex border-b border-[#e6e6e6]">
          <button
            type="button"
            onClick={() => {
              setMode('LOGIN');
              setErrorMessage(null);
            }}
            className={`flex-1 py-3 text-xs font-bold uppercase tracking-wider border-b-2 transition-colors ${
              mode === 'LOGIN'
                ? 'border-[#1c69d4] text-[#1c69d4] bg-[#fafafa]'
                : 'border-transparent text-[#6b6b6b]'
            }`}
          >
            LOGIN AKUN
          </button>
          <button
            type="button"
            onClick={() => {
              setMode('REGISTER');
              setErrorMessage(null);
            }}
            className={`flex-1 py-3 text-xs font-bold uppercase tracking-wider border-b-2 transition-colors ${
              mode === 'REGISTER'
                ? 'border-[#1c69d4] text-[#1c69d4] bg-[#fafafa]'
                : 'border-transparent text-[#6b6b6b]'
            }`}
          >
            REGISTRASI BARU
          </button>
        </div>

        {/* Manual Authentication Form */}
        <form onSubmit={handleSubmit} className="space-y-4 text-xs">
          {errorMessage && (
            <div className="p-3.5 bg-[#dc2626]/10 border border-[#dc2626]/40 text-[#b91c1c] font-bold text-xs flex items-start gap-2.5 shadow-xs">
              <AlertTriangle className="w-4 h-4 text-[#dc2626] shrink-0 mt-0.5" />
              <div className="leading-relaxed">{errorMessage}</div>
            </div>
          )}

          {mode === 'REGISTER' && (
            <>
              <div className="p-3 bg-[#1c69d4]/10 border border-[#1c69d4]/30 text-[#1c69d4] text-[11px] font-bold flex items-start gap-2">
                <Info className="w-4 h-4 shrink-0 mt-0.5" />
                <span>
                  Registrasi publik khusus untuk Akun Pengemudi EV (User Driver). Akun Admin & Operator dikelola secara privat oleh sistem.
                </span>
              </div>

              <div>
                <label className="text-[#262626] font-bold uppercase tracking-wider block mb-1">Nama Lengkap</label>
                <div className="relative">
                  <User className="w-4 h-4 absolute left-3 top-3 text-[#9a9a9a]" />
                  <input
                    type="text"
                    required
                    placeholder="Contoh: Ferdi Ardiansyah"
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                    className="w-full h-11 pl-10 pr-3 bg-white border border-[#cccccc] text-[#262626] focus:border-[#1c69d4]"
                  />
                </div>
              </div>

              <div>
                <label className="text-[#262626] font-bold uppercase tracking-wider block mb-1">Nomor Telepon (WhatsApp)</label>
                <div className="relative">
                  <Phone className="w-4 h-4 absolute left-3 top-3 text-[#9a9a9a]" />
                  <input
                    type="text"
                    required
                    placeholder="Contoh: 081234567890"
                    value={phone}
                    onChange={(e) => setPhone(e.target.value)}
                    className="w-full h-11 pl-10 pr-3 bg-white border border-[#cccccc] text-[#262626] focus:border-[#1c69d4]"
                  />
                </div>
              </div>
            </>
          )}

          <div>
            <label className="text-[#262626] font-bold uppercase tracking-wider block mb-1">Alamat Email</label>
            <div className="relative">
              <Mail className="w-4 h-4 absolute left-3 top-3 text-[#9a9a9a]" />
              <input
                type="email"
                required
                placeholder={mode === 'LOGIN' ? 'Masukkan email (User, Admin, atau Operator)' : 'contoh: driver@gmail.com'}
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                className="w-full h-11 pl-10 pr-3 bg-white border border-[#cccccc] text-[#262626] focus:border-[#1c69d4]"
              />
            </div>
          </div>

          <div>
            <label className="text-[#262626] font-bold uppercase tracking-wider block mb-1">Kata Sandi (Password)</label>
            <div className="relative">
              <Lock className="w-4 h-4 absolute left-3 top-3 text-[#9a9a9a]" />
              <input
                type="password"
                required
                placeholder="••••••••"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                className="w-full h-11 pl-10 pr-3 bg-white border border-[#cccccc] text-[#262626] focus:border-[#1c69d4]"
              />
            </div>
          </div>

          <button
            type="submit"
            disabled={isLoading}
            className="bmw-btn-primary w-full mt-2"
          >
            {isLoading ? (
              <span>MEMPROSES AUTHENTICATION...</span>
            ) : (
              <span>{mode === 'LOGIN' ? 'MASUK SEKARANG' : 'DAFTAR AKUN DRIVER EV'}</span>
            )}
          </button>
        </form>
      </div>
    </div>
  );
};
