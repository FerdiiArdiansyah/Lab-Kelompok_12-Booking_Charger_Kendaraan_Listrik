import { useState, useEffect } from 'react';
import { Navbar } from './components/Navbar';
import { StationExplorer } from './components/StationExplorer';
import { BookingModal } from './components/BookingModal';
import { LiveSessionViewer } from './components/LiveSessionViewer';
import { BillingManager } from './components/BillingManager';
import { UserGarage } from './components/UserGarage';
import { UserBookingsList } from './components/UserBookingsList';
import { AdminPortal } from './components/AdminPortal';
import { OperatorPortal } from './components/OperatorPortal';
import { AuthModal } from './components/AuthModal';

import { apiService } from './services/api';
import type { Station, ChargerSlot, UserVehicle, Booking, ChargingSession, Invoice, User } from './types';

export function App() {
  const [currentUser, setCurrentUser] = useState<User | null>(null);
  const [activeTab, setActiveTab] = useState<string>('finder');
  const [isAuthModalOpen, setIsAuthModalOpen] = useState<boolean>(false);

  // Public State (Available to all visitors)
  const [stations, setStations] = useState<Station[]>([]);

  // User-Scoped Private States (Empty until user logs in)
  const [vehicles, setVehicles] = useState<UserVehicle[]>([]);
  const [bookings, setBookings] = useState<Booking[]>([]);
  const [activeSession, setActiveSession] = useState<ChargingSession | null>(null);
  const [invoices, setInvoices] = useState<Invoice[]>([]);

  // Modal State
  const [selectedStationForBooking, setSelectedStationForBooking] = useState<{
    station: Station;
    slot?: ChargerSlot;
  } | null>(null);

  // Toast notification alert
  const [toastMessage, setToastMessage] = useState<string | null>(null);

  const showToast = (msg: string) => {
    setToastMessage(msg);
    setTimeout(() => setToastMessage(null), 3500);
  };

  // 1. Load Public Station Data & Restore Active Session on Mount
  useEffect(() => {
    async function initApp() {
      try {
        const [fetchedStations, fetchedAllBookings] = await Promise.all([
          apiService.getStations(),
          apiService.getAllBookings(),
        ]);
        if (fetchedStations && fetchedStations.length > 0) {
          setStations(fetchedStations);
        }
        if (fetchedAllBookings && fetchedAllBookings.length > 0) {
          setBookings(fetchedAllBookings);
        }
      } catch (err) {
        console.error('Error fetching public data:', err);
      }

      // Restore session from localStorage if exists
      const savedUserStr = localStorage.getItem('volt_user') || localStorage.getItem('user');
      if (savedUserStr) {
        try {
          const parsed = JSON.parse(savedUserStr);
          const userObj: User = {
            id: parsed.id || parsed.userId || 'usr-admin',
            name: parsed.name || parsed.nama || 'Administrator VoltHub',
            email: parsed.email || 'admin@spklu.co.id',
            phone: parsed.phone || '08123456789',
            role: (parsed.role ? parsed.role.toUpperCase() : 'USER') as 'USER' | 'ADMIN' | 'OPERATOR',
            status: 'ACTIVE',
          };
          setCurrentUser(userObj);
          if (userObj.role === 'ADMIN') {
            setActiveTab('admin');
          } else if (userObj.role === 'OPERATOR') {
            setActiveTab('operator');
          }
        } catch (e) {
          console.error('Error restoring session:', e);
        }
      }
    }
    initApp();
  }, []);

  // 2. Sync Public and User-Scoped Data when currentUser changes
  useEffect(() => {
    async function loadAppAndUserData() {
      try {
        const allBookingsBackend = await apiService.getAllBookings();

        if (currentUser) {
          const [userVehicles, userBookings, userSession, userInvoices] = await Promise.all([
            apiService.getVehicles(currentUser.id),
            apiService.getBookings(currentUser.id),
            apiService.getActiveSession(currentUser.id),
            apiService.getInvoices(currentUser.id),
          ]);

          setVehicles(userVehicles || []);
          
          // Merge user bookings and all backend bookings dynamically
          const combined = [...(allBookingsBackend || []), ...(userBookings || [])];
          const uniqueBookingsMap = new Map();
          combined.forEach(b => uniqueBookingsMap.set(b.id, b));
          setBookings(Array.from(uniqueBookingsMap.values()));

          const validActiveSession = userSession && userSession.id && userSession.status === 'IN_PROGRESS' ? userSession : null;
          setActiveSession(validActiveSession);
          setInvoices(userInvoices || []);
        } else {
          // Visitor mode (logged out): Always show public station booking schedule matrix!
          setBookings(allBookingsBackend || []);
          setVehicles([]);
          setActiveSession(null);
          setInvoices([]);
          setActiveTab('finder');
        }
      } catch (err) {
        console.error('Error fetching app data:', err);
      }
    }

    loadAppAndUserData();
  }, [currentUser]);

  // 3. Auto-Logout on 20 Seconds of User Inactivity (Auto Hapus Token & Reset Sesi)
  useEffect(() => {
    if (!currentUser) return;

    const INACTIVITY_LIMIT_MS = 20000; // 20 detik tanpa interaksi
    let timer: number;

    const performAutoLogout = async () => {
      localStorage.removeItem('volt_token');
      localStorage.removeItem('token');
      localStorage.removeItem('volt_user');
      localStorage.removeItem('user');
      setCurrentUser(null);
      setVehicles([]);
      const publicBookings = await apiService.getAllBookings();
      setBookings(publicBookings || []);
      setActiveSession(null);
      setInvoices([]);
      setActiveTab('finder');
      showToast('⚠️ Sesi & Token telah dihapus otomatis karena tidak ada interaksi pengguna selama 20 detik!');
    };

    const resetInactivityTimer = () => {
      if (timer) window.clearTimeout(timer);
      timer = window.setTimeout(performAutoLogout, INACTIVITY_LIMIT_MS);
    };

    // Jalankan timer saat komponen terdeteksi user login
    resetInactivityTimer();

    // Deteksi seluruh interaksi pengguna (mouse, keyboard, click, scroll, touch)
    const events = ['mousemove', 'keydown', 'mousedown', 'scroll', 'touchstart', 'pointerdown'];
    events.forEach((evt) => window.addEventListener(evt, resetInactivityTimer, { passive: true }));

    return () => {
      if (timer) window.clearTimeout(timer);
      events.forEach((evt) => window.removeEventListener(evt, resetInactivityTimer));
    };
  }, [currentUser]);

  // Tab Guard Navigation Helper
  const handleSelectTab = (tab: string) => {
    if (!currentUser && tab !== 'finder') {
      setIsAuthModalOpen(true);
      showToast('Silakan Login atau Registrasi terlebih dahulu untuk melihat data Anda.');
      return;
    }
    setActiveTab(tab);
  };

  const handleLogout = async () => {
    setCurrentUser(null);
    localStorage.removeItem('volt_token');
    localStorage.removeItem('token');
    localStorage.removeItem('volt_user');
    localStorage.removeItem('user');
    setVehicles([]);
    setActiveSession(null);
    setInvoices([]);
    setActiveTab('finder');
    const publicBookings = await apiService.getAllBookings();
    setBookings(publicBookings || []);
    showToast('Anda telah keluar (Logout). Session & Hak Akses Privat dibersihkan.');
  };

  // Auth Handlers
  const handleRegister = async (data: { name: string; email: string; password: string; phone: string; role: string }) => {
    const user = await apiService.registerUser(data);
    const defaultTab = user.role === 'ADMIN' ? 'admin' : user.role === 'OPERATOR' ? 'operator' : 'finder';
    setActiveTab(defaultTab);
    showToast(`Registrasi berhasil! Selamat datang, ${user.name}`);
    return user;
  };

  const handleLogin = async (email: string, pass: string) => {
    const user = await apiService.loginUser(email, pass);
    if (user) {
      const defaultTab = user.role === 'ADMIN' ? 'admin' : user.role === 'OPERATOR' ? 'operator' : 'finder';
      setActiveTab(defaultTab);
      showToast(`Login Berhasil! Selamat datang ${user.name} (${user.role})`);
    }
    return user;
  };

  // 1. Create Booking handler
  const handleConfirmBooking = async (data: {
    slotId: string;
    vehicleId: string;
    startTime: string;
    endTime: string;
  }) => {
    if (!selectedStationForBooking) return;
    if (!currentUser) {
      setIsAuthModalOpen(true);
      showToast('Silakan Login atau Registrasi terlebih dahulu untuk reservasi.');
      return;
    }

    const bookingPayload = {
      stationId: selectedStationForBooking.station.id,
      slotId: data.slotId,
      vehicleId: data.vehicleId,
      startTime: data.startTime,
      endTime: data.endTime,
      userId: currentUser.id,
    };

    const createdBooking = await apiService.createBooking(bookingPayload);

    setBookings((prev) => [createdBooking, ...prev]);
    setSelectedStationForBooking(null);
    setActiveTab('bookings');
    showToast(`Booking #${createdBooking.id} BERHASIL dibuat! Bebas Overlap.`);
  };

  // 1b. Direct Start Session from Modal
  const handleConfirmAndStartSession = async (data: {
    slotId: string;
    vehicleId: string;
    startTime: string;
    endTime: string;
  }) => {
    if (!selectedStationForBooking) return;
    if (!currentUser) {
      setIsAuthModalOpen(true);
      showToast('Silakan Registrasi atau Login terlebih dahulu untuk cas.');
      return;
    }

    const bookingPayload = {
      stationId: selectedStationForBooking.station.id,
      slotId: data.slotId,
      vehicleId: data.vehicleId,
      startTime: data.startTime,
      endTime: data.endTime,
      userId: currentUser.id,
    };

    const createdBooking = await apiService.createBooking(bookingPayload);
    setBookings((prev) => [createdBooking, ...prev]);

    const newSession = await apiService.startSession(createdBooking.id, createdBooking.slotId, currentUser.id);

    // Dynamic Backend Lock to IN_USE (station-service)
    await apiService.updateSlotStatus(selectedStationForBooking.station.id, data.slotId, 'IN_USE');

    setStations((prev) =>
      prev.map((st) => ({
        ...st,
        slots: st.slots
          ? st.slots.map((sl) => (sl.id === data.slotId ? { ...sl, status: 'IN_USE' } : sl))
          : [],
      }))
    );

    setSelectedStationForBooking(null);
    setActiveSession(newSession);
    setActiveTab('session');
    showToast(`⚡ Sesi pengisian daya langsung dimulai! Slot #${data.slotId} terkunci IN_USE & Telemetry aktif.`);
  };

  // 2. Join Waitlist handler
  const handleJoinWaitlist = () => {
    setSelectedStationForBooking(null);
    showToast('Anda BERHASIL masuk dalam daftar antrean Waitlist!');
  };

  // 3. Start Session from Booking (Check-In)
  const handleStartSessionFromBooking = async (booking: Booking) => {
    if (!currentUser) return;
    const newSession = await apiService.startSession(booking.id, booking.slotId, currentUser.id);

    setBookings((prev) =>
      prev.map((b) => (b.id === booking.id ? { ...b, status: 'IN_PROGRESS' } : b))
    );

    // Dynamic Backend Lock to IN_USE (station-service)
    const targetStation = stations.find((st) => st.slots && st.slots.some((sl) => sl.id === booking.slotId));
    if (targetStation) {
      await apiService.updateSlotStatus(targetStation.id, booking.slotId, 'IN_USE');
    }

    // Immediately lock slot status to IN_USE upon check-in
    setStations((prev) =>
      prev.map((st) => ({
        ...st,
        slots: st.slots
          ? st.slots.map((sl) => (sl.id === booking.slotId ? { ...sl, status: 'IN_USE' } : sl))
          : [],
      }))
    );
    setActiveSession(newSession);
    setActiveTab('session');
    showToast(`Check-In berhasil! Status slot terkunci IN_USE & Telemetry charging #${newSession.id} aktif.`);
  };

  // 4. Finish Session & Generate Invoice
  const handleFinishSession = async (sessionId: string, finalKwh: number) => {
    const station = stations.find((st) => st.slots && st.slots.some((sl) => sl.id === activeSession?.slotId));
    const pricePerKwh = Math.round(station?.activeTariff?.pricePerKwh || 2467);
    const isUltraFast = (station?.totalPower || 0) >= 150;

    let serviceFee = 0;
    if (finalKwh < 0.5) {
      serviceFee = isUltraFast ? 5000 : 2500;
    } else {
      serviceFee = isUltraFast ? 20000 : 10000;
    }

    const activeSlotId = activeSession?.slotId;
    const activeBookingId = activeSession?.bookingId;
    
    // 1. Dynamic Backend Finish Session (session-service)
    await apiService.finishSession(sessionId, finalKwh);

    // 2. Dynamic Backend Complete Booking (booking-service)
    if (activeBookingId) {
      await apiService.completeBooking(activeBookingId);
    }

    // 3. Dynamic Backend Unlock Slot to AVAILABLE (station-service)
    if (station && activeSlotId) {
      await apiService.updateSlotStatus(station.id, activeSlotId, 'AVAILABLE');
    }

    // 4. Update local booking state to COMPLETED
    setBookings((prev) =>
      prev.map((b) =>
        b.id === activeBookingId || (activeSlotId && b.slotId === activeSlotId && (b.status === 'IN_PROGRESS' || b.status === 'CONFIRMED'))
          ? { ...b, status: 'COMPLETED' }
          : b
      )
    );

    const createdInvoice = await apiService.createInvoice({
      sessionId,
      userId: currentUser?.id || 'usr-driver',
      tariffId: station?.activeTariff?.id || 'trf-001',
      consumedKwh: finalKwh,
      pricePerKwh,
      serviceFee,
    });

    setInvoices((prev) => [createdInvoice, ...prev]);

    if (activeSlotId) {
      setStations((prev) =>
        prev.map((st) => ({
          ...st,
          slots: st.slots
            ? st.slots.map((sl) => (sl.id === activeSlotId ? { ...sl, status: 'AVAILABLE' } : sl))
            : [],
        }))
      );
    }

    // Refresh all bookings from backend to immediately update 24-hour matrix on dashboard
    apiService.getAllBookings().then((updatedList) => {
      if (updatedList) setBookings(updatedList);
    });

    setActiveSession(null);
    setActiveTab('billing');
    showToast(`Charging selesai! Sesi #${sessionId} tercatat & Invoice #${createdInvoice.id} terbit.`);
  };

  // 5. Confirm Payment
  const handleConfirmPayment = async (invoiceId: string, paymentMethod: any) => {
    const invObj = invoices.find((i) => i.id === invoiceId);
    const amount = invObj ? Math.round(invObj.total) : 50000;

    await apiService.payInvoice(invoiceId, paymentMethod, amount);
    setInvoices((prev) =>
      prev.map((inv) => (inv.id === invoiceId ? { ...inv, status: 'PAID' } : inv))
    );
    showToast(`Pembayaran Invoice #${invoiceId} via ${paymentMethod} LUNAS!`);
  };

  // 6. Garage Operations
  const handleAddVehicle = async (vehicleData: Omit<UserVehicle, 'id' | 'createdAt'>) => {
    if (!currentUser) return;
    try {
      const newVehicle = await apiService.addVehicle({ ...vehicleData, userId: currentUser.id });
      setVehicles((prev) => [newVehicle, ...prev]);
      showToast(`Mobil EV ${newVehicle.brand} ${newVehicle.model} (${newVehicle.licensePlate}) ditambahkan!`);
    } catch (err: any) {
      showToast(`⚠️ ${err.message || 'Gagal menambahkan kendaraan.'}`);
    }
  };

  const handleDeleteVehicle = async (id: string) => {
    if (!currentUser) return;
    await apiService.deleteVehicle(id, currentUser.id);
    setVehicles((prev) => prev.filter((v) => v.id !== id));
    showToast('Mobil EV berhasil dihapus dari garasi.');
  };

  // 7. Admin & Operator Slot Status Update
  const handleUpdateSlotStatus = async (stationId: string, slotId: string, newStatus: any) => {
    // Prevent manual override if slot is currently IN_USE by an active session
    const targetStation = stations.find((st) => st.id === stationId);
    const targetSlot = targetStation?.slots?.find((sl) => sl.id === slotId);

    if (targetSlot?.status === 'IN_USE') {
      showToast('⚠️ Status slot IN_USE sedang digunakan oleh pengisian aktif dan tidak dapat diubah hingga pengemudi menyelesaikan pengisian.');
      return;
    }

    await apiService.updateSlotStatus(stationId, slotId, newStatus);
    setStations((prev) =>
      prev.map((st) => {
        if (st.id !== stationId) return st;
        return {
          ...st,
          slots: st.slots ? st.slots.map((sl) => (sl.id === slotId ? { ...sl, status: newStatus } : sl)) : [],
        };
      })
    );
    showToast(`Status Slot #${slotId} diperbarui menjadi ${newStatus}`);
  };

  const handleAddStation = async (stationData: { name: string; location: string; mapUrl?: string; imageUrl?: string; totalPower: number }) => {
    const newStation = await apiService.addStation(stationData);
    setStations((prev) => [newStation, ...prev]);
    showToast(`Stasiun SPKLU ${newStation.name} berhasil ditambahkan!`);
  };

  const handleUpdateStation = async (stationId: string, data: { name?: string; location?: string; mapUrl?: string; imageUrl?: string; totalPower?: number }) => {
    await apiService.updateStation(stationId, data);
    setStations((prev) =>
      prev.map((st) => (st.id === stationId ? { ...st, ...data } : st))
    );
    showToast(`Stasiun SPKLU ${data.name || stationId} berhasil diperbarui!`);
  };

  return (
    <div className="min-h-screen bg-white text-[#262626] selection:bg-[#1c69d4] selection:text-white">
      {/* Top Sticky Header */}
      <Navbar
        activeTab={activeTab}
        setActiveTab={handleSelectTab}
        currentUser={currentUser}
        onOpenAuthModal={() => setIsAuthModalOpen(true)}
        onLogout={handleLogout}
        hasActiveSession={!!activeSession}
      />

      {/* Toast Alert Banner */}
      {toastMessage && (
        <div className="fixed bottom-20 md:bottom-6 left-4 right-4 md:left-auto md:right-6 z-50 bg-[#1c69d4] text-[#ffffff] font-extrabold px-5 py-3 shadow-2xl text-xs uppercase tracking-wider flex items-center justify-center md:justify-start gap-2 border border-white/20">
          <span>⚡ {toastMessage}</span>
        </div>
      )}

      {/* Main Container */}
      <main className="max-w-[1440px] mx-auto px-3 sm:px-4 lg:px-8 py-4 sm:py-8 pb-24 md:pb-8">
        {activeTab === 'finder' && (
          <StationExplorer
            stations={stations}
            allBookings={bookings}
            onSelectStationForBooking={(st, slot) => {
              if (!currentUser) {
                setIsAuthModalOpen(true);
                showToast('Silakan Registrasi atau Login terlebih dahulu untuk reservasi.');
              } else {
                setSelectedStationForBooking({ station: st, slot });
              }
            }}
          />
        )}

        {activeTab === 'garage' && (
          currentUser ? (
            <UserGarage
              vehicles={vehicles}
              onAddVehicle={handleAddVehicle}
              onDeleteVehicle={handleDeleteVehicle}
            />
          ) : (
            <div className="p-8 text-center bg-[#fafafa] border border-[#e6e6e6] my-8 max-w-lg mx-auto shadow-xs">
              <h3 className="text-base font-bold text-[#1a1a1a] uppercase tracking-wider mb-2">Akses Garasi EV Terbatas</h3>
              <p className="text-xs text-[#6b6b6b] leading-relaxed mb-6">
                Silakan Login atau Registrasi terlebih dahulu untuk mengelola daftar kendaraan listrik pribadi Anda.
              </p>
              <button
                onClick={() => setIsAuthModalOpen(true)}
                className="bmw-btn-primary text-xs"
              >
                LOGIN / REGISTRASI SEKARANG
              </button>
            </div>
          )
        )}

        {activeTab === 'bookings' && (
          currentUser ? (
            <UserBookingsList
              bookings={bookings.filter((b) => b.userId === currentUser?.id)}
              stations={stations}
              vehicles={vehicles}
              onStartSessionFromBooking={handleStartSessionFromBooking}
              onCancelBooking={(id) => {
                setBookings((prev) => prev.map((b) => (b.id === id ? { ...b, status: 'CANCELLED' } : b)));
                showToast('Booking dibatalkan.');
              }}
              onNavigateToFinder={() => setActiveTab('finder')}
            />
          ) : (
            <div className="p-8 text-center bg-[#fafafa] border border-[#e6e6e6] my-8 max-w-lg mx-auto shadow-xs">
              <h3 className="text-base font-bold text-[#1a1a1a] uppercase tracking-wider mb-2">Akses Booking Terbatas</h3>
              <p className="text-xs text-[#6b6b6b] leading-relaxed mb-6">
                Silakan Login atau Registrasi terlebih dahulu untuk melihat dan mengelola riwayat reservasi Anda.
              </p>
              <button
                onClick={() => setIsAuthModalOpen(true)}
                className="bmw-btn-primary text-xs"
              >
                LOGIN / REGISTRASI SEKARANG
              </button>
            </div>
          )
        )}

        {activeTab === 'session' && (
          currentUser ? (
            <LiveSessionViewer
              session={activeSession}
              stations={stations}
              onFinishSession={handleFinishSession}
            />
          ) : (
            <div className="p-8 text-center bg-[#fafafa] border border-[#e6e6e6] my-8 max-w-lg mx-auto shadow-xs">
              <h3 className="text-base font-bold text-[#1a1a1a] uppercase tracking-wider mb-2">Live Telemetry Terbatas</h3>
              <p className="text-xs text-[#6b6b6b] leading-relaxed mb-6">
                Silakan Login terlebih dahulu untuk memantau status charging aktif kendaraan Anda.
              </p>
              <button
                onClick={() => setIsAuthModalOpen(true)}
                className="bmw-btn-primary text-xs"
              >
                LOGIN / REGISTRASI SEKARANG
              </button>
            </div>
          )
        )}

        {activeTab === 'billing' && (
          currentUser ? (
            <BillingManager
              invoices={invoices}
              onConfirmPayment={handleConfirmPayment}
            />
          ) : (
            <div className="p-8 text-center bg-[#fafafa] border border-[#e6e6e6] my-8 max-w-lg mx-auto shadow-xs">
              <h3 className="text-base font-bold text-[#1a1a1a] uppercase tracking-wider mb-2">Akses Tagihan Terbatas</h3>
              <p className="text-xs text-[#6b6b6b] leading-relaxed mb-6">
                Silakan Login terlebih dahulu untuk melakukan pembayaran dan mengunduh invoice transaksi.
              </p>
              <button
                onClick={() => setIsAuthModalOpen(true)}
                className="bmw-btn-primary text-xs"
              >
                LOGIN / REGISTRASI SEKARANG
              </button>
            </div>
          )
        )}

        {activeTab === 'operator' && (
          currentUser?.role === 'OPERATOR' ? (
            <OperatorPortal
              stations={stations}
              onUpdateSlotStatus={handleUpdateSlotStatus}
              bookings={bookings}
              invoices={invoices}
              activeSession={activeSession}
              currentUser={currentUser}
              vehicles={vehicles}
            />
          ) : (
            <div className="p-8 text-center bg-[#fafafa] border border-[#e6e6e6] my-8 max-w-lg mx-auto shadow-xs">
              <h3 className="text-base font-bold text-[#1a1a1a] uppercase tracking-wider mb-2">Akses Operator Terbatas</h3>
              <p className="text-xs text-[#6b6b6b] leading-relaxed mb-6">
                Hanya Akun Operator yang berhak mengakses Laman Monitoring Realtime.
              </p>
              <button
                onClick={() => setIsAuthModalOpen(true)}
                className="bmw-btn-primary text-xs"
              >
                LOGIN SEBAGAI OPERATOR
              </button>
            </div>
          )
        )}

        {activeTab === 'admin' && (
          currentUser?.role === 'ADMIN' ? (
            <AdminPortal
              stations={stations}
              onUpdateSlotStatus={handleUpdateSlotStatus}
              onAddStation={handleAddStation}
              onUpdateStation={handleUpdateStation}
            />
          ) : (
            <div className="p-8 text-center bg-[#fafafa] border border-[#e6e6e6] my-8 max-w-lg mx-auto shadow-xs">
              <h3 className="text-base font-bold text-[#1a1a1a] uppercase tracking-wider mb-2">Akses Administrator Terbatas</h3>
              <p className="text-xs text-[#6b6b6b] leading-relaxed mb-6">
                Hanya Akun Admin Pusat SPKLU yang berhak mengelola stasiun dan data sistem.
              </p>
              <button
                onClick={() => setIsAuthModalOpen(true)}
                className="bmw-btn-primary text-xs"
              >
                LOGIN SEBAGAI ADMIN
              </button>
            </div>
          )
        )}
      </main>

      {/* Booking Modal Wizard */}
      {selectedStationForBooking && (
        <BookingModal
          station={selectedStationForBooking.station}
          initialSlot={selectedStationForBooking.slot}
          vehicles={vehicles}
          existingBookings={bookings}
          onClose={() => setSelectedStationForBooking(null)}
          onConfirmBooking={handleConfirmBooking}
          onConfirmAndStartSession={handleConfirmAndStartSession}
          onJoinWaitlist={handleJoinWaitlist}
        />
      )}

      {/* Auth Modal Login & Registration */}
      <AuthModal
        isOpen={isAuthModalOpen}
        onClose={() => setIsAuthModalOpen(false)}
        onLoginSuccess={(user) => {
          setCurrentUser(user);
          setIsAuthModalOpen(false);
        }}
        onRegister={handleRegister}
        onLogin={handleLogin}
      />
    </div>
  );
}

export default App;
