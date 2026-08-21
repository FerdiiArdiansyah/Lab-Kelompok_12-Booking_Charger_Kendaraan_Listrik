import React, { useState } from 'react';
import { Car, Plus, Trash2, ShieldCheck } from 'lucide-react';
import type { UserVehicle } from '../types';

interface UserGarageProps {
  vehicles: UserVehicle[];
  onAddVehicle: (vehicle: Omit<UserVehicle, 'id' | 'createdAt'>) => void;
  onDeleteVehicle: (id: string) => void;
}

export const UserGarage: React.FC<UserGarageProps> = ({
  vehicles,
  onAddVehicle,
  onDeleteVehicle,
}) => {
  const [showAddForm, setShowAddForm] = useState(false);
  const [brand, setBrand] = useState('');
  const [model, setModel] = useState('');
  const [licensePlate, setLicensePlate] = useState('');
  const [connectorType, setConnectorType] = useState<UserVehicle['connectorType']>('CCS2');
  const [batteryCapacityKwh, setBatteryCapacityKwh] = useState<number>(60.0);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onAddVehicle({
      userId: '',
      brand,
      model,
      licensePlate,
      connectorType,
      batteryCapacityKwh,
    });
    setBrand('');
    setModel('');
    setLicensePlate('');
    setBatteryCapacityKwh(60.0);
    setShowAddForm(false);
  };

  return (
    <div className="space-y-6 max-w-5xl mx-auto">
      {/* Header */}
      <div className="bmw-card p-6 sm:p-8 flex flex-col sm:flex-row items-center justify-between gap-4">
        <div>
          <span className="px-2.5 py-0.5 text-[10px] font-extrabold bg-[#1c69d4] text-white uppercase tracking-wider">
            FLEET VEHICLE REGISTRY
          </span>
          <h1 className="text-2xl sm:text-4xl font-black text-[#262626] mt-1">GARASI KENDARAAN EV</h1>
          <p className="text-xs text-[#6b6b6b] font-light">Kelola informasi mobil listrik Anda untuk pemetaan kompatibilitas otomatis saat booking.</p>
        </div>
        <button
          onClick={() => setShowAddForm(true)}
          className="bmw-btn-primary shrink-0 text-xs"
        >
          <Plus className="w-4 h-4" />
          <span>TAMBAH KENDARAAN EV</span>
        </button>
      </div>

      {/* Vehicles Grid / Empty State */}
      {vehicles.length === 0 ? (
        <div className="bmw-card p-12 text-center space-y-3 border-dashed">
          <Car className="w-12 h-12 text-[#9a9a9a] mx-auto" />
          <h3 className="text-lg font-bold text-[#262626]">Belum Ada Kendaraan EV Terdaftar</h3>
          <p className="text-xs text-[#6b6b6b] max-w-md mx-auto">
            Garasi EV Anda masih kosong. Tambahkan mobil listrik Anda untuk pemetaan pengisian daya dan kustomisasi reservasi stasiun SPKLU.
          </p>
          <button
            onClick={() => setShowAddForm(true)}
            className="bmw-btn-primary text-xs mx-auto mt-2"
          >
            <Plus className="w-4 h-4" />
            <span>TAMBAH KENDARAAN EV</span>
          </button>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {vehicles.map((v) => (
            <div
              key={v.id}
              className="bmw-card p-6 space-y-4 hover:border-[#1c69d4] transition-colors group flex flex-col justify-between"
            >
              <div className="space-y-3">
                <div className="flex items-start justify-between">
                  <div className="p-3 bg-[#fafafa] border border-[#e6e6e6] text-[#1c69d4]">
                    <Car className="w-6 h-6" />
                  </div>
                  <button
                    onClick={() => onDeleteVehicle(v.id)}
                    className="p-2 text-[#9a9a9a] hover:text-[#dc2626] transition-colors"
                    title="Hapus Kendaraan"
                  >
                    <Trash2 className="w-4 h-4" />
                  </button>
                </div>

                <div>
                  <span className="text-[10px] font-bold uppercase text-[#6b6b6b] tracking-wider">{v.brand}</span>
                  <h3 className="text-xl font-bold text-[#262626] group-hover:text-[#1c69d4] transition-colors">
                    {v.model}
                  </h3>
                  <span className="inline-block mt-1 px-2.5 py-0.5 bg-[#fafafa] border border-[#cccccc] text-xs font-mono font-bold text-[#262626]">
                    {v.licensePlate}
                  </span>
                </div>

                <div className="grid grid-cols-2 gap-3 bg-[#fafafa] p-3 border border-[#e6e6e6] text-xs">
                  <div>
                    <span className="text-[#6b6b6b] text-[10px] block font-bold uppercase tracking-wider">Konektor</span>
                    <span className="font-extrabold text-[#1c69d4]">{v.connectorType}</span>
                  </div>
                  <div>
                    <span className="text-[#6b6b6b] text-[10px] block font-bold uppercase tracking-wider">Baterai</span>
                    <span className="font-extrabold text-[#262626]">{v.batteryCapacityKwh} kWh</span>
                  </div>
                </div>
              </div>

              <div className="pt-3 border-t border-[#e6e6e6] text-[11px] text-[#6b6b6b] flex items-center gap-1.5 font-light">
                <ShieldCheck className="w-3.5 h-3.5 text-[#22c55e]" /> Terverifikasi oleh User Microservice
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Add Vehicle Modal Form */}
      {showAddForm && (
        <div className="fixed inset-0 z-50 bg-black/60 backdrop-blur-xs flex items-center justify-center p-4">
          <form
            onSubmit={handleSubmit}
            className="bg-white border border-[#262626] max-w-md w-full p-6 space-y-5 relative shadow-2xl"
          >
            <div className="flex items-center justify-between pb-3 border-b border-[#e6e6e6]">
              <h2 className="text-lg font-bold text-[#262626] flex items-center gap-2">
                <Car className="w-5 h-5 text-[#1c69d4]" /> REGISTRASI EV BARU
              </h2>
              <button
                type="button"
                onClick={() => setShowAddForm(false)}
                className="text-[#6b6b6b] hover:text-[#262626]"
              >
                ✕
              </button>
            </div>

            <div className="space-y-3 text-xs">
              <div>
                <label className="text-[#262626] font-bold uppercase tracking-wider block mb-1">Merek (Brand)</label>
                <input
                  type="text"
                  required
                  placeholder="Contoh: Hyundai, Wuling, BMW, Tesla"
                  value={brand}
                  onChange={(e) => setBrand(e.target.value)}
                  className="w-full h-11 px-3 bg-white border border-[#cccccc] text-[#262626]"
                />
              </div>

              <div>
                <label className="text-[#262626] font-bold uppercase tracking-wider block mb-1">Model / Varian</label>
                <input
                  type="text"
                  required
                  placeholder="Contoh: Ioniq 5 Signature Long Range"
                  value={model}
                  onChange={(e) => setModel(e.target.value)}
                  className="w-full h-11 px-3 bg-white border border-[#cccccc] text-[#262626]"
                />
              </div>

              <div>
                <label className="text-[#262626] font-bold uppercase tracking-wider block mb-1">Nomor Plat Polisi</label>
                <input
                  type="text"
                  required
                  placeholder="Contoh: B 1234 EV"
                  value={licensePlate}
                  onChange={(e) => setLicensePlate(e.target.value)}
                  className="w-full h-11 px-3 bg-white border border-[#cccccc] text-[#262626]"
                />
              </div>

              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="text-[#262626] font-bold uppercase tracking-wider block mb-1">Tipe Konektor</label>
                  <select
                    value={connectorType}
                    onChange={(e) => setConnectorType(e.target.value as any)}
                    className="w-full h-11 px-3 bg-white border border-[#cccccc] text-[#262626]"
                  >
                    <option value="CCS2">CCS2</option>
                    <option value="CHAdeMO">CHAdeMO</option>
                    <option value="Type 2">Type 2</option>
                    <option value="GB/T">GB/T</option>
                  </select>
                </div>
                <div>
                  <label className="text-[#262626] font-bold uppercase tracking-wider block mb-1">Baterai (kWh)</label>
                  <input
                    type="number"
                    step="0.1"
                    required
                    value={batteryCapacityKwh}
                    onChange={(e) => setBatteryCapacityKwh(Number(e.target.value))}
                    className="w-full h-11 px-3 bg-white border border-[#cccccc] text-[#262626]"
                  />
                </div>
              </div>
            </div>

            <div className="flex justify-end gap-3 pt-3 border-t border-[#e6e6e6]">
              <button
                type="button"
                onClick={() => setShowAddForm(false)}
                className="bmw-btn-secondary"
              >
                BATAL
              </button>
              <button
                type="submit"
                className="bmw-btn-primary"
              >
                SIMPAN MOBIL EV
              </button>
            </div>
          </form>
        </div>
      )}
    </div>
  );
};
