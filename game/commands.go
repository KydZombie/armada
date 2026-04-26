package game

import (
	"fmt"
	"sort"
	"strconv"
	"unicode"
	"strings"

	"github.com/KydZombie/armada/core"
)

func initializeCommands() *core.CommandDB[Game] {
	db := core.NewCommandDB[Game]()

	registerGenericCommands(db)
	registerInfoCommands(db)
	registerUnitCommands(db)
	registerCombatCommands(db)
	registerDebugCommands(db)

	return db
}

func parseRoomIndex(arg string, roomCount int) (int, string, bool) {
	roomRunes := []rune(strings.ToLower(arg))
	if len(roomRunes) != 1 {
		return 0, "Invalid room", false
	}

	roomIdx := int(roomRunes[0] - 'a')
	if roomIdx < 0 || roomIdx >= roomCount {
		return 0, "Invalid room", false
	}

	return roomIdx, "", true
}

func parseWeaponTarget(arg string, roomCount int) (roomIdx int, targetHull bool, errText string, ok bool) {
	if strings.EqualFold(arg, "hull") {
		return 0, true, "", true
	}

	return 0, false, "Only hull targeting is available until enemy train rooms are implemented.", false
}

func parsePositiveInt(arg string, label string) (int, string, bool) {
	value, err := strconv.Atoi(arg)
	if err != nil || value <= 0 {
		return 0, fmt.Sprintf("Invalid %s", label), false
	}

	return value, "", true
}

func registerAliasCommand(db *core.CommandDB[Game], name string, onRun func(args []string, game *Game) (string, bool), description string) {
	db.RegisterCommand(core.Command[Game]{
		Name:        name,
		OnRun:       onRun,
		Description: []string{description},
	})
}

func commandCategory(name string) string {
	switch name {
	case "help", "clear":
		return "Utility"
	case "selc", "move":
		return "Crew"
	case "attackenemy", "attack", "atk", "enemystatus", "estatus", "status", "selw", "weapons", "fire", "trainstatus", "tstatus", "damagetrain", "dtrain", "dt", "damageroom", "droom", "dr", "damage":
		return "Combat"
	default:
		return "Other"
	}
}

func weaponListText(train *Train) string {
	if len(train.Weapons) == 0 {
		return "none"
	}

	parts := make([]string, 0, len(train.Weapons))
	for _, weapon := range train.Weapons {
		parts = append(parts, weapon.StatusText())
	}

	return strings.Join(parts, ", ")
}

func setCombatStatusMessage(game *Game, message string) {
	game.SetCombatStatus(message)
}

func combatFailure(game *Game, message string) (string, bool) {
	setCombatStatusMessage(game, message)
	return message, false
}

func parseDamageAmount(arg string) (int, string, bool) {
	return parsePositiveInt(arg, "damage amount")
}

func enemyStatusText(game *Game) (string, bool) {
	if game.Enemy == nil {
		return "There is no enemy to inspect.", false
	}

	statusText := "alive"
	if !game.Enemy.Alive() {
		statusText = "dead"
	}

	return fmt.Sprintf(
		"Enemy: %s. HP: %d/%d. Shields: %d. Status: %s.",
		game.Enemy.Name(),
		game.Enemy.Health(),
		game.Enemy.MaxHealth(),
		game.Enemy.ShieldLayers(),
		statusText,
	), true
}

func trainStatusText(game *Game) (string, bool) {
	lifeSupportText := "online"
	if !game.Train.LifeSupportOperational() {
		lifeSupportText = fmt.Sprintf("offline (%d dmg/tick)", game.Train.LifeSupportDamagePerTick())
	}

	return fmt.Sprintf(
		"Train hull: %d/%d. Ready damage: %d. Shields: %d. Evasion: %d%%. Medbay: %d HP/tick. Life Support: %s. Weapons: %s.",
		game.Train.Health,
		game.Train.MaxHealth,
		game.Train.TotalAttackPower(),
		game.Train.ShieldLayers(),
		game.Train.EvasionChance(),
		game.Train.MedbayHealingPerTick(),
		lifeSupportText,
		weaponListText(game.Train),
	), true
}

