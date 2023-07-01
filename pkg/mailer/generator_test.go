package mailer_test

import (
	"github.com/aasumitro/tix/pkg/mailer"
	"github.com/aasumitro/tix/pkg/mailer/template"
	"github.com/stretchr/testify/assert"
	"testing"
)

var testedThemes = []mailer.Theme{
	// Insert your new theme here
	new(template.Default),
	new(template.Flat),
}

/////////////////////////////////////////////////////
// Every theme should display the same information //
// Find below the tests to check that              //
/////////////////////////////////////////////////////

// Implement this interface when creating a new example checking a common feature of all themes
type Example interface {
	// Create the mailer example with data
	// Represents the "Given" step in Given/When/Then Workflow
	getExample() (m mailer.Mailer, email mailer.Email)
	// Checks the content of the generated HTML email by asserting content presence or not
	assertHTMLContent(t *testing.T, s string)
	// Checks the content of the generated Plaintext email by asserting content presence or not
	assertPlainTextContent(t *testing.T, s string)
}

// Scenario
type SimpleExample struct {
	theme mailer.Theme
}

func (ed *SimpleExample) getExample() (mailer.Mailer, mailer.Email) {
	h := mailer.Mailer{
		Theme: ed.theme,
		Product: mailer.Product{
			Name:      "MailerName",
			Link:      "http://mailer-link.com",
			Copyright: "Copyright © Mailer-Test",
			Logo:      "http://www.duchess-france.org/wp-content/uploads/2016/01/gopher.png",
		},
		TextDirection:      mailer.TDLeftToRight,
		DisableCSSInlining: true,
	}

	email := mailer.Email{
		Body: mailer.Body{
			Name: "Jon Snow",
			Intros: []string{
				"Welcome to Mailer! We're very excited to have you on board.",
			},
			Dictionary: []mailer.Entry{
				{"Firstname", "Jon"},
				{"Lastname", "Snow"},
				{"Birthday", "01/01/283"},
			},
			Table: mailer.Table{
				Data: [][]mailer.Entry{
					{
						{Key: "Item", Value: "Golang"},
						{Key: "Description", Value: "Open source programming language that makes it easy to build simple, reliable, and efficient software"},
						{Key: "Price", Value: "$10.99"},
					},
					{
						{Key: "Item", Value: "Mailer"},
						{Key: "Description", Value: "Programmatically create beautiful e-mails using Golang."},
						{Key: "Price", Value: "$1.99"},
					},
				},
				Columns: mailer.Columns{
					CustomWidth: map[string]string{
						"Item":  "20%",
						"Price": "15%",
					},
					CustomAlignment: map[string]string{
						"Price": "right",
					},
				},
			},
			Actions: []mailer.Action{
				{
					Instructions: "To get started with  Mailer, please click here:",
					Button: mailer.Button{
						Color: "#22BC66",
						Text:  "Confirm your account",
						Link:  "https://mailer-example.com/confirm?token=d9729feb74992cc3482b350163a1a010",
					},
				},
			},
			Outros: []string{
				"Need help, or have questions? Just reply to this email, we'd love to help.",
			},
		},
	}
	return h, email
}

