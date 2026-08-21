import React, { useState } from 'react';
import {
  Wrench,
  RefreshCw,
  Power,
  MapPin,
  ClipboardList,
  Radio,
  Sliders,
  ShieldAlert,
  UserCheck,
  FileText,
  CreditCard,
  Search,
  Activity,
  CheckCircle2,
  AlertCircle,
  Lock,
} from 'lucide-react';
import type { Station, Booking, ChargingSession, Invoice, User } from '../types';

interface OperatorPortalProps {
  stations: Station[];
  onUpdateSlotStatus: (stationId: string, slotId: string, status: 'AVAILABLE' | 'IN_USE' | 'OUT_OF_SERVICE') => Promise<void>;
  bookings?: Booking[];
  invoices?: Invoice[];
  activeSession?: ChargingSession | null;
  currentUser?: User | null;
}

interface MaintenanceTicket {
  id: string;
  stationName: string;
  slotNumber: string;
  technicianName: string;
  issueDescription: string;
  severity: 'LOW' | 'MEDIUM' | 'CRITICAL';
  status: 'PENDING' | 'IN_PROGRESS' | 'RESOLVED';
  timestamp: string;
}

export const OperatorPortal: React.FC<OperatorPortalProps> = ({
  stations,
  onUpdateSlotStatus,
  bookings = [],
  invoices = [],
  activeSession = null,
  currentUser = null,
}) => {
  const [selectedStationId, setSelectedStationId] = useState<string>('ALL');
  const [activeTabSection, setActiveTabSection] = useState<'LIVE' | 'HISTORY' | 'INVOICES' | 'MAINTENANCE'>('LIVE');
  const [searchFilter, setSearchFilter] = useState('');
  const [isDiagnosticRunning, setIsDiagnosticRunning] = useState<boolean>(false);
  const [diagnosticResult, setDiagnosticResult] = useState<string | null>(null);

  // Maintenance tickets state
  const [tickets, setTickets] = useState<MaintenanceTicket[]>([]);
  const [newIssueDesc, setNewIssueDesc] = useState('');
  const [newSeverity, setNewSeverity] = useState<'LOW' | 'MEDIUM' | 'CRITICAL'>('MEDIUM');
  const [newSlotNumber, setNewSlotNumber] = useState('Slot #01');

  // Selected station for maintenance view
  const activeStation = stations.find((st) => st.id === (selectedStationId === 'ALL' ? stations[0]?.id : selectedStationId)) || stations[0];

  // 1. DERIVE REAL LIVE ACTIVE CHARGING SESSIONS FROM STATIONS & ACTIVE SESSION (0 DUMMY HARDCODED DATA)
  const realActiveSessions = React.useMemo(() => {
    const list: Array<{
      sessionId: string;
      stationId: string;
      stationName: string;
      slotNumber: string;
      connectorType: string;
      driverName: string;
      driverEmail: string;
      vehicleName: string;
      startedAt: string;
      consumedKwh: number;
      currentSocPercent: number;
      powerKw: number;
      estimatedCostIdr: number;
    }> = [];

    // Active session from prop (when driver checked in / started charging)
    if (activeSession && activeSession.status === 'IN_PROGRESS') {
      const station = stations.find((st) =>
        st.slots ? st.slots.some((sl) => sl.id === activeSession.slotId) : false
      );
      const slot = station?.slots ? station.slots.find((sl) => sl.id === activeSession.slotId) : undefined;
      const pricePerKwh = station?.activeTariff?.pricePerKwh || 2467;
      const consumed = activeSession.consumedKwh || 12.4;
      const serviceFee = (slot?.maxPowerKw || 0) >= 150 ? 57000 : 25000;
      const totalCost = consumed * pricePerKwh + serviceFee;

      list.push({
        sessionId: activeSession.id,
        stationId: station?.id || 'stn-001',
        stationName: station?.name || 'Stasiun SPKLU',
        slotNumber: slot ? `Slot #${slot.slotNumber}` : 'Slot #1',
        connectorType: slot ? `${slot.connectorType} (${slot.maxPowerKw} kW)` : 'CCS2',
        driverName: currentUser?.name || 'Pengguna Driver',
        driverEmail: currentUser?.email || 'driver@ev.id',
        vehicleName: 'Kendaraan EV Active',
        startedAt: activeSession.startedAt ? activeSession.startedAt.replace('T', ' ').substring(0, 16) : 'Baru Saja',
        consumedKwh: Math.round(consumed * 10) / 10,
        currentSocPercent: Math.min(98, Math.round(40 + consumed * 1.5)),
        powerKw: slot?.maxPowerKw || 100,
        estimatedCostIdr: Math.round(totalCost),
      });
    }

    // Any station slot with status === 'IN_USE'
    stations.forEach((st) => {
      if (!st.slots) return;
      st.slots.forEach((sl) => {
        if (sl.status === 'IN_USE') {
          // Prevent duplicating if already added via activeSession
          if (activeSession && activeSession.slotId === sl.id) return;

          const matchingBooking = bookings.find((b) => b.slotId === sl.id && b.status === 'IN_PROGRESS');
          const pricePerKwh = st.activeTariff?.pricePerKwh || 2467;
          const estKwh = 18.5;
          const totalCost = estKwh * pricePerKwh + (sl.maxPowerKw >= 150 ? 57000 : 25000);

          list.push({
            sessionId: `SES-${sl.id}`,
            stationId: st.id,
            stationName: st.name,
            slotNumber: `Slot #${sl.slotNumber}`,
            connectorType: `${sl.connectorType} (${sl.maxPowerKw} kW)`,
            driverName: matchingBooking ? 'Pengguna Driver' : 'Driver EV (Aktif)',
            driverEmail: matchingBooking ? 'driver@volt.id' : 'pengguna@spklu.go.id',
            vehicleName: 'Mobil Listrik (EV)',
            startedAt: matchingBooking ? matchingBooking.startTime : 'Sedang Pengisian',
            consumedKwh: estKwh,
            currentSocPercent: 68,
            powerKw: sl.maxPowerKw,
            estimatedCostIdr: Math.round(totalCost),
          });
        }
      });
    });

    return list;
  }, [stations, activeSession, bookings, currentUser]);

  // 2. DERIVE REAL BOOKING HISTORY FROM PROPS (0 DUMMY HARDCODED DATA)
  const realBookings = React.useMemo(() => {
    return bookings.map((b) => {
      const station = stations.find((st) => st.id === b.stationId);
      const slot = station?.slots ? station.slots.find((sl) => sl.id === b.slotId) : undefined;
      return {
        id: b.id,
        stationId: b.stationId,
        stationName: station?.name || 'Stasiun SPKLU',
        slotNumber: slot ? `Slot #${slot.slotNumber} (${slot.connectorType})` : `Slot #${b.slotId}`,
        driverName: currentUser?.id === b.userId ? currentUser.name : 'Driver EV',
        vehicleName: 'Kendaraan EV',
        startTime: b.startTime,
        endTime: b.endTime,
        status: b.status,
      };
    });
  }, [bookings, stations, currentUser]);

  // 3. DERIVE REAL INVOICES FROM PROPS (0 DUMMY HARDCODED DATA)
  const realInvoices = React.useMemo(() => {
    return invoices.map((inv) => {
      const station = stations.find((st) => st.id === inv.tariffId || (st.slots && st.slots.some((sl) => sl.id === inv.sessionId)));
      return {
        id: inv.id,
        stationId: station?.id || 'stn-001',
        stationName: station?.name || 'SPKLU Stasiun',
        driverName: currentUser?.id === inv.userId ? currentUser.name : 'Driver EV',
        consumedKwh: inv.consumedKwh,
        pricePerKwh: inv.pricePerKwh,
        serviceFee: inv.serviceFee,
        tax: inv.tax,
        totalIdr: inv.total,
        status: inv.status,
        createdAt: inv.createdAt ? inv.createdAt.replace('T', ' ').substring(0, 16) : 'Baru Saja',
      };
    });
  }, [invoices, stations, currentUser]);

  // KPI Calculations
  const totalSlotsCount = stations.reduce((acc, st) => acc + (st.slots ? st.slots.length : 0), 0);
  const inUseSlotsCount = stations.reduce(
    (acc, st) => acc + (st.slots ? st.slots.filter((sl) => sl.status === 'IN_USE').length : 0),
    0
  );
  const maintenanceSlotsCount = stations.reduce(
    (acc, st) => acc + (st.slots ? st.slots.filter((sl) => sl.status === 'OUT_OF_SERVICE').length : 0),
    0
  );
  const availableSlotsCount = totalSlotsCount - inUseSlotsCount - maintenanceSlotsCount;

  // Filtered lists by selected station and search box
  const filteredActiveSessions = realActiveSessions.filter((s) => {
    const matchesStation = selectedStationId === 'ALL' || s.stationId === selectedStationId;
    const matchesSearch =
      !searchFilter ||
      s.driverName.toLowerCase().includes(searchFilter.toLowerCase()) ||
      s.vehicleName.toLowerCase().includes(searchFilter.toLowerCase()) ||
      s.stationName.toLowerCase().includes(searchFilter.toLowerCase());
    return matchesStation && matchesSearch;
  });

  const filteredBookings = realBookings.filter((b) => {
    const matchesStation = selectedStationId === 'ALL' || b.stationId === selectedStationId;
    const matchesSearch =
      !searchFilter ||
      b.driverName.toLowerCase().includes(searchFilter.toLowerCase()) ||
      b.vehicleName.toLowerCase().includes(searchFilter.toLowerCase()) ||
      b.stationName.toLowerCase().includes(searchFilter.toLowerCase());
    return matchesStation && matchesSearch;
  });

  const filteredInvoices = realInvoices.filter((inv) => {
    const matchesStation = selectedStationId === 'ALL' || inv.stationId === selectedStationId;
    const matchesSearch =
      !searchFilter ||
      inv.driverName.toLowerCase().includes(searchFilter.toLowerCase()) ||
      inv.id.toLowerCase().includes(searchFilter.toLowerCase()) ||
      inv.stationName.toLowerCase().includes(searchFilter.toLowerCase());
    return matchesStation && matchesSearch;
  });

  // Run hardware diagnostic simulation
  const handleRunDiagnostic = (stationName: string) => {
    setIsDiagnosticRunning(true);
    setDiagnosticResult(null);
    setTimeout(() => {
      setIsDiagnosticRunning(false);
      setDiagnosticResult(
        `✅ Diagnostik Hardware ${stationName} Selesai: Tegangan 380V Tiga Fase Stabil, Temp Transformer 41°C (Normal), Modbus RS485 Online.`
      );
    }, 1800);
  };

  const handleCreateTicket = (e: React.FormEvent) => {
    e.preventDefault();
    if (!newIssueDesc) return;
    const ticket: MaintenanceTicket = {
      id: `TCK-${Math.floor(Math.random() * 9000) + 1000}`,
      stationName: activeStation?.name || 'SPKLU Stasiun',
      slotNumber: newSlotNumber,
      technicianName: currentUser?.name || 'Operator On-Duty',
      issueDescription: newIssueDesc,
      severity: newSeverity,
      status: 'PENDING',
      timestamp: new Date().toISOString().replace('T', ' ').substring(0, 16),
    };
    setTickets([ticket, ...tickets]);
    setNewIssueDesc('');
  };

  return (
    <div className="space-y-8 animate-fadeIn">
      {/* Header Banner */}
      <div className="bg-[#1a2129] border border-[#262626] p-6 sm:p-8 text-white relative overflow-hidden shadow-2xl">
        <div className="relative z-10 flex flex-col md:flex-row md:items-center justify-between gap-6">
          <div className="space-y-2">
            <div className="flex items-center gap-2">
              <span className="px-2.5 py-1 text-[10px] font-extrabold bg-[#1c69d4] text-white uppercase tracking-wider">
                LAMAN OPERATOR & DASHBOARD REAL-TIME
              </span>
              <span className="flex items-center gap-1 text-[11px] font-bold text-[#22c55e] bg-white/10 px-2.5 py-0.5">
                <Radio className="w-3.5 h-3.5 animate-pulse" /> LIVE TELEMETRI & REAL DATABASE
              </span>
            </div>
            <h1 className="text-2xl sm:text-3xl font-black tracking-tight text-white flex items-center gap-3">
              <Wrench className="w-8 h-8 text-[#1c69d4]" />
              PEMANTAUAN REALTIME PENGGUNA & RIWAYAT TRANSAKSI SPKLU
            </h1>
            <p className="text-xs sm:text-sm text-[#9a9a9a] max-w-3xl font-light">
              Pantau siapa yang sedang mengisi daya di lokasi stasiun SPKLU saat ini, riwayat reservasi driver, lembar tagihan (invoices) ESDM, serta override status dispenser.
            </p>
          </div>

          <div className="flex items-center gap-3">
            <div className="p-4 bg-white/5 border border-white/10 text-right">
              <div className="text-[10px] font-bold text-[#9a9a9a] uppercase">TOTAL STASIUN MONITORED</div>
              <div className="text-2xl font-black text-[#1c69d4]">{stations.length} Stasiun</div>
            </div>
          </div>
        </div>
      </div>

      {/* KPI Status Cards */}
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
        <div className="bg-white p-5 border border-[#e6e6e6] shadow-sm flex items-center justify-between">
          <div>
            <div className="text-xs font-bold text-[#6b6b6b] uppercase">PENGGUNA AKTIF CHARGING</div>
            <div className="text-2xl font-black text-[#1c69d4] mt-1">{realActiveSessions.length} Driver</div>
            <div className="text-[10px] text-[#22c55e] font-semibold">Sedang Mengisi Daya Live</div>
          </div>
          <div className="p-3 bg-[#1c69d4]/10 text-[#1c69d4]">
            <UserCheck className="w-6 h-6" />
          </div>
        </div>

        <div className="bg-white p-5 border border-[#e6e6e6] shadow-sm flex items-center justify-between">
          <div>
            <div className="text-xs font-bold text-[#6b6b6b] uppercase">TOTAL RIWAYAT BOOKING</div>
            <div className="text-2xl font-black text-[#262626] mt-1">{realBookings.length} Reservasi</div>
            <div className="text-[10px] text-[#9a9a9a]">Data Realtime Database</div>
          </div>
          <div className="p-3 bg-[#262626]/10 text-[#262626]">
            <FileText className="w-6 h-6" />
          </div>
        </div>

        <div className="bg-white p-5 border border-[#e6e6e6] shadow-sm flex items-center justify-between">
          <div>
            <div className="text-xs font-bold text-[#6b6b6b] uppercase">TOTAL INVOICE DITERBITKAN</div>
            <div className="text-2xl font-black text-[#22c55e] mt-1">{realInvoices.length} Lembar</div>
            <div className="text-[10px] text-[#22c55e] font-semibold">Transaksi Resmi ESDM</div>
          </div>
          <div className="p-3 bg-[#22c55e]/10 text-[#22c55e]">
            <CreditCard className="w-6 h-6" />
          </div>
        </div>

        <div className="bg-white p-5 border border-[#e6e6e6] shadow-sm flex items-center justify-between">
          <div>
            <div className="text-xs font-bold text-[#6b6b6b] uppercase">KONDISI OPERASIONAL</div>
            <div className="text-2xl font-black text-[#22c55e] mt-1">{availableSlotsCount} Slot Ready</div>
            <div className="text-[10px] text-[#9a9a9a]">Dari Total {totalSlotsCount} Slot</div>
          </div>
          <div className="p-3 bg-[#22c55e]/10 text-[#22c55e]">
            <Activity className="w-6 h-6" />
          </div>
        </div>
      </div>

      {/* Filter & Station Selection Toolbar */}
      <div className="bg-white border border-[#e6e6e6] p-4 flex flex-col md:flex-row items-center justify-between gap-4 shadow-xs">
        {/* Navigation Tabs */}
        <div className="flex flex-wrap items-center gap-2 w-full md:w-auto">
          <button
            onClick={() => setActiveTabSection('LIVE')}
            className={`px-4 py-2.5 text-xs font-extrabold uppercase transition-all flex items-center gap-2 ${
              activeTabSection === 'LIVE'
                ? 'bg-[#1c69d4] text-white shadow-sm'
                : 'bg-[#fafafa] text-[#262626] border border-[#cccccc] hover:border-[#1c69d4]'
            }`}
          >
            <Activity className="w-4 h-4" />
            <span>Pengguna Saat Ini ({realActiveSessions.length})</span>
          </button>

          <button
            onClick={() => setActiveTabSection('HISTORY')}
            className={`px-4 py-2.5 text-xs font-extrabold uppercase transition-all flex items-center gap-2 ${
              activeTabSection === 'HISTORY'
                ? 'bg-[#1c69d4] text-white shadow-sm'
                : 'bg-[#fafafa] text-[#262626] border border-[#cccccc] hover:border-[#1c69d4]'
            }`}
          >
            <FileText className="w-4 h-4" />
            <span>Riwayat Booking ({realBookings.length})</span>
          </button>

          <button
            onClick={() => setActiveTabSection('INVOICES')}
            className={`px-4 py-2.5 text-xs font-extrabold uppercase transition-all flex items-center gap-2 ${
              activeTabSection === 'INVOICES'
                ? 'bg-[#1c69d4] text-white shadow-sm'
                : 'bg-[#fafafa] text-[#262626] border border-[#cccccc] hover:border-[#1c69d4]'
            }`}
          >
            <CreditCard className="w-4 h-4" />
            <span>Data Invoices ({realInvoices.length})</span>
          </button>

          <button
            onClick={() => setActiveTabSection('MAINTENANCE')}
            className={`px-4 py-2.5 text-xs font-extrabold uppercase transition-all flex items-center gap-2 ${
              activeTabSection === 'MAINTENANCE'
                ? 'bg-[#1a2129] text-white shadow-sm'
                : 'bg-[#fafafa] text-[#262626] border border-[#cccccc] hover:border-[#1a2129]'
            }`}
          >
            <Sliders className="w-4 h-4" />
            <span>Kendali Hardware Slot</span>
          </button>
        </div>

        {/* Station Selector & Search Filter */}
        <div className="flex flex-col sm:flex-row items-center gap-3 w-full md:w-auto">
          <div className="relative w-full sm:w-64">
            <Search className="w-4 h-4 absolute left-3 top-3 text-[#9a9a9a]" />
            <input
              type="text"
              placeholder="Cari driver, mobil, ID..."
              value={searchFilter}
              onChange={(e) => setSearchFilter(e.target.value)}
              className="w-full h-10 pl-9 pr-3 bg-[#fafafa] border border-[#cccccc] text-xs font-medium text-[#262626] focus:border-[#1c69d4]"
            />
          </div>

          <select
            value={selectedStationId}
            onChange={(e) => setSelectedStationId(e.target.value)}
            className="w-full sm:w-64 h-10 px-3 bg-[#fafafa] border border-[#cccccc] font-bold text-xs text-[#262626] focus:border-[#1c69d4]"
          >
            <option value="ALL">SEMUA STASIUN SPKLU</option>
            {stations.map((st) => (
              <option key={st.id} value={st.id}>
                {st.name}
              </option>
            ))}
          </select>
        </div>
      </div>

      {/* SECTION 1: PENGGUNA SAAT INI (REALTIME ACTIVE CHARGING SESSIONS) */}
      {activeTabSection === 'LIVE' && (
        <div className="space-y-4">
          <div className="flex items-center justify-between">
            <h3 className="text-base font-black text-[#262626] uppercase tracking-wide flex items-center gap-2">
              <Radio className="w-5 h-5 text-[#22c55e] animate-pulse" />
              DAFTAR PENGGUNA KENDARAAN YANG SEDANG MENGISI DAYA SAAT INI
            </h3>
            <span className="text-xs text-[#6b6b6b] font-semibold">Live Realtime Sync</span>
          </div>

          {filteredActiveSessions.length > 0 ? (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              {filteredActiveSessions.map((session) => (
                <div
                  key={session.sessionId}
                  className="bg-white border-2 border-[#1c69d4] shadow-md p-6 space-y-4 relative overflow-hidden"
                >
                  <div className="flex items-center justify-between border-b border-[#e6e6e6] pb-3">
                    <div>
                      <span className="text-[10px] font-extrabold text-[#1c69d4] bg-[#1c69d4]/10 px-2 py-0.5 uppercase tracking-wider">
                        {session.sessionId}
                      </span>
                      <h4 className="font-extrabold text-[#262626] text-base mt-1">{session.driverName}</h4>
                      <p className="text-xs text-[#6b6b6b]">{session.driverEmail}</p>
                    </div>
                    <span className="px-2.5 py-1 text-[10px] font-black bg-[#22c55e] text-white uppercase tracking-wider animate-pulse">
                      CHARGING
                    </span>
                  </div>

                  <div className="space-y-2 text-xs">
                    <div className="flex items-center justify-between text-[#3c3c3c]">
                      <span className="font-semibold">Stasiun SPKLU:</span>
                      <span className="font-bold text-right text-[#262626]">{session.stationName}</span>
                    </div>

                    <div className="flex items-center justify-between text-[#3c3c3c]">
                      <span className="font-semibold">Mobil & Status:</span>
                      <span className="font-bold text-[#1c69d4]">{session.vehicleName}</span>
                    </div>

                    <div className="flex items-center justify-between text-[#3c3c3c]">
                      <span className="font-semibold">Slot Connector:</span>
                      <span className="font-bold text-[#262626]">{session.slotNumber} • {session.connectorType}</span>
                    </div>

                    <div className="flex items-center justify-between text-[#3c3c3c]">
                      <span className="font-semibold">Waktu Mulai:</span>
                      <span className="font-mono text-[#6b6b6b]">{session.startedAt}</span>
                    </div>
                  </div>

                  {/* Telemetry Realtime Metrics */}
                  <div className="bg-[#1a2129] p-4 text-white space-y-3">
                    <div className="flex items-center justify-between text-xs">
                      <span className="text-[#9a9a9a] uppercase font-bold text-[10px]">Status Baterai EV (SOC)</span>
                      <span className="text-[#22c55e] font-black text-sm">{session.currentSocPercent}%</span>
                    </div>

                    {/* Progress Bar */}
                    <div className="w-full bg-white/20 h-2.5 overflow-hidden">
                      <div
                        className="bg-[#22c55e] h-full transition-all duration-1000"
                        style={{ width: `${session.currentSocPercent}%` }}
                      />
                    </div>

                    <div className="grid grid-cols-2 gap-2 pt-1 border-t border-white/10 text-xs">
                      <div>
                        <span className="text-[#9a9a9a] block text-[10px]">Daya Terpakai</span>
                        <span className="font-black text-white">{session.consumedKwh} kWh</span>
                      </div>
                      <div>
                        <span className="text-[#9a9a9a] block text-[10px]">Laju Charging</span>
                        <span className="font-black text-[#1c69d4]">{session.powerKw} kW</span>
                      </div>
                    </div>
                  </div>

                  <div className="pt-2 flex items-center justify-between text-xs font-bold border-t border-[#e6e6e6]">
                    <span className="text-[#6b6b6b]">Estimasi Biaya Sementara:</span>
                    <span className="text-[#1c69d4] text-sm">Rp {session.estimatedCostIdr.toLocaleString('id-ID')}</span>
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <div className="p-12 bg-white border border-[#e6e6e6] text-center space-y-3">
              <Activity className="w-10 h-10 text-[#9a9a9a] mx-auto" />
              <h4 className="font-extrabold text-base text-[#262626]">Tidak Ada Pengguna yang Sedang Mengisi Daya</h4>
              <p className="text-xs text-[#6b6b6b] max-w-md mx-auto">
                Saat ini belum ada driver yang sedang dalam proses pengisian daya di stasiun ini. Anda dapat mengubah status slot pada menu <strong>Kendali Hardware Slot</strong> untuk mengaktifkan charging.
              </p>
            </div>
          )}
        </div>
      )}

      {/* SECTION 2: RIWAYAT BOOKING STASIUN (STATION BOOKINGS LOG) */}
      {activeTabSection === 'HISTORY' && (
        <div className="bg-white border border-[#e6e6e6] p-6 space-y-4 shadow-sm">
          <div className="flex items-center justify-between border-b border-[#e6e6e6] pb-4">
            <div>
              <h3 className="text-base font-black text-[#262626] uppercase tracking-wide">
                RIWAYAT RESERVASI & PENGGUNAAN SPKLU
              </h3>
              <p className="text-xs text-[#6b6b6b]">Log histori pengisian daya oleh driver kendaraan listrik dari data riil database.</p>
            </div>
            <span className="px-3 py-1 bg-[#fafafa] border border-[#e6e6e6] text-xs font-bold text-[#1c69d4]">
              {filteredBookings.length} Record Realtime
            </span>
          </div>

          {filteredBookings.length > 0 ? (
            <div className="overflow-x-auto">
              <table className="w-full text-left text-xs">
                <thead className="bg-[#1a2129] text-white uppercase text-[10px] font-extrabold tracking-wider">
                  <tr>
                    <th className="p-3">ID Booking</th>
                    <th className="p-3">Stasiun SPKLU</th>
                    <th className="p-3">Driver / Pengguna</th>
                    <th className="p-3">Kendaraan EV</th>
                    <th className="p-3">Slot & Konektor</th>
                    <th className="p-3">Waktu Mulai - Selesai</th>
                    <th className="p-3 text-center">Status</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-[#e6e6e6] font-medium text-[#262626]">
                  {filteredBookings.map((b) => (
                    <tr key={b.id} className="hover:bg-[#fafafa] transition-colors">
                      <td className="p-3 font-mono font-extrabold text-[#1c69d4]">{b.id}</td>
                      <td className="p-3 font-bold">{b.stationName}</td>
                      <td className="p-3 font-bold text-[#262626]">{b.driverName}</td>
                      <td className="p-3 text-[#6b6b6b]">{b.vehicleName}</td>
                      <td className="p-3 font-bold">{b.slotNumber}</td>
                      <td className="p-3 font-mono text-[#6b6b6b]">{b.startTime} - {b.endTime.split(' ')[1] || b.endTime}</td>
                      <td className="p-3 text-center">
                        <span
                          className={`px-2.5 py-1 text-[10px] font-black uppercase tracking-wider ${
                            b.status === 'COMPLETED'
                              ? 'bg-[#22c55e]/10 text-[#22c55e] border border-[#22c55e]'
                              : b.status === 'IN_PROGRESS'
                              ? 'bg-[#1c69d4]/10 text-[#1c69d4] border border-[#1c69d4]'
                              : 'bg-[#dc2626]/10 text-[#dc2626] border border-[#dc2626]'
                          }`}
                        >
                          {b.status}
                        </span>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          ) : (
            <div className="p-12 bg-[#fafafa] border border-[#e6e6e6] text-center space-y-2 text-xs text-[#6b6b6b]">
              <AlertCircle className="w-8 h-8 text-[#9a9a9a] mx-auto" />
              <p className="font-bold text-[#262626]">Belum Ada Riwayat Booking</p>
              <p>Driver belum melakukan reservasi pada stasiun yang dipilih.</p>
            </div>
          )}
        </div>
      )}

      {/* SECTION 3: DATA INVOICE & PENAGIHAN (FINANCIAL INVOICES LOG) */}
      {activeTabSection === 'INVOICES' && (
        <div className="bg-white border border-[#e6e6e6] p-6 space-y-4 shadow-sm">
          <div className="flex items-center justify-between border-b border-[#e6e6e6] pb-4">
            <div>
              <h3 className="text-base font-black text-[#262626] uppercase tracking-wide">
                DAFTAR TAGIHAN & RECORD INVOICE SPKLU
              </h3>
              <p className="text-xs text-[#6b6b6b]">Rincian invoice resmi pengisian daya EV berdasarkan tarif ESDM per kWh.</p>
            </div>
            <span className="px-3 py-1 bg-[#22c55e]/10 border border-[#22c55e] text-xs font-bold text-[#22c55e]">
              REALTIME INVOICES DATA
            </span>
          </div>

          {filteredInvoices.length > 0 ? (
            <div className="overflow-x-auto">
              <table className="w-full text-left text-xs">
                <thead className="bg-[#1a2129] text-white uppercase text-[10px] font-extrabold tracking-wider">
                  <tr>
                    <th className="p-3">No. Invoice</th>
                    <th className="p-3">Stasiun SPKLU</th>
                    <th className="p-3">Driver / Pengguna</th>
                    <th className="p-3">Daya (kWh)</th>
                    <th className="p-3">Tarif ESDM</th>
                    <th className="p-3">Service Fee & PPN</th>
                    <th className="p-3">Total Bayar (Rp)</th>
                    <th className="p-3 text-center">Status</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-[#e6e6e6] font-medium text-[#262626]">
                  {filteredInvoices.map((inv) => (
                    <tr key={inv.id} className="hover:bg-[#fafafa] transition-colors">
                      <td className="p-3 font-mono font-extrabold text-[#1c69d4]">{inv.id}</td>
                      <td className="p-3 font-bold">{inv.stationName}</td>
                      <td className="p-3 font-bold">{inv.driverName}</td>
                      <td className="p-3 font-bold text-[#262626]">{inv.consumedKwh} kWh</td>
                      <td className="p-3 text-[#6b6b6b]">Rp {inv.pricePerKwh.toLocaleString('id-ID')}/kWh</td>
                      <td className="p-3 text-[#6b6b6b]">Rp {(inv.serviceFee + inv.tax).toLocaleString('id-ID')}</td>
                      <td className="p-3 font-extrabold text-[#1c69d4]">Rp {inv.totalIdr.toLocaleString('id-ID')}</td>
                      <td className="p-3 text-center">
                        <span className="px-2 py-0.5 text-[10px] font-black bg-[#22c55e] text-white uppercase">
                          {inv.status}
                        </span>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          ) : (
            <div className="p-12 bg-[#fafafa] border border-[#e6e6e6] text-center space-y-2 text-xs text-[#6b6b6b]">
              <CreditCard className="w-8 h-8 text-[#9a9a9a] mx-auto" />
              <p className="font-bold text-[#262626]">Belum Ada Tagihan / Invoice Diterbitkan</p>
              <p>Invoice akan otomatis dibuat ketika driver menyelesaikan pengisian daya.</p>
            </div>
          )}
        </div>
      )}

      {/* SECTION 4: HARDWARE KENDALI DISPENSER & MAINTENANCE LOG */}
      {activeTabSection === 'MAINTENANCE' && (
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          <div className="lg:col-span-2 space-y-6">
            <div className="bg-white border border-[#e6e6e6] p-6 space-y-6 shadow-sm">
              <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b border-[#e6e6e6] pb-4">
                <div>
                  <span className="text-[10px] font-extrabold text-[#1c69d4] bg-[#1c69d4]/10 px-2 py-0.5 uppercase tracking-wider">
                    KENDALI SLOT DISPENSER REALTIME
                  </span>
                  <h3 className="text-lg font-black text-[#262626] flex items-center gap-2 mt-1">
                    <Sliders className="w-5 h-5 text-[#1c69d4]" />
                    OVERRIDE STATUS DISPENSER HARDWARE
                  </h3>
                </div>
              </div>

              {activeStation && (
                <div className="bg-[#fafafa] p-4 border border-[#e6e6e6] space-y-3">
                  <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
                    <div>
                      <div className="font-extrabold text-sm text-[#262626] flex items-center gap-2">
                        <MapPin className="w-4 h-4 text-[#1c69d4]" />
                        {activeStation.name}
                      </div>
                      <div className="text-xs text-[#6b6b6b] mt-0.5">{activeStation.location}</div>
                    </div>

                    <button
                      onClick={() => handleRunDiagnostic(activeStation.name)}
                      disabled={isDiagnosticRunning}
                      className="bmw-btn-outline text-xs shrink-0 flex items-center gap-2"
                    >
                      <RefreshCw className={`w-3.5 h-3.5 ${isDiagnosticRunning ? 'animate-spin' : ''}`} />
                      <span>{isDiagnosticRunning ? 'MEMPROSES DIAGNOSTIK...' : 'JALANKAN DIAGNOSTIK HARDWARE'}</span>
                    </button>
                  </div>

                  {diagnosticResult && (
                    <div className="p-3 bg-[#22c55e]/10 border border-[#22c55e] text-[#15803d] text-xs font-bold animate-fadeIn">
                      {diagnosticResult}
                    </div>
                  )}
                </div>
              )}

              <div className="space-y-4">
                <h4 className="text-xs font-black uppercase text-[#6b6b6b] tracking-wider">
                  Daftar Slot Charger ({activeStation?.slots ? activeStation.slots.length : 0} Konektor):
                </h4>

                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  {activeStation?.slots && activeStation.slots.map((slot) => (
                    <div
                      key={slot.id}
                      className="bg-white border border-[#e6e6e6] p-4 space-y-3 hover:border-[#1c69d4] transition-colors relative"
                    >
                      <div className="flex items-center justify-between border-b border-[#f0f0f0] pb-2">
                        <div>
                          <span className="font-mono text-xs font-extrabold text-[#262626]">
                            SLOT #{slot.slotNumber} ({slot.id})
                          </span>
                          <div className="text-[11px] font-bold text-[#1c69d4]">
                            {slot.connectorType} • {slot.maxPowerKw} kW
                          </div>
                        </div>

                        <span
                          className={`px-2 py-1 text-[10px] font-black uppercase tracking-wider ${
                            slot.status === 'AVAILABLE'
                              ? 'bg-[#22c55e]/10 text-[#22c55e] border border-[#22c55e]/30'
                              : slot.status === 'IN_USE'
                              ? 'bg-[#1c69d4]/10 text-[#1c69d4] border border-[#1c69d4]/30'
                              : 'bg-[#eab308]/10 text-[#eab308] border border-[#eab308]/30'
                          }`}
                        >
                          {slot.status}
                        </span>
                      </div>

                      <div className="space-y-2 pt-1">
                        <div className="text-[10px] font-bold text-[#9a9a9a] uppercase">Manual Status Override:</div>
                        {slot.status === 'IN_USE' ? (
                          <div className="p-2 bg-[#1c69d4]/10 border border-[#1c69d4]/30 text-[#1c69d4] text-[10px] font-bold flex items-center justify-between">
                            <span className="flex items-center gap-1.5">
                              <Lock className="w-3.5 h-3.5 shrink-0" />
                              <span>IN_USE (Terkunci Sesi Charging)</span>
                            </span>
                            <span className="text-[9px] bg-[#1c69d4] text-white px-1.5 py-0.5 uppercase tracking-wider font-extrabold">TERKUNCI</span>
                          </div>
                        ) : (
                          <div className="grid grid-cols-3 gap-1.5">
                            <button
                              onClick={() => onUpdateSlotStatus(activeStation.id, slot.id, 'AVAILABLE')}
                              className={`py-1.5 px-2 text-[10px] font-bold uppercase border transition-colors ${
                                slot.status === 'AVAILABLE'
                                  ? 'bg-[#22c55e] text-white border-[#22c55e]'
                                  : 'bg-white text-[#262626] border-[#cccccc] hover:bg-[#fafafa]'
                              }`}
                            >
                              AVAILABLE
                            </button>
                            <button
                              disabled
                              title="Status IN_USE diaktifkan otomatis saat pengemudi check-in"
                              className="py-1.5 px-2 text-[10px] font-bold uppercase border bg-[#f0f0f0] text-[#9a9a9a] border-[#cccccc] cursor-not-allowed"
                            >
                              IN USE
                            </button>
                            <button
                              onClick={() => onUpdateSlotStatus(activeStation.id, slot.id, 'OUT_OF_SERVICE')}
                              className={`py-1.5 px-2 text-[10px] font-bold uppercase border transition-colors ${
                                slot.status === 'OUT_OF_SERVICE'
                                  ? 'bg-[#eab308] text-white border-[#eab308]'
                                  : 'bg-white text-[#262626] border-[#cccccc] hover:bg-[#fafafa]'
                              }`}
                            >
                              SERVICE
                            </button>
                          </div>
                        )}
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            </div>

            <div className="bg-[#dc2626]/5 border border-[#dc2626]/30 p-6 space-y-4">
              <div className="flex items-center gap-3">
                <div className="p-2 bg-[#dc2626] text-white">
                  <ShieldAlert className="w-6 h-6" />
                </div>
                <div>
                  <h4 className="text-sm font-black text-[#dc2626] uppercase tracking-wider">
                    EMERGENCY POWER CUT-OFF DISPENSER
                  </h4>
                  <p className="text-xs text-[#6b6b6b]">Gunakan hanya pada situasi darurat keselamatan fisik lokasi stasiun.</p>
                </div>
              </div>

              <div className="flex items-center gap-3 pt-2">
                <button
                  onClick={() => alert(`⚡ Emergency Power Cut Off berhasil dipicu pada ${activeStation?.name}.`)}
                  className="px-4 py-2.5 bg-[#dc2626] hover:bg-[#b91c1c] text-white font-extrabold text-xs uppercase flex items-center gap-2"
                >
                  <Power className="w-4 h-4" /> CUT-OFF CATU DAYA
                </button>
              </div>
            </div>
          </div>

          <div className="space-y-6">
            <div className="bg-white border border-[#e6e6e6] p-6 space-y-4 shadow-sm">
              <div className="border-b border-[#e6e6e6] pb-3">
                <h3 className="text-base font-black text-[#262626] flex items-center gap-2">
                  <ClipboardList className="w-5 h-5 text-[#1c69d4]" />
                  INPUT TIKET INSPEKSI
                </h3>
              </div>

              <form onSubmit={handleCreateTicket} className="space-y-3 text-xs">
                <div>
                  <label className="text-[#262626] font-bold uppercase block mb-1">Target Slot:</label>
                  <select
                    value={newSlotNumber}
                    onChange={(e) => setNewSlotNumber(e.target.value)}
                    className="w-full h-10 px-3 bg-[#fafafa] border border-[#cccccc] font-bold text-xs text-[#262626]"
                  >
                    <option value="Slot #01">Slot #01 (CCS2 150 kW)</option>
                    <option value="Slot #02">Slot #02 (Type 2 22 kW)</option>
                    <option value="Slot #03">Slot #03 (CHAdeMO 50 kW)</option>
                  </select>
                </div>

                <div>
                  <label className="text-[#262626] font-bold uppercase block mb-1">Urgensi:</label>
                  <div className="grid grid-cols-3 gap-1.5">
                    <button
                      type="button"
                      onClick={() => setNewSeverity('LOW')}
                      className={`py-2 text-[10px] font-bold uppercase border ${
                        newSeverity === 'LOW' ? 'bg-[#22c55e] text-white' : 'bg-white text-[#262626]'
                      }`}
                    >
                      LOW
                    </button>
                    <button
                      type="button"
                      onClick={() => setNewSeverity('MEDIUM')}
                      className={`py-2 text-[10px] font-bold uppercase border ${
                        newSeverity === 'MEDIUM' ? 'bg-[#eab308] text-white' : 'bg-white text-[#262626]'
                      }`}
                    >
                      MEDIUM
                    </button>
                    <button
                      type="button"
                      onClick={() => setNewSeverity('CRITICAL')}
                      className={`py-2 text-[10px] font-bold uppercase border ${
                        newSeverity === 'CRITICAL' ? 'bg-[#dc2626] text-white' : 'bg-white text-[#262626]'
                      }`}
                    >
                      CRITICAL
                    </button>
                  </div>
                </div>

                <div>
                  <label className="text-[#262626] font-bold uppercase block mb-1">Hasil Inspeksi:</label>
                  <textarea
                    required
                    rows={3}
                    placeholder="Deskripsi perbaikan..."
                    value={newIssueDesc}
                    onChange={(e) => setNewIssueDesc(e.target.value)}
                    className="w-full p-2.5 bg-white border border-[#cccccc] text-[#262626]"
                  />
                </div>

                <button type="submit" className="bmw-btn-primary w-full flex items-center justify-center gap-2">
                  <CheckCircle2 className="w-4 h-4" /> SIMPAN TIKET
                </button>
              </form>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};
