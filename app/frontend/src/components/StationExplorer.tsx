import React, { useState } from 'react';
import { Search, MapPin, Zap, ChevronRight, SlidersHorizontal, Info, ExternalLink, Image as ImageIcon, Map as MapIcon } from 'lucide-react';
import type { Station, ChargerSlot, Booking } from '../types';

interface StationExplorerProps {
  stations: Station[];
  allBookings?: Booking[];
  onSelectStationForBooking: (station: Station, slot?: ChargerSlot) => void;
}

export const StationExplorer: React.FC<StationExplorerProps> = ({
  stations,
  allBookings = [],
  onSelectStationForBooking,
}) => {
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedConnector, setSelectedConnector] = useState<string>('ALL');
  const [selectedPowerFilter, setSelectedPowerFilter] = useState<string>('ALL');
  const [selectedStationForModal, setSelectedStationForModal] = useState<Station | null>(null);

  // Track view mode per station: 'map' (Google Maps iframe) or 'photo' (Local/Custom Image)
  const [viewModes, setViewModes] = useState<Record<string, 'map' | 'photo'>>({});

  const toggleViewMode = (stationId: string) => {
    setViewModes((prev) => ({
      ...prev,
      [stationId]: prev[stationId] === 'photo' ? 'map' : 'photo',
    }));
  };

  // Filter stations strictly based on search query, connector type, and power capacity
  const filteredStations = stations.filter((st) => {
    const query = searchQuery.trim().toLowerCase();
    const matchesSearch =
      query === '' ||
      st.name.toLowerCase().includes(query) ||
      st.location.toLowerCase().includes(query);

    const matchesConnector =
      selectedConnector === 'ALL' ||
      (st.slots && st.slots.some((sl) => sl.connectorType === selectedConnector));

    const totalPower = st.totalPower || 0;
    const matchesPower =
      selectedPowerFilter === 'ALL' ||
      (selectedPowerFilter === 'FAST' && totalPower >= 50 && totalPower <= 120) ||
      (selectedPowerFilter === 'ULTRA_FAST' && totalPower > 120);

    return matchesSearch && matchesConnector && matchesPower;
  });

  return (
    <div className="space-y-8 animate-fadeIn">
      {/* BMW Corporate Hero Band */}
      <div className="bg-[#1a2129] p-8 md:p-12 border-b border-[#e6e6e6] relative overflow-hidden text-white shadow-xl">
        <div className="max-w-4xl space-y-4 relative z-10">
          <div className="inline-flex items-center gap-2 px-3 py-1 bg-[#1c69d4] text-white text-xs font-bold uppercase tracking-wider">
            <Zap className="w-3.5 h-3.5 fill-white" /> Fast Charging Network Indonesia
          </div>

          <h1 className="text-3xl md:text-5xl font-black tracking-tight text-white leading-tight">
            TEMUKAN STASIUN SPKLU & PESAN KONEKTOR CHARGING
          </h1>

          <p className="text-[#bbbbbb] text-sm md:text-base font-light leading-relaxed max-w-2xl">
            Infrastruktur reservasi charger Kendaraan Listrik interaktif berbasis TDD & GORM. Bebas bentrok jadwal dengan Anti-Overlap Engine.
          </p>

          {/* Search Bar */}
          <div className="pt-4 max-w-2xl">
            <div className="relative">
              <Search className="w-5 h-5 absolute left-4 top-3.5 text-[#9a9a9a]" />
              <input
                type="text"
                placeholder="Cari lokasi SPKLU (Jakarta, Makassar, Parepare, Palopo...)"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="w-full pl-12 pr-4 h-12 bg-white text-[#262626] border border-white text-sm placeholder-[#9a9a9a] focus:outline-none focus:ring-2 focus:ring-[#1c69d4] shadow-md font-medium"
              />
            </div>
          </div>
        </div>
      </div>

      {/* Filter Row (Chips) */}
      <div className="bg-[#fafafa] border border-[#e6e6e6] p-3 sm:p-4 flex flex-col md:flex-row items-start md:items-center justify-between gap-3 shadow-xs">
        {/* Connector Filter */}
        <div className="w-full md:w-auto flex items-center gap-2 overflow-x-auto pb-1 md:pb-0 scrollbar-none">
          <span className="text-xs font-bold text-[#262626] uppercase tracking-wider flex items-center gap-1.5 shrink-0 mr-1">
            <SlidersHorizontal className="w-4 h-4 text-[#1c69d4]" /> Konektor:
          </span>
          {['ALL', 'CCS2', 'CHAdeMO', 'Type 2'].map((conn) => (
            <button
              key={conn}
              onClick={() => setSelectedConnector(conn)}
              className={`px-3 sm:px-4 py-1.5 sm:py-2 text-xs font-bold shrink-0 transition-all ${
                selectedConnector === conn
                  ? 'bg-[#1c69d4] text-white shadow-sm'
                  : 'bg-white text-[#262626] border border-[#cccccc] hover:border-[#1c69d4]'
              }`}
            >
              {conn === 'ALL' ? 'SEMUA KONEKTOR' : conn}
            </button>
          ))}
        </div>

        {/* Power Filter */}
        <div className="w-full md:w-auto flex items-center gap-2 overflow-x-auto pb-1 md:pb-0 scrollbar-none">
          <span className="text-xs font-bold text-[#6b6b6b] uppercase tracking-wider shrink-0">Daya:</span>
          {['ALL', 'FAST', 'ULTRA_FAST'].map((pwr) => (
            <button
              key={pwr}
              onClick={() => setSelectedPowerFilter(pwr)}
              className={`px-3 sm:px-4 py-1.5 sm:py-2 text-xs font-bold shrink-0 transition-all ${
                selectedPowerFilter === pwr
                  ? 'bg-[#262626] text-white shadow-sm'
                  : 'bg-white text-[#262626] border border-[#cccccc] hover:border-[#262626]'
              }`}
            >
              {pwr === 'ALL' ? 'SEMUA DAYA' : pwr === 'FAST' ? 'FAST (50-120 kW)' : 'ULTRA FAST (>120 kW)'}
            </button>
          ))}
        </div>
      </div>

      {/* Filter Stats Bar */}
      <div className="flex items-center justify-between text-xs text-[#6b6b6b] font-medium px-1">
        <span>Menampilkan <strong className="text-[#262626] font-bold">{filteredStations.length}</strong> dari {stations.length} Stasiun SPKLU</span>
        {(searchQuery || selectedConnector !== 'ALL' || selectedPowerFilter !== 'ALL') && (
          <button
            onClick={() => {
              setSearchQuery('');
              setSelectedConnector('ALL');
              setSelectedPowerFilter('ALL');
            }}
            className="text-[#1c69d4] hover:underline font-bold"
          >
            Reset Filter ✕
          </button>
        )}
      </div>

      {/* Station Cards Grid */}
      {filteredStations.length > 0 ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {filteredStations.map((station) => {
            const availableCount = station.slots ? station.slots.filter((s) => s.status === 'AVAILABLE').length : 0;
            const isFull = availableCount === 0;
            const googleMapUrl = station.mapUrl || `https://www.google.com/maps/search/?api=1&query=${encodeURIComponent(station.name + ' ' + station.location)}`;
            const currentMode = viewModes[station.id] || (station.imageUrl && station.imageUrl.startsWith('/images/') ? 'photo' : 'map');

            const embedUrl = `https://maps.google.com/maps?q=${encodeURIComponent(station.name + ' ' + station.location)}&t=m&z=15&ie=UTF8&iwloc=&output=embed`;

            return (
              <div
                key={station.id}
                className="bg-white border border-[#e6e6e6] hover:border-[#1c69d4] transition-all flex flex-col justify-between group overflow-hidden shadow-xs hover:shadow-md"
              >
                {/* Station Location Header - Dynamic Google Maps View */}
                <div className="relative h-56 bg-[#1a2129] overflow-hidden">
                  {currentMode === 'map' ? (
                    <iframe
                      title={`Peta ${station.name}`}
                      width="100%"
                      height="100%"
                      frameBorder="0"
                      scrolling="no"
                      src={embedUrl}
                      className="w-full h-full border-0 filter brightness-95 contrast-105 pointer-events-auto"
                    />
                  ) : (
                    <img
                      src={station.imageUrl || '/images/spklu-mattoanging.png'}
                      alt={station.name}
                      className="w-full h-full object-cover transition-transform duration-500 group-hover:scale-105"
                    />
                  )}

                  {/* Gradient Overlay bottom for readable title */}
                  <div className="absolute inset-x-0 bottom-0 h-24 bg-gradient-to-t from-black/90 via-black/40 to-transparent pointer-events-none" />

                  {/* Power & Status Overlay Badge */}
                  <div className="absolute top-3 left-3 flex items-center gap-2 pointer-events-auto">
                    <span className="px-2.5 py-1 text-[10px] font-extrabold bg-[#1c69d4] text-white uppercase tracking-wider shadow-sm flex items-center gap-1">
                      <Zap className="w-3 h-3 fill-white" />
                      {station.totalPower >= 150 ? 'ULTRA FAST CHARGING' : 'FAST CHARGING'}
                    </span>
                  </div>

                  <div className="absolute top-3 right-3 flex items-center gap-1.5 pointer-events-auto">
                    {/* View Switcher Toggle Button */}
                    <button
                      onClick={() => toggleViewMode(station.id)}
                      className="px-2.5 py-1 text-[10px] font-extrabold bg-black/80 hover:bg-black text-white border border-white/20 uppercase tracking-wider shadow-sm flex items-center gap-1 transition-colors"
                      title="Ganti Tampilan Google Maps / Foto Stasiun"
                    >
                      {currentMode === 'map' ? (
                        <>
                          <ImageIcon className="w-3 h-3 text-[#1c69d4]" /> Foto Tempat
                        </>
                      ) : (
                        <>
                          <MapIcon className="w-3 h-3 text-[#22c55e]" /> Google Maps
                        </>
                      )}
                    </button>

                    <span
                      className={`px-2.5 py-1 text-[10px] font-extrabold tracking-wider uppercase border shadow-sm ${
                        station.status === 'ACTIVE'
                          ? 'bg-[#22c55e] text-white border-[#22c55e]'
                          : 'bg-[#dc2626] text-white border-[#dc2626]'
                      }`}
                    >
                      {station.status}
                    </span>
                  </div>

                  {/* Title overlay on image/map */}
                  <div className="absolute bottom-3 left-3 right-3 pointer-events-none">
                    <h3 className="text-lg font-black text-white leading-snug drop-shadow-md">
                      {station.name}
                    </h3>
                  </div>
                </div>

                <div className="p-6 space-y-4 flex-1 flex flex-col justify-between">
                  <div className="space-y-4">
                    {/* Location Address & Google Maps Link */}
                    <div className="space-y-1.5">
                      <div className="flex items-start gap-2 text-[#3c3c3c] text-xs font-light leading-relaxed">
                        <MapPin className="w-4 h-4 text-[#1c69d4] shrink-0 mt-0.5" />
                        <span>{station.location}</span>
                      </div>

                      <a
                        href={googleMapUrl}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="inline-flex items-center gap-1.5 text-xs text-[#1c69d4] hover:text-[#0653b6] font-extrabold transition-colors pt-1"
                      >
                        <ExternalLink className="w-3.5 h-3.5" />
                        <span>📌 Buka Navigasi di Google Maps</span>
                      </a>
                    </div>

                    {/* Specs Box */}
                    <div className="grid grid-cols-2 gap-3 bg-[#fafafa] p-3 border border-[#e6e6e6] text-xs">
                      <div>
                        <span className="text-[#6b6b6b] block text-[10px] uppercase font-bold tracking-wider">Total Daya</span>
                        <span className="font-extrabold text-[#262626] text-sm">{station.totalPower} kW</span>
                      </div>
                      <div>
                        <span className="text-[#6b6b6b] block text-[10px] uppercase font-bold tracking-wider">Tarif ESDM</span>
                        <span className="font-extrabold text-[#1c69d4] text-sm">
                          Rp {station.activeTariff ? station.activeTariff.pricePerKwh.toLocaleString('id-ID') : '2.467'}/kWh
                        </span>
                      </div>
                    </div>

                    {/* Slots Breakdown */}
                    <div className="space-y-2">
                      <div className="flex items-center justify-between text-xs">
                        <span className="text-[#6b6b6b] font-semibold">Slot Charger ({station.slots ? station.slots.length : 0})</span>
                        <span className={`font-bold ${isFull ? 'text-[#f59e0b]' : 'text-[#22c55e]'}`}>
                          {availableCount} Slot Tersedia
                        </span>
                      </div>

                      <div className="flex flex-wrap gap-2">
                        {station.slots && station.slots.map((slot) => (
                          <div
                            key={slot.id}
                            onClick={() => onSelectStationForBooking(station, slot)}
                            className={`px-3 py-1.5 text-xs font-bold cursor-pointer transition-colors border flex items-center gap-1.5 ${
                              slot.status === 'AVAILABLE'
                                ? 'bg-white text-[#262626] border-[#cccccc] hover:border-[#1c69d4] hover:text-[#1c69d4]'
                                : 'bg-[#ebebeb] text-[#9a9a9a] border-[#e6e6e6] cursor-not-allowed'
                            }`}
                          >
                            <Zap className="w-3.5 h-3.5 text-[#1c69d4]" />
                            <span>Slot #{slot.slotNumber} ({slot.connectorType})</span>
                          </div>
                        ))}
                      </div>

                      {/* Dynamic 24-Hour Schedule Matrix Bar Preview from Backend Data */}
                      {(() => {
                        const getLocalDateStr = (d: Date = new Date()) => {
                          const year = d.getFullYear();
                          const month = String(d.getMonth() + 1).padStart(2, '0');
                          const day = String(d.getDate()).padStart(2, '0');
                          return `${year}-${month}-${day}`;
                        };
                        const todayStr = getLocalDateStr();

                        // Filter active bookings for this specific station
                        const stationBookings = (allBookings || []).filter(
                          (b) => b.stationId === station.id && b.status !== 'CANCELLED' && b.status !== 'EXPIRED' && b.status !== 'COMPLETED'
                        );

                        // Calculate 24-hour occupancy for this station
                        const hourlyOccupancy = Array.from({ length: 24 }).map((_, h) => {
                          const hourSlotStart = new Date(`${todayStr}T${String(h).padStart(2, '0')}:00:00`).getTime();
                          const hourSlotEnd = hourSlotStart + 3600000;

                          return stationBookings.some((b) => {
                            const bStart = new Date(b.startTime);
                            const bEnd = new Date(b.endTime);

                            const startHour = bStart.getHours();
                            let endHour = bEnd.getHours();
                            if (endHour <= startHour && bEnd.getTime() > bStart.getTime()) {
                              endHour = startHour + 1;
                            }

                            const isHourMatch = (h >= startHour && h < endHour);
                            const isTimeOverlap = hourSlotStart < bEnd.getTime() && hourSlotEnd > bStart.getTime();

                            return isHourMatch || isTimeOverlap;
                          });
                        });

                        const occupiedHoursList: number[] = [];
                        const emptyHoursList: number[] = [];
                        hourlyOccupancy.forEach((isOccupied, h) => {
                          if (isOccupied) occupiedHoursList.push(h);
                          else emptyHoursList.push(h);
                        });

                        const occupiedSummary = occupiedHoursList.length > 0
                          ? occupiedHoursList.map((h) => `${String(h).padStart(2, '0')}:00`).join(', ')
                          : 'Tidak Ada (Semua Jam Kosong)';

                        const emptySummary = emptyHoursList.length === 24
                          ? '24 Jam Bebas (Semua Kosong)'
                          : emptyHoursList.length > 0
                          ? emptyHoursList.map((h) => `${String(h).padStart(2, '0')}:00`).join(', ')
                          : 'Tidak Ada (Semua Jam Terisi)';

                        return (
                          <div className="pt-2 bg-[#fafafa] p-2.5 border border-[#e6e6e6] space-y-1.5">
                            {/* 24 Hour Blocks Grid with printed hour numbers */}

                            {/* 24 Hour Blocks Grid with printed hour numbers */}
                            <div className="grid grid-cols-12 sm:grid-cols-24 gap-0.5 bg-white p-1 border border-[#cccccc]">
                              {hourlyOccupancy.map((isHourOccupied, h) => (
                                <div
                                  key={h}
                                  title={`Jam ${String(h).padStart(2, '0')}:00 WIB (${isHourOccupied ? 'TERISI / BENTROK' : 'KOSONG / TERSEDIA'})`}
                                  className={`h-6 text-[9px] font-extrabold flex items-center justify-center transition-colors ${
                                    isHourOccupied
                                      ? 'bg-red-500 text-white border border-red-600'
                                      : 'bg-[#22c55e] text-white border border-emerald-600'
                                  }`}
                                >
                                  {String(h).padStart(2, '0')}
                                </div>
                              ))}
                            </div>

                            {/* Explicit Empty & Occupied Hours Summary Text */}
                            <div className="space-y-0.5 pt-0.5 text-[10px]">
                              <div className="text-emerald-700 font-extrabold truncate">
                                🟢 Jam Kosong: {emptySummary}
                              </div>
                              <div className="text-red-600 font-extrabold truncate">
                                🔴 Jam Terisi: {occupiedSummary}
                              </div>
                            </div>
                          </div>
                        );
                      })()}
                    </div>
                  </div>

                  {/* Card Action CTAs */}
                  <div className="pt-4 border-t border-[#e6e6e6] flex items-center justify-between gap-3 mt-4">
                    <button
                      onClick={() => onSelectStationForBooking(station)}
                      className="bmw-btn-primary flex-1 text-xs"
                    >
                      <span>{isFull ? 'GABUNG WAITLIST' : 'BOOKING SLOT'}</span>
                      <ChevronRight className="w-4 h-4" />
                    </button>
                    <button
                      onClick={() => setSelectedStationForModal(station)}
                      className="p-3 bg-[#fafafa] hover:bg-[#ebebeb] text-[#262626] border border-[#cccccc] transition-colors"
                      title="Lihat Detail Peta Interaktif & Lokasi"
                    >
                      <Info className="w-4 h-4" />
                    </button>
                  </div>
                </div>
              </div>
            );
          })}
        </div>
      ) : (
        <div className="p-12 text-center bg-[#fafafa] border border-[#e6e6e6] space-y-3">
          <Zap className="w-12 h-12 text-[#9a9a9a] mx-auto" />
          <h3 className="text-lg font-bold text-[#262626]">Stasiun SPKLU Tidak Ditemukan</h3>
          <p className="text-xs text-[#6b6b6b]">Tidak ada stasiun yang cocok dengan filter pencarian "{searchQuery}". Silakan ubah filter atau kata kunci Anda.</p>
        </div>
      )}

      {/* Map & Image Modal */}
      {selectedStationForModal && (
        <div className="fixed inset-0 z-50 bg-black/75 backdrop-blur-xs flex items-center justify-center p-4">
          <div className="bg-white border border-[#262626] max-w-2xl w-full p-6 space-y-6 relative shadow-2xl">
            <div className="flex items-center justify-between border-b border-[#e6e6e6] pb-4">
              <div>
                <h2 className="text-xl font-bold text-[#262626]">{selectedStationForModal.name}</h2>
                <p className="text-xs text-[#6b6b6b] flex items-center gap-1 mt-1">
                  <MapPin className="w-3.5 h-3.5 text-[#1c69d4]" /> {selectedStationForModal.location}
                </p>
              </div>
              <button
                onClick={() => setSelectedStationForModal(null)}
                className="text-[#6b6b6b] hover:text-[#262626] p-2 text-xl font-bold"
              >
                ✕
              </button>
            </div>

            {/* Dynamic Interactive Google Maps Embed Modal View */}
            <div className="relative h-80 bg-[#1a2129] border border-[#262626] overflow-hidden">
              <iframe
                title={`Detail Peta ${selectedStationForModal.name}`}
                width="100%"
                height="100%"
                frameBorder="0"
                src={`https://maps.google.com/maps?q=${encodeURIComponent(selectedStationForModal.name + ' ' + selectedStationForModal.location)}&t=m&z=15&ie=UTF8&iwloc=&output=embed`}
                className="w-full h-full border-0"
              />
              <div className="absolute bottom-3 right-3 z-10">
                <a
                  href={selectedStationForModal.mapUrl || `https://www.google.com/maps/search/?api=1&query=${encodeURIComponent(selectedStationForModal.name + ' ' + selectedStationForModal.location)}`}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="bmw-btn-primary text-xs flex items-center gap-2 shadow-lg"
                >
                  <ExternalLink className="w-4 h-4" /> BUKA GOOGLE MAPS NAVIGASI
                </a>
              </div>
            </div>

            <div className="flex justify-end gap-3 pt-4 border-t border-[#e6e6e6]">
              <button
                onClick={() => setSelectedStationForModal(null)}
                className="bmw-btn-outline text-xs"
              >
                TUTUP
              </button>
              <button
                onClick={() => {
                  const st = selectedStationForModal;
                  setSelectedStationForModal(null);
                  onSelectStationForBooking(st);
                }}
                className="bmw-btn-primary text-xs flex items-center gap-2"
              >
                <span>RESERVASI SLOT</span>
                <ChevronRight className="w-4 h-4" />
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};
