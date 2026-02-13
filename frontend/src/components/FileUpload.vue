<script setup>
import { ref } from 'vue';

defineProps({
    studentFile: {
        type: Object,
        default: null
    }
});

const emit = defineEmits(['file-change']);

const isDragging = ref(false);
const showDownloadHelp = ref(false);
const fileInputRef = ref(null);

const handleDrop = (event) => {
    isDragging.value = false;
    const file = event.dataTransfer.files[0];
    if (file) emit('file-change', file);
};

const handleFileChange = (event) => {
    const file = event.target.files[0];
    emit('file-change', file);
};

const triggerFileInput = () => {
    fileInputRef.value.click();
};
</script>

<template>
    <div class="max-w-[612px] mx-auto mb-10 p-8 border-2 rounded-2xl transition-all duration-300 animate-entry delay-100 text-center relative"
        :class="[
            isDragging ? 'border-emerald-500 bg-emerald-50/60 scale-[1.01] shadow-xl border-dashed' : 'border-dashed border-stone-200 bg-stone-50/50 hover:bg-stone-100/50'
        ]" @dragover.prevent="!studentFile ? isDragging = true : null"
        @dragleave.prevent="!studentFile ? isDragging = false : null"
        @drop.prevent="!studentFile ? handleDrop($event) : null">

        <h2 class="text-2xl font-bold text-stone-700 mb-6 flex flex-col items-center justify-center gap-2 font-serif">
            <div class="flex items-center">📂 上傳全人資料</div>
            <span @click="showDownloadHelp = !showDownloadHelp"
                class="text-sm text-stone-500 underline cursor-pointer hover:text-stone-700 transition-colors select-none font-sans font-medium relative z-10">
                如何取得全人資料 JSON 檔案?
            </span>
        </h2>

        <div v-if="showDownloadHelp"
            class="mb-6 p-5 bg-white/80 backdrop-blur-sm rounded-xl shadow-sm text-stone-600 leading-relaxed border border-stone-200 text-left max-w-2xl mx-auto">
            <p class="mb-1"><span class="font-bold">Step 1️⃣：</span>進入政大首頁並且登入 iNCCU</p>
            <p class="mb-1"><span class="font-bold">Step 2️⃣：</span>點選「<a
                    href="https://i.nccu.edu.tw/sso_app/PersonalInfoSSO.aspx" target="_blank"
                    class="text-emerald-600 hover:text-emerald-700 underline decoration-dotted transition-colors">進入我的全人</a>（點擊文字可直接前往全人系統）」
            </p>
            <p class="mb-1"><span class="font-bold">Step 3️⃣：</span>下滑到底，在「相關連結」找到「資料格式化匯出」選項，進入後選擇「課業學習」後下載</p>
            <p><span class="font-bold">Step 4️⃣：</span>得到熱騰騰的🔥全人資料 JSON 檔案🔥！</p>
        </div>

        <div v-if="!studentFile" @click="triggerFileInput"
            class="w-full max-w-md mx-auto min-h-[140px] border-2 border-dashed border-stone-300 rounded-xl flex flex-col items-center justify-center cursor-pointer hover:border-emerald-500 hover:bg-white transition-all group relative z-10 bg-white/40">
            <div
                class="w-14 h-14 bg-stone-100 rounded-full flex items-center justify-center mb-3 group-hover:scale-110 transition-transform group-hover:bg-emerald-100 duration-300">
                <svg xmlns="http://www.w3.org/2000/svg"
                    class="h-7 w-7 text-stone-400 group-hover:text-emerald-600 transition-colors" fill="none"
                    viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
                </svg>
            </div>
            <p class="text-stone-600 font-bold text-lg group-hover:text-emerald-700 transition-colors">點擊選擇或拖曳檔案至此
            </p>
            <p class="text-xs text-stone-400 mt-1 font-mono">支援格式：.json</p>
        </div>

        <div v-else
            class="w-full max-w-md mx-auto min-h-[140px] bg-white border-2 border-emerald-500/30 rounded-xl flex flex-col items-center justify-center p-6 shadow-sm">
            <h3 class="text-2xl font-bold text-stone-800 font-serif mb-2 tracking-wide">{{ studentFile.name }}</h3>
            <div class="flex items-center gap-3 mb-6">
                <span class="px-3 py-1 bg-stone-100 text-stone-500 rounded-full text-xs font-mono font-bold">JSON</span>
                <span class="text-stone-400 font-mono text-sm">{{ (studentFile.size / 1024).toFixed(2) }} KB</span>
            </div>

            <button @click="triggerFileInput"
                class="px-8 py-2.5 bg-white border border-stone-200 text-stone-500 font-bold rounded-xl hover:bg-emerald-50 hover:text-emerald-600 hover:border-emerald-200 transition-all shadow-sm hover:shadow-md flex items-center gap-2 group">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 group-hover:scale-110 transition-transform"
                    fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                </svg>
                <span>重新選擇</span>
            </button>
        </div>
        <input type="file" ref="fileInputRef" id="jsonFile" accept=".json" @change="handleFileChange"
            @click="$event.target.value = null" class="hidden">
    </div>
</template>