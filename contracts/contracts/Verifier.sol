// SPDX-License-Identifier: AML
//
// Copyright 2017 Christian Reitwiessner
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.

// 2019 OKIMS

pragma solidity ^0.8.0;

library Pairing {
    uint256 constant PRIME_Q =
        21888242871839275222246405745257275088696311157297823662689037894645226208583;

    struct G1Point {
        uint256 X;
        uint256 Y;
    }

    // Encoding of field elements is: X[0] * z + X[1]
    struct G2Point {
        uint256[2] X;
        uint256[2] Y;
    }

    /*
     * @return The negation of p, i.e. p.plus(p.negate()) should be zero.
     */
    function negate(G1Point memory p) internal pure returns (G1Point memory) {
        // The prime q in the base field F_q for G1
        if (p.X == 0 && p.Y == 0) {
            return G1Point(0, 0);
        } else {
            return G1Point(p.X, PRIME_Q - (p.Y % PRIME_Q));
        }
    }

    /*
     * @return The sum of two points of G1
     */
    function plus(
        G1Point memory p1,
        G1Point memory p2
    ) internal view returns (G1Point memory r) {
        uint256[4] memory input;
        input[0] = p1.X;
        input[1] = p1.Y;
        input[2] = p2.X;
        input[3] = p2.Y;
        bool success;

        // solium-disable-next-line security/no-inline-assembly
        assembly {
            success := staticcall(sub(gas(), 2000), 6, input, 0xc0, r, 0x60)
            // Use "invalid" to make gas estimation work
            switch success
            case 0 {
                invalid()
            }
        }

        require(success, "pairing-add-failed");
    }

    /*
     * @return The product of a point on G1 and a scalar, i.e.
     *         p == p.scalar_mul(1) and p.plus(p) == p.scalar_mul(2) for all
     *         points p.
     */
    function scalar_mul(
        G1Point memory p,
        uint256 s
    ) internal view returns (G1Point memory r) {
        uint256[3] memory input;
        input[0] = p.X;
        input[1] = p.Y;
        input[2] = s;
        bool success;
        // solium-disable-next-line security/no-inline-assembly
        assembly {
            success := staticcall(sub(gas(), 2000), 7, input, 0x80, r, 0x60)
            // Use "invalid" to make gas estimation work
            switch success
            case 0 {
                invalid()
            }
        }
        require(success, "pairing-mul-failed");
    }

    /* @return The result of computing the pairing check
     *         e(p1[0], p2[0]) *  .... * e(p1[n], p2[n]) == 1
     *         For example,
     *         pairing([P1(), P1().negate()], [P2(), P2()]) should return true.
     */
    function pairing(
        G1Point memory a1,
        G2Point memory a2,
        G1Point memory b1,
        G2Point memory b2,
        G1Point memory c1,
        G2Point memory c2,
        G1Point memory d1,
        G2Point memory d2
    ) internal view returns (bool) {
        G1Point[4] memory p1 = [a1, b1, c1, d1];
        G2Point[4] memory p2 = [a2, b2, c2, d2];
        uint256 inputSize = 24;
        uint256[] memory input = new uint256[](inputSize);

        for (uint256 i = 0; i < 4; i++) {
            uint256 j = i * 6;
            input[j + 0] = p1[i].X;
            input[j + 1] = p1[i].Y;
            input[j + 2] = p2[i].X[0];
            input[j + 3] = p2[i].X[1];
            input[j + 4] = p2[i].Y[0];
            input[j + 5] = p2[i].Y[1];
        }

        uint256[1] memory out;
        bool success;

        // solium-disable-next-line security/no-inline-assembly
        assembly {
            success := staticcall(
                sub(gas(), 2000),
                8,
                add(input, 0x20),
                mul(inputSize, 0x20),
                out,
                0x20
            )
            // Use "invalid" to make gas estimation work
            switch success
            case 0 {
                invalid()
            }
        }

        require(success, "pairing-opcode-failed");

        return out[0] != 0;
    }
}

