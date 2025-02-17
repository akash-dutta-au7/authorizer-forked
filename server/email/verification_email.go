package email

import (
	"log"

	"github.com/authorizerdev/authorizer/server/constants"
	"github.com/authorizerdev/authorizer/server/envstore"
)

// SendVerificationMail to send verification email
func SendVerificationMail(toEmail, token, hostname string) error {
	// The receiver needs to be in slice as the receive supports multiple receiver
	Receiver := []string{toEmail}

	Subject := "Please verify your email"
	message := `
	<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
    <html xmlns="http://www.w3.org/1999/xhtml" xmlns:o="urn:schemas-microsoft-com:office:office">
        <head>
            <meta charset="UTF-8">
            <meta content="width=device-width, initial-scale=1" name="viewport">
            <meta name="x-apple-disable-message-reformatting">
            <meta http-equiv="X-UA-Compatible" content="IE=edge">
            <meta content="telephone=no" name="format-detection">
            <title></title>
            <!--[if (mso 16)]>
            <style type="text/css">
            a {}
            </style>
            <![endif]-->
            <!--[if gte mso 9]><style>sup { font-size: 100%% !important; }</style><![endif]-->
            <!--[if gte mso 9]>
        <xml>
            <o:OfficeDocumentSettings>
            <o:AllowPNG></o:AllowPNG>
            <o:PixelsPerInch>96</o:PixelsPerInch>
            </o:OfficeDocumentSettings>
        </xml>
        <![endif]-->
        </head>
        <body style="font-family: sans-serif;">
            <div class="es-wrapper-color">
                <!--[if gte mso 9]>
                    <v:background xmlns:v="urn:schemas-microsoft-com:vml" fill="t">
                        <v:fill type="tile" color="#ffffff"></v:fill>
                    </v:background>
                <![endif]-->
                <table class="es-wrapper" width="100%%" cellspacing="0" cellpadding="0">
                    <tbody>
                        <tr>
                            <td class="esd-email-paddings" valign="top">
                                <table class="es-content esd-footer-popover" cellspacing="0" cellpadding="0" align="center">
                                    <tbody>
                                        <tr>
                                            <td class="esd-stripe" align="center">
                                                <table class="es-content-body" style="border-left:1px solid transparent;border-right:1px solid transparent;border-top:1px solid transparent;border-bottom:1px solid transparent;padding:20px 0px;" width="600" cellspacing="0" cellpadding="0" bgcolor="#ffffff" align="center">
                                                    <tbody>
                                                        <tr>
                                                            <td class="esd-structure es-p20t es-p40b es-p40r es-p40l" esd-custom-block-id="8537" align="left">
                                                                <table width="100%%" cellspacing="0" cellpadding="0">
                                                                    <tbody>
                                                                        <tr>
                                                                            <td class="esd-container-frame" width="518" align="left">
                                                                                <table width="100%%" cellspacing="0" cellpadding="0">
                                                                                    <tbody>
                                                                                        <tr>
                                                                                            <td class="esd-block-image es-m-txt-c es-p5b" style="font-size:0;padding:10px" align="center"><a target="_blank" clicktracking="off"><img src="{{.org_logo}}" alt="icon" style="display: block;" title="icon" width="30"></a></td>
                                                                                        </tr>
                                                                                        
                                                                                        <tr style="background: rgb(249,250,251);padding: 10px;margin-bottom:10px;border-radius:5px;">
                                                                                            <td class="esd-block-text es-m-txt-c es-p15t" align="center" style="padding:10px;padding-bottom:30px;">
                                                                                                <p>Hey there 👋</p>
                                                                                                <p>We have received request to verify email for <b>{{.org_name}}</b>. If this is correct, please confirm your email address by clicking the button below.</p> <br/>
                                                                                                <a 
                                                                                                clicktracking="off" href="{{.verification_url}}" class="es-button" target="_blank" style="text-decoration: none;padding:10px 15px;background-color: rgba(59,130,246,1);color: #fff;font-size: 1em;border-radius:5px;">Confirm Email</a>
                                                                                            </td>
                                                                                        </tr>
                                                                                    </tbody>
                                                                                </table>
                                                                            </td>
                                                                        </tr>
                                                                    </tbody>
                                                                </table>
                                                            </td>
                                                        </tr>
                                                    </tbody>
                                                </table>
                                            </td>
                                        </tr>
                                    </tbody>
                                </table>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>
            <div style="position: absolute; left: -9999px; top: -9999px; margin: 0px;"></div>
        </body>
    </html>
	`
	data := make(map[string]interface{}, 3)
	data["org_logo"] = envstore.EnvStoreObj.GetStringStoreEnvVariable(constants.EnvKeyOrganizationLogo)
	data["org_name"] = envstore.EnvStoreObj.GetStringStoreEnvVariable(constants.EnvKeyOrganizationName)
	data["verification_url"] = hostname + "/verify_email?token=" + token
	message = addEmailTemplate(message, data, "verify_email.tmpl")
	// bodyMessage := sender.WriteHTMLEmail(Receiver, Subject, message)

	err := SendMail(Receiver, Subject, message)
	if err != nil {
		log.Println("=> error sending email:", err)
	}
	return err
}
