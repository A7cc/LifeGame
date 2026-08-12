package services

import "LifeGame/core"

func (a *App) setAnnounce(announce core.Announce) {
	a.Announce = &announce
}

func (a *App) currentAnnounce() core.Announce {
	if a.Announce != nil {
		return *a.Announce
	}
	if a.Gameinfo != nil {
		return a.Gameinfo.UpdateAnnounce(core.Announce{})
	}
	return core.Announce{}
}