func (ed *SimpleExample) assertHTMLContent(t *testing.T, r string) {
	// Assert on product
	assert.Contains(t, r, "MailerName", "Product: Should find the name of the product in email")
	assert.Contains(t, r, "http://mailer-link.com", "Product: Should find the link of the product in email")
	assert.Contains(t, r, "Copyright © Mailer-Test", "Product: Should find the Copyright of the product in email")
	assert.Contains(t, r, "http://www.duchess-france.org/wp-content/uploads/2016/01/gopher.png", "Product: Should find the logo of the product in email")
	assert.Contains(t, r, "If you’re having trouble with the button &#39;Confirm your account&#39;, copy and paste the URL below into your web browser.", "Product: Should find the trouble text in email")
	// Assert on email body
	assert.Contains(t, r, "Hi Jon Snow", "Name: Should find the name of the person")
	assert.Contains(t, r, "Welcome to Mailer", "Intro: Should have intro")
	assert.Contains(t, r, "Birthday", "Dictionary: Should have dictionary")
	assert.Contains(t, r, "Open source programming language", "Table: Should have table with first row and first column")
	assert.Contains(t, r, "Programmatically create beautiful e-mails using Golang", "Table: Should have table with second row and first column")
	assert.Contains(t, r, "$10.99", "Table: Should have table with first row and second column")
	assert.Contains(t, r, "$1.99", "Table: Should have table with second row and second column")
	assert.Contains(t, r, "Confirm your account", "Action: Should have button of action")
	assert.Contains(t, r, "#22BC66", "Action: Button should have given color")
	assert.Contains(t, r, "https://mailer-example.com/confirm?token=d9729feb74992cc3482b350163a1a010", "Action: Button should have link")
	assert.Contains(t, r, "Need help, or have questions", "Outro: Should have outro")
}

func (ed *SimpleExample) assertPlainTextContent(t *testing.T, r string) {
	// Assert on product
	assert.Contains(t, r, "MailerName", "Product: Should find the name of the product in email")
	assert.Contains(t, r, "http://mailer-link.com", "Product: Should find the link of the product in email")
	assert.Contains(t, r, "Copyright © Mailer-Test", "Product: Should find the Copyright of the product in email")
	assert.NotContains(t, r, "http://www.duchess-france.org/wp-content/uploads/2016/01/gopher.png", "Product: Should not find any logo in plain text")

	// Assert on email body
	assert.Contains(t, r, "Hi Jon Snow", "Name: Should find the name of the person")
	assert.Contains(t, r, "Welcome to Mailer", "Intro: Should have intro")
	assert.Contains(t, r, "Birthday", "Dictionary: Should have dictionary")
	assert.Contains(t, r, "Open source", "Table: Should have table content")
	assert.Contains(t, r, `+--------+--------------------------------+--------+
|  ITEM  |          DESCRIPTION           | PRICE  |
+--------+--------------------------------+--------+
| Golang | Open source programming        | $10.99 |
|        | language that makes it easy    |        |
|        | to build simple, reliable, and |        |
|        | efficient software             |        |
| Mailer | Programmatically create        | $1.99  |
|        | beautiful e-mails using        |        |
|        | Golang.                        |        |
+--------+--------------------------------+--------`, "Table: Should have pretty table content")
	assert.Contains(t, r, "started with Mailer", "Action: Should have instruction")
	assert.NotContains(t, r, "Confirm your account", "Action: Should not have button of action in plain text")
	assert.NotContains(t, r, "#22BC66", "Action: Button should not have color in plain text")
	assert.Contains(t, r, "https://mailer-example.com/confirm?token=d9729feb74992cc3482b350163a1a010", "Action: Even if button is not possible in plain text, it should have the link")
	assert.Contains(t, r, "Need help, or have questions", "Outro: Should have outro")
}

type WithTitleInsteadOfNameExample struct {
	theme mailer.Theme
}

func (ed *WithTitleInsteadOfNameExample) getExample() (mailer.Mailer, mailer.Email) {
	h := mailer.Mailer{
		Theme: ed.theme,
		Product: mailer.Product{
			Name: "Mailer",
			Link: "http://mailer.com",
		},
		DisableCSSInlining: true,
	}

	email := mailer.Email{
		Body: mailer.Body{
			Name:  "Jon Snow",
			Title: "A new e-mail",
		},
	}
	return h, email
}

func (ed *WithTitleInsteadOfNameExample) assertHTMLContent(t *testing.T, r string) {
	assert.NotContains(t, r, "Hi Jon Snow", "Name: should not find greetings from Jon Snow because title should be used")
	assert.Contains(t, r, "A new e-mail", "Title should be used instead of name")
}

func (ed *WithTitleInsteadOfNameExample) assertPlainTextContent(t *testing.T, r string) {
	assert.NotContains(t, r, "Hi Jon Snow", "Name: should not find greetings from Jon Snow because title should be used")
	assert.Contains(t, r, "A new e-mail", "Title shoud be used instead of name")
}

