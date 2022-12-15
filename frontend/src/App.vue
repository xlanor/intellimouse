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
  }

</script>

<template>
  <v-app color="#2E2252">
    <v-row fill-height>
      <transition name="fade" mode="out-in">
          <Device key="1" v-if="state.displaySelectDevices" @deviceSelected="deviceSelected"/>
          <DeviceSettings v-else  key="2" :checksum="state.selectedChecksum"/>
      </transition>
    </v-row>
  </v-app>
</template>

<style>
@import url('https://fonts.googleapis.com/css?family=Raleway');
@import url(@mdi/font/css/materialdesignicons.css);


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
