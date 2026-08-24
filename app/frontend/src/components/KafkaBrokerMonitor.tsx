import React, { useState, useEffect } from 'react';
import { Server, Database, Zap, RefreshCw, Send, CheckCircle2, FileText, CreditCard, Clock, Play } from 'lucide-react';
import { apiService } from '../services/api';

interface KafkaEvent {
  id: string;
  topic: string;
  source: string;
  aggregate_id: string;
  payload: any;
  timestamp: string;
}

export const KafkaBrokerMonitor: React.FC = () => {
  const [events, setEvents] = useState<KafkaEvent[]>([]);
  const [loading, setLoading] = useState(false);
  const [publishing, setPublishing] = useState(false);
  const [selectedTopic, setSelectedTopic] = useState('ChargingCompleted');
  const [customAggregateId, setCustomAggregateId] = useState('ses-sim-001');

  const topics = [
    { name: 'BookingConfirmed', color: 'bg-emerald-500', desc: 'Station DB & Booking DB -> Event Booking Terkonfirmasi' },
    { name: 'BookingExpired', color: 'bg-amber-500', desc: 'Booking DB -> Event Reservasi Kadaluarsa' },
    { name: 'BookingCancelled', color: 'bg-rose-500', desc: 'Booking DB -> Event Pembatalan Reservasi' },
    { name: 'ChargingStarted', color: 'bg-blue-500', desc: 'Session DB -> Event Pengisian Daya Dimulai' },
    { name: 'ChargingCompleted', color: 'bg-indigo-500', desc: 'Session DB -> Event Pengisian Daya Selesai' },
    { name: 'PaymentCreated', color: 'bg-violet-500', desc: 'Billing Service -> Event Tagihan / Invoice Diterbitkan' },
    { name: 'PaymentCompleted', color: 'bg-emerald-600', desc: 'Billing Service -> Event Pembayaran Lunas' },
    { name: 'PaymentFailed', color: 'bg-red-600', desc: 'Billing Service -> Event Pembayaran Gagal' },
  ];

  const fetchKafkaEvents = async () => {
    setLoading(true);
    try {
      const data = await apiService.getKafkaEvents();
      if (data && data.events) {
        setEvents(data.events.reverse());
      }
    } catch (err) {
      console.warn('Failed to fetch Kafka events:', err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchKafkaEvents();
    const interval = setInterval(fetchKafkaEvents, 5000);
    return () => clearInterval(interval);
  }, []);

  const handleSimulatePublish = async () => {
    setPublishing(true);
    try {
      await apiService.publishKafkaEvent({
        topic: selectedTopic,
        source: selectedTopic.startsWith('Booking') ? 'Booking DB' : selectedTopic.startsWith('Charging') ? 'Session DB' : 'Billing DB',
        aggregate_id: customAggregateId || `agg-${Date.now()}`,
        payload: {
          simulated_at: new Date().toISOString(),
          user_id: 'usr-driver',
          consumed_kwh: 24.8,
          price_per_kwh: 2467,
          message: `Simulated event ${selectedTopic} via Kafka Message Broker visualizer`,
        },
      });
      await fetchKafkaEvents();
    } catch (err) {
      console.error('Failed to publish Kafka event:', err);
    } finally {
      setPublishing(false);
    }
  };

  return (
    <div className="space-y-6">
      {/* Header Banner */}
      <div className="bg-[#1a2129] text-white p-6 border border-[#262e38] shadow-md">
        <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
          <div>
            <div className="flex items-center gap-2 mb-1">
              <span className="w-2.5 h-2.5 rounded-full bg-emerald-400 animate-ping" />
              <span className="text-xs font-mono tracking-widest text-emerald-400 uppercase font-bold">KAFKA MESSAGE BROKER ONLINE</span>
            </div>
            <h2 className="text-xl font-extrabold uppercase tracking-tight text-white flex items-center gap-2">
              <Server className="w-6 h-6 text-[#1c69d4]" />
              Event Streaming Architecture (Port 9092)
            </h2>
            <p className="text-xs text-[#bbbbbb] mt-1">
              Sinkronisasi data terdistribusi otomatis antara Station DB, Booking DB, Session DB, Kafka Message Broker, dan Billing Service.
            </p>
          </div>
          <button
            onClick={fetchKafkaEvents}
            disabled={loading}
            className="px-4 py-2 bg-[#1c69d4] hover:bg-[#0653b6] text-white text-xs font-bold uppercase tracking-wider flex items-center gap-2 transition-colors self-start md:self-auto shrink-0"
          >
            <RefreshCw className={`w-4 h-4 ${loading ? 'animate-spin' : ''}`} />
            REFRESH EVENT STREAM
          </button>
        </div>
      </div>

      {/* Architecture Flow Diagram Visualizer */}
      <div className="bg-white border border-[#e6e6e6] p-6 shadow-xs">
        <h3 className="text-xs font-extrabold text-[#262626] uppercase tracking-wider mb-4 flex items-center gap-2">
          <Zap className="w-4 h-4 text-[#1c69d4]" />
          Arsitektur Event-Driven (Distributed Message Flow)
        </h3>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
          {/* Source Databases */}
          <div className="bg-[#fafafa] border border-[#cccccc] p-4 flex flex-col items-center justify-center text-center">
            <div className="flex items-center gap-3 mb-2">
              <Database className="w-5 h-5 text-emerald-600" />
              <span className="text-xs font-bold uppercase text-[#262626]">Station DB</span>
            </div>
            <div className="flex items-center gap-3 mb-2">
              <Database className="w-5 h-5 text-blue-600" />
              <span className="text-xs font-bold uppercase text-[#262626]">Booking DB</span>
            </div>
            <div className="flex items-center gap-3">
              <Database className="w-5 h-5 text-orange-600" />
              <span className="text-xs font-bold uppercase text-[#262626]">Session DB</span>
            </div>
          </div>

          {/* Central Kafka Broker */}
          <div className="bg-[#1c69d4]/10 border-2 border-[#1c69d4] p-4 flex flex-col items-center justify-center text-center">
            <Server className="w-8 h-8 text-[#1c69d4] mb-2" />
            <h4 className="text-xs font-extrabold uppercase text-[#1c69d4] tracking-wider">Kafka Message Broker</h4>
            <span className="text-[10px] text-[#6b6b6b] font-mono mt-1">Pub/Sub Event Bus (Port 9092)</span>
            <div className="mt-2 text-[10px] font-bold text-emerald-700 bg-emerald-100 px-2 py-0.5 rounded-full">
              {events.length} Event Streamed
            </div>
          </div>

          {/* Billing Service & DB */}
          <div className="bg-[#fafafa] border border-[#cccccc] p-4 flex flex-col items-center justify-center text-center">
            <Server className="w-6 h-6 text-indigo-600 mb-1" />
            <h4 className="text-xs font-bold uppercase text-[#262626]">Billing Service</h4>
            <div className="grid grid-cols-2 gap-1 text-[10px] text-[#6b6b6b] mt-2 text-left w-full bg-white p-2 border border-[#e6e6e6]">
              <div className="flex items-center gap-1"><CheckCircle2 className="w-3 h-3 text-emerald-500" /> Calculate Charges</div>
              <div className="flex items-center gap-1"><FileText className="w-3 h-3 text-blue-500" /> Generate Invoice</div>
              <div className="flex items-center gap-1"><Clock className="w-3 h-3 text-amber-500" /> Process Payment</div>
              <div className="flex items-center gap-1"><CreditCard className="w-3 h-3 text-indigo-500" /> Payment Status</div>
            </div>
            <div className="flex items-center gap-2 mt-2 pt-2 border-t border-[#cccccc] w-full justify-center">
              <Database className="w-4 h-4 text-purple-600" />
              <span className="text-[11px] font-bold text-[#262626]">Billing DB</span>
            </div>
          </div>
        </div>

        {/* Kafka Event Simulator Control */}
        <div className="p-4 bg-[#f7f7f7] border border-[#e6e6e6]">
          <h4 className="text-xs font-extrabold uppercase tracking-wider text-[#262626] mb-3 flex items-center gap-2">
            <Play className="w-4 h-4 text-[#1c69d4]" />
            Simulasi Kirim Kafka Event
          </h4>
          <div className="flex flex-col sm:flex-row items-center gap-3">
            <select
              value={selectedTopic}
              onChange={(e) => setSelectedTopic(e.target.value)}
              className="w-full sm:w-1/3 h-10 px-3 bg-white border border-[#cccccc] text-xs font-bold text-[#262626]"
            >
              {topics.map((t) => (
                <option key={t.name} value={t.name}>{t.name}</option>
              ))}
            </select>
            <input
              type="text"
              placeholder="Aggregate ID (e.g. bkg-101 / ses-202)"
              value={customAggregateId}
              onChange={(e) => setCustomAggregateId(e.target.value)}
              className="w-full sm:w-1/3 h-10 px-3 bg-white border border-[#cccccc] text-xs text-[#262626]"
            />
            <button
              onClick={handleSimulatePublish}
              disabled={publishing}
              className="w-full sm:w-1/3 h-10 bg-[#1c69d4] hover:bg-[#0653b6] text-white font-bold text-xs uppercase tracking-wider flex items-center justify-center gap-2 transition-colors"
            >
              <Send className="w-3.5 h-3.5" />
              {publishing ? 'PUBLISHING...' : 'PUBLISH EVENT'}
            </button>
          </div>
        </div>
      </div>

      {/* Topics Summary Grid */}
      <div className="grid grid-cols-2 sm:grid-cols-4 gap-3">
        {topics.map((t) => {
          const count = events.filter((e) => e.topic === t.name).length;
          return (
            <div key={t.name} className="bg-white border border-[#e6e6e6] p-3.5 shadow-xs">
              <div className="flex items-center justify-between mb-1.5">
                <span className={`w-2.5 h-2.5 rounded-full ${t.color}`} />
                <span className="text-xs font-bold text-[#1c69d4] font-mono">{count} Events</span>
              </div>
              <h4 className="text-xs font-extrabold text-[#262626] font-mono">{t.name}</h4>
              <p className="text-[10px] text-[#6b6b6b] mt-1 leading-tight">{t.desc}</p>
            </div>
          );
        })}
      </div>

      {/* Live Kafka Event Stream Feed */}
      <div className="bg-white border border-[#e6e6e6] p-6 shadow-xs">
        <h3 className="text-xs font-extrabold text-[#262626] uppercase tracking-wider mb-4 flex items-center gap-2">
          <Server className="w-4 h-4 text-[#1c69d4]" />
          Live Kafka Message Feed ({events.length} Terrekam)
        </h3>

        {events.length === 0 ? (
          <div className="p-8 text-center bg-[#fafafa] border border-dashed border-[#cccccc] text-xs text-[#6b6b6b]">
            Belum ada event Kafka yang dipublikasikan. Gunakan tombol simulator di atas untuk mengirim event.
          </div>
        ) : (
          <div className="space-y-2.5 max-h-[420px] overflow-y-auto pr-1">
            {events.map((evt) => (
              <div key={evt.id} className="p-3 bg-[#fafafa] border border-[#e6e6e6] hover:border-[#1c69d4] transition-colors">
                <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-1 mb-1.5">
                  <div className="flex items-center gap-2">
                    <span className="text-[10px] font-mono bg-[#1c69d4] text-white px-2 py-0.5 font-bold uppercase">
                      TOPIC: {evt.topic}
                    </span>
                    <span className="text-xs font-extrabold text-[#262626] font-mono">ID: {evt.id}</span>
                  </div>
                  <span className="text-[11px] text-[#6b6b6b] font-mono">
                    {new Date(evt.timestamp).toLocaleString('id-ID')}
                  </span>
                </div>
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-2 text-xs text-[#3c3c3c] bg-white p-2 border border-[#e6e6e6] font-mono">
                  <div><span className="font-bold text-[#1a1a1a]">Source:</span> {evt.source}</div>
                  <div><span className="font-bold text-[#1a1a1a]">Aggregate ID:</span> {evt.aggregate_id}</div>
                  <div className="col-span-1 sm:col-span-2 truncate"><span className="font-bold text-[#1a1a1a]">Payload:</span> {JSON.stringify(evt.payload)}</div>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
};
