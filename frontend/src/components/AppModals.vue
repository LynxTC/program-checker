<script setup>
defineProps({
    showDisclaimer: Boolean,
    disclaimerAction: String,
    showCompletion: Boolean,
    completedPrograms: Array,
    showContact: Boolean,
    showPrivacy: Boolean,
    showTerms: Boolean
});

const emit = defineEmits([
    'update:showDisclaimer',
    'update:showCompletion',
    'update:showContact',
    'update:showPrivacy',
    'update:showTerms',
    'confirmDisclaimer'
]);
</script>

<template>
    <!-- 免責聲明 Modal -->
    <div v-if="showDisclaimer"
        class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50 backdrop-blur-sm p-4">
        <div
            class="bg-white rounded-2xl shadow-2xl max-w-md w-full p-8 transform transition-all scale-100 text-center border border-stone-100">
            <h3 class="text-2xl font-bold text-emerald-900 mb-4 flex items-center justify-center font-serif">
                <span class="text-amber-500 mr-2">⚠️</span> 注意事項
            </h3>
            <p class="text-stone-600 mb-8 leading-relaxed">
                分析師提供結果僅供參考<br>可能因申請年度不同或修習同名非認列課程產生誤差<br><br>
                <span
                    class="font-bold text-emerald-800 bg-emerald-50 px-2 py-1 rounded whitespace-nowrap text-xs sm:text-base">實際修習狀態以學程設置單位認定為準</span>
            </p>
            <div class="flex justify-center space-x-3">
                <button @click="emit('update:showDisclaimer', false)"
                    class="px-6 py-2.5 text-stone-500 hover:bg-stone-100 font-bold rounded-xl transition-colors">
                    取消
                </button>
                <button @click="emit('confirmDisclaimer')"
                    class="px-6 py-2.5 bg-emerald-700 hover:bg-emerald-800 text-white font-bold rounded-xl shadow-lg shadow-emerald-900/20 transition-colors">
                    {{ disclaimerAction === 'check' ? '確定' : '開始分析' }}
                </button>
            </div>
        </div>
    </div>

    <!-- 達標恭喜 Modal -->
    <Transition enter-active-class="transition ease-out duration-300" enter-from-class="opacity-0"
        enter-to-class="opacity-100" leave-active-class="transition ease-in duration-200" leave-from-class="opacity-100"
        leave-to-class="opacity-0">
        <div v-if="showCompletion"
            class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50 backdrop-blur-sm p-4">
            <div
                class="bg-white rounded-2xl shadow-2xl max-w-md w-full p-6 sm:p-8 max-h-[90vh] overflow-y-auto transform transition-all scale-100 text-center border border-stone-100">
                <div class="text-5xl mb-4">🎉</div>
                <h3 class="text-2xl font-bold text-emerald-900 mb-4 font-serif tracking-wide">
                    恭喜！學程修畢達成
                </h3>
                <p class="text-stone-600 mb-6 leading-relaxed">
                    您已達成以下 <span class="font-bold text-emerald-700">{{ completedPrograms.length }}</span> 個學程的修畢要求：
                </p>
                <div
                    class="bg-stone-50 rounded-xl p-4 mb-6 max-h-40 overflow-y-auto custom-scrollbar border border-stone-100">
                    <ul class="space-y-2">
                        <li v-for="name in completedPrograms" :key="name" class="text-emerald-800 font-bold font-serif">
                            {{ name }}
                        </li>
                    </ul>
                </div>
                <div
                    class="mb-6 text-sm text-amber-800 bg-amber-50 p-4 rounded-xl border border-amber-100 text-justify leading-relaxed shadow-sm">
                    <span class="font-bold block mb-1 text-amber-900">⚠️ 學分認列注意事項</span>
                    部分學程對於雙主修、輔系或原系所之學分認列可能有特殊限制（如：不得重複認列）。本系統僅進行初步檢核，<span
                        class="font-bold">實際認列結果請以各學程設置單位審核為準</span>。
                </div>
                <p class="text-sm text-stone-400 mb-8 italic">
                    請記得及時向相關學程設置單位提出證明申請
                </p>
                <button @click="emit('update:showCompletion', false)"
                    class="w-full px-6 py-3 bg-emerald-700 hover:bg-emerald-800 text-white font-bold rounded-xl shadow-lg shadow-emerald-900/20 transition-colors">
                    我知道了
                </button>
            </div>
        </div>
    </Transition>

    <!-- Contact Us Modal -->
    <div v-if="showContact"
        class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50 backdrop-blur-sm p-4"
        @click.self="emit('update:showContact', false)">
        <div
            class="bg-white rounded-2xl shadow-2xl max-w-3xl w-full flex flex-col max-h-[90vh] transform transition-all scale-100 border border-stone-100 overflow-hidden">
            <div class="flex justify-between items-center p-4 border-b flex-shrink-0">
                <h3 class="text-xl font-bold text-emerald-900 font-serif">聯絡我們</h3>
                <button @click="emit('update:showContact', false)" class="text-stone-400 hover:text-stone-600 text-2xl"
                    title="關閉">&times;</button>
            </div>
            <div class="p-0 overflow-y-auto flex-grow bg-stone-50">
                <div
                    class="p-4 bg-amber-50 border-b border-amber-100 text-amber-800 text-xs sm:text-sm flex justify-between items-center">
                    <span>表單請以政大帳號登入。若下方表單無法正常顯示，請點擊右側按鈕：</span>
                    <a href="https://docs.google.com/forms/d/e/1FAIpQLSfy53oPNPgDu_O1zwzWWjbbV4A3rn_6RA8FKwEsx8P9kv6r7A/viewform"
                        target="_blank"
                        class="bg-amber-600 text-white px-3 py-1 rounded-md font-bold hover:bg-amber-700 transition-colors">直接開啟表單
                        ↗</a>
                </div>
                <iframe
                    src="https://docs.google.com/forms/d/e/1FAIpQLSfy53oPNPgDu_O1zwzWWjbbV4A3rn_6RA8FKwEsx8P9kv6r7A/viewform"
                    class="w-full h-[80vh] sm:h-[70vh]" frameborder="0" marginheight="0" marginwidth="0">載入中…</iframe>
            </div>
        </div>
    </div>

    <!-- 隱私權政策 Modal -->
    <div v-if="showPrivacy"
        class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50 backdrop-blur-sm p-4"
        @click.self="emit('update:showPrivacy', false)">
        <div
            class="bg-white rounded-2xl shadow-2xl max-w-2xl w-full flex flex-col max-h-[90vh] transform transition-all scale-100 border border-stone-100">
            <div class="flex justify-between items-center p-6 sm:p-8 border-b flex-shrink-0">
                <h3 class="text-2xl font-bold text-emerald-900 font-serif">隱私權政策</h3>
                <button @click="emit('update:showPrivacy', false)" class="text-stone-400 hover:text-stone-600 text-2xl"
                    title="關閉">&times;</button>
            </div>
            <div class="p-6 sm:p-8 overflow-y-auto">
                <div class="text-stone-600 space-y-4 leading-relaxed text-justify">
                    <p>歡迎您使用「ProAnalyst」（以下簡稱本工具）。本工具由𤫹焈焈麀普羅安納利斯特團隊開發與維護。為了讓您能夠安心使用本工具的各項服務與資訊，特此向您說明本工具的隱私權保護政策，以保障您的權益，請您詳閱下列內容：
                    </p>

                    <h4 class="font-bold text-emerald-800 text-lg mt-4">一、隱私權保護政策的適用範圍</h4>
                    <p>隱私權保護政策內容，包括本工具如何處理在您使用網站服務時收集到的個人識別資料。本隱私權保護政策不適用於本工具以外的相關連結網站，也不適用於非本工具所委託或參與管理的人員。</p>

                    <h4 class="font-bold text-emerald-800 text-lg mt-4">二、個人資料的蒐集、處理及利用方式</h4>
                    <ul class="list-disc pl-5 space-y-2">
                        <li>當您使用本工具進行學程檢核時，我們需要您上傳個人的全人資料 JSON 檔案。</li>
                        <li><strong>資料不落地原則：</strong>您上傳的檔案僅會在伺服器的記憶體中進行暫時性的運算與分析，運算完成後即會將結果回傳給您，伺服器<strong>不會儲存</strong>您的檔案內容、成績資料或任何個人識別資訊。
                        </li>
                        <li>
                            學程推薦分析：在您同意免責聲明後，我們會將您上傳的全人資料 JSON 檔案，用於分析並推薦您可能感興趣的學程。分析完成後，相關資料將立即從伺服器記憶體中刪除，不作儲存。
                        </li>
                        <li>本工具不會將您的個人資料提供、交換、出租或出售給任何其他個人、團體、私人企業或公務機關。</li>
                        <li>若您使用「聯絡我們」功能填寫表單，該資料將透過 Google 表單收集與處理，相關權利義務請參閱 Google 隱私權政策。</li>
                    </ul>

                    <h4 class="font-bold text-emerald-800 text-lg mt-4">三、網站對外的相關連結</h4>
                    <p>本工具的網頁提供其他網站的網路連結，您也可經由本工具所提供的連結，點選進入其他網站。但該連結網站不適用本工具的隱私權保護政策，您必須參考該連結網站中的隱私權保護政策。</p>

                    <h4 class="font-bold text-emerald-800 text-lg mt-4">四、隱私權保護政策之修正</h4>
                    <p>本工具隱私權保護政策將因應需求隨時進行修正，修正後的條款將刊登於網站上。</p>
                </div>
            </div>
        </div>
    </div>

    <!-- 服務條款 Modal -->
    <div v-if="showTerms"
        class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50 backdrop-blur-sm p-4"
        @click.self="emit('update:showTerms', false)">
        <div
            class="bg-white rounded-2xl shadow-2xl max-w-2xl w-full flex flex-col max-h-[90vh] transform transition-all scale-100 border border-stone-100">
            <div class="flex justify-between items-center p-6 sm:p-8 border-b flex-shrink-0">
                <h3 class="text-2xl font-bold text-emerald-900 font-serif">服務條款</h3>
                <button @click="emit('update:showTerms', false)" class="text-stone-400 hover:text-stone-600 text-2xl"
                    title="關閉">&times;</button>
            </div>
            <div class="p-6 sm:p-8 overflow-y-auto">
                <div class="text-stone-600 space-y-4 leading-relaxed text-justify">
                    <p>歡迎使用「ProAnalyst」（以下簡稱本服務）。本服務由𤫹焈焈麀普羅安納利斯特團隊開發與維護。為了保障您的權益，請詳細閱讀本服務條款。</p>

                    <h4 class="font-bold text-emerald-800 text-lg mt-4">一、認知與接受條款</h4>
                    <p>當您開始使用本服務時，即表示您已閱讀、瞭解並同意接受本服務條款之所有內容。若您不同意本服務條款的任何部分，請立即停止使用本服務。</p>

                    <h4 class="font-bold text-emerald-800 text-lg mt-4">二、服務性質與免責聲明</h4>
                    <ul class="list-disc pl-5 space-y-2">
                        <li><strong>非官方聲明：</strong>本服務為個人 side project，與國立政治大學無隸屬或代理關係。</li>
                        <li>本服務旨在協助學生快速檢核學分學程修習進度，<strong>檢核結果僅供參考</strong>。</li>
                        <li>本服務所依據之學程規則與課程資料可能隨學校政策變動，我們盡力確保資料之即時性與正確性，但不保證內容完全無誤。</li>
                        <li><strong>最終修畢資格與學分認定，悉以國立政治大學教務處及各學程設置單位之正式審核結果為準。</strong></li>
                        <li>本工具原始碼採用 MIT 授權條款開源，歡迎自由取用、修改與散布，惟依據授權條款規定，<strong>使用時需保留原始著作權聲明與姓名標示
                                (Attribution)</strong>。
                        </li>
                        <li>對於因使用本服務或無法使用本服務而產生之任何直接、間接、附帶、特別、衍生性或懲罰性損害，開發團隊不負任何賠償責任。</li>
                    </ul>

                    <h4 class="font-bold text-emerald-800 text-lg mt-4">三、使用者的守法義務及承諾</h4>
                    <p>您承諾絕不為任何非法目的或以任何非法方式使用本服務，並承諾遵守中華民國相關法規及一切使用網際網路之國際慣例。</p>

                    <h4 class="font-bold text-emerald-800 text-lg mt-4">四、服務變更與終止</h4>
                    <p>我們保留隨時修改、暫停或終止本服務之權利，恕不另行通知。</p>

                    <h4 class="font-bold text-emerald-800 text-lg mt-4">五、生成式 AI 使用聲明</h4>
                    <p>本服務在開發過程中使用了生成式 AI 技術輔助程式碼撰寫與除錯。儘管開發團隊已盡力審核，但仍可能存在未預期的錯誤。</p>
                </div>
            </div>
        </div>
    </div>
</template>
