package algorithm

import (
	"math"
	"sort"
	"strconv"
	"strings"
)

// WordSearchHelper structure to hold search state
type WordSearchHelper struct {
	arLetters       [][]string
	arLettersStatic [][]string
	arWord          []string
	foundWord       map[int]map[string]string // index -> char -> "row col"
	iFirstColumn    int
	iFirstRow       int
	sFoundString    strings.Builder
	sWord           string
}

// NewWordSearchHelper creates a new helper instance
func NewWordSearchHelper(sWord string, arLetters [][]string) *WordSearchHelper {
	helper := &WordSearchHelper{
		sWord:           sWord,
		arWord:          strings.Split(sWord, ""),
		arLettersStatic: CopyArray(arLetters),
		foundWord:       make(map[int]map[string]string),
	}
	helper.arLetters = CopyArray(helper.arLettersStatic)
	return helper
}

// Search initiates the search for the word in the matrix
func (h *WordSearchHelper) Search() bool {
	startPositions := h.FindLetterLocations()

	for _, pos := range startPositions {
		h.sFoundString.Reset()
		h.foundWord = make(map[int]map[string]string) // Reset found word for new start position
		h.arLetters = CopyArray(h.arLettersStatic)    // Reset matrix

		zero := 0
		row := pos.Row
		col := pos.Col
		isWordFound := h.FindWord(&zero, &row, &col)
		if isWordFound {
			return true
		}
	}
	return false
}

// FindLetterLocations finds all locations of the first letter of the word
func (h *WordSearchHelper) FindLetterLocations() []Position {
	if len(h.sWord) == 0 {
		return nil
	}
	letter := h.arWord[0]
	var locations []Position
	for i := 0; i < len(h.arLetters); i++ {
		for j := 0; j < len(h.arLetters[i]); j++ {
			if h.arLetters[i][j] == letter {
				locations = append(locations, Position{Row: i, Col: j})
			}
		}
	}

	return locations
}

// FindWord recursive/iterative search logic
func (h *WordSearchHelper) FindWord(iWordIndexWrapper *int, iMatrixStartWrapper *int, jMatrixStartWrapper *int) bool {
	// Porting the logic directly is tricky because of the reference passing in C# (int?).
	// Adapting to Go's value/pointer semantics.

	iWordIndex := 0
	if iWordIndexWrapper != nil {
		iWordIndex = *iWordIndexWrapper
	}
	iMatrixStart := 0
	if iMatrixStartWrapper != nil {
		iMatrixStart = *iMatrixStartWrapper
	}
	jMatrixStart := 0
	if jMatrixStartWrapper != nil {
		jMatrixStart = *jMatrixStartWrapper
	}

	for w := iWordIndex; w < len(h.arWord) && iWordIndex <= len(h.foundWord); w++ {
		iWordIndex = w
		sSearchChar := h.arWord[w]
		bFound := false

		for i := iMatrixStart; i < len(h.arLetters) && !bFound; i++ {
			startJ := 0
			if i == iMatrixStart {
				startJ = jMatrixStart
			}

			for j := startJ; j < len(h.arLetters[i]) && !bFound; j++ {
				// string comparison
				if strings.EqualFold(sSearchChar, h.arLetters[i][j]) && h.IsNeighborToPrevLetter(i, j, iWordIndex, sSearchChar) {
					if h.IsNeighborToNextLetter(i, j, h.arWord, iWordIndex, CopyArray(h.arLetters)) {
						hLetLoc := map[string]string{
							sSearchChar: strconv.Itoa(i) + " " + strconv.Itoa(j),
						}
						h.foundWord[iWordIndex] = hLetLoc
						h.sFoundString.WriteString(h.arLetters[i][j])
						h.arLetters[i][j] = "*"
						bFound = true

						// Reset search space for next iteration to adjacent
						// Actual logic in C# sets iMatrixStart/jMatrixStart to i-1, i-1 but that seems to apply to the NEXT looping.
						// However, the loop continues.

						// The original C# code modifies iMatrixStart and jMatrixStart but they are loop variables 'i' and 'j' in the inner loop so it doesn't affect current iteration directly unless we break/continue.
						// But here, we just set bFound=true which breaks inner loops.

						iMatrixStart = i - 1
						if iMatrixStart < 0 {
							iMatrixStart = 0
						}
						jMatrixStart = j - 1
						if jMatrixStart < 0 {
							jMatrixStart = 0
						}

						if strings.EqualFold(h.sFoundString.String(), h.sWord) {
							return true
						}
					}
				}
			}
			jMatrixStart = 0 // Reset for next row
		}
		// iMatrixStart preserved for outer loop logic if needed, but self-assignment is redundant
		// Actually in C# FindWord recieves iWordIndex, iMatrixStart, jMatrixStart as nullable ints. If null they start at 0.
		// The loop for 'w' continues.

		if !bFound {
			// If not found in this pass
			break // Breaking the word building loop
		}
	}

	// The "GetNextFirstLetter" logic in C# is invoked if the standard loop fails
	// This part handles backtracking/retrying from a different start point if the greedy approach failed?
	// The C# code is quite specific. I'll attempt to simplify or map closely.

	// Due to complexity and potential recursion in "GetNextFirstLetter" -> calling "FindWord(1, null, null)",
	// creating a simplified version of standard backtracking might be safer than direct port if I don't fully replicate the state mutations.

	// However, sticking to direct port for safety.
	tempCol := h.iFirstColumn
	tempRow := h.iFirstRow
	h.UpdateNextFirstLetterStartPos()

	// Check if we need to backtrack/retry
	if !(tempCol == h.iFirstColumn && tempRow == h.iFirstRow) && h.GetNextFirstLetter(h.iFirstColumn, h.iFirstRow) {
		one := 1
		return h.FindWord(&one, nil, nil)
	}

	return false
}

