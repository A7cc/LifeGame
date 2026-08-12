package db

import (
	"database/sql"
	"encoding/json"

	"LifeGame/core"
)

// LoadItems 从数据库加载物资数据
func LoadItems(db *sql.DB) ([]core.Item, []core.Item, error) {
	rows, err := db.Query("SELECT id, name, price, price_min, price_max, effects, region FROM items")
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	domesticItems := []core.Item{}
	foreignItems := []core.Item{}

	for rows.Next() {
		item := core.Item{}
		effectsJSON := ""
		region := ""

		err := rows.Scan(&item.IId, &item.IName, &item.IPrice, &item.IPrice_min, &item.IPrice_max, &effectsJSON, &region)
		if err != nil {
			return nil, nil, err
		}

		// 解析 effects JSON
		if err := json.Unmarshal([]byte(effectsJSON), &item.IEffects); err != nil {
			return nil, nil, err
		}

		if region == "domestic" {
			domesticItems = append(domesticItems, item)
		} else {
			foreignItems = append(foreignItems, item)
		}
	}

	return domesticItems, foreignItems, nil
}

// LoadCompanies 从数据库加载公司数据
func LoadCompanies(db *sql.DB) ([]core.Company, error) {
	rows, err := db.Query("SELECT id, name, price, risk, profit, time FROM companies")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	companies := []core.Company{}

	for rows.Next() {
		c := core.Company{}

		if err := rows.Scan(&c.CId, &c.CName, &c.CPrice, &c.CRisk, &c.CProfit, &c.CTime); err != nil {
			return nil, err
		}
		companies = append(companies, c)
	}

	return companies, nil
}

// LoadAntiques 从数据库加载古董数据
func LoadAntiques(db *sql.DB) ([]core.Antique, error) {
	rows, err := db.Query("SELECT id, name, price, material, img, desc, level FROM antiques")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	antiques := []core.Antique{}

	for rows.Next() {
		a := core.Antique{}

		if err := rows.Scan(&a.AId, &a.AName, &a.APrice, &a.AMaterial, &a.AImg, &a.ADesc, &a.ALevel); err != nil {
			return nil, err
		}
		antiques = append(antiques, a)
	}

	return antiques, nil
}

