package mtgmatcher

import (
	"strings"

	"github.com/kodabb/go-mtgban/mtgmatcher/mtgjson"
)

func GetUUID(uuid string) (*CardObject, error) {
	if backend.UUIDs == nil {
		return nil, ErrDatastoreEmpty
	}

	co, found := backend.UUIDs[uuid]
	if !found {
		return nil, ErrCardUnknownId
	}

	return &co, nil
}

func GetSets() map[string]*mtgjson.Set {
	return backend.Sets
}

func GetSet(code string) (*mtgjson.Set, error) {
	if backend.Sets == nil {
		return nil, ErrDatastoreEmpty
	}

	set, found := backend.Sets[strings.ToUpper(code)]
	if !found {
		return nil, ErrCardUnknownId
	}

	return set, nil
}

func GetSetByName(edition string, flags ...bool) (*mtgjson.Set, error) {
	if backend.Sets == nil {
		return nil, ErrDatastoreEmpty
	}

	card := &Card{
		Edition: edition,
	}
	if len(flags) > 0 {
		card.Foil = flags[0]
	}
	adjustEdition(card)

	for _, set := range backend.Sets {
		if set.Name == card.Edition {
			return set, nil
		}
	}

	return nil, ErrCardUnknownId
}

func GetSetUUID(uuid string) (*mtgjson.Set, error) {
	if backend.UUIDs == nil || backend.Sets == nil {
		return nil, ErrDatastoreEmpty
	}

	co, found := backend.UUIDs[uuid]
	if !found {
		return nil, ErrCardUnknownId
	}

	set, found := backend.Sets[co.SetCode]
	if !found {
		return nil, ErrCardUnknownId
	}

	return set, nil
}

func HasPromoPackPrinting(name string) bool {
	return hasPrinting(name, mtgjson.PromoTypePromoPack)
}

func HasPrereleasePrinting(name string) bool {
	return hasPrinting(name, mtgjson.PromoTypePrerelease)
}

func hasPrinting(name, promo string) bool {
	if backend.Sets == nil {
		return false
	}

	card, found := backend.Cards[Normalize(name)]
	if !found {
		cc := &Card{
			Name: name,
		}
		adjustName(cc)
		card, found = backend.Cards[Normalize(cc.Name)]
		if !found {
			return false
		}
	}
	for _, code := range card.Printings {
		set, found := backend.Sets[code]
		if !found {
			continue
		}
		for _, in := range set.Cards {
			if (card.Name == in.Name) && in.HasPromoType(promo) {
				return true
			}
		}
	}

	return false
}
