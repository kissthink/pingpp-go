package pingpp

import (
	"encoding/json"
	"net/url"
	"strconv"
)

type ChargeClient struct {
	key     string
	backend Backend
}

func GetChargeClient(key string) *ChargeClient {
	return &ChargeClient{
		key:     key,
		backend: getBackend(),
	}
}

func (chargeClient *ChargeClient) SetKey(key string) {
	chargeClient.key = key
}

func (chargeClient *ChargeClient) New(params *ChargeParams) (*Charge, error) {
	var charge Charge
	body := &url.Values{
		"amount": {strconv.FormatUint(params.Amount, 10)},
	}

	body.Add("order_no", params.Order_no)
	appstring, _ := json.Marshal(params.App)
	body.Add("app", string(appstring))
	body.Add("channel", params.Channel)
	body.Add("currency", params.Currency)
	body.Add("subject", params.Subject)
	body.Add("body", params.Body)
	body.Add("client_ip", params.Client_ip)

	if len(params.Metadata) > 0 {
		metastring, _ := json.Marshal(params.Metadata)
		body.Add("metadata", string(metastring))
	}

	if params.Time_expire > 0 {
		body.Add("time_expire", strconv.FormatUint(params.Time_expire, 10))
	}

	if len(params.Description) > 0 {
		body.Add("description", string(params.Description))
	}

	chargestring, _ := json.Marshal(params)
	resp_byte, err := chargeClient.backend.CallJson("POST", "charges", chargeClient.key, chargestring, &charge)
	json.Unmarshal(resp_byte, &charge)
	return &charge, err
}

func (chargeClient *ChargeClient) Get(charge_id string) (*Charge, error) {
	var charge Charge
	body := &url.Values{}
	resp_byte, err := chargeClient.backend.Call("GET", "charges/"+charge_id, chargeClient.key, body, &charge)
	json.Unmarshal(resp_byte, &charge)
	return &charge, err

}

func (chargeClient *ChargeClient) List(params *ChargeListParams) (*ChargeList, error) {
	var charges ChargeList
	body := &url.Values{}
	if params.Limit > 100 || params.Limit <= 0 {
		body.Add("limit", strconv.FormatUint(10, 10))
	} else {
		body.Add("limit", strconv.FormatUint(params.Limit, 10))
	}
	if params.Start_after != "" {
		body.Add("starting_after", params.Start_after)
	}
	if params.End_before != "" {
		body.Add("ending_before", params.End_before)
	}
	if params.Createdgt != "" {
		body.Add("created[gt]", params.Createdgt)
	}
	if params.Createdgte != "" {
		body.Add("created[gte]", params.Createdgte)
	}
	if params.Createdlt != "" {
		body.Add("created[lt]", params.Createdlt)
	}
	if params.Createdlte != "" {
		body.Add("created[lte]", params.Createdlte)
	}
	if params.Appid != "" {
		body.Add("app[id]", params.Appid)
	}
	if params.Channel != "" {
		body.Add("channel", params.Channel)
	}
	if params.Paid == 1 {
		body.Add("paid", "ture")
	} else if params.Paid == 2 {
		body.Add("paid", "false")
	}
	if params.Refunded == 1 {
		body.Add("refunded", "ture")
	} else if params.Refunded == 2 {
		body.Add("refunded", "false")
	}
	// body.Add("metadata", charge.metadata)
	resp_byte, err := chargeClient.backend.Call("GET", "charges", chargeClient.key, body, &charges)
	json.Unmarshal(resp_byte, &charges)
	return &charges, err
}
