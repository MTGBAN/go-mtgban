package abugames

import (
	"errors"
	"strings"

	"github.com/kodabb/go-mtgban/mtgjson"
	"github.com/kodabb/go-mtgban/mtgmatcher"
)

var cardTable = map[string]string{
	// Typos
	"Bogart Brute":                        "Boggart Brute",
	"Deathgazer Cockatrice":               "Deathgaze Cockatrice",
	"Fireblade Artist Ravnica Allegiance": "Fireblade Artist",
	"Jace, the Mind Sculpor":              "Jace, the Mind Sculptor",
	"Mindblade Rendor":                    "Mindblade Render",
	"Neglected Hierloom / Ashmouth Blade": "Neglected Heirloom // Ashmouth Blade",
	"Rathi Berserker":                     "Aerathi Berserker",
	"Skin Invasion / Skin Shredder":       "Skin Invasion // Skin Shedder",
	"Soulmemder":                          "Soulmender",
	"Svagthos, the Restless Tomb":         "Svogthos, the Restless Tomb",
	"Visitor from Planet":                 "Visitor from Planet Q",

	// Funny cards
	"B.F.M. Big Furry Monster Left":   "B.F.M. (Big Furry Monster 28)",
	"B.F.M. Big Furry Monster Right":  "B.F.M. (Big Furry Monster 29)",
	"Surgeon Commander":               "Surgeon ~General~ Commander",
	"No Name":                         "_____",
	"Who What When Where Why":         "Who",
	"Absolute Longest Card Name Ever": "Our Market Research Shows That Players Like Really Long Card Names So We Made this Card to Have the Absolute Longest Card Name Ever Elemental",
}

var promoTags = []string{
	"Alternate Art Duelist",
	"Arena",
	"Book",
	"Buy-a-Box",
	"Convention",
	"Draft Weekend",
	"FNM",
	"Into Pack",
	"Judge",
	"Launch",
	"OPEN HOUSE FULL ART",
	"Open House",
	"Planeswalker Weekend",
	"Prerelease",
	"Preview",
	"Promo",
	"Release",
	"SDCC",
	"Store Championship",
	"TopDeck Magazine",
}

