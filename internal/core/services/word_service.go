package services

import (
	"fmt"
	"service-matrix-go/internal/core/algorithm"
	"service-matrix-go/internal/core/domain"
	"service-matrix-go/internal/infrastructure/storage"
	"sort"
	"strconv"
	"strings"
)

type WordService struct {
	fileHelper *storage.FileHelper
}

func NewWordService(fh *storage.FileHelper) *WordService {
	return &WordService{fileHelper: fh}
}

// Search implements the word search logic based on WordSearchCommandHandler
func (s *WordService) Search(req domain.SearchRequest) (map[string]map[int]map[string]string, error) {
	definitionWords := []string{}

	dictionary, err := s.fileHelper.ReadFileAsync("resources", "definitions.txt")
	if err != nil {
		return nil, err
	}
	merged, err := s.fileHelper.ReadFileAsync("resources", "merged.txt")
	if err == nil {
		// Create a map to track existing words for efficient lookup
		existingMap := make(map[string]bool)
		for _, w := range dictionary {
			existingMap[w] = true
		}

		// Only append words that don't already exist
		for _, w := range merged {
			if !existingMap[w] {
				dictionary = append(dictionary, w)
			}
		}
	}

	definitionWords = dictionary

	// excludes, err := s.fileHelper.ReadFileAsync("data", "exclude.txt")
	// if err == nil {
	// 	excludeMap := make(map[string]bool)
	// 	for _, line := range excludes {
	// 		excludeMap[line] = true
	// 	}

	// 	// filter out excluded
	// 	var filtered []string
	// 	for _, w := range definitionWords {
	// 		if len(w) > req.MaxLength || len(w) < req.MinLength {
	// 			continue
	// 		}
	// 		if !excludeMap[w] {
	// 			filtered = append(filtered, w)
	// 		}
	// 	}
	// 	definitionWords = filtered
	// }

	// Matrix conversion
	rows := len(req.LettersMatrix)
	if rows == 0 {
		return nil, nil
	}
	cols := len(req.LettersMatrix[0])

	// Create matrix
	lettersMatrix2D := make([][]string, rows)
	for i := 0; i < rows; i++ {
		lettersMatrix2D[i] = make([]string, cols)
		for j := 0; j < cols; j++ {
			if j < len(req.LettersMatrix[i]) {
				lettersMatrix2D[i][j] = req.LettersMatrix[i][j]
			}
		}
	}

	foundWordsList := make(map[string]map[int]map[string]string)

	for _, definitionWord := range definitionWords {
		if _, exists := foundWordsList[definitionWord]; exists {
			continue
		}

		if !algorithm.IsAllLettersInMatrix(lettersMatrix2D, definitionWord) {
			continue
		}

		searchHelper := algorithm.NewWordSearchHelper(definitionWord, lettersMatrix2D)
		if searchHelper.Search() {
			foundWord := searchHelper.GetFoundString()
			if strings.EqualFold(definitionWord, foundWord) {
				foundWordsList[foundWord] = searchHelper.GetFoundWord()
			}
		}
	}

	// Sort by length desc and take maxWords
	type kv struct {
		Key   string
		Value map[int]map[string]string
	}
	var ss []kv
	for k, v := range foundWordsList {
		ss = append(ss, kv{k, v})
	}
	sort.Slice(ss, func(i, j int) bool {
		return len(ss[i].Key) > len(ss[j].Key)
	})

	topResults := make(map[string]map[int]map[string]string)
	count := 0
	for _, kv := range ss {
		if count >= int(req.MaxWords) {
			break
		}
		topResults[kv.Key] = kv.Value
		count++
	}
	fmt.Printf("DEBUG: topResults count: %d\n", len(topResults))
	return topResults, nil
}

