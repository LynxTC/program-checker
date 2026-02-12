package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// preprocessRequirements 階段 1: 處理學程要求的預處理 (名稱解析、特殊學程的課程清單調整)
func preprocessRequirements(programID string, program Program) ([]ProgramRequirement, map[string]bool, map[string]bool, map[string]string) {
	// 建立 Requirements 的副本
	localRequirements := make([]ProgramRequirement, len(program.Requirements))
	for i, req := range program.Requirements {
		localRequirements[i] = req
		coursesCopy := make([]string, len(req.Courses))
		copy(coursesCopy, req.Courses)
		localRequirements[i].Courses = coursesCopy
	}

	programCourseNamesClean := make(map[string]bool)
	geCourseNames := make(map[string]bool)
	courseInstructorMap := make(map[string]string)

	isSoutheastAsianProgram := programID == "southeast_asian_area_studies"

	for i, req := range localRequirements {
		for j, courseName := range req.Courses {
			// 特殊處理：東南亞區域研究微學程 - 處理 "課程名稱(教師名)" 格式
			if isSoutheastAsianProgram && strings.Contains(courseName, "(") && strings.HasSuffix(courseName, ")") {
				start := strings.LastIndex(courseName, "(")
				realName := strings.TrimSpace(courseName[:start])
				instructor := courseName[start+1 : len(courseName)-1]

				// 更新副本中的課程名稱為真實名稱
				localRequirements[i].Courses[j] = realName

				// 記錄該課程對應的教師
				norm := normalizeCourseName(realName)
				courseInstructorMap[norm] = instructor
				programCourseNamesClean[norm] = true
			} else {
				programCourseNamesClean[normalizeCourseName(courseName)] = true
			}
		}
	}

	// 加入通識課程定義
	for _, courseName := range program.GeneralEducationCourses {
		norm := normalizeCourseName(courseName)
		programCourseNamesClean[norm] = true
		geCourseNames[norm] = true
	}

	// 特殊處理：近代社會的身體與性別跨領域學分學程 - 加入通識課程
	if programID == "modern_society_body_gender" {
		geCourses := []string{
			"自我、身體、文化",
			"同志生命美學",
			"臺灣電影與文學中的性別",
			"身心障礙與臺灣藝文",
			"藝術、自我探索與文化溯源",
		}
		for _, courseName := range geCourses {
			norm := normalizeCourseName(courseName)
			programCourseNamesClean[norm] = true
			geCourseNames[norm] = true
		}
	}

	return localRequirements, programCourseNamesClean, geCourseNames, courseInstructorMap
}

