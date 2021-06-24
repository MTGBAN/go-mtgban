package mtgmatcher

import (
	"log"
	"os"
	"testing"

	"github.com/kodabb/go-mtgban/mtgmatcher/mtgjson"
)

type MatchTest struct {
	Id   string
	Err  error
	Desc string
	In   Card
}

var MatchTests = []MatchTest{
	// Errors
	MatchTest{
		Desc: "card_does_not_exist",
		Err:  ErrCardDoesNotExist,
		In: Card{
			Name: "I do not exist",
		},
	},
	MatchTest{
		Desc: "wrong_card_number",
		Err:  ErrCardWrongVariant,
		In: Card{
			Name:      "Arcane Denial",
			Variation: "10",
			Edition:   "Alliances",
		},
	},
	MatchTest{
		Desc: "not_in_a_promo_pack",
		Err:  ErrCardNotInEdition,
		In: Card{
			Name:    "Demonic Tutor",
			Edition: "Promo Pack",
		},
	},
	MatchTest{
		Desc: "not_a_prerelease",
		Err:  ErrCardNotInEdition,
		In: Card{
			Name:      "Lobotomy",
			Variation: "Prerelease",
			Edition:   "ignored",
		},
	},
	MatchTest{
		Desc: "pure_aliasing",
		Err:  ErrAliasing,
		In: Card{
			Name:      "Forest",
			Variation: "",
			Edition:   "Zendikar",
		},
	},
	MatchTest{
		Desc: "missing_variant",
		Err:  ErrCardMissingVariant,
		In: Card{
			Name:      "Arcane Denial",
			Variation: "",
			Edition:   "Alliances",
		},
	},
	MatchTest{
		Desc: "invalid_id",
		Err:  ErrCardDoesNotExist,
		In: Card{
			Id: "not-an-id",
		},
	},

	// ID lookup
	MatchTest{
		Id:   "f3a94132-ce71-5556-bfd3-1461601a810d",
		Desc: "id_lookup_mtgjson",
		In: Card{
			Id: "f3a94132-ce71-5556-bfd3-1461601a810d",
		},
	},
	MatchTest{
		Id:   "f3a94132-ce71-5556-bfd3-1461601a810d_f",
		Desc: "id_lookup_foil",
		In: Card{
			Id:   "f3a94132-ce71-5556-bfd3-1461601a810d",
			Foil: true,
		},
	},
	MatchTest{
		Id:   "b2f56602-e85a-588f-a4be-40b6e56f44f7",
		Desc: "id_lookup_alt_foil",
		In: Card{
			Id:   "be665b02-1cf2-50c6-8861-85da921bc853",
			Foil: true,
		},
	},
	MatchTest{
		Id:   "095ecbbf-7e08-5bc7-85f9-7d830b9141a2_f",
		Desc: "id_lookup_variant_ea_foil",
		In: Card{
			Id:   "095ecbbf-7e08-5bc7-85f9-7d830b9141a2",
			Foil: true,
		},
	},
	MatchTest{
		Id:   "0b77872e-c4eb-54f4-8f63-57fcf68e682a",
		Desc: "id_lookup_non_foil",
		In: Card{
			Id: "7a6610cf-7d53-5e4b-bb76-9be1708d3892",
		},
	},
	MatchTest{
		Id:   "38fc8e52-821b-5701-a5bd-7402ee13a51c",
		Desc: "id_lookup_multi_variations",
		In: Card{
			Id: "a67c3cd2-d09d-5c2e-807b-5f4772c367e7",
		},
	},
	MatchTest{
		Id:   "91cff9f4-41ec-5359-8723-dd2e6ba6efc4",
		Desc: "id_lookup_scryfall",
		In: Card{
			Id: "0535c823-f6e9-4a2f-8adf-f69b6f0fea1f",
		},
	},

	// Number duplicates
	MatchTest{
		Id:   "fb083deb-30ea-5ff4-8aa8-cee8531cd7ec",
		Desc: "fullart_land",
		In: Card{
			Name:      "Swamp",
			Variation: "241",
			Edition:   "Zendikar",
		},
	},
	MatchTest{
		Id:   "aed5fe79-ddec-5bf7-93b3-63a042faf863",
		Desc: "fullart_land_could_be_confused_with_suffix",
		In: Card{
			Name:      "Island",
			Variation: "234 A - Full Art",
			Edition:   "Zendikar",
		},
	},
	MatchTest{
		Id:   "4fb5d3f7-cc7b-5502-8906-555ba919bd02",
		Desc: "nonfullart_land_could_be_confused_with_suffix",
		In: Card{
			Name:      "Forest",
			Variation: "274 A - Non-Full Art",
			Edition:   "Battle for Zendikar",
		},
	},
	MatchTest{
		Id:   "df3a4387-62c5-5fcc-a675-1c5e04d6103b",
		Desc: "complex_number_variant",
		In: Card{
			Name:      "Plains",
			Variation: "87 - A",
			Edition:   "Unsanctioned",
		},
	},
	MatchTest{
		Id:   "a7d7f03a-d876-52aa-97f6-44d371226533",
		Desc: "alternative_complex_number_variant",
		In: Card{
			Name:      "Brothers Yamazaki",
			Variation: "160 A",
			Edition:   "Champions of Kamigawa",
		},
	},
	MatchTest{
		Id:   "aacfd47d-9b20-52ad-a62c-cba3414357ad",
		Desc: "second_alternative_complex_number_variant",
		In: Card{
			Name:      "Brothers Yamazaki",
			Variation: "160 - B",
			Edition:   "Champions of Kamigawa",
		},
	},
	MatchTest{
		Id:   "c7f233d4-0770-5b10-9836-b4034047a9f8",
		Desc: "borderless_lands",
		In: Card{
			Name:    "Plains",
			Edition: "Unstable",
		},
	},
	MatchTest{
		Id:   "af8f1ee0-f235-5e76-9994-60c8d809da47",
		Desc: "weekend_lands",
		In: Card{
			Name:      "Island",
			Variation: "Ravnica Weekend B02",
			Edition:   "ignored",
		},
	},
	MatchTest{
		Id:   "7b64a5cf-a4c9-5391-bbba-0dd945281569",
		Desc: "mps_lands_2006",
		In: Card{
			Name:    "Swamp",
			Edition: "Magic Premiere Shop 2006",
			Foil:    true,
		},
	},
	MatchTest{
		Id:   "77321166-66e1-5e9a-b630-35dcddb4b818",
		Desc: "mps_lands_2009",
		In: Card{
			Name:      "Island",
			Variation: "Rob Alexander MPS 2009",
			Edition:   "Promos: MPS Lands",
		},
	},
	MatchTest{
		Id:   "053c8559-8ab8-5a1a-9444-6140b41470c4",
		Desc: "plains_from_set_with_special_cards_and_C_to_be_ignored",
		In: Card{
			Name:      "Plains",
			Variation: "366 C",
			Edition:   "Tenth Edition",
		},
	},
	MatchTest{
		Id:   "5675b6f8-ca15-5455-aaf7-56dfb038ec52",
		Desc: "nonfullart_land_with_incomplete_number",
		In: Card{
			Name:      "Forest",
			Variation: "Non-Full Art 271",
			Edition:   "Battle for Zendikar",
		},
	},
	MatchTest{
		Id:   "b2f53204-4357-56e7-a8f4-7f29ed9e674c",
		Desc: "land_with_letter",
		In: Card{
			Name:    "Forest B",
			Edition: "Intl. Collectors' Edition",
		},
	},
	MatchTest{
		Id:   "d0fbb33b-cd41-5fb7-8518-382dd07860d1",
		Desc: "forest_F",
		In: Card{
			Name:    "Forest F",
			Edition: "Battle Royale",
		},
	},
	MatchTest{
		Id:   "4fb5d3f7-cc7b-5502-8906-555ba919bd02",
		Desc: "intro_land",
		In: Card{
			Name:      "Forest",
			Variation: "274 Intro",
			Edition:   "Battle for Zendikar",
		},
	},
	MatchTest{
		Id:   "057df7fb-238d-55b1-93dd-ec76548a0fca",
		Desc: "land_with_collectors_number",
		In: Card{
			Name:    "Forest 277",
			Edition: "Ixalan",
		},
	},

	// Naming conventions
	MatchTest{
		Id:   "bf0aa055-3635-5efd-930d-4f0a7caaa411",
		Desc: "transform_card",
		In: Card{
			Name:    "Daybreak Ranger / Nightfall Predator",
			Edition: "Innistrad",
		},
	},
	MatchTest{
		Id:   "7e6e9448-42ab-58e1-828a-ebef7b5ada77",
		Desc: "aftermath_card",
		In: Card{
			Name:      "Commit to Memory",
			Variation: "Prerelease",
			Edition:   "ignored",
		},
	},
	MatchTest{
		Id:   "0ad9df53-068e-5bbd-9a83-d0dc4168ce6e",
		Desc: "split_card",
		In: Card{
			Name:    "Fire // Ice",
			Edition: "Apocalypse",
		},
	},
	MatchTest{
		Id:   "88400e25-72b6-54b9-8e0d-40851b42bcdd",
		Desc: "flip_card_not_really",
		In: Card{
			Name:    "Journey to Eternity",
			Edition: "Rivals of Ixalan",
		},
	},
	MatchTest{
		Id:   "7170634e-89fc-5e19-b7e6-ae4393b143d5",
		Desc: "flip_card",
		In: Card{
			Name:      "Startled Awake - Persistent Nightmare",
			Variation: "Prerelease",
			Edition:   "ignored",
		},
	},
	MatchTest{
		Id:   "7170634e-89fc-5e19-b7e6-ae4393b143d5",
		Desc: "flip_card_half",
		In: Card{
			Name:      "Startled Awake",
			Variation: "Prerelease",
			Edition:   "ignored",
		},
	},
	MatchTest{
		Id:   "6081d8fd-26e1-5e6a-9c98-417d03214856",
		Desc: "meld_card",
		In: Card{
			Name:    "Bruna, the Fading Light",
			Edition: "Eldritch Moon",
		},
	},
	MatchTest{
		Id:   "46b6d569-0deb-58e3-af91-f4652dd709bc",
		Desc: "meld_card_b",
		In: Card{
			Name:      "Gisela, the Broken Blade | Brisela, Voice of Nightmares",
			Variation: "Prerelease",
			Edition:   "ignored",
		},
	},
	MatchTest{
		Id:   "98a0d909-5db8-53b2-9e5a-080a9b7e94e8",
		Desc: "triple_card",
		In: Card{
			Name:    "Smelt // Herd // Saw",
			Edition: "Mystery Booster Playtest Ins",
		},
	},
	MatchTest{
		Id:   "17100b0d-3b74-5329-a832-dfad07d5c35b",
		Desc: "double_card_with_lots_of_matches_first",
		In: Card{
			Name:    "Start",
			Edition: "Amonkhet",
		},
	},
	MatchTest{
		Id:   "2b29b236-1f37-5d41-b15e-a4fafe265885",
		Desc: "double_card_with_lots_of_matches_second",
		In: Card{
			Name:    "Start",
			Edition: "Mystery Booster Playtest Cards",
		},
	},
	MatchTest{
		Id:   "25463956-7fc1-5781-88cb-abba28a59ddd",
		Desc: "incorrect_name_but_salvageable",
		In: Card{
			Name:    "B.O.B.",
			Edition: "Unsanctioned",
		},
	},
	MatchTest{
		Id:   "c98bf90b-e3b8-5a16-a797-73391ca6e4d6",
		Desc: "parenthesis_in_the_name",
		In: Card{
			Name:    "Erase (Not the Urza's Legacy One)",
			Edition: "Unhinged",
		},
	},
	MatchTest{
		Id:   "451ac233-9ba8-59db-8fff-6962e0b173f6",
		Desc: "parenthesis_in_the_name_and_variation",
		In: Card{
			Name:      "B.F.M. (Big Furry Monster)",
			Variation: "29",
			Edition:   "Unglued",
		},
	},
	MatchTest{
		Id:   "33b10155-90dc-5846-9c82-af414ba6c1e5",
		Desc: "custom_name_for_variation",
		In: Card{
			Name:    "B.F.M. Left",
			Edition: "Unglued",
		},
	},
	MatchTest{
		Id:   "fb3bdc21-d1c3-5fa2-8ea6-3ff48b11a5bc",
		Desc: "number_in_the_name",
		In: Card{
			Name:    "Serum Visions (30)",
			Edition: "Secret Lair Drop",
		},
	},
	MatchTest{
		Id:   "f193238a-07a8-53b6-8383-e30e95353891",
		Desc: "ignore_b_side_face_foil_consequences",
		In: Card{
			Name:    "Curse of the Fire Penguin",
			Edition: "Unhinged",
		},
	},
	MatchTest{
		Id:   "dd154c88-8235-5ca0-9925-11391b89669a",
		Desc: "relaxed_naming",
		In: Card{
			Name:        "Two-Headed Giant of",
			Edition:     "Beta",
			MaybePrefix: true,
		},
	},

	// Incorrect editions
	MatchTest{
		Id:   "0ebdbff9-e756-511f-a17d-43951169d0ea",
		Desc: "incorrect_edition_but_card_has_a_single_printing",
		In: Card{
			Name:    "Mirrodin Besieged",
			Edition: "Scars of Phyrexia",
		},
	},
	MatchTest{
		Id:   "c0439fdf-36e9-5578-94c0-36a056ede97d",
		Desc: "incorrect_edition_set_name_does_not_interfere_with_number",
		In: Card{
			Name:      "Death Baron",
			Variation: "Convention Foil M19",
			Edition:   "ignored",
			Foil:      true,
		},
	},
	MatchTest{
		Id:   "ff963a6c-5c2c-5d51-a75b-abce3e819db1",
		Desc: "incorrect_edition_but_salvageable",
		In: Card{
			Name:    "Polluted Mire",
			Edition: "Duel Decks Anthology",
		},
	},
	MatchTest{
		Id:   "b6698d85-bcd9-5262-a91d-2b3eb746e24c",
		Desc: "incorrect_edition_but_salvageable_and_could_alias_other_cards",
		In: Card{
			Name:    "Garruk Wildspeaker",
			Edition: "Duel Decks Anthology",
		},
	},
	MatchTest{
		Id:   "70d57c64-c2a1-52e8-a807-bc965fb2ddb7",
		Desc: "incorrect_edition_but_salvageable_missing_only_a_chunk",
		In: Card{
			Name:    "No One Will Hear Your Cries",
			Edition: "Archenemy: Nicol Bolas",
		},
	},
	MatchTest{
		Id:   "f0f9e5f9-17e9-5827-b97e-c56e693a5beb",
		Desc: "incorrect_edition_belongs_to_a_foilonly_subset",
		In: Card{
			Name:      "Zur's Weirding",
			Variation: "Foil",
			Edition:   "Mystery Booster",
		},
	},
	MatchTest{
		Id:   "8067d275-fb65-5277-b812-8cf33604a788",
		Desc: "incorrect_edition_year_should_not_interfere",
		In: Card{
			Name:      "Yule Ooze",
			Variation: "2011 Holiday",
			Edition:   "Happy Holidays",
		},
	},
	MatchTest{
		Id:   "fbdf9d58-2746-5ed5-ad97-0a9780d066ca",
		Desc: "incorrect_edition_champs_and_states",
		In: Card{
			Name:      "Mutavault",
			Variation: "Extended art",
			Edition:   "ignored",
		},
	},
	MatchTest{
		Id:   "15bd27f2-c842-5f2c-9862-d3a8d36143b7",
		Desc: "champ_in_variant_but_not_champs_and_states",
		In: Card{
			Name:      "Champion of Lambholt",
			Variation: "commander-anthology-2018-champion-of-lambholt",
			Edition:   "Commander Anthology 2018",
		},
	},
	MatchTest{
		Id:   "2126fb89-1eab-5d06-a7d9-953db1242849",
		Desc: "mismatching_year",
		In: Card{
			Name:      "Mountain",
			Variation: "Grand Prix 2018",
			Edition:   "MagicFest 2019",
			Foil:      true,
		},
	},
	MatchTest{
		Id:   "5e7c2039-a072-5907-8053-74f518384e8a",
		Desc: "not_grand_prix_but_magic_fest",
		In: Card{
			Name:      "Lightning Bolt",
			Variation: "Grand Prix",
			Edition:   "ignored",
		},
	},
	MatchTest{
		Id:   "a6e7bc06-ea7d-5186-8dd4-d95086e4e8d2",
		Desc: "apac_lands",
		In: Card{
			Name:      "Forest",
			Variation: "Pete Venters",
			Edition:   "Asia Pacific Land Program",
		},
	},
	MatchTest{
		Id:   "7760d40f-2afd-5552-b59f-7d395ed4b7be",
		Desc: "euro_lands",
		In: Card{
			Name:      "Plains",
			Variation: "EURO Land Steppe Tundra Ben Thompson art",
			Edition:   "ignored",
		},
	},
	MatchTest{
		Id:   "44c98b84-1a4b-5101-912a-c26a87463cc5",
		Desc: "euro_lands_comma",
		In: Card{
			Name:      "Island",
			Variation: "EURO Land, Venezia",
			Edition:   "ignored",
		},
	},

	// Promo pack
	MatchTest{
		Id:   "9d3e0596-a001-51ac-922a-e4b11cd09126",
		Desc: "m20_promo_packs_lands",
		In: Card{
			Name:    "Plains",
			Edition: "Promo Pack",
		},
	},
	MatchTest{
		Id:   "f37a6201-033a-58b3-a816-94bf813891b7",
		Desc: "m20_promo_packs_lands_variant",
		In: Card{
			Name:      "Swamp",
			Variation: "M20 Promo Pack",
			Edition:   "ignored",
		},
	},
	MatchTest{
		Id:   "338a6aea-1d26-5a37-87f8-42de7ce9dc2b",
		Desc: "promo_pack_in_promos_with_pw_stamp",
		In: Card{
			Name:    "Zendikar Resurgent",
			Edition: "Promo Pack",
		},
	},
	MatchTest{
		Id:   "61ca0be6-5672-5d82-a9df-b10d15fc6be1",
		Desc: "promo_pack_in_expansion_with_inverted_frame",
		In: Card{
			Name:    "Alseid of Life's Bounty",
			Edition: "Promo Pack",
		},
	},
	MatchTest{
		Id:   "d6a7fe98-1a0b-569f-9a49-f012698fe2ba",
		Desc: "promo_pack_in_promos_with_inverted_frame",
		In: Card{
			Name:    "Negate",
			Edition: "Promo Pack",
		},
	},
	MatchTest{
		Id:   "beb3990e-ee5c-51e0-9651-e1e5a5f336c0",
		Desc: "nonpromo_pack_card_that_may_have_a_promo_pack_version",
		In: Card{
			Name:    "Slaying Fire",
			Edition: "Throne of Eldraine",
		},
	},
	MatchTest{
		Id:   "d3ca5ef5-78d2-5bf1-b10b-10119df615a7",
		Desc: "promo_pack_with_duplication",
		In: Card{
			Name:      "Sorcerous Spyglass",
			Variation: "Promo Pack XLN",
			Edition:   "ignored",
		},
	},
	MatchTest{
		Id:   "3438d2d9-4e87-573a-bc3f-28e703fe47ee",
		Desc: "nonpromo_pack_is_fine_too",
		In: Card{
			Name:    "Sorcerous Spyglass",
			Edition: "XLN",
		},
	},
	MatchTest{
		Id:   "0b473d58-3a60-5ed6-8de6-02866f43d125",
		Desc: "so_many_varaiants_and_untagged",
		In: Card{
			Name:      "Teferi, Master of Time",
			Variation: "Prerelease",
			Edition:   "ignored",
		},
	},
	MatchTest{
		Id:   "4f1f9adc-1aee-5d28-97c3-e44e51e38879",
		Desc: "so_many_varaiants_and_tagged",
		In: Card{
			Name:      "Teferi, Master of Time",
			Variation: "#075 Promo Pack Core Set 2021",
			Edition:   "ignored",
		},
	},
	MatchTest{
		Id:   "f877637f-fde4-52c8-a7be-ca5c494c38f3",
		Desc: "core_promo_pack_2020",
		In: Card{
			Name:    "Temple of Triumph",
			Edition: "Core 2020 - Promo Pack",
		},
	},
	MatchTest{
		Id:   "351ae10e-91ac-5baa-90b2-cbf1555d9de6",
		Desc: "core_promo_pack_2021",
		In: Card{
			Name:    "Temple of Triumph",
			Edition: "Core 2021 - Promo Pack",
		},
	},
	MatchTest{
		Id:   "e2033df0-b1b5-5b9b-90cb-bf738cda501b",
		Desc: "core_promo_pack_2021_alt",
		In: Card{
			Name:      "Temple of Mystery",
			Variation: "Promo Pack - M21",
			Edition:   "Promo Pack",
		},
	},
	MatchTest{
		Id:   "24a266a1-c13e-5790-a649-b03139d47ccc",
		Desc: "core_promo_pack_old",
		In: Card{
			Name:    "Heroic Intervention",
			Edition: "Core 2021 - Promo Pack",
		},
	},

	// Prerelease
	MatchTest{
		Id:   "8d9f230c-ad3f-5b5f-9132-96368bb04ace",
		Desc: "old_prerelease",
		In: Card{
			Name:      "Glory",
			Variation: "Prerelease",
			Edition:   "ignored",
		},
	},
	MatchTest{
		Id:   "45f8e8cc-db7e-5448-b5b5-c49b9e43efbd",
		Desc: "prerelease_in_promos_before_the_date_but_without_s_suffix",
		In: Card{
			Name:      "Scourge of Fleets",
			Variation: "Prerelease",
			Edition:   "ignored",
		},
	},
	MatchTest{
		Id:   "51be1904-c38c-5303-b2b7-d889aeb66819",
		Desc: "prerelease_with_s_suffix",
		In: Card{
			Name:      "Pristine Skywise",
			Variation: "Prerelease",
			Edition:   "ignored",
		},
	},
	MatchTest{
		Id:   "09f6ce54-cfba-5e7f-b2c1-c36f95e26ab1",
		Desc: "JPN_prerelease_with_s_suffix",
		In: Card{
			Name:      "Ugin, the Ineffable",
			Variation: "JPN Prerelease",
			Edition:   "ignored",
		},
	},
	MatchTest{
		Id:   "a8a8aac9-b986-5fbe-be9d-511fb272f6d8",
		Desc: "JPN_prerelease_with_s_suffix_but_number_could_interfere",
		In: Card{
			Name:      "Teyo, the Shieldmage",
			Variation: "032 - JPN Alternate Art Prerelease Foil",
			Edition:   "ignored",
		},
	},
	MatchTest{
		Id:   "5849e009-a208-59e0-b2dc-230b053bf015",
		Desc: "prerelease_in_promos_after_the_date_but_without_s_suffix",
		In: Card{
			Name:      "Astral Drift",
			Variation: "Prerelease",
			Edition:   "ignored",
		},
	},
	MatchTest{
		Id:   "7efa76e4-a6a8-5f6c-a337-9d88acc1d593",
		Desc: "prerelease_with_duplication",
		In: Card{
			Name:      "Sorcerous Spyglass",
			Variation: "Ixalan Prerelease",
			Edition:   "ignored",
		},
	},
	MatchTest{
		Id:   "6b0ad772-5c2b-59b2-9a51-5296d3383077",
		Desc: "lubu_dedup_july",
		In: Card{
			Name:      "Lu Bu, Master-at-Arms",
			Variation: "July 4 Prerelease",
			Edition:   "ignored",
		},
	},

	// JPN alternate art
	MatchTest{
		Id:   "a0c99852-b08e-5f09-9f48-317f2253df15",
		Desc: "normal_nonJPN_version",
		In: Card{
			Name:    "Vraska, Swarm's Eminence",
			Edition: "War of the Spark",
		},
	},
	MatchTest{
		Id:   "bfbf2df1-007f-500a-bfe5-1310ad1bad5d",
		Desc: "JPN_variant",
		In: Card{
			Name:      "Vraska, Swarm's Eminence",
			Variation: "JPN Alternate Art",
			Edition:   "War of the Spark",
		},
	},
	MatchTest{
		Id:   "2f4942dd-d6d7-5b79-8dd9-91c9cd0daf1c",
		Desc: "JPN_variant_but_number_could_interfere",
		In: Card{
			Name:      "Teyo, the Shieldmage",
			Variation: "032 - JPN Alternate Art",
			Edition:   "War of the Spark",
		},
	},

	// Borderless cards
	MatchTest{
		Id:   "46153afe-5e05-5082-852a-648c03924bcf",
		Desc: "normal_nonborderless_variant",
		In: Card{
			Name:    "Oko, Thief of Crowns",
			Edition: "Throne of Eldraine",
		},
	},
	MatchTest{
		Id:   "f203bad8-9c07-507c-9699-fc8fec69e2d2",
		Desc: "borderless_variant",
		In: Card{
			Name:      "Oko, Thief of Crowns",
			Variation: "Borderless",
			Edition:   "Throne of Eldraine",
		},
	},
	MatchTest{
		Id:   "09618c19-c87e-51e2-959c-8b176156c9ca",
		Desc: "borderless_but_from_a_funny_set",
		In: Card{
			Name:    "Sap Sucker",
			Edition: "Unstable",
		},
	},
	MatchTest{
		Id:   "1b30ae75-338a-574f-9232-e0b119f8a6a5",
		Desc: "borderless_boxtopper",
		In: Card{
			Name:    "Ancient Tomb",
			Edition: "PUMA",
		},
	},
	MatchTest{
		Id:   "f2bfda59-3560-55a7-80a1-9d19f1a1d163",
		Desc: "borderless_but_also_alternative",
		In: Card{
			Name:    "Doubling Season",
			Edition: "Double Masters Box Toppers",
		},
	},

	// Box topper-style extended art
	MatchTest{
		Id:   "c67e23df-18de-5668-83f8-4fc9c23299bf",
		Desc: "normal_nonextendedart_variant",
		In: Card{
			Name:    "Heliod's Intervention",
			Edition: "Theros Beyond Death",
		},
	},
	MatchTest{
		Id:   "39b197ca-526c-5aa4-be3c-97b5db042efc",
		Desc: "extendedart_variant",
		In: Card{
			Name:      "Heliod's Intervention",
			Variation: "Extended Art",
			Edition:   "Theros Beyond Death",
		},
	},
	MatchTest{
		Id:   "7bab8861-40c1-5895-91d7-504b0eaaeae1",
		Desc: "antialias_extendedart_vs_bab",
		In: Card{
			Name:      "Orah, Skyclave Hierophant",
			Variation: "Extended Art",
			Edition:   "Zendikar Rising",
		},
	},
	MatchTest{
		Id:   "fce0a486-fb7e-5e33-bea2-573c4e6529be",
		Desc: "antialias_bab_vs_extendedart",
		In: Card{
			Name:      "Orah, Skyclave Hierophant",
			Variation: "Buy-A-Box Foil",
			Edition:   "Zendikar Rising",
			Foil:      true,
		},
	},

	// Showcase frame
	MatchTest{
		Id:   "51c8a322-0601-51ed-b5f9-bebb5d97b5d9",
		Desc: "normal_nonshowcase_variant",
		In: Card{
			Name:    "Brazen Borrower",
			Edition: "Throne of Eldraine",
		},
	},
	MatchTest{
		Id:   "66819bb9-e044-512b-921e-7a5a82be79f5",
		Desc: "showcase_variant",
		In: Card{
			Name:      "Brazen Borrower",
			Variation: "Showcase",
			Edition:   "Throne of Eldraine",
		},
	},
	MatchTest{
		Id:   "f082ca30-a227-5c88-b95a-37ca03cefcd9",
		Desc: "showcase_borderless",
		In: Card{
			Name:      "Zagoth Triome",
			Variation: "Showcase",
			Edition:   "Ikoria: Lair of Behemoths",
		},
	},
	MatchTest{
		Id:   "cb0847b3-ef9b-560c-9cbd-37e91acfd86d",
		Desc: "correct_number_but_no_showcase_tag",
		In: Card{
			Name:      "Renata, Called to the Hunt",
			Variation: "267",
			Edition:   "Theros Beyond Death",
		},
	},

	// Reskinned frame
	MatchTest{
		Id:   "f3a94132-ce71-5556-bfd3-1461601a810d",
		Desc: "nongodzilla_variant",
		In: Card{
			Name:    "Sprite Dragon",
			Edition: "Ikoria: Lair of Behemoths",
		},
	},
	MatchTest{
		Id:   "793cba6b-4f65-5834-b055-b83822576a86_f",
		Desc: "nongodzilla_notjapanese_variant",
		In: Card{
			Name:    "Dirge Bat",
			Edition: "Ikoria: Lair of Behemoths",
			Foil:    true,
		},
	},
	MatchTest{
		Id:   "7a8fdc89-bdd8-5f81-8fe1-af8c5663907f",
		Desc: "godzilla_variant",
		In: Card{
			Name:      "Sprite Dragon",
			Variation: "Godzilla",
			Edition:   "Ikoria: Lair of Behemoths",
		},
	},
	MatchTest{
		Id:   "7a8fdc89-bdd8-5f81-8fe1-af8c5663907f",
		Desc: "godzilla_variant_alt_name",
		In: Card{
			Name:    "Dorat, the Perfect Pet",
			Edition: "Ikoria: Lair of Behemoths",
		},
	},
	MatchTest{
		Id:   "eaf748e5-ed61-504f-8799-bd52eb3ea9cd",
		Desc: "godzilla_variant_promo",
		In: Card{
			Name:      "Hangarback Walker",
			Variation: "Godzilla",
			Edition:   "ignored",
		},
	},
	MatchTest{
		Id:   "eaf748e5-ed61-504f-8799-bd52eb3ea9cd",
		Desc: "godzilla_variant_reskin_name",
		In: Card{
			Name:    "Mechagodzilla Battle Fortress",
			Edition: "Promo Cards",
		},
	},
	MatchTest{
		Id:   "fd9ff43c-fee9-5bdb-ab0a-d4c55f2178fd",
		Desc: "godzilla_variant_japanese",
		In: Card{
			Name:      "Dirge Bat",
			Variation: "Godzilla Series Japanese",
			Edition:   "Ikoria",
		},
	},
	MatchTest{
		Id:   "5df3565b-ab85-5cc7-83c4-9cd3bb5674da",
		Desc: "godzilla_variant_jpn_name",
		In: Card{
			Name:    "Mothra's Giant Cocoon",
			Edition: "Ikoria: Lair",
		},
	},

	// Arabian Nights different mana symbol
	MatchTest{
		Id:   "d429117f-4b10-5e66-ad2f-e233252a034a",
		Desc: "ARN_light_variant",
		In: Card{
			Name:      "Wyluli Wolf",
			Variation: "light circle",
			Edition:   "Arabian Nights",
		},
	},
	MatchTest{
		Id:   "184eabef-2042-5e2d-a2b3-96921e251de0",
		Desc: "ARN_dark_variant",
		In: Card{
			Name:      "Oubliette",
			Variation: "dark circle",
			Edition:   "Arabian Nights",
		},
	},
	MatchTest{
		Id:   "184eabef-2042-5e2d-a2b3-96921e251de0",
		Desc: "ARN_dark_variant_implied",
		In: Card{
			Name:      "Oubliette",
			Variation: "",
			Edition:   "Arabian Nights",
		},
	},
	MatchTest{
		Id:   "3dd0bd56-5340-5542-8457-646b9acd58ff",
		Desc: "ARN_no_variant",
		In: Card{
			Name:    "Abu Ja'far",
			Edition: "Arabian Nights",
		},
	},

	// Same-set variants
	MatchTest{
		Id:   "e8cad79a-2808-52a0-9504-469eab1d2486",
		Desc: "single_variant_with_no_special_tag",
		In: Card{
			Name:    "Will Kenrith",
			Edition: "Battlebond",
			Foil:    true,
		},
	},
	MatchTest{
		Id:   "a5d41107-bb39-5a10-ad0a-66513a31aa4d",
		Desc: "kaya_is_special",
		In: Card{
			Name:    "Kaya, Ghost Assassin",
			Edition: "Conspiracy: Take the Crown",
		},
	},
	MatchTest{
		Id:   "e0a1d531-00b4-587e-bfde-49de60e78f8e",
		Desc: "kaya_is_very_special",
		In: Card{
			Name:    "Kaya, Ghost Assassin",
			Edition: "Conspiracy: Take the Crown",
			Foil:    true,
		},
	},
	MatchTest{
		Id:   "b33e380c-0615-5559-93da-2ba3610d2b68",
		Desc: "too_many_variations",
		In: Card{
			Name:    "Tamiyo's Journal",
			Edition: "Shadows over Innistrad",
		},
	},
	MatchTest{
		Id:   "eb08862e-454a-510e-90f4-04e5adc335d5",
		Desc: "too_many_variations_foil",
		In: Card{
			Name:    "Tamiyo's Journal",
			Edition: "Shadows over Innistrad",
			Foil:    true,
		},
	},
	MatchTest{
		Id:   "9b53ce45-735c-5247-b744-9fcac2dbdc4b",
		Desc: "too_many_variations_what_did_i_say",
		In: Card{
			Name:      "Tamiyo's Journal",
			Variation: "Entry 546",
			Edition:   "Shadows over Innistrad",
		},
	},
	MatchTest{
		Id:   "7a9a79d8-c997-55b7-8a52-f8b54d5b60ee",
		Desc: "custom_variant",
		In: Card{
			Name:      "Urza's Tower",
			Variation: "Mountains",
			Edition:   "Chronicles",
		},
	},
	MatchTest{
		Id:   "825770fb-6760-5b0d-993f-0dfe58b55aa6",
		Desc: "number_with_suffix_in_variant",
		In: Card{
			Name:      "Arcane Denial",
			Variation: "22b",
			Edition:   "Alliances",
		},
	},
	MatchTest{
		Id:   "0aad334d-ac52-54b1-a044-1dd3641a6569",
		Desc: "one_funny_variation",
		In: Card{
			Name:      "Secret Base",
			Variation: "Version 2",
			Edition:   "Unstable",
		},
	},
	MatchTest{
		Id:   "4e62e057-32c2-55e6-bbe7-ddf0d0391d6b",
		Desc: "single_printintg_multiple_variants",
		In: Card{
			Name:      "Taste of Paradise",
			Variation: "TasteOfParadise",
			Edition:   "Alliances",
		},
	},
	MatchTest{
		Id:   "a7aacf96-8097-51a0-a50b-09d258edbc51",
		Desc: "artist_last_name_too_many_s",
		In: Card{
			Name:      "Simic Signet",
			Variation: "Mike Sass",
			Edition:   "Commander Anthology Volume II",
		},
	},

	// FNM promos (often confused with set promos)
	MatchTest{
		Id:   "0b5151a5-ada5-5045-a7a9-ecbe69593f69",
		Desc: "fnm_normal",
		In: Card{
			Name:      "Aether Hub",
			Variation: "FNM",
			Edition:   "ignored",
		},
	},
	MatchTest{
		Id:   "c72fd89d-60f2-59a5-9f77-8d138aebd38c",
		Desc: "fnm_plus_year",
		In: Card{
			Name:      "Goblin Warchief",
			Variation: "FNM 2016",
			Edition:   "ignored",
		},
	},
	MatchTest{
		Id:   "43805f80-743b-57ef-8f99-3ad19631120e",
		Desc: "fnm_with_promo_alias",
		In: Card{
			Name:      "Reliquary Tower",
			Variation: "FNM",
			Edition:   "Promo",
		},
	},
	MatchTest{
		Id:   "ad497aaa-2367-55c3-a599-ad038a3b1b7e",
		Desc: "nonfnm_with_fnm_alias",
		In: Card{
			Name:      "Reliquary Tower",
			Variation: "Promo",
			Edition:   "Promo",
		},
	},
	MatchTest{
		Id:   "588fba1e-9027-500d-ba10-83b79404a8a3",
		Desc: "nonfnm_wrong_info",
		In: Card{
			Name:      "Shanna, Sisay's Legacy",
			Variation: "FNM Foil",
			Edition:   "Promos: FNM",
		},
	},
	MatchTest{
		Id:   "8c9f9179-3a9f-51a0-8443-dfb19508b74c",
		Desc: "nonfnm_with_inverted_frame",
		In: Card{
			Name:      "Dovin's Veto",
			Variation: "FNM",
			Edition:   "Promo",
		},
	},

	// Arena
	MatchTest{
		Id:   "c136db2b-2e01-5e83-9eea-c40c05a24efe",
		Desc: "arena_normal",
		In: Card{
			Name:      "Enlightened Tutor",
			Variation: "Arena",
			Edition:   "ignored",
		},
	},
	MatchTest{
		Id:   "6667e7b7-f067-573f-ba33-69f836a02b47",
		Desc: "arena_with_year",
		In: Card{
			Name:      "Mountain",
			Variation: "Arena 1999",
			Edition:   "ignored",
		},
	},
	MatchTest{
		Id:   "abcb6c78-34f0-5654-a3bd-1116bb870f76",
		Desc: "arena_no_year",
		In: Card{
			Name:      "Forest",
			Variation: "Arena Foil - Mercadian Masques",
			Edition:   "ignored",
		},
	},
	MatchTest{
		Id:   "c8914a27-d1d2-5a60-95f7-c5517ad91caa",
		Desc: "arena_misprint",
		In: Card{
			Name:      "Island",
			Variation: "Arena 1999 misprint",
			Edition:   "ignored",
		},
	},
	MatchTest{
		Id:   "c8914a27-d1d2-5a60-95f7-c5517ad91caa",
		Desc: "arena_misprint_alt",
		In: Card{
			Name:      "Island",
			Variation: "Urza's Saga Arena Foil NO SYMBOL",
			Edition:   "ignored",
		},
	},
	MatchTest{
		Id:   "d879ae07-3ac7-5f03-a873-ce66d38fd61b",
		Desc: "arena_land_with_number",
		In: Card{
			Name:      "Forest",
			Variation: "Arena 2001 1",
			Edition:   "ignored",
		},
	},

	// Various promos
	MatchTest{
		Id:   "8ba3c9bf-e5b0-5008-9755-267c97c4b81f",
		Desc: "judge_normal",
		In: Card{
			Name:      "Tradewind Rider",
			Variation: "Judge",
			Edition:   "ignored",
		},
	},
	MatchTest{
		Id:   "a5af462e-afe3-53b5-9d9f-ab836e00b5ce",
		Desc: "judge_with_year",
		In: Card{
			Name:      "Vindicate",
			Variation: "Judge 2007",
			Edition:   "ignored",
		},
	},
	MatchTest{
		Id:   "818ab3df-bb0b-58b0-a15e-3daa53dafb15",
		Desc: "judge_with_set",
		In: Card{
			Name:      "Vampiric Tutor",
			Variation: "J18",
			Edition:   "Judge Promos",
		},
	},
	MatchTest{
		Id:   "4b8eb39c-6b31-5da8-9b13-9a53b4772d90",
		Desc: "sdcc",
		In: Card{
			Name:      "Liliana Vess",
			Variation: "2014 SDCC",
			Edition:   "ignored",
		},
	},
	MatchTest{
		Id:   "4b8eb39c-6b31-5da8-9b13-9a53b4772d90",
		Desc: "sdcc_extended_name",
		In: Card{
			Name:      "Liliana Vess",
			Variation: "San Diego Comic Con",
			Edition:   "ignored",
		},
	},
	MatchTest{
		Id:   "64614b74-7b1a-5b4d-b1b1-72af112cc287",
		Desc: "textless_normal",
		In: Card{
			Name:      "Fireball",
			Variation: "Textless",
			Edition:   "ignored",
		},
	},
	MatchTest{
		Id:   "d44d615b-2902-549a-8821-14845925556d",
		Desc: "gateway_normal",
		In: Card{
			Name:      "Lava Axe",
			Variation: "Gateway",
			Edition:   "ignored",
		},
	},
	MatchTest{
		Id:   "a94df3c2-5f76-58c5-b34a-28662290ebf7",
		Desc: "wpn_normal",
		In: Card{
			Name:      "Curse of Thirst",
			Variation: "WPN",
			Edition:   "ignored",
		},
	},
	MatchTest{
		Id:   "a94df3c2-5f76-58c5-b34a-28662290ebf7",
		Desc: "maybe_gateway_or_wpn",
		In: Card{
			Name:      "Curse of Thirst",
			Variation: "Gateway WPN",
			Edition:   "ignored",
		},
	},
	MatchTest{
		Id:   "4c17977b-5aca-5b0a-a456-4d4a0d5e42a1",
		Desc: "heros_path_promo",
		In: Card{
			Name:      "The Explorer",
			Variation: "Hero's Path",
			Edition:   "ignored",
		},
	},
	MatchTest{
		Id:   "0dd114e6-2a55-5338-9b6d-a4da134e4660",
		Desc: "duels_of_the_pw",
		In: Card{
			Name:      "Vigor",
			Variation: "Duels",
			Edition:   "ignored",
		},
	},
	MatchTest{
		Id:   "47ccb61a-3662-5dde-968d-24ae77e072c5",
		Desc: "duels_of_the_pw_with_year",
		In: Card{
			Name:      "Ogre Battledriver",
			Variation: "Duels of the Planeswalkers 2014",
			Edition:   "ignored",
		},
	},
	MatchTest{
		Id:   "58625b64-eba5-5e6b-aa79-c5f06039214d",
		Desc: "clash_pack",
		In: Card{
			Name:      "Temple of Mystery",
			Variation: "Clash Pack",
			Edition:   "ignored",
		},
	},
	MatchTest{
		Id:   "58625b64-eba5-5e6b-aa79-c5f06039214d",
		Desc: "clash_pack_alt",
		In: Card{
			Name:    "Temple of Mystery",
			Edition: "Clash Pack Promos",
		},
	},
	MatchTest{
		Id:   "218dee6b-c9ca-5d09-bec2-5517467db69b",
		Desc: "variation_has_no_useful_info",
		In: Card{
			Name:      "Zombie Apocalypse",
			Variation: "Some Kind of Promo",
			Edition:   "ignored",
			Foil:      true,
		},
	},
	MatchTest{
		Id:   "73ba75f6-9f27-5c75-9646-634e2f26bea7",
		Desc: "variation_has_no_useful_info_may_trigger_dupes_if_incorrectly_handled",
		In: Card{
			Name:      "Unclaimed Territory",
			Variation: "League Promo",
			Edition:   "ignored",
			Foil:      true,
		},
	},
	MatchTest{
		Id:   "40a038b4-87c9-552e-966b-d4dacac5ec38",
		Desc: "unknown_promo",
		In: Card{
			Name:      "Trueheart Duelist",
			Variation: "Game Day Extended",
			Edition:   "ignored",
		},
	},

	// Release cards
	MatchTest{
		Id:   "ec6031e8-8d33-547c-b446-3bd0112a931d",
		Desc: "release_but_it_is_a_promo",
		In: Card{
			Name:      "Valakut, the Molten Pinnacle",
			Variation: "Release Event",
			Edition:   "ignored",
			Foil:      true,
		},
	},
	MatchTest{
		Id:   "674e74e5-e5ad-56b4-acc7-ffd550d9b36d",
		Desc: "release_but_it_is_from_launch_parties",
		In: Card{
			Name:      "Vexing Shusher",
			Variation: "Release Event",
			Edition:   "ignored",
			Foil:      true,
		},
	},
	MatchTest{
		Id:   "dbfdf3b7-4eb9-5a78-b1fb-6db73491d59f",
		Desc: "release_events",
		In: Card{
			Name:      "Shriekmaw",
			Variation: "Release Event",
			Edition:   "ignored",
			Foil:      true,
		},
	},
	MatchTest{
		Id:   "7f43c4bb-2b38-5a05-9ce2-e2042009af0e",
		Desc: "release_but_there_is_a_prerelease_too",
		In: Card{
			Name:      "Identity Thief",
			Variation: "Release Event",
			Edition:   "ignored",
			Foil:      true,
		},
	},
	MatchTest{
		Id:   "734e9eb3-d86b-564d-9553-4ffb48b5e13a",
		Desc: "prerelease_but_there_is_a_release_too",
		In: Card{
			Name:      "Identity Thief",
			Variation: "Prerelease Event",
			Edition:   "ignored",
			Foil:      true,
		},
	},
	MatchTest{
		Id:   "1498e66b-3ddf-54e1-a4c4-a76a664a4bf1",
		Desc: "non_release_non_prerelease_version",
		In: Card{
			Name:    "Identity Thief",
			Edition: "Eldritch Moon",
		},
	},
	MatchTest{
		Id:   "7f43c4bb-2b38-5a05-9ce2-e2042009af0e",
		Desc: "release_too_much_info",
		In: Card{
			Name:      "Identity Thief",
			Variation: "Eldritch Moon Launch Foil 22 July 2016",
			Edition:   "Promos: Miscellaneous",
		},
	},
	MatchTest{
		Id:   "412013dc-8a6a-54f8-beb7-77b56baa5057",
		Desc: "launch_in_the_set_itself",
		In: Card{
			Name:      "Scholar of the Lost Trove",
			Variation: "Launch Promo Foil",
			Edition:   "Jumpstart",
		},
	},
	MatchTest{
		Id:   "bcd861de-1206-5df9-a2c5-e6a3464b37b8",
		Desc: "launch_in_the_set_wrong_edition",
		In: Card{
			Name:      "Chord of Calling",
			Variation: "Launch",
			Edition:   "ignored",
		},
	},

	// Buy-a-Box promo
	MatchTest{
		Id:   "ad5c0740-144d-58fd-8fde-e2f3aee52fc8",
		Desc: "bab_marked_as_promo_but_it_is_really_in_the_set",
		In: Card{
			Name:      "Impervious Greatwurm",
			Variation: "BIBB",
			Edition:   "ignored",
			Foil:      true,
		},
	},
	MatchTest{
		Id:   "eaae1532-4362-54f9-8855-477ce59eab99",
		Desc: "bab_marked_as_promo_but_it_is_really_in_the_set_set_is_not_expansion",
		In: Card{
			Name:      "Flusterstorm",
			Variation: "buy a box",
			Edition:   "ignored",
		},
	},
	MatchTest{
		Id:   "c36c2a16-d6b7-51c9-8c94-e505ab618d66",
		Desc: "bab_marked_as_promo_but_it's_really_in_the_set_set_is_core",
		In: Card{
			Name:      "Rienne Angel of Rebirth",
			Variation: "M20 BIBB",
			Edition:   "ignored",
		},
	},
	MatchTest{
		Id:   "cf72208c-ac68-55d4-bb68-66487d682749",
		Desc: "bab_old_style_it_is_in_Promos",
		In: Card{
			Name:      "Sylvan Caryatid",
			Variation: "buy-a-box",
			Edition:   "ignored",
			Foil:      true,
		},
	},
	MatchTest{
		Id:   "e548ac74-e9a5-5a0e-9401-32cc5a74dc6b",
		Desc: "bab_but_also_pro_tour",
		In: Card{
			Name:      "Surgical Extraction",
			Variation: "BIBB Promo",
			Edition:   "ignored",
			Foil:      true,
		},
	},
	MatchTest{
		Id:   "9aa8535c-407e-5c43-a8bd-1adf83dbafcb",
		Desc: "bab_but_also_in_normal_set",
		In: Card{
			Name:      "Mirran Crusader",
			Variation: "Buy-a-Box",
			Edition:   "ignored",
			Foil:      true,
		},
	},
	MatchTest{
		Id:   "daa40fb0-cd50-5c75-bee6-1b6adafbb590",
		Desc: "bab_in_separate_set_with_wrong_info",
		In: Card{
			Name:      "Growing Rites of Itlimoc",
			Variation: "buy-a-box",
			Edition:   "Ixalan Promos",
		},
	},

	// Bundle promo
	MatchTest{
		Id:   "372145ed-c7a8-5494-b1e6-6f5aec74d7c0",
		Desc: "nonbundle_in_the_same_set",
		In: Card{
			Name:    "Piper of the Swarm",
			Edition: "Throne of Eldraine",
		},
	},
	MatchTest{
		Id:   "37fdc0d7-976d-5e5e-af6c-ee6b50795454",
		Desc: "bundle_in_the_same_set",
		In: Card{
			Name:      "Piper of the Swarm",
			Variation: "Bundle",
			Edition:   "Throne of Eldraine",
		},
	},
	MatchTest{
		Id:   "37fdc0d7-976d-5e5e-af6c-ee6b50795454",
		Desc: "bundle_in_the_same_set_but_unknown_set",
		In: Card{
			Name:      "Piper of the Swarm",
			Variation: "Bundle Promo",
			Edition:   "ignored",
		},
	},
	MatchTest{
		Id:   "cb8c4745-f3d2-51d6-87ab-71612430ae5f",
		Desc: "nonbundle_in_the_same_set_but_special_version",
		In: Card{
			Name:      "Piper of the Swarm",
			Variation: "Extended Art",
			Edition:   "Throne of Eldraine",
		},
	},
	MatchTest{
		Id:   "35c17fee-50a9-5273-ba96-492b156cbfff_f",
		Desc: "magicfest_normal",
		In: Card{
			Name:      "Path to Exile",
			Variation: "Magic Fest",
			Edition:   "ignored",
			Foil:      true,
		},
	},
	MatchTest{
		Id:   "35c17fee-50a9-5273-ba96-492b156cbfff",
		Desc: "magicfest_textless",
		In: Card{
			Name:      "Path to Exile",
			Variation: "MagicFest Textless",
			Edition:   "ignored",
		},
	},
	MatchTest{
		Id:   "a95b64f0-9976-592f-b7ad-aa09158aa63c",
		Desc: "bfz_std_with_wrong_info",
		In: Card{
			Name:      "Sunken Hollow",
			Variation: "alt art",
			Edition:   "Dominaria",
		},
	},
	MatchTest{
		Id:   "b8551573-9c09-5dc2-a440-d260fcbe6fad",
		Desc: "unstable_letter_variant",
		In: Card{
			Name:      "Very Cryptic Command",
			Variation: "E Counter/Return/Untap/Roll",
			Edition:   "Unstable",
		},
	},

	// Homelands and Fallen Empires
	MatchTest{
		Id:   "64a0d121-4a0b-5015-bcf2-985a996f196f",
		Desc: "homelands_flavor",
		In: Card{
			Name:      "Abbey Matron",
			Variation: "Quote Halina, Dwarven Trader",
			Edition:   "Homelands",
		},
	},
	MatchTest{
		Id:   "5a8aacf9-dfb4-5d3a-a36e-87b2ef6cca43",
		Desc: "homelands_flavor_alt",
		In: Card{
			Name:      "Folk of An-Havva",
			Variation: "Quote Joskun, An-Havva Constable",
			Edition:   "Homelands",
		},
	},
	MatchTest{
		Id:   "2f01f5b8-6fe8-510e-8c35-3b209dbd41ce",
		Desc: "homelands_flavor_with_extra",
		In: Card{
			Name:      "Memory Lapse",
			Variation: "Quote Chandler, Female Art",
			Edition:   "Homelands",
		},
	},
	MatchTest{
		Id:   "a860ebd5-1f8a-54b8-bf83-473fd7594a15",
		Desc: "fem_artist",
		In: Card{
			Name:      "Armor Thrull",
			Variation: "Jeff A. Menges",
			Edition:   "Fallen Empires",
		},
	},
	MatchTest{
		Id:   "d4324874-0c9d-5d2e-936c-d24b8f5de060",
		Desc: "fem_artist_incomplete",
		In: Card{
			Name:      "Icatian Javelineers",
			Variation: "Melissa Benson",
			Edition:   "Fallen Empires",
		},
	},
	MatchTest{
		Id:   "46015c37-2fbb-5d1c-b4f5-8378e777fb6f",
		Desc: "fem_variant_is_number_suffix",
		In: Card{
			Name:      "Homarid Warrior",
			Variation: "B",
			Edition:   "Fallen Empires",
		},
	},
	MatchTest{
		Id:   "1e5e8355-c4c0-552c-a000-ad80e603844e",
		Desc: "fem_variant_is_polluted",
		In: Card{
			Name:      "Basal Thrull",
			Variation: "Artist Phil Foglio",
			Edition:   "Fallen Empires",
		},
	},

	// Duel Decks
	MatchTest{
		Id:   "0ecf4a89-44f9-5c9a-9ecd-422702e44ef2",
		Desc: "duel_decks_variant",
		In: Card{
			Name:    "Goblin Rabblemaster",
			Edition: "DD: Merfolk vs Goblins",
		},
	},
	MatchTest{
		Id:   "f4ca3eba-d073-5a83-8732-c5d465b06a11",
		Desc: "duel_decks_variant_with_number",
		In: Card{
			Name:      "Forest",
			Variation: "#38",
			Edition:   "DD: Zendikar vs. Eldrazi",
		},
	},
	MatchTest{
		Id:   "15ce106a-8fde-5848-99b3-21eceb764be0",
		Desc: "duel_decks_variant_with_mb1_tag",
		In: Card{
			Name:      "Elvish Warrior",
			Variation: "Mystery Booster",
			Edition:   "Elves vs. Goblins",
		},
	},
	MatchTest{
		Id:   "12420bff-ce07-5ce1-8f72-03af7df3f1ef",
		Desc: "dda_in_variation",
		In: Card{
			Name:      "Flamewave Invoker",
			Variation: "Jace vs Chandra",
			Edition:   "Duel Decks Anthology",
		},
	},
	MatchTest{
		Id:   "3626181e-e565-5dfa-8c59-2d02e47beff3",
		Desc: "dda_in_variation_inverted",
		In: Card{
			Name:      "Flamewave Invoker",
			Variation: "Goblins vs Elves",
			Edition:   "Duel Decks Anthology",
		},
	},
	MatchTest{
		Id:   "8882a58b-a123-5c16-bd10-a86c6faf4155",
		Desc: "duel_decks_variant_jpn",
		In: Card{
			Name:      "Jace Beleren",
			Variation: "Japanese",
			Edition:   "Duel Decks: Jace vs. Chandra",
		},
	},

	// Deckmasters variants
	MatchTest{
		Id:   "3bc39616-8759-51eb-ab74-cf8271c339f6",
		Desc: "deckmasters_number_in_variation_with_other_text",
		In: Card{
			Name:      "Phyrexian War Beast",
			Variation: "37A Propeller Right",
			Edition:   "DKM",
		},
	},
	MatchTest{
		Id:   "3bc39616-8759-51eb-ab74-cf8271c339f6",
		Desc: "deckmasters_use_first_if_empty",
		In: Card{
			Name:      "Phyrexian War Beast",
			Variation: "",
			Edition:   "DKM",
		},
	},
	MatchTest{
		Id:   "7a6610cf-7d53-5e4b-bb76-9be1708d3892",
		Desc: "deckmasters_variant_foil",
		In: Card{
			Name:      "Incinerate",
			Variation: "Foil",
			Edition:   "Deckmasters",
			Foil:      true,
		},
	},
	MatchTest{
		Id:   "0b77872e-c4eb-54f4-8f63-57fcf68e682a",
		Desc: "deckmasters_variant_non_foil",
		In: Card{
			Name:    "Incinerate",
			Edition: "Deckmasters",
		},
	},
	MatchTest{
		Id:   "982dabaa-9e5c-5a15-9fc2-cb4de4f13f11",
		Desc: "variation_deckmasters_foil_but_untagged",
		In: Card{
			Name:      "Icy Manipulator",
			Variation: "Promo",
			Edition:   "Deckmasters",
		},
	},

	// Champs
	MatchTest{
		Id:   "59efdba1-25d8-56b0-8b82-d17839e19ff3",
		Desc: "states_but_is_gateway",
		In: Card{
			Name:      "Dauntless Dourbark",
			Variation: "2008 States",
			Edition:   "ignored",
		},
	},
	MatchTest{
		Id:   "79cbb056-116f-53b6-83ff-6af010bf6e49",
		Desc: "champs_and_states",
		In: Card{
			Name:      "Voidslime",
			Variation: "Champs",
			Edition:   "ignored",
		},
	},
	MatchTest{
		Id:   "0c99ba44-73d6-535d-b68e-3dcc690aadd6",
		Desc: "not_champs",
		In: Card{
			Name:      "Ghalta, Primal Hunger",
			Variation: "Champs / States",
			Edition:   "ignored",
		},
	},

	// IDW and Comic promos
	MatchTest{
		Id:   "095baf23-4c83-5be7-b12f-67fbb7deddb4",
		Desc: "book_promo",
		In: Card{
			Name:      "Jace Beleren",
			Variation: "Book",
			Edition:   "ignored",
		},
	},
	MatchTest{
		Id:   "55a4b8cf-4e8e-5c2c-abf6-3b1ed592d323",
		Desc: "idw_normal",
		In: Card{
			Name:      "Wash Out",
			Variation: "IDW",
			Edition:   "ignored",
		},
	},
	MatchTest{
		Id:   "d1f8693c-4cdc-5558-b678-5fcefc0d220d",
		Desc: "idw_also_magazine",
		In: Card{
			Name:      "Duress",
			Variation: "IDW Promo",
			Edition:   "ignored",
		},
	},
	MatchTest{
		Id:   "c964334d-1153-55da-87d2-b986282af243",
		Desc: "idw_but_also_magazine",
		In: Card{
			Name:      "Duress",
			Variation: "Japanese Magazine Promo",
			Edition:   "ignored",
		},
	},
	MatchTest{
		Id:   "2a0a421a-418d-519c-b316-04e8071c36d7",
		Desc: "japanese_magazine_insert",
		In: Card{
			Name:      "Shivan Dragon",
			Variation: "Japanese Gotta Comic Promo",
			Edition:   "ignored",
		},
	},
	MatchTest{
		Id:   "d1f8693c-4cdc-5558-b678-5fcefc0d220d",
		Desc: "duress_dupe_idw",
		In: Card{
			Name:      "Duress",
			Variation: "Comic promo",
			Edition:   "ignored",
		},
	},
	MatchTest{
		Id:   "7d5ea62e-47af-5802-a24f-8090ba028f18",
		Desc: "magazine_insert_not_tagged_jpn",
		In: Card{
			Name:    "Lava Coil",
			Edition: "Magazine Inserts",
		},
	},
	MatchTest{
		Id:   "cf72208c-ac68-55d4-bb68-66487d682749",
		Desc: "magazine_insert_in_generic_promo",
		In: Card{
			Name:    "Sylvan Caryatid",
			Edition: "Magazine Inserts",
		},
	},
	MatchTest{
		Id:   "857618a8-7f5f-5a1f-8435-b210a8e8250c",
		Desc: "magazine_insert_as_pres",
		In: Card{
			Name:    "Beast Whisperer",
			Edition: "Magazine Inserts",
		},
	},

	// Core sets
	MatchTest{
		Id:   "9a72fcfd-5d43-55fb-8ade-8476d38f506b",
		Desc: "coreset_normal",
		In: Card{
			Name:    "Guttersnipe",
			Edition: "Core Set 2019 / M19",
		},
	},
	MatchTest{
		Id:   "5561b2e3-b0fd-5f2c-be7e-0ccd448bb8e2",
		Desc: "coreset_confusing_promo",
		In: Card{
			Name:      "Naya Sojourners",
			Variation: "Magic 2010 Game Day",
			Edition:   "Promo Magic 2010 Game Day",
		},
	},

	// WCD
	MatchTest{
		Id:   "d36e37ff-9b07-50a4-9cb1-451caa554159",
		Desc: "wcd_pick_the_first_one_if_not_enough_info",
		In: Card{
			Name:      "Ancient Tomb",
			Variation: "Tokyo 1999 - Not Tournament Legal",
			Edition:   "World Championships",
		},
	},
	MatchTest{
		Id:   "d50a6f25-d392-55e6-a34e-f83ad2b89c33",
		Desc: "wcd_with_number",
		In: Card{
			Name:      "Plains",
			Variation: "8th Edition 332 Julien Nuijten 2004",
			Edition:   "World Championship",
		},
	},
	MatchTest{
		Id:   "ae03bc21-485c-5175-9a41-7fc5421d62ef",
		Desc: "wcd_with_variant",
		In: Card{
			Name:      "Memory Lapse",
			Variation: "Statue A Sideboard Shawn Hammer Regnier",
			Edition:   "World Championship",
		},
	},
	MatchTest{
		Id:   "09ec16b3-b9b3-5be4-ba2c-7b52dc330b05",
		Desc: "wcd_with_variant_embedded_in_number",
		In: Card{
			Name:      "Plains",
			Variation: "Odyssey 331 Brian Kibler 2002",
			Edition:   "World Championship",
		},
	},
	MatchTest{
		Id:   "70dc718c-b23b-51a8-b9d8-36433a438d79",
		Desc: "wcd_with_player_name_aliasing",
		In: Card{
			Name:      "Cursed Scroll",
			Variation: "Matt Linde 1999",
			Edition:   "World Championship",
		},
	},
	MatchTest{
		Id:   "19a00cd2-3fdd-5b82-9eab-b561d94e362a",
		Desc: "wcd_with_correct_number",
		In: Card{
			Name:      "Strip Mine",
			Variation: "ll363",
			Edition:   "Pro Tour Collector Set",
		},
	},
	MatchTest{
		Id:   "fbd3ea3e-80cb-5f55-b41e-5730aff46c6a",
		Desc: "wcd_with_correct_original_number",
		In: Card{
			Name:      "Order of Leitbur",
			Variation: "1996 Bertrand Lestree FEM Male 16b",
			Edition:   "World Championship Decks",
		},
	},
	MatchTest{
		Id:   "df552574-7dbd-5bc3-bcd9-ba8d9c745b41",
		Desc: "wcd_with_correct_sideboard_number",
		In: Card{
			Name:      "Krosan Reclamation",
			Variation: "dz122sb",
			Edition:   "WCD 2003: Daniel Zink",
		},
	},
	MatchTest{
		Id:   "0e1e5130-d97a-56f5-9188-5df2f9e965e1",
		Desc: "wcd_only_the_year",
		In: Card{
			Name:      "Karplusan Forest",
			Variation: "Brussels, August 2000",
			Edition:   "World Championships 2000",
		},
	},
	MatchTest{
		Id:   "15d4bf85-0746-59c5-83e7-1f40646672b8",
		Desc: "wcd_minimal_info",
		In: Card{
			Name:    "City of Brass",
			Edition: "WCD 1997: Paul McCabe",
		},
	},

	// Foil-only special category
	MatchTest{
		Id:   "a05c82be-c929-5720-a509-7b9f51156db9_f",
		Desc: "PLS_foil_only_booster_normal",
		In: Card{
			Name:    "Skyship Weatherlight",
			Edition: "Planeshift",
			Foil:    true,
		},
	},
	MatchTest{
		Id:   "dcc4ee11-6a61-55f0-966a-d19732010ffa",
		Desc: "PLS_foil_only_booster_alternate",
		In: Card{
			Name:      "Skyship Weatherlight",
			Variation: "Alternate Art",
			Edition:   "Planeshift",
			Foil:      true,
		},
	},
	MatchTest{
		Id:   "4b375581-c43d-5bd1-b990-a0fa496b8262",
		Desc: "10E_foil_only_booster_clean",
		In: Card{
			Name:    "Time Stop",
			Edition: "Tenth Edition",
			Foil:    true,
		},
	},
	MatchTest{
		Id:   "3b77bb52-4181-57f5-b3cd-f3a15b95aa29_f",
		Desc: "10E_foil_only_booster_normal",
		In: Card{
			Name:    "Angelic Chorus",
			Edition: "Tenth Edition",
			Foil:    true,
		},
	},

	// Portal variants
	MatchTest{
		Id:   "73b7e8ec-6b0c-5c35-92ca-dc0cd1156456",
		Desc: "portal_starter_deck",
		In: Card{
			Name:      "Blaze",
			Variation: "reminder text",
			Edition:   "Portal",
		},
	},
	MatchTest{
		Id:   "e410cea1-ca02-5fce-b4e9-ad8b6dfb6a30",
		Desc: "portal_starter_deck_alt",
		In: Card{
			Name:      "Raging Goblin",
			Variation: "No flavor text",
			Edition:   "Portal",
		},
	},
	MatchTest{
		Id:   "593d4dcc-f98f-5f35-84a0-a71014a9b3b4",
		Desc: "portal_demo_game",
		In: Card{
			Name:      "Cloud Pirates",
			Variation: "reminder text",
			Edition:   "Portal",
		},
	},
	MatchTest{
		Id:   "9c6f7fe2-544c-5d92-b5c7-777a5c944561",
		Desc: "portal_simplified_chinese_no_alias",
		In: Card{
			Name:    "Assassin's Blade",
			Edition: "Portal",
		},
	},

	// Mystery Booster - The List debacle
	MatchTest{
		Id:   "b66e5217-c523-5d72-b94b-252787083ec6",
		Desc: "mb1list_only_in_list",
		In: Card{
			Name:    "Scroll Rack",
			Edition: "Mystery Booster The List",
		},
	},
	MatchTest{
		Id:   "0c466efd-2923-5e49-b0d7-ae69c279296a",
		Desc: "mb1list_only_in_mb1",
		In: Card{
			Name:    "Mana Crypt",
			Edition: "Mystery Booster The List",
		},
	},
	MatchTest{
		Id:   "cafa10aa-9fa7-5b41-99d7-cbfd51eae016",
		Desc: "mb1list_in_both_undistinguishable",
		In: Card{
			Name:    "Living Death",
			Edition: "Mystery Booster The List",
		},
	},
	MatchTest{
		Id:   "65166320-c130-50e5-a78a-8f60bc8c453c",
		Desc: "mb1list_in_both_list_tagged",
		In: Card{
			Name:      "Imperious Perfect",
			Variation: "Extended Art",
			Edition:   "Mystery Booster The List",
		},
	},
	MatchTest{
		Id:   "ab89a3c7-9262-5060-bb1f-cbc8630613b2",
		Desc: "mb1list_in_both_mb1_untagged",
		In: Card{
			Name:    "Imperious Perfect",
			Edition: "Mystery Booster The List",
		},
	},
	MatchTest{
		Id:   "2ea9464a-c6f9-5e0c-a0bf-55f89fa42252",
		Desc: "mb1list_in_both_tagged_mb1",
		In: Card{
			Name:      "Bojuka Bog",
			Variation: "Worldwake",
			Edition:   "Mystery Booster The List",
		},
	},
	MatchTest{
		Id:   "706ca274-e840-56e6-b252-4d7da8f13063",
		Desc: "mb1list_in_both_tagged_plist",
		In: Card{
			Name:      "Bojuka Bog",
			Variation: "Commander 2018",
			Edition:   "Mystery Booster The List",
		},
	},
	MatchTest{
		Id:   "a057d0d7-11fb-5d77-9217-96d8cd26dd67",
		Desc: "mb1list_in_both_mb1_foil",
		In: Card{
			Name:      "Helix Pinnacle",
			Variation: "Foil",
			Edition:   "Mystery Booster The List",
			Foil:      true,
		},
	},
	MatchTest{
		Id:   "9ef93f93-52a8-5807-8c47-aa32048c25ce",
		Desc: "mb1list_in_both_list_nonfoil",
		In: Card{
			Name:      "Helix Pinnacle",
			Variation: "Non-Foil",
			Edition:   "Mystery Booster The List",
		},
	},
	MatchTest{
		Id:   "7cac5a8b-f4cc-5169-acf1-f4e0b2b553b6",
		Desc: "mb1list_explicit_list",
		In: Card{
			Name:    "Scourge of the Throne",
			Edition: "The List",
		},
	},
	MatchTest{
		Id:   "d26069ff-4e63-547f-a64b-cbf2d83cdbdf",
		Desc: "mb1list_explicit_variant_list",
		In: Card{
			Name:      "Thought Vessel",
			Variation: "The List",
			Edition:   "Ignored",
		},
	},

	// Triple card because of double misprint
	MatchTest{
		Id:   "89d7f595-d549-526a-ba83-6f976f6c542c",
		Desc: "prerelease_with_multiple_misprints_dark",
		In: Card{
			Name:      "Laquatus's Champion",
			Variation: "dark misprint",
			Edition:   "Prerelease Events",
		},
	},
	MatchTest{
		Id:   "7c3cac05-a2a1-5e74-8a5f-dc73a3f3b013",
		Desc: "prerelease_with_multiple_misprints_nondark",
		In: Card{
			Name:      "Laquatus's Champion",
			Variation: "misprint",
			Edition:   "Prerelease Events",
		},
	},
	MatchTest{
		Id:   "5351f6ca-4c00-57c4-ba43-272f6eb622df",
		Desc: "prerelease_with_multiple_misprints_normal",
		In: Card{
			Name:      "Laquatus's Champion",
			Variation: "Prerelease",
			Edition:   "ignored",
		},
	},

	// CMR duplicates
	MatchTest{
		Id:   "0dfd6f8d-f7e6-5aba-ab2b-88fb15e8aeb3",
		Desc: "cmr_duplicates_main",
		In: Card{
			Name:    "Return to Dust",
			Edition: "Commander Legends",
		},
	},
	MatchTest{
		Id:   "5d1ddb47-9191-5527-a695-458bf67ec2f9",
		Desc: "cmr_duplicates_decks_actually",
		In: Card{
			Name:      "Return to Dust",
			Variation: "not in draft",
			Edition:   "Commander Legends",
		},
	},
	MatchTest{
		Id:   "57ec163e-2368-5a73-abef-f1851339ae58",
		Desc: "cmr_duplicates_ea",
		In: Card{
			Name:      "Return to Dust",
			Variation: "Extended Art",
			Edition:   "Commander Legends",
		},
	},
	MatchTest{
		Id:   "cc3d585e-6835-5442-a227-7cdc23b6a839",
		Desc: "cmr_duplicates_decks_notreally",
		In: Card{
			Name:      "Acidic Slimes",
			Variation: "anything",
			Edition:   "Commander Legends",
		},
	},
	MatchTest{
		Id:   "0fba90e8-f8a6-534b-8d0c-096f5ff4d116_f",
		Desc: "cmr_normal_foil",
		In: Card{
			Name:    "Yurlok of Scorch Thrash",
			Edition: "Commander Legends",
			Foil:    true,
		},
	},
	MatchTest{
		Id:   "acaf7712-434a-5a18-a1c9-1f92081a7ca1",
		Desc: "cmr_etched_foil_one",
		In: Card{
			Name:      "Yurlok of Scorch Thrash",
			Variation: "Etched Foil",
			Edition:   "Commander Legends",
			Foil:      true,
		},
	},
	MatchTest{
		Id:   "3777b628-6807-51dc-855a-12224af0f4a8",
		Desc: "cmr_etched_foil_two",
		In: Card{
			Name:      "Nevinyrral, Urborg Tyrant",
			Variation: "Foil-Etched",
			Edition:   "Commander Legends",
		},
	},
	MatchTest{
		Id:   "0f8e79b8-c429-5e9d-bb51-849c114bf5de",
		Desc: "cmr_bab",
		In: Card{
			Name:      "Mana Confluence",
			Variation: "721",
			Edition:   "Buy A Box",
		},
	},

	MatchTest{
		Id:   "9fce4a02-01ec-5dc3-a3e9-92646a69e6fd",
		Desc: "planechase_planes",
		In: Card{
			Name:    "The Hippodrome",
			Edition: "Planechase",
		},
	},
	MatchTest{
		Id:   "a82e3d86-2ffb-593d-88b5-16db79ce2025",
		Desc: "boosterfun_variant",
		In: Card{
			Name:    "Shadowspear",
			Edition: "Theros Beyond Death: Extras",
		},
	},
	MatchTest{
		Id:   "4986aa22-de4e-5f16-88bd-f378f8071905",
		Desc: "boosterfun_variant_but_custom_setsize",
		In: Card{
			Name:    "Abrade",
			Edition: "Commander Legends: Variants",
		},
	},
	MatchTest{
		Id:   "3aa21396-8dcd-5f9a-9243-677182293d4a",
		Desc: "boosterfun_variant_but_needs_tags",
		In: Card{
			Name:    "Pollywog Symbiote",
			Edition: "Ikoria: Lair of Behemoths: Variants",
		},
	},
	MatchTest{
		Id:   "e1b8fc9f-e974-51cd-972f-624b3e9d0c1b",
		Desc: "intropack_dedup_normal",
		In: Card{
			Name:    "Wind Drake",
			Edition: "Kaladesh",
		},
	},
	MatchTest{
		Id:   "992c07f2-f9b9-5c7a-a66c-efd497c3c9b2",
		Desc: "intropack_dedup_intro",
		In: Card{
			Name:      "Wind Drake",
			Variation: "Intro Pack",
			Edition:   "Kaladesh",
		},
	},

	// Best effort mode, when the user gives up and sets the Promo
	// flag, this small subsets of promos should still be supported
	MatchTest{
		Id:   "f87da7e0-10e5-5cae-917e-6e5a92f7880c",
		Desc: "untagged_promo_anypromoset",
		In: Card{
			Name:    "Acidic Slime",
			Edition: "Promo",
			Promo:   true,
		},
	},
	MatchTest{
		Id:   "21beeb79-a0e0-5809-94df-25e44b45220f",
		Desc: "untagged_promo_clashpack",
		In: Card{
			Name:    "Courser of Kruphix",
			Edition: "Promo",
			Promo:   true,
		},
	},
	MatchTest{
		Id:   "3f9a0369-5fe7-5aee-85fe-3cfaacd275af",
		Desc: "untagged_promo_promopack_in_set",
		In: Card{
			Name:    "Bloodchief's Thirst",
			Edition: "Promo",
			Promo:   true,
		},
	},
	MatchTest{
		Id:   "c6c2a57b-1422-54f5-be68-1d64328e3679",
		Desc: "untagged_promo_expansion",
		In: Card{
			Name:    "Adorned Pouncer",
			Edition: "Promo",
			Promo:   true,
		},
	},
	MatchTest{
		Id:   "faf229ab-b135-5c37-abb7-dbf68a1dad9b",
		Desc: "untagged_promo_coreset",
		In: Card{
			Name:    "Alhammarret, High Arbiter",
			Edition: "Promo",
			Promo:   true,
		},
	},
	MatchTest{
		Id:   "bcd861de-1206-5df9-a2c5-e6a3464b37b8",
		Desc: "untagged_promo_masters",
		In: Card{
			Name:    "Chord of Calling",
			Edition: "Promo",
			Promo:   true,
		},
	},
	MatchTest{
		Id:   "2a113727-7f6e-5ab7-9dce-5e6f19ac62c8",
		Desc: "untagged_promo_skip_jpn",
		In: Card{
			Name:    "Dark Ritual",
			Edition: "Promo",
			Promo:   true,
		},
	},
	MatchTest{
		Id:   "0f8e79b8-c429-5e9d-bb51-849c114bf5de",
		Desc: "untagged_promo_draft",
		In: Card{
			Name:    "Mana Confluence",
			Edition: "Promo",
			Promo:   true,
		},
	},
	MatchTest{
		Id:   "dcc4ee11-6a61-55f0-966a-d19732010ffa",
		Desc: "untagged_promo_nonpromo",
		In: Card{
			Name:    "Skyship Weatherlight",
			Edition: "Promo",
			Promo:   true,
		},
	},
	MatchTest{
		Id:   "37fdc0d7-976d-5e5e-af6c-ee6b50795454",
		Desc: "untagged_promo_bundle",
		In: Card{
			Name:    "Piper of the Swarm",
			Edition: "Promo",
			Promo:   true,
		},
	},
	MatchTest{
		Id:   "63ef8762-36b8-5891-97ae-744d0e5d3f78",
		Desc: "untagged_promo_bab",
		In: Card{
			Name:    "Realmwalker",
			Edition: "Promo",
			Promo:   true,
		},
	},
	MatchTest{
		Id:   "1cb28add-17dd-5e08-9d02-7d6c2275a1fa",
		Desc: "untagged_promo_not_boxtopper",
		In: Card{
			Name:    "Batterskull",
			Edition: "Promo",
			Promo:   true,
		},
	},
	MatchTest{
		Id:   "7165c42c-7210-5564-b5a9-e666de59034a",
		Desc: "tagged_promopack_in_core",
		In: Card{
			Name:      "Eliminate",
			Variation: "Promo Pack",
			Edition:   "Promo Pack",
		},
	},

	// Non-English custom sets
	MatchTest{
		Id:   "dcb15a6c-4629-59a8-bdbf-a17eb57d6fe3_jpn",
		Desc: "custom_language_number",
		In: Card{
			Name:      "Urza's Power Plant",
			Variation: "115d",
			Edition:   "Chronicles Japanese",
		},
	},
	MatchTest{
		Id:   "9ecb374f-fc0d-54d3-a26b-f590306026a0_ita",
		Desc: "custom_language_tag",
		In: Card{
			Name:    "Sylvan Library",
			Edition: "Legends Italian",
		},
	},
	MatchTest{
		Id:   "9ecb374f-fc0d-54d3-a26b-f590306026a0_ita",
		Desc: "custom_language_id_lookup",
		In: Card{
			Id: "9ecb374f-fc0d-54d3-a26b-f590306026a0_ita",
		},
	},

	MatchTest{
		Id:   "f2aa69f2-db69-5564-8661-d51855c6796a",
		Desc: "sta_double_variant_normal",
		In: Card{
			Name:    "Demonic Tutor",
			Edition: "Mystical Archive",
		},
	},
	MatchTest{
		Id:   "b076bded-cee7-5f26-b9d8-3540918f68fa",
		Desc: "sta_double_variant_jpn",
		In: Card{
			Name:      "Demonic Tutor",
			Variation: "JPN",
			Edition:   "Strixhaven Mystical Archive",
		},
	},
	MatchTest{
		Id:   "b076bded-cee7-5f26-b9d8-3540918f68fa",
		Desc: "sta_double_variant_jpn_number",
		In: Card{
			Name:      "Demonic Tutor",
			Variation: "JPN 90",
			Edition:   "Strixhaven Mystical Archives",
		},
	},
	MatchTest{
		Id:   "3dc11b06-acdc-536d-8acf-af02c5672ab3",
		Desc: "sta_normal_etched",
		In: Card{
			Name:      "Demonic Tutor",
			Variation: "Foil Etched",
			Edition:   "Strixhaven Mystical Archives",
			Foil:      true,
		},
	},
	MatchTest{
		Id:   "a974d871-39bc-5bcb-947b-853836d3a6f1",
		Desc: "sta_jpn_etched",
		In: Card{
			Name:      "Demonic Tutor",
			Variation: "JP Alternate Art Foil Etched",
			Edition:   "Strixhaven Mystical Archives",
			Foil:      true,
		},
	},
	MatchTest{
		Id:   "6c82bc6b-0e97-5283-ae43-e5d16f53e19a",
		Desc: "sta_double_variant_japanese_noforeigndata",
		In: Card{
			Name:      "Negate",
			Variation: "Japanese",
			Edition:   "Strixhaven Mystical Archive",
		},
	},
	MatchTest{
		Id:   "06d4b87d-3030-5331-b7b7-bf7024512ff5",
		Desc: "sta_number_and_etched",
		In: Card{
			Name:      "Time Warp",
			Variation: "Etched 22",
			Edition:   "Strixhaven: School of Mages Etched",
		},
	},
}

