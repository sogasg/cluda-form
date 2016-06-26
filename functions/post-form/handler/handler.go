package handler

import (
	"errors"
	"fmt"

	"github.com/cluda/cluda-form/functions/post-form/util"
)

// Handle will handel a new event
func Handle(e Event, conf Config, cli Clients) (interface{}, error) {

	resp, err := cli.Dynamo.GetItem(util.FormDataRequest(conf.FormFreeTable, e.Receiver))

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	if len(resp.Item) == 0 {
		// is a new email

		// add to db
		secret := util.RandString(10)
		_, err := cli.Dynamo.PutItem(util.NewFormDataPut(conf.FormFreeTable, e.Receiver, secret))

		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}

		// send confirm email
		templateData := util.EmailData{
			Text1:  "",
			Text2:  "To activate your form, please confirm your email address by clicking the link below.",
			Button: "Confirm email address",
			Secret: secret,
		}

		body, err := util.ParseTemplate("../email-templates/action.html", templateData)
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}

		resp, err := cli.Ses.SendEmail(util.SendEmialInput(conf.EmailFromAddres, "sogasg@gmail.com", "Test 22", body))
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}

		println(resp)

		return "verification email sent", nil
	} else if *resp.Item["verifyed"].BOOL {

		// add to submission table

		// send submission to asosiated email

		return "submission handled", nil
	} else {
		fmt.Println(e.Receiver, " not verifyed")
		return "", errors.New("receiver not verifyed")
	}

}