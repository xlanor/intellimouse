
<script lang="ts" setup>
    import { defineProps, reactive, computed, onMounted } from 'vue';
    import type { Mouse } from "../models/mouse.models";
    import type { ButtonEvent } from "../models/mouse_event.models"
    import { button_map } from "../helpers/map_back_button"
    import { SelectDevice } from "../../wailsjs/go/backend/App"
    import SetLedColor from "./SetLedColor.vue"
    import SetDpi from "./SetDpi.vue"
    import SetMouseButton from './SetMouseButton.vue';
    import * as runtime from "../../wailsjs/runtime/runtime.js";
    const props: any = defineProps({
        checksum: String,
    })

    const onMouseEvent = async(loadedMouse: Mouse) => {
        console.log(`GOT MOUSE ${JSON.stringify(loadedMouse)}`)
        state.mouse_loaded = true;
        state.dpi = loadedMouse.dpi;
        state.back_button = loadedMouse.back_button;
        state.led = loadedMouse.led;
    }

    const state = reactive({
        mouse_loaded: false,
        dpi: 0,
        back_button: "",
        forward_button: "",
        middle_button: "",
        led: "",
        show_information_page: true,
        show_select_led: false,
        show_select_dpi: false,
        show_select_back_button: false,
        show_select_forward_button: false,
        show_select_middle_button: false,
        validShowDevices: computed(():any => {
            let rs: any = {  'width':'36px','min-height': '36px', 'border-style': 'solid' }
            if (state.led !== "") {
                return { ...rs, 'background-color':state.led}
            }
            return rs
        }),
        showButtonPage: computed((): any => {
            if (!state.show_information_page && (state.show_select_back_button || state.show_select_forward_button || state.show_select_middle_button)) {
                return true;
            }
            return false;
        }),
        showButtonKey :computed((): any => {
            if (state.show_information_page) {
                return ""
            }
            if (state.show_select_back_button) {
                return "back"
            }
            if (state.show_select_forward_button) {
                return "front"
            }
            if (state.show_select_middle_button) {
                return "middle"
            }
        }),
        showButtonValue: computed((): any => {

            if (state.show_information_page) {
                return ""
            }
            if (state.show_select_back_button) {
                return state.back_button
            }
            if (state.show_select_forward_button) {
                return state.forward_button
            }
            if (state.show_select_middle_button) {
                return state.middle_button
            }
        })

    })

    const toggle_show_select_led = () => {
        state.show_information_page = !state.show_information_page;
        state.show_select_led = !state.show_select_led;
    }

    const toggle_show_select_dpi = () => {
        state.show_information_page = !state.show_information_page;
        state.show_select_dpi = !state.show_select_dpi;
    }

    const toggle_show_back_button = () => {
        state.show_information_page = !state.show_information_page;
        state.show_select_back_button = !state.show_select_back_button;

    }

    const toggle_show_forward_button = () => {
        state.show_information_page = !state.show_information_page;
        state.show_select_forward_button = !state.show_select_forward_button;
    }

    const toggle_show_middle_button = () => {
        state.show_information_page = !state.show_information_page;
        state.show_select_middle_button = !state.show_select_middle_button;
    }

    const updateDpi = (new_dpi: number) => {
        state.dpi = new_dpi
        toggle_show_select_dpi()
    }

    const updateLed = (new_led: string) => {
        state.led = new_led
        toggle_show_select_led()
    }
    const showButtonBinding = (new_button: ButtonEvent) => {
        if (new_button.button_type === "front") {
            state.forward_button = new_button.new_value
            toggle_show_forward_button()
        } else if (new_button.button_type === "middle") {
            state.middle_button = new_button.new_value
            toggle_show_middle_button()
        } else {
            state.back_button = new_button.new_value
            toggle_show_back_button()
        }
    }
    // use callbacks here, because we want to tightly couple these things together.


    onMounted(() => {
        SelectDevice(props.checksum);
        /*
        state.mouse_loaded = true;
        state.dpi = 16000;
        state.back_button = "INTELLIMOUSE_PRO_BACK_BUTTON_SET_BACK_BUTTON";
        state.led = "#D40B58";
        */
        runtime.EventsOn("mouseinformation", onMouseEvent);
    })
</script>

<template>
    <v-row fill-height align="center" justify="center" >
    <transition name="fade" mode="out-in">
        <v-col key="1" v-if="state.show_information_page" cols="5" class="device-select">
            <v-row
                align="start"
                no-gutters
            >
                <v-col cols="7" align="right" class="column-pad-right">
                    DPI:
                </v-col>
                <v-col align="left" cols="5">
                    <div class="select-box" @click="toggle_show_select_dpi">
                        {{state.dpi}}
                    </div>
                </v-col>
            </v-row>
            <v-row
                align="start"
                no-gutters
            >
                <v-col cols="7" align="right" class="column-pad-right">
                    Back Button:
                </v-col>
                <v-col align="left" cols="5">
                    <div class="select-box" @click="toggle_show_back_button">
                    {{button_map.get(state.back_button)}}
                    </div>
                </v-col>
            </v-row>
            <v-row
                align="start"
                no-gutters
            >
                <v-col cols="7" align="right" class="column-pad-right">
                    LED Lighting:
                </v-col>
                <v-col align="left" class="led-box">
                    <div style="height:36px">
                        <div :style="state.validShowDevices"
                            @click="toggle_show_select_led"
                        />
                    </div>
                </v-col>
            </v-row>
        </v-col>
        <v-col key="2" v-else-if="state.show_select_dpi" class="device-select">
            <SetDpi :key="state.dpi" :dpi="state.dpi" @new_dpi="updateDpi"/>
        </v-col>
        <v-col key="3" v-else-if="state.show_select_led" cols="5" class="device-select" > 
            <SetLedColor :key="state.led" :hashcolor="state.led" @new_led="updateLed"/>
        </v-col>
        <v-col key="4" v-else="state.showButtonPage" cols="5" class="device-select">
            <SetMouseButton :key="state.showButtonKey" :button_type="state.showButtonKey" @new_binding="showButtonBinding"/>
        </v-col>
    </transition>
    </v-row>
</template>

<style>
.container-box {
  height: 100%;
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}
.column-pad-right{
  padding-right: 20px !important;
}
.device-select{
  margin-top: 10px;
  background-color: rgba(255, 255, 255, 0.2);
  backdrop-filter: blur(20px) saturate(160%) contrast(45%) brightness(140%);
  -webkit-backdrop-filter: blur(20px) saturate(160%) contrast(45%) brightness(140%);
  border-radius: 15px;
  transition: 0.5s;
  border: 1px solid #b5b5b5;
  border-radius: 10px;
    font-size: 1.5em; 
}
.select-box:hover{
    cursor: grab;
}
.led-box:hover{
    cursor: grab;
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
