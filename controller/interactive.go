package controller

/*
TODO:
[ ] handleInteractiveMessages is trying to get the API from the request but I'm not sure its looking in the right place
	Probably should look under `payload`
*/
import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/edwardIshaq/slack"
)

type interactive struct {
	components []interactiveComponent
}

func (i *interactive) addComponent(comp interactiveComponent) {
	i.components = append(i.components, comp)
}

type interactiveComponent interface {
	canHandleCallback(string) bool
	handleCallback(http.ResponseWriter, *http.Request, slack.AttachmentActionCallback)
}

func (i *interactive) registerRoutes() {
	http.HandleFunc("/interactive", i.handleInteractiveMessages)
}

func (i *interactive) callbackHandler(callbackID string) interactiveComponent {
	for _, comp := range i.components {
		if comp.canHandleCallback(callbackID) {
			return comp
		}
	}
	return nil
}

func (i *interactive) handleInteractiveMessages(w http.ResponseWriter, r *http.Request) {
	// Scan `components` to see who can hanle this callbackID
	actionCallback := parseAttachmentActionCallback(r)
	fmt.Println("Received interactive callback", debugJSON(actionCallback))
	if comp := i.callbackHandler(actionCallback.CallbackID); comp != nil {
		comp.handleCallback(w, r, actionCallback)
		return
	}

	// payload := r.PostFormValue("payload")
	// fmt.Printf("payload= %v\n", payload)
	http.NotFound(w, r)
}

func debugJSON(v interface{}) string {
	b, err := json.MarshalIndent(v, " ", "\t")
	if err != nil {
		return fmt.Sprintf("indent JSON failed with error: %v", err)
	}
	return string(b)
}

func replyMessage(params *slack.Msg, responseURL string) {
	data, _ := json.Marshal(params)
	bodyReader := bytes.NewReader(data)
	req, err := http.NewRequest(http.MethodPost, responseURL, bodyReader)

	//Fire the request
	resp, err := slack.HTTPClient.Do(req)
	if err != nil {
		fmt.Println("\nResponseError: ", err)
		return
	}
	defer resp.Body.Close()
}

func postDialog(dialog slack.Dialog, triggerID, token string) error {
	dialogjson, err := json.Marshal(dialog)
	if err != nil {
		return fmt.Errorf("Error converting dialog to json: %v", err)
	}

	values := url.Values{
		"token":      {token},
		"trigger_id": {triggerID},
		"dialog":     {string(dialogjson)},
	}

	reqBody := strings.NewReader(values.Encode())
	req, err := http.NewRequest("POST", slack.SLACK_API+"dialog.open", reqBody)
	if err != nil {
		fmt.Println("error happened: ", err)
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req = req.WithContext(context.Background())

	//Fire the request
	resp, err := slack.HTTPClient.Do(req)
	if err != nil {
		fmt.Println("\nResponseError: ", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Status Code not OK! , got: %d", resp.StatusCode)
	}

	bodyBytes, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		fmt.Printf("error reading body %v", err2)
		return err2
	}

	bodyString := string(bodyBytes)
	fmt.Println("\nbodyString: ", bodyString)
	dialogResponse := &slack.DialogOpenResponse{}
	json.Unmarshal(bodyBytes, &dialogResponse)
	fmt.Printf("\n\nDialogOpenResponse: %v\n\n", dialogResponse)
	return nil

}

func parseAttachmentActionCallback(r *http.Request) slack.AttachmentActionCallback {
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
