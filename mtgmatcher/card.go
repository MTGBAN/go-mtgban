package mtgmatcher

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/mtgban/go-mtgban/mtgmatcher/mtgjson"
	"golang.org/x/exp/slices"
)

// Card is a generic card representation using fields defined by the MTGJSON project.
type InputCard struct {
	// The mtgjson unique identifier of the card
	// When used as input it can host mtgjson or scryfall id
	Id string `json:"id,omitempty"`

	// The canonical name of the card
	Name string `json:"name,omitempty"`

	// The hint or commonly know variation
	Variation string `json:"variant,omitempty"`

	// The set the card comes from, or a portion of it
	Edition string `json:"edition,omitempty"`

	// Whether the card is foil or not
	Foil bool `json:"foil,omitempty"`

	// The card belongs to the extended side of the set, usually containing
	// variants with the same name of existing cards in the same set, but with
	// different frames or border effects
	beyondBaseSet bool

	// In case edition information is not accurate, use this flag to
	// perform a best-effor search, which will try to isolate promo
	// printings from the others
	promoWildcard bool

	// In case card got renamed in some way, this contains the original
	// card name, instead of the sanitized version
	originalName string

	// The language as parsed
	Language string `json:"language,omitempty"`
}

// Card implements the Stringer interface
func (c *InputCard) String() string {
	name := c.Name
	edition := c.Edition

	if name == "" {
		co, err := GetUUID(c.Id)
		if err == nil {
			name = co.Name
			edition = co.Edition
		}
	}

	if c.Variation != "" {
		name = fmt.Sprintf("%s ('%s')", name, c.Variation)
	}
	finish := ""
	if c.isEtched() {
		finish = " (etched)"
	} else if c.Foil {
		finish = " (foil)"
	}
	lang := ""
	if c.Language != "" && c.Language != "English" {
		lang = " {" + c.Language + "}"
	}
	return fmt.Sprintf("%s [%s%s]%s", name, edition, finish, lang)
}

func output(card Card, flags ...bool) string {
	hasNonfoil := card.HasFinish(mtgjson.FinishNonfoil)
	hasFoil := card.HasFinish(mtgjson.FinishFoil)
	hasEtched := card.HasFinish(mtgjson.FinishEtched)

	etched := len(flags) > 1 && flags[1]
	foil := len(flags) > 0 && flags[0] && !etched

	// In case the foiling information is incorrect
	if !foil && !hasNonfoil && !hasEtched {
		foil = true
	} else if foil && !hasFoil {
		foil = false
	}
	if hasFoil && !hasNonfoil && !hasEtched {
		foil = true
	} else if !hasFoil && (hasNonfoil || hasEtched) {
		foil = false
	}

	// In case the etching information is incorrect
	if !etched && !hasNonfoil && !hasFoil {
		etched = true
	} else if etched && !hasEtched {
		etched = false
	}
	if hasEtched && !hasNonfoil && !hasFoil {
		etched = true
	} else if !hasEtched && (hasNonfoil || hasFoil) {
		etched = false
	}

	// Prepare the output card
	id := card.UUID
	// Append suffixes to the Id to distinguish cards among finishes
	if etched && (hasNonfoil || hasFoil) {
		id += suffixEtched
	} else if foil && hasNonfoil {
		id += suffixFoil
	}

	return id
}

func (c *InputCard) addToVariant(tag string) {
	if c.Variation != "" {
		c.Variation += " "
	}
	c.Variation += tag
}

