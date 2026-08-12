package db

import (
	"database/sql"
	"encoding/json"
	"path/filepath"
	"testing"
)

func TestInitDBCreatesCurrentSchemaAndReferenceData(t *testing.T) {
	database, err := InitDB(filepath.Join(t.TempDir(), "lifegame.db"))
	if err != nil {
		t.Fatalf("InitDB() error = %v", err)
	}
	t.Cleanup(func() { CloseDB() })

	var migrationTables int
	if err := database.QueryRow(`
		SELECT COUNT(*) FROM sqlite_master
		WHERE type = 'table' AND name = 'schema_migrations'
	`).Scan(&migrationTables); err != nil {
		t.Fatal(err)
	}
	if migrationTables != 0 {
		t.Fatal("current database unexpectedly contains schema_migrations")
	}

	for _, column := range []string{"gender", "nationality", "meet_scene"} {
		var count int
		if err := database.QueryRow(`SELECT COUNT(*) FROM pragma_table_info('datings') WHERE name = ?`, column).Scan(&count); err != nil {
			t.Fatal(err)
		}
		if count != 1 {
			t.Fatalf("datings.%s column count = %d, want 1", column, count)
		}
	}
	var saveVersionColumns int
	if err := database.QueryRow(`SELECT COUNT(*) FROM pragma_table_info('saves') WHERE name = 'save_version'`).Scan(&saveVersionColumns); err != nil {
		t.Fatal(err)
	}
	if saveVersionColumns != 1 {
		t.Fatalf("saves.save_version column count = %d, want 1", saveVersionColumns)
	}

	var miniGameCount, antiqueCount, houseCount, carCount int
	if err := database.QueryRow("SELECT COUNT(*) FROM minigames").Scan(&miniGameCount); err != nil {
		t.Fatal(err)
	}
	if err := database.QueryRow("SELECT COUNT(*) FROM houses").Scan(&houseCount); err != nil {
		t.Fatal(err)
	}
	if err := database.QueryRow("SELECT COUNT(*) FROM antiques").Scan(&antiqueCount); err != nil {
		t.Fatal(err)
	}
	if err := database.QueryRow("SELECT COUNT(*) FROM cars").Scan(&carCount); err != nil {
		t.Fatal(err)
	}
	if miniGameCount != 23 || antiqueCount != 50 || houseCount != 50 || carCount != 50 {
		t.Fatalf("minigames/antiques/houses/cars = %d/%d/%d/%d, want 23/50/50/50", miniGameCount, antiqueCount, houseCount, carCount)
	}

	var invalidAssetImages int
	if err := database.QueryRow(`
		SELECT (SELECT COUNT(*) FROM houses WHERE img NOT LIKE '/images/houseinfo/houses/%.webp')
		     + (SELECT COUNT(*) FROM cars WHERE img NOT LIKE '/images/carinfo/cars/%.webp')
	`).Scan(&invalidAssetImages); err != nil {
		t.Fatal(err)
	}
	if invalidAssetImages != 0 {
		t.Fatalf("non-current house/car images = %d, want 0", invalidAssetImages)
	}

	var invalidAntiqueImages, distinctAntiqueImages int
	if err := database.QueryRow(`
		SELECT COUNT(*) FROM antiques
		WHERE img NOT GLOB '/images/antiqueinfo/[0-9][0-9]-*.webp'
	`).Scan(&invalidAntiqueImages); err != nil {
		t.Fatal(err)
	}
	if err := database.QueryRow("SELECT COUNT(DISTINCT img) FROM antiques").Scan(&distinctAntiqueImages); err != nil {
		t.Fatal(err)
	}
	if invalidAntiqueImages != 0 || distinctAntiqueImages != 50 {
		t.Fatalf("invalid/distinct antique images = %d/%d, want 0/50", invalidAntiqueImages, distinctAntiqueImages)
	}

	var taxiDifficultyJSON string
	if err := database.QueryRow("SELECT difficulty FROM minigames WHERE name = 'taxi'").Scan(&taxiDifficultyJSON); err != nil {
		t.Fatal(err)
	}
	var taxiDifficulties map[int]struct {
		MinRunTime int `json:"smgminruntime"`
	}
	if err := json.Unmarshal([]byte(taxiDifficultyJSON), &taxiDifficulties); err != nil {
		t.Fatal(err)
	}
	if taxiDifficulties[0].MinRunTime != 1 {
		t.Fatalf("taxi minimum runtime = %d, want 1", taxiDifficulties[0].MinRunTime)
	}

	var femaleDatings, maleDatings int
	if err := database.QueryRow("SELECT COUNT(*) FROM datings WHERE gender = 0").Scan(&femaleDatings); err != nil {
		t.Fatal(err)
	}
	if err := database.QueryRow("SELECT COUNT(*) FROM datings WHERE gender = 1").Scan(&maleDatings); err != nil {
		t.Fatal(err)
	}
	if femaleDatings != 50 || maleDatings != 50 {
		t.Fatalf("dating candidates female/male = %d/%d, want 50/50", femaleDatings, maleDatings)
	}

	for gender, wantImage := range map[int]string{
		0: "/images/datinginfo/dating-partner/female/01.webp",
		1: "/images/datinginfo/dating-partner/male/01.webp",
	} {
		var image string
		if err := database.QueryRow("SELECT image FROM datings WHERE gender = ? ORDER BY id LIMIT 1", gender).Scan(&image); err != nil {
			t.Fatal(err)
		}
		if image != wantImage {
			t.Fatalf("gender %d first image = %q, want %q", gender, image, wantImage)
		}
		assertDatingRegionDistribution(t, database, gender)
	}

	var sceneCandidates int
	if err := database.QueryRow("SELECT COUNT(*) FROM datings WHERE meet_scene <> ''").Scan(&sceneCandidates); err != nil {
		t.Fatal(err)
	}
	if sceneCandidates != 14 {
		t.Fatalf("scene candidates = %d, want 14", sceneCandidates)
	}
	assertCurrentDatingGameConditions(t, database)

	if err := initializeCurrentDatabase(database, true); err != nil {
		t.Fatalf("second initializeCurrentDatabase() error = %v", err)
	}
	var datingCount int
	if err := database.QueryRow("SELECT COUNT(*) FROM datings").Scan(&datingCount); err != nil {
		t.Fatal(err)
	}
	if datingCount != 100 {
		t.Fatalf("dating count after second initialization = %d, want 100", datingCount)
	}
}

