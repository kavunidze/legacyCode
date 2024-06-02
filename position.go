package gopuzzlegame

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func (p *Position) CompareTo(other Position) int {
	if p.Y < other.Y {
		return -1
	} else if p.Y > other.Y {
		return 1
	} else {
		if p.X < other.X {
			return -1
		} else if p.X > other.X {
			return 1
		} else {
			return 0
		}
	}
}

func (p *Position) CompareToBool(other Position) bool {
	if p.Y < other.Y {
		return false
	} else if p.Y > other.Y {
		return true
	} else {
		if p.X < other.X {
			return false
		} else if p.X > other.X {
			return true
		} else {
			return false
		}
	}
}