// Returns whether the input string may represent a token
func IsToken(name string) bool {
	// Check main table first
	if slices.Contains(backend.Tokens, name) {
		return true
	}
	switch name {
	// Custom token names
	case "A Threat to Alara: Nicol Bolas",
		"Fun Format: Pack Wars",
		"On An Adventure",
		"Pyromantic Pixels",
		"Theme: The Gold Standard",
		"Theme: WUBRG Cards":
		return true
	// WCD extra cards
	case "Biography",
		"Blank",
		"Overview":
		return true
	}
	switch {
	// Avoid confusion with Monarch and Emblem below
	case HasPrefix(name, "Emblem of the Warmind"),
		HasPrefix(name, "Kavu Monarch"),
		HasPrefix(name, "Leering Emblem"),
		// and with the `card` wildcard
		HasPrefix(name, "Our Market Research"):
		return false
	// Anything token
	case strings.Contains(name, " Card"),
		strings.Contains(name, "Card "),
		strings.HasPrefix(name, "Bounty"),
		Contains(name, "Arena Code"),
		Contains(name, "Art Series"),
		Contains(name, "Charlie Brown"),
		Contains(name, "Checklist"),
		Contains(name, "Copy"),
		Contains(name, "Decklist"),
		Contains(name, "DFC Helper"),
		Contains(name, "Dungeon of the Mad Mage"),
		Contains(name, "Emblem"),
		Contains(name, "Experience C"),
		Contains(name, "Giant Teddy Bear"),
		Contains(name, "Guild Symbol"),
		Contains(name, "Magic Minigame"),
		Contains(name, "The Monarch"),
		strings.Contains(name, "The Initiative"),
		Contains(name, "Morph Overlay"),
		Contains(name, "On Your Turn"),
		Contains(name, "Online Code"),
		Contains(name, "Oversize"),
		Contains(name, "Punch Out"),
		Contains(name, "Token"),
		Contains(name, "Rules Tip"):
		return true
	// Alternative rules tip card names found on mkm
	case strings.HasPrefix(name, "Build a Deck: "),
		strings.HasPrefix(name, "Tip: "):
		return true
	}

	return false
}

func (c *InputCard) isUnsupported() bool {
	return c.Contains("Art Series") ||
		strings.HasSuffix(c.Edition, "Art Variants") || // toa
		(c.Contains("Art Card") && !c.Contains("Chinese")) || // Art Series, except a well-known edition
		c.Contains("Complete") || // a complete collection
		c.Contains("Fallen Empires: Wyvern Misprints") ||
		c.Contains("Ultra-Pro Puzzle") ||
		c.Contains("Player Cards") || // scg pro players
		c.Contains("Foreign White Border") || // for REV and 4ED
		c.Contains("Filler Cards") || // Misprints from mkm and ct
		c.Contains("Salvat") || // Salvat-Hachette 2005/2011
		c.Contains("Redemption Program") || // PRES
		c.Contains("Heroes of the Realm") || // HTR*
		c.Contains("Memorabilia") ||
		c.Contains("Front Card") || // Jumpstart
		(c.Contains("Duel Masters") && c.Contains("Not Tournament Legal")) || // scg
		c.Contains("Sealed") ||
		c.Contains("Un-Known Event Playtest") ||
		c.Contains("Charlie Brown") || // abu
		// Oversized are usually ok, but 8th and 9th ed box topper variants
		// conflict with the actual edition name, so skip them
		(c.Contains("Oversize") && (c.Contains("8th") || c.Contains("9th")))
}

func (c *InputCard) isSpecificUnsupported() bool {
	switch c.Name {
	case "Spined Wurm":
		return Contains(c.Edition, "Starter 2000")
	case "Drudge Skeletons",
		"Emerald Medallion",
		"Forest",
		"Sapphire Medallion",
		"Serra Angel",
		"Time Elemental",
		"Winged Sliver":
		return c.Contains("Misprint")
	// Erraneous release information
	case "Zombify":
		return c.Contains("Game Night")
	}
	return false
}

