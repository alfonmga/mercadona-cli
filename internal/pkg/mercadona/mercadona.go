package mercadona

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/mitchellh/mapstructure"
)

const apiBaseEndpoint string = "https://tienda.mercadona.es/api/"
const requiredMinOrderPrice = 50.00

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
	Product             MercadonaProduct `mapstructure:"product"`
	RecommendedQuantity int              `mapstructure:"recommended_quantity"`
}
type MercadonaProduct struct {
	Id string `mapstructure:"id"`
}
type MercadonaOrder struct {
	Id           int    `mapstructure:"id"`
	ChangesUntil string `mapstructure:"changes_until"`
	StatusUI     string `mapstructure:"status_ui"`
}
type MercadonaCustomerInfo struct {
	CartId            string `json:"cart_id"`
	CurrentPostalCode string `json:"current_postal_code"`
	Email             string `json:"email"`
}
type addProductToCartResponse struct {
	Summary struct {
		Total string `json:"total"`
	} `json:"summary"`
}
type cartProductLine struct {
	Version   int      `json:"version"`
	ProductId string   `json:"product_id"`
	Quantity  int      `json:"quantity"`
	Sources   []string `json:"sources"`
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
	if _, err := os.Stat(getMercadonaCLIAuthenticationDataPath()); os.IsNotExist(err) {
		return
	}

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

func fetchOrders(page string, statusCode string) string {
	client := getAuthenticatedClient()
	resp, err := client.R().Get(fmt.Sprintf("%s%s", apiBaseEndpoint, fmt.Sprintf("customers/%s/orders/?status=%s&page=%s", currentAccountAuthData.CustomerId, statusCode, page)))
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

func addProductsToCart(cartId string, lines []cartProductLine) addProductToCartResponse {
	reqBody, err := json.Marshal(map[string]interface{}{
		"id":      cartId,
		"version": 1,
		"lines":   lines,
	})
	if err != nil {
		log.Fatalf("logInReqBody err: %v", err)
	}
	// log.Printf("%s", reqBody)

	var result addProductToCartResponse
	client := getAuthenticatedClient()
	resp, err := client.R().SetResult(&result).SetHeader("Content-Type", "application/json").SetBody(bytes.NewBuffer(reqBody)).Put(fmt.Sprintf("%s%s", apiBaseEndpoint, fmt.Sprintf("customers/%s/cart/", currentAccountAuthData.CustomerId)))
	if err != nil {
		log.Fatalf("Error add product to cart request: %s", err)
	}

	if resp.StatusCode() != 200 {
		log.Printf("addProductToCart resp body: %s", string(resp.Body()))
		log.Fatalf("Error adding product to cart: %v", resp.StatusCode())
	}

	return result
}

func init() {
	loadAuthenticationDataFromFile()
}

func Authenticate(credentials MercadonaLogInCredentialsBodyRequest) {
	authData := logIn(credentials)
	persistAuthentication(authData)
	log.Printf("Now you are authenticated! customer id: %v", authData.CustomerId)
}

func CustomerInfo() MercadonaCustomerInfo {
	assertAccessToken()

	client := getAuthenticatedClient()
	var result MercadonaCustomerInfo
	resp, err := client.R().SetResult(&result).Get(fmt.Sprintf("%s%s", apiBaseEndpoint, fmt.Sprintf("customers/%s/", currentAccountAuthData.CustomerId)))
	if err != nil {
		log.Fatalf("CustomerInfo request error: %v", err)
	}

	if resp.StatusCode() != 200 {
		log.Printf("resp body: %s", string(resp.Body()))
		log.Fatalf("resp.StatusCode error: %v", resp.StatusCode())
	}

	return result
}

type mercadonaCheckoutResult struct {
	Id      int `json:"id"`
	Address struct {
		Id         int    `json:"id"`
		Name       string `json:"address"`
		Details    string `json:"address_detail"`
		Comments   string `json:"comments"`
		PostalCode string `json:"postal_code"`
		Town       string `json:"town"`
	} `json:"address"`
	CustomerPhone string `json:"customer_phone"`
	PaymentMethod struct {
		CreditCardNumber string `json:"credit_card_number"`
	} `json:"payment_method"`
	OrderId int    `json:"order_id"`
	Price   string `json:"price"`
	Summary struct {
		Total string `json:"total"`
	} `json:"summary"`
	WarehouseCode string `json:"warehouse_code"`
	ChangesUntil  string `json:"changes_until"`
	Slot          struct {
		Start string `json:"start"`
		End   string `json:"end"`
	} `json:"slot"`
}

func checkoutCart(cartId string, lines []cartProductLine) mercadonaCheckoutResult {
	reqBody, err := json.Marshal(map[string]interface{}{
		"cart": map[string]interface{}{
			"id":      cartId,
			"version": 25,
			"lines":   lines,
		},
	})
	if err != nil {
		log.Fatalf("logInReqBody err: %v", err)
	}
	// log.Printf("%s", reqBody)

	var result mercadonaCheckoutResult
	client := getAuthenticatedClient()
	resp, err := client.R().SetResult(&result).SetHeader("Content-Type", "application/json").SetBody(bytes.NewBuffer(reqBody)).Post(fmt.Sprintf("%s%s", apiBaseEndpoint, fmt.Sprintf("customers/%s/checkouts/", currentAccountAuthData.CustomerId)))
	if err != nil {
		log.Fatalf("Error checkoutCart request: %s", err)
	}

	if resp.StatusCode() != 201 {
		log.Printf("checkoutCart resp body: %s", string(resp.Body()))
		log.Fatalf("Error checkoutCart: %v", resp.StatusCode())
	}

	return result
}

type mercadonaAddressSlot struct {
	Id        string `json:"id"`
	Start     string `json:"start"`
	End       string `json:"end"`
	Open      bool   `json:"open"`
	Available bool   `json:"available"`
}

func fetchAddressSlots(addressId int) []mercadonaAddressSlot {
	var paginatedResults mercadonaPagination

	client := getAuthenticatedClient()
	resp, err := client.R().SetResult(&paginatedResults).Get(fmt.Sprintf("%s%s", apiBaseEndpoint, fmt.Sprintf("customers/%s/addresses/%v/slots/", currentAccountAuthData.CustomerId, addressId)))
	if err != nil {
		log.Fatalf("fetchAddressSlots request error: %v", err)
	}

	if resp.StatusCode() != 200 {
		log.Printf("fetchAddressSlots resp body: %s", string(resp.Body()))
		log.Fatalf("fetchAddressSlots resp.StatusCode error: %v", resp.StatusCode())
	}

	if len(paginatedResults.Results) == 0 {
		log.Fatal("No slots found.")
	}

	var products []mercadonaAddressSlot
	err = mapstructure.Decode(paginatedResults.Results, &products)
	if err != nil {
		log.Fatalf("Unable to decode fetchAddressSlots paginatedResults.Results: %v", err)
	}

	return products
}

func setCheckoutOrderDeliveryInfo(checkoutId int, addressId int, slotId string) {
	reqBody, err := json.Marshal(map[string]interface{}{
		"address": map[string]int{
			"id": addressId,
		},
		"slot": map[string]string{
			"id": slotId,
		},
	})
	if err != nil {
		log.Fatalf("setCheckoutOrderDeliveryInfo reqbody err: %v", err)
	}
	// log.Printf("%s", reqBody)

	var result addProductToCartResponse
	client := getAuthenticatedClient()
	resp, err := client.R().SetResult(&result).SetHeader("Content-Type", "application/json").SetBody(bytes.NewBuffer(reqBody)).Put(fmt.Sprintf("%s%s", apiBaseEndpoint, fmt.Sprintf("customers/%s/checkouts/%v/delivery-info/", currentAccountAuthData.CustomerId, checkoutId)))
	if err != nil {
		log.Fatalf("Error setCheckoutOrderDeliveryInfo request: %s", err)
	}

	if resp.StatusCode() != 200 {
		log.Printf("setCheckoutOrderDeliveryInfo resp body: %s", string(resp.Body()))
		log.Fatalf("Error setCheckoutOrderDeliveryInfo: %v", resp.StatusCode())
	}
}

func submitCheckoutOrder(checkoutId int) mercadonaCheckoutResult {
	var result mercadonaCheckoutResult
	client := getAuthenticatedClient()
	resp, err := client.R().SetResult(&result).SetHeader("Content-Type", "application/json").Post(fmt.Sprintf("%s%s", apiBaseEndpoint, fmt.Sprintf("customers/%s/checkouts/%v/orders/", currentAccountAuthData.CustomerId, checkoutId)))
	if err != nil {
		log.Fatalf("Error setCheckoutOrderDeliveryInfo request: %s", err)
	}

	if resp.StatusCode() != 201 {
		log.Printf("setCheckoutOrderDeliveryInfo resp body: %s", string(resp.Body()))
		log.Fatalf("Error setCheckoutOrderDeliveryInfo: %v", resp.StatusCode())
	}

	return result
}

func MakeNewOrder(slotDate string) {
	assertAccessToken()

	customerInfo := CustomerInfo()
	cartId := customerInfo.CartId

	recommendedProducts := fetchRecommendedProducts()
	// fmt.Printf("%v", recommendedProducts)

	var products []cartProductLine

	for _, _p := range recommendedProducts {
		products = append(products, cartProductLine{Version: 25, ProductId: _p.Product.Id, Quantity: _p.RecommendedQuantity, Sources: make([]string, 0)})
		cartRes := addProductsToCart(cartId, products)

		cartValueAfterProductAdded, err := strconv.ParseFloat(cartRes.Summary.Total, 64)
		if err != nil {
			log.Fatalf("Error parsing float64: %s", err)
		}
		// log.Printf("#%v Cart price after last product added: %v", i, cartValueAfterProductAdded)

		if cartValueAfterProductAdded >= requiredMinOrderPrice {
			break
		}
	}

	checkoutRes := checkoutCart(cartId, products)
	availableSlots := fetchAddressSlots(checkoutRes.Address.Id)

	tLoc, err := time.LoadLocation("Europe/Madrid")
	if err != nil {
		log.Fatalf("Err loading Europe/Madrid location: %s", err)
	}

	tTargetSlot, err := time.ParseInLocation("2006-01-02 15:04", slotDate, tLoc)
	if err != nil {
		log.Fatalf("Error parsing slot target date: %s", err)
	}

	var targetSlotId string = ""
	for _, _slot := range availableSlots {
		t, err := time.ParseInLocation(time.RFC3339, _slot.Start, tLoc)
		if err != nil {
			log.Fatalf("Error parsing slot start date: %s", err)
		}

		if t.Equal(tTargetSlot) {
			if !_slot.Open {
				log.Fatal("Slot date is not open")
			}
			if !_slot.Available {
				log.Fatal("Slot date is not available")
			}
			targetSlotId = _slot.Id
		}
	}
	if targetSlotId == "" {
		log.Fatal("Target slot id not found")
	}

	setCheckoutOrderDeliveryInfo(checkoutRes.Id, checkoutRes.Address.Id, targetSlotId)
	checkoutOrderSubmissionRes := submitCheckoutOrder(checkoutRes.Id)

	log.Println("-- New order created successfully!! --")
	log.Printf("Shipping to: %s, %s, %s, %s", checkoutOrderSubmissionRes.Address.Name, checkoutOrderSubmissionRes.Address.Details, checkoutOrderSubmissionRes.Address.PostalCode, checkoutOrderSubmissionRes.Address.Town)
	log.Printf("Shipping from Warehouse: %s", checkoutOrderSubmissionRes.WarehouseCode)

	shippingTime, err := time.ParseInLocation(time.RFC3339, checkoutOrderSubmissionRes.Slot.Start, tLoc)
	if err != nil {
		log.Fatalf("Error parsing shipping time: %s", err)
	}
	log.Printf("Shipping date: %s", shippingTime.Format("2006-01-02 15:04"))

	log.Printf("Payment method: ****%s", checkoutOrderSubmissionRes.PaymentMethod.CreditCardNumber)
	log.Printf("Given contact phone number: %s", checkoutOrderSubmissionRes.CustomerPhone)
	log.Printf("Total order price (shipping included): %v EUR", checkoutOrderSubmissionRes.Summary.Total)

	changesUntilTime, err := time.ParseInLocation(time.RFC3339, checkoutOrderSubmissionRes.ChangesUntil, tLoc)
	if err != nil {
		log.Fatalf("Error parsing changes_until time: %s", err)
	}
	log.Printf("Changes until date: %s", changesUntilTime.Format("2006-01-02 15:04"))

	log.Printf("https://tienda.mercadona.es/user-area/orders/%v", checkoutRes.Id)
}

func ListAllOrders(page string) string {
	assertAccessToken()
	return fetchOrders(page, "")
}

func GetActiveOrderModifyURL() string {
	assertAccessToken()

	ordersFetchResp := fetchOrders("1", "2")

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
