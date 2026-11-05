package items

func NewWeapon(id int, name string, level int, weight float64, quality Quality) *Weapon {
	return &Weapon{id: id, name: name, level: level, weight: weight, quality: quality}
}

func NewArmor(id int, name string, level int, weight float64, quality Quality) *Armor {
	return &Armor{id: id, name: name, level: level, weight: weight, quality: quality}
}

func NewPotion(id int, name string, weight float64, effect string, charges int) *Potion {
	return &Potion{id: id, name: name, weight: weight, Effect: effect, Charges: charges}
}