// Returns whether the input string may represent a basic land
func IsBasicLand(name string) bool {
	switch {
	case strings.Contains(name, "Bear") && !strings.Contains(name, "Beard"), // G
		strings.Contains(name, "Mosquito"),                                     // B
		strings.Contains(name, "Stronghold"), strings.Contains(name, "Bandit"), // R
		strings.Contains(name, "Yeti"), strings.Contains(name, "Titan"), // R
		strings.Contains(name, "Valley"), strings.Contains(name, "Goat"), // R
		strings.Contains(name, "Fish"), strings.Contains(name, "Sanctuary"), // U
		strings.Contains(name, "Wak-Wak"): // U
	case strings.HasPrefix(name, "Plains"),
		strings.HasPrefix(name, "Island"),
		strings.HasPrefix(name, "Swamp"),
		strings.HasPrefix(name, "Mountain"),
		strings.HasPrefix(name, "Forest"),
		strings.HasPrefix(name, "Wastes"):
		return true
	case HasPrefix(name, "Snow-Covered"):
		return true
	}
	return false
}

// Returns whether the card is a basic land
func (c *InputCard) IsBasicLand() bool {
	return IsBasicLand(c.Name)
}

// More specific version of the above, for internal use only
func (c *InputCard) isBasicLand() bool {
	switch c.Name {
	case "Plains", "Island", "Swamp", "Mountain", "Forest", "Wastes":
		return true
	}
	return false
}

// Returns whether the cards is a "generic" promo, that probably needs
// further analysis to be fully categorized. Tokens are excluded.
func (c *InputCard) isGenericPromo() bool {
	return !c.isBaB() && !c.isPromoPack() && !c.isPrerelease() && !c.isSDCC() &&
		!c.isRetro() &&
		!c.Contains("Year of the") && // tcg
		!c.Contains("Deckmasters") && // no real promos here, just foils
		!c.Contains("Token") && !IsToken(c.Name) &&
		(Contains(c.Variation, "Promo") || // catch-all (*not* Edition)
			c.Contains("Gift Box") || // ck+scg
			(c.Contains("Promo") && c.Contains("Intro Pack")) || // scg
			c.Contains("League") ||
			c.Contains("Play Draft") || // scg
			c.Contains("Miscellaneous") ||
			c.Contains("Open House") || // tcg
			(c.Contains("Other") && !c.Contains("Brother")) ||
			c.Contains("Planeswalker Event") || // tcg
			c.Contains("Planeswalker Weekend") || // scg
			c.Contains("Store Challenge") || // scg
			c.Contains("Unique")) // mtgs
}

func (c *InputCard) isDCIPromo() bool {
	return c.Contains("DCI") && !c.Contains("Judge")
}

func (c *InputCard) isGenericAltArt() bool {
	// "Alt" includes Alternative
	return c.Contains("Alt") && c.Contains("Art")
}

func (c *InputCard) isGenericExtendedArt() bool {
	return Contains(c.Variation, "Art") &&
		(Contains(c.Variation, "Extended") ||
			Contains(c.Variation, "Full"))
}

func (c *InputCard) isPrerelease() bool {
	return c.Contains("Prerelease") ||
		c.Contains("Preview") // scg
}

func (c *InputCard) isPromoPack() bool {
	return c.Contains("Promo Pack") ||
		c.Variation == "Dark Frame Promo" ||
		Contains(c.Variation, "Planeswalker Stamp") ||
		Contains(c.Variation, "Silver Stamped") ||
		(strings.HasSuffix(ExtractNumber(c.Variation), "p") && !c.Contains("30th"))
}

func (c *InputCard) isBorderless() bool {
	return Contains(c.Variation, "Borderless")
}

func (c *InputCard) isExtendedArt() bool {
	return Contains(c.Variation, "Extended")
}

func (c *InputCard) isShowcase() bool {
	return Contains(c.Variation, "Showcase") ||
		Contains(c.Variation, "Sketch") // binderpos
}

func (c *InputCard) isReskin() bool {
	return (Contains(c.Variation, "Reskin") ||
		Contains(c.Variation, "Dracula") ||
		Contains(c.Variation, "Godzilla")) &&
		// Needed to distinguish the SLD godizlla lands
		!c.isBasicLand()
}

