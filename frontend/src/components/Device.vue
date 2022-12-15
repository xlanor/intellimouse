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

  const getDeviceName = (vendor_id: number, product_id: number):string => {
    if (vendor_id === 1118 && product_id === 2090) {
      return "Microsoft Intellimouse Explorer Pro (2019)"
    }
    return "unknown"
  }

  onMounted(()=> {
    runtime.EventsOn("devices", onImportEvent);
  })
    
</script>

<template>
  <v-row fill-height align="center" justify="center" >
  <v-container>
      <v-row align="center" justify="center" class="title-text"> 
        SELECT DEVICE
      </v-row>
      <v-row class="line"/>
      <v-row 
        align="center" justify="center"
        v-for="(rates, index) in devices"
        :key="index + 'box'"
        @click="$emit('deviceSelected', rates.checksum)"
        v-if="state.validShowDevices"
      >
        <div class="select-device-text">
        {{ getDeviceName (rates.vendor_id, rates.product_id)}}
        </div>
      </v-row>
  </v-container>
  </v-row>
</template>

<style>
.title-text {
  font-size: 3.5em;
    text-shadow: 2px 4px 4px rgba(46,91,173,0.6);
}
.container-box {
  height: 100%;
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}
.box1{
  margin-top: 10px !important;
  height: 5em;
  width: 30em;
  background-color: rgba(255, 255, 255, 0.2);
  backdrop-filter: blur(20px) saturate(160%) contrast(45%) brightness(140%);
  -webkit-backdrop-filter: blur(20px) saturate(160%) contrast(45%) brightness(140%);
  border-radius: 15px;
  transition: 0.5s;
  border: 1px solid #b5b5b5;
  border-radius: 10px;
}
.box1:hover{
  box-shadow: 0 0 3px black;
  margin-top: 0px;
  width: 40em;
  height: 7em;
}
.select-device-text {
  cursor: grab;
  top: 0;
  height: 3em;
  width: 25em;
  display: flex;
  font-size: 1.5em;
  margin-top: 10px !important;
  background-color: rgba(255, 255, 255, 0.2);
  backdrop-filter: blur(20px) saturate(160%) contrast(45%) brightness(140%);
  -webkit-backdrop-filter: blur(20px) saturate(160%) contrast(45%) brightness(140%);

  transition: 0.5s;
  border: 1px solid #b5b5b5;
  border-radius: 10px;
  align-items: center; /** Y-axis align **/
  justify-content: center; /** X-axis align **/
  text-shadow: 2px 4px 4px rgba(18, 23, 31, 0.6);
}
.select-device-text:hover{
  box-shadow: 0 0 3px black;
  margin-top: 0px;
  width: 30em;
  height: 3.5em;
}
.line {
  border: 0;
  background-color: rgba(255, 255, 255, 0.2);
  height: 2px;
  margin-bottom: 10px;
}
</style>