import React, { useState } from 'react';
import { ShieldCheck, Plus, MapPin, Activity, CheckCircle2, Link as LinkIcon, ExternalLink, Image as ImageIcon, Edit3, Lock } from 'lucide-react';
import type { Station, ChargerSlot } from '../types';

interface AdminPortalProps {
  stations: Station[];
  onUpdateSlotStatus: (stationId: string, slotId: string, newStatus: ChargerSlot['status']) => void;
  onAddStation: (stationData: { name: string; location: string; mapUrl?: string; imageUrl?: string; totalPower: number }) => Promise<void>;
  onUpdateStation?: (stationId: string, stationData: { name?: string; location?: string; mapUrl?: string; imageUrl?: string; totalPower?: number }) => Promise<void>;
}

const StationThumbnail: React.FC<{ imageUrl?: string; name: string }> = ({ imageUrl, name }) => {
  const [hasError, setHasError] = useState(false);

  if (hasError || !imageUrl) {
    return (
      <div className="w-16 h-12 bg-[#1a2129] border border-[#262626] flex items-center justify-center text-white shrink-0 shadow-xs">
        <Activity className="w-5 h-5 text-[#1c69d4]" />
      </div>
    );
  }

  return (
    <img
      src={imageUrl}
      alt={name}
      onError={() => setHasError(true)}
      className="w-16 h-12 object-cover border border-[#cccccc] shrink-0"
    />
  );
};

