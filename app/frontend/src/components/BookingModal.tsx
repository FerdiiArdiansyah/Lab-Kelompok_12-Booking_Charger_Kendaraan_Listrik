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

const defaultVehiclesList: UserVehicle[] = [
  { id: 'vhc-001', userId: 'usr-driver', brand: 'BMW', model: 'iX xDrive50', licensePlate: 'B 8888 EV', connectorType: 'CCS2', batteryCapacityKwh: 105.2, createdAt: new Date().toISOString() },
  { id: 'vhc-002', userId: 'usr-driver', brand: 'Hyundai', model: 'Ioniq 5 Long Range', licensePlate: 'B 5678 EV', connectorType: 'CCS2', batteryCapacityKwh: 72.6, createdAt: new Date().toISOString() },
  { id: 'vhc-003', userId: 'usr-driver', brand: 'Wuling', model: 'Air EV Long Range', licensePlate: 'B 1234 EV', connectorType: 'Type 2', batteryCapacityKwh: 26.7, createdAt: new Date().toISOString() },
];

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
  const availableVehicles = vehicles && vehicles.length > 0 ? vehicles : defaultVehiclesList;

  const [selectedSlotId, setSelectedSlotId] = useState<string>(
    initialSlot ? initialSlot.id : station.slots[0]?.id || ''
  );
  const [selectedVehicleId, setSelectedVehicleId] = useState<string>(
    availableVehicles[0]?.id || ''
  );
  const [bookingDate, setBookingDate] = useState<string>(
    new Date().toISOString().split('T')[0]
  );
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
              <select
                value={selectedVehicleId}
                onChange={(e) => setSelectedVehicleId(e.target.value)}
                className="w-full h-12 px-3 bg-white border border-[#cccccc] text-xs font-bold text-[#262626] focus:border-[#1c69d4] focus:outline-none"
              >
                {availableVehicles.map((v) => (
                  <option key={v.id} value={v.id}>
                    {v.brand} {v.model} ({v.connectorType}) - {v.licensePlate}
                  </option>
                ))}
              </select>
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

        {/* Step 2: Date & Time */}
        <div className="space-y-3 bg-[#fafafa] p-4 border border-[#e6e6e6]">
          <div className="text-xs font-bold text-[#262626] uppercase tracking-wider flex items-center gap-2">
            <Clock className="w-4 h-4 text-[#1c69d4]" /> Waktu & Durasi Cas:
          </div>

          <div className="grid grid-cols-3 gap-3">
            <div>
              <label className="text-[11px] text-[#6b6b6b] font-bold block mb-1">Tanggal</label>
              <input
                type="date"
                value={bookingDate}
                onChange={(e) => setBookingDate(e.target.value)}
                className="w-full h-10 px-2 bg-white border border-[#cccccc] text-xs text-[#262626]"
              />
            </div>
            <div>
              <label className="text-[11px] text-[#6b6b6b] font-bold block mb-1">Jam Mulai</label>
              <select
                value={startHour}
                onChange={(e) => setStartHour(Number(e.target.value))}
                className="w-full h-10 px-2 bg-white border border-[#cccccc] text-xs text-[#262626]"
              >
                {Array.from({ length: 24 }).map((_, i) => (
                  <option key={i} value={i}>
                    {String(i).padStart(2, '0')}:00 WIB
                  </option>
                ))}
              </select>
            </div>
            <div>
              <label className="text-[11px] text-[#6b6b6b] font-bold block mb-1">Durasi</label>
              <select
                value={durationMinutes}
                onChange={(e) => setDurationMinutes(Number(e.target.value))}
                className="w-full h-10 px-2 bg-white border border-[#cccccc] text-xs text-[#262626]"
              >
                <option value={30}>30 Menit</option>
                <option value={60}>60 Menit (1 Jam)</option>
                <option value={90}>90 Menit</option>
                <option value={120}>120 Menit (2 Jam)</option>
              </select>
            </div>
          </div>

          {/* Overlap Status Visual Indicator */}
          {hasOverlap ? (
            <div className="p-3 bg-[#f59e0b]/10 border border-[#f59e0b] text-[#262626] text-xs flex items-center justify-between">
              <span className="flex items-center gap-2 font-bold text-[#262626]">
                <ShieldAlert className="w-4 h-4 text-[#f59e0b]" /> Slot telah terisi pada jam ini.
              </span>
              <span className="text-[10px] bg-[#f59e0b] px-2 py-0.5 font-bold text-white uppercase">
                Fitur Waitlist Aktif
              </span>
            </div>
          ) : (
            <div className="p-3 bg-[#22c55e]/10 border border-[#22c55e] text-[#22c55e] text-xs flex items-center gap-2 font-bold">
              <CheckCircle2 className="w-4 h-4 text-[#22c55e]" /> Slot TERSEDIA dan siap dipesan!
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
