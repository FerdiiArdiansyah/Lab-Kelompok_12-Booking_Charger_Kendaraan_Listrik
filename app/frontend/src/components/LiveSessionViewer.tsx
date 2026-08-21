import React, { useState, useEffect } from 'react';
import { Activity, ShieldCheck, Square, DollarSign, Clock, BatteryCharging, Radio } from 'lucide-react';
import type { ChargingSession, Station } from '../types';

interface LiveSessionViewerProps {
  session: ChargingSession | null;
  stations: Station[];
  onFinishSession: (sessionId: string, finalKwh: number) => void;
}

export const LiveSessionViewer: React.FC<LiveSessionViewerProps> = ({
  session,
  stations,
  onFinishSession,
}) => {
  const [liveKwh, setLiveKwh] = useState<number>(session?.consumedKwh ?? 0);
  const [livePower, setLivePower] = useState<number>(49.5);
  const [liveVoltage, setLiveVoltage] = useState<number>(399.8);
  const [liveAmpere, setLiveAmpere] = useState<number>(123.8);

  useEffect(() => {
    if (session) {
      setLiveKwh(session.consumedKwh ?? 0);
    }
  }, [session]);

  // Live telemetry pulse simulator
  useEffect(() => {
    if (!session || session.status !== 'IN_PROGRESS') return;

    const interval = setInterval(() => {
      setLiveKwh((prev) => +((prev || 0) + 0.08).toFixed(2));
      setLivePower(+(49.0 + Math.random() * 2.0).toFixed(1));
      setLiveVoltage(+(398.0 + Math.random() * 4.0).toFixed(1));
      setLiveAmpere(+(122.0 + Math.random() * 3.0).toFixed(1));
    }, 2000);

    return () => clearInterval(interval);
  }, [session]);

  if (!session) {
    return (
      <div className="bmw-card p-12 text-center space-y-4 max-w-xl mx-auto my-12">
        <div className="w-16 h-16 bg-[#fafafa] border border-[#cccccc] flex items-center justify-center mx-auto text-[#1c69d4]">
          <Activity className="w-8 h-8" />
        </div>
        <h2 className="text-2xl font-bold text-[#262626]">TIDAK ADA SESI CHARGING AKTIF</h2>
        <p className="text-xs text-[#6b6b6b] max-w-md mx-auto font-light">
          Silakan pesan slot charger di menu <strong>Cari SPKLU</strong> dan lakukan Check-In untuk memulai stream telemetry charging secara real-time.
        </p>
      </div>
    );
  }

  const currentKwhVal = liveKwh ?? 0;
  const station = (stations || []).find((st) => st.slots && st.slots.some((sl) => sl.id === session.slotId));
  const activeTariff = station?.activeTariff?.pricePerKwh || 2467;
  const currentCost = currentKwhVal * activeTariff + 25000;
  const socPercentage = Math.min(100, Math.round((currentKwhVal / 60) * 100));

  return (
    <div className="space-y-6 max-w-5xl mx-auto">
      {/* BMW Hero Banner Status */}
      <div className="bmw-hero-dark p-6 sm:p-8 flex flex-col md:flex-row items-center justify-between gap-6 border-b border-[#e6e6e6]">
        <div className="space-y-2 text-center md:text-left">
          <div className="inline-flex items-center gap-2 px-3 py-1 bg-[#22c55e] text-white text-xs font-bold uppercase tracking-wider">
            <Radio className="w-3.5 h-3.5 animate-pulse" /> Telemetry Stream Active (ID: {session.id})
          </div>
          <h1 className="text-2xl sm:text-4xl font-black text-white">
            SESI PENGISIAN LISTRIK SPKLU
          </h1>
          <p className="text-xs text-[#bbbbbb] font-light">
            {station?.name || 'SPKLU Station'} • Slot ID: {session.slotId}
          </p>
        </div>

        <button
          onClick={() => onFinishSession(session.id, liveKwh)}
          className="px-6 py-3.5 bg-[#dc2626] hover:bg-[#b91c1c] text-white font-extrabold text-xs tracking-wider uppercase transition-colors flex items-center gap-2 shrink-0"
        >
          <Square className="w-4 h-4 fill-white" />
          <span>HENTIKAN CHARGING & BUAT TAGIHAN</span>
        </button>
      </div>

      {/* Main Telemetry Gauges Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Circle SoC Meter */}
        <div className="bmw-card p-6 flex flex-col items-center justify-center text-center space-y-4">
          <span className="text-xs font-bold text-[#6b6b6b] uppercase tracking-wider">State of Charge (SoC)</span>

          <div className="relative w-44 h-44 flex items-center justify-center">
            <div className="absolute inset-0 border-4 border-[#e6e6e6]" />
            <div className="absolute inset-0 border-4 border-[#1c69d4] border-t-transparent animate-spin-slow" />

            <div className="text-center space-y-1">
              <BatteryCharging className="w-8 h-8 text-[#1c69d4] mx-auto animate-bounce" />
              <div className="text-4xl font-black text-[#262626]">{socPercentage}%</div>
              <span className="text-[10px] bg-[#22c55e] text-white px-2.5 py-0.5 font-bold uppercase tracking-wider">
                CHARGING ACTIVE
              </span>
            </div>
          </div>

          <div className="text-xs text-[#6b6b6b] font-light">
            Estimasi Baterai Terisi: <strong className="text-[#262626] font-bold">{currentKwhVal.toFixed(1)} / 60 kWh</strong>
          </div>
        </div>

        {/* 4 Spec Cell Cards */}
        <div className="lg:col-span-2 grid grid-cols-2 gap-4">
          <div className="bmw-card p-5 flex flex-col justify-between">
            <div className="flex items-center justify-between text-[#6b6b6b]">
              <span className="text-xs font-bold uppercase tracking-wider">Energi Terpakai</span>
              <Activity className="w-4 h-4 text-[#1c69d4]" />
            </div>
            <div className="mt-4">
              <div className="text-3xl font-black text-[#1c69d4]">{currentKwhVal.toFixed(2)}</div>
              <span className="text-xs text-[#6b6b6b] font-light">KiloWatt Hour (kWh)</span>
            </div>
          </div>

          <div className="bmw-card p-5 flex flex-col justify-between">
            <div className="flex items-center justify-between text-[#6b6b6b]">
              <span className="text-xs font-bold uppercase tracking-wider">Daya Pengisian</span>
              <Activity className="w-4 h-4 text-[#22c55e]" />
            </div>
            <div className="mt-4">
              <div className="text-3xl font-black text-[#22c55e]">{livePower}</div>
              <span className="text-xs text-[#6b6b6b] font-light">kW (KiloWatt)</span>
            </div>
          </div>

          <div className="bmw-card p-5 flex flex-col justify-between">
            <div className="flex items-center justify-between text-[#6b6b6b]">
              <span className="text-xs font-bold uppercase tracking-wider">Tegangan & Arus</span>
              <ShieldCheck className="w-4 h-4 text-[#262626]" />
            </div>
            <div className="mt-4">
              <div className="text-2xl font-black text-[#262626]">{liveVoltage}V / {liveAmpere}A</div>
              <span className="text-xs text-[#6b6b6b] font-light">Volt & Ampere DC</span>
            </div>
          </div>

          <div className="bmw-card p-5 flex flex-col justify-between">
            <div className="flex items-center justify-between text-[#6b6b6b]">
              <span className="text-xs font-bold uppercase tracking-wider">Biaya Berjalan</span>
              <DollarSign className="w-4 h-4 text-[#f59e0b]" />
            </div>
            <div className="mt-4">
              <div className="text-2xl font-black text-[#f59e0b]">
                Rp {currentCost.toLocaleString('id-ID', { maximumFractionDigits: 0 })}
              </div>
              <span className="text-xs text-[#6b6b6b] font-light">Termasuk Biaya Layanan ESDM</span>
            </div>
          </div>
        </div>
      </div>

      {/* Stream Table */}
      <div className="bmw-card p-6 space-y-4">
        <h3 className="text-xs font-bold text-[#262626] uppercase tracking-wider flex items-center gap-2">
          <Clock className="w-4 h-4 text-[#1c69d4]" /> Stream Record Telemetry Log
        </h3>

        <div className="overflow-x-auto">
          <table className="w-full text-left text-xs">
            <thead>
              <tr className="border-b border-[#cccccc] text-[#6b6b6b] uppercase font-bold">
                <th className="pb-3">Waktu Record</th>
                <th className="pb-3">Accumulator (kWh)</th>
                <th className="pb-3">Daya (kW)</th>
                <th className="pb-3">Tegangan (V)</th>
                <th className="pb-3">Arus (A)</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-[#e6e6e6] text-[#3c3c3c] font-mono">
              <tr>
                <td className="py-3 font-bold text-[#1c69d4]">Terbaru ({new Date().toLocaleTimeString()})</td>
                <td className="py-3 font-bold text-[#262626]">{currentKwhVal.toFixed(2)} kWh</td>
                <td className="py-3">{livePower} kW</td>
                <td className="py-3">{liveVoltage} V</td>
                <td className="py-3">{liveAmpere} A</td>
              </tr>
              {session.readings?.map((r) => (
                <tr key={r.id}>
                  <td className="py-3 text-[#6b6b6b]">{new Date(r.recordedAt).toLocaleTimeString()}</td>
                  <td className="py-3">{r.currentKwh} kWh</td>
                  <td className="py-3">{r.powerKw} kW</td>
                  <td className="py-3">{r.voltage} V</td>
                  <td className="py-3">{r.currentAmpere} A</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
};