export const AdminPortal: React.FC<AdminPortalProps> = ({
  stations,
  onUpdateSlotStatus,
  onAddStation,
  onUpdateStation,
}) => {
  const [selectedStationId, setSelectedStationId] = useState<string>(
    stations[0]?.id || ''
  );
  const [isAddModalOpen, setIsAddModalOpen] = useState(false);
  const [isEditModalOpen, setIsEditModalOpen] = useState(false);

  // New Station Form State - Map Link & Image URL stored dynamically in DB
  const [name, setName] = useState('');
  const [location, setLocation] = useState('');
  const [mapUrl, setMapUrl] = useState('');
  const [imageUrl, setImageUrl] = useState('');
  const [totalPower, setTotalPower] = useState<number>(150);

  // Edit Station Form State
  const [editName, setEditName] = useState('');
  const [editLocation, setEditLocation] = useState('');
  const [editMapUrl, setEditMapUrl] = useState('');
  const [editImageUrl, setEditImageUrl] = useState('');
  const [editTotalPower, setEditTotalPower] = useState<number>(150);

  const selectedStation = stations.find((st) => st.id === selectedStationId) || stations[0];

  const handleOpenEditModal = () => {
    if (!selectedStation) return;
    setEditName(selectedStation.name);
    setEditLocation(selectedStation.location);
    setEditMapUrl(selectedStation.mapUrl || '');
    setEditImageUrl(selectedStation.imageUrl || '');
    setEditTotalPower(selectedStation.totalPower || 150);
    setIsEditModalOpen(true);
  };

  const handleCreateStation = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name || !location) return;

    await onAddStation({
      name,
      location,
      mapUrl: mapUrl || `https://www.google.com/maps/search/?api=1&query=${encodeURIComponent(name + ' ' + location)}`,
      imageUrl: imageUrl || '/images/spklu-mattoanging.png',
      totalPower,
    });

    setName('');
    setLocation('');
    setMapUrl('');
    setImageUrl('');
    setIsAddModalOpen(false);
  };

  const handleSaveEditStation = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedStation || !onUpdateStation) return;

    await onUpdateStation(selectedStation.id, {
      name: editName,
      location: editLocation,
      mapUrl: editMapUrl,
      imageUrl: editImageUrl,
      totalPower: editTotalPower,
    });

    setIsEditModalOpen(false);
  };

  return (
    <div className="space-y-6 max-w-6xl mx-auto animate-fadeIn">
      {/* Admin Banner */}
      <div className="bg-[#1a2129] border border-[#262626] p-6 sm:p-8 text-white relative overflow-hidden shadow-2xl">
        <div className="relative z-10 flex flex-col md:flex-row md:items-center justify-between gap-4">
          <div className="flex items-center gap-4">
            <div className="p-3 bg-[#1c69d4] text-white shadow-md">
              <ShieldCheck className="w-8 h-8" />
            </div>
            <div>
              <span className="px-2.5 py-0.5 text-[10px] font-extrabold bg-white/10 text-white uppercase tracking-wider">
                ADMIN CONTROL CENTER
              </span>
              <h1 className="text-2xl sm:text-3xl font-black text-white mt-1">
                SISTEM MANAJEMEN MASTER DATA SPKLU
              </h1>
              <p className="text-xs text-[#9a9a9a] font-light">
                Pengelolaan dinamis stasiun SPKLU Indonesia, pemantauan slot charger, gambar fisik tempat, link peta Google Maps, dan database PostgreSQL.
              </p>
            </div>
          </div>

          <button
            onClick={() => setIsAddModalOpen(true)}
            className="bmw-btn-primary text-xs shrink-0 flex items-center gap-2"
          >
            <Plus className="w-4 h-4" /> TAMBAH STASIUN SPKLU BARU
          </button>
        </div>
      </div>

      {/* Admin Metrics Cards */}
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
        <div className="bg-white p-5 border border-[#e6e6e6] shadow-sm">
          <span className="text-[10px] font-bold text-[#6b6b6b] uppercase tracking-wider">Total Stasiun Active</span>
          <div className="text-2xl font-black text-[#262626] mt-1">{stations.length} Lokasi</div>
          <span className="text-[11px] text-[#1c69d4] font-semibold">Resmi ESDM & PLN</span>
        </div>
        <div className="bg-white p-5 border border-[#e6e6e6] shadow-sm">
          <span className="text-[10px] font-bold text-[#6b6b6b] uppercase tracking-wider">Total Slot Charger</span>
          <div className="text-2xl font-black text-[#1c69d4] mt-1">
            {stations.reduce((acc, st) => acc + (st.slots ? st.slots.length : 0), 0)} Slot
          </div>
          <span className="text-[11px] text-[#22c55e] font-semibold">CCS2, CHAdeMO & Type 2</span>
        </div>
        <div className="bg-white p-5 border border-[#e6e6e6] shadow-sm">
          <span className="text-[10px] font-bold text-[#6b6b6b] uppercase tracking-wider">Database Backend</span>
          <div className="text-2xl font-black text-[#262626] mt-1">PostgreSQL DB</div>
          <span className="text-[11px] text-[#6b6b6b] font-light">GORM AutoMigrate Active</span>
        </div>
        <div className="bg-white p-5 border border-[#e6e6e6] shadow-sm">
          <span className="text-[10px] font-bold text-[#6b6b6b] uppercase tracking-wider">Status Integrasi</span>
          <div className="text-2xl font-black text-[#22c55e] mt-1">LIVE REAL DATA</div>
          <span className="text-[11px] text-[#22c55e] font-semibold">0 Dummy Mock Fallback</span>
        </div>
      </div>

      {/* Station Slots & Master Data Management */}
      <div className="bg-white border border-[#e6e6e6] p-6 sm:p-8 space-y-6 shadow-sm">
        <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4 border-b border-[#e6e6e6] pb-4">
          <div>
            <h2 className="text-lg font-bold text-[#262626]">Kelola Data Stasiun & Status Slot Charger</h2>
            <p className="text-xs text-[#6b6b6b] font-light">Ubah nama stasiun, alamat, link Google Maps, link foto fisik tempat, atau status slot operasional.</p>
          </div>

          <select
            value={selectedStationId}
            onChange={(e) => setSelectedStationId(e.target.value)}
            className="p-3 bg-[#fafafa] border border-[#cccccc] text-xs font-bold text-[#262626] focus:border-[#1c69d4]"
          >
            {stations.map((st) => (
              <option key={st.id} value={st.id}>
                {st.name} ({st.slots ? st.slots.length : 0} Slots)
              </option>
            ))}
          </select>
        </div>

        {selectedStation && selectedStation.slots ? (
          <div className="space-y-4">
            <div className="p-4 bg-[#fafafa] border border-[#e6e6e6] flex flex-col sm:flex-row sm:items-center justify-between gap-4">
              <div className="flex items-center gap-3">
                <StationThumbnail imageUrl={selectedStation.imageUrl} name={selectedStation.name} />
                <div>
                  <h4 className="font-extrabold text-[#262626] text-sm">{selectedStation.name}</h4>
                  <p className="text-xs text-[#6b6b6b]">{selectedStation.location}</p>
                </div>
              </div>

              <div className="flex items-center gap-2 self-start sm:self-auto">
                <button
                  onClick={handleOpenEditModal}
                  className="bmw-btn-primary text-xs flex items-center gap-1.5"
                >
                  <Edit3 className="w-3.5 h-3.5" /> Edit Data & Link Foto/Peta
                </button>
                <a
                  href={selectedStation.mapUrl || `https://www.google.com/maps/search/?api=1&query=${encodeURIComponent(selectedStation.name + ' ' + selectedStation.location)}`}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="bmw-btn-outline text-xs flex items-center gap-1.5"
                >
                  <ExternalLink className="w-3.5 h-3.5" /> Buka Google Maps
                </a>
              </div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {selectedStation.slots.map((slot) => (
                <div
                  key={slot.id}
                  className="bg-[#fafafa] p-4 border border-[#e6e6e6] space-y-3"
                >
                  <div className="flex items-center justify-between">
                    <span className="font-extrabold text-[#262626] text-sm">
                      Slot #{slot.slotNumber} ({slot.connectorType})
                    </span>
                    <span className="text-xs text-[#1c69d4] font-bold">{slot.maxPowerKw} kW</span>
                  </div>

                  <div className="space-y-1">
                    <label className="text-[10px] text-[#6b6b6b] uppercase font-bold tracking-wider block">Status Operasional</label>
                    {slot.status === 'IN_USE' ? (
                      <div className="p-2.5 bg-[#1c69d4]/10 border border-[#1c69d4]/30 text-[#1c69d4] text-[10px] font-bold flex items-center justify-between">
                        <span className="flex items-center gap-1.5">
                          <Lock className="w-3.5 h-3.5 shrink-0" />
                          <span>IN_USE (Terkunci Sesi Charging)</span>
                        </span>
                        <span className="text-[9px] bg-[#1c69d4] text-white px-1.5 py-0.5 uppercase tracking-wider font-extrabold">TERKUNCI</span>
                      </div>
                    ) : (
                      <select
                        value={slot.status}
                        onChange={(e) =>
                          onUpdateSlotStatus(selectedStation.id, slot.id, e.target.value as any)
                        }
                        className={`w-full p-2.5 text-xs font-bold border bg-white ${
                          slot.status === 'AVAILABLE'
                            ? 'border-[#22c55e] text-[#22c55e]'
                            : 'border-[#dc2626] text-[#dc2626]'
                        }`}
                      >
                        <option value="AVAILABLE">AVAILABLE (Tersedia)</option>
                        <option value="OUT_OF_SERVICE">OUT_OF_SERVICE (Pemeliharaan)</option>
                      </select>
                    )}
                  </div>
                </div>
              ))}
            </div>
          </div>
        ) : (
          <div className="p-4 bg-[#fafafa] text-xs text-[#6b6b6b]">Pilih stasiun di atas untuk mengelola slot charger.</div>
        )}
      </div>

      {/* Edit Station Modal */}
      {isEditModalOpen && selectedStation && (
        <div className="fixed inset-0 z-50 bg-black/75 backdrop-blur-xs flex items-center justify-center p-4">
          <div className="bg-white border border-[#262626] max-w-lg w-full p-6 sm:p-8 space-y-6 shadow-2xl">
            <div className="flex items-center justify-between border-b border-[#e6e6e6] pb-3">
              <div className="flex items-center gap-2">
                <Edit3 className="w-5 h-5 text-[#1c69d4]" />
                <h3 className="font-extrabold text-lg text-[#262626]">EDIT DATA STASIUN & LINK FOTO/PETA</h3>
              </div>
              <button
                onClick={() => setIsEditModalOpen(false)}
                className="text-[#9a9a9a] hover:text-[#262626] text-xl font-bold"
              >
                ✕
              </button>
            </div>

            <form onSubmit={handleSaveEditStation} className="space-y-4 text-xs">
              <div>
                <label className="text-[#262626] font-bold uppercase tracking-wider block mb-1">Nama Stasiun SPKLU</label>
                <div className="relative">
                  <Activity className="w-4 h-4 absolute left-3 top-3 text-[#9a9a9a]" />
                  <input
                    type="text"
                    required
                    value={editName}
                    onChange={(e) => setEditName(e.target.value)}
                    className="w-full h-11 pl-10 pr-3 bg-white border border-[#cccccc] text-[#262626] focus:border-[#1c69d4]"
                  />
                </div>
              </div>

              <div>
                <label className="text-[#262626] font-bold uppercase tracking-wider block mb-1">Alamat Lokasi Lengkap</label>
                <div className="relative">
                  <MapPin className="w-4 h-4 absolute left-3 top-3 text-[#9a9a9a]" />
                  <input
                    type="text"
                    required
                    value={editLocation}
                    onChange={(e) => setEditLocation(e.target.value)}
                    className="w-full h-11 pl-10 pr-3 bg-white border border-[#cccccc] text-[#262626] focus:border-[#1c69d4]"
                  />
                </div>
              </div>

              <div>
                <label className="text-[#262626] font-bold uppercase tracking-wider block mb-1">Link Lokasi Google Maps</label>
                <div className="relative">
                  <LinkIcon className="w-4 h-4 absolute left-3 top-3 text-[#9a9a9a]" />
                  <input
                    type="url"
                    value={editMapUrl}
                    onChange={(e) => setEditMapUrl(e.target.value)}
                    className="w-full h-11 pl-10 pr-3 bg-white border border-[#cccccc] text-[#262626] focus:border-[#1c69d4]"
                  />
                </div>
              </div>

              <div>
                <label className="text-[#262626] font-bold uppercase tracking-wider block mb-1">URL Gambar Fisik Stasiun (Image URL)</label>
                <div className="relative">
                  <ImageIcon className="w-4 h-4 absolute left-3 top-3 text-[#9a9a9a]" />
                  <input
                    type="text"
                    placeholder="/images/spklu-mattoanging.png atau https://..."
                    value={editImageUrl}
                    onChange={(e) => setEditImageUrl(e.target.value)}
                    className="w-full h-11 pl-10 pr-3 bg-white border border-[#cccccc] text-[#262626] focus:border-[#1c69d4]"
                  />
                </div>
                <p className="text-[10px] text-[#6b6b6b] mt-1">Anda dapat memasukkan link URL publik atau link lokal seperti `/images/spklu-mattoanging.png`</p>
              </div>

              <div>
                <label className="text-[#262626] font-bold uppercase tracking-wider block mb-1">Kapasitas Total Power (kW)</label>
                <input
                  type="number"
                  required
                  value={editTotalPower}
                  onChange={(e) => setEditTotalPower(parseFloat(e.target.value))}
                  className="w-full h-11 px-3 bg-white border border-[#cccccc] text-[#262626] focus:border-[#1c69d4]"
                />
              </div>

              <div className="pt-2 flex justify-end gap-2">
                <button
                  type="button"
                  onClick={() => setIsEditModalOpen(false)}
                  className="bmw-btn-outline text-xs"
                >
                  BATAL
                </button>
                <button type="submit" className="bmw-btn-primary text-xs flex items-center gap-2">
                  <CheckCircle2 className="w-4 h-4" /> SIMPAN PERUBAHAN
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Add New Station Modal */}
      {isAddModalOpen && (
        <div className="fixed inset-0 z-50 bg-black/75 backdrop-blur-xs flex items-center justify-center p-4">
          <div className="bg-white border border-[#262626] max-w-lg w-full p-6 sm:p-8 space-y-6 shadow-2xl">
            <div className="flex items-center justify-between border-b border-[#e6e6e6] pb-3">
              <div className="flex items-center gap-2">
                <Plus className="w-5 h-5 text-[#1c69d4]" />
                <h3 className="font-extrabold text-lg text-[#262626]">TAMBAH STASIUN SPKLU BARU</h3>
              </div>
              <button
                onClick={() => setIsAddModalOpen(false)}
                className="text-[#9a9a9a] hover:text-[#262626] text-xl font-bold"
              >
                ✕
              </button>
            </div>

            <form onSubmit={handleCreateStation} className="space-y-4 text-xs">
              <div>
                <label className="text-[#262626] font-bold uppercase tracking-wider block mb-1">Nama Stasiun SPKLU</label>
                <div className="relative">
                  <Activity className="w-4 h-4 absolute left-3 top-3 text-[#9a9a9a]" />
                  <input
                    type="text"
                    required
                    placeholder="Contoh: SPKLU PLN UID Makassar Panakkukang"
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                    className="w-full h-11 pl-10 pr-3 bg-white border border-[#cccccc] text-[#262626] focus:border-[#1c69d4]"
                  />
                </div>
              </div>

              <div>
                <label className="text-[#262626] font-bold uppercase tracking-wider block mb-1">Alamat Lokasi Lengkap</label>
                <div className="relative">
                  <MapPin className="w-4 h-4 absolute left-3 top-3 text-[#9a9a9a]" />
                  <input
                    type="text"
                    required
                    placeholder="Jl. Boulevard No.8, Panakkukang, Makassar"
                    value={location}
                    onChange={(e) => setLocation(e.target.value)}
                    className="w-full h-11 pl-10 pr-3 bg-white border border-[#cccccc] text-[#262626] focus:border-[#1c69d4]"
                  />
                </div>
              </div>

              <div>
                <label className="text-[#262626] font-bold uppercase tracking-wider block mb-1">Link Lokasi Google Maps</label>
                <div className="relative">
                  <LinkIcon className="w-4 h-4 absolute left-3 top-3 text-[#9a9a9a]" />
                  <input
                    type="url"
                    placeholder="https://maps.google.com/?q=SPKLU+Makassar"
                    value={mapUrl}
                    onChange={(e) => setMapUrl(e.target.value)}
                    className="w-full h-11 pl-10 pr-3 bg-white border border-[#cccccc] text-[#262626] focus:border-[#1c69d4]"
                  />
                </div>
              </div>

              <div>
                <label className="text-[#262626] font-bold uppercase tracking-wider block mb-1">URL Gambar Fisik Stasiun (Image URL)</label>
                <div className="relative">
                  <ImageIcon className="w-4 h-4 absolute left-3 top-3 text-[#9a9a9a]" />
                  <input
                    type="text"
                    placeholder="/images/spklu-mattoanging.png atau https://..."
                    value={imageUrl}
                    onChange={(e) => setImageUrl(e.target.value)}
                    className="w-full h-11 pl-10 pr-3 bg-white border border-[#cccccc] text-[#262626] focus:border-[#1c69d4]"
                  />
                </div>
              </div>

              <div>
                <label className="text-[#262626] font-bold uppercase tracking-wider block mb-1">Kapasitas Total Power (kW)</label>
                <input
                  type="number"
                  required
                  value={totalPower}
                  onChange={(e) => setTotalPower(parseFloat(e.target.value))}
                  className="w-full h-11 px-3 bg-white border border-[#cccccc] text-[#262626] focus:border-[#1c69d4]"
                />
              </div>

              <div className="pt-2 flex justify-end gap-2">
                <button
                  type="button"
                  onClick={() => setIsAddModalOpen(false)}
                  className="bmw-btn-outline text-xs"
                >
                  BATAL
                </button>
                <button type="submit" className="bmw-btn-primary text-xs flex items-center gap-2">
                  <CheckCircle2 className="w-4 h-4" /> SIMPAN STASIUN SPKLU
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};
