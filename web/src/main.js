import { createApp } from 'vue'
// 引入封装的router
import App from './App.vue'
import { store } from '@/pinia'


// 整合element-plus
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
// 全局引入字体图标
import * as ElIcons from '@element-plus/icons-vue'

const app = createApp(App)
for (const name in ElIcons) {
  app.component(name, (ElIcons)[name])
}
app
  .use(ElementPlus)
  .use(store)
  .mount('#app')