// filterAndProcessCourses 階段 2: 篩選並處理課程 (含特殊學程的學分上限、群組調整等)
func filterAndProcessCourses(programID string, rawCourses []StudentCourse, localRequirements *[]ProgramRequirement, programCourseNamesClean map[string]bool, geCourseNames map[string]bool, courseInstructorMap map[string]string) ([]StudentCourse, []StudentCourse) {
	var relevantPassed []StudentCourse
	var inProgressCourses []StudentCourse

	// 基礎篩選
	for _, course := range rawCourses {
		if _, ok := programCourseNamesClean[course.Name]; ok {
			if course.IsPassed {
				relevantPassed = append(relevantPassed, course)
			} else if course.IsInProgress {
				inProgressCourses = append(inProgressCourses, course)
			}
		}
	}

	// 特殊處理：東南亞區域研究微學程 - 同一名老師開設課程至多認列兩門
	if programID == "southeast_asian_area_studies" {
		instructorCounts := make(map[string]int)
		var filteredByInstructor []StudentCourse

		// 先對已通過課程排序 (學分高者優先)
		sort.Slice(relevantPassed, func(i, j int) bool {
			return relevantPassed[i].Credit > relevantPassed[j].Credit
		})

		for _, c := range relevantPassed {
			norm := normalizeCourseName(c.Name)
			if instructor, ok := courseInstructorMap[norm]; ok {
				if instructorCounts[instructor] < 2 {
					instructorCounts[instructor]++
					filteredByInstructor = append(filteredByInstructor, c)
				}
			} else {
				filteredByInstructor = append(filteredByInstructor, c)
			}
		}
		relevantPassed = filteredByInstructor
	}

	// 特殊處理：CFA核心學程 - 特定先修課程即使無成績也視為已修
	if programID == "CFA" {
		specialCourses := map[string]bool{
			"中級會計學（二）": true,
			"投資學":      true,
			"商事法":      true,
			"民法概要":     true,
		}

		var newInProgress []StudentCourse
		for _, c := range inProgressCourses {
			if specialCourses[c.Name] {
				c.IsPassed = true
				relevantPassed = append(relevantPassed, c)
			} else {
				newInProgress = append(newInProgress, c)
			}
		}
		inProgressCourses = newInProgress
	}

	// 特殊處理：專利學分學程 (patent)
	if programID == "patent" {
		// 1. 處理學分上限 (Capping)
		type CapGroup struct {
			Courses []string
			Limit   float64
		}
		capGroups := []CapGroup{
			{[]string{"微積分"}, 2.0},
			{[]string{"民法概要", "民法總則", "民法債編總論（一）", "民法債編總論（二）"}, 6.0},
			{[]string{"普通物理學實驗", "普通物理學實驗（一）", "普通物理學實驗（二）"}, 2.0},
		}

		for _, group := range capGroups {
			var groupCourses []*StudentCourse
			for i := range relevantPassed {
				c := &relevantPassed[i]
				for _, target := range group.Courses {
					if c.Name == target {
						groupCourses = append(groupCourses, c)
						break
					}
				}
			}

			currentTotal := 0.0
			for _, c := range groupCourses {
				if currentTotal >= group.Limit {
					c.Credit = 0
					c.IsCapped = true
				} else if currentTotal+c.Credit > group.Limit {
					allowed := group.Limit - currentTotal
					c.Credit = allowed
					c.IsCapped = true
					currentTotal += allowed
				} else {
					currentTotal += c.Credit
				}
			}
		}

		// 2. 處理 "民法概要" 的歸屬 (法學院 vs 商學院)
		hasCivilLawOverview := false
		for _, c := range relevantPassed {
			if c.Name == "民法概要" && c.Credit > 0 {
				hasCivilLawOverview = true
				break
			}
		}

		if hasCivilLawOverview {
			var businessReqIndex, lawReqIndex int = -1, -1
			for i, req := range *localRequirements {
				switch req.Category {
				case "商學院":
					businessReqIndex = i
				case "法學院":
					lawReqIndex = i
				}
			}

			if businessReqIndex != -1 && lawReqIndex != -1 {
				businessHasOther := false
				for _, c := range relevantPassed {
					if c.Name == "民法概要" {
						continue
					}
					for _, target := range (*localRequirements)[businessReqIndex].Courses {
						if c.Name == target {
							businessHasOther = true
							break
						}
					}
					if businessHasOther {
						break
					}
				}

				targetToRemoveIndex := businessReqIndex
				if !businessHasOther {
					targetToRemoveIndex = lawReqIndex
				}

				newCourses := []string{}
				for _, c := range (*localRequirements)[targetToRemoveIndex].Courses {
					if c != "民法概要" {
						newCourses = append(newCourses, c)
					}
				}
				(*localRequirements)[targetToRemoveIndex].Courses = newCourses
			}
		}
	}

	// 特殊處理：金融科技專長學程 (fintech)
	if programID == "fintech" {
		// Rule 2: 計算機概論/計算機程式設計/計算機程式 三門課僅能擇一門認列
		compIntroGroup := []string{"計算機概論", "計算機程式設計", "計算機程式"}
		var bestCompIntro StudentCourse
		foundCompIntro := false
		var otherCourses []StudentCourse

		for _, c := range relevantPassed {
			isCompIntro := false
			for _, name := range compIntroGroup {
				if c.Name == name {
					isCompIntro = true
					break
				}
			}
			if isCompIntro {
				if !foundCompIntro || c.Credit > bestCompIntro.Credit {
					bestCompIntro = c
					foundCompIntro = true
				}
			} else {
				otherCourses = append(otherCourses, c)
			}
		}
		relevantPassed = otherCourses
		if foundCompIntro {
			relevantPassed = append(relevantPassed, bestCompIntro)
		}

		// Rule 3: 處理 A/B/C 群組重疊課程的歸屬
		overlapNames := map[string]bool{
			"機器學習與人工智慧個案實作":       true,
			"商業資料分析基礎：Python （一）": true,
			"商業資料分析：Python（1）":    true,
			"程式設計與統計軟體(實務)":       true,
			"用Python學財務計量":        true,
		}

		idxA, idxB, idxC := -1, -1, -1
		for i, req := range *localRequirements {
			if strings.Contains(req.Category, "群A") {
				idxA = i
			}
			if strings.Contains(req.Category, "群B") {
				idxB = i
			}
			if strings.Contains(req.Category, "選修C") {
				idxC = i
			}
		}

		if idxA != -1 && idxB != -1 && idxC != -1 {
			var passedPureA, passedPureB int
			var passedOverlap []string

			isInList := func(name string, list []string) bool {
				for _, v := range list {
					if v == name {
						return true
					}
				}
				return false
			}

			for _, c := range relevantPassed {
				if overlapNames[c.Name] {
					passedOverlap = append(passedOverlap, c.Name)
				} else {
					if isInList(c.Name, (*localRequirements)[idxA].Courses) {
						passedPureA++
					}
					if isInList(c.Name, (*localRequirements)[idxB].Courses) {
						passedPureB++
					}
				}
			}

			assignedToA := []string{}
			assignedToC := []string{}
			currentA := passedPureA
			currentB := passedPureB

			for _, name := range passedOverlap {
				if currentA < 1 {
					assignedToA = append(assignedToA, name)
					currentA++
				} else if (currentA + currentB) < 3 {
					assignedToA = append(assignedToA, name)
					currentA++
				} else {
					assignedToC = append(assignedToC, name)
				}
			}

			// 更新 localRequirements
			updateReqCourses := func(idx int, assigned []string) {
				newCourses := []string{}
				for _, name := range (*localRequirements)[idx].Courses {
					if !overlapNames[name] {
						newCourses = append(newCourses, name)
					}
				}
				newCourses = append(newCourses, assigned...)
				(*localRequirements)[idx].Courses = newCourses
			}

			updateReqCourses(idxA, assignedToA)
			updateReqCourses(idxC, assignedToC)
		}
	}

	// 處理通識課程：若超過一門，僅保留學分最高者
	var gePassed []StudentCourse
	var otherPassed []StudentCourse
	for _, c := range relevantPassed {
		if geCourseNames[c.Name] {
			gePassed = append(gePassed, c)
		} else {
			otherPassed = append(otherPassed, c)
		}
	}
	if len(gePassed) > 1 {
		sort.Slice(gePassed, func(i, j int) bool {
			return gePassed[i].Credit > gePassed[j].Credit
		})
		gePassed = gePassed[:1]
	}
	completedCourses := append(otherPassed, gePassed...)

	return completedCourses, inProgressCourses
}

