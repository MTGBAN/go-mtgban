package mtgmatcher

import (
	"strconv"
	"strings"

	"github.com/kodabb/go-mtgban/mtgmatcher/mtgjson"
)

func Match(inCard *Card) (cardId string, err error) {
	if backend.Sets == nil {
		return "", ErrDatastoreEmpty
	}

	// Look up by uuid
	if inCard.Id != "" {
		co, found := backend.UUIDs[inCard.Id]
		if found {
			return output(co.Card, inCard.Foil), nil
		}
	}

	// Get the card basic info to retrieve the Printings array
	entry, found := backend.Cards[Normalize(inCard.Name)]
	if !found {
		// Fixup up the name and try again
		adjustName(inCard)

		entry, found = backend.Cards[Normalize(inCard.Name)]
		if !found {
			return "", ErrCardDoesNotExist
		}
	}

	// Restore the card to the canonical MTGJSON name
	inCard.Name = entry.Name

	// Fix up edition
	adjustEdition(inCard)

	logger.Println("Processing", inCard, entry.Printings)

	// If there are multiple printings of the card, filter out to the
	// minimum common elements, using the rules defined.
	printings := entry.Printings
	if len(printings) > 1 {
		printings = filterPrintings(inCard, printings)
		logger.Println("Filtered printings:", printings)

		// Filtering was too aggressive or wrong data fed,
		// in either case, nothing else to be done here.
		if len(printings) == 0 {
			return "", ErrCardNotInEdition
		}
	}

	// This map will contain the setCode and an array of possible matches for
	// each edition.
	cardSet := map[string][]mtgjson.Card{}

	// Only one printing, it *has* to be it
	if len(printings) == 1 {
		cardSet[printings[0]] = matchInSet(inCard, printings[0])
	} else {
		// If multiple printing, try filtering to the closest name
		// described by the inCard.Edition.
		logger.Println("Several printings found, iterating over edition name")

		// First loop, search for a perfect match
		for _, setCode := range printings {
			// Perfect match, the card *has* to be present in the set
			if Equals(backend.Sets[setCode].Name, inCard.Edition) {
				logger.Println("Found a perfect match with", inCard.Edition, setCode)
				cardSet[setCode] = matchInSet(inCard, setCode)
			}
		}

		// Second loop, hope that a portion of the edition is in the set Name
		// This may result in false positives under certain circumstances.
		if len(cardSet) == 0 {
			logger.Println("No perfect match found, trying with heuristics")
			for _, setCode := range printings {
				set := backend.Sets[setCode]
				if Contains(set.Name, inCard.Edition) ||
					(inCard.isGenericPromo() && strings.HasSuffix(set.Name, "Promos")) {
					logger.Println("Found a possible match with", inCard.Edition, setCode)
					cardSet[setCode] = matchInSet(inCard, setCode)
				}
			}
		}

		// Third loop, YOLO
		// Let's consider every edition and hope the second pass will filter
		// duplicates out. This may result in false positives of course.
		if len(cardSet) == 0 {
			logger.Println("No loose match found, trying all")
			for _, setCode := range printings {
				cardSet[setCode] = matchInSet(inCard, setCode)
			}
		}
	}

	// Determine if any deduplication needs to be performed
	logger.Println("Found these possible matches")
	single := len(cardSet) == 1
	for _, dupCards := range cardSet {
		single = single && len(dupCards) == 1
		for _, card := range dupCards {
			logger.Println(card.SetCode, card.Name, card.Number)
		}
	}

	// Use the result as-is if it comes from a single card in a single set
	var outCards []mtgjson.Card
	if single {
		logger.Println("Single printing, using it right away")
		for _, outCards = range cardSet {
		}
	} else {
		// Otherwise do a second pass filter, using all inCard details
		logger.Println("Now filtering...")
		outCards = filterCards(inCard, cardSet)

		for _, card := range outCards {
			logger.Println(card.SetCode, card.Name, card.Number)
		}
	}

	// Just keep the first card found for gold-bordered sets
	if len(outCards) > 1 {
		if inCard.isWorldChamp() {
			logger.Println("Dropping a few extra entries...")
			logger.Println(outCards[1:])
			outCards = []mtgjson.Card{outCards[0]}
		}
	}

	// Finish line
	switch len(outCards) {
	// Not found, rip
	case 0:
		logger.Println("No matches...")
		err = ErrCardWrongVariant
	// Victory
	case 1:
		logger.Println("Found it!")
		cardId = output(outCards[0], inCard.Foil)
	// FOR SHAME
	default:
		logger.Println("Aliasing...")
		alias := newAliasingError()
		for i := range outCards {
			alias.dupes = append(alias.dupes, output(outCards[i], inCard.Foil))
		}
		err = alias
	}

	return
}

