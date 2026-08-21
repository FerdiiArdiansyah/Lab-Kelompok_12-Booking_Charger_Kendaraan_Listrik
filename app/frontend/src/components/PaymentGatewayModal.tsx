import React, { useState, useEffect } from 'react';
import {
  CreditCard,
  QrCode,
  ShieldCheck,
  CheckCircle2,
  Copy,
  Check,
  Lock,
  X,
  Clock,
  Building,
  Sparkles
} from 'lucide-react';
import type { Invoice, Payment } from '../types';

interface PaymentGatewayModalProps {
  invoice: Invoice;
  paymentMethod: Payment['paymentMethod'];
  onClose: () => void;
  onSuccess: (paymentMethod: Payment['paymentMethod']) => void;
}

export const PaymentGatewayModal: React.FC<PaymentGatewayModalProps> = ({
  invoice,
  paymentMethod,
  onClose,
  onSuccess,
}) => {
  const [step, setStep] = useState<'PAYMENT' | 'PROCESSING' | 'SUCCESS'>('PAYMENT');
  const [copied, setCopied] = useState(false);
  const [timeLeft, setTimeLeft] = useState(899); // 14 mins 59 secs
  const [cardNumber, setCardNumber] = useState('');
  const [cardExpiry, setCardExpiry] = useState('');
  const [cardCvv, setCardCvv] = useState('');
  const [cardHolder, setCardHolder] = useState('');
  const [otpCode, setOtpCode] = useState('');
  const [showOtp, setShowOtp] = useState(false);

  const totalAmount = Math.round(invoice.total);
  const vaNumber = paymentMethod === 'VA_BCA'
    ? `88012${invoice.id.replace(/[^0-9]/g, '').padStart(6, '0')}99`
    : paymentMethod === 'VA_MANDIRI'
    ? `89008${invoice.id.replace(/[^0-9]/g, '').padStart(6, '0')}77`
    : `88099${invoice.id.replace(/[^0-9]/g, '').padStart(6, '0')}11`;

  useEffect(() => {
    const timer = setInterval(() => {
      setTimeLeft((prev) => (prev > 0 ? prev - 1 : 0));
    }, 1000);
    return () => clearInterval(timer);
  }, []);

  const formatTime = (seconds: number) => {
    const m = Math.floor(seconds / 60);
    const s = seconds % 60;
    return `${m.toString().padStart(2, '0')}:${s.toString().padStart(2, '0')}`;
  };

  const handleCopyVa = () => {
    navigator.clipboard.writeText(vaNumber);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  const handleSimulatePayment = () => {
    if (paymentMethod === 'CREDIT_CARD' && !showOtp) {
      if (!cardNumber || !cardExpiry || !cardCvv) {
        alert('Mohon lengkapi data kartu kredit Anda.');
        return;
      }
      setShowOtp(true);
      return;
    }

    setStep('PROCESSING');
    setTimeout(() => {
      setStep('SUCCESS');
      setTimeout(() => {
        onSuccess(paymentMethod);
      }, 1500);
    }, 1500);
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/70 backdrop-blur-sm p-4 animate-in fade-in duration-200">
      <div className="bg-white border-2 border-[#1c69d4] w-full max-w-lg shadow-2xl overflow-hidden relative">
        {/* Payment Gateway Header */}
        <div className="bg-[#1c69d4] text-white p-4 flex items-center justify-between">
          <div className="flex items-center gap-2">
            <div className="bg-white text-[#1c69d4] px-2 py-0.5 font-black text-xs uppercase tracking-wider">
              VoltHub Gateway
            </div>
            <span className="text-[11px] font-medium opacity-90 flex items-center gap-1">
              <Lock className="w-3 h-3" /> 256-Bit SSL Encrypted
            </span>
          </div>
          <button
            onClick={onClose}
            className="p-1 hover:bg-white/10 text-white transition-colors cursor-pointer"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        {/* Merchant & Order Summary Bar */}
        <div className="bg-[#fafafa] border-b border-[#e6e6e6] p-4 flex items-center justify-between text-xs">
          <div>
            <div className="text-[#6b6b6b] font-medium">Merchant:</div>
            <div className="font-bold text-[#262626]">PT PLN VoltHub Charging Services</div>
          </div>
          <div className="text-right">
            <div className="text-[#6b6b6b] font-medium">Tagihan #{invoice.id}:</div>
            <div className="font-black text-[#1c69d4] text-base">
              Rp {totalAmount.toLocaleString('id-ID')}
            </div>
          </div>
        </div>

        {/* Modal Body depending on step */}
        {step === 'PROCESSING' && (
          <div className="p-12 text-center space-y-4">
            <div className="w-16 h-16 border-4 border-[#1c69d4] border-t-transparent rounded-full animate-spin mx-auto" />
            <h3 className="text-lg font-extrabold text-[#262626]">Verifikasi Pembayaran Instan...</h3>
            <p className="text-xs text-[#6b6b6b]">
              Menghubungkan ke server Payment Gateway & Billing Microservice (Port 8085).
            </p>
          </div>
        )}

        {step === 'SUCCESS' && (
          <div className="p-12 text-center space-y-4 animate-in zoom-in-95 duration-200">
            <div className="w-20 h-20 bg-[#22c55e]/10 border-2 border-[#22c55e] rounded-full flex items-center justify-center mx-auto text-[#22c55e]">
              <CheckCircle2 className="w-12 h-12" />
            </div>
            <div>
              <h3 className="text-xl font-black text-[#262626]">PEMBAYARAN LUNAS!</h3>
              <p className="text-xs text-[#6b6b6b] mt-1">
                Struk transaksi resmi telah dikonfirmasi oleh Billing Service.
              </p>
            </div>
            <div className="p-3 bg-[#fafafa] border border-[#e6e6e6] text-[11px] font-mono text-[#3c3c3c]">
              Ref: TX-{paymentMethod}-{Date.now()}
            </div>
          </div>
        )}

        {step === 'PAYMENT' && (
          <div className="p-6 space-y-6">
            {/* Timer Banner */}
            <div className="bg-[#fffbe0] border border-[#ffe58f] p-2.5 text-xs text-[#d48806] flex items-center justify-between font-medium">
              <span className="flex items-center gap-1.5">
                <Clock className="w-4 h-4" /> Selesaikan pembayaran dalam:
              </span>
              <span className="font-mono font-bold text-sm text-[#d48806]">{formatTime(timeLeft)}</span>
            </div>

            {/* Selected Gateway Detail */}
            {paymentMethod === 'QRIS' && (
              <div className="text-center space-y-4">
                <div className="inline-block p-4 bg-white border-2 border-[#262626] shadow-md relative">
                  <QrCode className="w-44 h-44 text-[#262626] mx-auto" />
                  <div className="mt-2 text-[10px] font-bold text-[#6b6b6b] uppercase tracking-wider">
                    NMID: ID1020392817293 • QRIS STANDAR NASIONAL
                  </div>
                </div>

                <p className="text-xs text-[#6b6b6b] max-w-xs mx-auto">
                  Scan kode QRIS menggunakan GoPay, OVO, Dana, ShopeePay, BCA Mobile, atau aplikasi Mobile Banking pilihan Anda.
                </p>

                <div className="p-3 bg-[#f0f7ff] border border-[#1c69d4]/30 text-xs text-[#1c69d4] flex items-center justify-center gap-2 font-bold">
                  <Sparkles className="w-4 h-4" /> Auto-Settlement Mode Aktif
                </div>
              </div>
            )}

            {(paymentMethod === 'VA_BCA' || paymentMethod === 'VA_MANDIRI' || paymentMethod === 'E_WALLET_GOPAY') && (
              <div className="space-y-4">
                <div className="p-4 bg-[#fafafa] border border-[#e6e6e6] space-y-3">
                  <div className="flex items-center justify-between text-xs">
                    <span className="text-[#6b6b6b] font-medium flex items-center gap-1.5">
                      <Building className="w-4 h-4 text-[#1c69d4]" />
                      Nomor Virtual Account ({paymentMethod.replace('_', ' ')}):
                    </span>
                  </div>

                  <div className="flex items-center justify-between bg-white p-3 border border-[#262626]">
                    <span className="font-mono text-lg font-black text-[#262626] tracking-wider">
                      {vaNumber}
                    </span>
                    <button
                      onClick={handleCopyVa}
                      className="px-3 py-1.5 bg-[#1c69d4] hover:bg-[#1552a8] text-white text-xs font-bold flex items-center gap-1 transition-colors cursor-pointer"
                    >
                      {copied ? <Check className="w-3.5 h-3.5" /> : <Copy className="w-3.5 h-3.5" />}
                      <span>{copied ? 'TERSALIN' : 'SALIN'}</span>
                    </button>
                  </div>
                </div>

                <div className="space-y-2 text-xs text-[#6b6b6b]">
                  <div className="font-bold text-[#262626]">Petunjuk Transfer M-Banking / ATM:</div>
                  <ol className="list-decimal list-inside space-y-1 text-[11px] bg-[#fafafa] p-3 border border-[#e6e6e6]">
                    <li>Buka aplikasi Mobile Banking atau datangi ATM terdekat.</li>
                    <li>Pilih menu <strong>Transfer / Pembayaran Virtual Account</strong>.</li>
                    <li>Masukkan nomor Virtual Account: <strong className="text-[#1c69d4]">{vaNumber}</strong>.</li>
                    <li>Periksa nominal <strong className="text-[#262626]">Rp {totalAmount.toLocaleString('id-ID')}</strong> lalu konfirmasi.</li>
                  </ol>
                </div>
              </div>
            )}

            {paymentMethod === 'CREDIT_CARD' && (
              <div className="space-y-4">
                {!showOtp ? (
                  <div className="space-y-3">
                    <div className="p-3 bg-[#fafafa] border border-[#e6e6e6] space-y-2 text-xs">
                      <label className="block font-bold text-[#262626]">Nomor Kartu Kredit / Debit:</label>
                      <div className="relative">
                        <input
                          type="text"
                          placeholder="4532 8912 3456 7890"
                          value={cardNumber}
                          onChange={(e) => setCardNumber(e.target.value)}
                          className="w-full p-2.5 bg-white border border-[#cccccc] font-mono text-sm focus:border-[#1c69d4] focus:outline-none"
                        />
                        <CreditCard className="w-4 h-4 text-[#6b6b6b] absolute right-3 top-3" />
                      </div>

                      <div className="grid grid-cols-2 gap-2 pt-1">
                        <div>
                          <label className="block font-bold text-[#262626]">Masa Berlaku:</label>
                          <input
                            type="text"
                            placeholder="MM/YY"
                            value={cardExpiry}
                            onChange={(e) => setCardExpiry(e.target.value)}
                            className="w-full p-2.5 bg-white border border-[#cccccc] font-mono text-xs focus:border-[#1c69d4] focus:outline-none"
                          />
                        </div>
                        <div>
                          <label className="block font-bold text-[#262626]">CVV / CVC:</label>
                          <input
                            type="password"
                            maxLength={4}
                            placeholder="123"
                            value={cardCvv}
                            onChange={(e) => setCardCvv(e.target.value)}
                            className="w-full p-2.5 bg-white border border-[#cccccc] font-mono text-xs focus:border-[#1c69d4] focus:outline-none"
                          />
                        </div>
                      </div>

                      <div>
                        <label className="block font-bold text-[#262626] mt-1">Nama Pemegang Kartu:</label>
                        <input
                          type="text"
                          placeholder="NAMA LENGKAP PADA KARTU"
                          value={cardHolder}
                          onChange={(e) => setCardHolder(e.target.value)}
                          className="w-full p-2.5 bg-white border border-[#cccccc] text-xs focus:border-[#1c69d4] focus:outline-none uppercase"
                        />
                      </div>
                    </div>
                  </div>
                ) : (
                  <div className="p-4 bg-[#f0f7ff] border border-[#1c69d4] space-y-3 text-center">
                    <div className="text-xs font-bold text-[#1c69d4] uppercase tracking-wider flex items-center justify-center gap-1.5">
                      <ShieldCheck className="w-4 h-4" /> Otentikasi 3D Secure (OTP)
                    </div>
                    <p className="text-xs text-[#6b6b6b]">
                      Kode otentikasi transaksi telah dikirim via SMS ke HP Anda (+62 812-****-889).
                    </p>
                    <input
                      type="text"
                      maxLength={6}
                      placeholder="Masukkan 6 Digit OTP (contoh: 123456)"
                      value={otpCode}
                      onChange={(e) => setOtpCode(e.target.value)}
                      className="w-48 text-center p-2.5 bg-white border border-[#262626] font-mono text-lg font-bold mx-auto block focus:outline-none focus:ring-2 focus:ring-[#1c69d4]"
                    />
                  </div>
                )}
              </div>
            )}

            {/* Action Payment Trigger Button */}
            <div className="pt-2">
              <button
                onClick={handleSimulatePayment}
                className="bmw-btn-primary w-full py-3 text-sm flex items-center justify-center gap-2 cursor-pointer shadow-md"
              >
                <ShieldCheck className="w-5 h-5" />
                <span>
                  {paymentMethod === 'CREDIT_CARD' && !showOtp
                    ? 'LANJUTKAN KE 3D SECURE OTP'
                    : `PROSES SIMULASI PEMBAYARAN RP ${totalAmount.toLocaleString('id-ID')}`}
                </span>
              </button>
            </div>

            <div className="text-center text-[10px] text-[#6b6b6b]">
              ⚡ Transaksi dijamin aman secara otomatis oleh VoltHub Central Payment Gateway Server.
            </div>
          </div>
        )}
      </div>
    </div>
  );
};