// processManagementAccounting 特殊處理：管理會計專業學程的計算邏輯
func processManagementAccounting(program Program, completedCourses []StudentCourse) ([]CategoryResult, bool, float64) {
	econCredits := 0.0
	totalPassedCredits := 0.0
	for _, c := range completedCourses {
		totalPassedCredits += c.Credit
		if c.Name == "經濟學" {
			econCredits += c.Credit
		}
	}

	if econCredits < 6.0 {
		totalPassedCredits -= econCredits
		var newCompleted []StudentCourse
		for _, c := range completedCourses {
			if c.Name != "經濟學" {
				newCompleted = append(newCompleted, c)
			}
		}
		completedCourses = newCompleted
	}

	isMet := totalPassedCredits >= program.MinCredits

	uniquePassedCourseNames := make(map[string]bool)
	for _, c := range completedCourses {
		uniquePassedCourseNames[c.Name] = true
	}
	passedCount := len(uniquePassedCourseNames)

	results := []CategoryResult{{
		Category:        program.Requirements[0].Category,
		RequiredCount:   0,
		RequiredCredits: program.MinCredits,
		PassedCount:     passedCount,
		PassedCredits:   totalPassedCredits,
		IsMet:           isMet,
		PassedCourses:   completedCourses,
	}}

	return results, isMet, totalPassedCredits
}

