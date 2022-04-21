package strikezone

import (
	"errors"
	"strings"

	"github.com/kodabb/go-mtgban/mtgmatcher"
)

func preprocess(cardName, edition, notes string) (*mtgmatcher.Card, error) {
	var variation string

	cn, found := cardTable[cardName]
	if found {
		cardName = cn
	}

	// Sometimes the buylist specifies tags at the end of the card name,
	// but without parenthesis, so make sure they are present.
	for _, tag := range tagsTable {
		if strings.HasSuffix(cardName, tag) {
			cardName = strings.Replace(cardName, tag, "("+tag+")", 1)
			break
		}
	}

	cardName = strings.Replace(cardName, "[", "(", 1)
	cardName = strings.Replace(cardName, "]", ")", 1)

	switch {
	// Tyop
	case strings.HasPrefix(cardName, "Snow-Cover "):
		cardName = strings.Replace(cardName, "Snow-Cover ", "Snow-Covered ", 1)
	// Found at beginning, move it to variation
	case strings.HasPrefix(cardName, "Borderless"):
		cardName = strings.TrimPrefix(cardName, "Borderless ")
		cardName = strings.TrimPrefix(cardName, "Alt Art ")
		variation = "Borderless"
	// Found at beginning, move it to variation
	case strings.HasPrefix(cardName, "Showcase"):
		cardName = strings.TrimPrefix(cardName, "Showcase ")
		variation = "Showcase"
	// Found at beginning, move it to variation
	case strings.HasPrefix(cardName, "Extended Art"):
		cardName = strings.TrimPrefix(cardName, "Extended Art ")
		variation = "Extended Art"
	// Found at end, move it to variation
	case strings.HasSuffix(cardName, "JPN ALT ART PRERELEASE"):
		cardName = strings.TrimSuffix(cardName, " JPN ALT ART PRERELEASE")
		variation = "Prerelease Japanese"
	// Found at end, move it to edition
	case strings.HasSuffix(cardName, "Ultimate Edition"):
		cardName = strings.TrimSuffix(cardName, " Ultimate Edition")
		edition = "Secret Lair: Ultimate Edition"
	// Found at end, move it to edition
	case strings.HasSuffix(cardName, "Godzilla") && mtgmatcher.IsBasicLand(cardName):
		cardName = strings.TrimSuffix(cardName, " Godzilla")
		edition = "SLD"
	// Found at beginning, just drop it
	case strings.HasPrefix(cardName, "Alt Art"):
		cardName = strings.TrimPrefix(cardName, "Alt Art ")
		edition = "SLD"
	// APAC and EURO lands, drop specifier
	case strings.Contains(cardName, "APAC") || strings.Contains(cardName, "EURO"):
		edition = "European Land Program"
		if strings.Contains(cardName, "APAC") {
			edition = "Asia Pacific Land Program"
		}
		variants := mtgmatcher.SplitVariants(cardName)
		cardName = variants[0]
		fields := strings.Fields(cardName)
		cardName = fields[0]
		variation = variants[1]
	// Ravnica weekend lands, move to variation
	case strings.Contains(cardName, "Ravnica Weekend"):
		fields := strings.Fields(cardName)
		cardName = fields[0]
		variation = strings.Join(fields[1:], " ")
	case strings.HasPrefix(cardName, "B.F.M.") && strings.Contains(cardName, "#"):
		for _, field := range strings.Fields(cardName) {
			if strings.HasPrefix(field, "#") {
				cardName = "B.F.M."
				variation = field[1:]
				break
			}
		}
	}

	ed, found := card2setTable[cardName]
	if found {
		edition = ed
	}

	variants := mtgmatcher.SplitVariants(cardName)
	cardName = variants[0]
	if len(variants) > 1 {
		if variation != "" {
			variation += " "
		}
		variation = strings.Join(variants[1:], " ")
	}

	// Repeat to catch numbers
	if mtgmatcher.IsBasicLand(cardName) {
		num := mtgmatcher.ExtractNumber(cardName)
		if num != "" {
			cardName = strings.Replace(cardName, num, "", 1)
			cardName = strings.TrimSpace(cardName)
			if variation != "" {
				variation += " "
			}
			variation += num
		}

		for _, tag := range tagsTable {
			if strings.HasSuffix(cardName, tag) {
				cardName = strings.TrimSuffix(cardName, " "+tag)
				if variation != "" {
					variation += " "
				}
				variation += tag
				break
			}
		}
	}

	switch variation {
	case "6E", "VI DCI", "DCI", "US":
		switch cardName {
		case "Crusade",
			"Lord of Atlantis",
			"Serra Avatar",
			"Thran Quarry",
			"Two-Headed Dragon":
			edition = "Junior Super Series"
		case "Forest",
			"Island",
			"Mountain",
			"Swamp",
			"Plains":
			edition = "Arena League 1999"
		case "Argothian Enchantress",
			"Balance",
			"Ball Lightning",
			"Gaea's Cradle",
			"Intuition",
			"Living Death",
			"Memory Lapse",
			"Oath of Druids",
			"Stroke of Genius",
			"Tradewind Rider":
			variation = "Judge"
		case "Vampiric Tutor":
			variation = "Judge 2000"
		case "Arc Lightning",
			"Chill",
			"Duress",
			"Enlightened Tutor",
			"Karn, Silver Golem",
			"Mana Leak",
			"Rewind",
			"Serum Visions",
			"Stupor":
			variation = "Arena"
		case "Powder Keg",
			"Voidmage Prodigy",
			"Wasteland":
			variation = "Magic Player Rewards"
		case "Zoetic Cavern":
			variation = "Gateway"
		default:
			variation = "FNM"
		}
	case "DotP 2012 - Xbox", "X Box Promo 2013", "X Box Promo",
		"Playstation", "PS3 Promo",
		"Duels of the Planeswalkers - PC",
		"Duel of the Planeswalkers":
		variation = "Duels"
	}

	switch {
	case strings.HasPrefix(variation, "The "):
		edition = "Magic Premiere Shop 2005"
	case strings.Contains(variation, "Holiday"):
		edition = "Happy Holidays"
	case mtgmatcher.HasPrefix(cardName, "Teferi Master of Time"):
		if edition == "Promotional Cards" {
			variation += "s"
		} else if edition == "Promo Pack" {
			variation += "p"
		}
	}

	if edition == "Promotional Cards" && variation == "" {
		ed, found = promo2setTable[cardName]
		if found {
			edition = ed
		}
	} else if edition == "Ikoria: Lair of Behemoths" && strings.Contains(cardName, " - ") {
		s := strings.Split(cardName, " - ")
		cardName = s[1]
		variation = "Godzilla"
		if strings.Contains(notes, "Japanese") {
			variation += " Japanese"
		}
	} else if edition == "Eternal Masters" && strings.Contains(notes, "Japanese") {
		return nil, errors.New("non-english card")
	}

	// Second pass in case some tags interfered with the lookup
	cn, found = cardTable[cardName]
	if found {
		cardName = cn
	}

	if variation == "Extemded Art" {
		variation = "Extended Art"
	}

	// Set finish
	isFoil := strings.Contains(notes, "Foil")
	if strings.Contains(notes, "Etched") {
		if variation != "" {
			variation += " "
		}
		variation += "Etched"
	}

	return &mtgmatcher.Card{
		Name:      cardName,
		Variation: variation,
		Edition:   edition,
		Foil:      isFoil,
	}, nil
}
