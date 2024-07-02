package teleblog

// v1.
// Problem: you still don't know is it bold or italic

// type MarkupTypeLiteral string

// const (
// 	Bold   MarkupTypeLiteral = "bold"
// 	Italic MarkupTypeLiteral = "italic"
// )

// func MarkupTypeLiteralSwitch(
// 	mt MarkupTypeLiteral,
// 	boldFn func(mt MarkupTypeLiteral) error,
// 	italicFn func(mt MarkupTypeLiteral) error,
// ) error {
// 	switch mt {
// 	case Bold:
// 		return boldFn(mt)
// 	case Italic:
// 		return italicFn(mt)
// 	default:
// 		return errors.New("Unknown markup type")
// 	}
// }

// func ttt() {
// 	some := "bold"
// 	MarkupTypeLiteralSwitch(
// 		MarkupTypeLiteral(some),
// 		func(mt MarkupTypeLiteral) error {
// 			return nil
// 		},
// 		func(mt MarkupTypeLiteral) error {
// 			return nil
// 		},
// 	)
// }

// // v2
// // Problem: ...

// type MarkupTypeLiteral interface {
// 	IsMarkupTypeLiteral()
// }

// type MarkupTypeLiteralBold string

// func (MarkupTypeLiteralBold) IsMarkupTypeLiteral() {}

// type MarkupTypeLiteralItalic string

// func (MarkupTypeLiteralItalic) IsMarkupTypeLiteral() {}

// const (
// 	Bold   MarkupTypeLiteralBold   = "bold"
// 	Italic MarkupTypeLiteralItalic = "italic"
// )

// func MarkupTypeLiteralSwitch(
// 	mt MarkupTypeLiteral,
// 	boldFn func(mt MarkupTypeLiteralBold) error,
// 	italicFn func(mt MarkupTypeLiteralItalic) error,
// ) error {
// 	switch mt {
// 	case Bold:
// 		return boldFn(mt.(MarkupTypeLiteralBold))
// 	case Italic:
// 		return italicFn(mt.(MarkupTypeLiteralItalic))
// 	default:
// 		return errors.New("Unknown markup type")
// 	}
// }

// func ttt() {
// 	some := "bold"
// 	MarkupTypeLiteralSwitch(
// 		some,
// 		func(mt MarkupTypeLiteralBold) error {
// 			return nil
// 		},
// 		func(mt MarkupTypeLiteralItalic) error {
// 			return nil
// 		},
// 	)
// }

// type MarkupTypeLiteralBold string

// func (MarkupTypeLiteralBold) IsMarkupTypeLiteral() {}

// func (MarkupTypeLiteralBold) String() string {
// 	return "bold"
// }

// type MarkupTypeLiteralItalic string

// func (MarkupTypeLiteralItalic) IsMarkupTypeLiteral() {}

// func (MarkupTypeLiteralItalic) String() string {
// 	return "italic"
// }