type WithGreetingDifferentThanDefault struct {
	theme mailer.Theme
}

func (ed *WithGreetingDifferentThanDefault) getExample() (mailer.Mailer, mailer.Email) {
	h := mailer.Mailer{
		Theme: ed.theme,
		Product: mailer.Product{
			Name: "Mailer",
			Link: "http://mailer.com",
		},
		DisableCSSInlining: true,
	}

	email := mailer.Email{
		Body: mailer.Body{
			Greeting: "Dear",
			Name:     "Jon Snow",
		},
	}
	return h, email
}

func (ed *WithGreetingDifferentThanDefault) assertHTMLContent(t *testing.T, r string) {
	assert.NotContains(t, r, "Hi Jon Snow", "Should not find greetings with 'Hi' which is default")
	assert.Contains(t, r, "Dear Jon Snow", "Should have greeting with Dear")
}

func (ed *WithGreetingDifferentThanDefault) assertPlainTextContent(t *testing.T, r string) {
	assert.NotContains(t, r, "Hi Jon Snow", "Should not find greetings with 'Hi' which is default")
	assert.Contains(t, r, "Dear Jon Snow", "Should have greeting with Dear")
}

type WithSignatureDifferentThanDefault struct {
	theme mailer.Theme
}

func (ed *WithSignatureDifferentThanDefault) getExample() (mailer.Mailer, mailer.Email) {
	h := mailer.Mailer{
		Theme: ed.theme,
		Product: mailer.Product{
			Name: "Mailer",
			Link: "http://mailer.com",
		},
		DisableCSSInlining: true,
	}

	email := mailer.Email{
		Body: mailer.Body{
			Name:      "Jon Snow",
			Signature: "Best regards",
		},
	}
	return h, email
}

func (ed *WithSignatureDifferentThanDefault) assertHTMLContent(t *testing.T, r string) {
	assert.NotContains(t, r, "Yours truly", "Should not find signature with 'Yours truly' which is default")
	assert.Contains(t, r, "Best regards", "Should have greeting with Dear")
}

func (ed *WithSignatureDifferentThanDefault) assertPlainTextContent(t *testing.T, r string) {
	assert.NotContains(t, r, "Yours truly", "Should not find signature with 'Yours truly' which is default")
	assert.Contains(t, r, "Best regards", "Should have greeting with Dear")
}

type WithInviteCode struct {
	theme mailer.Theme
}

func (ed *WithInviteCode) getExample() (mailer.Mailer, mailer.Email) {
	h := mailer.Mailer{
		Theme: ed.theme,
		Product: mailer.Product{
			Name: "Mailer",
			Link: "http://mailer.com",
		},
		DisableCSSInlining: true,
	}

	email := mailer.Email{
		Body: mailer.Body{
			Name: "Jon Snow",
			Actions: []mailer.Action{
				{
					Instructions: "Here is your invite code:",
					InviteCode:   "123456",
				},
			},
		},
	}
	return h, email
}

func (ed *WithInviteCode) assertHTMLContent(t *testing.T, r string) {
	assert.Contains(t, r, "Here is your invite code", "Should contains the instruction")
	assert.Contains(t, r, "123456", "Should contain the short code")
}

func (ed *WithInviteCode) assertPlainTextContent(t *testing.T, r string) {
	assert.Contains(t, r, "Here is your invite code", "Should contains the instruction")
	assert.Contains(t, r, "123456", "Should contain the short code")
}

type WithFreeMarkdownContent struct {
	theme mailer.Theme
}