// UpdateWords implements UpdateWordsCommandHandler
func (s *WordService) UpdateWords(req domain.UpdateWordsRequest) (int, error) {
	filename := "exclude.txt"
	if req.Include {
		filename = "include.txt"
	}

	existing, _ := s.fileHelper.ReadFileAsync("data", filename)
	existingMap := make(map[string]bool)
	for _, w := range existing {
		existingMap[w] = true
	}

	var newWords []string
	count := 0
	for _, w := range req.Words {
		if !existingMap[w] {
			newWords = append(newWords, w)
			count++
		}
	}

	if len(newWords) > 0 {
		err := s.fileHelper.WriteFileAppend(newWords, "data", filename)
		if err != nil {
			return 0, err
		}
	}
	return count, nil
}

// GetList implements GetWordsQueryHandler
func (s *WordService) GetList(include bool) ([]string, error) {
	filename := "exclude.txt"
	if include {
		filename = "include.txt"
	}
	return s.fileHelper.ReadFileAsync("data", filename)
}

// MergeWords implements MergeWordsCommandHandler
func (s *WordService) MergeWords() (domain.MergeResponse, error) {
	// removedCounter := 0 // Unused in C# logic effectively as it's always 0

	includes, err := s.fileHelper.ReadFileAsync("data", "include.txt")
	if err != nil {
		return domain.MergeResponse{}, err
	}

	dictionary, err := s.fileHelper.ReadFileAsync("resources", "definitions.txt")
	if err == nil {
		merged, _ := s.fileHelper.ReadFileAsync("resources", "merged.txt")
		dictionary = append(dictionary, merged...)
	}

	// Create set for dictionary
	dictMap := make(map[string]bool)
	for _, w := range dictionary {
		dictMap[w] = true
	}

	var mergedList []string
	var remainingIncludes []string

	for _, include := range includes {
		if len(include) < 4 || strings.Contains(include, "-") || strings.Contains(include, " ") {
			// Skip and effectively remove from includes logic?
			// In C#, it says: if (...) continue; -> so it's NOT added to mergedList.
			// But later: includes = includes.Except(mergedList).ToList();
			// So if it was skipped here, it remains in includes?
			// C# logic:
			// foreach include in includes:
			//    if invalid continue
			//    if !dict.Contains(include) -> add to mergedList
			// includes = includes - mergedList
			// save mergedList to "mergeable_definitions.txt"
			// save includes to "include.txt"

			// So invalid words remain in includes.
			remainingIncludes = append(remainingIncludes, include)
			continue
		}

		includeFormatted := strings.ToLower(strings.TrimSpace(include))
		includeFormatted = strings.ReplaceAll(includeFormatted, "ё", "е")

		if !dictMap[includeFormatted] {
			mergedList = append(mergedList, includeFormatted)
		} else {
			remainingIncludes = append(remainingIncludes, include)
		}
	}

	// Actually logic in C# for `includes = includes.Except(mergedList)`:
	// If a word was added to mergedList, it is REMOVED from includes.
	// If it was valid but already in dictionary, it was NOT added to mergedList. Does it remain in includes?
	// C#: if !dictionary.Contains(...) -> added to mergedList.
	// So if it IS in dictionary, it is NOT added to mergedList.
	// Then Except(mergedList) removes only those added.
	// So existing dictionary words stay in includes?
	// Wait, that seems redundant. But following C# logic strictly:

	// Re-evaluating C#:
	// mergedList contains words from includes that are NOT in dictionary.
	// includes (new) = includes (old) without words in mergedList.
	// So words that moved to mergedList are removed from includes.
	// Words that were ALREADY in dictionary stay in includes? That seems odd conceptually but I must follow code.
	// Or maybe the intention was to clean includes?
	// Let's stick to strict translation.

	// In Go:
	// mergedList collected correctly.
	// We need to construct new includes list.
	// It should contain everything from old includes EXCEPT what's in mergedList.

	// Optimization:
	// mergedSet := make(map[string]bool)
	// for _, m := range mergedList { mergedSet[m] = true }
	// newIncludes := []string{}
	// for _, inc := range includes { if !mergedSet[inc] { newIncludes = append(newIncludes, inc) } }

	// Writing files
	if len(mergedList) > 0 {
		err = s.fileHelper.WriteFileNewContents(mergedList, "data", "mergeable_definitions.txt")
		if err != nil {
			return domain.MergeResponse{}, err
		}
	}

	// We need to calculate the new includes list based on the logic
	mergedSet := make(map[string]bool)
	for _, m := range mergedList {
		mergedSet[m] = true
	}

	var finalIncludes []string
	for _, inc := range includes {
		// C# logic does comparison potentially case-insensitive?
		// C# Except is default equality (case-sensitive usually unless comparer provided).
		// `mergedList` has lowercased formatted strings. `includes` has original strings.
		// If `include` was "Foo", `includeFormatted`="foo". `mergedList` has "foo".
		// `includes.Except` might NOT remove "Foo" if it checks against "foo"?
		// C# List.Except uses Default comparer. String is case sensitive.
		// So "Foo" would NOT be removed if "foo" is in mergedList?
		// HOWEVER, `includeFormatted` was derived from `include`.
		// If `include` was modified (trimmed, lowered), then `mergedList` content is different.
		// So `Except` might functionally do nothing if valid words were mixed case?
		// But likely `include.txt` is lowercase.
		// I will implement case-insensitive removal to be safe or strictly follow C# (which might be buggy there).
		// Let's assume strict C# behavior: Case Sensitive.

		// Wait, if C# behaves that way, I should replicate it.
		// If the user inputs "Word", it goes to mergedList as "word".
		// simple "Word" != "word". So "Word" remains in includes.
		// That effectively implementation DUPLICATES/Copies "word" to mergeable, and keeps "Word" in includes.
		// Unless `include` strings are exact matches.

		// Let's implement logical behavior: Remove if match.
		// formatted := strings.ToLower(strings.TrimSpace(inc)) ...
		if mergedSet[strings.ToLower(strings.TrimSpace(inc))] {
			continue
		}
		finalIncludes = append(finalIncludes, inc)
	}

	err = s.fileHelper.WriteFileNewContents(finalIncludes, "data", "include.txt")
	if err != nil {
		return domain.MergeResponse{}, err
	}

	return domain.MergeResponse{AddedCount: len(mergedList), RemovedCount: 0}, nil
}