func (h *WordSearchHelper) GetNextFirstLetter(iMatrixStart, jMatrixStart int) bool {
	i3Start := iMatrixStart
	jStart := jMatrixStart

	sSearchChar := h.arWord[0]
	bFound := false
	h.arLetters = CopyArray(h.arLettersStatic)
	h.foundWord = make(map[int]map[string]string)
	h.sFoundString.Reset()

	for i3 := i3Start; i3 < len(h.arLetters) && !bFound; i3++ {
		j := 0
		if i3 == i3Start {
			j = jStart
		}
		for ; j < len(h.arLetters[i3]); j++ {
			if bFound {
				break
			}
			if h.arLetters[i3][j] == sSearchChar {
				if h.IsNeighborToNextLetter(i3, j, h.arWord, 0, CopyArray(h.arLetters)) {
					h.sFoundString.WriteString(sSearchChar)
					hLetLoc := map[string]string{
						sSearchChar: strconv.Itoa(i3) + " " + strconv.Itoa(j),
					}
					h.foundWord[0] = hLetLoc
					h.arLetters[i3][j] = "*"
					bFound = true
					break
				}
			}
		}
	}
	return bFound
}

func (h *WordSearchHelper) UpdateNextFirstLetterStartPos() {
	tempCol := h.iFirstColumn
	tempRow := h.iFirstRow
	if h.iFirstColumn < 5 && h.iFirstRow < 5 { // Hardcoded 5? Assuming matrix size constraint or logic from C#
		h.iFirstRow++
		if h.iFirstRow > 4 {
			h.iFirstRow = h.iFirstRow % 4 // Original Code: _iFirstRow %= 4; wait, 5%4 is 1.
			h.iFirstColumn++
		}
		// Original: if (_iFirstColumn > 4 || _iFirstRow > 4) -> reset
		if h.iFirstColumn > 4 || h.iFirstRow > 4 {
			h.iFirstColumn = tempCol
			h.iFirstRow = tempRow
		}
	}
}

func (h *WordSearchHelper) IsNeighborToNextLetter(iCurrentX, iCurrentY int, arWord2 []string, iWordIndex int, arLettersLoc [][]string) bool {
	if iWordIndex == len(arWord2)-1 || iWordIndex == 0 { // 0 check in C# IsNeighborToNextLetter seems to allow first letter?
		// Actually original C# says: if (iWordIndex == arWord2.Length - 1 || iWordIndex == 0) return true;
		return true
	}
	sNextLetter := arWord2[iWordIndex+1]
	for dX := 1; dX <= 3; dX++ {
		for dY := 1; dY <= 3; dY++ {
			neighborX := (iCurrentX + dX) - 2
			neighborY := (iCurrentY + dY) - 2

			// Boundary checks
			rows := len(arLettersLoc)
			cols := len(arLettersLoc[0])

			if !(neighborX == iCurrentX && neighborY == iCurrentY) &&
				neighborX >= 0 && neighborX < rows &&
				neighborY >= 0 && neighborY < cols &&
				strings.EqualFold(sNextLetter, arLettersLoc[neighborX][neighborY]) {

				secondNextNeighbor := true
				if len(arWord2) > iWordIndex+2 {
					arLettersLocTemp := CopyArray(arLettersLoc)
					arLettersLocTemp[iCurrentX][iCurrentY] = "*"
					secondNextNeighbor = h.IsNeighborToNextLetter(neighborX, neighborY, arWord2, iWordIndex+1, arLettersLocTemp)
				}
				if secondNextNeighbor {
					return true
				}
			}
		}
	}
	return false
}

