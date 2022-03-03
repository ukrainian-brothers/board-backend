# Aktorzy domeny

## Ogłoszenie

- Ogłoszenie ma kategorie (np. mieszkanie/transport/prawnik/praca)
- Ogłoszenie jest przypisane do użytkownika
- Ogłoszenie zawiera opis wraz z tłumaczeniami
- Ogłoszenie może zawierać lokalizacje
- Ogłoszenie zawiera tytuł wraz z tłumaczeniami
- Ogłoszenie zawiera detale
    * w przypadku mieszkania istnieje możliwośzć sprecyzowania miejsc ile osób może przyjąć
- Obrazek na liście ogłoszeń to avatar usera/placeholder

## Userzy

- Imie, naziwsko
- Numer telefonu
- Email
- Konto połączone z facebooka/XYZ
- Uprawnienia (RBAC)

# Scenariusze

## Scenariusze Użytkownika

### Rejestracja Lokalna

Użytkownik podaje login, hasło, email i numer telefonu

### Rejestracja przez Sociale (e.g Facebook)

Użytkownik klika przycisk rejestracji poprzez facebooka, po sukcesywnej weryfikacji na facebooku zostaje zwrócony na
frontend gdzie musi dodatkowe dane uzupełnić (takie których nie dostaliśmy od facebooka)

## Scenariusze ogłoszenia

### Dodawanie ogłoszenia

Dodawanie ogłoszenia wymaga zalogowania, i podania konkretnych danych (wyamganych i takich zależnych od typu ogłoszenia)
, zabezpieczone captcha

### Usuwanie ogłoszenia

Usuniecie ogłoszenia wymaga konkretnych dostępu do danego ogłoszenia

### Edytowanie ogłoszenia

Edytowanie wymaga konkretnych dostepów do danego ogłoszenia

### Pobranie ogłoszeń (z filtracja lokalizacji/tytuł/kategoria)


Można pobrać ogłoszenia filtrując filtrujac lokalizacje, tytuł, opis, kategorie

# Libki

- https://justsend.pl/
- https://www.hcaptcha.com/
- https://github.com/casbin/casbin

