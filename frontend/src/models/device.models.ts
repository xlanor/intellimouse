export interface Device {
  path: string;
  vendor_id: number;
  product_id: number;
  serial: string;
  manufacturer: string;
  product: string;
  interface: number;
  timestamp: number;
}