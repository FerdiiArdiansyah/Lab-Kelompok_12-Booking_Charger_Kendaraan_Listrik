import React, { useState } from 'react';
import { Clock, Car, ShieldAlert, CheckCircle2, Zap, AlertTriangle, ArrowRight } from 'lucide-react';
import type { Station, ChargerSlot, UserVehicle, Booking } from '../types';

interface BookingModalProps {
  station: Station;
  initialSlot?: ChargerSlot;
  vehicles: UserVehicle[];
  existingBookings: Booking[];
  onClose: () => void;
  onConfirmBooking: (bookingData: {
    slotId: string;
    vehicleId: string;
    startTime: string;
    endTime: string;
  }) => void;
  onConfirmAndStartSession?: (bookingData: {
    slotId: string;
    vehicleId: string;
    startTime: string;
    endTime: string;
  }) => void;
  onJoinWaitlist: (waitlistData: {
    slotId: string;
    preferredStartTime: string;
    preferredEndTime: string;
  }) => void;
}

export const BookingModal: React.FC<BookingModalProps> = ({
  station,
  initialSlot,
  vehicles,
  existingBookings,
  onClose,
  onConfirmBooking,
  onConfirmAndStartSession,
  onJoinWaitlist,
}) => {
  const availableVehicles = vehicles || [];

  const [selectedSlotId, setSelectedSlotId] = useState<string>(
    initialSlot ? initialSlot.id : station.slots[0]?.id || ''
  );
  const [selectedVehicleId, setSelectedVehicleId] = useState<string>(
    availableVehicles[0]?.id || ''
  );
  const getLocalDateStr = (d: Date = new Date()) => {
    const year = d.getFullYear();
    const month = String(d.getMonth() + 1).padStart(2, '0');
    const day = String(d.getDate()).padStart(2, '0');
    return `${year}-${month}-${day}`;
  };

  const [bookingDate, setBookingDate] = useState<string>(getLocalDateStr());
  const [startHour, setStartHour] = useState<number>(new Date().getHours() + 1);
  const [durationMinutes, setDurationMinutes] = useState<number>(60);

  const selectedSlot = station.slots.find((s) => s.id === selectedSlotId);
  const selectedVehicle = availableVehicles.find((v) => v.id === selectedVehicleId);

  // Calculate Start & End ISO string
  const startTimeObj = new Date(`${bookingDate}T${String(startHour).padStart(2, '0')}:00:00`);
  const endTimeObj = new Date(startTimeObj.getTime() + durationMinutes * 60000);

  const startTimeIso = startTimeObj.toISOString();
  const endTimeIso = endTimeObj.toISOString();

  // 1. Vehicle Connector Compatibility Check (TICKET-07)
  const isConnectorCompatible =
    !selectedVehicle || !selectedSlot
      ? true
      : selectedVehicle.connectorType === selectedSlot.connectorType;

  // 2. Anti-Overlap Booking Check (TICKET-01)
  const hasOverlap = existingBookings.some((b) => {
    if (b.slotId !== selectedSlotId || b.status === 'CANCELLED' || b.status === 'EXPIRED') {
      return false;
    }
    const bStart = new Date(b.startTime).getTime();
    const bEnd = new Date(b.endTime).getTime();
    const reqStart = startTimeObj.getTime();
    const reqEnd = endTimeObj.getTime();
    return reqStart < bEnd && reqEnd > bStart;
  });

  // Calculate ESDM Tariff Estimation
  const estimatedKwh = selectedVehicle
    ? Math.min(selectedVehicle.batteryCapacityKwh, (selectedSlot?.maxPowerKw || 50) * (durationMinutes / 60))
    : 30;

  const basePricePerKwh = station.activeTariff ? station.activeTariff.pricePerKwh : 2467.0;
  const isUltraFast = (selectedSlot?.maxPowerKw || 0) >= 150;
  const serviceFee = isUltraFast ? 57000 : 25000;
  const subtotal = estimatedKwh * basePricePerKwh + serviceFee;
  const tax = subtotal * 0.11;
  const grandTotal = subtotal + tax;

  const handleAction = () => {
    if (hasOverlap) {
      onJoinWaitlist({
        slotId: selectedSlotId,
        preferredStartTime: startTimeIso,
        preferredEndTime: endTimeIso,
      });
    } else {
      onConfirmBooking({
        slotId: selectedSlotId,
        vehicleId: selectedVehicleId,
        startTime: startTimeIso,
        endTime: endTimeIso,
      });
    }
  };

  const handleConfirmAndStart = () => {
    if (onConfirmAndStartSession) {
      onConfirmAndStartSession({
        slotId: selectedSlotId,
        vehicleId: selectedVehicleId,
        startTime: startTimeIso,
        endTime: endTimeIso,
      });
    }
  };

  return (
    <div className="fixed inset-0 z-50 bg-black/60 backdrop-blur-xs flex items-center justify-center p-4 overflow-y-auto">
      <div className="bg-white border border-[#262626] max-w-2xl w-full p-6 sm:p-8 space-y-6 relative my-8 shadow-2xl">
        {/* Header */}
        <div className="flex items-center justify-between pb-4 border-b border-[#e6e6e6]">
          <div>
            <span className="px-2.5 py-0.5 text-[10px] font-extrabold bg-[#1c69d4] text-white uppercase tracking-wider">
              ANTI-OVERLAP ENGINE
            </span>
            <h2 className="text-2xl font-black text-[#262626] mt-1">Reservasi Slot Charger SPKLU</h2>
            <p className="text-xs text-[#6b6b6b] font-light">{station.name}</p>
          </div>
          <button
            onClick={onClose}
            className="text-[#6b6b6b] hover:text-[#262626] text-xl font-bold p-2"
          >
            ✕
          </button>
        </div>

        {/* Step 1: Select Slot & Vehicle */}
        <div className="space-y-4">
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
            {/* Slot Selector */}
            <div className="space-y-1.5">
              <label className="text-xs font-bold text-[#262626] uppercase tracking-wider flex items-center gap-1.5">
                <Zap className="w-4 h-4 text-[#1c69d4]" /> Pilih Slot Charger:
              </label>
              <select
                value={selectedSlotId}
                onChange={(e) => setSelectedSlotId(e.target.value)}
                className="w-full h-12 px-3 bg-white border border-[#cccccc] text-xs font-bold text-[#262626] focus:border-[#1c69d4] focus:outline-none"
              >
                {station.slots.map((sl) => (
                  <option key={sl.id} value={sl.id}>
                    Slot #{sl.slotNumber} - {sl.connectorType} ({sl.maxPowerKw} kW)
                  </option>
                ))}
              </select>
            </div>

            {/* Vehicle Selector */}
            <div className="space-y-1.5">
              <label className="text-xs font-bold text-[#262626] uppercase tracking-wider flex items-center gap-1.5">
                <Car className="w-4 h-4 text-[#1c69d4]" /> Pilih Kendaraan EV:
              </label>
              {availableVehicles.length > 0 ? (
                <select
                  value={selectedVehicleId}
                  onChange={(e) => setSelectedVehicleId(e.target.value)}
                  className="w-full h-12 px-3 bg-white border border-[#cccccc] text-xs font-bold text-[#262626] focus:border-[#1c69d4] focus:outline-none"
                >
                  {availableVehicles.map((v) => (
                    <option key={v.id} value={v.id}>
                      {v.brand} {v.model} ({v.connectorType}) - Plat: {v.licensePlate}
                    </option>
                  ))}
                </select>
              ) : (
                <div className="p-3 bg-amber-50 border border-amber-300 text-amber-900 text-xs font-bold flex items-center gap-2">
                  <AlertTriangle className="w-4 h-4 text-amber-600 shrink-0" />
                  <span>Garasi EV Anda masih kosong. Silakan mendaftarkan mobil EV Anda di menu Garasi EV terlebih dahulu.</span>
                </div>
              )}
            </div>
          </div>

          {/* Connector Warning Alert */}
          {!isConnectorCompatible && (
            <div className="p-4 bg-[#dc2626]/10 border border-[#dc2626] text-[#dc2626] text-xs flex items-start gap-3">
              <AlertTriangle className="w-5 h-5 shrink-0 mt-0.5" />
              <div>
                <span className="font-bold block text-sm">INKOMPATIBILITAS KONEKTOR!</span>
                Konektor slot ini adalah <strong>{selectedSlot?.connectorType}</strong>, sedangkan tipe konektor kendaraan Anda adalah <strong>{selectedVehicle?.connectorType}</strong>. Mohon pilih slot atau kendaraan yang sesuai.
              </div>
            </div>
          )}
        </div>

        {/* Step 2: 24-Hour Interactive Schedule & Time Selection */}
        <div className="space-y-3 bg-[#fafafa] p-4 border border-[#e6e6e6]">
          <div className="flex items-center justify-between">
            <div className="text-xs font-bold text-[#262626] uppercase tracking-wider flex items-center gap-2">
              <Clock className="w-4 h-4 text-[#1c69d4]" /> Jadwal 24-Jam (Bebas Bentrok / Anti-Overlap):
            </div>
            <div className="flex items-center gap-3 text-[10px] font-bold uppercase">
              <span className="flex items-center gap-1"><span className="w-2.5 h-2.5 bg-[#22c55e] inline-block border border-emerald-600" /> Tersedia</span>
              <span className="flex items-center gap-1"><span className="w-2.5 h-2.5 bg-[#1c69d4] inline-block" /> Pilihan Anda</span>
              <span className="flex items-center gap-1"><span className="w-2.5 h-2.5 bg-[#dc2626] inline-block" /> Terisi</span>
            </div>
          </div>

          {/* 24-Hour Visual Schedule Timeline Matrix */}
          <div className="space-y-1 pt-1">
            <div className="text-[10px] text-[#6b6b6b] font-bold flex justify-between">
              <span>00:00</span>
              <span>06:00</span>
              <span>12:00</span>
              <span>18:00</span>
              <span>23:59</span>
            </div>
            <div className="grid grid-cols-12 sm:grid-cols-24 gap-1">
              {Array.from({ length: 24 }).map((_, hour) => {
                const hourSlotStart = new Date(`${bookingDate}T${String(hour).padStart(2, '0')}:00:00`).getTime();
                const hourSlotEnd = hourSlotStart + 3600000;

                const isHourOccupied = existingBookings.some((b) => {
                  if (b.slotId !== selectedSlotId || b.status === 'CANCELLED' || b.status === 'EXPIRED') return false;
                  const bStart = new Date(b.startTime).getTime();
                  const bEnd = new Date(b.endTime).getTime();
                  return hourSlotStart < bEnd && hourSlotEnd > bStart;
                });

                const reqStart = startTimeObj.getTime();
                const reqEnd = endTimeObj.getTime();
                const isSelectedRange = hourSlotStart < reqEnd && hourSlotEnd > reqStart;

                let btnStyle = 'bg-emerald-50 text-emerald-800 border-emerald-300 hover:bg-emerald-100';
                if (isHourOccupied) {
                  btnStyle = 'bg-red-500 text-white border-red-600 cursor-not-allowed opacity-90';
                } else if (isSelectedRange) {
                  btnStyle = 'bg-[#1c69d4] text-white border-[#1c69d4] font-bold shadow-xs';
                }

                return (
                  <button
                    key={hour}
                    type="button"
                    title={`Jam ${String(hour).padStart(2, '0')}:00 WIB ${isHourOccupied ? '(TERISI)' : '(TERSEDIA)'}`}
                    onClick={() => {
                      if (!isHourOccupied) {
                        setStartHour(hour);
                      }
                    }}
                    className={`h-9 border text-[10px] flex items-center justify-center transition-all ${btnStyle}`}
                  >
                    {String(hour).padStart(2, '0')}
                  </button>
                );
              })}
            </div>
          </div>

          {/* Time Picker Controls */}
          <div className="grid grid-cols-1 sm:grid-cols-3 gap-3 pt-2">
            <div>
              <label className="text-[11px] text-[#6b6b6b] font-bold block mb-1">Tanggal Reservasi</label>
              <input
                type="date"
                value={bookingDate}
                onChange={(e) => setBookingDate(e.target.value)}
                className="w-full h-10 px-2 bg-white border border-[#cccccc] text-xs font-bold text-[#262626]"
              />
            </div>
            <div>
              <label className="text-[11px] text-[#6b6b6b] font-bold block mb-1">Jam & Menit Mulai</label>
              <div className="flex gap-1">
                <select
                  value={startHour}
                  onChange={(e) => setStartHour(Number(e.target.value))}
                  className="w-full h-10 px-2 bg-white border border-[#cccccc] text-xs font-bold text-[#262626]"
                >
                  {Array.from({ length: 24 }).map((_, i) => (
                    <option key={i} value={i}>
                      {String(i).padStart(2, '0')}:00 WIB
                    </option>
                  ))}
                </select>
              </div>
            </div>
            <div>
              <label className="text-[11px] text-[#6b6b6b] font-bold block mb-1">Durasi Pengisian</label>
              <select
                value={durationMinutes}
                onChange={(e) => setDurationMinutes(Number(e.target.value))}
                className="w-full h-10 px-2 bg-white border border-[#cccccc] text-xs font-bold text-[#262626]"
              >
                <option value={30}>30 Menit</option>
                <option value={60}>60 Menit (1 Jam)</option>
                <option value={90}>90 Menit (1.5 Jam)</option>
                <option value={120}>120 Menit (2 Jam)</option>
                <option value={180}>180 Menit (3 Jam)</option>
              </select>
            </div>
          </div>

          {/* Overlap Status Visual Indicator & Auto Suggestion */}
          {hasOverlap ? (
            <div className="p-3 bg-[#dc2626]/10 border border-[#dc2626] text-[#262626] text-xs space-y-2">
              <div className="flex items-center justify-between">
                <span className="flex items-center gap-2 font-bold text-[#dc2626]">
                  <ShieldAlert className="w-4 h-4 text-[#dc2626]" /> BENTROK JADWAL 24-JAM: Slot #{selectedSlot?.slotNumber} telah dipesan pada jam ini.
                </span>
                <span className="text-[10px] bg-[#f59e0b] px-2 py-0.5 font-bold text-white uppercase">
                  Waitlist System
                </span>
              </div>
              <p className="text-[11px] text-[#6b6b6b]">
                Sistem Anti-Overlap mencegah reservasi ganda pada slot yang sama. Anda dapat memilih jam lain yang berwarna hijau pada grafik 24-jam di atas atau bergabung ke antrean Waitlist.
              </p>
            </div>
          ) : (
            <div className="p-3 bg-[#22c55e]/10 border border-[#22c55e] text-[#22c55e] text-xs flex items-center justify-between font-bold">
              <span className="flex items-center gap-2">
                <CheckCircle2 className="w-4 h-4 text-[#22c55e]" /> Slot TERSEDIA (Bebas Bentrok: {String(startHour).padStart(2, '0')}:00 - {String(endTimeObj.getHours()).padStart(2, '0')}:{String(endTimeObj.getMinutes()).padStart(2, '0')} WIB)
              </span>
              <span className="text-[10px] bg-[#22c55e] px-2 py-0.5 text-white uppercase font-extrabold">
                Conflicted FREE
              </span>
            </div>
          )}
        </div>

        {/* Step 3: Tariff Line Items */}
        <div className="space-y-2 bg-[#ffffff] p-4 border border-[#e6e6e6] text-xs">
          <div className="flex items-center justify-between text-[#6b6b6b] font-bold border-b border-[#e6e6e6] pb-2">
            <span>Rincian Biaya Listrik (Permen ESDM No 7/2023)</span>
            <span className="text-[#1c69d4]">~{estimatedKwh.toFixed(1)} kWh</span>
          </div>

          <div className="space-y-1.5 pt-1 text-[#3c3c3c]">
            <div className="flex justify-between">
              <span>Listrik ({estimatedKwh.toFixed(1)} kWh x Rp {basePricePerKwh.toLocaleString('id-ID')})</span>
              <span>Rp {(estimatedKwh * basePricePerKwh).toLocaleString('id-ID')}</span>
            </div>
            <div className="flex justify-between">
              <span>Biaya Layanan SPKLU ({isUltraFast ? 'Ultra Fast' : 'Fast Charging'})</span>
              <span>Rp {serviceFee.toLocaleString('id-ID')}</span>
            </div>
            <div className="flex justify-between text-[#6b6b6b]">
              <span>PPN (11%)</span>
              <span>Rp {tax.toLocaleString('id-ID', { maximumFractionDigits: 0 })}</span>
            </div>
            <div className="flex justify-between text-base font-extrabold text-[#262626] pt-2 border-t border-[#e6e6e6]">
              <span>Total Estimasi</span>
              <span className="text-[#1c69d4]">Rp {grandTotal.toLocaleString('id-ID', { maximumFractionDigits: 0 })}</span>
            </div>
          </div>
        </div>

        {/* Actions */}
        <div className="flex flex-wrap items-center justify-end gap-3 pt-2">
          <button
            type="button"
            onClick={onClose}
            className="bmw-btn-secondary text-xs"
          >
            BATAL
          </button>
          {!hasOverlap && onConfirmAndStartSession && (
            <button
              type="button"
              onClick={handleConfirmAndStart}
              disabled={!isConnectorCompatible}
              className="bmw-btn-primary bg-[#22c55e] hover:bg-[#16a34a] text-white text-xs flex items-center gap-1.5"
            >
              <Zap className="w-4 h-4 fill-white" />
              <span>⚡ MULAI CAS SEKARANG</span>
            </button>
          )}
          <button
            type="button"
            onClick={handleAction}
            disabled={!isConnectorCompatible}
            className={`bmw-btn-primary text-xs ${
              !isConnectorCompatible
                ? 'opacity-50 cursor-not-allowed bg-[#d6d6d6]'
                : hasOverlap
                ? 'bg-[#f59e0b] hover:bg-[#d97706]'
                : ''
            }`}
          >
            <span>{hasOverlap ? 'MASUK ANTREAN WAITLIST' : 'RESERVASI JADWAL'}</span>
            <ArrowRight className="w-4 h-4" />
          </button>
        </div>
      </div>
    </div>
  );
};
