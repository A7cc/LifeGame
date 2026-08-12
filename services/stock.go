package services

import (
	"LifeGame/core"
	"time"
)

const stockUpdateMinInterval = 4 * time.Second

func cloneStocks(stocks []core.StockInfo) []core.StockInfo {
	result := make([]core.StockInfo, len(stocks))
	for i, stock := range stocks {
		result[i] = stock
		result[i].SIHistory = append([]int(nil), stock.SIHistory...)
		result[i].SIKlineHistory = append([][4]int(nil), stock.SIKlineHistory...)
	}
	return result
}

// 更新股票数据
func (a *App) UpdateStock() StockUpdateResponse {
	a.stateMu.Lock()
	defer a.stateMu.Unlock()

	if errResp := a.requireGame(); errResp != nil {
		return StockUpdateResponse{Code: -1, Msg: responseMessage(errResp)}
	}
	now := time.Now()
	if a.Gameinfo.GStockUpdateCount >= core.MaxStockUpdatesPerYear {
		return StockUpdateResponse{
			Code: 200, Msg: "本年度行情已收盘，请进入下一年",
			Stocks: cloneStocks(a.Gameinfo.GStockInfo), News: append([]string(nil), a.Gameinfo.GStockNews...),
			Epoch: a.stockEpoch, Version: a.stockVersion, UpdatedAt: a.stockUpdatedAt,
			MarketClosed: true,
		}
	}
	if a.stockUpdatedAt > 0 && now.Sub(time.UnixMilli(a.stockUpdatedAt)) < stockUpdateMinInterval {
		return StockUpdateResponse{
			Code: 200, Msg: "行情更新过于频繁",
			Stocks: cloneStocks(a.Gameinfo.GStockInfo), News: append([]string(nil), a.Gameinfo.GStockNews...),
			Epoch: a.stockEpoch, Version: a.stockVersion, UpdatedAt: a.stockUpdatedAt,
			Remaining: core.MaxStockUpdatesPerYear - a.Gameinfo.GStockUpdateCount,
		}
	}
	// 更新股票价格信息
	a.Gameinfo.UpdateStockPrices()
	a.Gameinfo.GStockUpdateCount++
	a.stockVersion++
	a.stockUpdatedAt = now.UnixMilli()
	return StockUpdateResponse{
		Code: 200, Stocks: cloneStocks(a.Gameinfo.GStockInfo),
		News: append([]string(nil), a.Gameinfo.GStockNews...), Epoch: a.stockEpoch,
		Version: a.stockVersion, UpdatedAt: a.stockUpdatedAt,
		Remaining:    core.MaxStockUpdatesPerYear - a.Gameinfo.GStockUpdateCount,
		MarketClosed: a.Gameinfo.GStockUpdateCount >= core.MaxStockUpdatesPerYear,
	}
}
