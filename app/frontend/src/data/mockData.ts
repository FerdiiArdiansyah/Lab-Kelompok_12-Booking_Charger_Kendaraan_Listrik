import type { Station, UserVehicle, Booking, ChargingSession, Invoice } from '../types';

// ALL DUMMY MOCK DATA REMOVED - DATA IS FETCHED 100% DYNAMICALLY FROM REAL BACKEND POSTGRES DB
export const SAMPLE_STATIONS: Station[] = [];
export const SAMPLE_VEHICLES: UserVehicle[] = [];
export const SAMPLE_BOOKINGS: Booking[] = [];
export const SAMPLE_SESSION: ChargingSession | null = null;
export const SAMPLE_INVOICE: Invoice = {
  id: '',
  sessionId: '',
  userId: '',
  tariffId: '',
  consumedKwh: 0,
  pricePerKwh: 0,
  serviceFee: 0,
  subtotal: 0,
  tax: 0,
  total: 0,
  status: 'UNPAID',
  createdAt: new Date().toISOString(),
};