// LoadStocks 从数据库加载股票数据
func LoadStocks(db *sql.DB) ([]core.Stock, error) {
	rows, err := db.Query("SELECT id, name, price, risk FROM stocks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stocks := []core.Stock{}

	for rows.Next() {
		s := core.Stock{}

		if err := rows.Scan(&s.SId, &s.SName, &s.SPrice, &s.SRisk); err != nil {
			return nil, err
		}
		stocks = append(stocks, s)
	}

	return stocks, nil
}

// LoadStockNews 从数据库加载股票新闻数据
func LoadStockNews(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SELECT content FROM stock_news ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	news := []string{}

	for rows.Next() {
		content := ""

		if err := rows.Scan(&content); err != nil {
			return nil, err
		}
		news = append(news, content)
	}

	return news, nil
}

// LoadHouses 从数据库加载房产数据
func LoadHouses(db *sql.DB) ([]core.House, error) {
	rows, err := db.Query("SELECT id, name, price, price_max, price_min, health, fame, img FROM houses")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	houses := []core.House{}

	for rows.Next() {
		h := core.House{}

		if err := rows.Scan(&h.HId, &h.HName, &h.HPrice, &h.HPrice_max, &h.HPrice_min, &h.HHealth, &h.HFame, &h.HImg); err != nil {
			return nil, err
		}
		houses = append(houses, h)
	}

	return houses, nil
}

// LoadCars 从数据库加载车辆数据
func LoadCars(db *sql.DB) ([]core.Car, error) {
	rows, err := db.Query("SELECT id, name, price, price_max, price_min, health, fame, img FROM cars")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cars := []core.Car{}

	for rows.Next() {
		c := core.Car{}

		if err := rows.Scan(&c.CId, &c.CName, &c.CPrice, &c.CPrice_max, &c.CPrice_min, &c.CHealth, &c.CFame, &c.CImg); err != nil {
			return nil, err
		}
		cars = append(cars, c)
	}

	return cars, nil
}

// LoadDatings 从数据库加载约会对象数据
func LoadDatings(db *sql.DB) ([]core.DatingInfo, error) {
	rows, err := db.Query(`SELECT id, name, image, age, gender, nationality, occupation, description, cost, meet_conditions, gifts, locations, meet_scene FROM datings ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	datings := []core.DatingInfo{}

	for rows.Next() {
		d := core.DatingInfo{}
		var meetConditionsJSON, giftsJSON, locationsJSON string

		if err := rows.Scan(&d.DId, &d.DName, &d.DImage, &d.DAge, &d.DSex, &d.DNationality, &d.DOccup, &d.DDesc, &d.DCost, &meetConditionsJSON, &giftsJSON, &locationsJSON, &d.DMeetScene); err != nil {
			return nil, err
		}

		// 解析 JSON 字段
		if err := json.Unmarshal([]byte(meetConditionsJSON), &d.DMeetConditions); err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(giftsJSON), &d.DGifts); err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(locationsJSON), &d.DLocations); err != nil {
			return nil, err
		}
		datings = append(datings, d)
	}
	return datings, nil
}

// LoadBankTasks 从数据库加载银行任务数据
func LoadBankTasks(db *sql.DB) ([]core.BankTask, error) {
	rows, err := db.Query("SELECT id, name, description, type, target_value, reward FROM bank_tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []core.BankTask{}

	for rows.Next() {
		t := core.BankTask{}
		err := rows.Scan(&t.TaskId, &t.TaskName, &t.TaskDesc, &t.TaskType, &t.TargetValue, &t.Reward)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}

// LoadMiniGames 从数据库加载小游戏数据
func LoadMiniGames(db *sql.DB) ([]core.MiniGame, error) {
	rows, err := db.Query(`SELECT id, name, cnname, icon, description, type, difficulty FROM minigames`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	games := []core.MiniGame{}
	for rows.Next() {
		g := core.MiniGame{}
		difficultyJSON := ""
		if err := rows.Scan(&g.MGId, &g.MGName, &g.MGCName, &g.MGIcon, &g.MGDesc, &g.MGType, &difficultyJSON); err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(difficultyJSON), &g.MGDifficulty); err != nil {
			return nil, err
		}
		games = append(games, g)
	}

	return games, nil
}

// LoadDiseases 从数据库加载疾病数据
func LoadDiseases(db *sql.DB) ([]core.DiseaseInfo, error) {
	rows, err := db.Query(`SELECT id, name, type, symptoms, healthimpact, upgradedays, treatment FROM diseases`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	diseases := []core.DiseaseInfo{}

	for rows.Next() {
		d := core.DiseaseInfo{}
		var treatmentJSON string

		err := rows.Scan(&d.DId, &d.DName, &d.DType, &d.DSymptoms, &d.DHealthImpact, &d.DUpgradeDays, &treatmentJSON)
		if err != nil {
			return nil, err
		}
		// 解析 JSON 字段
		if err := json.Unmarshal([]byte(treatmentJSON), &d.DTreatments); err != nil {
			return nil, err
		}
		diseases = append(diseases, d)
	}

	return diseases, nil
}

// LoadTreats 从数据库加载药品数据
func LoadTreats(db *sql.DB) ([]core.TreatInfo, error) {
	rows, err := db.Query(`SELECT id, name, description, price, heal, sideeffect, source FROM treats`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	treats := []core.TreatInfo{}

	for rows.Next() {
		m := core.TreatInfo{}

		if err := rows.Scan(&m.TId, &m.TName, &m.TDesc, &m.TPrice, &m.THeal, &m.TSideEffect, &m.TSource); err != nil {
			return nil, err
		}
		treats = append(treats, m)
	}

	return treats, nil
}

// LoadHospitals 从数据库加载医院数据
func LoadHospitals(db *sql.DB) ([]core.HospitalInfo, error) {
	rows, err := db.Query(`SELECT id, name, type, icon, description, services FROM hospitals`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	hospitals := []core.HospitalInfo{}
	for rows.Next() {
		h := core.HospitalInfo{}
		var servicesJSON string

		if err := rows.Scan(&h.HId, &h.HName, &h.HType, &h.HIcon, &h.HDescription, &servicesJSON); err != nil {
			return nil, err
		}
		// 解析 JSON 字段
		if err := json.Unmarshal([]byte(servicesJSON), &h.HServices); err != nil {
			return nil, err
		}
		hospitals = append(hospitals, h)
	}
	return hospitals, nil
}

// LoadAllGameData 加载所有游戏数据
func LoadAllGameData(db *sql.DB) error {
	var err error

	// 加载物资
	core.CachedItemIns, core.CachedItemOut, err = LoadItems(db)
	if err != nil {
		return err
	}

	// 加载公司
	core.CachedCompanies, err = LoadCompanies(db)
	if err != nil {
		return err
	}

	// 加载古董
	core.CachedAntiques, err = LoadAntiques(db)
	if err != nil {
		return err
	}

	// 加载股票
	core.CachedStocks, err = LoadStocks(db)
	if err != nil {
		return err
	}

	// 加载股票新闻
	core.CachedStockNews, err = LoadStockNews(db)
	if err != nil {
		return err
	}

	// 加载房产
	core.CachedHouses, err = LoadHouses(db)
	if err != nil {
		return err
	}

	// 加载车辆
	core.CachedCars, err = LoadCars(db)
	if err != nil {
		return err
	}

	// 加载约会对象信息
	core.CachedDatings, err = LoadDatings(db)
	if err != nil {
		return err
	}

	// 加载银行任务
	core.CachedBankTasks, err = LoadBankTasks(db)
	if err != nil {
		return err
	}

	// 加载小游戏
	core.CachedMiniGames, err = LoadMiniGames(db)
	if err != nil {
		return err
	}

	// 加载疾病信息
	core.CachedDiseases, err = LoadDiseases(db)
	if err != nil {
		return err
	}

	// 加载药品信息
	core.CachedTreats, err = LoadTreats(db)
	if err != nil {
		return err
	}

	// 加载医院信息
	core.CachedHospitals, err = LoadHospitals(db)
	if err != nil {
		return err
	}

	return nil
}