func weaponsStatusText(game *Game) (string, bool) {
	return fmt.Sprintf(
		"Weapons: %s. Weapons system: %t.",
		weaponListText(game.Train),
		game.Train.WeaponsOperational(),
	), true
}

func applyTrainDamageCommand(game *Game, damageArg string) (string, bool) {
	damage, errText, ok := parseDamageAmount(damageArg)
	if !ok {
		return combatFailure(game, errText)
	}

	game.Train.ApplyHullDamage(damage)
	message := fmt.Sprintf(
		"Train takes %d damage. Hull is now %d/%d.",
		damage,
		game.Train.Health,
		game.Train.MaxHealth,
	)
	setCombatStatusMessage(game, message)

	return message, true
}

func applyRoomDamageCommand(game *Game, roomArg string, damageArg string) (string, bool) {
	roomIdx, errText, ok := parseRoomIndex(roomArg, len(game.Train.Rooms))
	if !ok {
		return combatFailure(game, errText)
	}

	damage, errText, ok := parseDamageAmount(damageArg)
	if !ok {
		return combatFailure(game, errText)
	}

	room := &game.Train.Rooms[roomIdx]
	room.ApplyDamage(damage)
	message := fmt.Sprintf(
		"Room %c takes %d damage. Room health is now %d/%d.",
		room.GetRune(),
		damage,
		room.Health,
		room.MaxHealth,
	)
	setCombatStatusMessage(game, message)

	return message, true
}

func resolveEnemyWeaponHit(enemy Enemy, weapon *Weapon) (damageApplied int, shieldsAbsorbed bool) {
	if enemy == nil || weapon == nil {
		return 0, false
	}

	return enemy.ResolveWeaponHit(*weapon)
}

func resolvePlayerWeaponAttack(game *Game, weapon *Weapon, target string) (string, bool) {
	if game.Enemy == nil {
		return combatFailure(game, "There is no enemy to attack.")
	}

	if !game.Enemy.Alive() {
		return combatFailure(game, fmt.Sprintf("%s is already defeated.", game.Enemy.Name()))
	}

	if !game.Train.WeaponsOperational() {
		return combatFailure(game, "Weapons system is offline.")
	}

	if !weapon.Ready() {
		return combatFailure(game, fmt.Sprintf("%s is cooling down for %d more second(s).", weapon.Name, weapon.CooldownDisplaySeconds()))
	}

	_, targetHull, errText, ok := parseWeaponTarget(target, len(game.Train.Rooms))
	if !ok {
		return combatFailure(game, errText)
	}

	var resultLines []string
	var combatStatusLines []string
	if targetHull {
		damageApplied, shieldsAbsorbed := resolveEnemyWeaponHit(game.Enemy, weapon)
		playerHitText := fmt.Sprintf("Your %s hits the %s hull for %d damage.", weapon.Name, game.Enemy.Name(), damageApplied)
		if shieldsAbsorbed {
			playerHitText = fmt.Sprintf("Your %s hits the %s shields and deals no hull damage.", weapon.Name, game.Enemy.Name())
		}
		resultLines = append(resultLines, playerHitText)
		combatStatusLines = append(combatStatusLines, playerHitText)
	}

	weapon.StartCooldown()

	if game.Enemy.Alive() {
		enemyAttack := game.Enemy.Attack()
		hullDamage, evaded, shielded := game.Train.ResolveIncomingAttack(enemyAttack)

		switch {
		case evaded:
			enemyResponseText := fmt.Sprintf("The %s attacks back, but the train evades the hit.", game.Enemy.Name())
			resultLines = append(resultLines, enemyResponseText)
			combatStatusLines = append(combatStatusLines, enemyResponseText)
		case shielded:
			enemyResponseText := fmt.Sprintf("The %s attacks for %d, but your shields absorb it.", game.Enemy.Name(), enemyAttack)
			resultLines = append(resultLines, enemyResponseText)
			combatStatusLines = append(combatStatusLines, enemyResponseText)
		default:
			enemyResponseText := fmt.Sprintf("The %s attacks back for %d. Train hull is now %d/%d.", game.Enemy.Name(), hullDamage, game.Train.Health, game.Train.MaxHealth)
			resultLines = append(resultLines, enemyResponseText)
			combatStatusLines = append(combatStatusLines, enemyResponseText)
		}

		game.SetCombatStatus(combatStatusLines...)
		return strings.Join(resultLines, " "), true
	}

	enemyDefeatedText := fmt.Sprintf("%s is defeated.", game.Enemy.Name())
	resultLines = append(resultLines, enemyDefeatedText)
	combatStatusLines = append(combatStatusLines, enemyDefeatedText)
	game.SetCombatStatus(combatStatusLines...)
	return strings.Join(resultLines, " "), true
}