func (h *WordSearchHelper) IsNeighborToPrevLetter(iCol, iRow, iWordIndex int, sLet string) bool {
	if iWordIndex == 0 {
		h.iFirstColumn = iCol
		h.iFirstRow = iRow
		return true
	}
	hLetterIdx, ok := h.foundWord[iWordIndex-1]
	if !ok {
		return false
	}
	var value string
	for _, v := range hLetterIdx {
		value = v
		break
	}
	sIndex := strings.Split(value, " ")
	if len(sIndex) < 2 {
		return false
	}

	// In C#  { sSearchChar, $"{i} {j}" } where i is row, j is col.
	// C# retrieval:
	// int iPrevRow = int.Parse(sIndex[1]);
	// int xDelta = Math.Abs(Math.Abs(int.Parse(sIndex[0])) - Math.Abs(iCol));
	// int yDelta = Math.Abs(Math.Abs(iPrevRow) - Math.Abs(iRow));
	// Wait, in C# loops are: for (int i ...), for (int j ...). 'i' is usually row, 'j' is col.
	// Stored as $"{i3} {j}".
	// Then parsed: sIndex[0] is i (row), sIndex[1] is j (col).
	// The C# logic: iPrevRow = sIndex[1] (which is COL).
	// xDelta = abs(sIndex[0] (ROW) - iCol (ROW?? or COL?)).
	// In C params: (int iCol, int iRow).
	// Usually i=Row, j=Col.
	// So if function called as IsNeighborToPrevLetter(i, j...) -> iCol=Row, iRow=Col.
	// Let's verify usage: IsNeighborToPrevLetter(i, j, ...) -> i is loop var 1 (row), j is loop var 2 (col).
	// So iCol is Row, iRow is Col.

	// C# Logic:
	// iPrevRow = sIndex[1] (The COL of prev)
	// xDelta = abs(sIndex[0] (The ROW of prev) - iCol (Current ROW))
	// yDelta = abs(iPrevRow (The COL of prev) - iRow (Current COL))
	// Check: xDelta < 2 && yDelta < 2.
	// This checks adjacency in X and Y.

	prevRow, _ := strconv.Atoi(sIndex[0])
	prevCol, _ := strconv.Atoi(sIndex[1])

	// iCol is Current Row, iRow is Current Col based on call site usage (i, j)
	currRow := iCol
	currCol := iRow

	xDelta := int(math.Abs(float64(prevRow - currRow)))
	yDelta := int(math.Abs(float64(prevCol - currCol)))

	return xDelta < 2 && yDelta < 2
}

func (h *WordSearchHelper) GetFoundString() string {
	return h.sFoundString.String()
}

func (h *WordSearchHelper) GetFoundWord() map[int]map[string]string {
	return h.foundWord
}

// Helper functions

type Position struct {
	Row int
	Col int
}

func CopyArray(source [][]string) [][]string {
	if source == nil {
		return nil
	}
	dest := make([][]string, len(source))
	for i := range source {
		dest[i] = make([]string, len(source[i]))
		copy(dest[i], source[i])
	}
	return dest
}

// IsAllLettersInMatrix checks if all letters of the word are in the matrix
func IsAllLettersInMatrix(matrix [][]string, wholeWord string) bool {
	allArrayLetters := make(map[rune]bool)

	for i := 0; i < len(matrix); i++ {
		for j := 0; j < len(matrix[i]); j++ {
			runes := []rune(matrix[i][j])
			if len(runes) > 0 {
				allArrayLetters[runes[0]] = true
			} else {
				allArrayLetters['*'] = true
			}
		}
	}

	for _, c := range wholeWord {
		if !allArrayLetters[c] {
			return false
		}
	}
	return true
}

// CleanWords filters words based on rules
func CleanWords(input []string) []string {
	var output []string
	for _, word := range input {
		if len(word) < 8 || len(word) > 24 || strings.Contains(word, " ") || strings.Contains(word, "-") {
			continue
		}
		output = append(output, word)
	}

	// Sort descending by length
	sort.Slice(output, func(i, j int) bool {
		return len(output[i]) > len(output[j])
	})

	return output
}