// postprocessResults 階段 3: 處理計算後的特殊規則 (跨群檢核、平均成績、系所限制等)
func postprocessResults(programID string, program Program, studentMajor string, categoryResults []CategoryResult, effectiveTotalCredits float64) ([]CategoryResult, bool, string, bool, string, bool, string, float64) {
	allCategoriesMet := true
	for _, res := range categoryResults {
		if !res.IsMet {
			allCategoriesMet = false
			break
		}
	}

	// 特殊處理：金融科技專長學程 (fintech) - Rule 1: A+B >= 3
	if programID == "fintech" {
		countA, countB := 0, 0
		foundA, foundB := false, false
		for _, res := range categoryResults {
			if strings.Contains(res.Category, "群A") {
				countA = res.PassedCount
				foundA = true
			}
			if strings.Contains(res.Category, "群B") {
				countB = res.PassedCount
				foundB = true
			}
		}
		if foundA && foundB {
			total := countA + countB
			isMet := total >= 3
			msg := ""
			if !isMet {
				allCategoriesMet = false
				msg = "群A與群B合計須至少修習 3 門"
			}
			newResult := CategoryResult{
				Category:        "群A + 群B 總修習門數",
				RequiredCount:   3,
				PassedCount:     total,
				PassedCredits:   0,
				IsMet:           isMet,
				LimitExceeded:   false,
				ExceededMessage: msg,
				PassedCourses:   []StudentCourse{},
			}
			// 插入在 "群B" 之後
			insertIdx := -1
			for i, res := range categoryResults {
				if strings.Contains(res.Category, "群B") {
					insertIdx = i + 1
					break
				}
			}
			if insertIdx != -1 && insertIdx <= len(categoryResults) {
				categoryResults = append(categoryResults[:insertIdx], append([]CategoryResult{newResult}, categoryResults[insertIdx:]...)...)
			} else {
				categoryResults = append(categoryResults, newResult)
			}
		}
	}

	// 特殊處理：跨領域精準健康學分學程
	if program.Name == "跨領域精準健康學分學程" {
		targetGroups := []string{"群A", "群B", "群C", "群D"}
		metGroups := 0
		for _, group := range targetGroups {
			for _, res := range categoryResults {
				if strings.Contains(res.Category, group) && res.PassedCount > 0 {
					metGroups++
					break
				}
			}
		}

		isMet := metGroups >= 2
		msg := ""
		if !isMet {
			allCategoriesMet = false
			msg = "須於群A至群D中至少修習兩群課程"
		}

		newResult := CategoryResult{
			Category:        "跨群選修要求 (A-D群至少兩群)",
			RequiredCount:   2,
			PassedCount:     metGroups,
			PassedCredits:   0,
			IsMet:           isMet,
			LimitExceeded:   false,
			ExceededMessage: msg,
			PassedCourses:   []StudentCourse{},
		}

		insertIdx := -1
		for i, res := range categoryResults {
			if strings.Contains(res.Category, "群D") {
				insertIdx = i + 1
			}
		}

		if insertIdx != -1 && insertIdx <= len(categoryResults) {
			categoryResults = append(categoryResults[:insertIdx], append([]CategoryResult{newResult}, categoryResults[insertIdx:]...)...)
		} else {
			categoryResults = append(categoryResults, newResult)
		}
	}

	// 特殊處理：東南亞文化與宗教跨領域學分學程
	if programID == "southeast_asia_culture_religion_interdisciplinary" {
		for i := range categoryResults {
			if strings.Contains(categoryResults[i].Category, "語言領域") {
				passed := categoryResults[i].PassedCourses

				hasViet := false
				hasIndo := false
				hasThai := false

				checkPair := func(name string) bool {
					found1 := false
					found2 := false
					for _, c := range passed {
						if c.Name == name {
							if c.Semester == "1" {
								found1 = true
							}
							if c.Semester == "2" {
								found2 = true
							}
						}
					}
					return found1 && found2
				}

				if checkPair("初級越語") {
					hasViet = true
				}
				if checkPair("初級印尼語") {
					hasIndo = true
				}
				if checkPair("初級泰語") {
					hasThai = true
				}

				if !hasViet && !hasIndo && !hasThai {
					categoryResults[i].IsMet = false
					allCategoriesMet = false
					categoryResults[i].LimitExceeded = true
					categoryResults[i].ExceededMessage = "須修畢同一語言之第一學期及第二學期課程（如：初級越語 上/下學期）"
				}
				break
			}
		}
	}

	// 特殊處理：人力資源管理學程 (學士及碩士)
	if programID == "human_resource_management_undergraduate" || programID == "human_resource_management_master" {
		targetCats := []string{"程序課程：管理類", "程序課程：勞工關係類", "程序課程：行為類"}
		uniquePassed := make(map[string]bool)
		foundAny := false

		for _, res := range categoryResults {
			for _, target := range targetCats {
				if res.Category == target {
					foundAny = true
					for _, c := range res.PassedCourses {
						uniquePassed[c.Name] = true
					}
				}
			}
		}

		if foundAny {
			count := len(uniquePassed)
			isMet := count >= 2
			msg := ""
			if !isMet {
				allCategoriesMet = false
				msg = "程序課程三類（管理類、勞工關係類、行為類）總共須至少修習 2 門"
			}

			categoryResults = append(categoryResults, CategoryResult{
				Category:        "程序課程總門數檢核",
				RequiredCount:   2,
				PassedCount:     count,
				PassedCredits:   0,
				IsMet:           isMet,
				PassedCourses:   []StudentCourse{},
				LimitExceeded:   false,
				ExceededMessage: msg,
			})
		}
	}

	// 特殊處理：人力資源管理學程 (碩士班)
	if programID == "human_resource_management_master" {
		var reqCatIndex, behCatIndex int = -1, -1
		for i, res := range categoryResults {
			switch res.Category {
			case "必修：管理心理學":
				reqCatIndex = i
			case "程序課程：行為類":
				behCatIndex = i
			}
		}

		if reqCatIndex != -1 {
			hasTarget := false
			var targetCredit float64
			for _, c := range categoryResults[reqCatIndex].PassedCourses {
				if c.Name == "組織行為專題研究" {
					hasTarget = true
					targetCredit = c.Credit
					break
				}
			}

			if hasTarget {
				hasBehavioral := false
				if behCatIndex != -1 && categoryResults[behCatIndex].PassedCount > 0 {
					hasBehavioral = true
				}

				if !hasBehavioral {
					newPassed := []StudentCourse{}
					for _, c := range categoryResults[reqCatIndex].PassedCourses {
						if c.Name != "組織行為專題研究" {
							newPassed = append(newPassed, c)
						}
					}
					categoryResults[reqCatIndex].PassedCourses = newPassed
					categoryResults[reqCatIndex].PassedCount = len(newPassed)
					categoryResults[reqCatIndex].PassedCredits -= targetCredit
					effectiveTotalCredits -= targetCredit

					if categoryResults[reqCatIndex].PassedCount < categoryResults[reqCatIndex].RequiredCount {
						categoryResults[reqCatIndex].IsMet = false
						allCategoriesMet = false
						categoryResults[reqCatIndex].LimitExceeded = true
						categoryResults[reqCatIndex].ExceededMessage = "修習「組織行為專題研究」須另修習至少一門行為類程序課程始得認列"
					}
				}
			}
		}
	}

	// 特殊處理：行銷學程 (學士班及碩士班)
	if programID == "marketing_undergraduate" || programID == "marketing_master" {
		restrictedGroups := [][]string{
			{"公共關係管理", "公共關係概論", "公共關係理論", "公關管理專題－危機溝通"},
			{"服務業行銷", "服務行銷管理"},
			{"多變量分析", "多變量統計分析"},
			{"財務行銷", "財務行銷實務專題"},
			{"品牌行銷專題研究", "專題研究－品牌行銷"},
		}

		for i, res := range categoryResults {
			if res.Category == "選修課程" {
				courseToGroup := make(map[string]int)
				for gIdx, group := range restrictedGroups {
					for _, name := range group {
						courseToGroup[name] = gIdx
					}
				}

				usedGroups := make(map[int]bool)
				var newPassed []StudentCourse
				var newPassedCredits float64

				for _, c := range res.PassedCourses {
					if gIdx, exists := courseToGroup[c.Name]; exists {
						if !usedGroups[gIdx] {
							usedGroups[gIdx] = true
							newPassed = append(newPassed, c)
							newPassedCredits += c.Credit
						} else {
							effectiveTotalCredits -= c.Credit
						}
					} else {
						newPassed = append(newPassed, c)
						newPassedCredits += c.Credit
					}
				}

				categoryResults[i].PassedCourses = newPassed
				categoryResults[i].PassedCount = len(newPassed)
				categoryResults[i].PassedCredits = newPassedCredits
				categoryResults[i].IsMet = categoryResults[i].PassedCount >= categoryResults[i].RequiredCount
				break
			}
		}
	}

	// 平均成績檢核
	avgScoreRequired := false
	avgScoreStr := "0.0"
	avgScoreMet := false
	avgScoreThreshold := ""

	if programID == "marketing_undergraduate" || programID == "marketing_master" || programID == "CFA" || programID == "real_property_financial_management" {
		avgScoreRequired = true
		threshold := 80.0
		if programID == "real_property_financial_management" {
			threshold = 70.0
		}
		avgScoreThreshold = fmt.Sprintf("%g", threshold)

		totalScoreCredit := 0.0
		totalCreditForAvg := 0.0
		uniqueCoursesForAvg := make(map[string]bool)

		for _, res := range categoryResults {
			if strings.Contains(res.Category, "先修") {
				continue
			}
			for _, c := range res.PassedCourses {
				key := c.Name + "-" + c.Semester
				if !uniqueCoursesForAvg[key] {
					uniqueCoursesForAvg[key] = true
					s, err := strconv.ParseFloat(c.Score, 64)
					if err == nil {
						totalScoreCredit += s * c.Credit
						totalCreditForAvg += c.Credit
					}
				}
			}
		}

		avg := 0.0
		if totalCreditForAvg > 0 {
			avg = totalScoreCredit / totalCreditForAvg
		}
		avgScoreStr = fmt.Sprintf("%.2f", avg)
		avgScoreMet = avg >= threshold
	}

	restrictionMessage := ""
	// 特殊處理：外語專長商管學分學程 - 商學院學生不得修習
	if programID == "foreign_language_student_business_primer" && businessMajors[studentMajor] {
		allCategoriesMet = false
		restrictionMessage = "本學程限定非商學院學生修習（商學院學生無法申請）"
	}
	// 特殊處理：管理會計專業學程 - 非商學院學生不得修習
	if programID == "management_accounting" && !businessMajors[studentMajor] {
		allCategoriesMet = false
		restrictionMessage = "本學程限定商學院學生修習（非商學院學生無法申請）"
	}

	return categoryResults, allCategoriesMet, restrictionMessage, avgScoreRequired, avgScoreStr, avgScoreMet, avgScoreThreshold, effectiveTotalCredits
}

