<script lang="ts" setup>
import { defineProps, reactive, watch, defineEmits } from 'vue';
import { SetLEDWrapper } from "../../wailsjs/go/backend/App"

    const props: any = defineProps({
        hashcolor: String,
    })

    const state = reactive({
        localhash:  props.hashcolor,

    })

    const emit = defineEmits(['new_led'])

    const swatches = [
        ['#FF0000', '#AA0000', '#550000'],
        ['#FFFF00', '#AAAA00', '#555500'],
        ['#00FF00', '#00AA00', '#005500'],
        ['#00FFFF', '#00AAAA', '#005555'],
        ['#0000FF', '#0000AA', '#000055'],
    ]

    watch(()=> state.localhash, (newVal :string , oldVal :string) => {
        console.log(`Change to ${newVal} from ${oldVal}`)
        SetLEDWrapper(newVal)
    })

    const onClose = () => {
        emit('new_led', state.localhash)
    }

</script>

<template>
    <v-container>
        <v-row class="flex-column" align="center" justify="center"> 
            LED Color 
            <v-btn @click="onClose">
            </v-btn>
        </v-row>
        <v-row>
         <div class="line"></div>
        </v-row>
        <v-row>
            <v-col
                class="d-flex justify-center"
            >
                <v-color-picker  
                    class="ma-2"
                    hide-inputs
                    hide-sliders
                    :swatches="swatches"
                    show-swatches
                    mode="hexa"
                    v-model="state.localhash"
                />
            </v-col>
        </v-row>
    </v-container>
</template>

<style>
.line {
  border: 0;
  background-color: rgba(255, 255, 255, 0.2);
  height: 2px;
  margin-bottom: 10px;
}
</style>