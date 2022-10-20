package toamagic

import (
	"errors"
	"strings"

	"github.com/kodabb/go-mtgban/mtgmatcher"
)

var cardTable = map[string]string{
	"Scornful Ã†ther-Lich": "Scornful Aether-Lich",
}

var promoTags = []string{
	"2016 Welcome Deck",
	"Foil",
	"Godzilla Lands",
}

func preprocess(cardName, edition, variant string) (*mtgmatcher.Card, error) {
	s := mtgmatcher.SplitVariants(cardName)
	cardName = s[0]
	if len(s) > 1 {
		if variant != "" {
			variant += " "
		}
		variant += strings.Join(s[1:], " ")
	}

	for _, tag := range promoTags {
		if strings.HasSuffix(cardName, tag) {
			if cardName == "Foil" {
				continue
			}
			cardName = strings.TrimSuffix(cardName, " "+tag)
			if variant != "" {
				variant += " "
			}
			variant += tag
		}
	}

	isFoil := strings.Contains(variant, "Foil")
	if isFoil {
		variant = strings.Replace(variant, "Foil", "", 1)
		variant = strings.TrimSpace(variant)
	}

	switch edition {
	case "Portal":
		if variant == "1" {
			variant = "Flavor Text"
		} else if variant == "2" {
			variant = "Reminder Text"
		}
	case "Archenemy Nicol Bolas":
		if cardName == "Highland Lake" || cardName == "Submerged Boneyard" {
			return nil, errors.New("does not exist")
		}
	}

	lutName, found := cardTable[cardName]
	if found {
		cardName = lutName
	}

	return &mtgmatcher.Card{
		Name:      cardName,
		Variation: variant,
		Edition:   edition,
		Foil:      isFoil,
	}, nil
}
