package product_test

import (
	"reflect"
	"testing"

	"example.com/online-shop/product"
)

var products = []product.Product{
	{Name: "Laptop", Price: 800.00, Stock: 10},
	{Name: "Gaming Laptop", Price: 1200.00, Stock: 5},
	{Name: "Mouse", Price: 30.00, Stock: 50},
	{Name: "Keyboard", Price: 75.00, Stock: 25},
	{Name: "Monitor", Price: 250.00, Stock: 15},
	{Name: "Headphones", Price: 100.00, Stock: 0},
	{Name: "Webcam", Price: 60.00, Stock: 12},
	{Name: "Microphone", Price: 150.00, Stock: 8},
	{Name: "USB Cable", Price: 10.00, Stock: 100},
	{Name: "External SSD", Price: 180.00, Stock: 20},
	{Name: "Graphics Card", Price: 650.00, Stock: 4},
	{Name: "Motherboard", Price: 220.00, Stock: 7},
	{Name: "CPU", Price: 350.00, Stock: 9},
	{Name: "RAM 16GB", Price: 55.00, Stock: 30},
	{Name: "RAM 32GB", Price: 95.00, Stock: 18},
	{Name: "Power Supply", Price: 110.00, Stock: 6},
	{Name: "PC Case", Price: 85.00, Stock: 11},
	{Name: "CPU Cooler", Price: 45.00, Stock: 0},
	{Name: "Laptop Stand", Price: 40.00, Stock: 22},
	{Name: "Mouse Pad", Price: 15.00, Stock: 60},
	{Name: "Wi-Fi Adapter", Price: 25.00, Stock: 14},
	{Name: "Router", Price: 130.00, Stock: 5},
	{Name: "Speakers", Price: 90.00, Stock: 16},
	{Name: "Flash Drive 64GB", Price: 20.00, Stock: 40},
	{Name: "Docking Station", Price: 190.00, Stock: 3},
}

var testProducts = []product.Product{
	{Name: "Laptop Pro", Price: 900, Stock: 0},
}

func TestGetCatalogPage(t *testing.T) {
	tests := []struct {
		name          string
		keyword       string
		input         []product.Product
		minPrice      float64
		maxPrice      float64
		expected      []product.Product
		availableOnly bool
		ascending     bool
		page          int
		pageSize      int
	}{
		{
			name:     "Search for laptops",
			keyword:  "laptop",
			input:    products,
			minPrice: 500,
			maxPrice: 1500,
			expected: []product.Product{
				{Name: "Laptop", Price: 800.00, Stock: 10},
				{Name: "Gaming Laptop", Price: 1200.00, Stock: 5},
			},
			availableOnly: true,
			ascending:     true,
			page:          1,
			pageSize:      10,
		},

		{
			name:          "Search for nonexistent product",
			keyword:       "nonexistent",
			input:         products,
			minPrice:      0,
			maxPrice:      1500,
			expected:      []product.Product{},
			availableOnly: true,
			ascending:     true,
			page:          1,
			pageSize:      10,
		},

		{
			name:          "Invalid price range",
			keyword:       "laptop",
			input:         products,
			minPrice:      1500,
			maxPrice:      500,
			expected:      []product.Product{},
			availableOnly: true,
			ascending:     true,
			page:          1,
			pageSize:      10,
		},

		{
			name:     "Empty search keyword",
			keyword:  "",
			input:    products,
			minPrice: 500,
			maxPrice: 1500,
			expected: []product.Product{
				{Name: "Graphics Card", Price: 650.00, Stock: 4},
				{Name: "Laptop", Price: 800.00, Stock: 10},
				{Name: "Gaming Laptop", Price: 1200.00, Stock: 5},
			},
			availableOnly: true,
			ascending:     true,
			page:          1,
			pageSize:      100,
		},

		{
			name:     "Out of stock products",
			keyword:  "Laptop",
			input:    testProducts,
			minPrice: 500,
			maxPrice: 1500,
			expected: []product.Product{
				{Name: "Laptop Pro", Price: 900, Stock: 0},
			},
			availableOnly: false,
			ascending:     true,
			page:          1,
			pageSize:      100,
		},

		{
			name:     "Laptops sorted by descending price",
			keyword:  "laptop",
			input:    products,
			minPrice: 500,
			maxPrice: 1500,
			expected: []product.Product{
				{Name: "Gaming Laptop", Price: 1200.00, Stock: 5},
				{Name: "Laptop", Price: 800.00, Stock: 10},
			},
			availableOnly: true,
			ascending:     false,
			page:          1,
			pageSize:      10,
		},

		{
			name:     "Pagination check for the second page",
			keyword:  "",
			input:    products,
			minPrice: 0,
			maxPrice: 200,
			expected: []product.Product{
				{Name: "Wi-Fi Adapter", Price: 25.00, Stock: 14},
				{Name: "Mouse", Price: 30.00, Stock: 50},
				{Name: "Laptop Stand", Price: 40.00, Stock: 22},
			},
			availableOnly: true,
			ascending:     true,
			page:          2,
			pageSize:      3,
		},

		{
			name:          "Out of range pages testing",
			keyword:       "",
			input:         products,
			minPrice:      0,
			maxPrice:      200,
			expected:      []product.Product{},
			availableOnly: true,
			ascending:     true,
			page:          100,
			pageSize:      3,
		},

		{
			name:          "Out of range pages testing",
			keyword:       "",
			input:         products,
			minPrice:      0,
			maxPrice:      200,
			expected:      []product.Product{},
			availableOnly: true,
			ascending:     true,
			page:          0,
			pageSize:      3,
		},

		{
			name:          "Out of range pages testing",
			keyword:       "",
			input:         products,
			minPrice:      0,
			maxPrice:      200,
			expected:      []product.Product{},
			availableOnly: true,
			ascending:     true,
			page:          -1,
			pageSize:      3,
		},

		{
			name:          "Out of range pages testing",
			keyword:       "",
			input:         products,
			minPrice:      0,
			maxPrice:      200,
			expected:      []product.Product{},
			availableOnly: true,
			ascending:     true,
			page:          1,
			pageSize:      0,
		},

		{
			name:          "Out of range pages testing",
			keyword:       "",
			input:         products,
			minPrice:      0,
			maxPrice:      200,
			expected:      []product.Product{},
			availableOnly: true,
			ascending:     true,
			page:          1,
			pageSize:      -3,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := product.GetCatalogPage(
				test.input,
				test.keyword,
				test.minPrice,
				test.maxPrice,
				test.availableOnly,
				test.ascending,
				test.page,
				test.pageSize,
			)

			if !reflect.DeepEqual(result, test.expected) {
				t.Errorf("Expected %v, gave %v", test.expected, result)
			}
		})
	}
}