func registerGenericCommands(db *core.CommandDB[Game]) {
	// TODO: Need to be able to do more than just this
	db.RegisterCommand(core.Command[Game]{
		Name: "help",
		OnRun: func(args []string, game *Game) (string, bool) {
			categoryNames := make(map[string][]string)
			for name, command := range db.CmdMap {
				if command.Description != nil {
					category := commandCategory(name)
					categoryNames[category] = append(categoryNames[category], name)
				}
			}

			categories := make([]string, 0, len(categoryNames))
			totalCommands := 0
			for category, names := range categoryNames {
				sort.Strings(names)
				categoryNames[category] = names
				categories = append(categories, category)
				totalCommands += len(names)
			}

			sort.Strings(categories)
			game.Terminal.OutputText(fmt.Sprintf("Available commands (%d):", totalCommands), StandardOutputMessage)
			for _, category := range categories {
				game.Terminal.OutputText(fmt.Sprintf("[%s]", category), StandardOutputMessage)
				for _, name := range categoryNames[category] {
					command := db.CmdMap[name]
					game.Terminal.OutputText(fmt.Sprint(command.Name, " ^"), StandardOutputMessage)
					for _, line := range command.Description {
						game.Terminal.OutputText(fmt.Sprint("  ", line), StandardOutputMessage)
					}
				}
			}

			return "", true
		},
		Description: []string{"Help command!"},
	})

	db.RegisterCommand(core.Command[Game]{
		Name: "clear",
		OnRun: func(args []string, game *Game) (string, bool) {
			if len(args) != 0 {
				return "clear does not take any arguments", false
			}

			game.Terminal.history = nil
			game.Terminal.currentHistoryLine = 0

			return "", true
		},
		Description: []string{"Clear the terminal history."},
	})
}

func registerInfoCommands(db *core.CommandDB[Game]) {

}

