package main

var _ Item = (*Weapon)(nil)
var _ Item = (*Armor)(nil)
var _ Item = (*Potion)(nil)

var _ Storable = (*Weapon)(nil)
var _ Storable = (*Armor)(nil)
var _ Storable = (*Potion)(nil)
