package ninetyfive

import (
	"errors"
	"strings"

	"github.com/mtgban/go-mtgban/mtgmatcher"
)

func preprocess(allCards NFCard, key, lang string, foil bool) (*mtgmatcher.InputCard, error) {
	card, found := allCards[key]
	if !found {
		return nil, errors.New("key not found")
	}

	cardName := card.CardName
	variant := card.CardNum
	edition := card.SetCode

	switch edition {
	case "pARL96":
		edition = "PARL"
	case "MB1", "pHED", "pCTB", "PAGL":
		edition = "PLST"
	case "pPOD":
		edition = "POR"
		variant = card.SetName
	default:
		_, err := mtgmatcher.GetSet(card.SetCode)
		if err != nil {
			edition = card.SetName
		}
	}

	if strings.HasSuffix(variant, "b") && strings.Contains(cardName, "Rainbow Foil") {
		variant = strings.TrimSuffix(variant, "b")
	}
	if strings.HasSuffix(variant, "r") {
		variant = strings.TrimSuffix(variant, "r")
		variant += " Rainbow Foil"
	}
	if strings.HasSuffix(variant, "u") {
		variant = strings.TrimSuffix(variant, "u")
		variant += " Surge Foil"
	}

	return &mtgmatcher.InputCard{
		Name:      cardName,
		Edition:   edition,
		Variation: variant,
		Foil:      foil,
		Language:  lang,
	}, nil
}
