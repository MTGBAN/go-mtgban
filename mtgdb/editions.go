package mtgdb

var EditionTable = map[string]string{
	// Main expansions
	"10th Edition":              "Tenth Edition",
	"3rd Edition":               "Revised Edition",
	"3rd Edition/Revised":       "Revised Edition",
	"4th Edition":               "Fourth Edition",
	"5th Edition":               "Fifth Edition",
	"6th Edition":               "Classic Sixth Edition",
	"7th Edition":               "Seventh Edition",
	"8th Edition":               "Eighth Edition",
	"9th Edition":               "Ninth Edition",
	"Alpha":                     "Limited Edition Alpha",
	"Beta":                      "Limited Edition Beta",
	"Betrayers":                 "Betrayers of Kamigawa",
	"Champions":                 "Champions of Kamigawa",
	"Classic 6th Edition":       "Classic Sixth Edition",
	"Futuresight":               "Future Sight",
	"Hours of Devestation":      "Hour of Devastation",
	"Ravnica":                   "Ravnica: City of Guilds",
	"Revised":                   "Revised Edition",
	"Saviors":                   "Saviors of Kamigawa",
	"Time Spiral (Timeshifted)": "Time Spiral Timeshifted",
	"Time Spiral - Timeshifted": "Time Spiral Timeshifted",
	"Time Spiral Time Shifted":  "Time Spiral Timeshifted",
	"TimeShifted":               "Time Spiral Timeshifted",
	"Timeshifted":               "Time Spiral Timeshifted",
	"Unlimited":                 "Unlimited Edition",

	// JPN planeswalkers and similar
	"War of the Spark (Japanese Alternate Art)":              "War of the Spark",
	"War of the Spark: Japanese Alternate-Art Planeswalkers": "War of the Spark",
	"War of the Spark JPN Planeswalkers":                     "War of the Spark",

	// Gift pack
	"2017 Gift Pack":       "2017 Gift Pack",
	"2018 Gift Pack":       "M19 Gift Pack",
	"Gift Box 2017":        "2017 Gift Pack",
	"Gift Pack 2017":       "2017 Gift Pack",
	"Gift Pack 2018":       "M19 Gift Pack",
	"Shooting Star Promo":  "2017 Gift Pack",
	"Mark Poole Art Promo": "2017 Gift Pack",

	// Treasure Chest
	"Treasure Chest Promo": "XLN Treasure Chest",
	"Treare Map Promo":     "XLN Treasure Chest",
	"Treasure Map":         "XLN Treasure Chest",

	// Game Night
	"Game Night 2018":               "Game Night",
	"Magic Game Night":              "Game Night",
	"Magic Game Night 2019":         "Game Night 2019",
	"Game Night: 2018":              "Game Night",
	"Game Night: 2019":              "Game Night: 2019",
	"Magic Game Night 2018 Box Set": "Game Night",
	"Magic Game Night 2019 Box Set": "Game Night: 2019",

	// Old school lands
	"APAC Land":          "Asia Pacific Land Program",
	"Promos: Apac Lands": "Asia Pacific Land Program",
	"GURU":               "Guru",
	"Guru Land":          "Guru",
	"Guru":               "Guru",
	"Promos: Guru Lands": "Guru",
	"Promos: Euro Lands": "European Land Program",

	// Mystery Booster
	"Mystery Booster":             "Mystery Booster",
	"Mystery Booster Test Print":  "Mystery Booster Playtest Cards",
	"Mystery Booster Test Prints": "Mystery Booster Playtest Cards",
	"Mystery Booster - Test Card": "Mystery Booster Playtest Cards",
	"Playtest Card":               "Mystery Booster Playtest Cards",

	// Secre Lair
	"Secret Lair":             "Secret Lair Drop",
	"Secret Lair Drop Series": "Secret Lair Drop",
	"Secret Lair Full Art":    "Secret Lair Drop",
	"Stained Glass":           "Secret Lair Drop Promos",

	// Various
	"DCI Legend Membership": "DCI Legend Membership",
	"Legend Promo":          "DCI Legend Membership",
	"Pones: The Galloping":  "Ponies: The Galloping",
	"Ponies: The Galloping": "Ponies: The Galloping",
	"Champs / States Promo": "Champs and States",
	"Champs":                "Champs and States",

	// Welcome decks
	"Amonkhet Welcome Deck": "Welcome Deck 2017",
	"Magic 2016":            "Welcome Deck 2016",
	"Magic 2017":            "Welcome Deck 2017",
	"Welcome Deck 2016":     "Welcome Deck 2016",
	"Welcome Deck 2017":     "Welcome Deck 2017",

	// Holiday cards
	"Happy Holidays":     "Happy Holidays",
	"Holiday Foil":       "Happy Holidays",
	"Holiday Promo":      "Happy Holidays",
	"WOTC Employee Card": "Happy Holidays",

	// Standard Series
	"Standard Series":                          "BFZ Standard Series",
	"Standard Series Promo":                    "BFZ Standard Series",
	"2017 Standard Showdown":                   "XLN Standard Showdown",
	"2018 Standard Showdown":                   "M19 Standard Showdown",
	"2017 Standard Showdown Guay":              "XLN Standard Showdown",
	"2018 Standard Showdown Danner":            "M19 Standard Showdown",
	"Rebecca Guay Standard Showdown 2017":      "XLN Standard Showdown",
	"Alayna Danner Standard Showdown 2018":     "M19 Standard Showdown",
	"Rebecca Guay Standard Showdown":           "XLN Standard Showdown",
	"Alayna Danner Standard Showdown":          "M19 Standard Showdown",
	"Standard Showdown 2017":                   "XLN Standard Showdown",
	"Standard Showdown 2018":                   "M19 Standard Showdown",
	"Standard Showdown Rebecca Guay":           "XLN Standard Showdown",
	"Standard Showdown Alayna Danner":          "M19 Standard Showdown",
	"Alayna Danner Art":                        "M19 Standard Showdown",
	"Rebecca Guay Art Standard Showdown Promo": "XLN Standard Showdown",

	// Guild kits
	"Guild Kits: Guilds of Ravnica":  "GRN Guild Kit",
	"Guild Kits: Ravnica Allegiance": "RNA Guild Kit",
	"Guilds of Ravnica: Guild Kits":  "GRN Guild Kit",
	"Ravnica Allegiance: Guild Kits": "RNA Guild Kit",
	"Guild Kit: Boros":               "GRN Guild Kit",
	"Guild Kit: Dimir":               "GRN Guild Kit",
	"Guild Kit: Golgari":             "GRN Guild Kit",
	"Guild Kit: Izzet":               "GRN Guild Kit",
	"Guild Kit: Selesnya":            "GRN Guild Kit",
	"Guild Kit: Azorius":             "RNA Guild Kit",
	"Guild Kit: Gruul":               "RNA Guild Kit",
	"Guild Kit: Orzhov":              "RNA Guild Kit",
	"Guild Kit: Rakdos":              "RNA Guild Kit",
	"Guild Kit: Simic":               "RNA Guild Kit",
	"Boros Guild Kit":                "GRN Guild Kit",
	"Dimir Guild Kit":                "GRN Guild Kit",
	"Golgari Guild Kit":              "GRN Guild Kit",
	"Izzet Guild Kit":                "GRN Guild Kit",
	"Selesnya Guild Kit":             "GRN Guild Kit",
	"Azorius Guild Kit":              "RNA Guild Kit",
	"Gruul Guild Kit":                "RNA Guild Kit",
	"Orzhov Guild Kit":               "RNA Guild Kit",
	"Rakdos Guild Kit":               "RNA Guild Kit",
	"Simic Guild Kit":                "RNA Guild Kit",
	"Ravnica Weekend Boros":          "GRN Ravnica Weekend",
	"Ravnica Weekend Dimir":          "GRN Ravnica Weekend",
	"Ravnica Weekend Golgari":        "GRN Ravnica Weekend",
	"Ravnica Weekend Izzet":          "GRN Ravnica Weekend",
	"Ravnica Weekend Selesnya":       "GRN Ravnica Weekend",
	"Ravnica Weekend Azorius":        "RNA Ravnica Weekend",
	"Ravnica Weekend Gruul":          "RNA Ravnica Weekend",
	"Ravnica Weekend Orzhov":         "RNA Ravnica Weekend",
	"Ravnica Weekend Rakdos":         "RNA Ravnica Weekend",
	"Ravnica Weekend Simic":          "RNA Ravnica Weekend",
	"Guilds of Ravnica Guild Kits":   "GRN Ravnica Weekend",
	"Ravnica Allegiance Guild Kits":  "RNA Ravnica Weekend",

	// Commander
	"Commander 2011 Edition":       "Commander 2011",
	"Commander 2013 Edition":       "Commander 2013",
	"Commander 2014 Edition":       "Commander 2014",
	"Commander 2015 Edition":       "Commander 2015",
	"Commander 2016 Edition":       "Commander 2016",
	"Commander 2017 Edition":       "Commander 2017",
	"Commander 2018 Edition":       "Commander 2018",
	"Commander 2019 Edition":       "Commander 2019",
	"Commander 2020 Edition":       "Commander 2020",
	"Commander 2020: Ikoria":       "Commander 2020",
	"Commander Anthology 2018":     "Commander Anthology Volume II",
	"Commander Anthology VOL. II":  "Commander Anthology Volume II",
	"Commander Anthology Vol. II":  "Commander Anthology Volume II",
	"Commander Anthology Volume 2": "Commander Anthology Volume II",
	"Commander Singles":            "Commander 2011",
	"Commander Decks":              "Commander 2011",
	"Commander":                    "Commander 2011",
	"Commander: 2011 Edition":      "Commander 2011",
	"Commander: 2013 Edition":      "Commander 2013",
	"Commander: 2014 Edition":      "Commander 2014",
	"Commander: 2015 Edition":      "Commander 2015",
	"Commander: 2016 Edition":      "Commander 2016",
	"Commander: 2017 Edition":      "Commander 2017",
	"Commander: 2018 Edition":      "Commander 2018",
	"Commander: 2019 Edition":      "Commander 2019",
	"Commander: 2020 Edition":      "Commander 2020",
	"Commander: Ikoria":            "Commander 2020",

	// Modern Masters
	"Modern Masters 2013":            "Modern Masters",
	"Modern Masters 2013 Edition":    "Modern Masters",
	"Modern Masters 2015 Edition":    "Modern Masters 2015",
	"Modern Masters 2017 Edition":    "Modern Masters 2017",
	"Modern Masters: 2013 Edition":   "Modern Masters",
	"Modern Masters: 2015 Edition":   "Modern Masters 2015",
	"Modern Masters: 2017 Edition":   "Modern Masters 2017",
	"Ultimate Box Toppers":           "Ultimate Box Topper",
	"Ultimate Masters - Box Toppers": "Ultimate Box Topper",
	"Ultimate Masters Box Toppers":   "Ultimate Box Topper",
	"Ultimate Masters: Box Topper":   "Ultimate Box Topper",
	"Ultimate Masters: Box Toppers":  "Ultimate Box Topper",

	// CE and IE editions
	"Collector's Edition - International": "Intl. Collectors' Edition",
	"Collectors Ed Intl":                  "Intl. Collectors' Edition",
	"Collectors' Edition - International": "Intl. Collectors' Edition",
	"International Collector's Edition":   "Intl. Collectors' Edition",
	"International Collectors’ Edition":   "Intl. Collectors' Edition",
	"International Edition":               "Intl. Collectors' Edition",
	"Collector's Edition (Domestic)":      "Collectors' Edition",
	"Collector's Edition - Domestic":      "Collectors' Edition",
	"Collector's Edition":                 "Collectors' Edition",
	"Collectors Ed":                       "Collectors' Edition",
	"Collectors' Edition":                 "Collectors' Edition",
	"Collectors’ Edition":                 "Collectors' Edition",

	// Portal
	"Portal 1":          "Portal",
	"Portal II":         "Portal Second Age",
	"Portal 3K":         "Portal Three Kingdoms",
	"Portal 3 Kingdoms": "Portal Three Kingdoms",

	// Duel Decks
	"Japanese Jace vs. Chandra Foil": "Duel Decks: Jace vs. Chandra",
	"Duel Deck Heros VS Monsters":    "Duel Decks: Heroes vs. Monsters",
	"Duel Decks: Heros vs. Monsters": "Duel Decks: Heroes vs. Monsters",
	"Duel Decks: Kiora vs. Elspeth":  "Duel Decks: Elspeth vs. Kiora",
	"Duel Decks: Kiora vs Elspeth":   "Duel Decks: Elspeth vs. Kiora",
	"DD: Anthology":                  "Duel Decks Anthology",

	// Various series
	"Global Series: Jiangg Yanggu & Mu Yanling":  "Global Series Jiang Yanggu & Mu Yanling",
	"Global Series: Jiang Yanggu & Mu Yanling":   "Global Series Jiang Yanggu & Mu Yanling",
	"Global Series: Jiang Yanggu and Mu Yanling": "Global Series Jiang Yanggu & Mu Yanling",
	"Fire & Lightning":                           "Premium Deck Series: Fire and Lightning",
	"PDS: Fire & Lightning":                      "Premium Deck Series: Fire and Lightning",
	"Premium Deck Fire and Lightning":            "Premium Deck Series: Fire and Lightning",
	"Premium Deck: Fire and Lightning":           "Premium Deck Series: Fire and Lightning",
	"Premium Deck Series: Fire & Lightning":      "Premium Deck Series: Fire and Lightning",
	"Graveborn":                                  "Premium Deck Series: Graveborn",
	"PDS: Graveborn":                             "Premium Deck Series: Graveborn",
	"Premium Deck Graveborn":                     "Premium Deck Series: Graveborn",
	"Premium Deck: Graveborn":                    "Premium Deck Series: Graveborn",
	"Slivers":                                    "Premium Deck Series: Slivers",
	"PDS: Slivers":                               "Premium Deck Series: Slivers",
	"Premium Deck Slivers":                       "Premium Deck Series: Slivers",
	"Premium Deck: Slivers":                      "Premium Deck Series: Slivers",

	// Planechase
	"Planechase 2009":                 "Planechase",
	"Planechase (2009 Edition)":       "Planechase",
	"Planechase (2012 Edition)":       "Planechase 2012",
	"Planechase 2009 Edition":         "Planechase",
	"Planechase 2012 Edition":         "Planechase 2012",
	"Planechase: 2009 Edition":        "Planechase",
	"Planechase: 2012 Edition":        "Planechase 2012",
	"Planechase Planes: 2009 Edition": "Planechase Planes",
	"Planechase Planes: 2012 Edition": "Planechase 2012 Planes",

	// Deckmasters
	"Deckmaster Promo": "Deckmasters",
	"Deckmaster":       "Deckmasters",
	"Deckmasters":      "Deckmasters",

	// Junior Super/Europe/APAC Series
	"European Junior Series":         "Junior Series Europe",
	"Junior Series Promo":            "Junior Series Europe",
	"Junior Series Promos":           "Junior Series Europe",
	"Euro JSS Promo":                 "Junior Series Europe",
	"Junior Series":                  "Junior Series Europe",
	"Japan JSS":                      "Junior APAC Series",
	"Japan Junior Tournament Promo":  "Junior APAC Series",
	"Junior APAC Series":             "Junior APAC Series",
	"Junior APAC Series Promos":      "Junior APAC Series",
	"Junior Super Series":            "Junior Super Series",
	"Magic Scholarship Series":       "Junior Super Series",
	"Magic Scholarship Series Promo": "Junior Super Series",
	"Scholarship Series":             "Junior Super Series",
	"Scholarship Series Promo":       "Junior Super Series",
	"MSS":           "Junior Super Series",
	"MSS Promo":     "Junior Super Series",
	"JSS":           "Junior Super Series",
	"JSS DCI PROMO": "Junior Super Series",
	"JSS Foil":      "Junior Super Series",
	"JSS Promo":     "Junior Super Series",

	// GP Promos
	"2010 Grand Prix Promo":              "Grand Prix Promos",
	"2018 Grand Prix Promo":              "Grand Prix Promos",
	"GP Promo":                           "Grand Prix Promos",
	"Gran Prix Promo":                    "Grand Prix Promos",
	"Grand Prix":                         "Grand Prix Promos",
	"Grand Prix Foil":                    "Grand Prix Promos",
	"Grand Prix Promo":                   "Grand Prix Promos",
	"Promos: Grand Prix":                 "Grand Prix Promos",
	"Grand Prix 2018":                    "MagicFest 2019",
	"MagicFest 2019":                     "MagicFest 2019",
	"MagicFest 2020":                     "MagicFest 2020",
	"MagicFest Foil - 2020":              "MagicFest 2020",
	"FOIL 2019 MF MagicFest GP Promo":    "MagicFest 2019",
	"NONFOIL 2019 MF MagicFest GP Promo": "MagicFest 2019",
	"FOIL 2020 MF MagicFest GP Promo":    "MagicFest 2020",
	"NONFOIL 2020 MF MagicFest GP Promo": "MagicFest 2020",

	// Nationals
	"2018 Nationals Promo": "Nationals Promos",
	"Nationals":            "Nationals Promos",

	// Pro Tour Promos
	"2011 Pro Tour Promo": "Pro Tour Promos",
	"MCQ Promo":           "Pro Tour Promos",
	"MCQ":                 "Pro Tour Promos",
	"Mythic Championship Qualifier Promo": "Pro Tour Promos",
	"Mythic Championship":                 "Pro Tour Promos",
	"Players Tour Qualifier PTQ Promo":    "Pro Tour Promos",
	"Players Tour Qualifier":              "Pro Tour Promos",
	"Pro Tour Foil":                       "Pro Tour Promos",
	"Pro Tour Promo":                      "Pro Tour Promos",
	"Pro Tour Promos":                     "Pro Tour Promos",
	"Pro Tour":                            "Pro Tour Promos",
	"RPTQ Promo":                          "Pro Tour Promos",
	"RPTQ":                                "Pro Tour Promos",
	"Regional PTQ Promo Foil": "Pro Tour Promos",
	"Regional PTQ Promo":      "Pro Tour Promos",
	"Regional PTQ":            "Pro Tour Promos",

	// Worlds
	"2015 World Magic Cup Qualifier":      "World Magic Cup Qualifiers",
	"2016 WMCQ Promo":                     "World Magic Cup Qualifiers",
	"2017 WMCQ Promo":                     "World Magic Cup Qualifiers",
	"DCI Promo World Magic Cup Qualifier": "World Magic Cup Qualifiers",
	"WCQ":                             "World Magic Cup Qualifiers",
	"WMC Promo":                       "World Magic Cup Qualifiers",
	"WMC Qualifier":                   "World Magic Cup Qualifiers",
	"WMC":                             "World Magic Cup Qualifiers",
	"WMCQ Foil":                       "World Magic Cup Qualifiers",
	"WMCQ Promo":                      "World Magic Cup Qualifiers",
	"WMCQ Promo 2016":                 "World Magic Cup Qualifiers",
	"WMCQ Promo 2017":                 "World Magic Cup Qualifiers",
	"WMCQ":                            "World Magic Cup Qualifiers",
	"World Magic Cup":                 "World Magic Cup Qualifiers",
	"World Magic Cup Promo":           "World Magic Cup Qualifiers",
	"World Magic Cup Qualifier Promo": "World Magic Cup Qualifiers",
	"World Championship Foil":         "World Championship Promos",
	"World Cup Qualifier Promo":       "World Magic Cup Qualifiers",

	// Tarkir extra sets
	"Dragonfury":                              "Tarkir Dragonfury",
	"Dragonfury Promo":                        "Tarkir Dragonfury",
	"Dragons of Tarkir Dragonfury Game Promo": "Tarkir Dragonfury",
	"Tarkir Dragonfury":                       "Tarkir Dragonfury",
	"Tarkir Dragonfury Promo":                 "Tarkir Dragonfury",
	"Promos: Ugin's Fate":                     "Ugin's Fate",
	"Ugin's Fate Promo":                       "Ugin's Fate",
	"Ugin's Fate":                             "Ugin's Fate",
	"Ugins Fate":                              "Ugin's Fate",
	"Ugins Fate Promo":                        "Ugin's Fate",
	"Ugin's Fate Alternate Art Promo":         "Ugin's Fate",
	"Ugin’s Fate Promo":                       "Ugin's Fate",

	// Clash packs
	"Fate Reforged Clash Pack": "Fate Reforged Clash Pack",
	"Magic 2015 Clash Deck":    "Magic 2015 Clash Pack",
	"Magic 2015 Clash Pack":    "Magic 2015 Clash Pack",
	"Magic Origins Clash Pack": "Magic Origins Clash Pack",

	// Resale
	"Media Promo":     "Resale Promos",
	"Repack Insert":   "Resale Promos",
	"Resale Foil":     "Resale Promos",
	"Resale Promo":    "Resale Promos",
	"Resale Walmart ": "Resale Promos",
	"Resale Walmart":  "Resale Promos",
	"Resale":          "Resale Promos",
	"Store Foil":      "Resale Promos",
	"Walmart Resale":  "Resale Promos",

	// 15th Anniversary
	"15th Anniversary":           "15th Anniversary Cards",
	"15th Anniversary Foil":      "15th Anniversary Cards",
	"15th Anniversary Promo":     "15th Anniversary Cards",
	"DCI Promo 15th Anniversary": "15th Anniversary Cards",
	"MTG 15th Anniversary":       "15th Anniversary Cards",

	// Convention
	"2HG":                        "Two-Headed Giant Tournament",
	"2HG Foil":                   "Two-Headed Giant Tournament",
	"DCI Promo Two-Headed Giant": "Two-Headed Giant Tournament",
	"Two-Headed Giant":           "Two-Headed Giant Tournament",
	"Dragon'Con 1994":            "Dragon Con",
	"HASCON":                     "HasCon 2017",
	"HASCON 2017":                "HasCon 2017",
	"Hascon Promo":               "HasCon 2017",
	"Hascon 2017 Promo":          "HasCon 2017",
	"PAX Prime Promo":            "URL/Convention Promos",
	"2012 Convention Promo":      "URL/Convention Promos",
	"URL Convention Promo":       "URL/Convention Promos",

	// Summer
	"Summer Foil":           "Summer of Magic",
	"Summer of Magic Promo": "Summer of Magic",
	"Summer of Magic":       "Summer of Magic",

	// Release/Prerelease
	"Release Event": "Release Events",

	// Wotc Store
	"Foil Beta Picture":       "Wizards of the Coast Online Store",
	"Redemption Original Art": "Wizards of the Coast Online Store",

	// SDCC
	"San Diego Comic-Con": "San Diego Comic-Con 2019",

	// Archenemy
	"Archenemy: 2010 Edition":          "Archenemy",
	"Archenemy - Nicol Bolas":          "Archenemy: Nicol Bolas",
	"Archenemy Schemes (2010 Edition)": "Archenemy Schemes",

	// Various
	"Arena IA":                     "Arena League 2001",
	"Beatdown":                     "Beatdown Box Set",
	"Battle Royale":                "Battle Royale Box Set",
	"Conspiracy: 2014 Edition":     "Conspiracy",
	"Coldsnap Reprints":            "Coldsnap Theme Decks",
	"Coldsnap Theme Deck Reprints": "Coldsnap Theme Decks",
	"Introductory 4th Edition":     "Introductory Two-Player Set",
	"PS3 Promo":                    "Duels of the Planeswalkers 2012 Promos",
	"Starter":                      "Starter 1999",
	"Vanguard":                     "Vanguard Series",

	"Modern Event Deck":                           "Modern Event Deck 2014",
	"Modern Event Deck - March of the Multitudes": "Modern Event Deck 2014",

	// Foreign-only
	"3rd Edition (Foreign Black Border)": "Foreign Black Border",

	// Foreign translations
	"Alleanze":                   "Alliances",
	"Apocalisse":                 "Apocalypse",
	"Ascesa Oscura":              "Dark Ascension",
	"Ascesa degli Eldrazi":       "Rise of the Eldrazi",
	"Assalto":                    "Onslaught",
	"Aurora":                     "Morningtide",
	"Battaglia per Zendikar":     "Battle for Zendikar",
	"Campioni di Kamigawa":       "Champions of Kamigawa",
	"Caos Dimensionale":          "Planar Chaos",
	"Cavalcavento":               "Weatherlight",
	"Cicatrici di Mirrodin":      "Scars of Mirrodin",
	"Commander Arsenal":          "Commander's Arsenal",
	"Congiunzione":               "Planeshift",
	"Decima Edizione":            "Tenth Edition",
	"Destino di Urza":            "Urza's Destiny",
	"Discordia":                  "Dissension",
	"Draghi di Tarkir":           "Dragons of Tarkir",
	"Era Glaciale":               "Ice Age",
	"Eredità di Urza":            "Urza's Legacy",
	"Esodo":                      "Exodus",
	"Fedeltà di Ravnica":         "Ravnica Allegiance",
	"Figli degli Dei":            "Born of the Gods",
	"Flagello":                   "Scourge",
	"Fortezza":                   "Stronghold",
	"Frammenti di Alara":         "Shards of Alara",
	"Gilde di Ravnica":           "Guilds of Ravnica",
	"Giuramento dei Guardiani":   "Oath of the Gatewatch",
	"I Khan di Tarkir":           "Khans of Tarkir",
	"Il Patto delle Gilde":       "Guildpact",
	"Invasione":                  "Invasion",
	"Irruzione":                  "Gatecrash",
	"L'Era della Rovina":         "Hour of Devastation",
	"La Guerra della Scintilla":  "War of the Spark",
	"Labirinto del Drago":        "Dragon's Maze",
	"Landa Tenebrosa":            "Shadowmoor",
	"Leggende":                   "Legends",
	"Legioni":                    "Legions",
	"Liberatori di Kamigawa":     "Saviors of Kamigawa",
	"Luna Spettrale":             "Eldritch Moon",
	"Maschere di Mercadia":       "Mercadian Masques",
	"Mirrodin Assediato":         "Mirrodin Besieged",
	"Nona Edizione":              "Ninth Edition",
	"Nuova Phyrexia":             "New Phyrexia",
	"Odissea":                    "Odyssey",
	"Ombre su Innistrad":         "Shadows over Innistrad",
	"Ondata Glaciale":            "Coldsnap",
	"Origini":                    "Homelands",
	"Orizzonti di Modern":        "Modern Horizons",
	"Ottava Edizione":            "Eighth Edition",
	"Profezia":                   "Prophecy",
	"Quarta Edizione":            "Fourth Edition",
	"Quinta Alba":                "Fifth Dawn",
	"Quinta Edizione":            "Fifth Edition",
	"Ravnica: Città delle Gilde": "Ravnica: City of Guilds",
	"Revised EU FBB":             "Foreign Black Border",
	"Revised EU FWB":             "Foreign White Border",
	"Riforgiare il Destino":      "Fate Reforged",
	"Rinascita di Alara":         "Alara Reborn",
	"Ritorno a Ravnica":          "Return to Ravnica",
	"Ritorno di Avacyn":          "Avacyn Restored",
	"Rivali di Ixalan":           "Rivals of Ixalan",
	"Rivolta dell'Etere":         "Aether Revolt",
	"Saga di Urza":               "Urza's Saga",
	"Sentenza":                   "Judgment",
	"Sesta Edizione":             "Classic Sixth Edition",
	"Settima Edizione":           "Seventh Edition",
	"Spirale Temporale":          "Time Spiral",
	"Tempesta":                   "Tempest",
	"Theros: Oltre la Morte":     "Theros Beyond Death",
	"Tormento":                   "Torment",
	"Traditori di Kamigawa":      "Betrayers of Kamigawa",
	"Trono di Eldraine":          "Throne of Eldraine",
	"Vespro":                     "Eventide",
	"Viaggio Verso Nyx":          "Journey into Nyx",
	"Visione Futura":             "Future Sight",
	"Visioni":                    "Visions",

	"Duel Deck: Ajani Vs Bolas":        "Duel Decks: Ajani vs. Nicol Bolas",
	"Duel Deck: Cavalieri vs Draghi":   "Duel Decks: Knights vs. Dragons",
	"Duel Deck: Elfi vs Goblin":        "Duel Decks: Elves vs. Goblins",
	"Duel Deck: Elspeth vs Tezzereth":  "Duel Decks: Elspeth vs. Tezzeret",
	"Duel Decks: Cavalieri vs. Draghi": "Duel Decks: Knights vs. Dragons",
}