func (c *InputCard) isStepAndCompleat() bool {
	return Contains(c.Variation, "Compleat")
}

func (c *InputCard) isOilSlick() bool {
	return strings.Contains(strings.ToLower(c.Variation), "slick") ||
		strings.Contains(strings.ToLower(c.Edition), "slick")
}

func (c *InputCard) isFNM() bool {
	return c.Contains("FNM") ||
		c.Contains("Friday Night Magic")
}

func (c *InputCard) isJPN() bool {
	return strings.Contains(c.Variation, "JPN") ||
		strings.Contains(c.Variation, "JP") ||
		c.Contains("Japanese") ||
		Contains(c.Variation, "Gotta") ||
		Contains(c.Variation, "Dengeki")
}

func (c *InputCard) isChineseAltArt() bool {
	return (c.Contains("Chinese") || strings.Contains(c.Variation, "CS")) && c.isGenericAltArt()
}

func (c *InputCard) isRelease() bool {
	return !c.Contains("Prerelease") &&
		(c.Contains("Release") ||
			c.Contains("Draft Weekend") ||
			c.Contains("Launch"))
}

func (c *InputCard) isWPNGateway() bool {
	return c.Contains("WPN") ||
		c.Contains("Gateway") ||
		Contains(c.Variation, "Wizards Play Network") ||
		Contains(c.Variation, "Commander Party") || // scg
		Contains(c.Variation, "Moonlit Lands") // ck
}

func (c *InputCard) isIDWMagazineBook() bool {
	return strings.HasPrefix(c.Variation, "IDW") || strings.HasPrefix(c.Edition, "IDW") ||
		c.Contains("Magazine") ||
		c.Contains("Duelist") ||
		// Catches Comic and Comics, but skips San Diego Comic-Con
		(c.Contains("Comic") && !c.Contains("Diego")) ||
		// Cannot use Contains because it may trigger a false positive
		// for cards with "book" in their variation (insidious bookworms)
		c.Variation == "Book" ||
		c.Variation == "Insert" || // mmc
		c.Variation == "Japanese Promo" || // tcg
		c.Contains("Book Insert") ||
		c.Contains("Walmart") ||
		c.Contains("Coro Coro") || // stks
		c.Contains("Graphic Novel") || // stks
		strings.Contains(c.Variation, "Book Promo") || // sz
		c.Contains("Top Deck") || // csi
		c.Contains("Hobby Japan") || // abu+tcg
		Contains(c.Edition, "CardZ") || // mkm
		Contains(c.Edition, "Dengeki") || // mkm
		c.Variation == "Insert Foil" || // ck
		c.Contains("Beadle & Grimm Phyrexian") || // scg
		c.Contains("Stance Socks") || // scg
		c.Contains("Manga Promo") || // csi
		c.Contains("Media Promo") || // tcg
		c.Contains("Media Insert") // mm+nf
}

func (c *InputCard) isResale() bool {
	return !c.Contains("Championship") && (c.Contains("Repack") || c.Contains("Store") || c.Contains("Resale"))
}

func (c *InputCard) isJudge() bool {
	return c.Contains("Judge")
}

func (c *InputCard) isRewards() bool {
	return (Contains(c.Variation, "Textless") &&
		!Contains(c.Variation, "Year of") &&
		!Contains(c.Variation, "Lunar") &&
		!Contains(c.Variation, "Store")) ||
		(c.Contains("Reward") && !c.isJudge())
}

func (c *InputCard) isMagicFest() bool {
	return c.Contains("Magic Fest") ||
		c.Contains("MagicCon") || // scg
		strings.Contains(c.Edition, "MFP") || // tcg collection
		strings.Contains(c.Variation, "MFP") // tcg collection
}

func (c *InputCard) isBaB() bool {
	return c.Contains("Buy a Box") ||
		strings.Contains(c.Variation, "BABP") || // tcg collection
		strings.Contains(c.Variation, "BIBB") || // sz
		(c.Contains("Box Promos") && // ha+sz
			!c.Contains("Xbox") && // ck+abu
			!c.Contains("Gift")) // csi
}