func preprocess(card *ABUCard) (*mtgmatcher.Card, error) {
	lang := ""
	if len(card.Language) > 0 {
		switch card.Language[0] {
		case "English":
			lang = "EN"
		case "French":
			lang = "FR"
		case "German":
			lang = "DE"
		case "Italian":
			lang = "IT"
		case "Spanish":
			lang = "ES"
		case "Portuguese":
			lang = "PT"
		case "Japanese":
			lang = "JP"
		case "Korean":
			lang = "KR"
		case "Chinese Simplified":
			lang = "CH"
		case "Russian":
			lang = "RU"
		default:
			lang = card.Language[0]
		}
	}

	if lang != "EN" || strings.Contains(card.Title, "Non-English") {
		return nil, errors.New("non-english card")
	}

	// Non-Singles magic cards
	switch card.Layout {
	case "Scheme", "Plane", "Phenomenon":
		return nil, errors.New("non-single card")
	}
	if strings.Contains(card.DisplayTitle, "Oversized") ||
		strings.Contains(card.DisplayTitle, "Charlie Brown") {
		return nil, errors.New("non-single card")
	}
	// Non-existing cards
	switch card.DisplayTitle {
	case "Steward of Valeron (Dengeki Character Festival) - FOIL",
		"Captain's Claws (Goldnight Castigator Shadow) - FOIL",
		"Beast of Burden (Prerelease - No Expansion Symbol) - FOIL",
		"Hymn to Tourach (B - Mark Justice - 1996)",
		"Mountain (6th Edition 343 - Mark Le Pine - 1999)":
		return nil, errors.New("untracked card")
	}
	// Unique cards
	if strings.HasPrefix(card.Title, "ID#") {
		return nil, errors.New("unique card")
	}
	switch card.Id {
	case "1604919", "1604921", "1604922", // Living Twister
		"1604802", "1604801", "1604799": // Commence the Endgame
		return nil, errors.New("duplicated card")
	}

	isFoil := strings.Contains(strings.ToLower(card.DisplayTitle), " foil")

	// Split by -, rebuild the cardname in a standardized way
	variation := ""
	vars := strings.Split(card.DisplayTitle, " - ")
	cardName := vars[0]
	if len(vars) > 1 {
		if vars[len(vars)-1] == card.Edition {
			vars = vars[:len(vars)-1]
		}

		variation = strings.Join(vars[1:], " ")

		// Fix some untagged prerelease cards
		// Nahiri's Wrath, Tendershoot Dryad
		if strings.Contains(variation, card.Edition+" FOIL") {
			variation = strings.Replace(variation, card.Edition+" FOIL", "Prerelease", 1)
		}
	}

	// Split by ()
	vars = mtgmatcher.SplitVariants(cardName)
	cardName = vars[0]
	if len(vars) > 1 {
		oldVariation := variation
		variation = strings.Join(vars[1:], " ")
		if oldVariation != "" {
			variation += " " + oldVariation
		}
	}

	// Cleanup variation as necessary
	if variation != "" {
		variation = strings.Replace(variation, "FOIL", "", -1)

		variation = strings.Replace(variation, "(", "", -1)
		variation = strings.Replace(variation, ")", "", -1)

		variation = strings.Replace(variation, "Not Tournament Legal", "", 1)

		variation = strings.TrimSpace(variation)

		// These are the definition of generic promos
		if strings.Contains(variation, "Magic League") ||
			mtgjson.NormContains(variation, "Game Day") ||
			mtgjson.NormContains(variation, "Gift Box") ||
			mtgjson.NormContains(variation, "Convention") ||
			strings.Contains(variation, "Intro Pack") {
			variation = "Promo"
		}

		isPromo := false
		for _, tag := range promoTags {
			if strings.Contains(variation, tag) {
				isPromo = true
				break
			}
		}
		if isPromo {
			// Handle promo cards appearing in multiple editions
			// like Sorcerous Spyglass
			switch card.Edition {
			case "Aether Revolt",
				"Ixalan",
				"Core Set 2020 / M20",
				"Core Set 2021 / M21",
				"Throne of Eldraine":
				variation += " " + card.Edition
			}
			card.Edition = "Promo"
		}
	}

	switch card.Edition {
	case "Promo":
		switch cardName {
		case "Skirk Marauder":
			card.Edition = "Arena League 2003"
		case "Captain Sisay":
			card.Edition = "Secret Lair Drop"
		case "Elvish Aberration":
			if variation == "FNM" {
				variation = "Arena"
			}
		case "Elvish Lyrist":
			if variation == "FNM" {
				variation = "JSS Foil"
			}
		case "Island":
			if variation == "Arena 1999 No Symbol Promo" {
				variation = "Arena 1999 misprint"
			}
		case "Stocking Tiger":
			if variation == "Target" {
				variation = "misprint"
			}
		case "Psychatog":
			if variation == "FNM" {
				variation = "Textless"
			}
		case "Rukh Egg":
			if variation == "Prerelease" {
				variation = "Release"
			}
		case "Sol Ring":
			if variation == "Commander" {
				variation = "MagicFest 2019"
			}
		case "Hall of Triumph":
			if variation == "Promo" {
				variation = "Hero's Path"
			}
		case "Disenchant":
			if variation == "Arena" && isFoil {
				variation = "FNM 2003"
			}
		}
		if strings.Contains(variation, "United Kingdom") {
			variation = strings.Replace(variation, "United Kingdom", "U.K.", 1)
		} else if strings.Contains(variation, "Scandanavia") {
			variation = strings.Replace(variation, "Scandanavia", "Scandinavia", 1)
		}
	case "Anthologies":
		if cardName == "Mountain" {
			if variation == "A" {
				variation = "B"
			} else if variation == "B" {
				variation = "A"
			}
		}
	case "World Championship":
		if cardName == "City of Brass" {
			if variation == "Leon Lindback 1996" {
				variation = "Sideboard Leon Lindback 1996"
			}
		}
	case "Ikoria: Lair of Behemoths":
		if strings.Contains(cardName, " / ") {
			s := strings.Split(cardName, " / ")
			cardName = s[0]
		}
	}

	name, found := cardTable[cardName]
	if found {
		cardName = name
	}

	return &mtgmatcher.Card{
		Name:      cardName,
		Variation: variation,
		Edition:   card.Edition,
		Foil:      isFoil,
	}, nil
}
