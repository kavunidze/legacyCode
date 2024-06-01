package gopuzzlegame

import (
	"fmt"
	"math/rand"

	"github.com/pkg/errors"
)

const (
	PuzzleGameStatusNew = iota
	PuzzleStatusInProgress
	PuzzleStatusFinished
	PuzzleStatusReachedStepLimit
)

type PuzzleController struct {
	PuzzleStatus int32
	Puzzle       *Puzzle
	StepsTaken   int32
	Steps        int32
}

func (p *PuzzleController) TapTile(tile *Tile) error {
	movable, err := p.Puzzle.IsTileMovable(tile)
	if err != nil {
		return err
	}
	if p.PuzzleStatus == PuzzleStatusInProgress && movable {
		mutablePuzzle := Puzzle{Tiles: p.Puzzle.Tiles}
		var tiles []*Tile
		puzzle, err := mutablePuzzle.MoveTiles(tile, tiles)
		if err != nil {
			return err
		}
		p.Puzzle = puzzle
		p.Puzzle.Sort()
		p.StepsTaken++
		isComplete, err := p.Puzzle.IsComplete()
		if err != nil {
			return err
		}
		if isComplete {
			p.PuzzleStatus = PuzzleStatusFinished
		} else if p.StepsTaken == p.Steps {
			p.PuzzleStatus = PuzzleStatusReachedStepLimit
		}
		return nil
	}
	return errors.New(fmt.Sprintf("tile can't change location: puzzle_status: %v, movable: %v", p.PuzzleStatus, movable))
}

func GeneratePuzzle(size int, shuffle bool) (*Puzzle, error) {
	var correctPositions []Position
	var currentPositions []Position
	whitespacePosition := Position{
		X: size,
		Y: size,
	}

	for y := 1; y <= size; y++ {
		for x := 1; x <= size; x++ {
			if x == size && y == size {
				correctPositions = append(correctPositions, whitespacePosition)
				currentPositions = append(currentPositions, whitespacePosition)
			} else {
				position := Position{
					X: x,
					Y: y,
				}
				correctPositions = append(correctPositions, position)
				currentPositions = append(currentPositions, position)
			}
		}
	}

	if shuffle {
		rand.Shuffle(len(currentPositions), func(i, j int) {
			currentPositions[i], currentPositions[j] = currentPositions[j], currentPositions[i]
		})
	}

	tiles := getTileListFromPositions(size, correctPositions, currentPositions)
	puzzle := &Puzzle{Tiles: tiles}

	if shuffle {
		isSolvable, err := puzzle.IsSolvable()
		if err != nil {
			return nil, err
		}
		numOfCorrectTiles, err := puzzle.GetNumberOfCorrectTiles()
		if err != nil {
			return nil, err
		}

		for !isSolvable || numOfCorrectTiles != 0 {
			rand.Shuffle(len(currentPositions), func(i, j int) {
				currentPositions[i], currentPositions[j] = currentPositions[j], currentPositions[i]
			})
			puzzle = &Puzzle{Tiles: getTileListFromPositions(size, correctPositions, currentPositions)}
			isSolvable, err = puzzle.IsSolvable()
			if err != nil {
				return nil, err
			}
			numOfCorrectTiles, err = puzzle.GetNumberOfCorrectTiles()
			if err != nil {
				return nil, err
			}
		}
	}

	return puzzle, nil
}

func getTileListFromPositions(size int, correctPositions, currentPositions []Position) []*Tile {
	whitespacePosition := Position{
		X: size,
		Y: size,
	}

	var result []*Tile
	for i := 1; i <= size*size; i++ {
		if i == size*size {
			result = append(result, &Tile{
				Value:           i,
				CorrectPosition: &whitespacePosition,
				CurrentPosition: currentPositions[i-1],
				IsWhitespace:    true,
			})
		} else {
			result = append(result, &Tile{
				Value:           i,
				CorrectPosition: &correctPositions[i-1],
				CurrentPosition: currentPositions[i-1],
			})
		}
	}

	return result
}
