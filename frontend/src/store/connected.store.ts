import { defineStore } from 'pinia'
import { reactive, toRaw } from 'vue';
import type { Device } from "../models/device.models"

export type RootState = {
    devices: Device[];
    timestamp: number;
  };

const getTimestamp = (): number => {
  return Math.floor(Date.now() / 1000)
}

export const useDeviceStore = defineStore({
    id: "useDeviceStore",
    state: () =>
      ({
        devices: reactive([]) as Device[],
        timestamp: getTimestamp(),
      } as RootState),
    getters: {
      getDevices:(state) => {
        return state.devices
      }
    },
    actions: {
      updateStoreTimestamp() {
        // We call this at the start of every action so that
        // all items will have the same timestamp
        // and we can flush all irrelevant usb devices
        // such as the ones that have been disconnected.
        this.timestamp = getTimestamp()
      },
      createNewItem (item: Device) {
        if (!item) return;
        item.timestamp = this.timestamp;
        this.devices.push(item);
      },
      flushItems() {
        this.devices = [];
      },
      createOrUpdateItem(vendor_id: number, product_id: number, payload: Device) {
        
        if (!vendor_id || !product_id || !payload) return;
        
        const index = this.findIndexById(vendor_id, product_id);
        payload.timestamp = this.timestamp;
        if (index !== -1) {
            this.devices[index] = {...payload};
        }else {
          // push it if it doesnt exist.
          this.devices.push(payload);
        }
      },
      flushOldItems() {
        // call this at the end of updating the store
        // purge non existent usb devices
        // TODO: THINK OF A MORE EFFICENT WAY TO DO THIS
        this.devices = this.devices.filter(x => x.timestamp == this.timestamp)
        
      },
      deleteItem(vendor_id: number, product_id: number) {
        const index = this.findIndexById(vendor_id, product_id);
  
        if (index === -1) return;
  
        this.devices.splice(index, 1);
      },
  
      findIndexById(vendor_id: number, product_id: number) {
        return this.devices.findIndex((item) => item.vendor_id === vendor_id && product_id === product_id);
      },
    },
  });