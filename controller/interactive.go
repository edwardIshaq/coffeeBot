package controller

/*
TODO:
[ ] handleInteractiveMessages is trying to get the API from the request but I'm not sure its looking in the right place
	Probably should look under `payload`
*/
import (
	"SlackPlatform/middleware"
	"SlackPlatform/models"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/nlopes/slack"
)

type interactive struct{}

func (i interactive) registerRoutes() {
	http.HandleFunc("/interactive", handleInteractiveMessages)
}

func handleInteractiveMessages(w http.ResponseWriter, r *http.Request) {
	token, ok := middleware.AccessToken(r.Context())
	if ok == false {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	//Process callback to extract `barista.dialog.<chosenBev>`
	actionCallback := parseAttachmentActionCallback(r)
	callbackID := actionCallback.CallbackID
	chosenBev := ""
	if strings.HasPrefix(actionCallback.CallbackID, "barista.dialog.") {
		callbackID = "barista.dialog"
		chosenBev = strings.Split(actionCallback.CallbackID, ".")[2]
	}

	switch callbackID {
	case "beverage_selection":
		fmt.Println("interacted with `menu`")
		if len(actionCallback.Actions) >= 1 {
			if len(actionCallback.Actions[0].SelectedOptions) >= 1 {
				chosenBeverage := actionCallback.Actions[0].SelectedOptions[0].Value
				postDialog(chosenBeverage, actionCallback.TriggerID, token)
			}
		}
		return

	case "barista.dialog":
		startTime := time.Now()
		dialogResponse := models.DialogSubmitCallback{}
		json.Unmarshal([]byte(r.PostFormValue("payload")), &dialogResponse)
		responseURL := dialogResponse.ResponseURL
		params := dialogResponse.FeedbackMessage(chosenBev)

		go func(params *slack.Msg, responseURL string) {
			data, _ := json.Marshal(params)
			bodyReader := bytes.NewReader(data)
			req, err := http.NewRequest(http.MethodPost, responseURL, bodyReader)
			fmt.Println("Request\n", req)

			//Fire the request
			resp, err := slack.HTTPClient.Do(req)
			if err != nil {
				fmt.Println("\nResponseError: ", err)
				return
			}
			defer resp.Body.Close()
			fmt.Printf("RESPONSE: %v", resp)
			fmt.Printf("\nprocessing Dialog took: %s\n", time.Since(startTime))
		}(params, responseURL)
	}
}

func postDialog(chosenBeverage, triggerID, token string) {
	dialog := makeDialog(chosenBeverage)

	if dialogjson, err := json.Marshal(dialog); err == nil {
		postBody := url.Values{
			"token":      {token},
			"trigger_id": {triggerID},
			"dialog":     {string(dialogjson)},
		}

		reqBody := strings.NewReader(postBody.Encode())
		req, err := http.NewRequest("POST", slack.SLACK_API+"dialog.open", reqBody)
		if err != nil {
			fmt.Println("error happened: ", err)
			return
		}

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req = req.WithContext(context.Background())

		//Fire the request
		resp, err := slack.HTTPClient.Do(req)
		if err != nil {
			fmt.Println("\nResponseError: ", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			bodyBytes, err2 := ioutil.ReadAll(resp.Body)
			if err2 != nil {
				fmt.Printf("error reading body %v", err2)
				return
			}
			bodyString := string(bodyBytes)
			fmt.Println("\nbodyString: ", bodyString)
		}
	}
}

func makeDialog(chosenBeverage string) models.Dialog {
	presetBeverage := models.BeverageByName(chosenBeverage)

	cupMenu := models.NewStaticSelectDialogInput("CupType", "Drink Size", models.AllDrinkSizes())
	cupMenu.Value = presetBeverage.CupType

	espressMenu := models.NewStaticSelectDialogInput("Espresso", "Espresso Options", models.AllEspressoOptions())
	espressMenu.Value = presetBeverage.Espresso

	syrupMenu := models.NewStaticSelectDialogInput("Syrup", "Syrup", models.AllSyrupOptions())
	syrupMenu.Value = presetBeverage.Syrup

	tempMenu := models.NewStaticSelectDialogInput("Temperture", "Temperture", models.AllTemps())
	tempMenu.Value = presetBeverage.Temperture

	commentInput := models.NewTextAreaInput("Comment", "Comments")
	commentInput.Optional = true

	callbackID := "barista.dialog." + chosenBeverage

	dialog := models.Dialog{
		CallbackID:  callbackID,
		Title:       models.DialogTitle(chosenBeverage),
		SubmitLabel: "Order",
		Elements: []interface{}{
			cupMenu,
			espressMenu,
			syrupMenu,
			tempMenu,
			commentInput,
		},
	}
	return dialog
}

func parseAttachmentActionCallback(r *http.Request) slack.AttachmentActionCallback {
	r.ParseForm()
	payload := r.PostFormValue("payload")
	actionCallback := slack.AttachmentActionCallback{}
	json.Unmarshal([]byte(payload), &actionCallback)
	return actionCallback
}

func textReply(w http.ResponseWriter, text string) {
	w.Header().Set("Content-Type", "application/json")
	message := fmt.Sprintf(`{"text" : "%s", "replace_original": false}`, text)
	w.Write([]byte(message))
}
