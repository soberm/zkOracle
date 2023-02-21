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
        Pairing.G1Point[8] IC;
    }

    struct Proof {
        Pairing.G1Point A;
        Pairing.G2Point B;
        Pairing.G1Point C;
    }

    function verifyingKey() internal pure returns (VerifyingKey memory vk) {
        vk.alfa1 = Pairing.G1Point(
            uint256(
                17016714854791908693644018086513653893215896208793597930656032031220211116663
            ),
            uint256(
                16757343191852003659912740290377003884466900501978478227413411532851887445810
            )
        );
        vk.beta2 = Pairing.G2Point(
            [
                uint256(
                    19212565775822376068504934871948603337307730859635958737440729495734031840036
                ),
                uint256(
                    1896543391238166986438220987277088144782812196277386671733168178820554382208
                )
            ],
            [
                uint256(
                    11947584091058415250859895731054729020711392610138929881882642037752792498713
                ),
                uint256(
                    11465144035447976761266281201698220721789501065660531481107526648700121226973
                )
            ]
        );
        vk.gamma2 = Pairing.G2Point(
            [
                uint256(
                    4009109974783295918805509956779967312567502219063850676499853676350966990220
                ),
                uint256(
                    7846204514765366285793409985043567303970307567628304670364829816297303805624
                )
            ],
            [
                uint256(
                    5054098396770660681876101541804786938813126118873529878396533988757685118249
                ),
                uint256(
                    2741876006875187891294838881824750459162564868948072786934058270297151775873
                )
            ]
        );
        vk.delta2 = Pairing.G2Point(
            [
                uint256(
                    12176980155069594746593494755512158120272217240194221533475601641381708018159
                ),
                uint256(
                    156972775932620918865638601238692672995338757556696868440856509861709189689
                )
            ],
            [
                uint256(
                    14916664236922439659724225137545667002288490327586423421123155950050267223865
                ),
                uint256(
                    7553054984308731649975066584478974199679777772737742814688188847266928890762
                )
            ]
        );
        vk.IC[0] = Pairing.G1Point(
            uint256(
                20117270796444870905859799217928779780148304684007526500872031227512575247512
            ),
            uint256(
                11486310064601523660107540032184823943481762026756112008610915179111506576680
            )
        );
        vk.IC[1] = Pairing.G1Point(
            uint256(
                2004763464392332167695673861791617341535184804925847109540271258880806635085
            ),
            uint256(
                8735821397930230487551240115453381667488493578850275527447868110414486189062
            )
        );
        vk.IC[2] = Pairing.G1Point(
            uint256(
                15422045512518054203532337178660398382842307480739473538599953860227071172034
            ),
            uint256(
                1217484772316514325632900185568153285673952194885438117368306231163963434182
            )
        );
        vk.IC[3] = Pairing.G1Point(
            uint256(
                6073372061216848115575814061815666802783463725978918906095205631652351424169
            ),
            uint256(
                5667930076566061051790650196965573499796177082104767418156192749360730565894
            )
        );
        vk.IC[4] = Pairing.G1Point(
            uint256(
                1616396465362148730240970413052748417892877791190797168136724768996759761498
            ),
            uint256(
                20341963073843661771421480081384253510749547577155245758194705855109018761101
            )
        );
        vk.IC[5] = Pairing.G1Point(
            uint256(
                12715124451256412298221193856425961831136227109988896903328795243885911037400
            ),
            uint256(
                10245070776793182161833956198540237235764031257745609331512562924248557357349
            )
        );
        vk.IC[6] = Pairing.G1Point(
            uint256(
                18223583606774593781487397707051117317662528813470043525537177319517949960959
            ),
            uint256(
                20306376578452095538897294947723936531995248854610056150519309711626985213028
            )
        );
        vk.IC[7] = Pairing.G1Point(
            uint256(
                12070056175836988588397626193992373827775695409467964982336922546723283603280
            ),
            uint256(
                4091417531773487206766422023583867878972285363894802200382920537615830678118
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
        uint256[7] memory input
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
