package server

type TradeSession struct {
	Active, Field4 uint32
	Units          [2]*Object
	Kind           uint32
	Stock          *TradeItem
	Accepted       [2]uint32
	Offers         [2]*TradeItem
	Totals         [2]uint32
	Gold           [2]*Object
	Next, Prev     *TradeSession
}
type TradeItem struct {
	Object     *Object
	Value      uint32
	Next, Prev *TradeItem
}
