package zkOracle

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/algebra/twistededwards"
	"github.com/consensys/gnark/std/signature/eddsa"
)

func checkOwnership(api frontend.API, curve twistededwards.Curve, sk frontend.Variable, pk eddsa.PublicKey) {
	base := curve.Params().Base
	g := twistededwards.Point{X: base[0], Y: base[1]}
	pubKey := curve.ScalarMul(g, sk)

	curve.AssertIsOnCurve(pubKey)
	api.AssertIsEqual(pubKey.X, pk.A.X)
	api.AssertIsEqual(pubKey.Y, pk.A.Y)
}
