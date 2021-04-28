package mercadona

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/mitchellh/mapstructure"
)

const apiBaseEndpoint string = "https://tienda.mercadona.es/api/"

var currentAccountAuthData authenticationData

type authenticationData struct {
	AccessToken string `json:"access_token"`
	CustomerId  string `json:"customer_id"`
}
type mercadonaPagination struct {
	NextPage string                   `json:"next_page"`
	Results  []map[string]interface{} `json:"results"`
}
type MercadonaLogInCredentialsBodyRequest struct {
	Email    string
	Password string
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

func getMercadonaCLIAuthenticationDataPath() string {
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("persistAccessToken unable to get home dir: %v", err)
	}
	return fmt.Sprintf("%s/.mercadona-cli-access_token", dirname)
}

func persistAuthentication(authData authenticationData) {
	data, err := json.Marshal(authData)
	if err != nil {
		log.Fatalf("Error unable to marshal authData: %s", err)
	}

	err = ioutil.WriteFile(getMercadonaCLIAuthenticationDataPath(), data, 0666)
	if err != nil {
		log.Fatalf("Error persisting authData: %s", err)
	}
}
func loadAuthenticationDataFromFile() {
	data, err := ioutil.ReadFile(getMercadonaCLIAuthenticationDataPath())
	if err != nil {
		log.Fatalf("Error reading authentication data file: %s", err)
	}

	var authData authenticationData
	err = json.Unmarshal(data, &authData)
	if err != nil {
		log.Fatalf("Error unmarshaling authentication data: %s", err)
	}

	currentAccountAuthData = authData
}

func assertAccessToken() {
	if currentAccountAuthData.AccessToken == "" {
		log.Fatal("Authentication is required")
	}
}

func logIn(credentials MercadonaLogInCredentialsBodyRequest) authenticationData {
	logInReqBody, err := json.Marshal(map[string]string{
		"username": credentials.Email,
		"password": credentials.Password,
	})
	if err != nil {
		log.Fatalf("logInReqBody err: %v", err)
	}

	var logInResponse authenticationData

	client := resty.New()
	logInResp, err := client.R().SetResult(&logInResponse).SetHeader("Content-Type", "application/json").SetBody(bytes.NewBuffer(logInReqBody)).Post(fmt.Sprintf("%s%s", apiBaseEndpoint, "auth/tokens/"))
	if err != nil {
		log.Fatalf("logIn req error: %v", err)
	}

	if logInResp.StatusCode() != 200 {
		log.Printf("logInResp body: %s", []byte(logInResp.Body()))
		log.Fatalf("logInResp.StatusCode error: %v", logInResp.StatusCode())
	}

	return logInResponse
}

func getAuthenticatedClient() *resty.Client {
	return resty.New().SetAuthToken(currentAccountAuthData.AccessToken)
}

func fetchOrders(page string) string {
	client := getAuthenticatedClient()
	resp, err := client.R().Get(fmt.Sprintf("%s%s", apiBaseEndpoint, fmt.Sprintf("customers/%s/orders/?page=%s", currentAccountAuthData.CustomerId, page)))
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
	var paginatedResults mercadonaPagination

	client := getAuthenticatedClient()
	resp, err := client.R().SetResult(&paginatedResults).Get(fmt.Sprintf("%s%s", apiBaseEndpoint, fmt.Sprintf("customers/%s/recommendations/myregulars/precision/", currentAccountAuthData.CustomerId)))
	if err != nil {
		log.Fatalf("fetchRegularProducts request error: %v", err)
	}

	if resp.StatusCode() != 200 {
		log.Printf("fetchRegularProducts resp body: %s", string(resp.Body()))
		log.Fatalf("fetchRegularProducts resp.StatusCode error: %v", resp.StatusCode())
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
	loadAuthenticationDataFromFile()
}

func Authenticate(credentials MercadonaLogInCredentialsBodyRequest) {
	authData := logIn(credentials)
	persistAuthentication(authData)
	log.Printf("Now you are authenticated! customer id: %v", authData.CustomerId)
}

func CustomerInfo() {
	assertAccessToken()

	client := getAuthenticatedClient()
	resp, err := client.R().Get(fmt.Sprintf("%s%s", apiBaseEndpoint, fmt.Sprintf("customers/%s/", currentAccountAuthData.CustomerId)))
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