func registerUnitCommands(db *core.CommandDB[Game]) {
	db.RegisterCommand(core.Command[Game]{
		Name: "selc",
		OnRun: func(args []string, game *Game) (string, bool) {
			if len(args) != 1 {
				return "Invalid number of arguments", false
			}

			characterIdx, err := strconv.Atoi(args[0])
			characterIdx -= 1
			if err != nil || characterIdx < 0 || characterIdx >= len(game.Train.Characters) {
				return fmt.Sprintf("Invalid character index (%s).", args[0]), false
			}
			game.SelectedCharacterIndex = characterIdx
			character := game.Train.Characters[characterIdx]

			return fmt.Sprint("Selected ", character.Name, "."), true
		},
		Description: []string{"Select a character"},
	})

	db.RegisterCommand(core.Command[Game]{
		Name: "move",
		OnRun: func(args []string, game *Game) (string, bool) {
			if len(args) != 3 {
				return "Invalid number of arguments", false
			}

			if game.SelectedCharacterIndex < 0 {
				return "No character selected", false
			}

			roomIdx, errText, ok := parseRoomIndex(args[0], len(game.Train.Rooms))
			if !ok {
				return errText, false
			}

			room := game.Train.Rooms[roomIdx]

			x, err := strconv.Atoi(args[1])
			x -= 1
			if err != nil || x < 0 || x >= room.Width {
				return fmt.Sprintf("Invalid x (%s). X starts at 1.", args[1]), false
			}
			y, err := strconv.Atoi(args[2])
			y -= 1
			if err != nil || y < 0 || y >= room.Height {
				return fmt.Sprintf("Invalid y (%s). Y starts at 1.", args[2]), false
			}

			character := game.Train.Characters[game.SelectedCharacterIndex]
			path, err := game.Train.MoveCharacter(character, RoomPos{RoomId: roomIdx, X: x, Y: y})
			if err != nil {
				return err.Error(), false
			}

			game.SelectedCharacterIndex = -1

			pathLabels := make([]string, 0, len(path))
			for _, roomId := range path {
				pathLabels = append(pathLabels, string(game.Train.Rooms[roomId].GetRune()))
			}

			// Debug: show animation path waypoints
			debugPath := make([]string, 0)
			for i, pos := range character.MovementPath {
				debugPath = append(debugPath, fmt.Sprintf("%s(%d,%d)", string(game.Train.Rooms[pos.RoomId].GetRune()), pos.X+1, pos.Y+1))
				if i < len(character.MovementPath)-1 {
					debugPath = append(debugPath, "->")
				}
			}

			if len(pathLabels) > 1 {
				return fmt.Sprintf("Moved %s through %s. Path: %s", character.Name, strings.Join(pathLabels, " -> "), strings.Join(debugPath, "")), true
			}

			return fmt.Sprint("Moved ", character.Name, " within room ", pathLabels[0], ". Path: ", strings.Join(debugPath, ""), "."), true
		},
		Description: []string{"Move a character", "Takes arguments [room], [x], and [y] (x and y start at 1)"},
	})
}

