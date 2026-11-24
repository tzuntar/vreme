# vreme

Modul v Go za delo z vremenskimi postajami (VPSA, domača naloga 5).

- paket `vreme` (koren modula): struktura `Postaja` ter funkcije `DodajPostajo`, `DodajMeritev`, `IzpisMeritev`, `IzpisPorocila`
- `cmd/vreme`: aplikacija z ukaznim vmesnikom ([urfave/cli/v3](https://github.com/urfave/cli))

## Uporaba

```bash
go run ./cmd/vreme --minTemp -50 --maxTemp 50 --stMeritev 3
go run ./cmd/vreme --help
```

| Stikalo       | Opis                                                     | Privzeto |
|---------------|----------------------------------------------------------|----------|
| `--minTemp`   | najnižja dovoljena temperatura meritve (°C)              | -50      |
| `--maxTemp`   | najvišja dovoljena temperatura meritve (°C)              | 50       |
| `--stMeritev` | najmanjše število meritev, potrebnih za izračun statistike | 3      |

Uvoz paketa v drug modul:

```bash
go get github.com/tzuntar/vreme@v1.0.0
```

Dokumentacija: https://pkg.go.dev/github.com/tzuntar/vreme
