# korturl

En enkel URL-förkortare med stöd för en slumpmässigt genererad kod eller valfritt eget alias.

## Funktioner

- Förkorta URL:er med slumpmässiga koder
- Välj egna alias
- Omdirigera till ursprungliga URL:er
- Enkelt och responsivt gränssnitt

## Tekniker

- [Go](https://go.dev/)
- [Templ](https://templ.guide/)
- [HTMX](https://htmx.org/)
- [PicoCSS](https://picocss.com/)
- [PostgreSQL](https://www.postgresql.org/)

## Eventuellt planerad funktionalitet

- Frivilliga användarkonton med användarnamn som inloggning
- Möjlighet att redigera och ta bort sina länkar
- Grundläggande besöksstatistik för inloggade användare
- **Integritet:** Utan konto kopplas inga länkar till en användare och ingen besöksstatistik sparas. Konton är helt frivilliga.