func TestInitDBPreservesExistingUserData(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "lifegame.db")
	database, err := InitDB(dbPath)
	if err != nil {
		t.Fatalf("first InitDB() error = %v", err)
	}
	if _, err := database.Exec("UPDATE houses SET name = ? WHERE id = 1", "自定义房产"); err != nil {
		t.Fatal(err)
	}
	if err := CloseDB(); err != nil {
		t.Fatal(err)
	}

	database, err = InitDB(dbPath)
	if err != nil {
		t.Fatalf("second InitDB() error = %v", err)
	}
	t.Cleanup(func() { CloseDB() })

	var name string
	if err := database.QueryRow("SELECT name FROM houses WHERE id = 1").Scan(&name); err != nil {
		t.Fatal(err)
	}
	if name != "自定义房产" {
		t.Fatalf("existing database value = %q, want user value", name)
	}
}

func assertCurrentDatingGameConditions(t *testing.T, database *sql.DB) {
	t.Helper()
	wantTargets := map[int][]string{
		5:  {""},
		6:  {"", ""},
		13: {""},
		14: {""},
		26: {"keepfit"},
		38: {"party"},
		39: {""},
		41: {"concert"},
		46: {"concert"},
	}
	for baseID, want := range wantTargets {
		for _, datingID := range []int{baseID, baseID + 50} {
			var conditionsJSON string
			if err := database.QueryRow("SELECT meet_conditions FROM datings WHERE id = ?", datingID).Scan(&conditionsJSON); err != nil {
				t.Fatal(err)
			}
			var conditions []struct {
				Type   string `json:"ctype"`
				Target string `json:"ctarget"`
			}
			if err := json.Unmarshal([]byte(conditionsJSON), &conditions); err != nil {
				t.Fatal(err)
			}
			got := make([]string, 0, len(want))
			for _, condition := range conditions {
				if condition.Type == "play_game" || condition.Type == "win_game" {
					got = append(got, condition.Target)
				}
			}
			if len(got) != len(want) {
				t.Fatalf("dating %d game targets = %#v, want %#v", datingID, got, want)
			}
			for index := range want {
				if got[index] != want[index] {
					t.Fatalf("dating %d game targets = %#v, want %#v", datingID, got, want)
				}
			}
		}
	}
}

func assertDatingRegionDistribution(t *testing.T, database *sql.DB, gender int) {
	t.Helper()
	regionsByNationality := map[string]string{
		"中国": "中国",
		"日本": "亚洲其他地区", "韩国": "亚洲其他地区", "印度": "亚洲其他地区", "菲律宾": "亚洲其他地区", "泰国": "亚洲其他地区",
		"法国": "欧洲", "意大利": "欧洲", "西班牙": "欧洲", "英国": "欧洲", "荷兰": "欧洲", "奥地利": "欧洲", "俄罗斯": "欧洲", "乌克兰": "欧洲", "瑞典": "欧洲", "波兰": "欧洲",
		"黎巴嫩": "中东", "埃及": "中东", "沙特阿拉伯": "中东", "以色列": "中东", "伊拉克": "中东", "阿联酋": "中东", "约旦": "中东", "土耳其": "中东", "伊朗": "中东", "卡塔尔": "中东",
		"巴西": "美洲", "美国": "美洲", "加拿大": "美洲", "阿根廷": "美洲", "墨西哥": "美洲", "智利": "美洲", "哥伦比亚": "美洲",
		"新西兰": "大洋洲", "澳大利亚": "大洋洲", "斐济": "大洋洲",
		"尼日利亚": "非洲", "南非": "非洲", "摩洛哥": "非洲",
	}
	want := map[string]int{
		"中国": 8, "亚洲其他地区": 5, "欧洲": 10, "中东": 10,
		"美洲": 10, "大洋洲": 4, "非洲": 3,
	}
	got := map[string]int{}
	distinct := map[string]struct{}{}
	rows, err := database.Query("SELECT nationality FROM datings WHERE gender = ?", gender)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var nationality string
		if err := rows.Scan(&nationality); err != nil {
			t.Fatal(err)
		}
		region, ok := regionsByNationality[nationality]
		if !ok {
			t.Fatalf("gender %d has unclassified nationality %q", gender, nationality)
		}
		got[region]++
		distinct[nationality] = struct{}{}
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	for region, wantCount := range want {
		if got[region] != wantCount {
			t.Fatalf("gender %d region %s count = %d, want %d", gender, region, got[region], wantCount)
		}
	}
	if len(distinct) != 39 {
		t.Fatalf("gender %d distinct nationalities = %d, want 39", gender, len(distinct))
	}
}
