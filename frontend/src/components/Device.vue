<script lang="ts">
export default {

}
</script>
<script lang="ts" setup>

  import { ref, onMounted, reactive, computed } from 'vue'
  import { storeToRefs } from 'pinia'
  import {
    useDeviceStore
  } from "../store/connected.store"
  import type { Device } from '../models/device.models'
  import * as runtime from "../../wailsjs/runtime/runtime.js";

  const  { createOrUpdateItem, flushOldItems,updateStoreTimestamp } = useDeviceStore()
  const { devices } = storeToRefs(useDeviceStore())
  const state: any = reactive({
    loading: false, 
    showDeviceState: ref(true),
    validShowDevices: computed(() => {
      return state.showDeviceState && devices
    }),
  })
  const activeName = ref('1')
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

  const showDevice = () => {
    console.log(state.showDeviceState)
    state.showDeviceState = !state.showDeviceState
  }
    
</script>

<template>
  <TransitionGroup name="fade" mode="out-in">
      <v-card
        elevation="2"
        class="device-select"
        v-for="(rates, index) in devices"
        :key="index + 'box'"
        @click="showDevice"
        v-if="state.validShowDevices"
      >
        {{rates.vendor_id}}:{{rates.product_id}}
      </v-card>
  </TransitionGroup>
</template>

<style>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.5s ease-out;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
.device-select {
  height: 4em;
  width: 30em;
}
</style>