<script lang="ts" setup>

  import { ref, onMounted, reactive, computed, defineEmits } from 'vue'
  import { storeToRefs } from 'pinia'
  import {
    useDeviceStore
  } from "../store/connected.store"
  import type { Device } from '../models/device.models'
  import * as runtime from "../../wailsjs/runtime/runtime.js";

  const emits = defineEmits(["deviceSelected"])
  const  { createOrUpdateItem, flushOldItems,updateStoreTimestamp } = useDeviceStore()
  const { devices } = storeToRefs(useDeviceStore())

  const state: any = reactive({
    loading: false, 
    showDeviceState: ref(true),
    validShowDevices: computed(() => {
      return state.showDeviceState && devices
    }),
  })

  const onImportEvent = async (message: Device[]) => {
    message.forEach((d: Device) => {
      updateStoreTimestamp()
      createOrUpdateItem(
        d.vendor_id, d.product_id, d
      )
      flushOldItems()
      
    })
    if (message.length === 0) {
      updateStoreTimestamp()
      flushOldItems()
    }
  };

  onMounted(()=> {
    runtime.EventsOn("devices", onImportEvent);
  })
    
</script>

<template>
    <v-card
      elevation="2"
      class="device-select"
      v-for="(rates, index) in devices"
      :key="index + 'box'"
      @click="$emit('deviceSelected', rates.checksum)"
      v-if="state.validShowDevices"
    >
      {{rates.vendor_id}}:{{rates.product_id}}
    </v-card>
</template>

<style>
.device-select {
  height: 4em;
  width: 30em;
}
</style>