func (ed *WithFreeMarkdownContent) getExample() (mailer.Mailer, mailer.Email) {
	h := mailer.Mailer{
		Theme: ed.theme,
		Product: mailer.Product{
			Name: "Mailer",
			Link: "http://mailer.com",
		},
		DisableCSSInlining: true,
	}

	email := mailer.Email{
		Body: mailer.Body{
			Name: "Jon Snow",
			FreeMarkdown: `
> _Mailer_ service will shutdown the **1st August 2017** for maintenance operations. 

Services will be unavailable based on the following schedule:

| Services | Downtime |
| :------:| :-----------: |
| Service A | 2AM to 3AM |
| Service B | 4AM to 5AM |
| Service C | 5AM to 6AM |

---

Feel free to contact us for any question regarding this matter at [support@mailer-example.com](mailto:support@mailer-example.com) or in our [Gitter](https://gitter.im/)

`,
			Intros: []string{
				"An intro that should be kept even with FreeMarkdown",
			},
			Dictionary: []mailer.Entry{
				{"Dictionary that should not be displayed", "Because of FreeMarkdown"},
			},
			Table: mailer.Table{
				Data: [][]mailer.Entry{
					{
						{Key: "Item", Value: "Golang"},
					},
					{
						{Key: "Item", Value: "Mailer"},
					},
				},
			},
			Actions: []mailer.Action{
				{
					Instructions: "Action that should not be displayed, because of FreeMarkdown:",
					Button: mailer.Button{
						Color: "#22BC66",
						Text:  "Button",
						Link:  "https://mailer-example.com/confirm?token=d9729feb74992cc3482b350163a1a010",
					},
				},
			},
			Outros: []string{
				"An outro that should be kept even with FreeMarkdown",
			},
		},
	}
	return h, email
}

func (ed *WithFreeMarkdownContent) assertHTMLContent(t *testing.T, r string) {
	assert.Contains(t, r, "Yours truly", "Should find signature with 'Yours truly' which is default")
	assert.Contains(t, r, "Jon Snow", "Should find title with 'Jon Snow'")
	assert.Contains(t, r, "<em>Mailer</em> service will shutdown", "Should find quote as HTML formatted content")
	assert.Contains(t, r, "<td align=\"center\">2AM to 3AM</td>", "Should find cell content as HTML formatted content")
	assert.Contains(t, r, "<a href=\"mailto:support@mailer-example.com\">support@mailer-example.com</a>", "Should find link of mailto as HTML formatted content")
	assert.Contains(t, r, "An intro that should be kept even with FreeMarkdown", "Should find intro even with FreeMarkdown")
	assert.Contains(t, r, "An outro that should be kept even with FreeMarkdown", "Should find outro even with FreeMarkdown")
	assert.NotContains(t, r, "should not be displayed", "Should find any other content that the one from FreeMarkdown object")
}

func (ed *WithFreeMarkdownContent) assertPlainTextContent(t *testing.T, r string) {
	assert.Contains(t, r, "Yours truly", "Should find signature with 'Yours truly' which is default")
	assert.Contains(t, r, "Jon Snow", "Should find title with 'Jon Snow'")
	assert.Contains(t, r, "> Mailer service will shutdown", "Should find quote as plain text with quote emphaze on sentence")
	assert.Contains(t, r, "2AM to 3AM", "Should find cell content as plain text")
	assert.Contains(t, r, `+-----------+------------+
| SERVICES  |  DOWNTIME  |
+-----------+------------+
| Service A | 2AM to 3AM |
| Service B | 4AM to 5AM |
| Service C | 5AM to 6AM |
+-----------+------------+`, "Should find pretty table as plain text")
	assert.Contains(t, r, "support@mailer-example.com", "Should find link of mailto as plain text")
	assert.NotContains(t, r, "<table>", "Should not find html table tags")
	assert.NotContains(t, r, "<tr>", "Should not find html tr tags")
	assert.NotContains(t, r, "<a>", "Should not find html link tags")
	assert.NotContains(t, r, "should not be displayed", "Should find any other content that the one from FreeMarkdown object")
}

// Test all the themes for the features

func TestThemeSimple(t *testing.T) {
	for _, theme := range testedThemes {
		checkExample(t, &SimpleExample{theme})
	}
}

