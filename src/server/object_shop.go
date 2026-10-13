package server

type GoldInitData struct {
	Amount uint32 // 0, 0
}

// ShopStockEntry is an item a shopkeeper sells (28 bytes on 386).
type ShopStockEntry struct {
	Type      uint32          // 0, 0
	Count     byte            // 1, 4
	_         [3]byte         // 1, 5
	Reward    uint32          // 2, 8: spell, ability or field guide type
	Modifiers [4]*ModifierEff // 3, 12
}

type ShopkeeperInitData struct {
	StockCount byte               // 0, 0
	_          [3]byte            // 0, 1
	Stock      [60]ShopStockEntry // 1, 4
	ShopText   [32]byte           // 421, 1684
	Field1716  float32            // 429, 1716
	Field1720  float32            // 430, 1720
}
