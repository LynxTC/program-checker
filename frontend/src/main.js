// 1. 匯入 Vue 框架中的核心函式
import { createApp } from 'vue';

// 2. 匯入主要應用程式的根組件
import App from './App.vue';

// 3. (可選) 匯入全域 CSS/樣式檔案
// 這是 Vite 專案中匯入全域樣式（例如 Tailwind CSS 基礎樣式）的標準做法
import './index.css';
// 如果您將 Tailwind 的指令放在 src/index.css，則需要這行

// 4. 建立應用程式實例
// createApp(App) 會建立一個 Vue 應用程式的實例，以 App.vue 作為根組件。
const app = createApp(App);

// 5. (可選) 註冊全域套件或服務
// 如果您需要使用 Vue Router, Pinia (狀態管理), 或其他第三方套件，
// 則需在這裡使用 app.use() 進行註冊。

// 範例 (如果需要使用路由，通常會在這裡匯入並註冊):
// import router from './router';
// app.use(router);

// 6. 將應用程式掛載到 DOM 元素
// 將 Vue 應用程式掛載到 index.html 中 ID 為 'app' 的元素上。
app.mount('#app');