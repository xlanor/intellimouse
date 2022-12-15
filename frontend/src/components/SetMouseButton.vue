<script lang="ts" setup>
    import { reactive, defineProps, defineEmits } from "vue"

    import type { ButtonEvent } from "../models/mouse_event.models"
    import { button_map } from "../helpers/map_back_button"

    const props:any = defineProps({
        button_type: String,
        mouse_binding: String,
    })

    const emits = defineEmits(["new_binding"])

    const state = reactive({
        button_type: props.button_type,
        mouse_binding: props.mouse_binding,
    })

    const updateMouseBinding = (newBinding: String) => {
        console.log(newBinding)
    }

    const onClose = () => {
        let be = {} as ButtonEvent 
        be.button_type = state.button_type
        be.new_value = state.mouse_binding

        emits("new_binding", be)
    }
</script>

<template>

    <v-container>
        <v-row class="flex-column" align="center" justify="center"> 
            {{button_type?.toUpperCase()}} Button
            <v-btn @click="onClose">
            </v-btn>
        </v-row>
        <v-row>
         <div class="line"></div>
        </v-row>
        <v-row>
            <v-radio-group
                v-model="state.mouse_binding"
            >
                <v-radio
                    v-for="[key, value] in button_map"
                    @change="updateMouseBinding"
                    :key="key"
                    :label="value"
                    :value="key"
                    :checked="key === state.mouse_binding"
                >
            </v-radio>
            </v-radio-group>
        </v-row>
    </v-container>
</template>

<style>
</style>