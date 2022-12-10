    
<script lang="ts" setup>
  import { ref, onMounted } from 'vue'
  import { storeToRefs } from 'pinia'
  import {
    useDeviceStore
  } from "../store/connected.store"
  import type { Device } from '../models/device.models'
  import * as runtime from "../../wailsjs/runtime/runtime.js";

  const  { createOrUpdateItem, flushOldItems,updateStoreTimestamp } = useDeviceStore()
  const { devices } = storeToRefs(useDeviceStore())

  const activeName = ref('1')
  const onImportEvent = async (message: Device[]) => {
    message.forEach((d: Device) => {
      updateStoreTimestamp()
      createOrUpdateItem(
        d.vendor_id, d.product_id, d
      )
      flushOldItems()
      
    })
    if (message.length == 0) {
      updateStoreTimestamp()
      flushOldItems()
    }
  };
  onMounted(()=> {
    runtime.EventsOn("devices", onImportEvent);
  })
    
</script>

<template>
    <div class="text-subtitle-2 mt-4 mb-2">Accordion</div>
  
    <v-expansion-panels variant="accordion">
      <v-expansion-panel
        v-if="devices"
        v-for="i in devices"
        :key="`$index`"
        title="Item"
        text="Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat."
      ></v-expansion-panel>
    </v-expansion-panels>
  
  </template>