func TestMain(m *testing.M) {
	allprintingsPath := os.Getenv("ALLPRINTINGS5_PATH")
	if allprintingsPath == "" {
		log.Fatalln("Need ALLPRINTINGS5_PATH variable set to run tests")
	}

	allPrintingsReader, err := os.Open(allprintingsPath)
	if err != nil {
		log.Fatalln(err)
	}
	defer allPrintingsReader.Close()

	allprints, err := mtgjson.LoadAllPrintings(allPrintingsReader)
	if err != nil {
		log.Fatalln(err)
	}

	NewDatastore(allprints)

	SetGlobalLogger(log.New(os.Stderr, "", 0))

	os.Exit(m.Run())
}

func TestMatch(t *testing.T) {
	for _, probe := range MatchTests {
		test := probe
		t.Run(test.Desc, func(t *testing.T) {
			//t.Parallel()

			cardId, err := Match(&test.In)
			if err == nil && test.Err != nil {
				t.Errorf("FAIL: Expected error: %s", test.Err.Error())
				return
			}
			if err != nil {
				if test.Err == nil {
					t.Errorf("FAIL: Unexpected error: %s", err.Error())
					return
				}
				if test.Err.Error() != err.Error() {
					t.Errorf("FAIL: Mismatched error: expected '%s', got '%s'", test.Err.Error(), err.Error())
					return
				}
			} else if cardId != test.Id {
				t.Errorf("FAIL: Id mismatch: expected '%s', got '%s'", test.Id, cardId)
				return
			}

			t.Log("PASS:", test.Desc)
		})
	}
}