func (c *InputCard) isBundle() bool {
	return c.Contains("Bundle")
}

func (c *InputCard) isFoil() bool {
	return Contains(c.Variation, "Foil") && !Contains(c.Variation, "Non") && !c.isEtched()
}

func (c *InputCard) isEtched() bool {
	// Note this can't be just "etch" because it would catch the "sketch" cards
	return Contains(c.Variation, "Etched")
}

func (c *InputCard) isARNLightMana() bool {
	return Contains(c.Variation, "light") || strings.Contains(c.Variation, "†")
}

func (c *InputCard) isARNDarkMana() bool {
	return Contains(c.Variation, "dark")
}

func (c *InputCard) isArena() bool {
	return c.Contains("Arena")
}

func (c *InputCard) isSDCC() bool {
	return c.Contains("SDCC") ||
		c.Contains("San Diego Comic-Con")
}

func (c *InputCard) isRetro() bool {
	return c.Contains("Retro")
}

func (c *InputCard) playerRewardsYear(maybeYear string) string {
	if maybeYear == "" {
		switch c.Name {
		case "Bear":
			if c.Variation == "Odyssey" {
				maybeYear = "2001"
			} else if c.Variation == "Onslaught" {
				maybeYear = "2003"
			}
		case "Beast":
			if c.Variation == "Odyssey" {
				maybeYear = "2001"
			} else if c.Variation == "Darksteel" {
				maybeYear = "2004"
			}
		case "Elephant":
			if c.Variation == "Invasion" {
				maybeYear = "2001"
			} else if c.Variation == "Odyssey" {
				maybeYear = "2002"
			}
		case "Spirit":
			if c.Variation == "Planeshift" {
				maybeYear = "2001"
			} else if c.Variation == "Champions" {
				maybeYear = "2004"
			}
		case "Lightning Bolt":
			if c.Contains("Oversize") {
				maybeYear = "2009"
			} else {
				maybeYear = "2010"
			}
		}
	}
	return maybeYear
}

func (c *InputCard) arenaYear(maybeYear string) string {
	if maybeYear == "" {
		switch {
		case strings.Contains(c.Variation, "Tony Roberts"):
			maybeYear = "1996"
		case strings.Contains(c.Variation, "Urza"),
			strings.Contains(c.Variation, "Saga"),
			strings.Contains(c.Variation, "Anthony S. Waters"),
			strings.Contains(c.Variation, "Donato Giancola"):
			maybeYear = "1999"
		case strings.Contains(c.Variation, "Mercadian"),
			strings.Contains(c.Variation, "Masques"):
			maybeYear = "2000"
		case strings.Contains(c.Variation, "Ice Age"),
			strings.Contains(c.Variation, "IA"),
			strings.Contains(c.Variation, "Pat Morrissey"),
			strings.Contains(c.Variation, "Anson Maddocks"),
			strings.Contains(c.Variation, "Tom Wanerstrand"),
			strings.Contains(c.Variation, "Christopher Rush"),
			strings.Contains(c.Variation, "Douglas Shuler"):
			maybeYear = "2001"
		case strings.Contains(c.Variation, "Mark Poole"):
			maybeYear = "2002"
		case strings.Contains(c.Variation, "Rob Alexander"):
			maybeYear = "2003"
		case strings.Contains(c.Variation, "Don Thompson"):
			maybeYear = "2005"
		case strings.Contains(c.Variation, "Beta"):
			switch c.Name {
			case "Forest":
				maybeYear = "2001"
			case "Island":
				maybeYear = "2002"
			}
		}
	} else if c.Name == "Forest" && strings.Contains(maybeYear, "2002") {
		maybeYear = "2001"
	} else if c.Name == "Island" && strings.Contains(maybeYear, "2001") && strings.Contains(c.Variation, "Poole") {
		maybeYear = "2002"
	}
	return maybeYear
}

