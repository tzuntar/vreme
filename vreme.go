// Package vreme omogoča upravljanje vremenskih postaj in njihovih meritev
// temperature.
//
// Postaje so shranjene v slovarju, kjer je ključ oznaka postaje. Novo postajo
// dodamo s funkcijo [DodajPostajo], meritve pa s funkcijo [DodajMeritev].
// Vse meritve izpišemo s funkcijo [IzpisMeritev], poročilo s statistiko
// (najnižja, najvišja in povprečna temperatura) pa s funkcijo [IzpisPorocila].
//
// Primer uporabe:
//
//	postaje := make(map[string]vreme.Postaja)
//	vreme.DodajPostajo(postaje, "LJ01", "Ljubljana 1", "Ljubljana")
//	err := vreme.DodajMeritev(postaje, "LJ01", 21.5, -50, 50)
//	if err != nil {
//		fmt.Println(err)
//	}
//	vreme.IzpisMeritev(postaje)
//	vreme.IzpisPorocila(postaje, 3)
package vreme

import (
	"errors"
	"fmt"
	"sort"
)

// Napake, ki jih vračajo funkcije paketa.
var (
	// ErrPostajaNeObstaja pomeni, da postaja z dano oznako ni v slovarju.
	ErrPostajaNeObstaja = errors.New("postaja ne obstaja")
	// ErrPostajaObstaja pomeni, da postaja z dano oznako že obstaja.
	ErrPostajaObstaja = errors.New("postaja že obstaja")
	// ErrIzvenObmocja pomeni, da je temperatura izven dovoljenega območja.
	ErrIzvenObmocja = errors.New("temperatura izven dovoljenega območja")
)

// Postaja predstavlja vremensko postajo z imenom, krajem in seznamom
// izmerjenih temperatur v °C. Polja so neizvožena, zato postajo ustvarimo
// s funkcijo [DodajPostajo].
type Postaja struct {
	ime     string
	kraj    string
	meritve []float64
}

// DodajPostajo v slovar postaje doda novo postajo brez meritev z dano
// oznako, imenom in krajem. Če postaja z dano oznako že obstaja, vrne
// [ErrPostajaObstaja] in slovarja ne spremeni.
func DodajPostajo(postaje map[string]Postaja, oznaka, ime, kraj string) error {
	if _, ok := postaje[oznaka]; ok {
		return fmt.Errorf("%w: %s", ErrPostajaObstaja, oznaka)
	}
	postaje[oznaka] = Postaja{ime: ime, kraj: kraj}
	return nil
}

// DodajMeritev postaji z dano oznako doda meritev temperature. Meritev mora
// biti v zaprtem intervalu [minTemp, maxTemp], sicer funkcija vrne
// [ErrIzvenObmocja]. Če postaja ne obstaja, vrne [ErrPostajaNeObstaja].
func DodajMeritev(postaje map[string]Postaja, oznaka string, temperatura, minTemp, maxTemp float64) error {
	if temperatura < minTemp || temperatura > maxTemp {
		return fmt.Errorf("%w: %.1f °C ni v [%.1f, %.1f]", ErrIzvenObmocja, temperatura, minTemp, maxTemp)
	}
	postaja, ok := postaje[oznaka]
	if !ok {
		return fmt.Errorf("%w: %s", ErrPostajaNeObstaja, oznaka)
	}
	postaja.meritve = append(postaja.meritve, temperatura)
	postaje[oznaka] = postaja
	return nil
}

// statistika vrne najnižjo, najvišjo in povprečno temperaturo postaje.
// Vrednost ok je false, če postaja ne obstaja ali ima manj kot stMeritev
// meritev.
func statistika(postaje map[string]Postaja, oznaka string, stMeritev int) (najnizja, najvisja, povprecna float64, ok bool) {
	postaja, obstaja := postaje[oznaka]
	if !obstaja || len(postaja.meritve) == 0 || len(postaja.meritve) < stMeritev {
		return 0, 0, 0, false
	}

	najnizja = postaja.meritve[0]
	najvisja = postaja.meritve[0]
	vsota := 0.0
	for _, meritev := range postaja.meritve {
		vsota += meritev
		najnizja = min(najnizja, meritev)
		najvisja = max(najvisja, meritev)
	}
	povprecna = vsota / float64(len(postaja.meritve))
	return najnizja, najvisja, povprecna, true
}

// IzpisMeritev na standardni izhod izpiše vse postaje (urejene po oznaki)
// skupaj z njihovimi meritvami.
func IzpisMeritev(postaje map[string]Postaja) {
	fmt.Println("MERITVE:")
	for _, oznaka := range urejeneOznake(postaje) {
		postaja := postaje[oznaka]
		fmt.Printf("%s - %s (%s): [", oznaka, postaja.ime, postaja.kraj)
		for idx, meritev := range postaja.meritve {
			if idx > 0 {
				fmt.Print(" ")
			}
			fmt.Printf("%.1f", meritev)
		}
		fmt.Println("]")
	}
}

// IzpisPorocila na standardni izhod izpiše poročilo za vse postaje (urejene
// po oznaki): najnižjo, najvišjo in povprečno temperaturo ter oceno
// (Vroče, Toplo, Hladno ali Mrzlo) glede na povprečje. Za postaje z manj kot
// stMeritev meritvami izpiše opozorilo.
func IzpisPorocila(postaje map[string]Postaja, stMeritev int) {
	fmt.Println("POROČILO:")
	for _, oznaka := range urejeneOznake(postaje) {
		postaja := postaje[oznaka]
		najnizja, najvisja, povprecna, ok := statistika(postaje, oznaka, stMeritev)
		if !ok {
			fmt.Printf("%s: premalo meritev (%d, potrebnih vsaj %d)\n", postaja.ime, len(postaja.meritve), stMeritev)
			continue
		}
		fmt.Printf("%s: min %.1f °C, max %.1f °C, povprečje %.1f °C -> %s\n",
			postaja.ime, najnizja, najvisja, povprecna, ocena(povprecna))
	}
}

// ocena vrne opisno oceno povprečne temperature.
func ocena(povprecna float64) string {
	switch {
	case povprecna >= 25:
		return "Vroče"
	case povprecna >= 15:
		return "Toplo"
	case povprecna >= 5:
		return "Hladno"
	default:
		return "Mrzlo"
	}
}

// urejeneOznake vrne oznake postaj v abecednem vrstnem redu, da je izpis
// ponovljiv.
func urejeneOznake(postaje map[string]Postaja) []string {
	oznake := make([]string, 0, len(postaje))
	for oznaka := range postaje {
		oznake = append(oznake, oznaka)
	}
	sort.Strings(oznake)
	return oznake
}