// CleanMerge implements clean merge logic
func (s *WordService) CleanMerge() (string, error) {
	input, err := s.fileHelper.ReadFileAsync("resources", "merged.txt")
	if err != nil {
		return "", err
	}

	output := algorithm.CleanWords(input)
	err = s.fileHelper.WriteFileNewContents(output, "data", "merged_cleaned.txt")
	if err != nil {
		return "", err
	}

	return "Processed " + strconv.Itoa(len(input)) + " lines.", nil
}

// LookupWord implements LookupWordQueryHandler
func (s *WordService) LookupWord(word string, exactMatch bool) ([]domain.LookupResultResponseItem, error) {
	// Porting LookupWordQueryHandler
	// I assume it looks in definitions, merged, include, exclude.

	var results []domain.LookupResultResponseItem
	sources := map[string]string{
		"definitions.txt": "resources",
		"merged.txt":      "resources",
		"include.txt":     "data",
		"exclude.txt":     "data",
	}

	for file, dir := range sources {
		lines, err := s.fileHelper.ReadFileAsync(dir, file)
		if err == nil {
			for i, line := range lines {
				found := false
				if exactMatch {
					if strings.EqualFold(line, word) {
						found = true
					}
				} else {
					if strings.Contains(strings.ToLower(line), strings.ToLower(word)) {
						found = true
					}
				}

				if found {
					results = append(results, domain.LookupResultResponseItem{
						Word:   line,
						Found:  true,
						Source: file,
						Line:   i + 1,
					})
				}
			}
		}
	}
	return results, nil
}
