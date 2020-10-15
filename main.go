package main

import (
	"VoshodFetcher/db"
	"VoshodFetcher/libs"
	"log"
	"time"
)

type Manufacturer struct {
	SortOrder      int    `json:"sort_order"`
	ManufacturerId uint   `json:"manufacturer_id"`
	Name           string `json:"name"`
}

type Category struct {
	CategoryID   uint `gorm:"primary_key" json:"category_id"`
	ParentID     uint `json:"parent_id"`
	Status       int
	Top          int       `json:"top"`
	Column       int       `json:"column"`
	DateAdded    time.Time `json:"date_added"`
	DateModified time.Time `json:"date_modified"`
}

func (c Category) TableName() string {
	return "oc_category"
}

type CategoryDescription struct {
	CategoryID      uint   `json:"category_id"`
	LanguageId      int    `json:"language_id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	MetaTitle       string `json:"meta_title"`
	MetaDescription string `json:"meta_description"`
	MetaKeyword     string `json:"meta_keyword"`
}

type Product struct {
	ProductId      uint `gorm:"primary_key" json:"product_id"`
	Subtract       int  `json:"subtract"`
	Minimum        int  `json:"minimum"`
	Status         int
	StockStatusId  int     `json:"stock_status_id"`
	Quantity       int     `json:"quantity"`
	Price          float32 `json:"price"`
	Location       string  `json:"location"`
	ManufacturerId uint    `json:"manufacturer_id"`
	Model          string  `json:"model"`
	Sku            string
	Upc            string
	Ean            string
	Jan            string
	Isbn           string
	Mpn            string
	Shipping       int
	Points         int
}

func (p Product) TableName() string {
	return "oc_product"
}

func main() {

	conf, _ := libs.LoadDatabaseConfiguration()

	db.Connect(conf)

	defer db.Close()

	//resVen, _ := libs.FetchResult(libs.FetchTypeVendor, 0)
	//for _, v := range resVen.Response.Vendors {
	//	r := Manufacturer{Name: v.Name, SortOrder: 0}
	//	q := db.SQL().Table("oc_manufacturer").First(&r, "name = ?", v.Name)
	//	if q.RecordNotFound() {
	//		db.SQL().Table("oc_manufacturer").Save(r)
	//		resDB2 := db.SQL().Table("oc_manufacturer").First(&r, "name = ?", v.Name)
	//		log.Println(r.ManufacturerId)
	//		if resDB2.Error == nil && r.ManufacturerId > 0 {
	//			db.SQL().Table("oc_manufacturer_to_store").Omit("name").Save(&r)
	//		}
	//	}
	//}

	categories := make(map[string]uint)
	//resCat, _ := libs.FetchResult(libs.FetchTypeCatalogs, 0)
	//// Проверяем существует ли категория имя
	//for _, c := range resCat.Response.Catalogs {
	//	catDescr := CategoryDescription{Name: c.Name, LanguageId: 1, Description: "", MetaDescription: "", MetaTitle: "", MetaKeyword: ""}
	//	q := db.SQL().Table("oc_category_description").First(&catDescr, "name = ?", c.Name)
	//	if q.RecordNotFound() { // Не существует
	//		cat := Category{Status: 1, ParentID: uint(categories[c.ParentID]), Top: 1, Column: 1, DateAdded: time.Time{}, DateModified: time.Time{}}
	//		catSv := db.SQL().Create(&cat)
	//		if catSv.Error == nil {
	//			catDescr.CategoryID = cat.CategoryID
	//			db.SQL().Table("oc_category_description").Save(&catDescr)
	//
	//			c2s := struct {
	//				CategoryID uint `json:"category_id"`
	//				StoreId    int  `json:"store_id"`
	//			}{
	//				CategoryID: cat.CategoryID,
	//				StoreId:    0,
	//			}
	//			db.SQL().Table("oc_category_to_store").Save(&c2s)
	//		}
	//	}
	//	categories[c.ID] = catDescr.CategoryID
	//}

	resItems, _ := libs.FetchResult(libs.FetchTypeItems, 1)
	// Проверяем существует ли категория имя
	for _, p := range resItems.Response.Items {
		brand := uint(0)
		//cat := categories[p.CatalogID].CategoryID
		prod := Product{ManufacturerId: brand, Status: 1, Model: p.Name, Price: p.Price, Quantity: p.Count, StockStatusId: 5}
		q := db.SQL().Table("oc_product").First(&prod, "model = ? AND manufacturer_id = ?", p.Name, 0)
		if q.RecordNotFound() { // Не существует
			ProdSv := db.SQL().Create(&prod)
			log.Println(p.CatalogID)
			if ProdSv.Error == nil && prod.ProductId > 0 {
				if p.CatalogID != "" {
					cat := categories[p.CatalogID]
					if cat != 0 {
						p2c := struct {
							CategoryID uint `json:"category_id"`
							ProductId  uint `json:"product_id"`
						}{
							CategoryID: cat,
							ProductId:  prod.ProductId,
						}
						db.SQL().Table("oc_product_to_category").Create(&p2c)
					}

					p2s := struct {
						ProductId uint `json:"product_id"`
						StoreId   int  `json:"store_id"`
					}{
						ProductId: prod.ProductId,
						StoreId:   0,
					}
					db.SQL().Table("oc_product_to_store").Create(&p2s)
				}

			}
		}
	}
	//for _, v := range res.Response.Vendors {
	//	db.SQL().Create(map[string]interface{}{
	//		"name": v.Name,
	//	})
	//}
	//res := libs.Scrape(libs.FetchVendor2,0)
	//b, err := json.MarshalIndent(res, "", "  ")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Print(string(b))
	//log.Printf("V: %d, C: %d, I: %d", len(res.Response.Vendors), len(res.Response.Catalogs), len(res.Response.Items))
}
