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
        Pairing.G1Point[9] IC;
    }

    struct Proof {
        Pairing.G1Point A;
        Pairing.G2Point B;
        Pairing.G1Point C;
    }

    function verifyingKey() internal pure returns (VerifyingKey memory vk) {
        vk.alfa1 = Pairing.G1Point(
            uint256(
                8419839428627061714140357861691566134881474360230375835637626242249113971491
            ),
            uint256(
                5369914303335915293021960132391787429424226814756283866572207785517650417432
            )
        );
        vk.beta2 = Pairing.G2Point(
            [
                uint256(
                    19009913377917025483075954421342038277892326068865437758106408645407972082517
                ),
                uint256(
                    5712350926399007535668350109039185470024873401213525826801019840750152054645
                )
            ],
            [
                uint256(
                    7774524998890089253555488299111391048915807172321422145460480441911070521117
                ),
                uint256(
                    1898291857371482558238184686353179568425475491296461533864704681366594549067
                )
            ]
        );
        vk.gamma2 = Pairing.G2Point(
            [
                uint256(
                    4179097497429890041920475304744571890541180826964375186886726784014376361536
                ),
                uint256(
                    7479032434544106234168576505359503321484192138838198413166317621325954692839
                )
            ],
            [
                uint256(
                    10883021926955399778795150401203975520258392779066679955323926038959866643883
                ),
                uint256(
                    13689532342786134991354932498329617468384805847833098939640120073587625868126
                )
            ]
        );
        vk.delta2 = Pairing.G2Point(
            [
                uint256(
                    749439120848622610901381862586187011597302624964960181700278117242856654820
                ),
                uint256(
                    17919315660394259279180372023964845597956568685058543619600434820068235157703
                )
            ],
            [
                uint256(
                    5622534396772574868042950304588272201989112327356114599323265312112617538422
                ),
                uint256(
                    18365638076314249995861789916384952574839189558057501144243558039384435374987
                )
            ]
        );
        vk.IC[0] = Pairing.G1Point(
            uint256(
                3870771352078397066634270886853142633068611602898612803095818344335929408789
            ),
            uint256(
                4576997275453297631827037839411966570219246449165130669863578248136795112493
            )
        );
        vk.IC[1] = Pairing.G1Point(
            uint256(
                3805564428475478237133894108547788880848553291403118431811049318113795623775
            ),
            uint256(
                7257840960664096976215178814996502724012421923841047758568755983539506338979
            )
        );
        vk.IC[2] = Pairing.G1Point(
            uint256(
                13328113424530599910638004685381924874765968492599608184650029804860261613052
            ),
            uint256(
                6958927142627112682262726799069605656721727443021579766291042632207865115378
            )
        );
        vk.IC[3] = Pairing.G1Point(
            uint256(
                6597509748769820417450692435737302189056149143696214790589011365899322497472
            ),
            uint256(
                13296534437799163819182629280476741964792285036764721119425509507370636933071
            )
        );
        vk.IC[4] = Pairing.G1Point(
            uint256(
                4822756062619671938445259908205022495178861331732994056099817712898720148681
            ),
            uint256(
                11855011450028705935689849535934066257217705629910953505364855681857988782121
            )
        );
        vk.IC[5] = Pairing.G1Point(
            uint256(
                15508958813097495360175904980287907623727166313461656102042091125756710501663
            ),
            uint256(
                11108706040118200724866375058691383839890817062668415603558756010971368708291
            )
        );
        vk.IC[6] = Pairing.G1Point(
            uint256(
                16575093534833459127102807480619539024484577922522466211814900514613488004257
            ),
            uint256(
                20915317618085717981414552453704308984754349301134284490084201886707616436673
            )
        );
        vk.IC[7] = Pairing.G1Point(
            uint256(
                13981089170785756433341278073703141923658201398954199091794608230782308576089
            ),
            uint256(
                19001908840158667794901458455835857484529633998723379435423119644347815220876
            )
        );
        vk.IC[8] = Pairing.G1Point(
            uint256(
                267329300174150601828595109590060049519814390843224674342780054476238882895
            ),
            uint256(
                11663309273196450044422777236694173992137102283583658544088983389217198361911
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
        uint256[8] memory input
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
