package chia

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func CreateSqlSchema() error {
	db, err := sql.Open("sqlite3", "./chia.prices.db")
	if err != nil {
		return fmt.Errorf("Error when opening database: %s", err)
	}
	defer db.Close()

	createMapTable := `CREATE TABLE IF NOT EXISTS crypto_map (
		id INTEGER primary key,
		coin_id INTEGER,
		coin_rank INTEGER,
		coin_name VARCHAR,
		symbol VARHAR,
		coin_ts DATETIME,
		UNIQUE(id, coin_id) );
		`
	log.Println("Creating table: crypto_map")
	_, err = db.Exec(createMapTable)
	if err != nil {
		return fmt.Errorf("%q: %s\n", err, createMapTable)
	}
	log.Println("Creating table: crypto_map [complete]")

	createChiaPriceTable := `CREATE TABLE IF NOT EXISTS chia_price (
		id INTEGER primary key autoincrement,
		coin_id INTEGER,
		price_usd INTEGER,
		Volume24H INTEGER,
		PercentChange1H INTEGER,
		PercentChange24H INTEGER,
		MarketCap INTEGER,
		TotalSupply VARCHAR,
		MaxSupport VARCHAR,
		coin_ts VARCHAR );
		`
	log.Println("Creating table: chia_price")
	_, err = db.Exec(createChiaPriceTable)
	if err != nil {
		return fmt.Errorf("%q: %s\n", err, createChiaPriceTable)
	}
	log.Println("Creating table: chia_price [complete]")
	return nil
}

func TableChiaPricesInsert(data *CoinMarketCapSymbolResponse) error {
	db, err := sql.Open("sqlite3", "./chia.prices.db")
	if err != nil {
		return fmt.Errorf("Error when opening database: %s", err)
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("insert or ignore into chia_price(coin_id, price_usd, Volume24h, PercentChange1H, PercentChange24H, MarketCap, TotalSupply, MaxSupport, coin_ts) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?) ")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(data.Data.Num9258.ID, data.Data.Num9258.Quote.Usd.Price, data.Data.Num9258.Quote.Usd.Volume24H, data.Data.Num9258.Quote.Usd.PercentChange1H, data.Data.Num9258.Quote.Usd.PercentChange24H,
		data.Data.Num9258.Quote.Usd.MarketCap, data.Data.Num9258.TotalSupply, data.Data.Num9258.MaxSupply, data.Status.Timestamp)
	if err != nil {
		log.Fatal(err)
	}
	tx.Commit()
	return nil
}
func TableChiaMapInsert(id int, coin_id int, coin_rank int, coin_name string, symbol string, coin_ts string) error {
	db, err := sql.Open("sqlite3", "./chia.prices.db")
	if err != nil {
		return fmt.Errorf("Error when opening database: %s", err)
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("insert or ignore into crypto_map(id,coin_id, coin_rank, coin_name, symbol, coin_ts) VALUES(?,?, ?, ?, ?, ?) ")
	if err != nil {
		log.Fatal(err)
	}

	defer stmt.Close()
	_, err = stmt.Exec(id, coin_id, coin_rank, coin_name, symbol, coin_ts)
	if err != nil {
		log.Fatal(err)
	}
	tx.Commit()
	return nil
}
func FetchCoinData(coin_name string) (int, error) {
	db, err := sql.Open("sqlite3", "./chia.prices.db")
	if err != nil {
		return 0, fmt.Errorf("Error when opening database: %s", err)
	}
	defer db.Close()
	rows, err := db.Query("select coin_id from crypto_map where symbol = ?", coin_name)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var result int
		if err := rows.Scan(&result); err != nil {
			// Check for a scan error.
			// Query rows will be closed with defer.
			log.Fatal(err)
		}
		return result, nil
	}

	return 0, nil
}

/**

		coin_id INTEGER,
		price_usd INTEGER,
		Volume24H INTEGER,
		PercentChange1H INTEGER,
		PercentChange24H INTEGER,
		MarketCap INTEGER,
		TotalSupply VARCHAR,
		MaxSupport VARCHAR,
		coin_ts VARCHAR );
**/

func FetchChiaPriceDB() *ChiaTableDbResponse {
	db, err := sql.Open("sqlite3", "./chia.prices.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT id,price_usd,PercentChange1H,PercentChange24H,TotalSupply from chia_price ORDER BY ID DESC LIMIT 1")
	chiaDBResponse := &ChiaTableDbResponse{}
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		var price_usd float64
		var percent_change_1h float64
		var percent_change_24h float64
		var total_supply float64
		if err := rows.Scan(&id, &price_usd, &percent_change_1h, &percent_change_24h, &total_supply); err != nil {
			// Check for a scan error.
			// Query rows will be closed with defer.
			log.Fatal(err)
		}
		chiaDBResponse.UpdateId = id
		chiaDBResponse.ChiaPrice = price_usd
		chiaDBResponse.PercentChange1H = percent_change_1h
		chiaDBResponse.PercentChange24h = percent_change_24h
		chiaDBResponse.TotalSupply = total_supply
	}

	return chiaDBResponse
}
