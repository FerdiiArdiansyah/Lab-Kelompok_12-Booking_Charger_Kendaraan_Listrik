import React from 'react';
import { Clock, Play, Zap, Plus, Calendar } from 'lucide-react';
import type { Booking, Station, UserVehicle } from '../types';

interface UserBookingsListProps {
  bookings: Booking[];
  stations: Station[];
  vehicles: UserVehicle[];
  onStartSessionFromBooking: (booking: Booking) => void;
  onCancelBooking: (bookingId: string) => void;
  onNavigateToFinder: () => void;
}

export const UserBookingsList: React.FC<UserBookingsListProps> = ({
  bookings,
  stations,
  vehicles,
  onStartSessionFromBooking,
  onCancelBooking,
  onNavigateToFinder,
}) => {
  return (
    <div className="space-y-6 max-w-5xl mx-auto">
      {/* Header */}
      <div className="bmw-card p-6 sm:p-8 flex flex-col sm:flex-row items-center justify-between gap-4">
        <div>
          <span className="px-2.5 py-0.5 text-[10px] font-extrabold bg-[#1c69d4] text-white uppercase tracking-wider">
            RESERVATION MANAGEMENT
          </span>
          <h1 className="text-2xl sm:text-4xl font-black text-[#262626] mt-1">BOOKING SAYA</h1>
          <p className="text-xs text-[#6b6b6b] font-light">Daftar jadwal reservasi charger SPKLU Anda. Lakukan Check-In untuk memulai sesi pengisian.</p>
        </div>
        <button
          onClick={onNavigateToFinder}
          className="bmw-btn-primary shrink-0 text-xs flex items-center gap-2"
        >
          <Plus className="w-4 h-4" />
          <span>BUAT RESERVASI SPKLU</span>
        </button>
      </div>

      {/* Bookings Grid */}
      <div className="space-y-4">
        {bookings.length === 0 ? (
          <div className="bmw-card p-12 text-center space-y-4 border-dashed">
            <Calendar className="w-12 h-12 text-[#1c69d4] mx-auto" />
            <div>
              <h3 className="text-lg font-bold text-[#262626]">Belum Ada Reservasi Aktif</h3>
              <p className="text-xs text-[#6b6b6b] max-w-md mx-auto mt-1">
                Anda belum memiliki jadwal reservasi charger. Silakan pilih stasiun SPKLU terdekat di menu <strong>Cari SPKLU</strong> untuk memesan slot pengisian daya.
              </p>
            </div>
            <button
              onClick={onNavigateToFinder}
              className="bmw-btn-primary text-xs mx-auto flex items-center gap-2"
            >
              <Plus className="w-4 h-4" />
              <span>CARI SPKLU & BUAT RESERVASI</span>
            </button>
          </div>
        ) : (
          bookings.map((bkg) => {
            const station = stations.find((st) => st.id === bkg.stationId);
            const vehicle = vehicles.find((v) => v.id === bkg.vehicleId);

            return (
              <div
                key={bkg.id}
                className="bmw-card p-6 flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4 hover:border-[#1c69d4] transition-colors"
              >
                <div className="space-y-2">
                  <div className="flex items-center gap-2">
                    <span className="text-xs font-mono font-extrabold text-[#1c69d4]">#{bkg.id}</span>
                    <span
                      className={`px-2.5 py-0.5 text-[10px] font-extrabold uppercase border ${
                        bkg.status === 'CONFIRMED'
                          ? 'bg-[#1c69d4]/10 text-[#1c69d4] border-[#1c69d4]/30'
                          : bkg.status === 'IN_PROGRESS'
                          ? 'bg-[#22c55e]/10 text-[#22c55e] border-[#22c55e]/30'
                          : 'bg-[#ebebeb] text-[#6b6b6b] border-[#cccccc]'
                      }`}
                    >
                      {bkg.status}
                    </span>
                  </div>

                  <h3 className="text-lg font-bold text-[#262626]">
                    {station ? station.name : 'Stasiun SPKLU'}
                  </h3>

                  <div className="flex flex-wrap items-center gap-x-4 gap-y-1 text-xs text-[#3c3c3c]">
                    <span className="flex items-center gap-1 font-semibold">
                      <Clock className="w-3.5 h-3.5 text-[#1c69d4]" />
                      {new Date(bkg.startTime).toLocaleString('id-ID', { dateStyle: 'short', timeStyle: 'short' })}
                    </span>
                    <span className="flex items-center gap-1 text-[#6b6b6b]">
                      <Zap className="w-3.5 h-3.5 text-[#f59e0b]" />
                      Slot ID: {bkg.slotId}
                    </span>
                    {vehicle && (
                      <span className="text-[#6b6b6b] font-light">
                        {vehicle.brand} {vehicle.model} ({vehicle.licensePlate})
                      </span>
                    )}
                  </div>
                </div>

                {/* Actions */}
                <div className="flex items-center gap-2 shrink-0">
                  {bkg.status === 'CONFIRMED' && (
                    <button
                      onClick={() => onStartSessionFromBooking(bkg)}
                      className="bmw-btn-primary text-xs"
                    >
                      <Play className="w-3.5 h-3.5 fill-white" />
                      <span>CHECK-IN & START CAS</span>
                    </button>
                  )}
                  {(bkg.status === 'CONFIRMED' || bkg.status === 'PENDING') && (
                    <button
                      onClick={() => onCancelBooking(bkg.id)}
                      className="bmw-btn-secondary text-xs"
                    >
                      BATALKAN
                    </button>
                  )}
                </div>
              </div>
            );
          })
        )}
      </div>
    </div>
  );
};