func (c *InputCard) isWorldChamp() bool {
	return Contains(c.Edition, "Pro Tour Collect") ||
		Contains(c.Edition, "Pro Tour 1996") ||
		Contains(c.Edition, "World Championship") ||
		Contains(c.Edition, "Championship Deck") ||
		Contains(c.Edition, "WCD")
}

func parseWorldChampPrefix(variation string) (string, bool) {
	players := map[string]string{
		"Aeo Paquette":         "ap",
		"Alex Borteh":          "ab",
		"Antoine Ruel":         "ar",
		"Ben Rubin":            "br",
		"Bertrand Lestree":     "bl",
		"Brian Hacker":         "bh",
		"Brian Kibler":         "bk",
		"Brian Selden":         "bs",
		"Brian Seldon":         "bs",
		"Carlos Romao":         "cr",
		"Daniel Zink":          "dz",
		"Dave Humpherys":       "dh",
		"Eric Tam":             "et",
		"Gabriel Nassif":       "gn",
		"George Baxter":        "gb",
		"Jakub Slemr":          "js",
		"Jan Tomcani":          "jt",
		"Janosch Kuhn":         "jk",
		"Janosch Kuehn":        "jk",
		"Jon Finkel":           "jf",
		"Julien Nuijten":       "jn",
		"Kai Budde":            "kb",
		"Leon Lindback":        "ll",
		"Manuel Bevand":        "mb",
		"Mark Justice":         "mj",
		"Mark Le Pine":         "mlp",
		"Matt Linde":           "ml",
		"Michael Locanto":      "ml",
		"Michael Loconto":      "ml",
		"Nicolas Labarre":      "nl",
		"Paul McCabe":          "pm",
		"Peer Kroger":          "pk",
		"Preston Poulter":      "pp",
		"Randy Buehler":        "rb",
		"Raphael Levy":         "rl",
		"Shawn Regnier":        "shr",
		"Shawn Hammer Regnier": "shr",
		"Sim Han How":          "shh",
		"Svend Geertsen":       "sg",
		"Tom van de Logt":      "tvdl",
		"Wolfgang Eder":        "we",
	}

	// We cannot use HasPrefix for the second check due to mlp/ml aliasing
	variation = strings.ToLower(variation)
	idx := strings.IndexFunc(variation, func(c rune) bool {
		return unicode.IsDigit(c)
	})
	// Iterate over the player list and check if their name or their initials are present
	for player, tag := range players {
		if Contains(variation, player) || (idx > -1 && variation[:idx] == tag) {
			sb := strings.Contains(variation, "sb") || strings.Contains(variation, "sideboard")
			return tag, sb
		}
	}
	return "", false
}

func (c *InputCard) worldChampPrefix() (string, bool) {
	prefix, sideboard := parseWorldChampPrefix(c.Variation)
	if prefix == "" {
		return parseWorldChampPrefix(c.Edition)
	}
	return prefix, sideboard
}

func (c *InputCard) isDuelsOfThePW() bool {
	// XXX: do not use c.Contains here
	return strings.Contains(c.Variation, "Duels") ||
		strings.Contains(c.Edition, "Duels") ||
		Contains(c.Variation, "DotP") // tat
}

func (c *InputCard) isBasicFullArt() bool {
	return c.isBasicLand() &&
		(Contains(c.Variation, "full art") ||
			c.Variation == "FA") && // csi
		!Contains(c.Variation, "non") &&
		!Contains(c.Variation, "not") // csi
}

func (c *InputCard) isBasicNonFullArt() bool {
	return c.isBasicLand() &&
		Contains(c.Variation, "non-full art") ||
		Contains(c.Variation, "Intro") || // abu
		Contains(c.Variation, "NOT the full art") // csi
}

