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

import "./Pairing.sol";

contract AggregationVerifier {
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
                3905127527146025156765272416675927748697235547470357549998877821255107079302
            ),
            uint256(
                7217595824121294988391587829185308340818299956584565100281524467356795278770
            )
        );
        vk.beta2 = Pairing.G2Point(
            [
                uint256(
                    2050016030674581335349132539332138721337640568524380126468380965558999653970
                ),
                uint256(
                    6004838838877283062406786285615997169152324447514466623528169280420672612318
                )
            ],
            [
                uint256(
                    2846873459068736356096090262310257018829886523990005324718470166477142650379
                ),
                uint256(
                    10374790157743050354366555960474821559047524238182619629715154422033807121806
                )
            ]
        );
        vk.gamma2 = Pairing.G2Point(
            [
                uint256(
                    1401188429620841212971641960021147658638150072070343140148826932558661596207
                ),
                uint256(
                    19316725531656418851542153946160692253521195731279715908921537760642900456707
                )
            ],
            [
                uint256(
                    4148060427500415543254352327560973405608914164143415024591846526775374408517
                ),
                uint256(
                    169923194330591178388216467113352639990078698871679560227378460409611133953
                )
            ]
        );
        vk.delta2 = Pairing.G2Point(
            [
                uint256(
                    18213643027016509131751747433264959991043430258025573096635554961543791375540
                ),
                uint256(
                    20793963136790779207468229444230213196245873125220681985164221048980928126171
                )
            ],
            [
                uint256(
                    17153996667925124962767210485654632109838273753893346500600409710822174779953
                ),
                uint256(
                    16782821129859967859485870446650725915594803693926605278312005796952169697356
                )
            ]
        );
        vk.IC[0] = Pairing.G1Point(
            uint256(
                3171109462461123210855075007028447311251477384181886722434920576431961609150
            ),
            uint256(
                9616722045570555374922377613269001435657174188815548381684959617328551019110
            )
        );
        vk.IC[1] = Pairing.G1Point(
            uint256(
                3734225733189277493551943181920197736392968778449163086942597874631799261087
            ),
            uint256(
                7652857868257055216136875352931411255484225283787685255536706140449267851354
            )
        );
        vk.IC[2] = Pairing.G1Point(
            uint256(
                15305427723796387967328048455854273161268569423235139559088304680186758761551
            ),
            uint256(
                21210768702428996375430283145346108723518041517275671293304248553833695167015
            )
        );
        vk.IC[3] = Pairing.G1Point(
            uint256(
                19769272315112367728468799910279212122075467741300562777550125102466298869892
            ),
            uint256(
                712487361178050084517224400667424285051397514918123989971995244134915757359
            )
        );
        vk.IC[4] = Pairing.G1Point(
            uint256(
                21481362032799539065015814858460155640007876695561451088342899352361615356745
            ),
            uint256(
                7797754787476042956426538288419858373345271992681023614400609517021840224340
            )
        );
        vk.IC[5] = Pairing.G1Point(
            uint256(
                16394152555968793924643719186726101410814161167234962874582220358592158221363
            ),
            uint256(
                21099119347219066606368147215151410454497136154226403493286835250817771377136
            )
        );
        vk.IC[6] = Pairing.G1Point(
            uint256(
                1932368365570905238304952296575551306838658068277335401492654290011238415618
            ),
            uint256(
                12405471140852495926442440371036245699328340024548117639803290452070451315823
            )
        );
        vk.IC[7] = Pairing.G1Point(
            uint256(
                1080203897953333536971688242966798978248687482436236631171625104232572003910
            ),
            uint256(
                17505889457021925526364380088208424417606771072468569004262106670622875180338
            )
        );
        vk.IC[8] = Pairing.G1Point(
            uint256(
                2567518869527138702894787242743372643845743685759043349924022790407333542513
            ),
            uint256(
                17203132490273687355946211797177030582893766830797832113665450376928143671628
            )
        );
        vk.IC[9] = Pairing.G1Point(
            uint256(
                14592814282201220132504820064273492070704403512054852698517753901432295608930
            ),
            uint256(
                13831068820837033881108956454125234323682612220840964996137662838182188648866
            )
        );
        vk.IC[10] = Pairing.G1Point(
            uint256(
                4684685106721893087144691911302918950816558962957859411406563256976376493475
            ),
            uint256(
                10089799999323624477360699608000659349711575937231269824362447417240776261658
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
