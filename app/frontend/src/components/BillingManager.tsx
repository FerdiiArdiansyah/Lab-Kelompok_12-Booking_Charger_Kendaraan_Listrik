import React, { useState, useEffect } from 'react';
import { CreditCard, QrCode, CheckCircle, ShieldCheck, Download, Receipt } from 'lucide-react';
import type { Invoice, Payment } from '../types';
import { PaymentGatewayModal } from './PaymentGatewayModal';

interface BillingManagerProps {
  invoices: Invoice[];
  onConfirmPayment: (invoiceId: string, paymentMethod: Payment['paymentMethod']) => void;
}

export const BillingManager: React.FC<BillingManagerProps> = ({
  invoices,
  onConfirmPayment,
}) => {
  const [selectedInvoiceId, setSelectedInvoiceId] = useState<string>(
    invoices[0]?.id || ''
  );
  const [selectedMethod, setSelectedMethod] = useState<Payment['paymentMethod']>('QRIS');
  const [isGatewayOpen, setIsGatewayOpen] = useState(false);

  useEffect(() => {
    if (invoices.length > 0 && !invoices.some((inv) => inv.id === selectedInvoiceId)) {
      setSelectedInvoiceId(invoices[0].id);
    }
  }, [invoices, selectedInvoiceId]);

  const selectedInvoice = invoices.find((inv) => inv.id === selectedInvoiceId) || invoices[0] || null;

  const handleOpenGateway = () => {
    if (!selectedInvoice) return;
    setIsGatewayOpen(true);
  };

  const handleGatewaySuccess = (method: Payment['paymentMethod']) => {
    if (selectedInvoice) {
      onConfirmPayment(selectedInvoice.id, method);
    }
    setIsGatewayOpen(false);
  };

  const handleDownloadPdf = () => {
    if (!selectedInvoice) return;
    const rawKwhCost = Math.round(selectedInvoice.consumedKwh * selectedInvoice.pricePerKwh);
    const rawServiceFee = Math.round(selectedInvoice.serviceFee);
    const calcSubtotal = Math.round(selectedInvoice.subtotal || (rawKwhCost + rawServiceFee));
    const calcTax = Math.round(selectedInvoice.tax || (calcSubtotal * 0.11));
    const calcTotal = Math.round(selectedInvoice.total || (calcSubtotal + calcTax));

    const printWindow = window.open('', '_blank');
    if (!printWindow) return;
    printWindow.document.write(`
      <!DOCTYPE html>
      <html>
        <head>
          <title>Invoice SPKLU - ${selectedInvoice.id}</title>
          <style>
            body { font-family: 'Helvetica Neue', Arial, sans-serif; padding: 40px; color: #262626; background: #fff; }
            .header { border-bottom: 2px solid #1c69d4; padding-bottom: 16px; margin-bottom: 24px; display: flex; justify-content: space-between; align-items: flex-start; }
            .title { font-size: 22px; font-weight: 900; margin: 0; color: #1c69d4; }
            .subtitle { font-size: 12px; color: #6b6b6b; margin-top: 4px; }
            .meta { text-align: right; font-size: 12px; color: #3c3c3c; }
            .badge { display: inline-block; background: #22c55e; color: white; padding: 4px 10px; font-size: 11px; font-weight: bold; border-radius: 2px; text-transform: uppercase; margin-bottom: 8px; }
            table { width: 100%; border-collapse: collapse; margin-top: 24px; }
            th, td { text-align: left; padding: 12px; border-bottom: 1px solid #e6e6e6; font-size: 13px; }
            th { background: #fafafa; font-weight: 700; color: #6b6b6b; text-transform: uppercase; font-size: 11px; }
            .right { text-align: right; }
            .total-row td { font-size: 16px; font-weight: 900; border-top: 2px solid #262626; color: #262626; }
            .footer { margin-top: 40px; font-size: 11px; color: #6b6b6b; border-top: 1px solid #e6e6e6; padding-top: 16px; text-align: center; }
          </style>
        </head>
        <body>
          <div class="header">
            <div>
              <span class="badge">${selectedInvoice.status}</span>
              <h1 class="title">VoltHub SPKLU Billing Receipt</h1>
              <p class="subtitle">Bukti Resmi Transaksi Pengisian Daya Listrik EV (Permen ESDM No 7/2023)</p>
            </div>
            <div class="meta">
              <strong>Invoice ID: ${selectedInvoice.id}</strong><br/>
              <span>Tanggal: ${new Date(selectedInvoice.createdAt).toLocaleDateString('id-ID')}</span>
            </div>
          </div>
          <table>
            <thead>
              <tr>
                <th>Deskripsi Komponen</th>
                <th class="right">Jumlah (IDR)</th>
              </tr>
            </thead>
            <tbody>
              <tr>
                <td>Konsumsi Listrik (${selectedInvoice.consumedKwh} kWh @ Rp ${selectedInvoice.pricePerKwh.toLocaleString('id-ID')})</td>
                <td class="right">Rp ${rawKwhCost.toLocaleString('id-ID')}</td>
              </tr>
              <tr>
                <td>Biaya Layanan SPKLU (Fast Charging ESDM)</td>
                <td class="right">Rp ${rawServiceFee.toLocaleString('id-ID')}</td>
              </tr>
              <tr>
                <td>Subtotal</td>
                <td class="right">Rp ${calcSubtotal.toLocaleString('id-ID')}</td>
              </tr>
              <tr>
                <td>PPN (11%)</td>
                <td class="right">Rp ${calcTax.toLocaleString('id-ID')}</td>
              </tr>
              <tr class="total-row">
                <td>Total Tagihan LUNAS</td>
                <td class="right" style="color: #1c69d4;">Rp ${calcTotal.toLocaleString('id-ID')}</td>
              </tr>
            </tbody>
          </table>
          <div class="footer">
            Diterbitkan secara elektronik oleh VoltHub Billing Microservice (Port 8085). Dokumen ini merupakan bukti pembayaran yang sah.
          </div>
          <script>window.onload = function() { window.print(); }</script>
        </body>
      </html>
    `);
    printWindow.document.close();
  };

  return (
    <div className="space-y-6 max-w-5xl mx-auto">
      {/* Header */}
      <div className="bmw-card p-6 sm:p-8 flex flex-col md:flex-row items-center justify-between gap-4">
        <div>
          <span className="px-2.5 py-0.5 text-[10px] font-extrabold bg-[#1c69d4] text-white uppercase tracking-wider">
            PERMEN ESDM NO 7/2023 COMPLIANT
          </span>
          <h1 className="text-2xl sm:text-4xl font-black text-[#262626] mt-1">INVOICE & TAGIHAN SPKLU</h1>
          <p className="text-xs text-[#6b6b6b] font-light">Kelola riwayat pembayaran dan rincian transaksi pengisian daya listrik EV Anda.</p>
        </div>
      </div>

      {invoices.length === 0 ? (
        <div className="bmw-card p-12 text-center space-y-4 border-dashed">
          <Receipt className="w-12 h-12 text-[#1c69d4] mx-auto" />
          <div>
            <h3 className="text-lg font-bold text-[#262626]">Belum Ada Tagihan SPKLU</h3>
            <p className="text-xs text-[#6b6b6b] max-w-md mx-auto mt-1">
              Anda belum memiliki riwayat tagihan pengisian daya. Selesaikan sesi pengisian daya di menu <strong>Live Telemetry</strong> untuk menerbitkan invoice.
            </p>
          </div>
        </div>
      ) : (
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Left: Invoice Selection List */}
          <div className="space-y-3">
            <span className="text-xs font-bold text-[#6b6b6b] uppercase tracking-wider block">Daftar Tagihan ({invoices.length})</span>
            {invoices.map((inv) => {
              const isSelected = selectedInvoice?.id === inv.id;
              const rawTotal = Math.round(inv.total);

              return (
                <div
                  key={inv.id}
                  onClick={() => setSelectedInvoiceId(inv.id)}
                  className={`p-4 cursor-pointer transition-colors border ${
                    isSelected
                      ? 'bg-[#fafafa] border-[#1c69d4] shadow-sm'
                      : 'bg-white border-[#e6e6e6] hover:border-[#cccccc]'
                  }`}
                >
                  <div className="flex items-center justify-between">
                    <span className="text-xs font-bold text-[#262626]">Invoice #{inv.id}</span>
                    <span
                      className={`px-2 py-0.5 text-[10px] font-bold uppercase tracking-wider border ${
                        inv.status === 'PAID'
                          ? 'bg-[#22c55e]/10 text-[#22c55e] border-[#22c55e]/30'
                          : 'bg-[#f59e0b]/10 text-[#f59e0b] border-[#f59e0b]/30'
                      }`}
                    >
                      {inv.status}
                    </span>
                  </div>
                  <div className="mt-2 flex items-baseline justify-between">
                    <span className="text-xs text-[#6b6b6b] font-light">{inv.consumedKwh} kWh</span>
                    <span className="text-sm font-black text-[#1c69d4]">
                      Rp {rawTotal.toLocaleString('id-ID')}
                    </span>
                  </div>
                </div>
              );
            })}
          </div>

          {/* Right: Selected Invoice Receipt & Payment Gateway */}
          {selectedInvoice ? (
            <div className="lg:col-span-2 space-y-6">
              {/* Digital Receipt */}
              <div className="bmw-card p-6 sm:p-8 space-y-6">
                <div className="flex items-center justify-between border-b border-[#e6e6e6] pb-4">
                  <div>
                    <div className="text-xs text-[#6b6b6b]">Rincian Resmi Tagihan</div>
                    <h2 className="text-xl font-extrabold text-[#262626]">SPKLU Billing Receipt</h2>
                  </div>
                  <div className="text-right">
                    <div className="text-xs font-bold text-[#1c69d4]"># {selectedInvoice.id}</div>
                    <div className="text-[10px] text-[#6b6b6b] font-light">{new Date(selectedInvoice.createdAt).toLocaleDateString('id-ID')}</div>
                  </div>
                </div>

                {/* Line Items Table with Rounded Numbers */}
                {(() => {
                  const rawKwhCost = Math.round(selectedInvoice.consumedKwh * selectedInvoice.pricePerKwh);
                  const rawServiceFee = Math.round(selectedInvoice.serviceFee);
                  const calcSubtotal = Math.round(selectedInvoice.subtotal || (rawKwhCost + rawServiceFee));
                  const calcTax = Math.round(selectedInvoice.tax || (calcSubtotal * 0.11));
                  const calcTotal = Math.round(selectedInvoice.total || (calcSubtotal + calcTax));

                  return (
                    <div className="space-y-3 text-xs">
                      <div className="flex justify-between text-[#6b6b6b] font-bold border-b border-[#e6e6e6] pb-2 uppercase tracking-wider">
                        <span>Deskripsi Komponen</span>
                        <span>Jumlah (IDR)</span>
                      </div>
                      <div className="flex justify-between text-[#3c3c3c]">
                        <span>Konsumsi Listrik ({selectedInvoice.consumedKwh} kWh x Rp {Math.round(selectedInvoice.pricePerKwh).toLocaleString('id-ID')})</span>
                        <span className="font-bold text-[#262626]">Rp {rawKwhCost.toLocaleString('id-ID')}</span>
                      </div>
                      <div className="flex justify-between text-[#3c3c3c]">
                        <span>Biaya Layanan SPKLU (Fast Charging ESDM)</span>
                        <span className="font-bold text-[#262626]">Rp {rawServiceFee.toLocaleString('id-ID')}</span>
                      </div>
                      <div className="flex justify-between text-[#6b6b6b] pt-1">
                        <span>Subtotal</span>
                        <span>Rp {calcSubtotal.toLocaleString('id-ID')}</span>
                      </div>
                      <div className="flex justify-between text-[#6b6b6b]">
                        <span>PPN (11%)</span>
                        <span>Rp {calcTax.toLocaleString('id-ID')}</span>
                      </div>
                      <div className="flex justify-between text-lg font-black text-[#262626] pt-3 border-t border-[#262626]">
                        <span>Total Tagihan</span>
                        <span className="text-[#1c69d4]">
                          Rp {calcTotal.toLocaleString('id-ID')}
                        </span>
                      </div>
                    </div>
                  );
                })()}

                {/* Payment Gateway Options */}
                {selectedInvoice.status === 'UNPAID' ? (
                  <div className="bg-[#fafafa] p-5 border border-[#e6e6e6] space-y-4 pt-4">
                    <h3 className="text-xs font-bold text-[#262626] uppercase tracking-wider flex items-center gap-2">
                      <CreditCard className="w-4 h-4 text-[#1c69d4]" /> Pilih Metode Pembayaran Instan
                    </h3>

                    <div className="grid grid-cols-2 sm:grid-cols-5 gap-2">
                      {[
                        { id: 'QRIS', label: 'QRIS Standard', icon: QrCode },
                        { id: 'VA_BCA', label: 'VA BCA', icon: CreditCard },
                        { id: 'VA_MANDIRI', label: 'VA Mandiri', icon: CreditCard },
                        { id: 'E_WALLET_GOPAY', label: 'GoPay / E-Wallet', icon: QrCode },
                        { id: 'CREDIT_CARD', label: 'Kartu Kredit/Debit', icon: CreditCard },
                      ].map((m) => {
                        const Icon = m.icon;
                        return (
                          <button
                            key={m.id}
                            onClick={() => setSelectedMethod(m.id as any)}
                            className={`p-3 border text-xs font-bold text-center transition-colors flex flex-col items-center justify-center gap-1.5 cursor-pointer ${
                              selectedMethod === m.id
                                ? 'bg-[#1c69d4] text-white border-[#1c69d4] shadow-xs'
                                : 'bg-white text-[#262626] border-[#cccccc] hover:border-[#262626]'
                            }`}
                          >
                            <Icon className="w-4 h-4" />
                            <span>{m.label}</span>
                          </button>
                        );
                      })}
                    </div>

                    {selectedMethod === 'QRIS' && (
                      <div className="p-4 bg-white border border-[#e6e6e6] text-center space-y-2">
                        <div className="w-32 h-32 mx-auto bg-white p-2 border border-[#262626] flex items-center justify-center">
                          <QrCode className="w-28 h-28 text-[#262626]" />
                        </div>
                        <p className="text-[11px] text-[#6b6b6b]">Scan kode QRIS nasional untuk pembayaran otomatis 1-click</p>
                      </div>
                    )}

                    <button
                      onClick={handleOpenGateway}
                      className="bmw-btn-primary w-full py-3 text-sm font-black flex items-center justify-center gap-2 cursor-pointer shadow-md"
                    >
                      <ShieldCheck className="w-4 h-4" />
                      <span>BUKA PAYMENT GATEWAY CHECKOUT (RP {Math.round(selectedInvoice.total).toLocaleString('id-ID')})</span>
                    </button>
                  </div>
                ) : (
                  <div className="p-4 bg-[#22c55e]/10 border border-[#22c55e] text-[#22c55e] text-xs flex items-center justify-between font-bold">
                    <span className="flex items-center gap-2">
                      <CheckCircle className="w-5 h-5 text-[#22c55e]" /> Tagihan LUNAS terverifikasi oleh Billing Microservice.
                    </span>
                    <button
                      onClick={handleDownloadPdf}
                      className="px-3 py-1.5 bg-[#22c55e] hover:bg-[#16a34a] text-white text-[11px] font-bold flex items-center gap-1.5 transition-colors cursor-pointer"
                    >
                      <Download className="w-3.5 h-3.5" /> UNDUH STRUK PDF
                    </button>
                  </div>
                )}
              </div>
            </div>
          ) : null}
        </div>
      )}

      {/* Standard Payment Gateway Modal */}
      {isGatewayOpen && selectedInvoice && (
        <PaymentGatewayModal
          invoice={selectedInvoice}
          paymentMethod={selectedMethod}
          onClose={() => setIsGatewayOpen(false)}
          onSuccess={handleGatewaySuccess}
        />
      )}
    </div>
  );
};