func (c *InputCard) isPremiereShop() bool {
	return c.isBasicLand() &&
		// XXX: do not use c.Contains here
		(strings.Contains(c.Variation, "MPS") ||
			strings.Contains(c.Variation, "Premier") || // csi
			strings.Contains(c.Edition, "MPS") ||
			strings.Contains(c.Edition, "Premiere Shop")) // mkm
}

func (c *InputCard) isPortalAlt() bool {
	return (Contains(c.Variation, "Reminder Text") &&
		!Contains(c.Variation, "No")) ||
		Contains(c.Variation, "No Flavor Text") || // csi
		Contains(c.Variation, "Without Flavor Text") // csi
}

func (c *InputCard) isDuelDecks() bool {
	return ((c.Contains(" vs ")) ||
		(strings.Contains(c.Variation, " v. "))) && // tcg
		!c.Contains("Anthology")
}

func (c *InputCard) isDuelDecksAnthology() bool {
	return strings.Contains(c.Edition, "DDA") ||
		(Contains(c.Edition, "Duel Decks") && Contains(c.Edition, "Anthology"))
}

func (c *InputCard) duelDecksVariant() string {
	if !c.isDuelDecks() {
		return ""
	}

	// Variation might contain numbers, strip them away
	variant := c.Variation
	num := ExtractNumber(variant)
	variant = strings.TrimSpace(strings.Replace(variant, num, "", 1))
	if len(variant) < len("Duel Deck") {
		variant = c.Edition
	}

	if strings.Contains(variant, ": ") {
		fields := strings.Split(variant, ": ")
		variant = fields[len(fields)-1]
	}

	return variant
}

func (c *InputCard) isMysteryList() bool {
	return c.Contains("Mystery") || c.Contains("The List") ||
		c.Contains("Planeswalker Symbol Reprints")
}

func (c *InputCard) isSecretLair() bool {
	return c.Contains("Secret Lair") || strings.Contains(c.Edition, "SLD")
}

func (c *InputCard) hasSecretLairTag(code string) bool {
	var tag bool
	switch code {
	case "SLU":
		// SLU is mostly static and cards are unlikely to reappear elsewhere
		tag = c.Contains("Ultimate") || len(MatchInSet(c.Name, "SLU")) == 1
	case "SLX":
		// SLX only has plain cards, if they are reskinned, they are from SLD
		tag = !c.isReskin() || c.Contains("Within")
	case "SLC":
		// These cards are numbered after the year they represent
		tag = c.Contains("30th") || c.Contains("Countdown") || ExtractYear(c.Variation) != ""
	case "SLP":
		// Simple check the variations
		tag = c.Contains("Showdown") || c.Contains("Prize") || c.Contains("Finish") || c.Contains("Play")
	}

	return c.isSecretLair() && tag
}

func (c *InputCard) isThickDisplay() bool {
	return c.Contains("Display") || c.Contains("Thick")
}

func (c *InputCard) isPhyrexian() bool {
	return Contains(c.Variation, "Phyrexian")
}

func (c *InputCard) isGalaxyFoil() bool {
	return Contains(c.Variation, "Galaxy")
}

func (c *InputCard) isSurgeFoil() bool {
	return strings.Contains(strings.ToLower(c.Variation), "surge") ||
		strings.Contains(strings.ToLower(c.Edition), "surge")
}

func (c *InputCard) isSerialized() bool {
	return strings.Contains(strings.ToLower(c.Variation), "serial") ||
		strings.Contains(strings.ToLower(c.Edition), "serial")
}

func (c *InputCard) possibleNumberSuffix() string {
	fields := strings.Fields(c.Variation)
	for _, field := range fields {
		if len(field) == 1 && unicode.IsLetter(rune(field[0])) {
			return strings.ToLower(field)
		}
	}
	return ""
}

