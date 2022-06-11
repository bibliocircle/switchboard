Switchboard provides a templating engine which you can use to generate random data in mock service responses by embedding placeholders in the template using `{{placeholder}}` notation.

The following is an example to demonstrate how to use templating.

**Input Template Example**

```json
{
    "NamePrefix": "{{NamePrefix}}",
    "NameSuffix": "{{NameSuffix}}",
    "FirstName": "{{FirstName}}",
    "LastName": "{{LastName}}",
    "Gender": "{{Gender}}",
    "SSN": "{{SSN}}",
    "Hobby": "{{Hobby}}",
}
```

**Output**

```json
{
    "NamePrefix": "Miss",
    "NameSuffix": "IV",
    "FirstName": "Liliana",
    "LastName": "Jerde",
    "Gender": "male",
    "SSN": "265866522",
    "Hobby": "Iceboat racing",
}
```

## All Available Random Generator Placeholders

- `Adverb`
- `Animal`
- `AnimalType`
- `AppName`
- `Boolean`
- `CarFuelType`
- `CarMaker`
- `CarModel`
- `CarTransmissionType`
- `CarType`
- `CelebrityActor`
- `CelebrityBusiness`
- `CelebritySport`
- `City`
- `Colour`
- `CompanyName`
- `Country`
- `CountryAbr`
- `CreditCardCvv`
- `CreditCardExpiry`
- `CreditCardNumber`
- `CreditCardType`
- `CurrencyLong`
- `CurrencyShort`
- `DomainName`
- `DomainSuffix`
- `Email`
- `FirstName`
- `Fruit`
- `Gamertag`
- `Gender`
- `HexColour`
- `Hobby`
- `Hour`
- `HTTPMethod`
- `HTTPStatusCode`
- `IPv4`
- `IPv6`
- `ISODate`
- `JobTitle`
- `Language`
- `LanguageCode`
- `LastName`
- `Latitude`
- `LogLevel`
- `Longitude`
- `LoremIpsumParagraph`
- `LoremIpsumSentence`
- `LoremIpsumWord`
- `Minute`
- `Month`
- `MonthText`
- `NamePrefix`
- `NameSuffix`
- `Noun`
- `Number`
- `PetName`
- `Phone`
- `PhoneFormatted`
- `Preposition`
- `ProgrammingLanguage`
- `RandomString`
- `RGBColour`
- `Second`
- `SemVer`
- `SSN`
- `State`
- `StateAbr`
- `Street`
- `StreetName`
- `StreetNumber`
- `StreetPrefix`
- `StreetSuffix`
- `TimeZone`
- `TimeZoneFull`
- `URL`
- `UserAgent`
- `UUID`
- `Vegetable`
- `Verb`
- `Word`
- `Year`
- `Zip`

