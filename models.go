package main

import "sort"

type Wish struct {
	ID       int
	Date     string
	Link     string
	Name     string
	Price    float64
	Currency string
	Category string
}

var Wishes = []Wish{
	{ID: 1, Date: "2024-03-10", Link: "https://aliexpress.ru/item/1005005134514393.html", Name: "TIMEMORE Basic Plus Scale", Price: 17.99, Currency: "₽", Category: "Electronics"},
	{ID: 2, Date: "2024-04-22", Link: "https://aliexpress.ru/item/1005005751288141.html?sku_id=12000034222521117", Name: "ThinkRider Professional NL-15 PRO Bicycle Bike Torque Wrench", Price: 3641, Currency: "₽", Category: "Electronics"},
	{ID: 3, Date: "2024-06-06", Link: "https://aliexpress.ru/item/1000005654040.html", Name: "Household Digital Ultrasonic Cleaner", Price: 4432, Currency: "₽", Category: "Electronics"},
	{ID: 4, Date: "2024-02-21", Link: "https://www.amazon.com/Ninja-Professional-Countertop-Technology-BL610/dp/B00NGV4506", Name: "Ninja BL610 Professional Blender", Price: 73.05, Currency: "$", Category: "Electronics"},
	{ID: 5, Date: "2024-06-06", Link: "https://aliexpress.ru/item/1005004280875987.html", Name: "N2 Nylon Carbon Bottle Cage, 2pcx", Price: 852, Currency: "₽", Category: "Cycling Gear"},
	{ID: 6, Date: "2024-06-07", Link: "https://www.tradeinn.com/bikeinn/en/zipp-hanlebar-smooth-course-handlebar-tape/136088660/p", Name: "Zipp Hanlebar Smooth Course handlebar tape", Price: 1490, Currency: "€", Category: "Cycling Gear"},
	{ID: 7, Date: "2024-06-07", Link: "https://veter.cc/store/tproduct/494955126-923722664721-noski-velosipednie-classic-white", Name: "Носки Велосипедные Classic White, Size: 39-42", Price: 1490, Currency: "₽", Category: "Apparel"},
	{ID: 8, Date: "2024-02-23", Link: "https://www.rapha.cc/it/en/shop/mens-pro-team-bib-shorts-regular/product/BEP01XXBLW", Name: "Rapha MEN'S PRO TEAM BIB SHORTS - REGULAR, Size: S", Price: 260, Currency: "€", Category: "Apparel"},
	{ID: 9, Date: "2024-06-07", Link: "https://www.tradeinn.com/bikeinn/en/muc-off-bio-canister-140ml-tubeless-sealant/137180560/p", Name: "Muc Off Bio Canister 140ml Tubeless Sealant", Price: 801.99, Currency: "₽", Category: "Accessories"},
	{ID: 10, Date: "2024-06-07", Link: "https://www.tradeinn.com/bikeinn/en/muc-off-tubeless-rim-tape-10-meters/137682600/p", Name: "Muc Off Tubeless Rim Tape 10 Meters, 21 mm", Price: 1490, Currency: "₽", Category: "Accessories"},
	{ID: 9, Date: "2024-06-08", Link: "https://www.tradeinn.com/bikeinn/en/muc-off-bio-canister-140ml-tubeless-sealant/137180560/p", Name: "Muc Off Bio Canister 140ml Tubeless Sealant", Price: 817.99, Currency: "₽", Category: "Accessories"},
}

func groupWishesByCategory(wishes []Wish) map[string][]Wish {
	categories := make(map[string][]Wish)
	for _, wish := range wishes {
		categories[wish.Category] = append(categories[wish.Category], wish)
	}
	return categories
}

func sortCategories(categories map[string][]Wish) map[string][]Wish {
	for _, wishes := range categories {
		sort.Slice(wishes, func(i, j int) bool {
			return wishes[i].Price < wishes[j].Price
		})
	}
	return categories
}

func GetCategories() map[string][]Wish {
	categories := groupWishesByCategory(Wishes)
	categories = sortCategories(categories)
	return categories
}
