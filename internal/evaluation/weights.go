package eval

import (
	"os"

	"gopkg.in/yaml.v3"
)

const WEIGHT_PATH = "./weights/weights.yml"

func LoadWeights() (*Weights, error) {
	var weights Weights
	wFile, err := os.Open(WEIGHT_PATH)
	if err != nil {
		return nil, err
	}
	defer wFile.Close()

	d := yaml.NewDecoder(wFile)
	err = d.Decode(&weights)
	if err != nil {
		return nil, err
	}
	return &weights, nil
}

func getPieceWeight(piece uint8) int {
	switch piece % 6 {
	case 1:
		return weights.Pieces.Pawn
	case 2:
		return weights.Pieces.Bishop
	case 3:
		return weights.Pieces.Knight
	case 4:
		return weights.Pieces.Rook
	case 5:
		return weights.Pieces.Queen
	default:
		return 0
	}
}

type Weights struct {
	Moves  Moves  `yaml:"moves"`
	Pieces Pieces `yaml:"pieces"`
	Knight Knight `yaml:"knight"`
	Bishop Bishop `yaml:"bishop"`
	Pawn   Pawn   `yaml:"pawn"`
}

type Pieces struct {
	Pawn   int `yaml:"pawn"`
	Knight int `yaml:"knight"`
	Bishop int `yaml:"bishop"`
	Rook   int `yaml:"rook"`
	Queen  int `yaml:"queen"`
	King   int `yaml:"king"`
}

type Knight struct {
	Center22 int `yaml:"center22"`
	Center44 int `yaml:"center44"`
	InnerRim int `yaml:"innerRim"`
	OuterRim int `yaml:"outerRim"`
}

type Bishop struct {
	MajorDiag int `yaml:"majorDiag"`
	MinorDiag int `yaml:"minorDiag"`
}

type Pawn struct {
	Passed    int `yaml:"passed"`
	Protected int `yaml:"protected"`
	Doubled   int `yaml:"doubled"`
	Isolated  int `yaml:"isolated"`
	Center22  int `yaml:"center22"`
	Center44  int `yaml:"center44"`
	Advance   int `yaml:"advance"`
}

type Moves struct {
	Move    int `yaml:"move"`
	Capture int `yaml:"capture"`
}