// processStandardRequirements 處理一般學程的分類要求計算 (核心迴圈邏輯)
func processStandardRequirements(localRequirements []ProgramRequirement, completedCourses []StudentCourse) ([]CategoryResult, float64) {
	var categoryResults []CategoryResult
	effectiveTotalCredits := 0.0

	for _, req := range localRequirements {
		passedInThisCategory := []StudentCourse{} // 該分類下已通過的課程紀錄

		// 找出該類別已通過課程
		for _, c := range completedCourses {
			// 檢查該課程是否在當前分類的要求列表中
			isMatch := false
			for _, reqCourseName := range req.Courses {
				if c.Name == reqCourseName {
					isMatch = true
					break
				}
			}
			if isMatch {
				passedInThisCategory = append(passedInThisCategory, c)
			}
		}

		// 根據學分由高到低排序，以便在有 MaxCount 限制時優先採計高學分課程
		sort.Slice(passedInThisCategory, func(i, j int) bool {
			return passedInThisCategory[i].Credit > passedInThisCategory[j].Credit
		})

		// 計算不重複的已通過課程門數
		uniquePassedCourseNames := make(map[string]bool)
		passedCreditsInCategory := 0.0
		hasCappedCourse := false
		for _, c := range passedInThisCategory {
			uniquePassedCourseNames[c.Name] = true
			passedCreditsInCategory += c.Credit
			if c.IsCapped {
				hasCappedCourse = true
			}
		}
		passedCount := len(uniquePassedCourseNames)

		// 計算該分類貢獻的有效學分 (處理 MaxCount 和 MaxCredits)
		creditsContributingToTotal := 0.0
		countContributing := 0

		limitExceeded := false
		exceededMsg := ""

		for _, c := range passedInThisCategory {
			// Check MaxCount
			if req.MaxCount > 0 && countContributing >= req.MaxCount {
				limitExceeded = true
				exceededMsg = fmt.Sprintf("超過門數上限 (至多 %d 門)", req.MaxCount)
				break
			}

			// Check MaxCredits
			addedCredit := c.Credit
			if req.MaxCredits > 0 {
				if creditsContributingToTotal+addedCredit > req.MaxCredits {
					addedCredit = req.MaxCredits - creditsContributingToTotal
					limitExceeded = true
					exceededMsg = fmt.Sprintf("超過學分上限 (至多 %.1f 學分)", req.MaxCredits)
				}
			}

			if addedCredit > 0 {
				creditsContributingToTotal += addedCredit
				countContributing++
			} else if req.MaxCredits > 0 && creditsContributingToTotal >= req.MaxCredits {
				limitExceeded = true
				exceededMsg = fmt.Sprintf("超過學分上限 (至多 %.1f 學分)", req.MaxCredits)
				break
			}
		}
		if !strings.HasPrefix(req.Category, "先修課程") {
			effectiveTotalCredits += creditsContributingToTotal
		}

		// 檢查是否有被 Patent 邏輯 Cap 的課程
		if hasCappedCourse && !limitExceeded {
			limitExceeded = true
			exceededMsg = "部分課程因超過群組學分上限而不計分或減修"
		}
		// 檢查 MaxCount (針對總數)
		if req.MaxCount > 0 && passedCount > req.MaxCount && !limitExceeded {
			limitExceeded = true
			exceededMsg = fmt.Sprintf("超過門數上限 (至多 %d 門)", req.MaxCount)
		}
		// 檢查 MaxCredits (針對總學分)
		if req.MaxCredits > 0 && passedCreditsInCategory > req.MaxCredits && !limitExceeded {
			limitExceeded = true
			exceededMsg = fmt.Sprintf("超過學分上限 (至多 %.1f 學分)", req.MaxCredits)
		}

		// 判斷該分類是否符合要求 (MinCount, MinCredits)
		isMet := passedCount >= req.MinCount
		if req.MinCredits > 0 && passedCreditsInCategory < req.MinCredits {
			isMet = false
		}

		// 若先修課程類別課程已修滿該列所有已列課程，將狀態視為已達成
		if strings.HasPrefix(req.Category, "先修課程") && len(req.Courses) > 0 {
			if req.MinCount == 0 && req.MinCredits == 0 {
				req.MinCount = len(req.Courses)
				isMet = passedCount >= len(req.Courses)
			}
		}

		categoryResults = append(categoryResults, CategoryResult{
			Category:        req.Category,
			RequiredCount:   req.MinCount,
			RequiredCredits: req.MinCredits,
			PassedCount:     passedCount,
			PassedCredits:   passedCreditsInCategory,
			IsMet:           isMet,
			PassedCourses:   passedInThisCategory,
			LimitExceeded:   limitExceeded,
			ExceededMessage: exceededMsg,
		})
	}

	return categoryResults, effectiveTotalCredits
}