// Return an array of mtgjson.Card containing all the cards with the exact
// same name as the input inCard in the given mtgjson.Set.
func matchInSet(inCard *Card, setCode string) (outCards []mtgjson.Card) {
	set := backend.Sets[setCode]
	for _, card := range set.Cards {
		if inCard.Name == card.Name {
			outCards = append(outCards, card)
		}
	}
	return
}

// Try to fixup the name of the card or move extra varitions to the
// variant attribute. This should only be used in case the card name
// was not found.
func adjustName(inCard *Card) {
	// Move the card number from name to variation
	num := ExtractNumber(inCard.Name)
	if num != "" {
		fields := strings.Fields(inCard.Name)
		for i, field := range fields {
			if strings.Contains(field, num) {
				fields = append(fields[:i], fields[i+1:]...)
				break
			}
		}
		inCard.Name = strings.Join(fields, " ")
		inCard.addToVariant(num)
		return
	}

	// Move any single letter variation from name to beginning variation
	if inCard.IsBasicLand() {
		fields := strings.Fields(inCard.Name)
		if len(fields) > 1 {
			_, err := strconv.Atoi(strings.TrimPrefix(fields[1], "0"))
			isNum := err == nil
			isLetter := len(fields[1]) == 1

			if isNum || isLetter {
				oldVariation := inCard.Variation
				cuts := Cut(inCard.Name, " "+fields[1])

				inCard.Name = cuts[0]
				inCard.Variation = cuts[1]
				if oldVariation != "" {
					inCard.Variation += " " + oldVariation
				}
				return
			}
		}
	}

	// Check if the input name is the reskinned one
	// Currently appearing in IKO and some promo sets (PLGS and IKO BaB)
	if strings.Contains(inCard.Edition, "Ikoria") ||
		strings.Contains(inCard.Edition, "Promos") {
		for _, card := range backend.Sets["IKO"].Cards {
			if Equals(inCard.Name, card.FlavorName) {
				inCard.Name = card.Name
				inCard.addToVariant("Godzilla")
				return
			}
		}
	}

	// Special case for Un-sets that sometimes drop the parenthesis
	if strings.Contains(inCard.Edition, "Unglued") ||
		strings.Contains(inCard.Edition, "Unhinged") ||
		strings.Contains(inCard.Edition, "Unstable") ||
		strings.Contains(inCard.Edition, "Unsanctioned") {
		if HasPrefix(inCard.Name, "Our Market Research") {
			inCard.Name = LongestCardEver
			return
		}
		if HasPrefix(inCard.Name, "The Ultimate Nightmare") {
			inCard.Name = NightmareCard
			return
		}
		if Contains(inCard.Name, "Surgeon") && Contains(inCard.Name, "Commander") {
			inCard.Name = "Surgeon ~General~ Commander"
			return
		}

		for cardName, props := range backend.Cards {
			if HasPrefix(cardName, inCard.Name) {
				inCard.Name = props.Name
				return
			}
		}
	}

	// Altenatively try checking across any prefix, as long as it's a double
	// sided card, for some particular cases, like meld cards, or Treasure Chest
	for cardName, props := range backend.Cards {
		if props.Layout != "normal" && HasPrefix(cardName, inCard.Name) {
			inCard.Name = props.Name
			return
		}
	}
}

