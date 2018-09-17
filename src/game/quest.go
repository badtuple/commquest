package game

import "../quest"

func maybeStartOrProgressQuest() error {
	_ = quest.Gen()
	return nil
}
