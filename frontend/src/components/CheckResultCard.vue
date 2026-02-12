<script setup>
defineProps({
    result: {
        type: Object,
        required: true
    }
});
</script>

<template>
    <div class="border-2 p-5 rounded-xl shadow-md"
        :class="result.isCompleted ? 'bg-emerald-50/80 border-emerald-200 text-emerald-900' : 'bg-rose-50/80 border-rose-200 text-rose-900'">
        <div class="flex flex-col sm:flex-row sm:items-center justify-between mb-6 pb-4 border-b border-black/5">
            <h3 class="text-2xl font-bold font-serif tracking-wide">{{ result.programName }}</h3>
            <span
                class="mt-2 sm:mt-0 px-4 py-1.5 text-sm font-bold rounded-full border shadow-sm inline-block text-center"
                :class="result.isCompleted ? 'bg-emerald-100 text-emerald-700 border-emerald-200' : 'bg-rose-100 text-rose-700 border-rose-200'">
                {{ result.isCompleted ? '✓ 已修畢' : '✗ 未修畢' }}
            </span>
        </div>

        <p class="text-stone-700 mb-6 leading-relaxed">{{ result.programDescription }}</p>

        <!-- 顯示資格限制訊息 -->
        <div v-if="result.restrictionMessage"
            class="mb-6 p-4 bg-rose-50 border-l-4 border-rose-500 text-rose-700 rounded-r-lg shadow-sm">
            <div class="flex items-start">
                <span class="text-2xl mr-3">🚫</span>
                <div>
                    <h4 class="font-bold text-lg mb-1">資格限制</h4>
                    <p>{{ result.restrictionMessage }}</p>
                </div>
            </div>
        </div>

        <div class="mb-6 p-5 bg-white/60 rounded-xl border border-black/5 backdrop-blur-sm">
            <h4 class="text-lg font-bold text-stone-800 mb-3 font-serif">學程總學分要求</h4>
            <div class="flex justify-between text-sm">
                <span class="font-medium text-stone-600">應修總學分:</span>
                <span class="font-mono text-base"
                    :class="result.totalCreditsMet ? 'text-emerald-700 font-bold' : 'text-rose-700'">{{
                        result.minRequiredCredits }} 學分</span>
            </div>
            <div class="flex justify-between text-sm">
                <span class="font-medium text-stone-600">已通過學分:</span>
                <span class="font-mono text-base"
                    :class="result.totalCreditsMet ? 'text-emerald-700 font-bold' : 'text-rose-700'">{{
                        result.totalPassedCredits }} 學分</span>
            </div>
            <p class="mt-3 text-xs font-bold uppercase tracking-wide"
                :class="result.totalCreditsMet ? 'text-emerald-600' : 'text-rose-600'">
                {{ result.totalCreditsMet ? '總學分要求已達成' : '總學分要求尚未達成' }}
            </p>
        </div>

        <h4 class="text-lg font-bold text-stone-800 mb-4 font-serif">分項要求</h4>

        <div v-for="cat in result.categoryResults" :key="cat.category" class="mb-4 p-4 rounded-xl border transition-all"
            :class="((cat.requiredCount > 0 || cat.requiredCredits > 0) ? cat.isMet : result.isCompleted) ? 'border-emerald-200 bg-emerald-50/50 text-emerald-900' : 'border-rose-200 bg-rose-50/50 text-rose-900'">
            <div class="flex flex-col sm:flex-row justify-between sm:items-center text-sm font-bold mb-2">
                <span>{{ cat.category }}</span>
                <div class="text-left sm:text-right mt-1 sm:mt-0 font-mono text-xs sm:text-sm opacity-80">
                    <div v-if="cat.requiredCount > 0">
                        {{ cat.passedCount }} / {{ cat.requiredCount }} {{
                            cat.category.includes('跨群選修要求') ? '群'
                                : '門' }}
                    </div>
                    <div v-if="cat.requiredCredits > 0">
                        {{ cat.passedCredits.toFixed(1) }} / {{ cat.requiredCredits.toFixed(1) }} 學分
                    </div>
                    <div v-if="cat.requiredCount === 0 && cat.requiredCredits === 0">
                        門數/學分無強制要求 (依總學分認定)
                    </div>
                </div>
            </div>
            <div v-if="cat.limitExceeded"
                class="text-xs font-bold text-amber-600 mt-2 flex items-center bg-amber-50 p-1.5 rounded">
                <span class="mr-1">⚠️</span>
                {{ cat.exceededMessage }}
            </div>
            <p v-if="cat.requiredCount > 0 || cat.requiredCredits > 0" class="text-xs mt-1 opacity-70">
                狀態: <span class="font-bold uppercase">{{
                    cat.isMet ? '已達成' : '未達成' }}</span>
            </p>
            <div v-if="cat.category !== '群A + 群B 總修習門數' && cat.category !== '跨群選修要求 (A-D群至少兩群)'"
                class="mt-3 text-xs text-stone-600">
                <p class="font-bold mb-1 opacity-70">已通過課程 ({{ cat.passedCourses.length }}):</p>
                <ul
                    class="list-none space-y-1 max-h-32 overflow-y-auto custom-scrollbar bg-white/60 p-2 rounded border border-black/5">
                    <li v-if="cat.passedCourses.length === 0" class="italic opacity-50">無符合要求的已通過課程</li>
                    <li v-for="c in cat.passedCourses" :key="c.name + c.semester">{{ c.name }} ({{
                        c.credit.toFixed(1) }} 學分<span v-if="c.isCapped" class="text-amber-600 font-bold ml-1"
                            title="此課程因超過上限而被調整學分">*</span>, {{
                                c.score }} 分)</li>
                </ul>
            </div>
        </div>

        <!-- 平均成績檢核區塊 (僅針對特定學程顯示) -->
        <div v-if="result.avgScoreRequired && result.totalCreditsMet"
            class="mb-4 p-4 bg-white/60 rounded-xl border border-black/5 backdrop-blur-sm">
            <h4 class="text-md font-bold text-stone-800 mb-2 font-serif">認列課程平均成績檢核</h4>
            <div class="flex justify-between text-sm">
                <span class="font-medium text-stone-600">認列課程平均成績:</span>
                <span class="font-mono" :class="result.avgScoreMet ? 'text-emerald-700 font-bold' : 'text-rose-700'">{{
                    result.avgScore }} 分</span>
            </div>
            <p class="mt-2 text-xs font-bold" :class="result.avgScoreMet ? 'text-emerald-600' : 'text-rose-600'">
                {{ result.avgScoreMet ? `平均成績已達 ${result.avgScoreThreshold} 分標準` : `平均成績未達
                ${result.avgScoreThreshold} 分標準` }}
            </p>
        </div>

        <div v-if="result.inProgressCourses && result.inProgressCourses.length > 0"
            class="mt-6 p-4 border border-amber-300 bg-amber-50/80 rounded-xl">

            <h4 class="text-lg font-bold text-amber-900 mb-2 flex items-center">
                <span class="mr-2">⏳</span>
                修習中課程 ({{ result.inProgressCourses ? result.inProgressCourses.length : 0 }} 門)
            </h4>

            <p class="text-sm text-amber-800 mb-2">以下課程成績尚未送達，若及格可能影響學程完成狀態：</p>
            <ul
                class="list-disc list-inside ml-2 text-sm text-amber-900 max-h-32 overflow-y-auto custom-scrollbar bg-white/60 p-2 rounded border border-amber-200">
                <li v-for="c in result.inProgressCourses" :key="c.name + c.semester">
                    {{ c.name }} ({{ c.credit.toFixed(1) }} 學分) - {{ c.semester }}
                </li>
            </ul>
        </div>
    </div>
</template>
