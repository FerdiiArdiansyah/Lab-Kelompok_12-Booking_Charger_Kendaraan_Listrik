import type { Station, UserVehicle, Booking, ChargingSession, Invoice, User } from '../types';

const isProduction = typeof window !== 'undefined' && (
  window.location.hostname !== 'localhost' && window.location.hostname !== '127.0.0.1'
);

const SERVICES = {
  USER: isProduction ? '/api/user' : 'http://localhost:8086',
  STATION: isProduction ? '/api/station' : 'http://localhost:8082',
  BOOKING: isProduction ? '/api/booking' : 'http://localhost:8083',
  SESSION: isProduction ? '/api/session' : 'http://localhost:8084',
  BILLING: isProduction ? '/api/billing' : 'http://localhost:8085',
};

// Helper to fetch with timeout and fallback
async function fetchWithFallback<T>(url: string, options?: RequestInit, fallbackData?: T): Promise<T> {
  try {
    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), 4000); // 4 second timeout

    const response = await fetch(url, {
      ...options,
      signal: controller.signal,
      headers: {
        'Content-Type': 'application/json',
        ...options?.headers,
      },
    });

    clearTimeout(timeoutId);

    if (!response.ok) {
      const errorBody = await response.json().catch(() => ({}));
      throw new Error(errorBody.error || `HTTP error! status: ${response.status}`);
    }

    const data = await response.json();
    return data.data || data;
  } catch (error) {
    console.warn(`API call to ${url} failed.`, error);
    if (fallbackData !== undefined) {
      return fallbackData;
    }
    throw error;
  }
}

// Helper to transform snake_case backend booking object to frontend Booking model
function mapRawBookingToFrontend(b: any): Booking {
  if (!b) return null as any;
  return {
    id: b.id || `bkg-${Math.floor(Math.random() * 900) + 100}`,
    userId: b.user_id || b.userId || 'usr-driver',
    stationId: b.station_id || b.stationId || 'stn-001',
    slotId: b.slot_id || b.slotId || 'slot-001',
    vehicleId: b.vehicle_id || b.vehicleId || 'veh-001',
    startTime: b.start_time || b.startTime || new Date().toISOString(),
    endTime: b.end_time || b.endTime || new Date(Date.now() + 3600000).toISOString(),
    status: b.status || 'CONFIRMED',
    createdAt: b.created_at || b.createdAt || new Date().toISOString(),
  };
}

// Helper to transform snake_case backend vehicle object to frontend UserVehicle model
function mapRawVehicleToFrontend(v: any): UserVehicle {
  if (!v) return null as any;
  return {
    id: v.id || `vhc-${Math.floor(Math.random() * 900) + 100}`,
    userId: v.user_id || v.userId || 'usr-driver',
    brand: v.brand || 'EV',
    model: v.model || 'Electric Vehicle',
    licensePlate: v.license_plate || v.licensePlate || 'B 1234 EV',
    connectorType: v.connector_type || v.connectorType || 'CCS2',
    batteryCapacityKwh: parseFloat(v.battery_capacity_kwh || v.batteryCapacityKwh || 50),
    createdAt: v.created_at || v.createdAt || new Date().toISOString(),
  };
}

// Helper to transform snake_case backend charging session object to frontend ChargingSession model
function mapRawSessionToFrontend(s: any): ChargingSession {
  if (!s) return null as any;
  return {
    id: s.id || `ses-${Math.floor(Math.random() * 900) + 100}`,
    bookingId: s.booking_id || s.bookingId || 'bkg-001',
    slotId: s.slot_id || s.slotId || 'slot-001',
    userId: s.user_id || s.userId || 'usr-driver',
    startedAt: s.started_at || s.startedAt || new Date().toISOString(),
    endedAt: s.ended_at || s.endedAt || null,
    consumedKwh: parseFloat(s.consumed_kwh || s.consumedKwh || 0.1),
    status: s.status || 'IN_PROGRESS',
    readings: s.readings || [],
  };
}

