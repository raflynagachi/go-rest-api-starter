package config

type App struct {
	Name string `json:"name"`
	Port int    `json:"port"`
}
type Database struct {
	Name     string `json:"name"`
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
}

type Redis struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
}

type ServiceEndpoints struct {
	NatsURL      string `json:"nats_url"`
	OrderURL     string `json:"order_url"`
	PaymentURL   string `json:"payment_url"`
	InventoryURL string `json:"inventory_url"`
}
