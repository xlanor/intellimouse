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
  <div class="app-root">
  <transition name="fade" mode="out-in">
    <div class="display-card" key="1" v-if="state.displaySelectDevices">
      <Device @deviceSelected="deviceSelected"/>
    </div>
    <div key="2" v-else>
      <DeviceSettings :checksum="state.selectedChecksum"/>
    </div>
  </transition>
    
  </div>
</template>

<style>
@import url('https://fonts.googleapis.com/css?family=Raleway');
.app-root{
  height: 100%;
}

.display-card{
  height: 100%;
}

#app {
  font-family: "Raleway";
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