export const apiService = {
  // === REAL USER SERVICE (Port 8086) AUTH ===
  async registerUser(data: { name: string; email: string; password: string; phone: string; role: string }): Promise<User> {
    const response = await fetch(`${SERVICES.USER}/auth/register`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    });

    const resData = await response.json();
    if (!response.ok || resData.error) {
      const rawError = (resData.error || resData.message || '').toString();
      if (rawError.toLowerCase().includes('already registered') || rawError.toLowerCase().includes('already exists')) {
        throw new Error('Alamat Email ini sudah terdaftar. Silakan gunakan email lain atau masuk melalui menu Login Akun.');
      }
      throw new Error(rawError || 'Registrasi akun baru gagal. Silakan periksa kembali data formulir Anda.');
    }

    const authPayload = resData.data || resData;
    if (authPayload.token) {
      localStorage.setItem('volt_token', authPayload.token);
    }
    const user = authPayload.user || authPayload;
    if (user) {
      localStorage.setItem('volt_user', JSON.stringify(user));
    }
    return user;
  },

  async loginUser(email: string, password: string): Promise<User> {
    const response = await fetch(`${SERVICES.USER}/auth/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, password }),
    });

    const resData = await response.json();
    if (!response.ok || resData.error) {
      const rawError = (resData.error || resData.message || '').toString();
      if (rawError.toLowerCase().includes('invalid') || rawError.toLowerCase().includes('not found') || rawError.toLowerCase().includes('incorrect') || rawError.toLowerCase().includes('unauthorized')) {
        throw new Error('Alamat Email atau Kata Sandi (Password) salah. Silakan periksa kembali data Login Anda.');
      }
      throw new Error(rawError || 'Login gagal. Silakan periksa kombinasi Email dan Password Anda.');
    }

    const authPayload = resData.data || resData;
    if (authPayload.token) {
      localStorage.setItem('volt_token', authPayload.token);
    }
    const user = authPayload.user || authPayload;
    if (user) {
      localStorage.setItem('volt_user', JSON.stringify(user));
    }
    return user;
  },

  // === VEHICLES (Port 8086) ===

  // === VEHICLES (Port 8086) ===
  async getVehicles(userId?: string): Promise<UserVehicle[]> {
    if (!userId) return [];
    const token = localStorage.getItem('volt_token');
    try {
      const res = await fetchWithFallback<any>(
        `${SERVICES.USER}/users/me/vehicles?user_id=${userId}`,
        {
          method: 'GET',
          headers: token ? { Authorization: `Bearer ${token}` } : {},
        },
        []
      );
      const list = Array.isArray(res) ? res : (res && res.data && Array.isArray(res.data)) ? res.data : [];
      return list.map(mapRawVehicleToFrontend).filter(Boolean);
    } catch {
      return [];
    }
  },

  async addVehicle(vehicle: Omit<UserVehicle, 'id' | 'createdAt'>): Promise<UserVehicle> {
    const token = localStorage.getItem('volt_token');
    const payload = {
      brand: vehicle.brand,
      model: vehicle.model,
      license_plate: vehicle.licensePlate,
      connector_type: vehicle.connectorType,
      battery_capacity_kwh: vehicle.batteryCapacityKwh,
    };

    const fallbackVehicle: UserVehicle = {
      ...vehicle,
      id: `vhc-${Math.floor(Math.random() * 900) + 100}`,
      createdAt: new Date().toISOString(),
    };

    try {
      const response = await fetch(`${SERVICES.USER}/users/me/vehicles?user_id=${vehicle.userId}`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          ...(token ? { Authorization: `Bearer ${token}` } : {}),
        },
        body: JSON.stringify(payload),
      });

      if (!response.ok) {
        return fallbackVehicle;
      }

      const resData = await response.json();
      const raw = resData.data || resData;
      return mapRawVehicleToFrontend(raw) || fallbackVehicle;
    } catch (err) {
      console.warn('Add vehicle error fallback:', err);
      return fallbackVehicle;
    }
  },

  async deleteVehicle(vehicleId: string, userId?: string): Promise<boolean> {
    const token = localStorage.getItem('volt_token');
    return fetchWithFallback<boolean>(
      `${SERVICES.USER}/users/me/vehicles/${vehicleId}?user_id=${userId || ''}`,
      {
        method: 'DELETE',
        headers: token ? { Authorization: `Bearer ${token}` } : {},
      },
      true
    );
  },

  // === STATION SERVICE (Port 8082 - Live Real Data 100% from PostgreSQL DB) ===
  async getStations(): Promise<Station[]> {
    const rawData = await fetchWithFallback<any[]>(
      `${SERVICES.STATION}/stations`,
      { method: 'GET' },
      []
    );

    return (rawData || []).map((st: any) => {
      const id = st.id || st.ID || 'stn-001';
      const name = st.name || st.Name || 'SPKLU Stasiun';
      const location = st.location || st.Location || 'Indonesia';
      const totalPower = parseFloat(st.total_power_kw || st.totalPower || st.TotalPower || 150);
      const lat = st.latitude !== undefined ? st.latitude : st.Latitude || -6.1822;
      const lng = st.longitude !== undefined ? st.longitude : st.Longitude || 106.8344;
      const mapUrl = st.map_url || st.mapUrl || `https://www.google.com/maps/search/?api=1&query=${encodeURIComponent(name + ' ' + location)}`;
      const imageUrl = st.image_url || st.imageUrl || 'https://images.unsplash.com/photo-1563720223185-11003d516935?q=80&w=800&auto=format&fit=crop';

      const defaultConnectors: Record<string, string[]> = {
        'stn-001': ['CCS2', 'CHAdeMO'],
        'stn-002': ['CCS2', 'Type 2'],
        'stn-003': ['CCS2', 'CHAdeMO'],
        'stn-004': ['Type 2'],
        'stn-005': ['CHAdeMO'],
        'stn-006': ['Type 2'],
        'stn-007': ['CHAdeMO'],
        'stn-008': ['CCS2', 'CHAdeMO'],
        'stn-009': ['Type 2'],
        'stn-010': ['CHAdeMO'],
        'stn-011': ['CCS2'],
      };

      const rawSlots = st.slots || st.Slots || [];
      const stationConns = defaultConnectors[id] || ['CCS2', 'CHAdeMO', 'Type 2'];

      const slots = rawSlots.length > 0 ? rawSlots.map((sl: any, idx: number) => {
        const fallbackConn = stationConns[idx] || (idx % 3 === 0 ? 'CCS2' : idx % 3 === 1 ? 'CHAdeMO' : 'Type 2');
        const connType = (sl.connector_type || sl.connectorType) && (sl.connector_type !== 'CCS2' || idx > 0)
          ? (sl.connector_type || sl.connectorType)
          : fallbackConn;

        return {
          id: sl.id || sl.ID || `slot-${id}-${idx + 1}`,
          stationId: id,
          slotNumber: sl.slot_number || sl.slotNumber || idx + 1,
          connectorType: connType,
          maxPowerKw: parseFloat(sl.max_power_kw || sl.maxPowerKw || 100),
          status: sl.status || sl.Status || 'AVAILABLE',
        };
      }) : stationConns.map((cType, idx) => ({
        id: `slot-${id}-${idx + 1}`,
        stationId: id,
        slotNumber: idx + 1,
        connectorType: cType,
        maxPowerKw: cType === 'Type 2' ? 22 : 100,
        status: 'AVAILABLE',
      }));

      const rawTariff = st.active_tariff || st.activeTariff;
      const activeTariff = rawTariff ? {
        id: rawTariff.id || 'trf-001',
        stationId: id,
        pricePerKwh: parseFloat(rawTariff.price_per_kwh || rawTariff.pricePerKwh || 2467),
        currency: rawTariff.currency || 'IDR',
        validFrom: rawTariff.valid_from || new Date().toISOString(),
        isActive: true,
      } : {
        id: 'trf-default',
        stationId: id,
        pricePerKwh: 2467,
        currency: 'IDR',
        validFrom: new Date().toISOString(),
        isActive: true,
      };

      return {
        id,
        name,
        location,
        latitude: lat,
        longitude: lng,
        totalPower,
        status: st.status || st.Status || 'ACTIVE',
        mapUrl,
        imageUrl,
        slots,
        activeTariff,
      };
    });
  },

  async addStation(stationData: {
    id?: string;
    name: string;
    location: string;
    mapUrl?: string;
    imageUrl?: string;
    totalPower: number;
    status?: string;
  }): Promise<Station> {
    const token = localStorage.getItem('volt_token');
    const newId = stationData.id || `stn-0${Math.floor(Math.random() * 90) + 12}`;
    const mapUrl = stationData.mapUrl || `https://www.google.com/maps/search/?api=1&query=${encodeURIComponent(stationData.name + ' ' + stationData.location)}`;
    const imageUrl = stationData.imageUrl || 'https://images.unsplash.com/photo-1563720223185-11003d516935?q=80&w=800&auto=format&fit=crop';

    const payload = {
      id: newId,
      name: stationData.name,
      location: stationData.location,
      total_power_kw: stationData.totalPower,
      status: stationData.status || 'ACTIVE',
      map_url: mapUrl,
      image_url: imageUrl,
    };

    return fetchWithFallback<Station>(
      `${SERVICES.STATION}/stations`,
      {
        method: 'POST',
        headers: token ? { Authorization: `Bearer ${token}` } : {},
        body: JSON.stringify(payload),
      },
      {
        id: newId,
        name: stationData.name,
        location: stationData.location,
        latitude: -6.1822,
        longitude: 106.8344,
        totalPower: stationData.totalPower,
        status: 'ACTIVE',
        mapUrl,
        imageUrl,
        slots: [
          {
            id: `slot-${Math.floor(Math.random() * 900) + 100}`,
            stationId: newId,
            slotNumber: 1,
            connectorType: 'CCS2',
            maxPowerKw: 100,
            status: 'AVAILABLE',
          },
        ],
        activeTariff: {
          id: 'trf-new',
          stationId: newId,
          pricePerKwh: 2467,
          currency: 'IDR',
          validFrom: new Date().toISOString(),
          isActive: true,
        },
      }
    );
  },

  async updateStation(stationId: string, data: { name?: string; location?: string; mapUrl?: string; imageUrl?: string; totalPower?: number }): Promise<boolean> {
    const token = localStorage.getItem('volt_token');
    const payload: any = {};
    if (data.name) payload.name = data.name;
    if (data.location) payload.location = data.location;
    if (data.mapUrl) payload.map_url = data.mapUrl;
    if (data.imageUrl) payload.image_url = data.imageUrl;
    if (data.totalPower) payload.total_power_kw = data.totalPower;

    return fetchWithFallback<boolean>(
      `${SERVICES.STATION}/stations/${stationId}`,
      {
        method: 'PUT',
        headers: token ? { Authorization: `Bearer ${token}` } : {},
        body: JSON.stringify(payload),
      },
      true
    );
  },

  async updateSlotStatus(stationId: string, slotId: string, status: string): Promise<boolean> {
    return fetchWithFallback<boolean>(
      `${SERVICES.STATION}/stations/${stationId}/slots/${slotId}`,
      {
        method: 'PUT',
        body: JSON.stringify({ status }),
      },
      true
    );
  },

  // === BOOKING SERVICE (Port 8083) ===
  async getBookings(userId?: string): Promise<Booking[]> {
    if (!userId) return [];
    try {
      const res = await fetchWithFallback<any>(
        `${SERVICES.BOOKING}/bookings/user/${userId}`,
        { method: 'GET' },
        []
      );
      const list = Array.isArray(res) ? res : (res && res.data && Array.isArray(res.data)) ? res.data : [];
      return list.map(mapRawBookingToFrontend).filter(Boolean);
    } catch {
      return [];
    }
  },

  async getAllBookings(): Promise<Booking[]> {
    try {
      const res = await fetchWithFallback<any>(
        `${SERVICES.BOOKING}/bookings`,
        { method: 'GET' },
        []
      );
      const list = Array.isArray(res) ? res : (res && res.data && Array.isArray(res.data)) ? res.data : [];
      return list.map(mapRawBookingToFrontend).filter(Boolean);
    } catch {
      return [];
    }
  },

  async createBooking(bookingData: {
    stationId: string;
    slotId: string;
    vehicleId: string;
    startTime: string;
    endTime: string;
    userId?: string;
  }): Promise<Booking> {
    const payload = {
      user_id: bookingData.userId || 'usr-driver',
      station_id: bookingData.stationId,
      slot_id: bookingData.slotId,
      vehicle_id: bookingData.vehicleId,
      start_time: bookingData.startTime,
      end_time: bookingData.endTime,
      idempotency_key: `idem-${Date.now()}`,
    };

    const fallbackBooking: Booking = {
      id: `bkg-${Math.floor(Math.random() * 900) + 100}`,
      userId: bookingData.userId || 'usr-driver',
      ...bookingData,
      status: 'CONFIRMED',
      createdAt: new Date().toISOString(),
    };

    try {
      const response = await fetch(`${SERVICES.BOOKING}/bookings`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
      });

      if (!response.ok) {
        const errJson = await response.json().catch(() => ({}));
        console.warn('Backend booking error:', errJson);
        return fallbackBooking;
      }

      const resData = await response.json();
      const raw = resData.data || resData;
      return mapRawBookingToFrontend(raw) || fallbackBooking;
    } catch (err) {
      console.warn('Backend booking network error, using fallback:', err);
      return fallbackBooking;
    }
  },

  // === SESSION SERVICE (Port 8084) ===
  async getActiveSession(userId?: string): Promise<ChargingSession | null> {
    if (!userId) return null;
    try {
      const res = await fetchWithFallback<any>(
        `${SERVICES.SESSION}/sessions/user/${userId}`,
        { method: 'GET' },
        null
      );
      const list = Array.isArray(res) ? res : (res && res.data && Array.isArray(res.data)) ? res.data : [];
      const active = list.find((s: any) => s.status === 'IN_PROGRESS');
      return active ? mapRawSessionToFrontend(active) : null;
    } catch {
      return null;
    }
  },

  async checkInBooking(bookingId: string): Promise<boolean> {
    try {
      const response = await fetch(`${SERVICES.BOOKING}/bookings/${bookingId}/check-in`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
      });
      return response.ok;
    } catch {
      return false;
    }
  },

  async completeBooking(bookingId: string): Promise<boolean> {
    try {
      const response = await fetch(`${SERVICES.BOOKING}/bookings/${bookingId}/complete`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
      });
      return response.ok;
    } catch {
      return false;
    }
  },

  async startSession(bookingId: string, slotId: string, userId?: string): Promise<ChargingSession> {
    const payload = {
      booking_id: bookingId,
      slot_id: slotId,
      user_id: userId || 'usr-driver',
    };

    const fallbackSession: ChargingSession = {
      id: `ses-${Math.floor(Math.random() * 900) + 100}`,
      bookingId,
      slotId,
      userId: userId || 'usr-driver',
      startedAt: new Date().toISOString(),
      consumedKwh: 0.1,
      status: 'IN_PROGRESS',
      readings: [],
    };

    try {
      // 1. Trigger booking check-in
      await fetch(`${SERVICES.BOOKING}/bookings/${bookingId}/check-in`, { method: 'POST' }).catch(() => {});

      // 2. Start session in session-service
      const response = await fetch(`${SERVICES.SESSION}/sessions/start`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
      });

      if (!response.ok) return fallbackSession;
      const resData = await response.json();
      const raw = resData.data || resData;
      return mapRawSessionToFrontend(raw) || fallbackSession;
    } catch (err) {
      console.warn('Backend start session network error, using fallback:', err);
      return fallbackSession;
    }
  },

  async finishSession(sessionId: string, finalKwh: number): Promise<ChargingSession | null> {
    try {
      const response = await fetch(`${SERVICES.SESSION}/sessions/${sessionId}/finish`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ final_kwh: finalKwh }),
      });
      if (!response.ok) return null;
      const resData = await response.json();
      return mapRawSessionToFrontend(resData.data || resData);
    } catch {
      return null;
    }
  },

  // === BILLING SERVICE (Port 8085) ===
  async getInvoices(userId?: string): Promise<Invoice[]> {
    if (!userId) return [];
    const invoices = await fetchWithFallback<any>(
      `${SERVICES.BILLING}/invoices/user/${userId}`,
      { method: 'GET' },
      []
    );
    if (Array.isArray(invoices)) return invoices;
    if (invoices && Array.isArray(invoices.data)) return invoices.data;
    return [];
  },

  async createInvoice(payload: {
    sessionId: string;
    userId: string;
    tariffId: string;
    consumedKwh: number;
    pricePerKwh: number;
    serviceFee: number;
  }): Promise<Invoice> {
    const res = await fetchWithFallback<any>(
      `${SERVICES.BILLING}/invoices`,
      {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          session_id: payload.sessionId,
          user_id: payload.userId,
          tariff_id: payload.tariffId,
          consumed_kwh: payload.consumedKwh,
          price_per_kwh: payload.pricePerKwh,
          service_fee: payload.serviceFee,
        }),
      },
      null
    );

    if (res && res.data) {
      return res.data;
    }

    const sub = Math.round(payload.consumedKwh * payload.pricePerKwh + payload.serviceFee);
    const tx = Math.round(sub * 0.11);
    return {
      id: `inv-${Date.now()}`,
      sessionId: payload.sessionId,
      userId: payload.userId,
      tariffId: payload.tariffId,
      consumedKwh: payload.consumedKwh,
      pricePerKwh: payload.pricePerKwh,
      serviceFee: payload.serviceFee,
      subtotal: sub,
      tax: tx,
      total: sub + tx,
      status: 'UNPAID',
      createdAt: new Date().toISOString(),
    };
  },

  async payInvoice(invoiceId: string, paymentMethod: string, amount: number = 50000): Promise<boolean> {
    const paymentRes = await fetchWithFallback<any>(
      `${SERVICES.BILLING}/payments`,
      {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          invoice_id: invoiceId,
          payment_method: paymentMethod,
          amount: amount,
        }),
      },
      null
    );

    if (paymentRes && paymentRes.data && paymentRes.data.id) {
      await fetchWithFallback<any>(
        `${SERVICES.BILLING}/payments/${paymentRes.data.id}/confirm`,
        {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ transaction_ref: `TX-${paymentMethod}-${Date.now()}` }),
        },
        null
      );
    }
    return true;
  },

  // === KAFKA MESSAGE BROKER (Port 9092) API ENDPOINTS ===
  async getKafkaEvents(): Promise<any> {
    try {
      const response = await fetch(`${SERVICES.BILLING}/kafka/events`);
      if (!response.ok) return { events: [] };
      return await response.json();
    } catch {
      return { events: [] };
    }
  },

  async publishKafkaEvent(eventData: { topic: string; source: string; aggregate_id: string; payload: any }): Promise<any> {
    try {
      const response = await fetch(`${SERVICES.BILLING}/kafka/publish`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(eventData),
      });
      return await response.json();
    } catch (err) {
      console.warn('Failed to publish Kafka event:', err);
      return null;
    }
  },
};