contract Verifier {
    using Pairing for *;

    uint256 constant SNARK_SCALAR_FIELD =
        21888242871839275222246405745257275088548364400416034343698204186575808495617;
    uint256 constant PRIME_Q =
        21888242871839275222246405745257275088696311157297823662689037894645226208583;

    struct VerifyingKey {
        Pairing.G1Point alfa1;
        Pairing.G2Point beta2;
        Pairing.G2Point gamma2;
        Pairing.G2Point delta2;
        Pairing.G1Point[11] IC;
    }

    struct Proof {
        Pairing.G1Point A;
        Pairing.G2Point B;
        Pairing.G1Point C;
    }

    function verifyingKey() internal pure returns (VerifyingKey memory vk) {
        vk.alfa1 = Pairing.G1Point(
            uint256(
                20743957995240333580170798864030367388043445301818510921568676180194904049694
            ),
            uint256(
                20905534378722221910254884038399036677327390721917982598802318145698136027341
            )
        );
        vk.beta2 = Pairing.G2Point(
            [
                uint256(
                    15848774384288714248558924597291258342220171038964542301120687627351223796311
                ),
                uint256(
                    5784886905210371246523447256373913078373012957579428468705174199792699041079
                )
            ],
            [
                uint256(
                    12381408612144995627167075580662253011115321198991890155335702609962744996026
                ),
                uint256(
                    13311225309212570056634919318538794825364903984983381861857897674426918065581
                )
            ]
        );
        vk.gamma2 = Pairing.G2Point(
            [
                uint256(
                    17719664661464732015611755358633618881053948876815553034210847929169596004969
                ),
                uint256(
                    8895760285573977580007191114580725823505543974921781387826709614331208048964
                )
            ],
            [
                uint256(
                    15999600839487866123196814839181201808956345575325080697063390276569408669441
                ),
                uint256(
                    11081962220322473721733471389198242949572551573827427556553691104404974565544
                )
            ]
        );
        vk.delta2 = Pairing.G2Point(
            [
                uint256(
                    18412922374839090327871052017831345216151928127253473490543351223630920360797
                ),
                uint256(
                    13577049657873460566000258116995419313851842930480993780796907146133392889209
                )
            ],
            [
                uint256(
                    11284054175159328143503693586584498811221143178077110202719631876862013012769
                ),
                uint256(
                    18661162242073199685560883353387280976177701091768417147921186087558246414253
                )
            ]
        );
        vk.IC[0] = Pairing.G1Point(
            uint256(
                8720443674563639744804387150152213527973461565624785678244466218623307395018
            ),
            uint256(
                6892786361160456278555613052733201668635070577805730594774320738235928466577
            )
        );
        vk.IC[1] = Pairing.G1Point(
            uint256(
                761250913433311676583895998310430825462153549353843885364196636886115350128
            ),
            uint256(
                1846731293688344132634631212081204677993371587256494040768666889178508406991
            )
        );
        vk.IC[2] = Pairing.G1Point(
            uint256(
                9053827860213174771394336968855974601916436224762073389373253999523432661728
            ),
            uint256(
                15525135389983450928310698198900824230934055711558763479358933979152684409901
            )
        );
        vk.IC[3] = Pairing.G1Point(
            uint256(
                7580571741099592096686785442378111939855999474985400993891701942306676118288
            ),
            uint256(
                16970326573078629753897975871246717249471018819839023736901041341867009673280
            )
        );
        vk.IC[4] = Pairing.G1Point(
            uint256(
                19687055720547218893816467789490487484809832435257977795254996275733127178700
            ),
            uint256(
                17862758906887228423257018958090673216358081800209346357489870155542294000497
            )
        );
        vk.IC[5] = Pairing.G1Point(
            uint256(
                6298759144030209478390558784322148303644866942571171566856448818827293034524
            ),
            uint256(
                8956179589354080402814560577492180116365062534440642294422661365656024961456
            )
        );
        vk.IC[6] = Pairing.G1Point(
            uint256(
                15801379661377989616854543383747499345404829394492658688041675780887887403917
            ),
            uint256(
                10263869644605820627554265629219381589447674165366094807614288275962162516157
            )
        );
        vk.IC[7] = Pairing.G1Point(
            uint256(
                20122396595366427927258489466736410425256137121023276076151779244786085674311
            ),
            uint256(
                8261030299334640763383804335491417065230059820132997474630660394085363761776
            )
        );
        vk.IC[8] = Pairing.G1Point(
            uint256(
                11802180486660100719886362028338208789719113471463789570313380387917168189768
            ),
            uint256(
                10924598976299286672750339629100378088744713928227534721061166055870476334603
            )
        );
        vk.IC[9] = Pairing.G1Point(
            uint256(
                11042356441843154044533461089643023107964331899588258559058073678755166158967
            ),
            uint256(
                15316502506160907603480758384661657015560913651394616961772143589971071261248
            )
        );
        vk.IC[10] = Pairing.G1Point(
            uint256(
                12617767177189648980556055310090588367110560387248249026768252047202295319832
            ),
            uint256(
                18905588440476182388238895995393506600463372039720126360746067120034072924916
            )
        );
    }

    /*
     * @returns Whether the proof is valid given the hardcoded verifying key
     *          above and the public inputs
     */
    function verifyProof(
        uint256[2] memory a,
        uint256[2][2] memory b,
        uint256[2] memory c,
        uint256[10] memory input
    ) public view returns (bool r) {
        Proof memory proof;
        proof.A = Pairing.G1Point(a[0], a[1]);
        proof.B = Pairing.G2Point([b[0][0], b[0][1]], [b[1][0], b[1][1]]);
        proof.C = Pairing.G1Point(c[0], c[1]);

        VerifyingKey memory vk = verifyingKey();

        // Compute the linear combination vk_x
        Pairing.G1Point memory vk_x = Pairing.G1Point(0, 0);

        // Make sure that proof.A, B, and C are each less than the prime q
        require(proof.A.X < PRIME_Q, "verifier-aX-gte-prime-q");
        require(proof.A.Y < PRIME_Q, "verifier-aY-gte-prime-q");

        require(proof.B.X[0] < PRIME_Q, "verifier-bX0-gte-prime-q");
        require(proof.B.Y[0] < PRIME_Q, "verifier-bY0-gte-prime-q");

        require(proof.B.X[1] < PRIME_Q, "verifier-bX1-gte-prime-q");
        require(proof.B.Y[1] < PRIME_Q, "verifier-bY1-gte-prime-q");

        require(proof.C.X < PRIME_Q, "verifier-cX-gte-prime-q");
        require(proof.C.Y < PRIME_Q, "verifier-cY-gte-prime-q");

        // Make sure that every input is less than the snark scalar field
        for (uint256 i = 0; i < input.length; i++) {
            require(
                input[i] < SNARK_SCALAR_FIELD,
                "verifier-gte-snark-scalar-field"
            );
            vk_x = Pairing.plus(
                vk_x,
                Pairing.scalar_mul(vk.IC[i + 1], input[i])
            );
        }

        vk_x = Pairing.plus(vk_x, vk.IC[0]);

        return
            Pairing.pairing(
                Pairing.negate(proof.A),
                proof.B,
                vk.alfa1,
                vk.beta2,
                vk_x,
                vk.gamma2,
                proof.C,
                vk.delta2
            );
    }
}
