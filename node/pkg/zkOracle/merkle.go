package zkOracle

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash/mimc"
)

func leafSum(api frontend.API, h mimc.MiMC, data frontend.Variable) frontend.Variable {

	h.Write(data)
	res := h.Sum()

	return res
}

func nodeSum(api frontend.API, h mimc.MiMC, a, b frontend.Variable) frontend.Variable {

	h.Write(a, b)
	res := h.Sum()

	return res
}

func ComputeRootFromPath(api frontend.API, h mimc.MiMC, proofSet, helper []frontend.Variable) frontend.Variable {

	sum := leafSum(api, h, proofSet[0])

	for i := 1; i < len(proofSet); i++ {
		api.AssertIsBoolean(helper[i-1])
		d1 := api.Select(helper[i-1], sum, proofSet[i])
		d2 := api.Select(helper[i-1], proofSet[i], sum)
		sum = nodeSum(api, h, d1, d2)
	}

	return sum
}