func TestThemeWithTitleInsteadOfName(t *testing.T) {
	for _, theme := range testedThemes {
		checkExample(t, &WithTitleInsteadOfNameExample{theme})
	}
}

func TestThemeWithGreetingDifferentThanDefault(t *testing.T) {
	for _, theme := range testedThemes {
		checkExample(t, &WithGreetingDifferentThanDefault{theme})
	}
}

func TestThemeWithGreetingDiffrentThanDefault(t *testing.T) {
	for _, theme := range testedThemes {
		checkExample(t, &WithSignatureDifferentThanDefault{theme})
	}
}

func TestThemeWithFreeMarkdownContent(t *testing.T) {
	for _, theme := range testedThemes {
		checkExample(t, &WithFreeMarkdownContent{theme})
	}
}

func TestThemeWithInviteCode(t *testing.T) {
	for _, theme := range testedThemes {
		checkExample(t, &WithInviteCode{theme})
	}
}

func checkExample(t *testing.T, ex Example) {
	// Given an example
	h, email := ex.getExample()

	// When generating HTML template
	r, err := h.GenerateHTML(&email)
	t.Log(r)
	assert.Nil(t, err)
	assert.NotEmpty(t, r)

	// Then asserting HTML is OK
	ex.assertHTMLContent(t, r)

	// When generating plain text template
	r, err = h.GeneratePlainText(&email)
	t.Log(r)
	assert.Nil(t, err)
	assert.NotEmpty(t, r)

	// Then asserting plain text is OK
	ex.assertPlainTextContent(t, r)
}

////////////////////////////////////////////
// Tests on default values for all themes //
// It does not check email content        //
////////////////////////////////////////////

func TestMailer_TextDirectionAsDefault(t *testing.T) {
	h := mailer.Mailer{
		Product: mailer.Product{
			Name: "Mailer",
			Link: "http://mailer.com",
		},
		TextDirection:      "not-existing", // Wrong text-direction
		DisableCSSInlining: true,
	}

	email := mailer.Email{
		Body: mailer.Body{
			Name: "Jon Snow",
			Intros: []string{
				"Welcome to Mailer! We're very excited to have you on board.",
			},
			Actions: []mailer.Action{
				{
					Instructions: "To get started with Mailer, please click here:",
					Button: mailer.Button{
						Color: "#22BC66",
						Text:  "Confirm your account",
						Link:  "https://mailer-example.com/confirm?token=d9729feb74992cc3482b350163a1a010",
					},
				},
			},
			Outros: []string{
				"Need help, or have questions? Just reply to this email, we'd love to help.",
			},
		},
	}

	_, err := h.GenerateHTML(&email)
	assert.Nil(t, err)
	assert.Equal(t, h.TextDirection, mailer.TDLeftToRight)
	assert.Equal(t, h.Theme.Name(), "default")
}

func TestMailer_Default(t *testing.T) {
	h := mailer.Mailer{}
	_ = mailer.SetDefaultMailerValues(&h)
	email := mailer.Email{}
	_ = mailer.SetDefaultEmailValues(&email)

	assert.Equal(t, h.TextDirection, mailer.TDLeftToRight)
	assert.Equal(t, h.Theme, new(template.Default))
	assert.Equal(t, h.Product.Name, "Mailer")
	assert.Equal(t, h.Product.Copyright, "Copyright © 2020 BAKODE. All rights reserved.")

	assert.Empty(t, email.Body.Actions)
	assert.Empty(t, email.Body.Dictionary)
	assert.Empty(t, email.Body.Intros)
	assert.Empty(t, email.Body.Outros)
	assert.Empty(t, email.Body.Table.Data)
	assert.Empty(t, email.Body.Table.Columns.CustomWidth)
	assert.Empty(t, email.Body.Table.Columns.CustomAlignment)
	assert.Empty(t, string(email.Body.FreeMarkdown))

	assert.Equal(t, email.Body.Greeting, "Hi")
	assert.Equal(t, email.Body.Signature, "Yours truly")
	assert.Empty(t, email.Body.Title)
}
