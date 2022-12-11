
<script lang="ts" setup>
    import { defineProps, reactive, toRefs, onMounted } from 'vue';
    import { SelectDevice } from "../../wailsjs/go/backend/App"
    const props: any = defineProps({
        checksum: String,
    })


    onMounted(() => {
        let err = SelectDevice(props.checksum);
        console.log(err)
    })
    
</script>

<template>
    <v-card
      elevation="2"
      class="device-select"
    >
        {{ props.checksum }}
    </v-card>
</template>
