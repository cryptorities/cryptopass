package cmd

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/dialog"
	mapp "github.com/cryptorities/cryptopass/pkg/app"
	"github.com/cryptorities/cryptopass/pkg/crypto"
	"path/filepath"
	"strings"
)

/**
	Alex Shvid
*/

type uiCommand struct {
}

func (t *uiCommand) Desc() string {
	return "user interface (ui)"
}

func (t *uiCommand) Usage() string {
	return "ui"
}

func (t *uiCommand) Run(args []string) error {
	a := app.New()
	win := a.NewWindow("CryptoPass")

	description := "Place public key on your server.\nUse private key on your desktop only for token generation.\nGenerated tokens use as passwords for rigs.\nCommand line interface usage: './cryptopass help'."
	hello := widget.NewLabel("Welcome to token generation program!")
	hello.Alignment = fyne.TextAlignCenter
	output := widget.NewMultiLineEntry()
	output.SetMinRowsVisible(5)
	output.SetText(description)
	win.SetContent(container.NewVBox(
		hello,
		widget.NewButton("About", func() {
			appInfo := mapp.GetAppInfo()
			message := fmt.Sprintf("%s [Version %s, Build %s]\n%s\n", mapp.ApplicationName, appInfo.Version, appInfo.Build, mapp.Copyright)
			dialog.ShowInformation(mapp.ApplicationName, message, win)
		}),
		widget.NewButton("Generate", func() {
			output.SetText("")
			publicKey, privateKey, err := generateKeyPair()
			if err != nil {
				output.SetText(err.Error())
			} else {
				output.SetText(fmt.Sprintf("PublicKey:\n%s\nPrivateKey:\n%s\n", publicKey, privateKey))
			}
		}),
		widget.NewButton("Issue Token", func() {
			output.SetText("")
			username := widget.NewEntry()
			username.Validator = validation.NewRegexp(`^[A-Za-z0-9_-]+$`, "username can only contain letters, numbers, '_', and '-'")
			expirationDate := widget.NewEntry()
			expirationDate.Validator = validation.NewRegexp(`^\d{4}\-(0?[1-9]|1[012])\-(0?[1-9]|[12][0-9]|3[01])$`, "expirationDate can only be in format YYYY-mm-dd")
			privateKey := widget.NewPasswordEntry()
			privateKey.Validator = validation.NewRegexp(`^[A-Za-z0-9_-]+$`, "private key must be Base64 url encoding")
			items := []*widget.FormItem{
				widget.NewFormItem("Username", username),
				widget.NewFormItem("Expiration Date YYYY-MM-DD", expirationDate),
				widget.NewFormItem("Private Key", privateKey),
			}
			form := dialog.NewForm("Issue Token...", "Issue", "Cancel", items, func(b bool) {
				if !b {
					return
				}
				token, err := crypto.Issue(strings.TrimSpace(username.Text), strings.TrimSpace(expirationDate.Text), func() (s string, err error) {
					return strings.TrimSpace(privateKey.Text), nil
				})
				if err != nil {
					output.SetText(err.Error())
				} else {
					output.SetText(fmt.Sprintf("Token:\n%s\n", token))
				}
			}, win)
			form.Resize(fyne.NewSize(440, 260))
			form.Show()
		}),
		widget.NewButton("Revoke Token", func() {
			output.SetText("")
			username := widget.NewEntry()
			username.Validator = validation.NewRegexp(`^[A-Za-z0-9_-]+$`, "username can only contain letters, numbers, '_', and '-'")
			expirationDate := widget.NewEntry()
			expirationDate.Validator = validation.NewRegexp(`^\d{4}\-(0?[1-9]|1[012])\-(0?[1-9]|[12][0-9]|3[01])$`, "expirationDate can only be in format YYYY-mm-dd")
			privateKey := widget.NewPasswordEntry()
			privateKey.Validator = validation.NewRegexp(`^[A-Za-z0-9_-]+$`, "private key must be Base64 url encoding")
			items := []*widget.FormItem{
				widget.NewFormItem("Username", username),
				widget.NewFormItem("Expiration Date YYYY-MM-DD", expirationDate),
				widget.NewFormItem("Private Key", privateKey),
			}
			form := dialog.NewForm("Revoke Token...", "Revoke", "Cancel", items, func(b bool) {
				if !b {
					return
				}
				token, err := crypto.Revoke(strings.TrimSpace(username.Text), strings.TrimSpace(expirationDate.Text), func() (s string, err error) {
					return strings.TrimSpace(privateKey.Text), nil
				})
				if err != nil {
					output.SetText(err.Error())
				} else {
					output.SetText(fmt.Sprintf("Revoke Token:\n%s\n", token))
				}
			}, win)
			form.Resize(fyne.NewSize(440, 260))
			form.Show()
		}),
		widget.NewButton("Verify Token", func() {
			output.SetText("")
			username := widget.NewEntry()
			username.Validator = validation.NewRegexp(`^[A-Za-z0-9_-]+$`, "username can only contain letters, numbers, '_', and '-'")
			token := widget.NewEntry()
			token.Validator = validation.NewRegexp(`^[A-Za-z0-9_-]+\.[0-9]+$`, "token base64")
			publicKey := widget.NewEntry()
			publicKey.Validator = validation.NewRegexp(`^[A-Za-z0-9_-]+$`, "public key must be Base64 url encoding")
			items := []*widget.FormItem{
				widget.NewFormItem("Username", username),
				widget.NewFormItem("Token", token),
				widget.NewFormItem("Public Key", publicKey),
			}
			form := dialog.NewForm("Verify Token...", "Verify", "Cancel", items, func(b bool) {
				if !b {
					return
				}
				var out strings.Builder
				valid, expiration, err := crypto.VerifyIssued(strings.TrimSpace(username.Text), strings.TrimSpace(token.Text), func() (s string, err error) {
					return strings.TrimSpace(publicKey.Text), nil
				})
				if err != nil {
					out.WriteString(fmt.Sprintf("Verify Issued Error: %s\n", err.Error()))
				} else if valid {
					out.WriteString(fmt.Sprintf("Valid Till %s\n", expiration))
				} else {
					out.WriteString("Invalid\n")
				}

				valid, expiration, err = crypto.VerifyRevoked(strings.TrimSpace(username.Text), strings.TrimSpace(token.Text), func() (s string, err error) {
					return strings.TrimSpace(publicKey.Text), nil
				})
				if err != nil {
					out.WriteString(fmt.Sprintf("Verify Revoked Error: %s\n", err.Error()))
				} else if valid {
					out.WriteString(fmt.Sprintf("Revoked On %s\n", expiration))
				}

				output.SetText(out.String())

			}, win)
			form.Resize(fyne.NewSize(440, 260))
			form.Show()
		}),
		widget.NewButton("Encrypt File", func() {
			output.SetText("")
			fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
				if err != nil {
					output.SetText(err.Error())
					return
				}
				if reader == nil {
					output.SetText("user cancelled")
					return
				}

				inputFile := reader.URI().String()
				if len(inputFile) == 0 {
					output.SetText("user did not select input file")
					return
				}
				if strings.HasPrefix(inputFile, "file://") {
					inputFile = inputFile[7:]
				}

				output.SetText(fmt.Sprintf("Input File: %s\n", inputFile))
				reader.Close()

				outputFile := fmt.Sprintf("%s.cp", inputFile)
				fs := dialog.NewFileSave(func(closer fyne.URIWriteCloser, err error) {
					if err != nil {
						output.SetText(err.Error())
						return
					}
					if closer == nil {
						output.SetText("user cancelled")
						return
					}
					outputFile = closer.URI().String()
					if len(outputFile) == 0 {
						output.SetText("user did not select output file")
						return
					}
					if strings.HasPrefix(outputFile, "file://") {
						outputFile = outputFile[7:]
					}

					output.SetText(fmt.Sprintf("Output File: %s\n", outputFile))
					closer.Close()

					promptDialog("Enter...", "Recipient Public Key",`^[A-Za-z0-9_-]+$`, false, win, func(publicKey string) {
						output.SetText(fmt.Sprintf("Using Recipient Public Key: %s\n", publicKey))
						n, err := crypto.EncryptFile(inputFile, outputFile, func() (s string, err error) {
							return publicKey, nil
						})
						if err != nil {
							output.SetText(err.Error())
						} else {
							output.SetText(fmt.Sprintf("Input File '%s'\nWas encrypted to\nOutput File '%s'\nWritten %d bytes.", inputFile, outputFile, n))
						}
					})

				}, win)
				_, fileName := filepath.Split(outputFile)
				fs.SetFileName(fileName)
				fs.Show()

			}, win)
			fd.Show()
		}),
		widget.NewButton("Decrypt File", func() {
			output.SetText("")
			fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
				if err != nil {
					dialog.ShowError(err, win)
					return
				}
				if reader == nil {
					output.SetText("user cancelled")
					return
				}

				inputFile := reader.URI().String()
				if len(inputFile) == 0 {
					output.SetText("user did not select file")
					return
				}
				if strings.HasPrefix(inputFile, "file://") {
					inputFile = inputFile[7:]
				}

				output.SetText(fmt.Sprintf("Selected Input File: %s\n", inputFile))
				reader.Close()

				var outputFile string
				if strings.HasSuffix(inputFile, ".cp") {
					outputFile = inputFile[:len(inputFile)-3]
				} else {
					outputFile = inputFile
				}

				fs := dialog.NewFileSave(func(closer fyne.URIWriteCloser, err error) {
					if err != nil {
						output.SetText(err.Error())
						return
					}
					if closer == nil {
						output.SetText("user cancelled")
						return
					}
					outputFile = closer.URI().String()
					if len(outputFile) == 0 {
						output.SetText("user did not select output file")
						return
					}
					if strings.HasPrefix(outputFile, "file://") {
						outputFile = outputFile[7:]
					}

					output.SetText(fmt.Sprintf("Output File: %s\n", outputFile))
					closer.Close()

					promptDialog("Enter...", "Recipient Private Key",`^[A-Za-z0-9_-]+$`, true, win, func(privateKey string) {
						n, err := crypto.DecryptFile(inputFile, outputFile, func() (s string, err error) {
							return privateKey, nil
						})
						if err != nil {
							output.SetText(err.Error())
						} else {
							output.SetText(fmt.Sprintf("Input File '%s'\nWas decrypted to\nOutput File '%s'\nWritten %d bytes.", inputFile, outputFile, n))
						}
					})

				}, win)
				_, fileName := filepath.Split(outputFile)
				fs.SetFileName(fileName)
				fs.Show()

			}, win)
			fd.SetFilter(storage.NewExtensionFileFilter([]string{".cp"}))
			fd.Show()
		}),
		output,
	))

	win.Resize(fyne.NewSize(640, 460))
	win.ShowAndRun()
	return nil
}

func promptDialog(welcome, field, regex string, password bool, parent fyne.Window, cb func(string)) {
	var entry *widget.Entry
	if password {
		entry = widget.NewPasswordEntry()
	} else {
		entry = widget.NewEntry()
	}
	entry.Validator = validation.NewRegexp(regex, field)
	items := []*widget.FormItem{
		widget.NewFormItem(field, entry),
	}
	form := dialog.NewForm(welcome, "Submit", "Cancel", items, func(b bool) {
		if !b {
			return
		}
		cb(strings.TrimSpace(entry.Text))
	}, parent)
	form.Resize(fyne.NewSize(440, 260))
	form.Show()
}

