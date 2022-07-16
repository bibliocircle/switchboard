package common

import (
	"bytes"
	"html/template"
	"reflect"
	"regexp"

	"github.com/brianvoe/gofakeit/v6"
)

type fakeData struct {
	Adverb              string `fake:"{adverb}"`
	Animal              string `fake:"{animal}"`
	AnimalType          string `fake:"{animaltype}"`
	AppName             string `fake:"{appname}"`
	Boolean             string `fake:"{bool}"`
	CarFuelType         string `fake:"{carfueltype}"`
	CarMaker            string `fake:"{carmaker}"`
	CarModel            string `fake:"{carmodel}"`
	CarTransmissionType string `fake:"{cartransmissiontype}"`
	CarType             string `fake:"{cartype}"`
	CelebrityActor      string `fake:"{celebrityactor}"`
	CelebrityBusiness   string `fake:"{celebritybusiness}"`
	CelebritySport      string `fake:"{celebritysport}"`
	City                string `fake:"{city}"`
	Colour              string `fake:"{color}"`
	CompanyName         string `fake:"{company}"`
	Country             string `fake:"{country}"`
	CountryAbr          string `fake:"{countryabr}"`
	CreditCardCvv       string `fake:"{creditcardcvv}"`
	CreditCardExpiry    string `fake:"{creditcardexp}"`
	CreditCardNumber    string `fake:"{creditcardnumber}"`
	CreditCardType      string `fake:"{creditcardtype}"`
	CurrencyLong        string `fake:"{currencylong}"`
	CurrencyShort       string `fake:"{currencyshort}"`
	DomainName          string `fake:"{domainname}"`
	DomainSuffix        string `fake:"{domainsuffix}"`
	Email               string `fake:"{email}"`
	FirstName           string `fake:"{firstname}"`
	Fruit               string `fake:"{fruit}"`
	Gamertag            string `fake:"{gamertag}"`
	Gender              string `fake:"{gender}"`
	HexColour           string `fake:"{hexcolor}"`
	Hobby               string `fake:"{hobby}"`
	Hour                string `fake:"{hour}"`
	HTTPMethod          string `fake:"{httpmethod}"`
	HTTPStatusCode      string `fake:"{httpstatuscode}"`
	IPv4                string `fake:"{ipv4address}"`
	IPv6                string `fake:"{ipv6address}"`
	ISODate             string `fake:"{date}"`
	JobTitle            string `fake:"{jobtitle}"`
	Language            string `fake:"{language}"`
	LanguageCode        string `fake:"{languageabbreviation}"`
	LastName            string `fake:"{lastname}"`
	Latitude            string `fake:"{latitude}"`
	LogLevel            string `fake:"{loglevel}"`
	Longitude           string `fake:"{longitude}"`
	LoremIpsumParagraph string `fake:"{loremipsumparagraph:1,5,20}"`
	LoremIpsumSentence  string `fake:"{loremipsumsentence:5}"`
	LoremIpsumWord      string `fake:"{loremipsumword}"`
	Minute              string `fake:"{minute}"`
	Month               string `fake:"{month}"`
	MonthText           string `fake:"{monthstring}"`
	NamePrefix          string `fake:"{nameprefix}"`
	NameSuffix          string `fake:"{namesuffix}"`
	Noun                string `fake:"{noun}"`
	Number              int
	PetName             string `fake:"{petname}"`
	Phone               string `fake:"{phone}"`
	PhoneFormatted      string `fake:"{phoneformatted}"`
	Preposition         string `fake:"{preposition}"`
	ProgrammingLanguage string `fake:"{programminglanguage}"`
	RandomString        string `fake:"{randomstring}"`
	RGBColour           string `fake:"{rgbcolor}"`
	Second              string `fake:"{second}"`
	SemVer              string `fake:"{appversion}"`
	SSN                 string `fake:"{ssn}"`
	State               string `fake:"{state}"`
	StateAbr            string `fake:"{stateabr}"`
	Street              string `fake:"{street}"`
	StreetName          string `fake:"{streetname}"`
	StreetNumber        string `fake:"{streetnumber}"`
	StreetPrefix        string `fake:"{streetprefix}"`
	StreetSuffix        string `fake:"{streetsuffix}"`
	TimeZone            string `fake:"{timezone}"`
	TimeZoneFull        string `fake:"{timezonefull}"`
	URL                 string `fake:"{url}"`
	UserAgent           string `fake:"{useragent}"`
	UUID                string `fake:"{uuid}"`
	Vegetable           string `fake:"{vegetable}"`
	Verb                string `fake:"{verb}"`
	Word                string `fake:"{word}"`
	Year                string `fake:"{year}"`
	Zip                 string `fake:"{zip}"`
}

func getFakeMap() *map[string]interface{} {
	var fd fakeData
	ferr := gofakeit.Struct(&fd)

	if ferr != nil {
		return nil
	}

	v := reflect.ValueOf(fd)
	fakeMap := make(map[string]interface{}, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		fakeMap[v.Type().Field(i).Name] = v.Field(i).Interface()
	}
	return &fakeMap
}

func GenFakeData(input string) (string, error) {
	if input == "" {
		return input, nil
	}
	rg := regexp.MustCompile(`{{\s*`)
	replaced := rg.ReplaceAll([]byte(input), []byte("{{."))
	tmpl, tmplErr := template.New("").Option("missingkey=zero").Parse(string(replaced))
	if tmplErr != nil {
		return "", tmplErr
	}
	buf := new(bytes.Buffer)
	tmpl.Execute(buf, getFakeMap())
	return buf.String(), nil
}