func registerCombatCommands(db *core.CommandDB[Game]) {
	attackEnemyCommand := core.Command[Game]{
		Name: "attackenemy",
		OnRun: func(args []string, game *Game) (string, bool) {
			if game.SelectedWeaponIndex >= len(game.Train.Weapons) {
				return "Invalid weapon selected", false
			}
			//weapon, ok := game.Train.GetWeaponByName("cannon")
			//if !ok {
			//	return "No cannon is installed.", false
			//}

			weapon := &game.Train.Weapons[game.SelectedWeaponIndex]

			return resolvePlayerWeaponAttack(game, weapon, "hull")
		},
		Description: []string{"Fire the cannon at the enemy hull."},
	}
	db.RegisterCommand(attackEnemyCommand)
	registerAliasCommand(db, "attack", attackEnemyCommand.OnRun, "Alias for attackenemy.")
	registerAliasCommand(db, "atk", attackEnemyCommand.OnRun, "Short alias for attackenemy.")

	enemyStatusCommand := core.Command[Game]{
		Name:        "enemystatus",
		OnRun:       func(args []string, game *Game) (string, bool) { return enemyStatusText(game) },
		Description: []string{"Get the status of the enemy."},
	}
	db.RegisterCommand(enemyStatusCommand)
	registerAliasCommand(db, "estatus", enemyStatusCommand.OnRun, "Short alias for enemystatus.")
	db.RegisterCommand(core.Command[Game]{
		Name:        "weapons",
		OnRun:       func(args []string, game *Game) (string, bool) { return weaponsStatusText(game) },
		Description: []string{"List installed weapons and cooldowns."},
	})
	db.RegisterCommand(core.Command[Game]{
		Name: "fire",
		OnRun: func(args []string, game *Game) (string, bool) {
			if len(args) != 2 {
				return "fire takes arguments [weapon] and [target]", false
			}

			weapon, ok := game.Train.GetWeaponByName(strings.ToLower(args[0]))
			if !ok {
				return fmt.Sprintf("Unknown weapon: %s.", args[0]), false
			}

			return resolvePlayerWeaponAttack(game, weapon, args[1])
		},
		Description: []string{"Fire a weapon at the enemy hull.", "Examples: fire cannon hull, fire missile hull"},
	})
	db.RegisterCommand(core.Command[Game]{
		Name:        "trainstatus",
		OnRun:       func(args []string, game *Game) (string, bool) { return trainStatusText(game) },
		Description: []string{"Get the status of the train hull and total damage."},
	})
	registerAliasCommand(db, "tstatus", db.CmdMap["trainstatus"].OnRun, "Short alias for trainstatus.")
	db.RegisterCommand(core.Command[Game]{
		Name: "status",
		OnRun: func(args []string, game *Game) (string, bool) {
			if len(args) == 0 {
				if game.Enemy != nil {
					return enemyStatusText(game)
				}
				return trainStatusText(game)
			}

			switch strings.ToLower(args[0]) {
			case "enemy":
				return enemyStatusText(game)
			case "train":
				return trainStatusText(game)
			case "weapons", "weapon":
				return weaponsStatusText(game)
			default:
				return "status takes one of: enemy, train, weapons", false
			}
		},
		Description: []string{"Get combat/system status.", "Examples: status enemy, status train, status weapons"},
	})

	db.RegisterCommand(core.Command[Game]{
		Name: "damagetrain",
		OnRun: func(args []string, game *Game) (string, bool) {
			if len(args) != 1 {
				return "Invalid number of arguments", false
			}

			return applyTrainDamageCommand(game, args[0])
		},
		Description: []string{"Damage the train hull.", "Takes arguments [amount]"},
	})
	registerAliasCommand(db, "dtrain", db.CmdMap["damagetrain"].OnRun, "Short alias for damagetrain.")

	db.RegisterCommand(core.Command[Game]{
		Name: "damageroom",
		OnRun: func(args []string, game *Game) (string, bool) {
			if len(args) != 2 {
				return "Invalid number of arguments", false
			}

			return applyRoomDamageCommand(game, args[0], args[1])
		},
		Description: []string{"Damage a train room.", "Takes arguments [room] and [amount]"},
	})
	registerAliasCommand(db, "droom", db.CmdMap["damageroom"].OnRun, "Short alias for damageroom.")
	registerAliasCommand(db, "dr", db.CmdMap["damageroom"].OnRun, "Very short alias for damageroom.")
	registerAliasCommand(db, "dt", db.CmdMap["damagetrain"].OnRun, "Very short alias for damagetrain.")
	db.RegisterCommand(core.Command[Game]{
		Name: "damage",
		OnRun: func(args []string, game *Game) (string, bool) {
			if len(args) == 0 {
				return "damage takes a target: room or train", false
			}

			switch strings.ToLower(args[0]) {
			case "train":
				if len(args) != 2 {
					return "damage train takes [amount]", false
				}
				return applyTrainDamageCommand(game, args[1])
			case "room":
				if len(args) != 3 {
					return "damage room takes [room] and [amount]", false
				}
				return applyRoomDamageCommand(game, args[1], args[2])
			default:
				return "damage takes one of: room, train", false
			}
		},
		Description: []string{"Damage train systems for testing.", "Examples: damage room a 2, damage train 5"},
	})

	db.RegisterCommand(core.Command[Game]{
		Name: "selw",
		OnRun: func(args []string, game *Game) (string, bool) {
			if len(args) == 0 {
				return "Must give weapon name", false
			}

			weaponName := args[0]

			for i, weapon := range game.Train.Weapons {
				if strings.Compare(weaponName, weapon.Name) == 0 {
					game.SelectedWeaponIndex = i
					return fmt.Sprint("Selected weapon ", weaponName), true
				}
			}

			return fmt.Sprint("Couldn't find a weapon named ", weaponName), false
		},
		Description: []string{"Select weapon using its name."},
	})
}

func registerDebugCommands(db *core.CommandDB[Game]) {
}
