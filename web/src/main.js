import { createApp } from 'vue'
import App from './App.vue'

const app = createApp(App)

// 整合element-plus
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
// 全局引入字体图标
import * as ElIcons from '@element-plus/icons-vue'
for (const name in ElIcons) {
    app.component(name, (ElIcons)[name])
}
app.use(ElementPlus)


app.mount('#app')