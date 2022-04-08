package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"
)

// Elavon -
func Elavon(w http.ResponseWriter, r *http.Request) {

	// Provide Converge Credentials
	//Converge 6-Digit Account ID *Not the 10-Digit Elavon Merchant ID*
	//const MERCHANTID = 631103
	const MERCHANTID = "0022805"
	//Converge User ID *MUST FLAG AS HOSTED API USER IN CONVERGE UI*
	//const USERID = "webpage"
	const USERID = "apiuser"
	//Converge PIN (64 CHAR A/N)
	//const PIN = "80KYG17V8IBW89MTJYJZIQ3C31DCCG9BJRYP9IYZ4D83ZGQEHCUQDVZB2YBSIG7S"
	const PIN = "DZKLNA2M7ZAAE8X6W4AD6T3IY4X1LFR41IG0B2ZBZDAYC6DO3JECM4IC4IM8QTR9"
	const CVVINDICATOR = '1' //means "present"
	//demo url
	const ELAVONURL = "https://demo.myvirtualmerchant.com/VirtualMerchantDemo/process.do"
	//const ELAVONURL = "https://www.myvirtualmerchant.com/VirtualMerchant/process.do"
	const PAYSUCCESSURL = "/paysuccess.php"
	const PAYFAILURL = "/payerror.php"
	const CARDTYPE = "CREDITCARD"

	type elavonPayload struct {
		sslmerchantid      string     //`json:"ssl_merchant_id"`
		ssluserid          string  //`json:"ssl_user_id"`
		sslpin             string  //`json:"ssl_pin"`
		ssltransactiontype string  //`json:"ssl_transaction_type"`
		sslamount          float32 //`json:"ssl_amount"`
	}

	// URL to Converge demo session token server
	url := "https://api.demo.convergepay.com/hosted-payments/transaction_token"
	// URL to Converge production session token server
	//url := "https://api.convergepay.com/hosted-payments/transaction_token"

	// URL to the demo Hosted Payments Page
	//hppurl := "https://api.demo.convergepay.com/hosted-payments"
	// URL to the production Hosted Payments Page
	//hppurl := "https://api.convergepay.com/hosted-payments"

	payload := elavonPayload{}

	payload.sslmerchantid = MERCHANTID
	payload.ssluserid = USERID
	payload.sslpin = PIN
	payload.ssltransactiontype = "CCSALE"

	//hardcoded amount for testing
	payload.sslamount = 1.00

	err := r.ParseForm()
	//Capture ssl_amount as POST data
	//amount := r.FormValue("ssl_amount")

	//Capture ssl_first_name as POST data
	//firstname  := r.FormValue("ssl_first_name")

	//Capture ssl_last_name as POST data
	//lastname  := r.FormValue("ssl_last_name")

	//Capture ssl_merchant_txn_id as POST data
	//merchanttxid := r.FormValue("ssl_merchant_txn_id")

	//Capture ssl_invoice_number as POST data
	//invoicenumber := r.FormValue("ssl_invoice_number")

	//Follow the above pattern to add additional fields to be sent in curl request below.

	payloadBytes, err := json.Marshal(payload)

	body := bytes.NewReader(payloadBytes)

	if err != nil {
		log.Printf("Error marshaling json %s", err)
	}

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	req, err := http.NewRequest("POST", url, body)

	if err != nil {
		log.Printf("Error creating post request to %s: %s", url, err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Printf("Error sending http request")
	}

	defer resp.Body.Close()

}
