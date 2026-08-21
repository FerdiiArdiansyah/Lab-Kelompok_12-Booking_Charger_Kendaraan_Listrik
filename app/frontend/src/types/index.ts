export interface User {
  id: string;
  name: string;
  email: string;
  phone: string;
  role: 'USER' | 'ADMIN' | 'OPERATOR';
  status: 'ACTIVE' | 'SUSPENDED';
}

export interface UserVehicle {
  id: string;
  userId: string;
  brand: string;
  model: string;
  licensePlate: string;
  connectorType: 'CCS2' | 'CHAdeMO' | 'Type 2' | 'GB/T';
  batteryCapacityKwh: number;
  createdAt: string;
}

export interface ChargerSlot {
  id: string;
  stationId: string;
  slotNumber: number;
  connectorType: 'CCS2' | 'CHAdeMO' | 'Type 2' | 'GB/T';
  maxPowerKw: number;
  status: 'AVAILABLE' | 'IN_USE' | 'OUT_OF_SERVICE';
}

export interface Tariff {
  id: string;
  stationId: string;
  pricePerKwh: number;
  currency: string;
  validFrom: string;
  isActive: boolean;
}

export interface Station {
  id: string;
  name: string;
  location: string;
  latitude: number;
  longitude: number;
  totalPower: number;
  status: 'ACTIVE' | 'MAINTENANCE' | 'OFFLINE';
  slots: ChargerSlot[];
  activeTariff?: Tariff;
  mapUrl?: string;
  imageUrl?: string;
}

export interface Booking {
  id: string;
  userId: string;
  slotId: string;
  stationId: string;
  vehicleId: string;
  startTime: string;
  endTime: string;
  status: 'PENDING' | 'CONFIRMED' | 'IN_PROGRESS' | 'COMPLETED' | 'CANCELLED' | 'EXPIRED';
  createdAt: string;
}

export interface Waitlist {
  id: string;
  userId: string;
  stationId: string;
  slotId: string;
  preferredStartTime: string;
  preferredEndTime: string;
  status: 'WAITING' | 'NOTIFIED' | 'EXPIRED' | 'CANCELLED';
  createdAt: string;
}

export interface MeterReading {
  id: number;
  sessionId: string;
  recordedAt: string;
  currentKwh: number;
  powerKw: number;
  voltage: number;
  currentAmpere: number;
}

export interface ChargingSession {
  id: string;
  bookingId: string;
  slotId: string;
  userId: string;
  startedAt: string;
  endedAt?: string;
  consumedKwh: number;
  status: 'IN_PROGRESS' | 'COMPLETED' | 'INTERRUPTED' | 'CANCELLED';
  readings?: MeterReading[];
}

export interface Payment {
  id: string;
  invoiceId: string;
  paymentMethod: 'QRIS' | 'VA_BCA' | 'VA_MANDIRI' | 'E_WALLET_GOPAY' | 'CREDIT_CARD';
  amount: number;
  status: 'PENDING' | 'SUCCESS' | 'FAILED';
  transactionRef?: string;
  paidAt?: string;
  createdAt: string;
}

export interface Invoice {
  id: string;
  sessionId: string;
  userId: string;
  tariffId: string;
  consumedKwh: number;
  pricePerKwh: number;
  serviceFee: number;
  subtotal: number;
  tax: number;
  total: number;
  status: 'UNPAID' | 'PAID' | 'VOID' | 'REFUNDED';
  createdAt: string;
  payments?: Payment[];
}
