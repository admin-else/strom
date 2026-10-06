package shapes

// BooleanOp mirrors net.minecraft.world.phys.shapes.BooleanOp. Java models it as
// a functional interface with named constants; Go uses a function type plus
// package-level values.
type BooleanOp func(first bool, second bool) bool

// BooleanOp constants mirror the static fields on BooleanOp.
var (
	BooleanOpFALSE       BooleanOp = func(first bool, second bool) bool { return false }
	BooleanOpNOT_OR      BooleanOp = func(first bool, second bool) bool { return !first && !second }
	BooleanOpONLY_SECOND BooleanOp = func(first bool, second bool) bool { return second && !first }
	BooleanOpNOT_FIRST   BooleanOp = func(first bool, second bool) bool { return !first }
	BooleanOpONLY_FIRST  BooleanOp = func(first bool, second bool) bool { return first && !second }
	BooleanOpNOT_SECOND  BooleanOp = func(first bool, second bool) bool { return !second }
	BooleanOpNOT_SAME    BooleanOp = func(first bool, second bool) bool { return first != second }
	BooleanOpNOT_AND     BooleanOp = func(first bool, second bool) bool { return !first || !second }
	BooleanOpAND         BooleanOp = func(first bool, second bool) bool { return first && second }
	BooleanOpSAME        BooleanOp = func(first bool, second bool) bool { return first == second }
	BooleanOpSECOND      BooleanOp = func(first bool, second bool) bool { return second }
	BooleanOpCAUSES      BooleanOp = func(first bool, second bool) bool { return !first || second }
	BooleanOpFIRST       BooleanOp = func(first bool, second bool) bool { return first }
	BooleanOpCAUSED_BY   BooleanOp = func(first bool, second bool) bool { return first || !second }
	BooleanOpOR          BooleanOp = func(first bool, second bool) bool { return first || second }
	BooleanOpTRUE        BooleanOp = func(first bool, second bool) bool { return true }
)
