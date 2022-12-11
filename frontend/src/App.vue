<script lang="ts" setup>
  import { reactive, ref, computed } from 'vue'
  import Device from './components/Device.vue'
  import DeviceSettings from './components/DeviceSettings.vue'

  // state

  const state: any = reactive({
    selectedChecksum: ref(""),
    displaySelectDevices: ref(true),
    showMouseInformation: computed(()=> {
      return state.displaySelectDevices === false
    })
  })

  const deviceSelected = (selectedChecksum: string) => {
    state.selectedChecksum = selectedChecksum
    state.displaySelectDevices = ref(false)
    console.log(`deviceSelected called, state ${JSON.stringify(state)}`)
  }

</script>

<template>
  
  <TransitionGroup name="fade" mode="out-in">
    <div v-if="state.displaySelectDevices">
      <Device @deviceSelected="deviceSelected"/>
    </div>
  </TransitionGroup>
  <TransitionGroup name="fade" mode="out-in">
    <div v-if="state.showMouseInformation">
      <DeviceSettings :checksum="state.selectedChecksum"/>
    </div>
  </TransitionGroup>
</template>

<style>

#app {
  font-family: "Roboto Mono";
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}
.fade-enter-active,
.fade-leave-active {
    transition: opacity 1s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