func (c *InputCard) ravnicaWeekend() (string, string) {
	num := ExtractNumber(c.Variation)
	if strings.HasPrefix(num, "a") {
		return "GRN Ravnica Weekend", num
	} else if strings.HasPrefix(num, "b") {
		return "RNA Ravnica Weekend", num
	}

	for _, guild := range GRNGuilds {
		if c.Contains(guild) {
			return "GRN Ravnica Weekend", prwkVariants[c.Name][strings.ToLower(guild)]
		}
	}
	for _, guild := range ARNGuilds {
		if c.Contains(guild) {
			return "RNA Ravnica Weekend", prw2Variants[c.Name][strings.ToLower(guild)]
		}
	}
	return "", ""
}

func (c *InputCard) ravnicaGuidKit() string {
	if !c.Contains("Guild Kit") {
		return ""
	}

	if c.Contains("Guilds of Ravnica") || c.Contains("GRN") {
		return "GRN Guild Kit"
	}
	if c.Contains("Ravnica Allegiance") || c.Contains("RNA") {
		return "RNA Guild Kit"
	}

	for _, guild := range GRNGuilds {
		if c.Contains(guild) {
			return "GRN Guild Kit"
		}
	}
	for _, guild := range ARNGuilds {
		if c.Contains(guild) {
			return "RNA Guild Kit"
		}
	}

	if !c.isBasicLand() {
		if len(MatchInSet(c.Name, "GK1")) > 0 {
			return "GRN Guild Kit"
		}
		if len(MatchInSet(c.Name, "GK2")) > 0 {
			return "RNA Guild Kit"
		}
	} else {
		return "Guild Kit"
	}

	return ""
}

func (c *InputCard) Contains(prop string) bool {
	return Contains(c.Edition, prop) || Contains(c.Variation, prop)
}

func (c *InputCard) Equals(prop string) bool {
	return Equals(c.Edition, prop) || Equals(c.Variation, prop)
}

func ParseCommanderEdition(edition, variant string) string {
	if !strings.Contains(edition, "Commander") {
		return ""
	}

	// Append a custom display tag to avoid including the main set during filtering
	if strings.Contains(edition, "Display") || strings.Contains(edition, "Thick") ||
		strings.Contains(variant, "Display") || strings.Contains(variant, "Thick") {
		return edition + " Display"
	}

	// Legends series
	if strings.Contains(edition, "Legends") {
		if edition == "Commander Legends" {
			return "Commander Legends"
		} else if strings.Contains(edition, "Baldur's Gate") {
			edition = "Commander Legends: Battle for Baldur's Gate"
			return edition
		}
	}

	// Well-known extra tags
	perSetCommander := map[string]string{
		"Launch":     "Commander 2011 Launch Party",
		"Arsenal":    "Commander's Arsenal",
		"Ikoria":     "Commander 2020",
		"Strixhaven": "Commander 2021",
		"Starter":    "Starter Commander Decks",
	}
	for key, ed := range perSetCommander {
		if strings.Contains(edition, key) {
			return ed
		}
	}
	for key, ed := range backend.CommanderKeywordMap {
		if strings.Contains(edition, key) {
			if strings.Contains(edition, "Promo") || strings.Contains(variant, "Promo") {
				ed += " Promos"
			}
			return ed
		}
	}

	// Collection series
	if strings.Contains(edition, "Collection") {
		for _, color := range []string{"Green", "Black"} {
			if strings.Contains(edition, color) {
				return "Commander Collection: " + color
			}
		}
	}

	// Check Anthology, but decouple from volume 2
	if strings.Contains(edition, "Anthology") {
		for _, tag := range []string{"2018", "II", "Vol"} {
			if strings.Contains(edition, tag) {
				return "Commander Anthology Volume II"
			}
		}
		return "Commander Anthology"
	}

	// Is there a year available?
	year := ExtractYear(edition)
	if year != "" {
		return "Commander " + year
	}

	// Special fallbacks
	switch edition {
	case "Commander",
		"Commander Decks",
		"Commander Singles":
		return "Commander 2011"
	}

	return ""
}
