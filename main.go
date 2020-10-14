package main

import (
	"VoshodFetcher/db"
	"VoshodFetcher/libs"
	"log"
)

type Manufacturer struct {
	ManufacturerId int `json:"manufacturer_id"`
	Name           string
}

func main() {

	conf, _ := libs.LoadDatabaseConfiguration()

	db.Connect(conf)

	defer db.Close()

	//vendors := map[string]Manufacturer{}
	res, _ := libs.FetchResult(libs.FetchTypeVendor, 1)
	for _, v := range res.Response.Vendors {
		r := Manufacturer{Name: v.Name}
		q := db.SQL().Table("oc_manufacturer").First(&r, "name = ?", v.Name)
		if q.RecordNotFound() {
			db.SQL().Table("oc_manufacturer").Save(r)
			resDB2 := db.SQL().Table("oc_manufacturer").First(&r, "name = ?", v.Name)
			log.Println(r.ManufacturerId)
			if resDB2.Error == nil && r.ManufacturerId > 0 {
				//r2 := struct {
				//	ManufacturerId int `json:"manufacturer_id"`
				//}{ManufacturerId: r.ManufacturerId}
				db.SQL().Table("oc_manufacturer_to_store").Omit("name").Save(&r)
			}
		}
	}

	//res, _ := libs.FetchResult(libs.FetchTypeVendor,0)

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
	log.Printf("V: %d, C: %d, I: %d", len(res.Response.Vendors), len(res.Response.Catalogs), len(res.Response.Items))
}
