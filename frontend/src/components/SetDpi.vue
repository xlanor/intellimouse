<script lang="ts" setup>
    import { defineProps, reactive, watch, defineEmits } from 'vue';

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
        
        // TODO: Replace below with backend binding to wails
        //SetLEDWrapper(newVal)
    })

    const onClose = () => {
        emit('new_dpi', state.currentdpi)
    }

</script>

<template>

</template>

<style>
</style>