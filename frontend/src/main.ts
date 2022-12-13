import {createApp} from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import 'vuetify/styles'
import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'
import { aliases, mdi } from 'vuetify/iconsets/mdi-svg'
import 'material-design-icons-iconfont/dist/material-design-icons.css'

const vuetify = createVuetify({
    components,
    directives,
    theme: { 
      defaultTheme: 'navy',
      themes: {
        navy: {
          dark: true,
          colors: {
            background: '#050A30',
          }
        }
      },
    },
    icons: {
      defaultSet: 'mdi',
      aliases,
      sets: {
        mdi,
      }
    },
  })
const pinia = createPinia()

createApp(App).use(pinia).use(vuetify).mount('#app')