// TODO: store customer id too

package mercadona

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/mitchellh/mapstructure"
)

const apiBaseEndpoint string = "https://tienda.mercadona.es/api/"

var accessToken string

type mercadonaPagination struct {
	NextPage string                   `json:"next_page"`
	Results  []map[string]interface{} `json:"results"`
}
type MercadonaLogInCredentialsBodyRequest struct {
	Email    string
	Password string
}
type mercadonaLogInBodyResponse struct {
	AccessToken string `json:"access_token"`
	CustomerId  string `json:"customer_id"`
}
type MercadonaRecommendedProduct struct {
	Product MercadonaProduct `mapstructure:"product"`
}
type MercadonaProduct struct {
	Id string `mapstructure:"id"`
}
type MercadonaOrder struct {
	Id           int    `mapstructure:"id"`
	ChangesUntil string `mapstructure:"changes_until"`
	StatusUI     string `mapstructure:"status_ui"`
}

func persistAccessToken(accessToken string) {
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("persistAccessToken unable to get home dir: %v", err)
	}
	ioutil.WriteFile(fmt.Sprintf("%s/.mercadona-cli-access_token", dirname), []byte(accessToken), 0666)
}

func loadStoredAccessToken() {
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("getStoredAccessToken unable to get home dir: %v", err)
	}
	accessTokenBlob, err := ioutil.ReadFile(fmt.Sprintf("%s/.mercadona-cli-access_token", dirname))
	if err != nil {
		log.Fatalf("getStoredAccessToken unable to read stored access token: %v", err)
	}

	accessToken = string(accessTokenBlob)
}

func assertAccessToken() {
	if accessToken == "" {
		log.Fatal("Authentication is required")
	}
}

func logIn(credentials MercadonaLogInCredentialsBodyRequest) mercadonaLogInBodyResponse {
	client := &http.Client{}

	logInReqBody, err := json.Marshal(map[string]string{
		"username": credentials.Email,
		"password": credentials.Password,
	})
	if err != nil {
		log.Fatalf("logInReqBody err: %v", err)
	}

	logInResp, err := client.Post(fmt.Sprintf("%s%s", apiBaseEndpoint, "auth/tokens/"), "application/json", bytes.NewBuffer(logInReqBody))
	if err != nil {
		log.Fatalf("logIn req error: %v", err)
	}

	defer logInResp.Body.Close()
	body, err := io.ReadAll(logInResp.Body)
	if err != nil {
		log.Fatalf("logInResp body read err: %v", err)
	}

	if logInResp.StatusCode != 200 {
		log.Printf("logInResp body: %s", []byte(body))
		log.Fatalf("logInResp.StatusCode error: %v", logInResp.StatusCode)
	}

	var logInResponse mercadonaLogInBodyResponse
	err = json.Unmarshal([]byte(body), &logInResponse)
	if err != nil {
		log.Fatalf("logInResponse unmarshal error: %v", err)
	}

	return logInResponse
}

func fetchOrders(page string) string {
	client := resty.New()
	resp, err := client.R().SetAuthToken(accessToken).Get(fmt.Sprintf("%s%s", apiBaseEndpoint, fmt.Sprintf("customers/%s/orders/?page=%s", "<customer_id>", page)))
	if err != nil {
		log.Fatalf("ListAllOrders request error: %v", err)
	}

	if resp.StatusCode() != 200 {
		log.Printf("resp body: %s", string(resp.Body()))
		log.Fatalf("resp.StatusCode error: %v", resp.StatusCode())
	}

	return string(resp.Body())
}

func fetchRecommendedProducts() []MercadonaRecommendedProduct {
	client := resty.New()
	resp, err := client.R().SetAuthToken(accessToken).Get(fmt.Sprintf("%s%s", apiBaseEndpoint, fmt.Sprintf("customers/%s/recommendations/myregulars/precision/", "<customer_id>")))
	if err != nil {
		log.Fatalf("fetchRegularProducts request error: %v", err)
	}

	if resp.StatusCode() != 200 {
		log.Printf("fetchRegularProducts resp body: %s", string(resp.Body()))
		log.Fatalf("fetchRegularProducts resp.StatusCode error: %v", resp.StatusCode())
	}

	var paginatedResults mercadonaPagination
	err = json.Unmarshal([]byte(resp.Body()), &paginatedResults)
	if err != nil {
		log.Fatalf("resp unmarshal error: %v", err)
	}

	if len(paginatedResults.Results) == 0 {
		log.Fatal("No favorite products found.")
	}

	var products []MercadonaRecommendedProduct
	err = mapstructure.Decode(paginatedResults.Results, &products)
	if err != nil {
		log.Fatalf("Unable to decode paginatedResults.Results: %v", err)
	}

	return products
}

func addProductToCart(productId string, quantity int) {
	// TODO: …
}

func checkoutCart(cartId string) {
	// TODO: …
}

func getCheckoutSlots(checkoutId string) {
	// TODO: …
}

func authorizeCheckoutPayment(checkoutId string) {
	// TODO: …
}

func init() {
	loadStoredAccessToken()
}

func Authenticate(credentials MercadonaLogInCredentialsBodyRequest) {
	resp := logIn(credentials)
	persistAccessToken(resp.AccessToken)
	log.Printf("Now you are authenticated! customer id: %v", resp.CustomerId)
}

func CustomerInfo() {
	assertAccessToken()

	client := resty.New()
	resp, err := client.R().SetAuthToken(accessToken).Get(fmt.Sprintf("%s%s", apiBaseEndpoint, fmt.Sprintf("customers/%s/", "<customer_id>")))
	if err != nil {
		log.Fatalf("CustomerInfo request error: %v", err)
	}

	if resp.StatusCode() != 200 {
		log.Printf("resp body: %s", string(resp.Body()))
		log.Fatalf("resp.StatusCode error: %v", resp.StatusCode())
	}

	fmt.Println(string(resp.Body()))
}

func MakeNewOrder() {
	assertAccessToken()

	recommendedProducts := fetchRecommendedProducts()
	fmt.Printf("%v", recommendedProducts)
	// TODO: …
}

func ListAllOrders(page string) string {
	assertAccessToken()
	return fetchOrders(page)
}

func GetActiveOrderModifyURL() string {
	assertAccessToken()

	ordersFetchResp := fetchOrders("1")

	var paginatedOrders mercadonaPagination
	err := json.Unmarshal([]byte(ordersFetchResp), &paginatedOrders)
	if err != nil {
		log.Fatalf("fetchOrdersResponse unmarshal error: %v", err)
	}

	if len(paginatedOrders.Results) == 0 {
		log.Println("No orders found.")
		return ""
	}

	var latestOrder MercadonaOrder
	err = mapstructure.Decode(paginatedOrders.Results[0], &latestOrder)
	if err != nil {
		log.Fatalf("Unable to decode paginatedOrders.Results[0]: %v", err)
	}

	if latestOrder.StatusUI != "confirmed" {
		log.Println("No active order found.")
		return ""
	}

	return fmt.Sprintf("https://tienda.mercadona.es/orders/%v/edit/products?my-regulars=true", latestOrder.Id)
}
