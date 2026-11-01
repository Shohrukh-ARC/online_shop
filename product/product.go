package product

import (
	"errors"
	"sort"
	"strings"
)

type Product struct {
	Name  string
	Price float64
	Stock int
}

type ProductUpdate struct {
	Name  *string  `json:"name"`
	Price *float64 `json:"price"`
	Stock *int     `json:"stock"`
}

func (p *Product) Sell(quantity int) error {
	if quantity <= 0 {
		return errors.New("Quantity must be positive")
	}

	if quantity > p.Stock {
		return errors.New("Not enough stock")
	}

	p.Stock -= quantity
	return nil
}

func (p *Product) Restock(quantity int) (int, error) {
	if quantity <= 0 {
		return 0, errors.New("Quantity must be positive")
	}

	p.Stock += quantity
	return p.Stock, nil
}

func CountAvailableProducts(products []Product) int {
	count := 0

	for _, product := range products {
		if product.Stock > 0 {
			count++
		}
	}

	return count
}

func FindProductByName(products []Product, name string) (Product, bool) {
	for _, product := range products {
		if strings.EqualFold(product.Name, name) {
			return product, true
		}
	}

	return Product{}, false
}

func FilterProductsByPrice(products []Product, minPrice, maxPrice float64) []Product {
	filteredItems := []Product{}
	for _, product := range products {
		if product.Price <= maxPrice && product.Price >= minPrice {
			filteredItems = append(filteredItems, product)
		}
	}

	return filteredItems
}

func FilterAvailableProducts(products []Product) []Product {
	availableProducts := []Product{}

	for _, product := range products {
		if product.Stock > 0 {
			availableProducts = append(availableProducts, product)
		}
	}

	return availableProducts
}

func SortProductsByPrice(products []Product, ascending bool) []Product {
	sortedProducts := make([]Product, len(products))
	copy(sortedProducts, products)

	if ascending {
		sort.Slice(sortedProducts, func(i, j int) bool {
			return sortedProducts[i].Price < sortedProducts[j].Price
		})
	} else {
		sort.Slice(sortedProducts, func(i, j int) bool {
			return sortedProducts[i].Price > sortedProducts[j].Price
		})
	}

	return sortedProducts
}

func SearchProductsByName(products []Product, keyword string) []Product {
	keyword = strings.ToLower(keyword)

	foundProducts := []Product{}

	for _, product := range products {
		if len(keyword) == 0 {
			continue
		}

		if strings.Contains(strings.ToLower(product.Name), keyword) {
			foundProducts = append(foundProducts, product)
		}
	}

	return foundProducts
}

func GetProductsPage(products []Product, page, pageSize int) []Product {
	requestedProducts := []Product{}

	if page <= 0 || pageSize <= 0 {
		return requestedProducts
	}

	startIndex := (page - 1) * pageSize
	endIndex := startIndex + pageSize

	if startIndex >= len(products) {
		return requestedProducts
	}

	if endIndex >= len(products) {
		endIndex = len(products)
	}

	for i := startIndex; i < endIndex; i++ {
		requestedProducts = append(requestedProducts, products[i])
	}

	return requestedProducts
}

func CountProductPages(products []Product, pageSize int) int {
	if len(products) == 0 || pageSize <= 0 {
		return 0
	}

	pages := len(products) / pageSize

	if len(products)%pageSize != 0 {
		pages++
	}

	return pages
}

func GetCatalogPage(
	products []Product,
	keyword string,
	minPrice, maxPrice float64,
	availableOnly bool,
	ascending bool,
	page, pageSize int,
) []Product {
	if minPrice > maxPrice {
		return []Product{}
	}

	filteredProducts := products

	if len(keyword) != 0 {
		filteredProducts = SearchProductsByName(filteredProducts, keyword)
	}

	filteredByPrice := FilterProductsByPrice(filteredProducts, minPrice, maxPrice)

	filteredByAvailability := filteredByPrice
	if availableOnly {
		filteredByAvailability = FilterAvailableProducts(filteredByPrice)
	}

	sortedByPrice := SortProductsByPrice(filteredByAvailability, ascending)
	sortedToPages := GetProductsPage(sortedByPrice, page, pageSize)

	return sortedToPages
}
