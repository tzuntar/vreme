// Ukaz vreme ustvari nekaj vremenskih postaj, jim doda meritve ter izpiše
// meritve in poročilo s statistiko.
//
// Uporaba:
//
//	vreme [--minTemp -50] [--maxTemp 50] [--stMeritev 3]
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/tzuntar/vreme"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "vreme",
		Usage: "izpis meritev in poročila vremenskih postaj",
		Flags: []cli.Flag{
			&cli.FloatFlag{
				Name:  "minTemp",
				Value: -50,
				Usage: "najnižja dovoljena temperatura meritve v °C",
			},
			&cli.FloatFlag{
				Name:  "maxTemp",
				Value: 50,
				Usage: "najvišja dovoljena temperatura meritve v °C",
			},
			&cli.IntFlag{
				Name:  "stMeritev",
				Value: 3,
				Usage: "najmanjše število meritev, potrebnih za izračun statistike",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return zazeni(cmd.Float("minTemp"), cmd.Float("maxTemp"), int(cmd.Int("stMeritev")))
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func zazeni(minTemp, maxTemp float64, stMeritev int) error {
	if minTemp > maxTemp {
		return fmt.Errorf("minTemp (%.1f) ne sme biti večja od maxTemp (%.1f)", minTemp, maxTemp)
	}

	postaje := make(map[string]vreme.Postaja)

	zacetne := []struct {
		oznaka, ime, kraj string
		meritve           []float64
	}{
		{"LJ01", "Ljubljana 1", "Ljubljana", []float64{2.0, 3.3, 2.3}},
		{"LJ02", "Ljubljana 2", "Ljubljana", []float64{11.1, 34.1, 27.12}},
		{"KR", "Kranj mesto", "Kranj", []float64{3.3, -2.3, 1.1}},
		{"KRBR", "Kranj Brnik", "Kranj", []float64{2.3, 4.1}},
	}
	for _, p := range zacetne {
		if err := vreme.DodajPostajo(postaje, p.oznaka, p.ime, p.kraj); err != nil {
			fmt.Println(err)
			continue
		}
		for _, m := range p.meritve {
			if err := vreme.DodajMeritev(postaje, p.oznaka, m, minTemp, maxTemp); err != nil {
				fmt.Println(err)
			}
		}
	}

	// Primera neveljavnih meritev.
	if err := vreme.DodajMeritev(postaje, "LJ01", 34.2, minTemp, maxTemp); err != nil {
		fmt.Println(err)
	}
	if err := vreme.DodajMeritev(postaje, "Krneki", 12.1, minTemp, maxTemp); err != nil {
		fmt.Println(err)
	}
	if err := vreme.DodajMeritev(postaje, "KR", 123.1, minTemp, maxTemp); err != nil {
		fmt.Println(err)
	}

	fmt.Println()
	vreme.IzpisMeritev(postaje)
	fmt.Println()
	vreme.IzpisPorocila(postaje, stMeritev)
	return nil
}
