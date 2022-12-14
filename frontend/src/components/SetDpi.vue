<script lang="ts" setup>
    import { defineProps, reactive, watch, defineEmits } from 'vue';
    import { SetDpiWrapper } from "../../wailsjs/go/backend/App"

    const props: any = defineProps({
        dpi: Number,
    })

    const state = reactive({
        currentdpi:  props.dpi,
    })

    const emit = defineEmits(['new_dpi'])

    watch(()=> state.currentdpi, (newVal :number , oldVal :number) => {
        console.log(`Change to ${newVal} from ${oldVal}`)
        // perform front end validation here
        // even though we already do this on the backend.
        // JS purely uses int64 without uint64.
        // we get this emitted as a number but we need to make sure
        // its unsigned
        // pretty much we just clamp
        if (newVal < 200) {
            newVal = 200
        }
        if (newVal > 16000){
            newVal = 16000
        }
        SetDpiWrapper(newVal)
    })

    const onClose = () => {
        emit('new_dpi', state.currentdpi)
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
            <v-slider
                v-model="state.currentdpi"
                :min="200"
                :max="16000"
                :step="100"
                thumb-label
            />
        </v-row>
    </v-container>
</template>

<style>
</style>