// Try to fixup the edition and variant of the card, using well-known variantions,
// or use edition/variant attributes to determine the correct edition/variant combo,
// or look up known cards in small sets.
func adjustEdition(inCard *Card) {
	edition := inCard.Edition
	variation := inCard.Variation

	// Need to decouple The List and Mystery booster first or it will confuse
	// later matching. For an uptodate list of aliased cards visit this link:
	// https://scryfall.com/search?q=in%3Aplist+%28in%3Amb1+or+in%3Afmb1%29+%28e%3Amb1+or+e%3Aplist+or+e%3Afmb1%29&unique=prints&as=grid&order=name
	// Skip only if the edition or variation are explictly set as The List
	if edition != "The List" && variation != "The List" &&
		(Contains(edition, "Mystery Booster") || Contains(edition, "The List") ||
			Contains(variation, "Mystery Booster") || Contains(variation, "The List")) {
		if inCard.Foil || (Contains(edition, "Foil") && !Contains(edition, "Non")) || (Contains(variation, "Foil") && !Contains(variation, "Non")) {
			edition = "FMB1"
		} else if Contains(edition, "Test") || Contains(variation, "Test") {
			edition = "CMB1"
		} else {
			// Check if card is is only one of these two sets
			mb1s := matchInSet(inCard, "MB1")
			plists := matchInSet(inCard, "PLIST")
			if len(mb1s) == 1 && len(plists) == 0 {
				edition = "MB1"
			} else if len(mb1s) == 0 && len(plists) == 1 {
				edition = "PLIST"
			} else if len(mb1s) == 1 && len(plists) == 1 {
				switch variation {
				// If it has one of these special treatments it's PLIST definitely
				case "Player Rewards",
					"MagicFest",
					"Commander",
					"Extended Art",
					"Signature Spellbook: Jace",
					"The List Textless",
					"Player Rewards Promo",
					"RNA MagicFest Promo",
					"State Champs Promo":
					edition = "PLIST"
				default:
					// Otherwise it's probably MB1, including the indistinguishable
					// ones, unless variation has additional information
					edition = "MB1"

					// Adjust variation to get a correct edition name
					ed, found := EditionTable[variation]
					if found {
						variation = ed
					}

					// Check if the card name has the appropriate variation that
					// lets us determine it's from PLIST
					if AliasedPLISTTable[inCard.Name][variation] {
						edition = "PLIST"
					}
				}
			}
		}
		// Ignore this, we have all we need now
		variation = ""
	}

	set, found := backend.Sets[edition]
	if found {
		edition = set.Name
	}
	ed, found := EditionTable[edition]
	if found {
		edition = ed
	}
	ed, found = EditionTable[variation]
	if found {
		edition = ed
	}

	// Adjust box set
	switch {
	case Equals(edition, "Double Masters Box Toppers"),
		Equals(edition, "Double Masters: Extras"),
		Equals(edition, "Double Masters: Variants"):
		edition = "Double Masters"
		if !inCard.isBasicLand() {
			variation = "Borderless"
		}
	case strings.Contains(edition, "Mythic Edition"),
		strings.Contains(inCard.Variation, "Mythic Edition"):
		edition = "Mythic Edition"
	case strings.Contains(edition, "Invocations") ||
		((edition == "Hour of Devastation" || edition == "Amonkhet") &&
			strings.Contains(inCard.Variation, "Invocation")):
		edition = "Amonkhet Invocations"
	case strings.Contains(edition, "Inventions"):
		edition = "Kaladesh Inventions"
	case strings.Contains(edition, "Expeditions") && !strings.Contains(edition, "Rising"):
		edition = "Zendikar Expeditions"
	case strings.Contains(edition, "Expeditions") && strings.Contains(edition, "Rising"):
		edition = "Zendikar Rising Expeditions"
	default:
		for _, tag := range []string{
			"(Collector Edition)", "Collectors", "Extras", "Variants",
		} {
			if strings.HasSuffix(edition, tag) {
				edition = strings.TrimSuffix(edition, tag)
				edition = strings.TrimSpace(edition)
				edition = strings.TrimSuffix(edition, ":")
			}
		}
	}

	switch {
	case strings.Contains(edition, "Commander"):
		edition = inCard.commanderEdition()
	case strings.Contains(variation, "Ravnica Weekend") || strings.Contains(edition, "Weekend"):
		edition, variation = inCard.ravnicaWeekend()
	case inCard.Contains("Guild Kit"):
		edition = inCard.ravnicaGuidKit()
	case strings.Contains(variation, "APAC Set") || strings.Contains(variation, "Euro Set"):
		num := ExtractNumber(variation)
		if num != "" {
			variation = strings.Replace(variation, num+" ", "", 1)
		}
	case strings.HasPrefix(variation, "Junior") && strings.Contains(variation, "APAC"),
		strings.HasPrefix(variation, "Junior APAC Series U"):
		edition = "Junior APAC Series"
	case strings.HasPrefix(variation, "Junior Super Series"),
		strings.HasPrefix(variation, "MSS Foil"),
		strings.HasPrefix(variation, "MSS #J"),
		strings.HasPrefix(variation, "MSS Promo J"),
		strings.HasPrefix(variation, "JSS #J"),
		strings.Contains(variation, "JSS Foil") && !Contains(variation, "euro"):
		edition = "Junior Super Series"
	case strings.HasPrefix(variation, "Junior Series Europe"),
		strings.HasPrefix(variation, "Junior Series E"),
		strings.HasPrefix(variation, "Junior Series #E"),
		strings.HasPrefix(variation, "Junior Series Promo E"),
		strings.HasPrefix(variation, "Junior Series Promo Foil E"),
		strings.HasPrefix(variation, "ESS Foil E"),
		strings.HasPrefix(variation, "European JrS E"),
		strings.HasPrefix(variation, "European JSS Foil E"):
		edition = "Junior Series Europe"
	case Equals(inCard.Name, "Teferi, Master of Time"):
		num := ExtractNumber(variation)
		_, err := strconv.Atoi(num)
		if err == nil {
			if inCard.isPrerelease() {
				variation = num + "s"
			} else if inCard.isPromoPack() {
				variation = num + "p"
			}
		}
		if num == "" {
			if inCard.isPrerelease() {
				variation = "75s"
			} else if inCard.isPromoPack() {
				variation = "75p"
			} else if inCard.isBorderless() {
				variation = "281"
			} else if inCard.isShowcase() {
				variation = "290"
			} else {
				variation = "75"
			}
		}
	}
	inCard.Edition = edition
	inCard.Variation = variation

	// Special handling since so many providers get this wrong
	switch {
	// XLN Treasure Chest
	case inCard.isBaB() && len(matchInSet(inCard, "PXTC")) != 0:
		inCard.Edition = backend.Sets["PXTC"].Name
	// BFZ Standard Series
	case inCard.isGenericAltArt() && len(matchInSet(inCard, "PSS1")) != 0:
		inCard.Edition = backend.Sets["PSS1"].Name
	// Champs and States
	case inCard.isGenericExtendedArt() && len(matchInSet(inCard, "PCMP")) != 0:
		inCard.Edition = backend.Sets["PCMP"].Name
	// Portal Demo Game
	case inCard.isPortalAlt() && len(matchInSet(inCard, "PPOD")) != 0:
		inCard.Edition = backend.Sets["PPOD"].Name
	// Secret Lair Ultimate
	case strings.Contains(inCard.Edition, "Secret Lair") &&
		len(matchInSet(inCard, "SLU")) != 0:
		inCard.Edition = backend.Sets["SLU"].Name
	// Summer of Magic
	case (inCard.isWPNGateway() || strings.Contains(inCard.Variation, "Summer")) &&
		len(matchInSet(inCard, "PSUM")) != 0:
		inCard.Edition = backend.Sets["PSUM"].Name

	// Single card mismatches
	case Equals(inCard.Name, "Rhox") && inCard.isGenericAltArt():
		inCard.Edition = "Starter 2000"
	case Equals(inCard.Name, "Balduvian Horde") && (strings.Contains(inCard.Variation, "Judge") || strings.Contains(inCard.Edition, "Promo") || inCard.Contains("DCI")):
		inCard.Edition = "World Championship Promos"
	case Equals(inCard.Name, "Nalathni Dragon") && inCard.isIDWMagazineBook():
		inCard.Edition = "Dragon Con"
	case Equals(inCard.Name, "Ass Whuppin'") && inCard.isPrerelease():
		inCard.Edition = "Release Events"
	case Equals(inCard.Name, "Celestine Reef") && inCard.isPrerelease():
		inCard.Edition = "Promotional Planes"
	case Equals(inCard.Name, "Ajani Vengeant") && inCard.isRelease():
		inCard.Variation = "Prerelease"
	case Equals(inCard.Name, "Tamiyo's Journal") && inCard.Variation == "" && inCard.Foil:
		inCard.Variation = "Foil"
	case Equals(inCard.Name, "Underworld Dreams") && inCard.Contains("DCI"):
		inCard.Edition = "Two-Headed Giant Tournament"
	case Equals(inCard.Name, "Jace Beleren") && inCard.Contains("DCI"):
		inCard.Edition = "Miscellaneous Book Promos"
	case Equals(inCard.Name, "Serra Angel") && inCard.Contains("DCI"):
		inCard.Edition = "Wizards of the Coast Online Store"

	case Equals(inCard.Name, "Incinerate") && inCard.Contains("DCI"):
		inCard.Edition = "DCI Legend Membership"
	case Equals(inCard.Name, "Counterspell") && inCard.Contains("DCI"):
		inCard.Edition = "DCI Legend Membership"

	case Equals(inCard.Name, "Kamahl, Pit Fighter") && inCard.Contains("DCI"):
		inCard.Edition = "15th Anniversary Cards"
	case Equals(inCard.Name, "Char") && inCard.Contains("DCI"):
		inCard.Edition = "15th Anniversary Cards"

	case Equals(inCard.Name, "Sigarda, Host of Herons") && inCard.isPrerelease():
		inCard.Edition = "Open the Helvault"
	case Equals(inCard.Name, "Griselbrand") && inCard.isPrerelease():
		inCard.Edition = "Open the Helvault"
	case Equals(inCard.Name, "Gisela, Blade of Goldnight") && inCard.isPrerelease():
		inCard.Edition = "Open the Helvault"
	case Equals(inCard.Name, "Bruna, Light of Alabaster") && inCard.isPrerelease():
		inCard.Edition = "Open the Helvault"
	case Equals(inCard.Name, "Avacyn, Angel of Hope") && inCard.isPrerelease():
		inCard.Edition = "Open the Helvault"
	}
}
